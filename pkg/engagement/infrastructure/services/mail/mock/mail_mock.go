package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

//FakeMailService ...
type FakeMailService struct {
	SendInBlueFn func(ctx context.Context, subject, text string, to ...string) (string, string, error)

	SendMailgunFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)

	SendEmailFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)

	SimpleEmailFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, error)

	SaveOutgoingEmailsFn func(ctx context.Context, payload *dto.OutgoingEmailsLog) error

	UpdateMailgunDeliveryStatusFn func(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)

	GenerateEmailTemplateFn func(name string, templateName string) string
}

//SendInBlue mock function
func (f *FakeMailService) SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error) {
	return f.SendInBlueFn(ctx, subject, text, to...)
}

//SendMailgun mock function
func (f *FakeMailService) SendMailgun(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendMailgunFn(ctx, subject, text, body, to...)
}

//SendEmail mock function
func (f *FakeMailService) SendEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendEmailFn(ctx, subject, text, body, to...)
}

//SimpleEmail mock function
func (f *FakeMailService) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return f.SimpleEmailFn(ctx, subject, text, body, to...)
}

//SaveOutgoingEmails mock function
func (f *FakeMailService) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return f.SaveOutgoingEmailsFn(ctx, payload)
}

//UpdateMailgunDeliveryStatus mock function
func (f *FakeMailService) UpdateMailgunDeliveryStatus(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error) {
	return f.UpdateMailgunDeliveryStatusFn(ctx, payload)
}

//GenerateEmailTemplate mock function
func (f *FakeMailService) GenerateEmailTemplate(name string, templateName string) string {
	return f.GenerateEmailTemplateFn(name, templateName)
}
