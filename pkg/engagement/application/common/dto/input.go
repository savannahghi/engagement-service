package dto

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

// SendSMSPayload is used to serialise an SMS sent through the AIT service REST API
type SendSMSPayload struct {
	To      []string           `json:"to"`
	Message string             `json:"message"`
	Sender  enumutils.SenderID `json:"sender"`
	Segment *string            `json:"segment"`
}

// EMailMessage holds data required to send emails
type EMailMessage struct {
	Subject string   `json:"subject,omitempty"`
	Text    string   `json:"text,omitempty"`
	To      []string `json:"to,omitempty"`
}

// FeedbackInput is reason a user gave a certain NPS score
// Its stored as question answer in plain text
type FeedbackInput struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// Feedback is reason a user gave a certain NPS score
// Its stored as question answer in plain text
type Feedback struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// NPSInput is the input for a survey
type NPSInput struct {
	Name        string           `json:"name"`
	Score       int              `json:"score"`
	SladeCode   string           `json:"sladeCode"`
	Email       *string          `json:"email"`
	PhoneNumber *string          `json:"phoneNumber"`
	Feedback    []*FeedbackInput `json:"feedback"`
}

// ListID is a HubSpot Contact List ID
type ListID struct {
	ListID int `json:"listId"`
}

// SetBewellAwareInput payload to set bewell aware
type SetBewellAwareInput struct {
	EmailAddress string `json:"email"`
}

// LoadCampgainDataInput input
type LoadCampgainDataInput struct {
	PhoneNumber *string  `json:"phoneNumber"`
	Emails      []string `json:"emails"`
}

// MarketingSMS represents marketing SMS data
type MarketingSMS struct {
	ID                   string                `json:"id"`
	PhoneNumber          string                `json:"phoneNumber"`
	SenderID             enumutils.SenderID    `json:"senderId"`
	MessageSentTimeStamp time.Time             `json:"messageSentTimeStamp"`
	Message              string                `json:"message"`
	DeliveryReport       *ATDeliveryReport     `json:"deliveryReport"`
	Status               string                `json:"status"`
	IsSynced             bool                  `json:"isSynced"`
	TimeSynced           *time.Time            `json:"timeSynced"`
	Engagement           domain.EngagementData `json:"engagement"`
}

// ATDeliveryReport callback delivery reports
type ATDeliveryReport struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	PhoneNumber             string    `json:"phoneNumber"`
	NetworkCode             *string   `json:"networkCode"`
	FailureReason           *string   `json:"failureReason"`
	RetryCount              int       `json:"retryCount"`
	DeliveryReportTimeStamp time.Time `json:"deliveryReportTimeStamp"`
}

// PrimaryEmailAddressPayload used when collecting HubSpot user email.
type PrimaryEmailAddressPayload struct {
	PhoneNumber  string `json:"phone"`
	EmailAddress string `json:"email"`
}

// UpdateContactPSMessage represents CRM update contact Pub/Sub message
type UpdateContactPSMessage struct {
	Properties domain.ContactProperties `json:"properties"`
	Phone      string                   `json:"phone"`
}

// UIDPayload is the user ID used in some inter-service requests
type UIDPayload struct {
	UID *string `json:"uid"`
}

// ContactLeadInput ...
// todo there should be better management of this @mathenge
type ContactLeadInput struct {
	ContactType    string                   `json:"contact_type,omitempty"`
	ContactValue   string                   `json:"contact_value,omitempty"`
	FirstName      string                   `json:"first_name,omitempty"`
	LastName       string                   `json:"last_name,omitempty"`
	DateOfBirth    scalarutils.Date         `json:"date_of_birth,omitempty"`
	IsSync         bool                     `json:"isSync"  firestore:"IsSync"`
	TimeSync       *time.Time               `json:"timeSync"  firestore:"TimeSync"`
	OptOut         domain.GeneralOptionType `json:"opt_out,omitempty"`
	WantCover      bool                     `json:"wantCover" firestore:"wantCover"`
	ContactChannel string                   `json:"contact_channel,omitempty"`
	IsRegistered   bool                     `json:"is_registered,omitempty"`
}

// OutgoingEmailsLog contains the content of the sent email message sent via MailGun
type OutgoingEmailsLog struct {
	UUID    string   `json:"uuid" firestore:"uuid"`
	To      []string `json:"to" firestore:"to"`
	From    string   `json:"from" firestore:"from"`
	Subject string   `json:"subject" firestore:"subject"`
	Text    string   `json:"text" firestore:"text"`
	// MessageID is a unique identifier of mailgun's message
	MessageID   string              `json:"message-id" firestore:"messageID"`
	EmailSentOn time.Time           `json:"emailSentOn" firestore:"emailSentOn"`
	Event       *MailgunEventOutput `json:"mailgunEvent" firestore:"mailgunEvent"`
}

// MailgunEvent represents mailgun event input e.g delivered, rejected etc
type MailgunEvent struct {
	EventName   string `json:"event" firestore:"event"`
	DeliveredOn string `json:"timestamp" firestore:"deliveredOn"`
	// MessageID is a unique identifier of mailgun's message
	MessageID string `json:"message-id" firestore:"messageID"`
}

// EngagementPubSubMessage represents engagement payload published to pubsub
type EngagementPubSubMessage struct {
	Engagement  domain.EngagementData `json:"engagement"`
	PhoneNumber string                `json:"phoneNumber"`
	MessageID   string                `json:"messageID"`
}

// RetrieveUserProfileInput used to retrieve user profile info using either email address or phone
type RetrieveUserProfileInput struct {
	PhoneNumber  *string `json:"phone"`
	EmailAddress *string `json:"email"`
}

//TemporaryPIN input used to send temporary PIN message
type TemporaryPIN struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	PIN         string `json:"pin,omitempty"`
	Channel     int    `json:"channel,omitempty"`
}
