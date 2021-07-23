// Package mail implements an email sending API that uses MailGun and SendInBlue.
package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/mail")

// Mail configuration constants
const (
	MailGunAPIKeyEnvVarName     = "MAILGUN_API_KEY"
	MailGunAPIBaseURLEnvVarName = "MAILGUN_API_BASE_URL"
	MailGunDomainEnvVarName     = "MAILGUN_DOMAIN"
	MailGunFromEnvVarName       = "MAILGUN_FROM"
	MailGunTimeoutSeconds       = 15

	SendInBlueAPIKeyEnvVarName  = "SEND_IN_BLUE_API_KEY"
	SendInBlueEnabledEnvVarName = "SEND_IN_BLUE_ENABLED"

	appName           = "Slade 360 HealthCloud"
	defaultUser       = "HealthCloud User"
	sendInBlueBaseURL = "https://api.sendinblue.com/v3/smtp/email"
)

// ServiceMail defines a Mail service interface
type ServiceMail interface {
	SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error)
	SendMailgun(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)
	SendEmail(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)
	SimpleEmail(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, error)
	SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error
	UpdateMailgunDeliveryStatus(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)
	GenerateEmailTemplate(name string, templateName string) string
}

// NewService initializes a new MailGun service
func NewService(repository repository.Repository) *Service {
	apiKey := serverutils.MustGetEnvVar(MailGunAPIKeyEnvVarName)
	domain := serverutils.MustGetEnvVar(MailGunDomainEnvVarName)
	from := serverutils.MustGetEnvVar(MailGunFromEnvVarName)

	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	// special case for sandbox
	if strings.Contains(domain, "sandbox") {
		mg.SetAPIBase(mailgun.APIBase)
	}

	sendInBlueEnabled := false
	if serverutils.MustGetEnvVar(SendInBlueEnabledEnvVarName) == "true" {
		sendInBlueEnabled = true
	}

	return &Service{
		Mg:                mg,
		From:              from,
		SendInBlueAPIKey:  serverutils.MustGetEnvVar(SendInBlueAPIKeyEnvVarName),
		SendInBlueEnabled: sendInBlueEnabled,
		Repository:        repository,
	}
}

// Service is an email sending service
type Service struct {
	Mg                *mailgun.MailgunImpl
	From              string
	SendInBlueEnabled bool
	SendInBlueAPIKey  string
	Repository        repository.Repository
}

// CheckPreconditions checks that all the required preconditions are satisfied
func (s Service) CheckPreconditions() {
	if s.Mg == nil {
		log.Panicf("uninitialized MailGun")
	}
	if s.From == "" {
		log.Panicf("uninitialized email from")
	}
	if s.SendInBlueAPIKey == "" {
		log.Panicf("uninitialized sendInBlueAPIKey")
	}
	if s.Repository == nil {
		log.Panicf("uninitialized repository in mail service")
	}
}

// MakeSendInBlueRequest makes a request to SendInBlue
func (s Service) MakeSendInBlueRequest(
	ctx context.Context,
	data map[string]interface{},
	target interface{},
) error {
	_, span := tracer.Start(ctx, "MakeSendInBlueRequest")
	defer span.End()

	bs, err := json.Marshal(data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("makeSendInBlueRequest: can't marshal data [%#v] to JSON: %w", data, err)
	}

	r := bytes.NewBuffer(bs)
	req, reqErr := http.NewRequest(http.MethodPost, sendInBlueBaseURL, r)
	if reqErr != nil {
		helpers.RecordSpanError(span, err)
		return reqErr
	}

	sendInBlueAPIKey := serverutils.MustGetEnvVar(SendInBlueAPIKeyEnvVarName)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", sendInBlueAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("SendInBlue API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("SendInBlue content error: %w", err)
	}

	if resp.StatusCode > 201 {
		return fmt.Errorf("SendInBlue API Error: %s", string(respBs))
	}

	err = json.Unmarshal(respBs, target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to unmarshal SendInBlue resp: %w", err)
	}

	return nil
}

