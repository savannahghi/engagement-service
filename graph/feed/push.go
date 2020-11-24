package feed

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

// FCMPushService defines the behavior of our FCM push implementation
type FCMPushService interface {
	Push(
		ctx context.Context,
		sender string,
		payload base.SendNotificationPayload,
	) error
}

// NewRemoteFCMPushService initializes an FCM push service
func NewRemoteFCMPushService(
	ctx context.Context,
) (*RemoteFCMPushService, error) {
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
	rfs := &RemoteFCMPushService{
		pubsubClient: pubsubClient,
	}
	rfs.checkPreconditions()
	return rfs, nil
}

// RemoteFCMPushService sends instructions to a remote FCM service over
// Google Cloud Pub-Sub
type RemoteFCMPushService struct {
	pubsubClient *pubsub.Client
}

func (rfs RemoteFCMPushService) checkPreconditions() {
	if rfs.pubsubClient == nil {
		log.Panicf(
			"attempt to use a remote FCM push service with a nil pubsubClient")
	}
}

// Push instructs a remote FCM service to send a push notification.
//
// This is done over Google Cloud Pub-Sub.
func (rfs RemoteFCMPushService) Push(
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
