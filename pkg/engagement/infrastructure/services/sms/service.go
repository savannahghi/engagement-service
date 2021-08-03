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
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/crm"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/edi"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/sms")

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
		ctx context.Context,
		message string,
		to []string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)
	Send(
		ctx context.Context,
		to, message string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)
	SendMarketingSMS(
		ctx context.Context,
		to []string,
		message string,
		from enumutils.SenderID,
		segment string,
	) (*dto.SendMessageResponse, error)
	SaveMarketingMessage(
		ctx context.Context,
		data dto.MarketingSMS,
	) (*dto.MarketingSMS, error)
	UpdateMarketingMessage(
		ctx context.Context,
		data *dto.MarketingSMS,
	) (*dto.MarketingSMS, error)
	GetMarketingSMSByPhone(
		ctx context.Context,
		phoneNumber string,
	) (*dto.MarketingSMS, error)
}

// Service defines a sms service struct
type Service struct {
	Env        string
	Repository repository.Repository
	Crm        crm.ServiceCrm
	PubSub     messaging.NotificationService
	Edi        edi.ServiceEdi
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
func NewService(
	repository repository.Repository,
	crm crm.ServiceCrm,
	pubsub messaging.NotificationService,
	edi edi.ServiceEdi,
) *Service {
	env := serverutils.MustGetEnvVar(AITEnvVarName)
	return &Service{env, repository, crm, pubsub, edi}
}

// SaveMarketingMessage saves the callback data for future analysis
func (s Service) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return s.Repository.SaveMarketingMessage(ctx, data)
}

// UpdateMarketingMessage adds a delivery report to an AIT SMS
func (s Service) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return s.Repository.UpdateMarketingMessage(ctx, data)
}

// GetMarketingSMSByPhone returns the latest message given a phone number
func (s Service) GetMarketingSMSByPhone(
	ctx context.Context,
	phoneNumber string,
) (*dto.MarketingSMS, error) {
	return s.Repository.GetMarketingSMSByID(ctx, phoneNumber)
}

// UpdateMessageSentStatus updates the message sent field to true when a message
// is sent to a user
func (s Service) UpdateMessageSentStatus(
	ctx context.Context,
	phonenumber string,
	segment string,
) error {
	return s.Repository.UpdateMessageSentStatus(ctx, phonenumber, segment)
}

// SendToMany is a utility method to send to many recipients at the same time
func (s Service) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	recipients := strings.Join(to, ",")
	return s.Send(ctx, recipients, message, from)
}

// Send is a method used to send to a single recipient
func (s Service) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {

	switch from {
	case enumutils.SenderIDSLADE360:
		return s.SendSMS(
			ctx,
			to,
			message,
			serverutils.MustGetEnvVar(APISenderIDEnvVarName),
			serverutils.MustGetEnvVar(APIUsernameEnvVarName),
			serverutils.MustGetEnvVar(APIKeyEnvVarName),
		)
	case enumutils.SenderIDBewell:
		return s.SendSMS(
			ctx,
			to,
			message,
			serverutils.MustGetEnvVar(BeWellAITSenderID),
			serverutils.MustGetEnvVar(BeWellAITUsername),
			serverutils.MustGetEnvVar(BeWellAITAPIKey),
		)
	}
	return nil, fmt.Errorf("unknown AIT sender")
}

// SendSMS is a method used to send SMS
func (s Service) SendSMS(
	ctx context.Context,
	to, message string,
	from string,
	username string,
	key string,
) (*dto.SendMessageResponse, error) {
	ctx, span := tracer.Start(ctx, "SendSMS")
	defer span.End()
	values := url.Values{}
	values.Set("username", username)
	values.Set("to", to)
	values.Set("message", message)
	values.Set("from", from)

	smsURL := GetSmsURL(s.Env)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, err := s.newPostRequest(ctx, smsURL, values, headers, key)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	// read the response body to a variable
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
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
	ctx context.Context,
	url string,
	values url.Values,
	headers map[string]string,
	key string,
) (*http.Response, error) {
	_, span := tracer.Start(ctx, "newPostRequest")
	defer span.End()
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		helpers.RecordSpanError(span, err)
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
	ctx context.Context,
	to []string,
	message string,
	from enumutils.SenderID,
	segment string,
) (*dto.SendMessageResponse, error) {
	ctx, span := tracer.Start(ctx, "SendMarketingSMS")
	defer span.End()
	var whitelistedNumbers []string
	for _, phone := range to {
		optedOut, err := s.Crm.IsOptedOut(ctx, phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check if the number %s is opted out: %w", phone, err)
		}
		if !optedOut {
			whitelistedNumbers = append(whitelistedNumbers, phone)
		}
	}

	if len(whitelistedNumbers) == 0 {
		return nil, nil
	}

	recipients := strings.Join(whitelistedNumbers, ",")
	resp, err := s.Send(ctx, recipients, message, from)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to send SMS: %v", err)
	}

	smsMsgData := resp.SMSMessageData
	if smsMsgData == nil {
		return nil, nil
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
		}
		savedSms, err := s.SaveMarketingMessage(
			ctx,
			data,
		)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return resp, fmt.Errorf(
				"failed to create a message in our data store: %v",
				err,
			)
		}

		if _, err := s.Edi.UpdateMessageSent(
			ctx,
			data.PhoneNumber,
			segment,
		); err != nil {
			helpers.RecordSpanError(span, err)
			return resp, fmt.Errorf(
				"failed to update message sent status to true: %v",
				err,
			)
		}

		metadata := map[string]interface{}{
			"body": message,
		}
		if err = s.PubSub.NotifyEngagementCreate(
			ctx,
			data.PhoneNumber,
			savedSms.ID,
			"NOTE",
			metadata,
			common.EngagementCreateTopic,
		); err != nil {
			helpers.RecordSpanError(span, err)
			return resp, fmt.Errorf(
				"failed to publish to %v pub/sub topic with error: %v",
				common.EngagementCreateTopic,
				err,
			)
		}
	}

	return resp, nil
}
