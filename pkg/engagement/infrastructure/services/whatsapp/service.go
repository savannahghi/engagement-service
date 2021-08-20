package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement/pkg/engagement/repository"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement/pkg/engagement/services/whatsapp")

// Twilio Whatsapp API contants
const (
	// TwilioHTTPClientTimeoutSeconds determines how long to wait (in seconds) before giving up on a
	// request to the Twilio API
	TwilioHTTPClientTimeoutSeconds = 30
	TwilioWhatsappSIDEnvVarName    = "TWILIO_WHATSAPP_SID"

	// gosec false positive
	TwilioWhatsappAuthTokenEnvVarName = "TWILIO_WHATSAPP_AUTH_TOKEN" /* #nosec */

	TwilioWhatsappSenderEnvVarName = "TWILIO_WHATSAPP_SENDER"

	twilioWhatsappBaseURL = "https://api.twilio.com/2010-04-01/Accounts/"
)

// NewService initializes a properly set up WhatsApp service
func NewService(r repository.Repository, crm crm.ServiceCrm) *Service {
	sid := serverutils.MustGetEnvVar(TwilioWhatsappSIDEnvVarName)
	authToken := serverutils.MustGetEnvVar(TwilioWhatsappAuthTokenEnvVarName)
	sender := serverutils.MustGetEnvVar(TwilioWhatsappSenderEnvVarName)
	httpClient := &http.Client{
		Timeout: time.Second * TwilioHTTPClientTimeoutSeconds,
	}
	return &Service{
		BaseURL:          twilioWhatsappBaseURL,
		AccountSID:       sid,
		AccountAuthToken: authToken,
		Sender:           sender,
		HTTPClient:       httpClient,
		Repository:       r,
		Crm:              crm,
	}
}

// ServiceWhatsapp defines the interactions with the whatsapp service
type ServiceWhatsapp interface {
	PhoneNumberVerificationCode(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	ReceiveInboundMessages(
		ctx context.Context,
		message *dto.TwilioMessage,
	) (*dto.TwilioMessage, error)

	WellnessCardActivationDependant(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationPrincipal(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		minorAgeThreshold string,
		marketingMessage string,
	) (bool, error)

	BillNotification(
		ctx context.Context,
		to string,
		productName string,
		billingPeriod string,
		billAmount string,
		paymentInstruction string,
		marketingMessage string,
	) (bool, error)

	VirtualCards(
		ctx context.Context,
		to string,
		wellnessCardFamily string,
		virtualCardLink string,
		marketingMessage string,
	) (bool, error)

	VisitStart(
		ctx context.Context,
		to string,
		memberName string,
		benefitName string,
		locationName string,
		startTime string,
		balance string,
		marketingMessage string,
	) (bool, error)

	ClaimNotification(
		ctx context.Context,
		to string,
		claimReference string,
		claimTypeParenthesized string,
		provider string,
		visitType string,
		claimTime string,
		marketingMessage string,
	) (bool, error)

	PreauthApproval(
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

	PreauthRequest(
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

	SladeOTP(
		ctx context.Context,
		to string,
		name string,
		otp string,
		marketingMessage string,
	) (bool, error)

	SaveTwilioCallbackResponse(
		ctx context.Context,
		data dto.Message,
	) error
}

// Service is a WhatsApp service. The receivers implement the query and mutation resolvers.
type Service struct {
	BaseURL          string
	AccountSID       string
	AccountAuthToken string
	Sender           string
	HTTPClient       *http.Client
	Repository       repository.Repository
	Crm              crm.ServiceCrm
}

// CheckPreconditions ...
func (s Service) CheckPreconditions() {
	if s.HTTPClient == nil {
		log.Panicf("nil http client in Twilio WhatsApp service")
	}

	if s.BaseURL == "" {
		log.Panicf("blank base URL in Twilio WhatsApp service")
	}

	if s.AccountSID == "" {
		log.Panicf("blank accountSID in Twilio WhatsApp service")
	}

	if s.AccountAuthToken == "" {
		log.Panicf("blank account auth token in Twilio WhatsApp service")
	}

	if s.Sender == "" {
		log.Panicf("blank sender in Twilio WhatsApp service")
	}
}

// MakeTwilioRequest makes a twilio request
func (s Service) MakeTwilioRequest(
	ctx context.Context,
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	_, span := tracer.Start(ctx, "MakeTwilioRequest")
	defer span.End()
	s.CheckPreconditions()

	if serverutils.IsDebug() {
		log.Printf("Twilio request data: \n%s\n", content)
	}

	r := strings.NewReader(content.Encode())
	req, err := http.NewRequest(method, s.BaseURL+urlPath, r)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.AccountSID, s.AccountAuthToken)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio room content error: %w", err)
	}

	if resp.StatusCode > 201 {
		return fmt.Errorf("twilio API Error: %s", string(respBs))
	}

	if serverutils.IsDebug() {
		log.Printf("Twilio response: \n%s\n", string(respBs))
	}
	err = json.Unmarshal(respBs, target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to unmarshal Twilio resp: %w", err)
	}

	return nil
}

// PhoneNumberVerificationCode sends Phone Number verification codes via WhatsApp
func (s Service) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "PhoneNumberVerificationCode")
	defer span.End()
	s.CheckPreconditions()

	normalizedPhoneNo, err := converterandformatter.NormalizeMSISDN(to)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("%s is not a valid E164 phone number: %w", to, err)
	}

	msgFrom := fmt.Sprintf("whatsapp:%s", s.Sender)
	msgTo := fmt.Sprintf("whatsapp:%s", *normalizedPhoneNo)
	msg := fmt.Sprintf("Your phone number verification code is %s", code)

	payload := url.Values{}
	payload.Add("From", msgFrom)
	payload.Add("Body", msg)
	payload.Add("To", msgTo)

	target := dto.Message{}
	path := fmt.Sprintf("%s/Messages.json", s.AccountSID)
	err = s.MakeTwilioRequest(
		ctx,
		http.MethodPost,
		path,
		payload,
		&target,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("error from Twilio: %w", err)
	}

	// save Twilio response for audit purposes
	_, _, err = firebasetools.CreateNode(ctx, &target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to save Twilio response: %w", err)
	}
	// TODO Find out why /ide is not working (401s)
	// TODO deploy UAT, deploy prod, tag (semver)
	return true, nil
}

