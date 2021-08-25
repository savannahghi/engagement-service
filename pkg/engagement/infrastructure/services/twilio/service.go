package twilio

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

	"github.com/asaskevich/govalidator"
	"github.com/kevinburke/twilio-go"
	"github.com/kevinburke/twilio-go/token"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement-service/pkg/engagement/repository"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/otel"
	"moul.io/http2curl"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement-service/pkg/engagement/services/twilio")

/* #nosec */
// DefaultTwilioRegion is set to global low latency auto-selection
const (
	TwilioRegionEnvVarName            = "TWILIO_REGION"
	TwilioVideoAPIURLEnvVarName       = "TWILIO_VIDEO_API_URL"
	TwilioVideoAPIKeySIDEnvVarName    = "TWILIO_VIDEO_SID"
	TwilioVideoAPIKeySecretEnvVarName = "TWILIO_VIDEO_SECRET"
	TwilioAccountSIDEnvVarName        = "TWILIO_ACCOUNT_SID"
	TwilioAccountAuthTokenEnvVarName  = "TWILIO_ACCOUNT_AUTH_TOKEN"
	TwilioSMSNumberEnvVarName         = "TWILIO_SMS_NUMBER"
	ServerPublicDomainEnvVarName      = "SERVER_PUBLIC_DOMAIN"
	TwilioCallbackPath                = "/twilio_callback"
	TwilioHTTPClientTimeoutSeconds    = 10
	TwilioPeerToPeerMaxParticipants   = 3
	TwilioAccessTokenTTL              = 14400
)

// ServiceTwilio defines the interaction with the twilio service
type ServiceTwilio interface {
	MakeTwilioRequest(
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error

	Room(ctx context.Context) (*dto.Room, error)

	TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error)

	SendSMS(ctx context.Context, to string, msg string) error

	SaveTwilioVideoCallbackStatus(
		ctx context.Context,
		data dto.CallbackData,
	) error
}

// NewService initializes a service to interact with Twilio
func NewService(sms sms.ServiceSMS, repo repository.Repository) *Service {
	region := serverutils.MustGetEnvVar(TwilioRegionEnvVarName)
	videoBaseURL := serverutils.MustGetEnvVar(TwilioVideoAPIURLEnvVarName)
	videoAPIKeySID := serverutils.MustGetEnvVar(TwilioVideoAPIKeySIDEnvVarName)
	videoAPIKeySecret := serverutils.MustGetEnvVar(TwilioVideoAPIKeySecretEnvVarName)
	accountSID := serverutils.MustGetEnvVar(TwilioAccountSIDEnvVarName)
	accountAuthToken := serverutils.MustGetEnvVar(TwilioAccountAuthTokenEnvVarName)
	httpClient := &http.Client{
		Timeout: time.Second * TwilioHTTPClientTimeoutSeconds,
	}
	publicDomain := serverutils.MustGetEnvVar(ServerPublicDomainEnvVarName)
	callbackURL := publicDomain + TwilioCallbackPath
	smsNumber := serverutils.MustGetEnvVar(TwilioSMSNumberEnvVarName)

	srv := &Service{
		region:            region,
		videoBaseURL:      videoBaseURL,
		videoAPIKeySID:    videoAPIKeySID,
		videoAPIKeySecret: videoAPIKeySecret,
		accountSID:        accountSID,
		accountAuthToken:  accountAuthToken,
		httpClient:        httpClient,
		twilioClient:      twilio.NewClient(accountSID, accountAuthToken, httpClient),
		callbackURL:       callbackURL,
		smsNumber:         smsNumber,
		sms:               sms,
		repository:        repo,
	}
	srv.checkPreconditions()
	return srv
}

// Service organizes methods needed to interact with Twilio for video, voice
// and text
type Service struct {
	region            string
	videoBaseURL      string
	videoAPIKeySID    string
	videoAPIKeySecret string
	accountSID        string
	accountAuthToken  string
	httpClient        *http.Client
	twilioClient      *twilio.Client
	callbackURL       string
	smsNumber         string
	sms               sms.ServiceSMS
	repository        repository.Repository
}

func (s Service) checkPreconditions() {
	if s.region == "" {
		log.Panicf("Twilio region not set")
	}

	if s.videoBaseURL == "" {
		log.Panicf("Twilio video base URL not set")
	}

	if !govalidator.IsURL(s.videoBaseURL) {
		log.Panicf("Twilio Video base URL (%s) is not a valid URL", s.videoBaseURL)
	}

	if s.videoAPIKeySID == "" {
		log.Panicf("Twilio Video API Key SID not set")
	}

	if s.videoAPIKeySecret == "" {
		log.Panicf("Twilio Video API Key secret not set")
	}

	if s.accountSID == "" {
		log.Panicf("Twilio Video account SID not set")
	}

	if s.accountAuthToken == "" {
		log.Panicf("Twilio Video account auth token not set")
	}

	if s.httpClient == nil {
		log.Panicf("nil HTTP client in Twilio service")
	}

	if s.twilioClient == nil {
		log.Panicf("nil Twilio client in Twilio service")
	}

	if s.callbackURL == "" {
		log.Panicf("empty Twilio callback URL")
	}
}

