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

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

// AIT environment variables names
const (
	APIKeyEnvVarName       = "AIT_API_KEY"
	APIUsernameEnvVarName  = "AIT_USERNAME"
	APISenderIDEnvVarName  = "AIT_SENDER_ID"
	AITEnvVarName          = "AIT_ENVIRONMENT"
	AITAuthenticationError = "The supplied authentication is invalid"
	//AITCallbackCollectionName is the name of a Cloud Firestore collection into which AIT
	// callback data will be saved for future analysis
	AITCallbackCollectionName = "ait_callbacks"
)

// ServiceSMS defines the interactions with sms service
type ServiceSMS interface {
	SendToMany(message string, to []string) (*dto.SendMessageResponse, error)
	Send(to, message string) (*dto.SendMessageResponse, error)
	SaveAITCallbackResponse(ctx context.Context, data dto.CallbackData) error
}

// Service defines a sms service struct
type Service struct {
	Username   string
	APIKey     string
	Env        string
	From       string
	Repository repository.Repository
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
	return fmt.Sprintf("https://%s.sandbox.africastalking.com", service)

}

// NewService returns a new service
func NewService(repository repository.Repository) *Service {
	username := base.MustGetEnvVar(APIUsernameEnvVarName)
	apiKey := base.MustGetEnvVar(APIKeyEnvVarName)
	from := base.MustGetEnvVar(APISenderIDEnvVarName)
	env := base.MustGetEnvVar(AITEnvVarName)
	return &Service{username, apiKey, env, from, repository}
}

// SaveAITCallbackResponse saves the callback data for future analysis
func (s Service) SaveAITCallbackResponse(ctx context.Context, data dto.CallbackData) error {
	return s.Repository.SaveAITCallbackResponse(ctx, data)
}

// SendToMany is a utility method to send to many recipients at the same time
func (s Service) SendToMany(message string, to []string) (*dto.SendMessageResponse, error) {
	recipients := strings.Join(to, ",")
	return s.Send(recipients, message)
}

// Send is a method used to send to a single recipient
func (s Service) Send(to, message string) (*dto.SendMessageResponse, error) {
	values := url.Values{}
	values.Set("username", s.Username)
	values.Set("to", to)
	values.Set("message", message)
	values.Set("from", s.From)

	smsURL := GetSmsURL(s.Env)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, err := s.newPostRequest(smsURL, values, headers)
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
		return nil, errors.New(fmt.Errorf("SMS Error : unable to parse sms response ; %v", err).Error())
	}

	return smsMessageResponse, nil
}

func (s Service) newPostRequest(url string, values url.Values, headers map[string]string) (*http.Response, error) {
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Length", strconv.Itoa(reader.Len()))
	req.Header.Set("apikey", s.APIKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
