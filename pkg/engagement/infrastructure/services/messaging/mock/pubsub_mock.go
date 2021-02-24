package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceMessaging is a mock implementation of the "messaging" service
type FakeServiceMessaging struct {
	NotifyFn func(
		ctx context.Context,
		topicID string,
		uid string,
		flavour base.Flavour,
		payload base.Element,
		metadata map[string]interface{},
	) error

	// Ask the notification service about the topics that it knows about
	TopicIDsFn func() []string

	SubscriptionIDsFn func() map[string]string

	ReverseSubscriptionIDsFn func() map[string]string
}

// Notify Sends a message to a topic
func (f *FakeServiceMessaging) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour base.Flavour,
	payload base.Element,
	metadata map[string]interface{},
) error {
	return f.NotifyFn(ctx, topicID, uid, flavour, payload, metadata)
}

// TopicIDs ...
func (f *FakeServiceMessaging) TopicIDs() []string {
	return f.TopicIDsFn()
}

// SubscriptionIDs ...
func (f *FakeServiceMessaging) SubscriptionIDs() map[string]string {
	return f.SubscriptionIDsFn()
}

// ReverseSubscriptionIDs ...
func (f *FakeServiceMessaging) ReverseSubscriptionIDs() map[string]string {
	return f.ReverseSubscriptionIDsFn()
}
