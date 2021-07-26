package crm

import (
	"context"
	"fmt"

	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/usecases")

// ServiceCrm represents commontools crm lib usecases extension
type ServiceCrm interface {
	IsOptedOut(ctx context.Context, phoneNumber string) (bool, error)
	CollectEmails(ctx context.Context, email string, phonenumber string) (*hubspotDomain.CRMContact, error)
	BeWellAware(ctx context.Context, email string) (*hubspotDomain.CRMContact, error)
}

// Hubspot interacts with `HubSpot` CRM usecases
type Hubspot struct {
	hubSpotUsecases hubspotUsecases.HubSpotUsecases
	mail            mail.ServiceMail
}

// NewCrmService inits a new crm instance
func NewCrmService(hubSpotUsecases hubspotUsecases.HubSpotUsecases, mail mail.ServiceMail) *Hubspot {
	return &Hubspot{
		hubSpotUsecases: hubSpotUsecases,
		mail:            mail,
	}
}

// IsOptedOut checks if a given phone number is opted out
func (h *Hubspot) IsOptedOut(ctx context.Context, phoneNumber string) (bool, error) {
	return h.hubSpotUsecases.IsOptedOut(ctx, phoneNumber)
}

// CollectEmails receives an email and update the email of the found record in our firestore and hubspot CRM
func (h *Hubspot) CollectEmails(ctx context.Context, email string, phonenumber string) (*hubspotDomain.CRMContact, error) {
	ctx, span := tracer.Start(ctx, "CollectEmails")
	defer span.End()

	contact, err := h.hubSpotUsecases.GetContactByPhone(ctx, phonenumber)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to get contact with phone number %s: %w", phonenumber, err)
	}
	if contact == nil {
		return nil, fmt.Errorf("contact with phonenumber %s not found", phonenumber)
	}
	contact.Properties.Email = email

	updatedContact, err := h.hubSpotUsecases.UpdateHubSpotContact(ctx, contact)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to update contact with phone number %s: %w", phonenumber, err)
	}

	name := "Kevin From Be.Well"
	body := h.mail.GenerateEmailTemplate(name, mail.MarketingEmailTemplate)
	subject := "Download the new Be.Well app to manage your insurance benefits"
	_, _, err = h.mail.SendEmail(
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
func (h *Hubspot) BeWellAware(ctx context.Context, email string) (*hubspotDomain.CRMContact, error) {
	ctx, span := tracer.Start(ctx, "BeWellAware")
	defer span.End()

	contact, err := h.hubSpotUsecases.GetContactByEmail(ctx, email)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to get contact with email %s: %w", email, err)
	}
	if contact == nil {
		return nil, fmt.Errorf("contact with email %s not found", email)
	}
	contact.Properties.BeWellAware = hubspotDomain.GeneralOptionTypeYes

	updatedContact, err := h.hubSpotUsecases.UpdateHubSpotContact(ctx, contact)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to update contact with email %s: %w", email, err)
	}
	return updatedContact, nil
}
