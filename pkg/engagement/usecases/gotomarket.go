package usecases

import (
	"context"
	"fmt"

	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
)

// GoToMarketUseCases represents all the marketing data business logic
type GoToMarketUseCases interface {
	CollectEmails(ctx context.Context, email string, phonenumber string) (*hubspotDomain.CRMContact, error)
	BeWellAware(ctx context.Context, email string) (*hubspotDomain.CRMContact, error)
}

// GoToMarketImpl represents the marketing usecase implementation
type GoToMarketImpl struct {
	hubspot hubspotUsecases.HubSpotUsecases
	mail    mail.ServiceMail
}

// NewGoToMarketUsecases initialises a marketing usecase
func NewGoToMarketUsecases(
	hubspot hubspotUsecases.HubSpotUsecases,
	mail mail.ServiceMail,
) *GoToMarketImpl {
	return &GoToMarketImpl{
		hubspot: hubspot,
		mail:    mail,
	}
}

// CollectEmails receives an email and update the email of the found record in our firestore and hubspot CRM
func (g *GoToMarketImpl) CollectEmails(ctx context.Context, email string, phonenumber string) (*hubspotDomain.CRMContact, error) {
	ctx, span := tracer.Start(ctx, "CollectEmails")
	defer span.End()

	contact, err := g.hubspot.GetContactByPhone(ctx, phonenumber)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to get contact with phone number %s: %w", phonenumber, err)
	}
	if contact == nil {
		return nil, fmt.Errorf("contact with phonenumber %s not found", phonenumber)
	}
	contact.Properties.Email = email

	updatedContact, err := g.hubspot.UpdateHubSpotContact(ctx, contact)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to update contact with phone number %s: %w", phonenumber, err)
	}

	name := "Kevin From Be.Well"
	body := g.mail.GenerateEmailTemplate(name, mail.MarketingEmailTemplate)
	subject := "Download the new Be.Well app to manage your insurance benefits"
	_, _, err = g.mail.SendEmail(
		ctx,
		subject,
		name,
		&body,
		email,
	)
	if err != nil {
		return contact, fmt.Errorf("failed to send welcome email to %s: %w", email, err)
	}

	return updatedContact, nil
}

//BeWellAware toggles the user identified by the provided email as bewell-aware
func (g *GoToMarketImpl) BeWellAware(ctx context.Context, email string) (*hubspotDomain.CRMContact, error) {
	ctx, span := tracer.Start(ctx, "BeWellAware")
	defer span.End()

	contact, err := g.hubspot.GetContactByEmail(ctx, email)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to get contact with email %s: %w", email, err)
	}
	if contact == nil {
		return nil, fmt.Errorf("contact with email %s not found", email)
	}
	contact.Properties.BeWellAware = hubspotDomain.GeneralOptionTypeYes

	updatedContact, err := g.hubspot.UpdateHubSpotContact(ctx, contact)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to update contact with email %s: %w", email, err)
	}
	return updatedContact, nil
}
