package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

// FakeServiceOTP is an interface that defines all interactions with the mock OTP service
type FakeServiceOTP struct {
	GenerateAndSendOTPFn   func(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmailFn       func(ctx context.Context, msisdn, email *string, appID *string) (string, error)
	SaveOTPToFirestoreFn   func(otp dto.OTP) error
	VerifyOtpFn            func(ctx context.Context, msisdn, verificationCode *string) (bool, error)
	VerifyEmailOtpFn       func(ctx context.Context, email, verificationCode *string) (bool, error)
	GenerateRetryOTPFn     func(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error)
	EmailVerificationOtpFn func(ctx context.Context, email *string) (string, error)
	GenerateOTPFn          func(ctx context.Context) (string, error)

	SendTemporaryPINFn func(ctx context.Context, input dto.TemporaryPIN) error
}

// GenerateAndSendOTP ...
func (f *FakeServiceOTP) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return f.GenerateAndSendOTPFn(ctx, msisdn, appID)
}

// SendOTPToEmail ...
func (f *FakeServiceOTP) SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error) {
	return f.SendOTPToEmailFn(ctx, msisdn, email, appID)
}

// SaveOTPToFirestore ...
func (f *FakeServiceOTP) SaveOTPToFirestore(otp dto.OTP) error {
	return f.SaveOTPToFirestoreFn(otp)
}

// VerifyOtp ...
func (f *FakeServiceOTP) VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error) {
	return f.VerifyOtpFn(ctx, msisdn, verificationCode)
}

// VerifyEmailOtp ...
func (f *FakeServiceOTP) VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error) {
	return f.VerifyEmailOtpFn(ctx, email, verificationCode)
}

// GenerateRetryOTP ...
func (f *FakeServiceOTP) GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error) {
	return f.GenerateRetryOTPFn(ctx, msisdn, retryStep, appID)
}

// EmailVerificationOtp ...
func (f *FakeServiceOTP) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	return f.EmailVerificationOtpFn(ctx, email)
}

// GenerateOTP ...
func (f *FakeServiceOTP) GenerateOTP(ctx context.Context) (string, error) {
	return f.GenerateOTPFn(ctx)
}

// SendTemporaryPIN ...
func (f *FakeServiceOTP) SendTemporaryPIN(ctx context.Context, input dto.TemporaryPIN) error {
	return f.SendTemporaryPINFn(ctx, input)
}
