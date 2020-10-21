package feed

import "context"

// NotificationService defines the behavior of our notifications
type NotificationService interface {

	// Send a message to a topic
	Notify(ctx context.Context, topicID string, el Element) error
}