// ReceiveInboundMessages save Twilio whatsapp's inbound messsage to firestore and hubspot
func (s Service) ReceiveInboundMessages(
	ctx context.Context,
	message *dto.TwilioMessage,
) (*dto.TwilioMessage, error) {
	if message == nil {
		return nil, fmt.Errorf("nil inbound message")
	}

	inboundMsg, err := s.Repository.SaveInboundWAMessages(ctx, *message)
	if err != nil {
		return nil, fmt.Errorf("failed to save inbound whatsapp message: %w", err)
	}

	phoneNumber := strings.ReplaceAll(message.From, "whatsapp:", "")
	contact, err := s.Crm.GetContactByPhone(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts by phone %s: %w", phoneNumber, err)
	}

	if contact.Properties.FirstChannelOfContact == "" {
		contact.Properties.FirstChannelOfContact = hubspotDomain.ChannelOfContact("WHATSAPP")
		_, err := s.Crm.UpdateHubSpotContact(ctx, contact)
		if err != nil {
			return nil, fmt.Errorf("failed to update contact with First channel of contact as WHATSAPP: %w", err)
		}
	}

	searchResp, err := s.Crm.SearchContactByPhone(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to search for phone %s: %w", phoneNumber, err)
	}
	contactID, err := strconv.Atoi(searchResp.Results[0].ContactID)
	if err != nil {
		return nil, fmt.Errorf("unable to convert contact ID to int")
	}

	engagementData := &hubspotDomain.EngagementData{
		Engagement: hubspotDomain.Engagement{
			Active:    true,
			Type:      hubspotDomain.EngagementTypeNOTE,
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		},
		Metadata: map[string]interface{}{
			"body": message.Body,
		},
		Associations: hubspotDomain.Associations{
			ContactIDs: []int{contactID},
		},
	}
	engagement, err := s.Crm.CreateHubspotEngagement(ctx, engagementData)
	if err != nil {
		return nil, fmt.Errorf("failed to create engagement for contact %s: %w", phoneNumber, err)
	}
	if engagement == nil {
		return nil, fmt.Errorf("nil engagement returned")
	}

	// add to the conversation inbox
	return inboundMsg, nil
}

// WellnessCardActivationDependant sends wellness card activation messages via WhatsApp
func (s Service) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Create a common path for Twilio messages
	// TODO Implement wellness card activation message
	return false, nil
}

// WellnessCardActivationPrincipal sends wellness card activation messages to principals via WhatsApp
func (s Service) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Implement wellness card activation message for principals
	return false, nil
}

// BillNotification sends bill notification messages via WhatsApp
func (s Service) BillNotification(
	ctx context.Context,
	to string,
	productName string,
	billingPeriod string,
	billAmount string,
	paymentInstruction string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Implement bill notification message
	return false, nil
}

// VirtualCards sends virtual card setup notifications
func (s Service) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Implement virtual card notification message
	return false, nil
}

// VisitStart sends visit start SMS messages to members
func (s Service) VisitStart(
	ctx context.Context,
	to string,
	memberName string,
	benefitName string,
	locationName string,
	startTime string,
	balance string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Implement virtual card notification message
	return false, nil
}

// ClaimNotification sends a claim notification message via WhatsApp
func (s Service) ClaimNotification(
	ctx context.Context,
	to string,
	claimReference string,
	claimTypeParenthesized string,
	provider string,
	visitType string,
	claimTime string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()
	// TODO Implement claim notification message
	return false, nil
}

// PreauthApproval sends a pre-authorization approval message via WhatsApp
func (s Service) PreauthApproval(
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
	s.CheckPreconditions()
	// TODO Implement preauth approval message
	return false, nil
}

// PreauthRequest sends a pre-authorization request message via WhatsApp
func (s Service) PreauthRequest(
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
	s.CheckPreconditions()
	// TODO Implement preauth request message
	return false, nil
}

// SladeOTP sends Slade ID OTP messages
func (s Service) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	s.CheckPreconditions()

	// TODO Implement Slade OTP
	return false, nil
}

// SaveTwilioCallbackResponse saves the twilio callback response for future
// analysis
func (s Service) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return s.Repository.SaveTwilioResponse(ctx, data)
}
