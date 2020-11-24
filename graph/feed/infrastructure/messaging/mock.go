package messaging

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/engagement/graph/feed"
)

// NewMockNotificationService initializes a mock notification service
func NewMockNotificationService() (*MockNotificationService, error) {
	mn := &MockNotificationService{}
	if err := mn.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"test notification service failed preconditions: %w", err)
	}
	return mn, nil
}

// MockNotificationService is used to mock notifications in-memory for tests
type MockNotificationService struct {
	notifications map[string][]feed.Element
	topicIDs      []string
	subscriptions map[string]string
}

func (mn MockNotificationService) checkPreconditions() error {
	if mn.notifications == nil {
		mn.notifications = make(map[string][]feed.Element)
	}

	if mn.topicIDs == nil {
		mn.topicIDs = make([]string, 1)
	}

	if mn.subscriptions == nil {
		mn.subscriptions = make(map[string]string)
	}

	return nil
}

// Notify MOCKS sending of a feed element to the message bus
func (mn MockNotificationService) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour feed.Flavour,
	payload feed.Element,
	metadata map[string]interface{},
) error {
	if err := mn.checkPreconditions(); err != nil {
		return fmt.Errorf("pubsub service precondition check failed: %w", err)
	}

	return nil
}

// TopicIDs returns topic IDs known to this mock notification service
func (mn MockNotificationService) TopicIDs() []string {
	if err := mn.checkPreconditions(); err != nil {
		panic("mock pubsub service precondition check failed")
	}

	return mn.topicIDs
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs for the topics
// and subscriptions known to this mock notification service
func (mn MockNotificationService) SubscriptionIDs() map[string]string {
	if err := mn.checkPreconditions(); err != nil {
		panic("mock pubsub service precondition check failed")
	}

	return mn.subscriptions
}

// ReverseSubscriptionIDs returns a map of subscription IDs to topic IDs for
// the subscriptions known to this mock service
func (mn MockNotificationService) ReverseSubscriptionIDs() map[string]string {
	if err := mn.checkPreconditions(); err != nil {
		panic("mock pubsub service precondition check failed")
	}
	return reverseMap(mn.subscriptions)
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}
