package otp_test

// TODO - restore tests
// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/otp"

// 	"github.com/stretchr/testify/assert"
// 	"gitlab.slade360emr.com/go/base"
// )

// const (
// 	InternationalTestUserPhoneNumber = "+12028569601"
// 	ValidTestEmail                   = "automated.test.user.bewell-app-ci@healthcloud.co.ke"
// )

// func TestMain(m *testing.M) {
// 	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
// 	os.Setenv("ENVIRONMENT", "staging")
// 	os.Exit(m.Run())
// }

// func TestNormalizeMSISDN(t *testing.T) {
// 	tests := map[string]struct {
// 		input                string
// 		expectError          bool
// 		expectedErrorMessage string
// 		expectedOutput       string
// 	}{
// 		"invalid_number": {
// 			input:                "1",
// 			expectError:          true,
// 			expectedErrorMessage: "invalid phone number: 1",
// 		},
// 		"international_with_plus": {
// 			input:          profileutils.TestUserPhoneNumber,
// 			expectError:    false,
// 			expectedOutput: profileutils.TestUserPhoneNumber,
// 		},
// 		"international_no_plus": {
// 			input:          profileutils.TestUserPhoneNumber[1:],
// 			expectError:    false,
// 			expectedOutput: profileutils.TestUserPhoneNumber,
// 		},
// 		"national_zero_prefix": {
// 			input:          "0" + profileutils.TestUserPhoneNumber[4:],
// 			expectError:    false,
// 			expectedOutput: profileutils.TestUserPhoneNumber,
// 		},
// 		"national_no_zero_prefix": {
// 			input:                profileutils.TestUserPhoneNumber[4:],
// 			expectError:          true,
// 			expectedErrorMessage: "invalid phone number: 711223344",
// 		},
// 	}
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			got, err := converterandformatter.NormalizeMSISDN(tc.input)
// 			if tc.expectError {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, tc.expectedErrorMessage, err.Error())
// 			}

// 			if !tc.expectError {
// 				assert.Nil(t, err)
// 				assert.Equal(t, tc.expectedOutput, *got)
// 			}
// 		})
// 	}
// }

// func TestGenerateAndSendOTP(t *testing.T) {
// 	otpService := otp.NewService()

// 	tests := map[string]struct {
// 		msisdn               string
// 		expectError          bool
// 		expectedErrorMessage string
// 	}{
// 		"valid_case : Kenyan number": {
// 			msisdn: profileutils.TestUserPhoneNumber,
// 		},
// 		"valid_case : american number": {
// 			msisdn: InternationalTestUserPhoneNumber,
// 		},
// 		"valid_case : Integration test number": {
// 			msisdn: otp.ITPhoneNumber,
// 		},
// 	}
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			code, err := otpService.GenerateAndSendOTP(tc.msisdn)
// 			if tc.expectError {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, tc.expectedErrorMessage, err.Error())
// 			}
// 			if !tc.expectError {
// 				assert.NotZero(t, code)
// 				assert.Nil(t, err)
// 				if tc.msisdn == otp.ITPhoneNumber {
// 					assert.Equal(t, code, otp.ITCode)
// 				}
// 			}
// 		})
// 	}
// }

// func TestService_SendOTPToEmail(t *testing.T) {
// 	phoneNumber := profileutils.TestUserPhoneNumber
// 	invalidPhoneNumber := "1"
// 	otpService := otp.NewService()
// 	ctx := context.Background()
// 	validEmail := ValidTestEmail
// 	invalidEmail := "This is not an email"
// 	testPhone := otp.ITPhoneNumber
// 	testEmail := otp.ITEmail
// 	tests := map[string]struct {
// 		msisdn               *string
// 		email                *string
// 		expectError          bool
// 		expectedErrorMessage string
// 	}{
// 		"valid_case": {
// 			msisdn:      &phoneNumber,
// 			email:       &validEmail,
// 			expectError: false,
// 		},
// 		"invalid_email": {
// 			msisdn:               &phoneNumber,
// 			email:                &invalidEmail,
// 			expectError:          true,
// 			expectedErrorMessage: "This is not an email is not a valid email",
// 		},
// 		"invalid phone": {
// 			msisdn:               &invalidPhoneNumber,
// 			email:                &validEmail,
// 			expectError:          true,
// 			expectedErrorMessage: "generateOTP > NormalizeMSISDN: invalid phone number: 1",
// 		},
// 		"valid phone and nil email": {
// 			msisdn:      &phoneNumber,
// 			email:       nil,
// 			expectError: false,
// 		},
// 		"valid_case_integration_test": {
// 			msisdn:      &testPhone,
// 			email:       &testEmail,
// 			expectError: false,
// 		},
// 	}
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			code, err := otpService.SendOTPToEmail(ctx, tc.msisdn, tc.email)
// 			if tc.expectError {
// 				assert.NotNil(t, err)
// 				assert.Contains(t, tc.expectedErrorMessage, err.Error())
// 			}
// 			if !tc.expectError {
// 				assert.NotZero(t, code)
// 				assert.Nil(t, err)
// 				if tc.msisdn == &testPhone && tc.email == &testEmail {
// 					assert.Equal(t, code, otp.ITCode)
// 				}
// 			}
// 		})
// 	}
// }

// // TODO: Test Mocking
// // disabled retry OTP tests due to service costs

