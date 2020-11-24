package pubsubhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging"
)

const (
	itemPublishSender   = "ITEM_PUBLISHED"
	itemDeleteSender    = "ITEM_DELETED"
	itemResolveSender   = "ITEM_RESOLVED"
	itemUnresolveSender = "ITEM_UNRESOLVED"
	itemHideSender      = "ITEM_HIDE"
	itemShowSender      = "ITEM_SHOW"
	itemPinSender       = "ITEM_PIN"
	itemUnpinSender     = "ITEM_UNPIN"

	nudgePublishSender   = "NUDGE_PUBLISHED"
	nudgeResolveSender   = "NUDGE_RESOLVED"
	nudgeUnresolveSender = "NUDGE_UNRESOLVED"
	nudgeShowSender      = "NUDGE_SHOW"
	nudgeHideSender      = "NUDGE_HIDE"

	feedUpdate       = "FEED_UPDATE"
	inboxCountUpdate = "INBOX_COUNT_CHANGED"
)

// HandlePubsubPayload defines the signature of a function that handles
// payloads received from Google Cloud Pubsub
type HandlePubsubPayload func(ctx context.Context, m *base.PubSubPayload) error

// HandleFeedRetrieval responds to feed retrieval messages
func HandleFeedRetrieval(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	return nil
}

// HandleThinFeedRetrieval responds to thin feed retrieval messages
func HandleThinFeedRetrieval(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	return nil
}

// HandleItemRetrieval responds to item retrieval messages
func HandleItemRetrieval(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	return nil
}

