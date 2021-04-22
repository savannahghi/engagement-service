package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *dummyResolver) ID(ctx context.Context, obj *resources.Dummy) (*string, error) {
	return nil, nil
}

func (r *mutationResolver) PhoneNumberVerificationCode(ctx context.Context, to string, code string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verificationCode, err := r.interactor.Whatsapp.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
	if err != nil {
		return false, fmt.Errorf("failed to send a verification code: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"phoneNumberVerificationCode",
		err,
	)

	return verificationCode, nil
}

func (r *mutationResolver) WellnessCardActivationDependant(ctx context.Context, to string, memberName string, cardName string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	wellnessCardActivationDependantMessage, err := r.interactor.Whatsapp.WellnessCardActivationDependant(
		ctx,
		to,
		memberName,
		cardName,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send wellness card activation messages to dependant via WhatsApp: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"wellnessCardActivationDependant",
		err,
	)

	return wellnessCardActivationDependantMessage, nil
}

func (r *mutationResolver) WellnessCardActivationPrincipal(ctx context.Context, to string, memberName string, cardName string, minorAgeThreshold string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	wellnessCardActivationPrincipal, err := r.interactor.Whatsapp.WellnessCardActivationPrincipal(
		ctx,
		to,
		memberName,
		cardName,
		minorAgeThreshold,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send wellness card activation messages to principal via WhatsApp: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"wellnessCardActivationPrincipal",
		err,
	)

	return wellnessCardActivationPrincipal, nil
}

func (r *mutationResolver) BillNotification(ctx context.Context, to string, productName string, billingPeriod string, billAmount string, paymentInstruction string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	billNotification, err := r.interactor.Whatsapp.BillNotification(
		ctx,
		to,
		productName,
		billingPeriod,
		billAmount,
		paymentInstruction,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send bill notification messages via WhatsApp: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"billNotification",
		err,
	)

	return billNotification, nil
}

func (r *mutationResolver) VirtualCards(ctx context.Context, to string, wellnessCardFamily string, virtualCardLink string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	virtualCardsNotification, err := r.interactor.Whatsapp.VirtualCards(
		ctx,
		to,
		wellnessCardFamily,
		virtualCardLink,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send virtual card setup notifications : %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"virtualCards",
		err,
	)

	return virtualCardsNotification, nil
}

func (r *mutationResolver) VisitStart(ctx context.Context, to string, memberName string, benefitName string, locationName string, startTime string, balance string, marketingMessage string) (bool, error) {
	beginTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	visitStart, err := r.interactor.Whatsapp.VisitStart(
		ctx,
		to,
		memberName,
		benefitName,
		locationName,
		startTime,
		balance,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send visit start SMS messages to members")
	}
	defer base.RecordGraphqlResolverMetrics(
		ctx,
		beginTime,
		"visitStart",
		err,
	)

	return visitStart, nil
}

func (r *mutationResolver) ClaimNotification(ctx context.Context, to string, claimReference string, claimTypeParenthesized string, provider string, visitType string, claimTime string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	claimNotification, err := r.interactor.Whatsapp.ClaimNotification(
		ctx,
		to,
		claimReference,
		claimTypeParenthesized,
		provider,
		visitType,
		claimTime,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send a claim notification message via WhatsApp")
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"claimNotification",
		err,
	)

	return claimNotification, nil
}

func (r *mutationResolver) PreauthApproval(ctx context.Context, to string, currency string, amount string, benefit string, provider string, member string, careContact string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	preauthApproval, err := r.interactor.Whatsapp.PreauthApproval(
		ctx,
		to,
		currency,
		amount,
		benefit,
		provider,
		member,
		careContact,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send a pre-authorization approval message via WhatsApp")
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"preauthApproval",
		err,
	)

	return preauthApproval, nil
}

func (r *mutationResolver) PreauthRequest(ctx context.Context, to string, currency string, amount string, benefit string, provider string, requestTime string, member string, careContact string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	preauthRequest, err := r.interactor.Whatsapp.PreauthRequest(
		ctx,
		to,
		currency,
		amount,
		benefit,
		provider,
		requestTime,
		member,
		careContact,
		marketingMessage,
	)

	if err != nil {
		return false, fmt.Errorf("failed to send a pre-authorization request message via WhatsApp")
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"preauthRequest",
		err,
	)

	return preauthRequest, nil
}

func (r *mutationResolver) SladeOtp(ctx context.Context, to string, name string, otp string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	sladeOTP, err := r.interactor.Whatsapp.SladeOTP(ctx, to, name, otp, marketingMessage)

	if err != nil {
		return false, fmt.Errorf("failed to send Slade ID OTP messages")
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sladeOTP",
		err,
	)

	return sladeOTP, nil
}

// Dummy returns generated.DummyResolver implementation.
func (r *Resolver) Dummy() generated.DummyResolver { return &dummyResolver{r} }

type dummyResolver struct{ *Resolver }