// func TestService_GenerateRetryOtp(t *testing.T) {
// 	phoneNumber := profileutils.TestUserPhoneNumber
// 	invalidPhoneNumber := "this is definitely not a number"
// 	service := otp.NewService()
// 	type args struct {
// 		ctx       context.Context
// 		msisdn    *string
// 		retryStep int
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Generate retry OTP whatsapp happy case",
// 			args: args{
// 				ctx:       context.Background(),
// 				msisdn:    &phoneNumber,
// 				retryStep: 1,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Generate retry OTP twilio happy case",
// 			args: args{
// 				ctx:       context.Background(),
// 				msisdn:    &phoneNumber,
// 				retryStep: 2,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Generate retry OTP sad case",
// 			args: args{
// 				ctx:    context.Background(),
// 				msisdn: &invalidPhoneNumber,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			otp, err := s.GenerateRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep)
// 			if err == nil {
// 				assert.NotNil(t, otp)
// 			}
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.GenerateRetryOtp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestService_GenerateOTP(t *testing.T) {
// 	service := otp.NewService()
// 	tests := []struct {
// 		name    string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Generate OTP",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			otp, err := s.GenerateOTP()
// 			if err == nil {
// 				assert.NotNil(t, otp)
// 			}
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestService_EmailVerificationOtp(t *testing.T) {
// 	email := ValidTestEmail
// 	invalidEmail := "not an email address"
// 	integrationTestEmail := otp.ITEmail

// 	service := otp.NewService()
// 	type args struct {
// 		ctx   context.Context
// 		email *string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid email",
// 			args: args{
// 				ctx:   context.Background(),
// 				email: &email,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "invalid email",
// 			args: args{
// 				ctx:   context.Background(),
// 				email: &invalidEmail,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "valid I.T email",
// 			args: args{
// 				ctx:   context.Background(),
// 				email: &integrationTestEmail,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			code, err := s.EmailVerificationOtp(tt.args.ctx, tt.args.email)
// 			if err == nil {
// 				assert.NotNil(t, code)
// 				if tt.args.email == &integrationTestEmail {
// 					assert.Equal(t, code, otp.ITCode)
// 				}
// 			}
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.EmailVerificationOtp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestService_VerifyOtp(t *testing.T) {
// 	phoneNumber := profileutils.TestUserPhoneNumber
// 	invalidNumber := "1111"
// 	srv := otp.NewService()
// 	assert.NotNil(t, srv, "service should not be bil")
// 	ctx := context.Background()
// 	// generate the otp
// 	otp_code, err := srv.GenerateRetryOTP(ctx, &phoneNumber, 1)
// 	if err == nil {
// 		assert.NotNil(t, otp_code)
// 	}

// 	testPhone := otp.ITPhoneNumber
// 	testCode := otp.ITCode

// 	type args struct {
// 		ctx              context.Context
// 		msisdn           *string
// 		verificationCode *string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "verify otp happy case",
// 			args: args{
// 				ctx:              ctx,
// 				msisdn:           &phoneNumber,
// 				verificationCode: &otp_code,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 		{
// 			name: "verify otp invalid phonenumber",
// 			args: args{
// 				ctx:              ctx,
// 				msisdn:           &invalidNumber,
// 				verificationCode: &otp_code,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 		{
// 			name: "verify otp I.T case",
// 			args: args{
// 				ctx:              ctx,
// 				msisdn:           &testPhone,
// 				verificationCode: &testCode,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := srv.VerifyOtp(tt.args.ctx, tt.args.msisdn, tt.args.verificationCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.VerifyOtp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("Service.VerifyOtp() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestService_VerifyEmailOtp(t *testing.T) {
// 	s := otp.NewService()
// 	ctx := context.Background()
// 	email := ValidTestEmail
// 	testEmail := otp.ITEmail
// 	testCode := otp.ITCode
// 	randomCode := "random"
// 	otp, err := s.EmailVerificationOtp(ctx, &email)
// 	if err != nil {
// 		t.Errorf("failed to send test email OTP")
// 		return
// 	}
// 	type args struct {
// 		ctx              context.Context
// 		email            *string
// 		verificationCode *string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 		want    bool
// 	}{
// 		{
// 			name: "happy case",
// 			args: args{
// 				ctx:              ctx,
// 				email:            &email,
// 				verificationCode: &otp,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 		{
// 			name: "happy case - integration tests",
// 			args: args{
// 				ctx:              ctx,
// 				email:            &testEmail,
// 				verificationCode: &testCode,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 		{
// 			name: "sad case",
// 			args: args{
// 				ctx:              ctx,
// 				email:            &testEmail,
// 				verificationCode: &randomCode,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			verify, err := s.VerifyEmailOtp(tt.args.ctx, tt.args.email, tt.args.verificationCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.VerifyEmailOtp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if verify != tt.want {
// 				t.Errorf("Service.VerifyEmailOtp() = %v, want %v", verify, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_sendOtp(t *testing.T) {
// 	ctx := context.Background()
// 	type args struct {
// 		normalizedPhoneNumber string
// 		code                  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid normalized Kenyan number",
// 			args: args{
// 				normalizedPhoneNumber: profileutils.TestUserPhoneNumber,
// 				code:                  "123456",
// 			},
// 			want:    "123456",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := otp.NewService()
// 			got, err := s.SendOTP(ctx, tt.args.normalizedPhoneNumber, tt.args.code)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("SendOTP() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("SendOTP() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
