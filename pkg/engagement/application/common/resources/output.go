package resources

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

// NotificationEnvelope is used to "wrap" elements with context and metadata
// before they are sent as notifications.
//
// This context and metadata allows the recipients of the notifications to
// process them intelligently.
type NotificationEnvelope struct {
	UID      string                 `json:"uid"`
	Flavour  base.Flavour           `json:"flavour"`
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

// CallbackData records data sent back from the AIT API to our HTTP callback URL
type CallbackData struct {
	Values map[string][]string `json:"values,omitempty" firestore:"values,omitempty"`
}

// Message is a Twilio WhatsApp or SMS message
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
func (c *Message) GetID() base.ID {
	return base.IDValue(c.ID)
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
	Msisdn string `json:"msisdn"`
}

// GenerateRetryOTP is an input struct for generating and
// sending fallback otp
type GenerateRetryOTP struct {
	Msisdn    *string `json:"msisdn"`
	RetryStep int     `json:"retryStep"`
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
