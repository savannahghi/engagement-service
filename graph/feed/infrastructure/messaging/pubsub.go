package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	"google.golang.org/api/iterator"
)

// messaging related constants
const (
	PubSubHandlerPath = "/pubsub"

	ackDeadlineSeconds  = 60
	maxBackoffSeconds   = 600
	maxDeliveryAttempts = 100 // go to the dead letter topic after this
	hoursInAWeek        = 24 * 7

	defaultPubsubTokenAudience = "bewell.co.ke"
	hostNameEnvVarName         = "SERVICE_HOST" // host at which this service is deployed
	serviceName                = "feed"
	subscriptionVersion        = "v6"
)

// PubSubMessage is a pub-sub message payload
type PubSubMessage struct {
	MessageID  string            `json:"messageId"`
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

// PubSubPayload is the payload of a Pub/Sub event.
type PubSubPayload struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

// NewPubSubNotificationService initializes a live notification service
func NewPubSubNotificationService(
	ctx context.Context,
	projectID string,
	projectNumber int,
) (feed.NotificationService, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	environment, err := base.GetEnvVar(base.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to get the environment variable `%s`: %w", base.Environment, err)
	}

	hostName, err := base.GetEnvVar(hostNameEnvVarName)
	if err != nil {
		return nil, fmt.Errorf("unable to get the %s environment variable: %w", hostNameEnvVarName, err)
	}

	callbackURL := fmt.Sprintf("%s/%s", hostName, PubSubHandlerPath)
	ns := &PubSubNotificationService{
		client:        client,
		environment:   environment,
		projectNumber: projectNumber,
		callbackURL:   callbackURL,
	}
	if err := ns.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"pubsub notification service failed preconditions: %w", err)
	}
	if err := ns.ensureTopicsExist(ctx); err != nil {
		return nil, fmt.Errorf(
			"error when ensuring that pubsub topics exist: %w", err)
	}
	if err := ns.ensureSubscriptionsExist(ctx); err != nil {
		return nil, fmt.Errorf(
			"error when ensuring that pubsub subscriptions exist: %w", err)
	}
	return ns, nil
}

// PubSubNotificationService sends "real" (production) notifications
type PubSubNotificationService struct {
	client        *pubsub.Client
	environment   string
	callbackURL   string
	projectNumber int
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

	if ps.projectNumber == 0 {
		return fmt.Errorf("project number is 0 in the notification service (invalid)")
	}

	return nil
}

func (ps PubSubNotificationService) ensureTopicsExist(
	ctx context.Context,
) error {
	// get a list of configured topic IDs from the project so that we don't recreate
	configuredTopics := []string{}
	it := ps.client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf(
				"error while iterating through pubsub topics: %w", err)
		}
		configuredTopics = append(configuredTopics, topic.ID())
	}

	// ensure that all our desired topics are all created
	for _, topicID := range ps.TopicIDs() {
		if !base.StringSliceContains(configuredTopics, topicID) {
			_, err := ps.client.CreateTopic(ctx, topicID)
			if err != nil {
				return fmt.Errorf("can't create topic %s: %w", topicID, err)
			}
		}
	}

	return nil
}

func (ps PubSubNotificationService) ensureSubscriptionsExist(
	ctx context.Context,
) error {
	subscriptionIDs := ps.SubscriptionIDs()
	for _, topicID := range ps.TopicIDs() {
		topic := ps.client.Topic(topicID)
		topicExists, err := topic.Exists(ctx)
		if err != nil {
			return fmt.Errorf("error when checking if topic %s exists: %w", topicID, err)
		}

		if !topicExists {
			return fmt.Errorf("no topic with ID %s exists", topicID)
		}

		subscriptionConfig, err := ps.getSubscriptionConfig(ctx, topicID)
		if err != nil {
			return fmt.Errorf(
				"can't initialize subscription config for topic %s: %w", topicID, err)
		}

		if subscriptionConfig == nil {
			return fmt.Errorf("nil subscription config")
		}

		subscriptionID, prs := subscriptionIDs[topicID]
		if !prs {
			return fmt.Errorf("no subscriptionID found in map for topic %s", topicID)
		}

		existingSubscription := ps.client.Subscription(subscriptionID)
		subscriptionExists, err := existingSubscription.Exists(ctx)
		if err != nil {
			return fmt.Errorf("error when checking if a subscription exists: %w", err)
		}
		if !subscriptionExists {
			_, err = ps.client.CreateSubscription(ctx, subscriptionID, *subscriptionConfig)
			if err != nil {
				log.Printf("Detailed error:\n%#v\n", err)

				return fmt.Errorf("can't create subscription %s: %w", topicID, err)
			}
		}
	}

	return nil
}

