package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

//FakeOtpService ...
type FakeOtpService struct {
	GenerateAndSendOTPFn   func(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmailFn       func(ctx context.Context, msisdn, email *string, appID *string) (string, error)
	SaveOTPToFirestoreFn   func(otp dto.OTP) error
	VerifyOtpFn            func(ctx context.Context, msisdn, verificationCode *string) (bool, error)
	VerifyEmailOtpFn       func(ctx context.Context, email, verificationCode *string) (bool, error)
	GenerateRetryOTPFn     func(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error)
	EmailVerificationOtpFn func(ctx context.Context, email *string) (string, error)
	GenerateOTPFn          func(ctx context.Context) (string, error)
	SendOTPFn              func(ctx context.Context, normalizedPhoneNumber string, code string, appID *string) (string, error)
}

//GenerateAndSendOTP mock function
func (f *FakeOtpService) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return f.GenerateAndSendOTPFn(ctx, msisdn, appID)
}

//SendOTPToEmail mock function
func (f *FakeOtpService) SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error) {
	return f.SendOTPToEmailFn(ctx, msisdn, email, appID)
}

//SaveOTPToFirestore mock function
func (f *FakeOtpService) SaveOTPToFirestore(otp dto.OTP) error {
	return f.SaveOTPToFirestoreFn(otp)
}

//VerifyOtp mock function
func (f *FakeOtpService) VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error) {
	return f.VerifyOtpFn(ctx, msisdn, verificationCode)
}

//VerifyEmailOtp mock function
func (f *FakeOtpService) VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error) {
	return f.VerifyEmailOtpFn(ctx, email, verificationCode)
}

//GenerateRetryOTP mock function
func (f *FakeOtpService) GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error) {
	return f.GenerateRetryOTPFn(ctx, msisdn, retryStep, appID)
}

//EmailVerificationOtp mock function
func (f *FakeOtpService) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	return f.EmailVerificationOtpFn(ctx, email)
}

//GenerateOTP mock function
func (f *FakeOtpService) GenerateOTP(ctx context.Context) (string, error) {
	return f.GenerateOTPFn(ctx)
}

//SendOTP mock function
func (f *FakeOtpService) SendOTP(ctx context.Context, normalizedPhoneNumber string, code string, appID *string) (string, error) {
	return f.SendOTPFn(ctx, normalizedPhoneNumber, code, appID)
}
