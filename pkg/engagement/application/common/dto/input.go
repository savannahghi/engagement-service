package dto

import (
	"time"

	"gitlab.slade360emr.com/go/base"
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
	EmailAddress string `json:"emailAddress"`
}

// MarketingSMS represents marketing SMS data
type MarketingSMS struct {
	ID                   string            `json:"id"`
	PhoneNumber          string            `json:"phoneNumber"`
	SenderID             base.SenderID     `json:"senderId"`
	MessageSentTimeStamp time.Time         `json:"messageSentTimeStamp"`
	Message              string            `json:"message"`
	DeliveryReport       *ATDeliveryReport `json:"deliveryReport"`
	Status               string            `json:"status"`
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