func (ps PubSubNotificationService) getSubscriptionConfig(
	ctx context.Context, topicID string,
) (*pubsub.SubscriptionConfig, error) {
	topic := ps.client.Topic(topicID)
	topicExists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when checking if topic %s exists: %w", topicID, err)
	}

	if !topicExists {
		return nil, fmt.Errorf("no topic with ID %s exists", topicID)
	}

	// This is a PUSH type subscription, because Cloud Run is a *serverless*
	// platform and we cannot keep long lived pull subscriptions there. In a
	// future where this service is no longer run on a serverless platform, we
	// should switch to the higher throughput pull subscriptions.
	//
	// Authentication is via Google signed OpenID Connect tokens. For the Cloud
	// Run deployment, this authentication is automatic (done by Google). If we
	// move this deployment to another environment, we have to do our own
	// verification in the HTTP handler.
	return &pubsub.SubscriptionConfig{
		Topic: topic,
		PushConfig: pubsub.PushConfig{
			Endpoint: ps.callbackURL,
			AuthenticationMethod: &pubsub.OIDCToken{
				Audience: defaultPubsubTokenAudience,
				ServiceAccountEmail: fmt.Sprintf(
					"%d-compute@developer.gserviceaccount.com", ps.projectNumber),
			},
		},
		AckDeadline:         ackDeadlineSeconds * time.Second,
		RetainAckedMessages: true,
		RetentionDuration:   time.Hour * hoursInAWeek,
		ExpirationPolicy:    time.Duration(0), // never expire
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: time.Second,
			MaximumBackoff: time.Second * maxBackoffSeconds,
		},
	}, nil
}

// Notify sends a notification to the specified topic.
// A search engine index job can be one of the listeners on this channel.
func (ps PubSubNotificationService) Notify(
	ctx context.Context,
	topicID string,
	el feed.Element,
) error {
	if err := ps.checkPreconditions(); err != nil {
		return fmt.Errorf(
			"pubsub service precondition check failed when notifying: %w", err)
	}

	if el == nil {
		return fmt.Errorf("can't publish nil element")
	}

	payload, err := el.ValidateAndMarshal()
	if err != nil {
		return fmt.Errorf("validation of element failed: %w", err)
	}

	t := ps.client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: payload,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err = result.Get(ctx) // message id ignored for now
	if err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	return nil
}

// TopicIDs returns the known (registered) topic IDs
func (ps PubSubNotificationService) TopicIDs() []string {
	return []string{
		feed.FeedRetrievalTopic,
		feed.ThinFeedRetrievalTopic,
		feed.ItemRetrievalTopic,
		feed.ItemPublishTopic,
		feed.ItemDeleteTopic,
		feed.ItemResolveTopic,
		feed.ItemUnresolveTopic,
		feed.ItemHideTopic,
		feed.ItemShowTopic,
		feed.ItemPinTopic,
		feed.ItemUnpinTopic,
		feed.NudgeRetrievalTopic,
		feed.NudgePublishTopic,
		feed.NudgeDeleteTopic,
		feed.NudgeResolveTopic,
		feed.NudgeUnresolveTopic,
		feed.NudgeHideTopic,
		feed.NudgeShowTopic,
		feed.ActionRetrievalTopic,
		feed.ActionPublishTopic,
		feed.ActionDeleteTopic,
		feed.MessagePostTopic,
		feed.MessageDeleteTopic,
		feed.IncomingEventTopic,
	}
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (ps PubSubNotificationService) SubscriptionIDs() map[string]string {
	output := map[string]string{}
	for _, topicID := range ps.TopicIDs() {
		subscriptionID := fmt.Sprintf(
			"%s-%s-%s-%s",
			serviceName,
			topicID,
			ps.environment,
			subscriptionVersion,
		)
		output[topicID] = subscriptionID
	}
	return output
}

// ReverseSubscriptionIDs returns a (reversed) map of subscription IDs to topicIDs
func (ps PubSubNotificationService) ReverseSubscriptionIDs() map[string]string {
	output := map[string]string{}
	for _, topicID := range ps.TopicIDs() {
		subscriptionID := fmt.Sprintf(
			"%s-%s-%s-%s",
			serviceName,
			topicID,
			ps.environment,
			subscriptionVersion,
		)
		output[subscriptionID] = topicID
	}
	return output
}
