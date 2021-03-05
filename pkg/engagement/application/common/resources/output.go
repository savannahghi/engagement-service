package resources

import "gitlab.slade360emr.com/go/base"

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