// HandleItemPublish responds to item publish messages
func HandleItemPublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemPublishSender, true, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemDelete responds to item delete messages
func HandleItemDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemDeleteSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemResolve responds to item resolve messages
func HandleItemResolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemResolveSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnresolve responds to item unresolve messages
func HandleItemUnresolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemUnresolveSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemHide responds to item hide messages
func HandleItemHide(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemHideSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemShow responds to item show messages
func HandleItemShow(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemShowSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemPin responds to item pin messages
func HandleItemPin(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemPinSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnpin responds to item unpin messages
func HandleItemUnpin(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyItemUpdate(ctx, itemUnpinSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeRetrieval responds to nudge retrieval messages
func HandleNudgeRetrieval(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	return nil
}

// HandleNudgePublish responds to nudge publish messages
func HandleNudgePublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyNudgeUpdate(ctx, nudgePublishSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeDelete responds to nudge delete messages
func HandleNudgeDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyFeedUpdate(ctx, m)
	if err != nil {
		return fmt.Errorf(
			"can't send thin feed notification over FCM: %w", err)
	}

	return nil
}

// HandleNudgeResolve responds to nudge resolve messages
func HandleNudgeResolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyNudgeUpdate(ctx, nudgeResolveSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeUnresolve responds to nudge unresolve messages
func HandleNudgeUnresolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyNudgeUpdate(ctx, nudgeUnresolveSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeHide responds to nudge hide messages
func HandleNudgeHide(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyNudgeUpdate(ctx, nudgeHideSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeShow responds to nudge hide messages
func HandleNudgeShow(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyNudgeUpdate(ctx, nudgeShowSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleActionRetrieval responds to action retrieval messages
func HandleActionRetrieval(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	return nil
}

// HandleActionPublish responds to action publish messages
func HandleActionPublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyFeedUpdate(ctx, m)
	if err != nil {
		return fmt.Errorf(
			"can't send thin feed notification over FCM: %w", err)
	}

	return nil
}

// HandleActionDelete responds to action publish messages
func HandleActionDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyFeedUpdate(ctx, m)
	if err != nil {
		return fmt.Errorf(
			"can't send thin feed notification over FCM: %w", err)
	}

	return nil
}

// HandleMessagePost responds to message post pubsub messages
func HandleMessagePost(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyFeedUpdate(ctx, m)
	if err != nil {
		return fmt.Errorf("can't notify feed update over FCM: %w", err)
	}

	return nil
}

// HandleMessageDelete responds to message delete pubsub messages
func HandleMessageDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := notifyFeedUpdate(ctx, m)
	if err != nil {
		return fmt.Errorf("can't notify feed update over FCM: %w", err)
	}

	return nil
}

// HandleIncomingEvent responds to message delete pubsub messages
func HandleIncomingEvent(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	log.Printf("incoming event data: \n%s\n", string(m.Message.Data))
	log.Printf("incoming event subscription: %s", m.Subscription)
	log.Printf("incoming event message ID: %s", m.Message.MessageID)
	log.Printf("incoming event message attributes: %#v", m.Message.Attributes)

	return nil
}

func notifyItemUpdate(
	ctx context.Context,
	sender string,
	includeNotification bool, // whether to show a tray notification
	m *base.PubSubPayload,
) error {
	var envelope feed.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		return fmt.Errorf(
			"can't unmarshal notification envelop from pubsub data: %w", err)
	}

	var item feed.Item
	err = json.Unmarshal(envelope.Payload, &item)
	if err != nil {
		return fmt.Errorf("can't unmarshal item from pubsub data: %w", err)
	}

	// include notifications for persistent items
	var notification *base.FirebaseSimpleNotificationInput
	iconURL := feed.DefaultIconPath
	if item.Persistent && includeNotification {
		// also include a notification
		notification = &base.FirebaseSimpleNotificationInput{
			Title:    item.Tagline,
			Body:     item.Summary,
			ImageURL: &iconURL,
		}
	}

	err = sendNotificationViaFCM(ctx, item.Users, sender, envelope, notification)
	if err != nil {
		return fmt.Errorf("unable to notify item: %w", err)
	}

	switch sender {
	case itemPublishSender:
		thinFeed, err := getThinFeed(ctx, envelope.UID, envelope.Flavour)
		if err != nil {
			return fmt.Errorf("can't instantiate new feed: %w", err)
		}

		existingLabels, err := thinFeed.Labels(ctx)
		if err != nil {
			return fmt.Errorf("can't fetch existing labels: %w", err)
		}

		if !base.StringSliceContains(existingLabels, item.Label) {
			err = thinFeed.SaveLabel(ctx, item.Label)
			if err != nil {
				return fmt.Errorf("can't save label: %w", err)
			}
		}

		err = updateInbox(ctx, envelope.UID, envelope.Flavour)
		if err != nil {
			return fmt.Errorf("unable to update inbox count: %w", err)
		}
	case itemDeleteSender,
		itemResolveSender,
		itemUnresolveSender,
		itemHideSender,
		itemShowSender,
		itemPinSender,
		itemUnpinSender:
		err := updateInbox(ctx, envelope.UID, envelope.Flavour)
		if err != nil {
			return fmt.Errorf("unable to update inbox count: %w", err)
		}
	}

	return nil
}

func updateInbox(ctx context.Context, uid string, flavour feed.Flavour) error {
	thinFeed, err := getThinFeed(ctx, uid, flavour)
	if err != nil {
		return fmt.Errorf("can't instantiate new feed: %w", err)
	}

	err = thinFeed.UpdateUnreadPersistentItemsCount(ctx)
	if err != nil {
		return fmt.Errorf("can't update inbox count: %w", err)
	}

	unread, err := thinFeed.UnreadPersistentItems(ctx)
	if err != nil {
		return fmt.Errorf("can't get inbox count: %w", err)
	}

	err = notifyInboxCountUpdate(ctx, uid, flavour, unread)
	if err != nil {
		return fmt.Errorf("can't notify inbox count: %w", err)
	}

	return nil
}

func notifyNudgeUpdate(
	ctx context.Context,
	sender string,
	m *base.PubSubPayload,
) error {
	var envelope feed.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		return fmt.Errorf(
			"can't unmarshal notification envelop from pubsub data: %w", err)
	}

	var nudge feed.Nudge
	err = json.Unmarshal(envelope.Payload, &nudge)
	if err != nil {
		return fmt.Errorf("can't unmarshal nudge from pubsub data: %w", err)
	}

	err = sendNotificationViaFCM(ctx, nudge.Users, sender, envelope, nil)
	if err != nil {
		return fmt.Errorf("unable to notify nudge: %w", err)
	}

	return nil
}

func notifyFeedUpdate(ctx context.Context, m *base.PubSubPayload) error {
	var notificationEnvelope feed.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &notificationEnvelope)
	if err != nil {
		return fmt.Errorf(
			"can't unmarshal notification envelope from pubsub data: %w", err)
	}

	notifyUIDs := []string{notificationEnvelope.UID}
	err = sendNotificationViaFCM(
		ctx, notifyUIDs, feedUpdate, notificationEnvelope, nil)
	if err != nil {
		return fmt.Errorf("unable to notify thin feed: %w", err)
	}

	return nil
}

func notifyInboxCountUpdate(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	count int,
) error {
	notificationEnvelope := feed.NotificationEnvelope{
		UID:     uid,
		Flavour: flavour,
		Payload: []byte(fmt.Sprintf("%d", count)),
		Metadata: map[string]interface{}{
			"sender": inboxCountUpdate,
			"count":  count,
		},
	}

	notifyUIDs := []string{uid}
	err := sendNotificationViaFCM(
		ctx, notifyUIDs, feedUpdate, notificationEnvelope, nil)
	if err != nil {
		return fmt.Errorf("unable to notify thin feed: %w", err)
	}

	return nil
}

func getUserTokens(uids []string) ([]string, error) {
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		return nil, fmt.Errorf("can't load inter-service config from YAML: %w", err)
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		return nil, fmt.Errorf(
			"can't set up profile interservice client: %w", err)
	}
	profileService := feed.NewRemoteProfileService(profileClient)
	userTokens, err := profileService.GetDeviceTokens(feed.UserUIDs{
		UIDs: uids,
	})
	if err != nil {
		return nil, fmt.Errorf("can't get push tokens: %w", err)
	}
	tokens := []string{}
	for _, toks := range userTokens {
		tokens = append(tokens, toks...)
	}
	return tokens, nil
}

func sendNotificationViaFCM(
	ctx context.Context,
	uids []string,
	sender string,
	pl feed.NotificationEnvelope,
	notification *base.FirebaseSimpleNotificationInput,
) error {
	tokens, err := getUserTokens(uids)
	if err != nil {
		return fmt.Errorf("can't get user tokens: %w", err)
	}
	if len(tokens) == 0 {
		return nil
	}
	pushService, err := feed.NewRemoteFCMPushService(ctx)
	if err != nil {
		log.Fatalf("unable to initialize push service: %s", err)
	}
	marshalled, err := json.Marshal(pl)
	if err != nil {
		return fmt.Errorf(
			"can't send element that failed validation over FCM: %w", err)
	}
	payload := base.SendNotificationPayload{
		RegistrationTokens: tokens,
		Data: map[string]string{
			sender: string(marshalled),
		},
	}
	if notification != nil {
		payload.Notification = notification
	}
	err = pushService.Push(ctx, sender, payload)
	if err != nil {
		return fmt.Errorf("can't send element over FCM: %w", err)
	}
	return nil
}

func getThinFeed(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
) (*feed.Feed, error) {
	repository, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate new firebase repository: %w", err)
	}

	notificationService, err := messaging.NewMockNotificationService()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate notification service")
	}

	agg, err := feed.NewCollection(repository, notificationService)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize feed collection: %w", err)
	}
	thinFeed, err := agg.GetThinFeed(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate thin feed: %w", err)
	}
	return thinFeed, nil
}
