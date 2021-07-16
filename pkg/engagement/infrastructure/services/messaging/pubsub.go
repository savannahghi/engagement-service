package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"go.opentelemetry.io/otel"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/base"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/messaging")

// messaging related constants
const (
	hostNameEnvVarName = "SERVICE_HOST" // host at which this service is deployed
)

// NotificationService represents logic required to communicate with pubsub
// it defines the behavior of our notifications
type NotificationService interface {

	// Send a message to a topic
	Notify(
		ctx context.Context,
		topicID string,
		uid string,
		flavour base.Flavour,
		payload base.Element,
		metadata map[string]interface{},
	) error

	// Ask the notification service about the topics that it knows about
	TopicIDs() []string

	SubscriptionIDs() map[string]string

	ReverseSubscriptionIDs() map[string]string
}

// NewPubSubNotificationService initializes a live notification service
func NewPubSubNotificationService(
	ctx context.Context,
	projectID string,
) (NotificationService, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	environment, err := serverutils.GetEnvVar(serverutils.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to get the environment variable `%s`: %w", serverutils.Environment, err)
	}

	hostName, err := serverutils.GetEnvVar(hostNameEnvVarName)
	if err != nil {
		return nil, fmt.Errorf("unable to get the %s environment variable: %w", hostNameEnvVarName, err)
	}

	callbackURL := fmt.Sprintf("%s%s", hostName, pubsubtools.PubSubHandlerPath)
	ns := &PubSubNotificationService{
		client:      client,
		environment: environment,
		callbackURL: callbackURL,
	}
	if err := ns.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"pubsub notification service failed preconditions: %w", err)
	}

	topicIDs := ns.TopicIDs()
	if err := pubsubtools.EnsureTopicsExist(ctx, client, topicIDs); err != nil {
		return nil, fmt.Errorf(
			"error when ensuring that pubsub topics exist: %w", err)
	}

	subscriptionIDs := pubsubtools.SubscriptionIDs(topicIDs)
	if err := pubsubtools.EnsureSubscriptionsExist(
		ctx,
		client,
		subscriptionIDs,
		ns.callbackURL,
	); err != nil {
		return nil, fmt.Errorf(
			"error when ensuring that pubsub subscriptions exist: %w", err)
	}
	return ns, nil
}

// PubSubNotificationService sends "real" (production) notifications
type PubSubNotificationService struct {
	client      *pubsub.Client
	environment string
	callbackURL string
}

func (ps PubSubNotificationService) checkPreconditions() error {
	if ps.client == nil {
		return fmt.Errorf("precondition check failed, nil pubsub client")
	}

	if ps.environment == "" {
		return fmt.Errorf("blank environment in notification service")
	}

	if ps.callbackURL == "" {
		return fmt.Errorf("blank callback URL in notification service")
	}

	return nil
}

// Notify sends a notification to the specified topic.
// A search engine index job can be one of the listeners on this channel.
func (ps PubSubNotificationService) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour base.Flavour,
	el base.Element,
	metadata map[string]interface{},
) error {
	ctx, span := tracer.Start(ctx, "Notify")
	defer span.End()
	if err := ps.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"pubsub service precondition check failed when notifying: %w", err)
	}

	if el == nil {
		return fmt.Errorf("can't publish nil element")
	}

	payload, err := el.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("validation of element failed: %w", err)
	}

	envelope := dto.NotificationEnvelope{
		UID:      uid,
		Flavour:  flavour,
		Payload:  payload,
		Metadata: metadata,
	}
	envelopePayload, err := json.Marshal(envelope)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"can't marshal notification envelope to JSON: %w", err)
	}

	return pubsubtools.PublishToPubsub(
		ctx,
		ps.client,
		topicID,
		ps.environment,
		helpers.ServiceName,
		helpers.TopicVersion,
		envelopePayload,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps PubSubNotificationService) TopicIDs() []string {
	return []string{
		helpers.AddPubSubNamespace(common.ItemPublishTopic),
		helpers.AddPubSubNamespace(ps.environment),
		helpers.AddPubSubNamespace(common.ItemResolveTopic),
		helpers.AddPubSubNamespace(common.ItemUnresolveTopic),
		helpers.AddPubSubNamespace(common.ItemHideTopic),
		helpers.AddPubSubNamespace(common.ItemShowTopic),
		helpers.AddPubSubNamespace(common.ItemPinTopic),
		helpers.AddPubSubNamespace(common.ItemUnpinTopic),
		helpers.AddPubSubNamespace(common.NudgePublishTopic),
		helpers.AddPubSubNamespace(common.NudgeDeleteTopic),
		helpers.AddPubSubNamespace(common.NudgeResolveTopic),
		helpers.AddPubSubNamespace(common.NudgeUnresolveTopic),
		helpers.AddPubSubNamespace(common.NudgeHideTopic),
		helpers.AddPubSubNamespace(common.NudgeShowTopic),
		helpers.AddPubSubNamespace(common.ActionPublishTopic),
		helpers.AddPubSubNamespace(common.ActionDeleteTopic),
		helpers.AddPubSubNamespace(common.MessagePostTopic),
		helpers.AddPubSubNamespace(common.MessageDeleteTopic),
		helpers.AddPubSubNamespace(common.IncomingEventTopic),
		helpers.AddPubSubNamespace(common.SentEmailTopic),
	}
}

// SubscriptionIDs ...
// TODO Implement this
func (ps PubSubNotificationService) SubscriptionIDs() map[string]string {
	return nil
}

// ReverseSubscriptionIDs ...
// TODO implement this
func (ps PubSubNotificationService) ReverseSubscriptionIDs() map[string]string {
	return nil
}