// SendInBlue sends email via the SendInBlue service
func (s Service) SendInBlue(
	ctx context.Context,
	subject, text string,
	to ...string,
) (string, string, error) {
	s.CheckPreconditions()

	ctx, span := tracer.Start(ctx, "SendInBlue")
	defer span.End()

	sender := map[string]string{
		"name":  appName,
		"email": s.From,
	}

	addresses := []map[string]string{}
	for _, address := range to {
		addresses = append(addresses, map[string]string{
			"email": address,
		})
	}

	data := map[string]interface{}{
		"sender":      sender,
		"to":          addresses,
		"subject":     subject,
		"htmlContent": text,
	}

	result := map[string]interface{}{}

	err := s.MakeSendInBlueRequest(ctx, data, &result)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", "", fmt.Errorf("unable to send email via sendInBlue: %w", err)
	}

	messageID, ok := result["messageId"].(string)
	if !ok {
		return "", "", fmt.Errorf("string messageID not found in SendInBlue response %#v", result)
	}

	return "ok", messageID, nil
}

// SendMailgun sends email via MailGun
func (s Service) SendMailgun(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	s.CheckPreconditions()

	ctx, span := tracer.Start(ctx, "SendMailgun")
	defer span.End()

	message := s.Mg.NewMessage(s.From, subject, text, to...)
	if body != nil {
		message.SetHtml(*body)
	} else {
		message.SetHtml(text)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*MailGunTimeoutSeconds)
	defer cancel()

	resp, id, err := s.Mg.Send(ctx, message)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return resp, id, fmt.Errorf("mailgun email sending error: %s", err)
	}

	messageID := helpers.StripMailGunIDSpecialCharacters(id)

	outgoingEmail := &dto.OutgoingEmailsLog{
		UUID:        uuid.NewString(),
		To:          to,
		From:        MailGunFromEnvVarName,
		Subject:     subject,
		Text:        text,
		MessageID:   messageID,
		EmailSentOn: time.Now(),
	}

	err = s.SaveOutgoingEmails(ctx, outgoingEmail)
	if err != nil {
		return resp, id, fmt.Errorf("unable to save outgoing email(s): %s", err)
	}

	return resp, id, err
}

// SendEmail sends the specified email to the recipient(s) specified in `to`
// and returns the status
func (s Service) SendEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "SendEmail")
	defer span.End()

	if s.SendInBlueEnabled {
		return s.SendInBlue(ctx, subject, text, to...)
	}

	return s.SendMailgun(ctx, subject, text, body, to...)
}

// SimpleEmail is a simplified API to send email.
// It returns only a status or error.
func (s Service) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	s.CheckPreconditions()

	ctx, span := tracer.Start(ctx, "SimpleEmail")
	defer span.End()

	status, _, err := s.SendEmail(ctx, subject, text, nil, to...)

	return status, err
}

// SaveOutgoingEmails saves all the outgoing emails
func (s Service) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	ctx, span := tracer.Start(ctx, "SaveOutgoingEmails")
	defer span.End()

	return s.Repository.SaveOutgoingEmails(ctx, payload)
}

// UpdateMailgunDeliveryStatus updates the status and delivery time of the sent email message
func (s Service) UpdateMailgunDeliveryStatus(
	ctx context.Context,
	payload *dto.MailgunEvent,
) (*dto.OutgoingEmailsLog, error) {
	ctx, span := tracer.Start(ctx, "UpdateMailgunDeliveryStatus")
	defer span.End()

	return s.Repository.UpdateMailgunDeliveryStatus(ctx, payload)
}

//GenerateEmailTemplate generates custom emails to be sent to the users
func (s Service) GenerateEmailTemplate(name string, templateName string) string {
	t := template.Must(template.New("Be.WellEmailTemplate").Parse(templateName))
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, name)
	return buf.String()
}
