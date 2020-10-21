package messaging

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/feed/graph/feed"
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
}

func (mn MockNotificationService) checkPreconditions() error {
	if mn.notifications == nil {
		mn.notifications = make(map[string][]feed.Element)
	}

	return nil
}

// Notify MOCKS sending of a feed element to the message bus
func (mn MockNotificationService) Notify(
	ctx context.Context,
	channel string,
	el feed.Element,
) error {
	if err := mn.checkPreconditions(); err != nil {
		return fmt.Errorf("pubsub service precondition check failed: %w", err)
	}

	return nil
}
