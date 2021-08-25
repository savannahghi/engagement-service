package fcm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (
	fcmServiceName = "fcm"
	fcmVersion     = ""
)

// PushService defines the behavior of our FCM push implementation
type PushService interface {
	Push(
		ctx context.Context,
		sender string,
		payload firebasetools.SendNotificationPayload,
	) error
}

// NewRemotePushService initializes an FCM push service
func NewRemotePushService(
	ctx context.Context,
) (*RemotePushService, error) {
	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
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
	notificationPayload firebasetools.SendNotificationPayload,
) error {
	ctx, span := tracer.Start(ctx, "Push")
	defer span.End()
	rfs.checkPreconditions()
	env := serverutils.GetRunningEnvironment()
	payload, err := json.Marshal(notificationPayload)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't marshal notification payload: %w", err)
	}

	err = pubsubtools.PublishToPubsub(
		ctx,
		rfs.pubsubClient,
		helpers.AddPubSubNamespace(common.FcmPublishTopic),
		env,
		fcmServiceName,
		fcmVersion,
		payload,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't publish FCM message to pubsub: %w", err)
	}

	return nil
}
