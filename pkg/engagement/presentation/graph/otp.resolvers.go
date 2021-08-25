package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/serverutils"
)

func (r *dummyResolver) ID(ctx context.Context, obj *dto.Dummy) (*string, error) {
	return nil, nil
}

func (r *mutationResolver) VerifyOtp(ctx context.Context, msisdn string, otp string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verifyOTP, err := r.interactor.OTP.VerifyOtp(ctx, &msisdn, &otp)
	if err != nil {
		return false, fmt.Errorf("failed to check for the validity of the supplied OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"verifyOTP",
		err,
	)

	return verifyOTP, nil
}

func (r *mutationResolver) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verifyEmailOTP, err := r.interactor.OTP.VerifyEmailOtp(ctx, &email, &otp)
	if err != nil {
		return false, fmt.Errorf("failed to check for the validity of the supplied OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"verifyEmailOTP",
		err,
	)

	return verifyEmailOTP, nil
}

func (r *queryResolver) GenerateOtp(ctx context.Context, msisdn string, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.interactor.OTP.GenerateAndSendOTP(ctx, msisdn, appID)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) GenerateAndEmailOtp(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.interactor.OTP.SendOTPToEmail(ctx, &msisdn, email, appID)
	if err != nil {
		return "", fmt.Errorf("failed to send the generated OTP to the provided email address")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateAndEmailOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) GenerateRetryOtp(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.interactor.OTP.GenerateRetryOTP(ctx, &msisdn, retryStep, appID)
	if err != nil {
		return "", fmt.Errorf("failed to generate fallback OTPs")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateRetryOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.interactor.OTP.EmailVerificationOtp(ctx, &email)
	if err != nil {
		return "", fmt.Errorf("failed to generate an OTP to the supplied email for verification")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"emailVerificationOTP",
		err,
	)

	return otp, nil
}

// Dummy returns generated.DummyResolver implementation.
func (r *Resolver) Dummy() generated.DummyResolver { return &dummyResolver{r} }

type dummyResolver struct{ *Resolver }
