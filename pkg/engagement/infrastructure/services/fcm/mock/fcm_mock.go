package mock

import (
	"context"

	"github.com/savannahghi/firebasetools"
)

// FakeServiceFcm simulates the behavior of our FCM push implementation
type FakeServiceFcm struct {
	PushFn func(
		ctx context.Context,
		sender string,
		payload firebasetools.SendNotificationPayload,
	) error
}

// Push instructs a remote FCM service to send a push notification.
//
// This is done over Google Cloud Pub-Sub.
func (f *FakeServiceFcm) Push(
	ctx context.Context,
	sender string,
	payload firebasetools.SendNotificationPayload,
) error {
	return f.PushFn(ctx, sender, payload)
}
