package feed

import "context"

// NotificationService defines the behavior of our notifications
type NotificationService interface {

	// Send a message to a topic
	Notify(ctx context.Context, topicID string, el Element) error

	// Ask the notification service about the topics that it knows about
	TopicIDs() []string

	// Ask the notification service about the subscriptions that it knows about,
	// as a map of topics to subscriptions
	SubscriptionIDs() map[string]string

	// Ask the notification service about the subscriptions that it knows about,
	// as a map of subscriptions to topics
	ReverseSubscriptionIDs() map[string]string
}
