package otp

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement/pkg/engagement/services/otp")

const (
	issuer       = "Savannah Informatics Limited"
	accountName  = "info@healthcloud.co.ke"
	subject      = "Be.Well Verification Code"
	whatsappStep = 1
	twilioStep   = 2
	otpMsg       = "%s is your Be.Well verification code %s"
)

// These constants are here to support Integration Testing
// done by the Frontend team.
// IT is shorthand for Integration Tests
const (
	ITCode              = "123456"
	ITPhoneNumber       = "0798000000"
	TestUserPhoneNumber = "+254798000000"
	ITEmail             = "test@healthcloud.co.ke"
)

// ServiceOTP is an interface that defines all interactions with OTP service
type ServiceOTP interface {
	GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error)
	SaveOTPToFirestore(otp dto.OTP) error
	VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error)
	VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error)
	GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error)
	EmailVerificationOtp(ctx context.Context, email *string) (string, error)
	GenerateOTP(ctx context.Context) (string, error)
}

// Service is an OTP generation and validation service
type Service struct {
	whatsapp whatsapp.ServiceWhatsapp
	mail     mail.ServiceMail
	sms      sms.ServiceSMS
	twilio   twilio.ServiceTwilio

	totpOpts             totp.GenerateOpts
	firestoreClient      *firestore.Client
	rootCollectionSuffix string
}

// NewService initializes a valid OTP service
// First we fetch the dependencies from dep.yaml file. Since this service has a predefined set
// of dependencies, the same dependecies defined in the yaml should be defined in the service
// struct definition explicitly, No guess work.
func NewService(
	whatsapp whatsapp.ServiceWhatsapp,
	mail mail.ServiceMail,
	sms sms.ServiceSMS,
	twilio twilio.ServiceTwilio,
) *Service {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {

		log.Panicf("unable to initialize Firebase app for OTP service: %s", err)
	}
	ctx := context.Background()

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {

		log.Panicf("unable to initialize Firestore client: %s", err)
	}

	if err != nil {

		log.Errorf("occurred while opening deps file %v", err)
		os.Exit(1)
	}

	return &Service{
		totpOpts: totp.GenerateOpts{
			Issuer:      issuer,
			AccountName: accountName,
		},
		firestoreClient:      firestoreClient,
		rootCollectionSuffix: serverutils.MustGetEnvVar("ROOT_COLLECTION_SUFFIX"),
		whatsapp:             whatsapp,
		mail:                 mail,
		sms:                  sms,
		twilio:               twilio,
	}
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("OTP service has a nil firestore client")
	}

	if s.twilio == nil {
		log.Panicf("OTP service needs to define a twilio client ")
	}

}

func (s Service) getOTPCollectionName() string {
	return firebasetools.SuffixCollection(converterandformatter.OTPCollectionName)
}

func cleanITPhoneNumber() (*string, error) {
	return converterandformatter.NormalizeMSISDN(ITPhoneNumber)
}

// SendOTP sends otp code message to specified number
func (s Service) SendOTP(
	ctx context.Context,
	normalizedPhoneNumber string,
	code string,
	appID *string,
) (string, error) {
	ctx, span := tracer.Start(ctx, "SendOTP")
	defer span.End()

	var appidentifier string
	if appID == nil {
		appID = &appidentifier
	}
	msg := fmt.Sprintf(otpMsg, code, *appID)

	if interserviceclient.IsKenyanNumber(normalizedPhoneNumber) {
		_, err := s.sms.Send(ctx, normalizedPhoneNumber, msg, enumutils.SenderIDBewell)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return "", fmt.Errorf("failed to send OTP verification message to recipient")
		}
	} else {
		// Make request to twilio
		err := s.twilio.SendSMS(ctx, normalizedPhoneNumber, msg)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return "", fmt.Errorf("sms not sent: %v", err)
		}
	}

	return code, nil
}

