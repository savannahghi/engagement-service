package dto

import (
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

// SendSMSPayload is used to serialise an SMS sent through the AIT service REST API
type SendSMSPayload struct {
	To      []string      `json:"to"`
	Message string        `json:"message"`
	Sender  base.SenderID `json:"sender"`
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

type LoadCampgainDataInput struct {
	PhoneNumber *string  `json:"phoneNumber"`
	Emails      []string `json:"emails"`
}

// MarketingSMS represents marketing SMS data
type MarketingSMS struct {
	ID                   string                `json:"id"`
	PhoneNumber          string                `json:"phoneNumber"`
	SenderID             base.SenderID         `json:"senderId"`
	MessageSentTimeStamp time.Time             `json:"messageSentTimeStamp"`
	Message              string                `json:"message"`
	DeliveryReport       *ATDeliveryReport     `json:"deliveryReport"`
	Status               string                `json:"status"`
	Engagement           domain.EngagementData `json:"engagement"`
	IsSynced             bool                  `json:"isSynced"`
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

// MarketingMessagePayload is used when retrieving the segmented data from the database
type MarketingMessagePayload struct {
	Wing           string `json:"wing" firestore:"wing"`
	InitialSegment string `json:"initialSegment" firestore:"initialSegment"`
}

// Segment represents the Segments data
type Segment struct {
	BeWellEnrolled        string `json:"be_well_enrolled" firestore:"be_well_enrolled"`
	OptOut                string `json:"opt_out" firestore:"opt_out"`
	BeWellAware           string `json:"be_well_aware" firestore:"be_well_aware"`
	BeWellPersona         string `json:"be_well_persona" firestore:"be_well_persona"`
	HasWellnessCard       string `json:"has_wellness_card" firestore:"has_wellness_card"`
	HasCover              string `json:"has_cover" firestore:"has_cover"`
	Payor                 string `json:"payor" firestore:"payor"`
	FirstChannelOfContact string `json:"first_channel_of_contact" firestore:"first_channel_of_contact"`
	InitialSegment        string `json:"initial_segment" firestore:"initial_segment"`
	HasVirtualCard        string `json:"has_virtual_card" firestore:"has_virtual_card"`
	Email                 string `json:"email" firestore:"email"`
	PhoneNumber           string `json:"phone" firestore:"phone"`
	FirstName             string `json:"firstname" firestore:"firstname"`
	LastName              string `json:"lastname" firestore:"lastname"`
	Wing                  string `json:"wing" firestore:"wing"`
	MessageSent           string `json:"message_sent" firestore:"message_sent"`
	IsSynced              string `json:"is_synced" firestore:"is_synced"`
	TimeSynced            string `json:"time_synced" firestore:"time_synced"`
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
	DateOfBirth    base.Date                `json:"date_of_birth,omitempty"`
	IsSync         bool                     `json:"isSync"  firestore:"IsSync"`
	TimeSync       *time.Time               `json:"timeSync"  firestore:"TimeSync"`
	OptOut         domain.GeneralOptionType `json:"opt_out,omitempty"`
	WantCover      bool                     `json:"wantCover" firestore:"wantCover"`
	ContactChannel string                   `json:"contact_channel,omitempty"`
	IsRegistered   bool                     `json:"is_registered,omitempty"`
}
