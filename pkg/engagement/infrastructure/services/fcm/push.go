package fcm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/base"
)

const (
	fcmPublishTopic = "fcm.send_notification"
	fcmServiceName  = "fcm"
	fcmVersion      = ""
)

// PushService defines the behavior of our FCM push implementation
type PushService interface {
	Push(
		ctx context.Context,
		sender string,
		payload base.SendNotificationPayload,
	) error
}

// NewRemotePushService initializes an FCM push service
func NewRemotePushService(
	ctx context.Context,
) (*RemotePushService, error) {
	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			base.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}
	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("can't initialize pubsub client: %w", err)
	}
	rfs := &RemotePushService{
		pubsubClient: pubsubClient,
	}
	rfs.checkPreconditions()
	return rfs, nil
}

// RemotePushService sends instructions to a remote FCM service over
// Google Cloud Pub-Sub
type RemotePushService struct {
	pubsubClient *pubsub.Client
}

func (rfs RemotePushService) checkPreconditions() {
	if rfs.pubsubClient == nil {
		log.Panicf(
			"attempt to use a remote FCM push service with a nil pubsubClient")
	}
}

// Push instructs a remote FCM service to send a push notification.
//
// This is done over Google Cloud Pub-Sub.
func (rfs RemotePushService) Push(
	ctx context.Context,
	sender string,
	notificationPayload base.SendNotificationPayload,
) error {
	rfs.checkPreconditions()
	env := base.GetRunningEnvironment()
	payload, err := json.Marshal(notificationPayload)
	if err != nil {
		return fmt.Errorf("can't marshal notification payload: %w", err)
	}

	err = base.PublishToPubsub(
		ctx,
		rfs.pubsubClient,
		fcmPublishTopic,
		env,
		fcmServiceName,
		fcmVersion,
		payload,
	)
	if err != nil {
		return fmt.Errorf("can't publish FCM message to pubsub: %w", err)
	}

	return nil
}