// GenerateAndSendOTP creates an OTP and sends it to the
// supplied phone number as a text message
func (s Service) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	cleanNo, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {

		return "", errors.Wrap(err, "generateOTP > NormalizeMSISDN")
	}

	// This is an alternate path that checks for the Constant
	// phone number used by our Frontend's Integration Testing.
	// This path returns a constant OTP familiar to both the teams.
	cleanITPhoneNumber, err := cleanITPhoneNumber()
	if err != nil {

		return "", errors.Wrap(err, "failed to normalize Integration Test number")
	}

	if *cleanNo == *cleanITPhoneNumber {
		return ITCode, nil
	}

	code, err := s.GenerateOTP(ctx)
	if err != nil {

		return "", errors.Wrap(err, "Unable to generate OTP")
	}

	var appidentifier string
	if appID == nil {
		appID = &appidentifier
	}
	msg := fmt.Sprintf(otpMsg, code, *appID)
	otp := dto.OTP{
		MSISDN:            msisdn,
		Message:           msg,
		AuthorizationCode: code,
		Timestamp:         time.Now(),
		IsValid:           true,
	}
	err = s.SaveOTPToFirestore(otp)
	if err != nil {

		return code, fmt.Errorf("unable to save OTP: %v", err)
	}

	code, err = s.SendOTP(ctx, *cleanNo, code, appID)
	if err != nil {

		log.Printf("OTP send error: %s", err)
		return code, err
	}
	return code, nil
}

//SendOTPToEmail is a companion to GenerateAndSendOTP function
//It will send the generated OTP to the provided email address
func (s Service) SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error) {
	_, span := tracer.Start(ctx, "SendOTPToEmail")
	defer span.End()
	code, err := s.GenerateAndSendOTP(ctx, *msisdn, appID)
	if err != nil {
		helpers.RecordSpanError(span, err)
		log.Printf("error: %s", err)
		return code, err
	}

	// If the code returned is the specified const for Integration Testing,
	// halt execution to prevent sending of a real email
	if code == ITCode {
		return code, nil
	}

	text := GenerateEmailFunc(code)

	if email == nil {
		return code, nil
	}

	emailstr := *email

	if !govalidator.IsEmail(emailstr) {
		return code, fmt.Errorf("%s is not a valid email", emailstr)
	}

	_, _, err = s.mail.SendEmail(ctx, subject, text, nil, emailstr)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return code, fmt.Errorf("unable to send OTP to email: %w", err)
	}

	return code, nil
}

// SaveOTPToFirestore persists the supplied OTP
func (s Service) SaveOTPToFirestore(otp dto.OTP) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.getOTPCollectionName()).Add(ctx, otp)
	return err
}

// VerifyOtp checks for the validity of the supplied OTP but does not invalidate it
func (s Service) VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error) {
	ctx, span := tracer.Start(ctx, "VerifyOtp")
	defer span.End()
	s.checkPreconditions()

	// ensure the phone number passed is correct
	normalizeMsisdn, err := converterandformatter.NormalizeMSISDN(*msisdn)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, errors.Wrap(err, "VerifyOtp > NormalizeMSISDN")
	}

	// This is an alternate path that checks for the Constant
	// phone number and OTP used by our Frontend's Integration Testing.
	// This path always verify the OTP if the condition matches i.e
	// returns a true and a nil.
	cleanITPhoneNumber, err := cleanITPhoneNumber()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, errors.Wrap(err, "failed to normalize Integration Test number")
	}

	if *normalizeMsisdn == *cleanITPhoneNumber && *verificationCode == ITCode {
		return true, nil
	}

	// check if the OTP is on file / known
	collection := s.firestoreClient.Collection(s.getOTPCollectionName())
	query := collection.Where(
		"isValid", "==", true,
	).Where(
		"msisdn", "==", normalizeMsisdn,
	).Where(
		"authorizationCode", "==", *verificationCode,
	)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return false, fmt.Errorf("no matching verification codes found")
	}
	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = firebasetools.UpdateRecordOnFirestore(
			s.firestoreClient, s.getOTPCollectionName(), doc.Ref.ID, otpData)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return false, fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}
	return true, nil
}

