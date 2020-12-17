package feed

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// NotificationService defines the behavior of our notifications
type NotificationService interface {

	// Send a message to a topic
	Notify(
		ctx context.Context,
		topicID string,
		uid string,
		flavour base.Flavour,
		payload base.Element,
		metadata map[string]interface{},
	) error

	// Ask the notification service about the topics that it knows about
	TopicIDs() []string
}

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
