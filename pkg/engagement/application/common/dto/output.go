package dto

import (
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
)

// NotificationEnvelope is used to "wrap" elements with context and metadata
// before they are sent as notifications.
//
// This context and metadata allows the recipients of the notifications to
// process them intelligently.
type NotificationEnvelope struct {
	UID      string                 `json:"uid"`
	Flavour  feedlib.Flavour        `json:"flavour"`
	Payload  []byte                 `json:"payload"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Recipient returns the details of a message recipient
type Recipient struct {
	Number    string `json:"number"`
	Cost      string `json:"cost"`
	Status    string `json:"status"`
	MessageID string `json:"messageID"`
}

// SMS returns the message details of a recipient
type SMS struct {
	Recipients []Recipient `json:"Recipients"`
}

// SendMessageResponse returns a message response with the recipient's details
type SendMessageResponse struct {
	SMSMessageData *SMS `json:"SMSMessageData"`
}

// Message is a Twilio WhatsApp or SMS message
// todo: clean this up in subsequent MR (@mathenge)
type Message struct {
	ID                  string            `json:"id" firestore:"id"`
	AccountSID          string            `json:"account_sid" firestore:"account_sid"`
	APIVersion          string            `json:"api_version" firestore:"api_version"`
	Body                string            `json:"body" firestore:"body"`
	DateCreated         string            `json:"date_created" firestore:"date_created"`
	DateSent            string            `json:"date_sent" firestore:"date_sent"`
	DateUpdated         string            `json:"date_updated" firestore:"date_updated"`
	Direction           string            `json:"direction" firestore:"direction"`
	ErrorCode           *string           `json:"error_code" firestore:"error_code"`
	ErrorMessage        *string           `json:"error_message" firestore:"error_message"`
	From                string            `json:"from" firestore:"from"`
	MessagingServiceSID string            `json:"messaging_service_sid" firestore:"messaging_service_sid"`
	NumMedia            string            `json:"num_media" firestore:"num_media"`
	NumSegments         string            `json:"num_segments" firestore:"num_segments"`
	Price               *string           `json:"price" firestore:"price"`
	PriceUnit           *string           `json:"price_unit" firestore:"price_unit"`
	SID                 string            `json:"sid" firestore:"sid"`
	Status              string            `json:"status" firestore:"status"`
	SubresourceURLs     map[string]string `json:"subresource_uris" firestore:"subresource_uris"`
	To                  string            `json:"to" firestore:"to"`
	URI                 string            `json:"uri" firestore:"uri"`
}

// IsNode is a "label" that marks this struct (and those that embed it) as
// implementations of the "Base" interface defined in our GraphQL schema.
func (c *Message) IsNode() {}

// GetID returns the struct's ID value
func (c *Message) GetID() firebasetools.ID {
	return firebasetools.IDValue(c.ID)
}

// SetID sets the struct's ID value
func (c *Message) SetID(id string) {
	c.ID = id
}

// Dummy ..
type Dummy struct {
	id string
}

//IsEntity ...
func (d Dummy) IsEntity() {}

// IsNode ..
func (d *Dummy) IsNode() {}

// SetID sets the trace's ID
func (d *Dummy) SetID(id string) {
	d.id = id
}

// OTP is used to persist and verify authorization codes
// (single use 'One Time PIN's)
type OTP struct {
	MSISDN            string    `json:"msisdn,omitempty" firestore:"msisdn"`
	Message           string    `json:"message,omitempty" firestore:"message"`
	AuthorizationCode string    `json:"authorizationCode,omitempty" firestore:"authorizationCode"`
	Timestamp         time.Time `json:"timestamp,omitempty" firestore:"timestamp"`
	IsValid           bool      `json:"isValid,omitempty" firestore:"isValid"`
	Email             string    `json:"email,omitempty" firestore:"email"`
}

// Msisdn is an input struct for a generating and sending an otp request
type Msisdn struct {
	Msisdn string  `json:"msisdn"`
	AppID  *string `json:"appId"`
}

// GenerateRetryOTP is an input struct for generating and
// sending fallback otp
type GenerateRetryOTP struct {
	Msisdn    *string `json:"msisdn"`
	RetryStep int     `json:"retryStep"`
	AppID     *string `json:"appId"`
}

// VerifyOTP is an input struct for confirming an otp
type VerifyOTP struct {
	Msisdn           *string `json:"msisdn"`
	Email            *string `json:"email"`
	VerificationCode *string `json:"verificationCode"`
}

// Room is used to serialize details of Twilio meeting rooms that are
// created via the Twilio Video REST API
type Room struct {
	SID                         string            `json:"sid"`
	Status                      string            `json:"status"`
	DateCreated                 time.Time         `json:"date_created"`
	DateUpdated                 time.Time         `json:"date_updated"`
	AccountSID                  string            `json:"account_sid"`
	EnableTurn                  bool              `json:"enable_turn"`
	UniqueName                  string            `json:"unique_name"`
	StatusCallbackMethod        string            `json:"status_callback_method"`
	Type                        string            `json:"type"`
	MaxParticipants             int               `json:"max_participants"`
	RecordParticipantsOnConnect bool              `json:"record_participants_on_connect"`
	VideoCodecs                 []string          `json:"video_codecs"`
	MediaRegion                 string            `json:"media_region"`
	URL                         string            `json:"url"`
	Links                       map[string]string `json:"links"`
	StatusCallback              *string           `json:"status_callback,omitempty"`
	Duration                    *int              `json:"duration,omitempty"`
	EndTime                     *time.Time        `json:"end_time,omitempty"`
}

// AccessToken is used to return the results of requesting an access token.
//
// In addition to the JWT, this includes the details that are needed in order to
// connect to the room on the client side
type AccessToken struct {
	JWT             string    `json:"jwt,omitempty"`
	UniqueName      string    `json:"uniqueName,omitempty"`
	SID             string    `json:"sid,omitempty"`
	DateUpdated     time.Time `json:"dateUpdated,omitempty"`
	Status          string    `json:"status,omitempty"`
	Type            string    `json:"type,omitempty"`
	MaxParticipants int       `json:"maxParticipants,omitempty"`
	Duration        *int      `json:"duration,omitempty"`
}

// IsEntity ...
func (a AccessToken) IsEntity() {}

// SMSPayload is used to serialise an SMS sent through the twilio service REST API
type SMSPayload struct {
	To      []string `json:"to"`
	Message string   `json:"message"`
}

// SavedNotification is used to serialize and save successful FCM notifications.
//
// It's the basis for a primitive "inbox" - a mechanism by which an app can
// request it's messages in bulk.
type SavedNotification struct {
	ID                string                      `json:"id,omitempty"`
	RegistrationToken string                      `json:"registrationToken,omitempty"`
	MessageID         string                      `json:"messageID,omitempty"`
	Timestamp         time.Time                   `json:"timestamp,omitempty"`
	Data              map[string]interface{}      `json:"data,omitempty"`
	Notification      *FirebaseSimpleNotification `json:"notification,omitempty"`
	AndroidConfig     *FirebaseAndroidConfig      `json:"androidConfig,omitempty"`
	WebpushConfig     *FirebaseWebpushConfig      `json:"webpushConfig,omitempty"`
	APNSConfig        *FirebaseAPNSConfig         `json:"apnsConfig,omitempty"`
}

// IsEntity ...
func (u SavedNotification) IsEntity() {}

// FirebaseSimpleNotification is used to serialize simple FCM notification.
// It is a mirror of Firebase messaging.Notification
type FirebaseSimpleNotification struct {
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	ImageURL *string                `json:"imageURL"`
	Data     map[string]interface{} `json:"data"`
}

// FirebaseAPNSConfig is a mirror of Firebase messaging.APNSConfig
type FirebaseAPNSConfig struct {
	Headers map[string]interface{} `json:"headers"`
}

// FirebaseWebpushConfig is a mirror of Firebase messaging.WebpushConfig
type FirebaseWebpushConfig struct {
	Headers map[string]interface{} `json:"headers"`
	Data    map[string]interface{} `json:"data"`
}

// FirebaseAndroidConfig is a mirror of Firebase messaging.AndroidConfig
type FirebaseAndroidConfig struct {
	Priority              string                 `json:"priority"` // one of "normal" or "high"
	CollapseKey           *string                `json:"collapseKey"`
	RestrictedPackageName *string                `json:"restrictedPackageName"`
	Data                  map[string]interface{} `json:"data"` // if specified, overrides the Data field on Message type
}

// NPSResponse represents a single user feedback
type NPSResponse struct {
	ID        string     `json:"id" firestore:"id"`
	Name      string     `json:"name" firestore:"name"`
	Score     int        `json:"score" firestore:"score"`
	SladeCode string     `json:"sladeCode" firestore:"sladeCode"`
	Email     *string    `json:"email" firestore:"email"`
	MSISDN    *string    `json:"msisdn" firestore:"msisdn"`
	Feedback  []Feedback `json:"feedback" firestore:"feedback"`
	Timestamp time.Time  `json:"timestamp,omitempty" firestore:"timestamp,omitempty"`
}

// IsNode is a "label" that marks this struct (and those that embed it) as
// implementations of the "Base" interface defined in our GraphQL schema.
func (e *NPSResponse) IsNode() {}

// GetID returns the struct's ID value
func (e *NPSResponse) GetID() firebasetools.ID {
	return firebasetools.IDValue(e.ID)
}

// SetID sets the struct's ID value
func (e *NPSResponse) SetID(id string) {
	e.ID = id
}

// OKResp is used to return OK responses in inter-service calls
type OKResp struct {
	Status   string      `json:"status,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// NewOKResp ...
func NewOKResp(rawResponse interface{}) *OKResp {
	return &OKResp{
		Status:   "OK",
		Response: rawResponse,
	}
}

// MarketingDataLoadOutput ...
type MarketingDataLoadOutput struct {
	LoadingError                  error                            `json:"loading_error"`
	StartedAt                     time.Time                        `json:"started_at"`
	StoppedAt                     time.Time                        `json:"stopped_at"`
	HoursTaken                    float64                          `json:"hours_taken"`
	EntriesUniqueLoadedOnFirebase int                              `json:"entries_unique_loaded_on_firebase"`
	EntriesUniqueLoadedOnCRM      int                              `json:"entries_unique_loaded_on_crm"`
	TotalEntriesFoundOnFile       int                              `json:"total_entries_found_on_file"`
	Entries                       []MarketingDataLoadEntriesOutput `json:"marketing_data_load_entries_output"`
}

// MarketingDataLoadEntriesOutput ...
type MarketingDataLoadEntriesOutput struct {
	Identifier                  string `json:"identifier"`
	HasLoadedToFirebase         bool   `json:"has_loaded_to_firebase"`
	HasBeenRollBackFromFirebase bool   `json:"has_been_rolled_back_from_firebase"`
	HasLoadedToCRM              bool   `json:"has_loaded_to_crm"`
	FirebaseLoadError           error  `json:"firebase_load_error"`
	CRMLoadError                error  `json:"crm_load_error"`
}

// MailgunEventOutput represents the MailGun's event name and delivery time in standardized time
// since mailgun gives us time as unixtimestamp
type MailgunEventOutput struct {
	// EventName is the name of every event that happens to your emails e.g delivered, rejected etc
	EventName   string    `json:"event" firestore:"event"`
	DeliveredOn time.Time `json:"timestamp" firestore:"deliveredOn"`
}

// CallbackData records data sent back from the Twilio API to our HTTP callback URL
type CallbackData struct {
	Values map[string][]string `json:"values,omitempty" firestore:"values,omitempty"`
}