// VerifyEmailOtp checks for the validity of the supplied OTP but does not invalidate it
func (s Service) VerifyEmailOtp(
	ctx context.Context,
	email, verificationCode *string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "VerifyEmailOtp")
	defer span.End()
	s.checkPreconditions()

	// This is an alternate path that checks for the Constant
	// email and OTP used by our Frontend's Integration Testing.
	// This path always verify the OTP if the condition matches i.e
	// returns a true and a nil.
	if *email == ITEmail && *verificationCode == ITCode {
		return true, nil
	}

	// check if the OTP is on file / known
	collection := s.firestoreClient.Collection(s.getOTPCollectionName())
	query := collection.Where(
		"isValid", "==", true,
	).Where(
		"email", "==", *email,
	).Where(
		"authorizationCode", "==", *verificationCode,
	)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return false, fmt.Errorf("no matching verification codes found")
	}
	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = firebasetools.UpdateRecordOnFirestore(
			s.firestoreClient, s.getOTPCollectionName(), doc.Ref.ID, otpData)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return false, fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}
	return true, nil
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (s Service) GenerateRetryOTP(
	ctx context.Context,
	msisdn *string,
	retryStep int,
	appID *string,
) (string, error) {
	ctx, span := tracer.Start(ctx, "GenerateRetryOTP")
	defer span.End()
	cleanNo, err := converterandformatter.NormalizeMSISDN(*msisdn)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", errors.Wrap(err, "generateOTP > NormalizeMSISDN")
	}

	// This is an alternate path that checks for the Constant
	// phone number used by our Frontend's Integration Testing.
	// This path returns a constant retry OTP familiar to both the teams.
	cleanITPhoneNumber, err := cleanITPhoneNumber()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", errors.Wrap(err, "failed to normalize Integration Test number")
	}

	if *cleanNo == *cleanITPhoneNumber {
		return ITCode, nil
	}

	code, err := s.GenerateOTP(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		log.Printf("error: %s", err)
		return "", fmt.Errorf("OTP generation failed: %w", err)
	}

	var appidentifier string
	if appID == nil {
		appID = &appidentifier
	}
	msg := fmt.Sprintf(otpMsg, code, *appID)

	otp := dto.OTP{
		MSISDN:            *msisdn,
		Message:           msg,
		AuthorizationCode: code,
		Timestamp:         time.Now(),
		IsValid:           true,
	}

	err = s.SaveOTPToFirestore(otp)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return code, fmt.Errorf("unable to save OTP: %v", err)
	}

	if retryStep == whatsappStep {

		sent, err := s.whatsapp.PhoneNumberVerificationCode(
			ctx,
			otp.MSISDN,
			otp.AuthorizationCode,
			otp.Message,
		)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return code, fmt.Errorf("unable to send a phone verification code :%w", err)
		}

		if !sent {
			return "", fmt.Errorf("unable to send OTP whatsapp message : %w", err)
		}

		return code, nil

	} else if retryStep == twilioStep {
		err := s.twilio.SendSMS(ctx, otp.MSISDN, otp.Message)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return code, fmt.Errorf("otp send retry failed: %w", err)
		}
		return code, nil

	} else {
		return "", fmt.Errorf("invalid retry step")
	}

}

// EmailVerificationOtp generates an OTP to the supplied email for verification
func (s Service) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	_, span := tracer.Start(ctx, "EmailVerificationOtp")
	defer span.End()
	// This is an alternate path that checks for the Constant
	// email used by our Frontend's Integration Testing.
	// This path returns a constant OTP familiar to both the teams.
	if *email == ITEmail {
		return ITCode, nil
	}
	code, err := s.GenerateOTP(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", errors.Wrap(err, "Unable to generate OTP")
	}

	text := GenerateEmailFunc(code)
	if !govalidator.IsEmail(*email) {
		return code, fmt.Errorf("%s is not a valid email", *email)
	}

	msg := fmt.Sprintf("Your phone number verification code is %s. ", code)

	otp := dto.OTP{
		Email:             *email,
		Message:           msg,
		AuthorizationCode: code,
		Timestamp:         time.Now(),
		IsValid:           true,
	}
	err = s.SaveOTPToFirestore(otp)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return code, fmt.Errorf("unable to save OTP: %v", err)
	}

	emailstr := *email

	_, _, err = s.mail.SendEmail(ctx, subject, text, nil, emailstr)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return code, fmt.Errorf("unable to send OTP to email: %w", err)
	}

	return code, nil
}

//GenerateOTP generates an OTP
func (s Service) GenerateOTP(ctx context.Context) (string, error) {
	_, span := tracer.Start(ctx, "GenerateOTP")
	defer span.End()
	key, err := totp.Generate(s.totpOpts)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", errors.Wrap(err, "generateOTP")
	}

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", errors.Wrap(err, "generateOTP > GenerateCode")
	}
	return code, nil
}
