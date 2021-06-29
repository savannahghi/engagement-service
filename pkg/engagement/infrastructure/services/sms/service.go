package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

// AIT environment variables names
const (
	APIKeyEnvVarName       = "AIT_API_KEY"
	APIUsernameEnvVarName  = "AIT_USERNAME"
	APISenderIDEnvVarName  = "AIT_SENDER_ID"
	AITEnvVarName          = "AIT_ENVIRONMENT"
	BeWellAITAPIKey        = "AIT_BEWELL_API_KEY"
	BeWellAITUsername      = "AIT_BEWELL_USERNAME"
	BeWellAITSenderID      = "AIT_BEWELL_SENDER_ID"
	AITAuthenticationError = "The supplied authentication is invalid"

	//AITCallbackCollectionName is the name of a Cloud Firestore collection into which AIT
	// callback data will be saved for future analysis
	AITCallbackCollectionName = "ait_callbacks"
)

// ServiceSMS defines the interactions with sms service
type ServiceSMS interface {
	SendToMany(
		message string,
		to []string,
		from base.SenderID,
	) (*dto.SendMessageResponse, error)
	Send(
		to, message string,
		from base.SenderID,
	) (*dto.SendMessageResponse, error)
	SendMarketingSMS(
		to []string,
		message string,
		from base.SenderID,
	) (*dto.SendMessageResponse, error)
	SaveMarketingMessage(
		ctx context.Context,
		data dto.MarketingSMS,
	) error
	UpdateMarketingMessage(
		ctx context.Context,
		phoneNumber string,
		deliveryReport *dto.ATDeliveryReport,
	) (*dto.MarketingSMS, error)
}

// Service defines a sms service struct
type Service struct {
	Env        string
	Repository repository.Repository
	Crm        hubspot.ServiceHubSpotInterface
}

// GetSmsURL is the sms endpoint
func GetSmsURL(env string) string {
	return GetAPIHost(env) + "/version1/messaging"
}

// GetAPIHost returns either sandbox or prod
func GetAPIHost(env string) string {
	return getHost(env, "api")
}

func getHost(env, service string) string {
	if env != "sandbox" {
		return fmt.Sprintf("https://%s.africastalking.com", service)
	}
	return fmt.Sprintf(
		"https://%s.sandbox.africastalking.com",
		service,
	)

}

// NewService returns a new service
func NewService(repository repository.Repository, crm hubspot.ServiceHubSpotInterface) *Service {
	env := base.MustGetEnvVar(AITEnvVarName)
	return &Service{env, repository, crm}
}

// SaveMarketingMessage saves the callback data for future analysis
func (s Service) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) error {
	return s.Repository.SaveMarketingMessage(ctx, data)
}

// UpdateMessageSentStatus updates the message sent field to true when a message
// is sent to a user
func (s Service) UpdateMessageSentStatus(
	ctx context.Context,
	phonenumber string,
) error {
	return s.Repository.UpdateMessageSentStatus(ctx, phonenumber)
}

// UpdateMarketingMessage adds a delivery reposrt to an AIT SMS
func (s Service) UpdateMarketingMessage(
	ctx context.Context,
	phoneNumber string,
	deliveryReport *dto.ATDeliveryReport,
) (*dto.MarketingSMS, error) {
	return s.Repository.UpdateMarketingMessage(ctx, phoneNumber, deliveryReport)
}

// SendToMany is a utility method to send to many recipients at the same time
func (s Service) SendToMany(
	message string,
	to []string,
	from base.SenderID,
) (*dto.SendMessageResponse, error) {
	recipients := strings.Join(to, ",")
	return s.Send(recipients, message, from)
}

