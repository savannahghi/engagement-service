package mock

import (
	"context"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
)

// FakeServiceMail defines a mock Mail service interface
type FakeServiceMail struct {
	SendInBlueFn  func(ctx context.Context, subject, text string, to ...string) (string, string, error)
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
	SaveOutgoingEmailsFn          func(ctx context.Context, payload *dto.OutgoingEmailsLog) error
	UpdateMailgunDeliveryStatusFn func(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)
	GenerateEmailTemplateFn func(name string, templateName string) string
}

// SendInBlue ...
func (f *FakeServiceMail) SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error) {
	return f.SendInBlueFn(ctx, subject, text, to...)
}

// SendMailgun ...
func (f *FakeServiceMail) SendMailgun(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendMailgunFn(ctx, subject, text, body, to...)
}

// SendEmail ...
func (f *FakeServiceMail) SendEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendEmailFn(ctx, subject, text, body, to...)
}

// SimpleEmail ...
func (f *FakeServiceMail) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return f.SimpleEmailFn(ctx, subject, text, body, to...)
}

// SaveOutgoingEmails ...
func (f *FakeServiceMail) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return f.SaveOutgoingEmailsFn(ctx, payload)
}

// UpdateMailgunDeliveryStatus ...
func (f *FakeServiceMail) UpdateMailgunDeliveryStatus(
	ctx context.Context,
	payload *dto.MailgunEvent,
) (*dto.OutgoingEmailsLog, error) {
	return f.UpdateMailgunDeliveryStatusFn(ctx, payload)
}

// GenerateEmailTemplate ...
func (f *FakeServiceMail) GenerateEmailTemplate(name string, templateName string) string {
	return f.GenerateEmailTemplateFn(name, templateName)
}
