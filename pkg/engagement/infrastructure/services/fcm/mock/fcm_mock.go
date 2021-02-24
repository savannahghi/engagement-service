package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceFcm simulates the behavior of our FCM push implementation
type FakeServiceFcm struct {
	PushFn func(
		ctx context.Context,
		sender string,
		payload base.SendNotificationPayload,
	) error
}

// Push instructs a remote FCM service to send a push notification.
//
// This is done over Google Cloud Pub-Sub.
func (f *FakeServiceFcm) Push(
	ctx context.Context,
	sender string,
	payload base.SendNotificationPayload,
) error {
	return f.PushFn(ctx, sender, payload)
}
