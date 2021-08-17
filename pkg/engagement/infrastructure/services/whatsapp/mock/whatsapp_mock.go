package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

// FakeServiceWhatsapp defines the interactions with the whatsapp service
type FakeServiceWhatsapp struct {
	PhoneNumberVerificationCodeFn func(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationDependantFn func(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationPrincipalFn func(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		minorAgeThreshold string,
		marketingMessage string,
	) (bool, error)

	BillNotificationFn func(
		ctx context.Context,
		to string,
		productName string,
		billingPeriod string,
		billAmount string,
		paymentInstruction string,
		marketingMessage string,
	) (bool, error)

	VirtualCardsFn func(
		ctx context.Context,
		to string,
		wellnessCardFamily string,
		virtualCardLink string,
		marketingMessage string,
	) (bool, error)

	VisitStartFn func(
		ctx context.Context,
		to string,
		memberName string,
		benefitName string,
		locationName string,
		startTime string,
		balance string,
		marketingMessage string,
	) (bool, error)

	ClaimNotificationFn func(
		ctx context.Context,
		to string,
		claimReference string,
		claimTypeParenthesized string,
		provider string,
		visitType string,
		claimTime string,
		marketingMessage string,
	) (bool, error)

	PreauthApprovalFn func(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	PreauthRequestFn func(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		requestTime string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	SladeOTPFn func(
		ctx context.Context,
		to string,
		name string,
		otp string,
		marketingMessage string,
	) (bool, error)

	SaveTwilioCallbackResponseFn func(
		ctx context.Context,
		data dto.Message,
	) error
}

// PhoneNumberVerificationCode ...
func (f *FakeServiceWhatsapp) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return f.PhoneNumberVerificationCodeFn(ctx, to, code, marketingMessage)
}

// WellnessCardActivationDependant ...
func (f *FakeServiceWhatsapp) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	return f.WellnessCardActivationDependantFn(ctx, to, memberName, cardName, marketingMessage)
}

// WellnessCardActivationPrincipal ...
func (f *FakeServiceWhatsapp) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	return f.WellnessCardActivationPrincipalFn(ctx, to, memberName, cardName, minorAgeThreshold, marketingMessage)
}

// BillNotification ...
func (f *FakeServiceWhatsapp) BillNotification(
	ctx context.Context,
	to string,
	productName string,
	billingPeriod string,
	billAmount string,
	paymentInstruction string,
	marketingMessage string,
) (bool, error) {
	return f.BillNotificationFn(
		ctx,
		to,
		productName,
		billingPeriod,
		billAmount,
		paymentInstruction,
		marketingMessage,
	)
}

// VirtualCards ...
func (f *FakeServiceWhatsapp) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	return f.VirtualCardsFn(ctx, to, wellnessCardFamily, virtualCardLink, marketingMessage)
}

// VisitStart ...
func (f *FakeServiceWhatsapp) VisitStart(
	ctx context.Context,
	to string,
	memberName string,
	benefitName string,
	locationName string,
	startTime string,
	balance string,
	marketingMessage string,
) (bool, error) {
	return f.VisitStartFn(
		ctx,
		to,
		memberName,
		benefitName,
		locationName,
		startTime,
		balance,
		marketingMessage,
	)
}

// ClaimNotification ...
func (f *FakeServiceWhatsapp) ClaimNotification(
	ctx context.Context,
	to string,
	claimReference string,
	claimTypeParenthesized string,
	provider string,
	visitType string,
	claimTime string,
	marketingMessage string,
) (bool, error) {
	return f.ClaimNotificationFn(
		ctx,
		to,
		claimReference,
		claimTypeParenthesized,
		provider,
		visitType,
		claimTime,
		marketingMessage,
	)
}

// PreauthApproval ...
func (f *FakeServiceWhatsapp) PreauthApproval(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return f.PreauthApprovalFn(
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
}

// PreauthRequest ...
func (f *FakeServiceWhatsapp) PreauthRequest(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	requestTime string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return f.PreauthRequestFn(
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
}

// SladeOTP ...
func (f *FakeServiceWhatsapp) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	return f.SladeOTPFn(ctx, to, name, otp, marketingMessage)
}

// SaveTwilioCallbackResponse ...
func (f *FakeServiceWhatsapp) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return f.SaveTwilioCallbackResponseFn(ctx, data)
}
