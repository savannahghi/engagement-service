package otp_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	mockMail "github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail/mock"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	mockSMS "github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms/mock"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio"
	mockTwilio "github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio/mock"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	mockWhatsapp "github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp/mock"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/interserviceclient"
)

const (
	InternationalTestUserPhoneNumber = "+14432967215"
	ValidTestEmail                   = "automated.test.user.bewell-app-ci@healthcloud.co.ke"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("ENVIRONMENT", "staging")
	os.Exit(m.Run())
}

var fakeSMS mockSMS.FakeSMSService
var smsServ sms.ServiceSMS = &fakeSMS

var fakeWhatsapp mockWhatsapp.FakeServiceWhatsapp
var whatsappServ whatsapp.ServiceWhatsapp = &fakeWhatsapp

var fakeTwilio mockTwilio.FakeTwilioService
var twilioServ twilio.ServiceTwilio = &fakeTwilio

var fakeMail mockMail.FakeMailService
var mailServ mail.ServiceMail = &fakeMail

func TestService_SendOTP(t *testing.T) {

	ctx := context.Background()

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	type args struct {
		ctx                   context.Context
		normalizedPhoneNumber string
		code                  string
		appID                 string
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		// test cases.
		{
			name: "valid case:Kenyan number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: interserviceclient.TestUserPhoneNumber,
				appID:                 "APP_IDENTIFIER",
			},
			wantErr: false,
		},
		{
			name: "valid case:American number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: InternationalTestUserPhoneNumber,
				appID:                 "APP_IDENTIFIER",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid case:American number" {

				fakeTwilio.SendSMSFn = func(ctx context.Context, to, msg string) error {
					return nil
				}

			}

			if tt.name == "valid case:Kenyan number" {

				fakeSMS.SendFn = func(ctx context.Context, to, message string, from enumutils.SenderID) (*dto.SendMessageResponse, error) {
					return &dto.SendMessageResponse{}, nil
				}
			}

			_, err := s.SendOTP(tt.args.ctx, tt.args.normalizedPhoneNumber, tt.args.code, &tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendOTP() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestService_GenerateAndSendOTP(t *testing.T) {

	ctx := context.Background()

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)
	internationalNumber := InternationalTestUserPhoneNumber

	type args struct {
		ctx    context.Context
		msisdn string
		appID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// test cases.
		{
			name: "valid case:Kenyan number",
			args: args{
				ctx:    ctx,
				msisdn: interserviceclient.TestUserPhoneNumber,
				appID:  "APP_IDENTIFIER",
			},
			wantErr: false,
		},
		{
			name: "valid case:American number",
			args: args{
				ctx:    ctx,
				msisdn: internationalNumber,
				appID:  "APP_IDENTIFIER",
			},
			wantErr: false,
		},
		{
			name: "invalid case:empty number",
			args: args{
				ctx:    ctx,
				msisdn: "",
				appID:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid case:American number" {

				fakeTwilio.SendSMSFn = func(ctx context.Context, to, msg string) error {
					return nil
				}

			}

			if tt.name == "valid case:Kenyan number" {

				fakeSMS.SendFn = func(ctx context.Context, to, message string, from enumutils.SenderID) (*dto.SendMessageResponse, error) {
					return &dto.SendMessageResponse{}, nil
				}
			}

			_, err := s.GenerateAndSendOTP(tt.args.ctx, tt.args.msisdn, &tt.args.appID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestService_SendOTPToEmail(t *testing.T) {

	ctx := context.Background()

	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidPhoneNumber := "gabriel"

	validEmail := ValidTestEmail
	invalidEmail := "gabriel"

	testPhone := otp.ITPhoneNumber

	testEmail := otp.ITEmail

	appID := "APP_IDENTIFIER"

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	type args struct {
		ctx    context.Context
		msisdn *string
		email  *string
		appID  *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//test cases.
		{
			name: "valid details",
			args: args{
				ctx:    ctx,
				msisdn: &phoneNumber,
				email:  &validEmail,
				appID:  &appID,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				ctx:    ctx,
				msisdn: &phoneNumber,
				email:  &invalidEmail,
				appID:  &appID,
			},
			wantErr: true,
		},
		{
			name: "invalid phoneNumber",
			args: args{
				ctx:    ctx,
				msisdn: &invalidPhoneNumber,
				email:  &validEmail,
				appID:  &appID,
			},
			wantErr: true,
		},
		{
			name: "valid case integration test details",
			args: args{
				ctx:    ctx,
				msisdn: &testPhone,
				email:  &testEmail,
				appID:  &appID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid details" {
				fakeSMS.SendFn = func(ctx context.Context, to, message string, from enumutils.SenderID) (*dto.SendMessageResponse, error) {
					return &dto.SendMessageResponse{}, nil
				}
				fakeMail.SendEmailFn = func(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error) {
					return "123", "123", nil
				}
			}

			_, err := s.SendOTPToEmail(tt.args.ctx, tt.args.msisdn, tt.args.email, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendOTPToEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}

		})
	}
}

func TestService_VerifyOtp(t *testing.T) {

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	phoneNumber := interserviceclient.TestUserPhoneNumber

	invalidNumber := "gabriel"

	testPhone := otp.ITPhoneNumber
	testCode := otp.ITCode

	ctx := context.Background()

	otpCode := "654321"

	// save the test OTP to firestore first
	err := s.SaveOTPToFirestore(dto.OTP{
		MSISDN:            phoneNumber,
		Message:           "this is a test otp",
		AuthorizationCode: otpCode,
		Timestamp:         time.Now(),
		IsValid:           true,
	})

	if err != nil {
		t.Errorf("failed to save test otp in the database")
		return
	}
	type args struct {
		ctx              context.Context
		msisdn           *string
		verificationCode *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		//test cases.
		{
			name: "verify otp happy case",
			args: args{
				ctx:              ctx,
				msisdn:           &phoneNumber,
				verificationCode: &otpCode,
			},
			wantErr: false,
		},
		{
			name: "verify otp invalid phonenumber",
			args: args{
				ctx:              ctx,
				msisdn:           &invalidNumber,
				verificationCode: &otpCode,
			},
			wantErr: true,
		},
		{
			name: "verify otp I.T case",
			args: args{
				ctx:              ctx,
				msisdn:           &testPhone,
				verificationCode: &testCode,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := s.VerifyOtp(tt.args.ctx, tt.args.msisdn, tt.args.verificationCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyOtp() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestService_VerifyEmailOtp(t *testing.T) {

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	ctx := context.Background()

	email := ValidTestEmail

	testEmail := otp.ITEmail
	testCode := otp.ITCode

	invalidCode := "gabriel"

	otpCode := "123456"

	//save the otp to firestore firsts
	err := s.SaveOTPToFirestore(dto.OTP{
		Email:             email,
		Message:           "this is a test otp",
		AuthorizationCode: otpCode,
		Timestamp:         time.Now(),
		IsValid:           true,
	})
	if err != nil {
		t.Errorf("failed to send test email OTP")
		return
	}
	type args struct {
		ctx              context.Context
		email            *string
		verificationCode *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// test cases.
		{
			name: "happy case",
			args: args{
				ctx:              ctx,
				email:            &email,
				verificationCode: &otpCode,
			},
			wantErr: false,
		},
		{
			name: "happy case - integration tests",
			args: args{
				ctx:              ctx,
				email:            &testEmail,
				verificationCode: &testCode,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:              ctx,
				email:            &testEmail,
				verificationCode: &invalidCode,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := s.VerifyEmailOtp(tt.args.ctx, tt.args.email, tt.args.verificationCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyEmailOtp() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}

		})
	}
}

func TestService_GenerateRetryOTP(t *testing.T) {

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	ctx := context.Background()

	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidPhoneNumber := "gabriel"

	type args struct {
		ctx       context.Context
		msisdn    *string
		retryStep int
		appID     *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// test cases.
		{
			name: "Generate retry OTP whatsapp happy case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 1,
			},
			wantErr: false,
		},
		{
			name: "Generate retry OTP twilio happy case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 2,
			},
			wantErr: false,
		},
		{
			name: "Generate retry OTP sad case",
			args: args{
				ctx:    context.Background(),
				msisdn: &invalidPhoneNumber,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Generate retry OTP whatsapp happy case" {
				fakeWhatsapp.PhoneNumberVerificationCodeFn = func(ctx context.Context, to, code, marketingMessage string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Generate retry OTP twilio happy case" {
				fakeTwilio.SendSMSFn = func(ctx context.Context, to, msg string) error {
					return nil
				}
			}

			_, err := s.GenerateRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateRetryOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}

		})
	}
}

func TestService_EmailVerificationOtp(t *testing.T) {

	email := ValidTestEmail
	invalidEmail := "gabriel"

	integrationTestEmail := otp.ITEmail

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)
	type args struct {
		ctx   context.Context
		email *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				ctx:   context.Background(),
				email: &email,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				ctx:   context.Background(),
				email: &invalidEmail,
			},
			wantErr: true,
		},
		{
			name: "valid I.T email",
			args: args{
				ctx:   context.Background(),
				email: &integrationTestEmail,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid email" {
				fakeMail.SendEmailFn = func(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error) {
					return "gabriel", "were", nil
				}
			}

			_, err := s.EmailVerificationOtp(tt.args.ctx, tt.args.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.EmailVerificationOtp() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestService_GenerateOTP(t *testing.T) {

	s := otp.NewService(whatsappServ, mailServ, smsServ, twilioServ)

	ctx := context.Background()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GenerateOTP(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