// Send is a method used to send to a single recipient
func (s Service) Send(
	to, message string,
	from base.SenderID,
) (*dto.SendMessageResponse, error) {
	switch from {
	case base.SenderIDSLADE360:
		return s.SendSMS(
			to,
			message,
			base.MustGetEnvVar(APISenderIDEnvVarName),
			base.MustGetEnvVar(APIUsernameEnvVarName),
			base.MustGetEnvVar(APIKeyEnvVarName),
		)
	case base.SenderIDBewell:
		return s.SendSMS(
			to,
			message,
			base.MustGetEnvVar(BeWellAITSenderID),
			base.MustGetEnvVar(BeWellAITUsername),
			base.MustGetEnvVar(BeWellAITAPIKey),
		)
	}
	return nil, fmt.Errorf("unknown AIT sender")
}

// SendSMS is a method used to send SMS
func (s Service) SendSMS(
	to, message string,
	from string,
	username string,
	key string,
) (*dto.SendMessageResponse, error) {
	values := url.Values{}
	values.Set("username", username)
	values.Set("to", to)
	values.Set("message", message)
	values.Set("from", from)

	smsURL := GetSmsURL(s.Env)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, err := s.newPostRequest(smsURL, values, headers, key)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	// read the response body to a variable
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	smsMessageResponse := &dto.SendMessageResponse{}

	bodyAsString := string(bodyBytes)
	if strings.Contains(bodyAsString, AITAuthenticationError) {
		// return so that other processes don't break
		log.Println("AIT Authentication error encountered")
		return smsMessageResponse, nil
	}

	// reset the response body to the original unread state so that decode can
	// continue
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(res.Body).Decode(smsMessageResponse); err != nil {
		return nil, errors.New(
			fmt.Errorf("SMS Error : unable to parse sms response ; %v", err).
				Error(),
		)
	}

	return smsMessageResponse, nil
}

func (s Service) newPostRequest(
	url string,
	values url.Values,
	headers map[string]string,
	key string,
) (*http.Response, error) {
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Length", strconv.Itoa(reader.Len()))
	req.Header.Set("apikey", key)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// SendMarketingSMS is a method to send marketing bulk SMS for Be.Well launch/campaigns.
// It interacts with our DB to save the message and update our CRM with engagements
func (s Service) SendMarketingSMS(
	to []string,
	message string,
	from base.SenderID,
) (*dto.SendMessageResponse, error) {
	recipients := strings.Join(to, ",")
	resp, err := s.Send(recipients, message, from)
	if err != nil {
		return nil, err
	}

	smsMsgData := resp.SMSMessageData
	if smsMsgData == nil {
		return nil, nil
	}

	engagement := domain.Engagement{
		Active:    true,
		Type:      "NOTE",
		Timestamp: time.Now().UnixNano() / 1000000,
	}
	engagementData := domain.EngagementData{
		Engagement: engagement,
		Metadata: map[string]interface{}{
			"body": message,
		},
	}

	smsMsgDataRecipients := smsMsgData.Recipients
	for _, recipient := range smsMsgDataRecipients {
		phone := recipient.Number
		data := dto.MarketingSMS{
			ID:                   uuid.New().String(),
			PhoneNumber:          phone,
			SenderID:             from,
			MessageSentTimeStamp: time.Now(),
			Message:              message,
			Status:               recipient.Status,
			Engagement:           engagementData,
		}

		// todo make this async
		resp, err := s.Crm.CreateEngagementByPhone(phone, engagementData)
		if err != nil {
			log.Print(err)
		}

		if resp != nil {
			data.IsSynced = true
		}

		if err := s.SaveMarketingMessage(
			context.Background(),
			data,
		); err != nil {
			return nil, err
		}

		// Toggle message sent value to TRUE
		if err := s.UpdateMessageSentStatus(
			context.Background(),
			data.PhoneNumber,
		); err != nil {
			return nil, err
		}

		// Sleep for 5 seconds to reduce the rate at which we call HubSpot's APIs
		// They have a rate limit of 100/10s
		time.Sleep(5 * time.Second)
	}

	return resp, nil
}
