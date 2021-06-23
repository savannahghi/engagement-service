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
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/sms"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/twilio"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/whatsapp"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

const (
	issuer       = "Savannah Informatics Limited"
	accountName  = "info@healthcloud.co.ke"
	subject      = "Be.Well Verification Code"
	whatsappStep = 1
	twilioStep   = 2
)

// These constants are here to support Integration Testing
// done by the Frontend team.
// IT is shorthand for Integration Tests
const (
	ITCode        = "123456"
	ITPhoneNumber = "0798000000"
	ITEmail       = "test@healthcloud.co.ke"
)

// ServiceOTP is an interface that defines all interactions with OTP service
type ServiceOTP interface {
	GenerateAndSendOTP(msisdn string) (string, error)
	SendOTPToEmail(ctx context.Context, msisdn, email *string) (string, error)
	SaveOTPToFirestore(otp dto.OTP) error
	VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error)
	VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error)
	GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int) (string, error)
	EmailVerificationOtp(ctx context.Context, email *string) (string, error)
	GenerateOTP() (string, error)
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
// struct definition explictly, No guess work.
func NewService() *Service {
	var repository repository.Repository

	whatsapp := whatsapp.NewService()

	mail := mail.NewService()

	sms := sms.NewService(repository)

	twilio := twilio.NewService()

	fc := &base.FirebaseClient{}
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
		rootCollectionSuffix: base.MustGetEnvVar("ROOT_COLLECTION_SUFFIX"),
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
	return base.SuffixCollection(base.OTPCollectionName)
}

func cleanITPhoneNumber() (*string, error) {
	return base.NormalizeMSISDN(ITPhoneNumber)
}

// SendOTP sends otp code message to specified number
func (s Service) SendOTP(ctx context.Context, normalizedPhoneNumber string, code string) (string, error) {
	msg := fmt.Sprintf("Your phone number verification code is %s. ", code)

	if base.IsKenyanNumber(normalizedPhoneNumber) {
		_, err := s.sms.Send(normalizedPhoneNumber, msg, base.SenderIDBewell)
		if err != nil {
			return "", fmt.Errorf("failed to send OTP verification message to recipient")
		}
	} else {
		// Make request to twilio
		err := s.twilio.SendSMS(ctx, normalizedPhoneNumber, msg)
		if err != nil {
			return "", fmt.Errorf("sms not sent: %v", err)
		}
	}

	return code, nil
}

// GenerateAndSendOTP creates an OTP and sends it to the
// supplied phone number as a text message
func (s Service) GenerateAndSendOTP(msisdn string) (string, error) {
	cleanNo, err := base.NormalizeMSISDN(msisdn)
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

	ctx := context.Background()
	code, err := s.GenerateOTP()
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate OTP")
	}

	msg := fmt.Sprintf("Your phone number verification code is %s. ", code)
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

	code, err = s.SendOTP(ctx, *cleanNo, code)
	if err != nil {
		log.Printf("OTP send error: %s", err)
		return code, err
	}
	return code, nil
}

//SendOTPToEmail is a companion to GenerateAndSendOTP function
//It will send the generated OTP to the provided email address
func (s Service) SendOTPToEmail(ctx context.Context, msisdn, email *string) (string, error) {
	code, err := s.GenerateAndSendOTP(*msisdn)
	if err != nil {
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

	_, _, err = s.mail.SendEmail(subject, text, emailstr)
	if err != nil {
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
	s.checkPreconditions()

	// ensure the phone number passed is correct
	normalizeMsisdn, err := base.NormalizeMSISDN(*msisdn)
	if err != nil {
		return false, errors.Wrap(err, "VerifyOtp > NormalizeMSISDN")
	}

	// This is an alternate path that checks for the Constant
	// phone number and OTP used by our Frontend's Integration Testing.
	// This path always verify the OTP if the condition matches i.e
	// returns a true and a nil.
	cleanITPhoneNumber, err := cleanITPhoneNumber()
	if err != nil {
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
		return false, fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return false, fmt.Errorf("no matching verification codes found")
	}
	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.getOTPCollectionName(), doc.Ref.ID, otpData)
		if err != nil {
			return false, fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}
	return true, nil
}

// VerifyEmailOtp checks for the validity of the supplied OTP but does not invalidate it
func (s Service) VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error) {
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
		return false, fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return false, fmt.Errorf("no matching verification codes found")
	}
	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.getOTPCollectionName(), doc.Ref.ID, otpData)
		if err != nil {
			return false, fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}
	return true, nil
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (s Service) GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int) (string, error) {
	cleanNo, err := base.NormalizeMSISDN(*msisdn)
	if err != nil {
		return "", errors.Wrap(err, "generateOTP > NormalizeMSISDN")
	}

	// This is an alternate path that checks for the Constant
	// phone number used by our Frontend's Integration Testing.
	// This path returns a constant retry OTP familiar to both the teams.
	cleanITPhoneNumber, err := cleanITPhoneNumber()
	if err != nil {
		return "", errors.Wrap(err, "failed to normalize Integration Test number")
	}

	if *cleanNo == *cleanITPhoneNumber {
		return ITCode, nil
	}

	code, err := s.GenerateOTP()
	if err != nil {
		log.Printf("error: %s", err)
		return "", fmt.Errorf("OTP generation failed: %w", err)
	}

	msg := fmt.Sprintf("Your phone number verification code is %s. ", code)

	otp := dto.OTP{
		MSISDN:            *msisdn,
		Message:           msg,
		AuthorizationCode: code,
		Timestamp:         time.Now(),
		IsValid:           true,
	}

	err = s.SaveOTPToFirestore(otp)
	if err != nil {
		return code, fmt.Errorf("unable to save OTP: %v", err)
	}

	if retryStep == whatsappStep {

		sent, err := s.whatsapp.PhoneNumberVerificationCode(ctx, otp.MSISDN, otp.AuthorizationCode, otp.Message)
		if err != nil {
			return code, fmt.Errorf("unable to send a phone verification code :%w", err)
		}

		if !sent {
			return "", fmt.Errorf("unable to send OTP whatsapp message : %w", err)
		}

		return code, nil

	} else if retryStep == twilioStep {
		err := s.twilio.SendSMS(ctx, otp.MSISDN, otp.Message)
		if err != nil {
			return code, fmt.Errorf("otp send retry failed: %w", err)
		}
		return code, nil

	} else {
		return "", fmt.Errorf("invalid retry step")
	}

}

// EmailVerificationOtp generates an OTP to the supplied email for verification
func (s Service) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	// This is an alternate path that checks for the Constant
	// email used by our Frontend's Integration Testing.
	// This path returns a constant OTP familiar to both the teams.
	if *email == ITEmail {
		return ITCode, nil
	}
	code, err := s.GenerateOTP()
	if err != nil {
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
		return code, fmt.Errorf("unable to save OTP: %v", err)
	}

	emailstr := *email

	_, _, err = s.mail.SendEmail(subject, text, emailstr)
	if err != nil {
		return code, fmt.Errorf("unable to send OTP to email: %w", err)
	}

	return code, nil
}

//GenerateOTP generates an OTP
func (s Service) GenerateOTP() (string, error) {
	key, err := totp.Generate(s.totpOpts)
	if err != nil {
		return "", errors.Wrap(err, "generateOTP")
	}

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", errors.Wrap(err, "generateOTP > GenerateCode")
	}
	return code, nil
}
