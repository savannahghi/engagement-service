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
const ackDeadlineSeconds = 60

func topicIDs() []string {
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

// NewPubSubNotificationService initializes a live notification service
func NewPubSubNotificationService(
	ctx context.Context,
	projectID string,
) (feed.NotificationService, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	ns := &PubSubNotificationService{
		client: client,
	}
	if err := ns.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"pubsub notification service failed preconditions: %w", err)
	}
	if err := ns.ensureTopicsExist(ctx); err != nil {
		return nil, fmt.Errorf(
			"error when ensuring that pubsub topics exist: %w", err)
	}
	return ns, nil
}

// PubSubNotificationService sends "real" (production) notifications
type PubSubNotificationService struct {
	client *pubsub.Client
}

func (ps PubSubNotificationService) ensureTopicsExist(
	ctx context.Context,
) error {
	if ps.client == nil {
		return fmt.Errorf("precondition check failed, nil pubsub client")
	}

	// get a list of configured topic IDs from the project
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
	// and that each topic has at least one subscription by the same name
	desiredTopicIDs := topicIDs()
	for _, topicID := range desiredTopicIDs {
		if !base.StringSliceContains(configuredTopics, topicID) {
			topic, err := ps.client.CreateTopic(ctx, topicID)
			if err != nil {
				return fmt.Errorf("can't create topic %s: %w", topicID, err)
			}

			_, err = ps.client.CreateSubscription(
				ctx,
				topicID,
				pubsub.SubscriptionConfig{
					Topic:               topic,
					AckDeadline:         ackDeadlineSeconds * time.Second,
					RetainAckedMessages: true,
					ExpirationPolicy:    time.Duration(0), // never expire
				})

			if err != nil {
				return fmt.Errorf(
					"can't create subscription %s: %w", topicID, err)
			}
		}
	}

	return nil
}

func (ps PubSubNotificationService) checkPreconditions() error {
	if ps.client == nil {
		return fmt.Errorf("precondition check failed, nil pubsub client")
	}
	return nil
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
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	log.Printf("published %T to %s, got ID %s", el, topicID, id)

	return nil
}