// MakeTwilioRequest makes a twilio request
func (s Service) MakeTwilioRequest(
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	s.checkPreconditions()
	if serverutils.IsDebug() {
		log.Printf("Twilio request data: \n%s\n", content)
	}

	r := strings.NewReader(content.Encode())
	req, reqErr := http.NewRequest(method, s.videoBaseURL+urlPath, r)
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.accountSID, s.accountAuthToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("twilio API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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
		return fmt.Errorf("unable to unmarshal Twilio resp: %w", err)
	}

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	return nil
}

// Room represents a real-time audio, data, video, and/or screen-share session,
// and is the basic building block for a Programmable Video application.
//
// In a Peer-to-peer Room, media flows directly between participants. This
// supports up to 10 participants in a mesh topology.
//
// In a Group Room, media is routed through Twilio's Media Servers. This
// supports up to 50 participants.
//
// Participants represent client applications that are connected to a Room and
// sharing audio, data, and/or video media with one another.
//
// Tracks represent the individual audio, data, and video media streams that
// are shared within a Room.
//
// LocalTracks represent the audio, data, and video captured from the local
// client's media sources (for example, microphone and camera).
//
// RemoteTracks represent the audio, data, and video tracks from other
// participants connected to the Room.
//
// Room names must be unique within an account.
//
// Rooms created via the REST API exist for five minutes to allow the first
// Participant to connect. If no Participants join within five minutes,
// the Room times out and a new Room must be created.
//
// Because of confidentiality issues in healthcare, we do not enable recording
// for these meetings.
func (s Service) Room(ctx context.Context) (*dto.Room, error) {
	_, span := tracer.Start(ctx, "Room")
	defer span.End()
	s.checkPreconditions()

	roomReqData := url.Values{}
	roomReqData.Set("Type", "peer-to-peer")
	roomReqData.Set("MaxParticipants", strconv.Itoa(TwilioPeerToPeerMaxParticipants))
	roomReqData.Set("StatusCallbackMethod", "POST")
	roomReqData.Set("StatusCallback", s.callbackURL)
	roomReqData.Set("EnableTurn", strconv.FormatBool(true))

	var room dto.Room
	err := s.MakeTwilioRequest("POST", "/v1/Rooms", roomReqData, &room)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("twilio room API call error: %w", err)
	}
	return &room, nil
}

// TwilioAccessToken is used to generate short-lived credentials used to authenticate
// the client-side application to Twilio.
//
// An access token should be generated for every user of the application.
//
// An access token can optionally encode a room name, which would allow the user
// to connect only to the room specified in the token.
//
// Access tokens are JSON Web Tokens (JWTs).
func (s Service) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	ctx, span := tracer.Start(ctx, "TwilioAccessToken")
	defer span.End()
	s.checkPreconditions()

	uid, err := firebasetools.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get logged in user uid: %w", err)
	}

	room, err := s.Room(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get room to issue a grant to: %w", err)
	}

	ttl := time.Second * TwilioAccessTokenTTL
	accessToken := token.New(s.accountSID, s.videoAPIKeySID, s.videoAPIKeySecret, uid, ttl)
	videoGrant := token.NewVideoGrant(room.SID)
	accessToken.AddGrant(videoGrant)

	jwt, err := accessToken.JWT()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to generate JWT for Twilio access token: %w", err)
	}
	payload := dto.AccessToken{
		JWT:             jwt,
		UniqueName:      room.UniqueName,
		SID:             room.SID,
		DateUpdated:     room.DateUpdated,
		Status:          room.Status,
		Type:            room.Type,
		MaxParticipants: room.MaxParticipants,
		Duration:        room.Duration,
	}
	return &payload, nil
}

// SendSMS sends a text message through Twilio's programmable SMS
func (s Service) SendSMS(ctx context.Context, to string, msg string) error {
	_, span := tracer.Start(ctx, "SendSMS")
	defer span.End()
	s.checkPreconditions()

	t, err := s.twilioClient.Messages.SendMessage(s.smsNumber, to, msg, nil)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio SMS API error: %w", err)
	}

	if t.ErrorCode != 0 {
		return fmt.Errorf("sms could not be sent: %v", t.ErrorMessage)
	}

	fmt.Printf("Raw Twilio SMS response: %v", t)
	return nil
}

// SaveTwilioVideoCallbackStatus saves status callback data
func (s Service) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return s.repository.SaveTwilioVideoCallbackStatus(ctx, data)
}
