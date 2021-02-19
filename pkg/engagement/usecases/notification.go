package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
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
	nudgeDeleteSender    = "NUDGE_DELETED"
	nudgeResolveSender   = "NUDGE_RESOLVED"
	nudgeUnresolveSender = "NUDGE_UNRESOLVED"
	nudgeShowSender      = "NUDGE_SHOW"
	nudgeHideSender      = "NUDGE_HIDE"

	feedUpdate       = "FEED_UPDATE"
	inboxCountUpdate = "INBOX_COUNT_CHANGED"
)

// NotificationUsecases represent logic required to make notification
type NotificationUsecases interface {
	HandleItemPublish(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemDelete(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemResolve(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemUnresolve(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemHide(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemShow(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemPin(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleItemUnpin(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgePublish(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgeDelete(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgeResolve(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgeUnresolve(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgeHide(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleNudgeShow(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleActionPublish(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleActionDelete(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleMessagePost(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleMessageDelete(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	HandleIncomingEvent(
		ctx context.Context,
		m *base.PubSubPayload,
	) error

	NotifyItemUpdate(
		ctx context.Context,
		sender string,
		includeNotification bool, // whether to show a tray notification
		m *base.PubSubPayload,
	) error

	UpdateInbox(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) error

	NotifyNudgeUpdate(
		ctx context.Context,
		sender string,
		m *base.PubSubPayload,
	) error

	NotifyInboxCountUpdate(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		count int,
	) error

	GetUserTokens(
		uids []string,
	) ([]string, error)

	SendNotificationViaFCM(
		ctx context.Context,
		uids []string,
		sender string,
		pl resources.NotificationEnvelope,
		notification *base.FirebaseSimpleNotificationInput,
	) error
}

// HandlePubsubPayload defines the signature of a function that handles
// payloads received from Google Cloud Pubsub
type HandlePubsubPayload func(ctx context.Context, m *base.PubSubPayload) error

// NotificationImpl represents the notification usecase implementation
type NotificationImpl struct {
	repository repository.Repository
	fcm        fcm.PushService
	onboarding onboarding.ProfileService
}

// NewNotification initializes a notification usecase
func NewNotification(
	repository repository.Repository,
	fcm fcm.PushService,
	onboarding onboarding.ProfileService,
) *NotificationImpl {
	return &NotificationImpl{
		repository: repository,
		fcm:        fcm,
		onboarding: onboarding,
	}
}

// HandleItemPublish responds to item publish messages
func (n NotificationImpl) HandleItemPublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	err := n.NotifyItemUpdate(ctx, itemPublishSender, true, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemDelete responds to item delete messages
func (n NotificationImpl) HandleItemDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemDeleteSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemResolve responds to item resolve messages
func (n NotificationImpl) HandleItemResolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemResolveSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnresolve responds to item unresolve messages
func (n NotificationImpl) HandleItemUnresolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemUnresolveSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemHide responds to item hide messages
func (n NotificationImpl) HandleItemHide(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemHideSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemShow responds to item show messages
func (n NotificationImpl) HandleItemShow(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemShowSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemPin responds to item pin messages
func (n NotificationImpl) HandleItemPin(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemPinSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnpin responds to item unpin messages
func (n NotificationImpl) HandleItemUnpin(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemUnpinSender, false, m)
	if err != nil {
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleNudgePublish responds to nudge publish messages
func (n NotificationImpl) HandleNudgePublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgePublishSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeDelete responds to nudge delete messages
func (n NotificationImpl) HandleNudgeDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeDeleteSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeResolve responds to nudge resolve messages
func (n NotificationImpl) HandleNudgeResolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeResolveSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeUnresolve responds to nudge unresolve messages
func (n NotificationImpl) HandleNudgeUnresolve(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeUnresolveSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeHide responds to nudge hide messages
func (n NotificationImpl) HandleNudgeHide(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeHideSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeShow responds to nudge hide messages
func (n NotificationImpl) HandleNudgeShow(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeShowSender, m)
	if err != nil {
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleActionPublish responds to action publish messages
func (n NotificationImpl) HandleActionPublish(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify action publish

	return nil
}

// HandleActionDelete responds to action publish messages
func (n NotificationImpl) HandleActionDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify action delete

	return nil
}

// HandleMessagePost responds to message post pubsub messages
func (n NotificationImpl) HandleMessagePost(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify the message and it's context i.e item and feed flavour

	return nil
}

// HandleMessageDelete responds to message delete pubsub messages
func (n NotificationImpl) HandleMessageDelete(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify the message delete and it's context i.e item and feed flavour

	return nil
}

// HandleIncomingEvent responds to message delete pubsub messages
func (n NotificationImpl) HandleIncomingEvent(ctx context.Context, m *base.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	log.Printf("incoming event data: \n%s\n", string(m.Message.Data))
	log.Printf("incoming event subscription: %s", m.Subscription)
	log.Printf("incoming event message ID: %s", m.Message.MessageID)
	log.Printf("incoming event message attributes: %#v", m.Message.Attributes)

	return nil
}

// NotifyItemUpdate sends a Firebase Cloud Messaging notification
func (n NotificationImpl) NotifyItemUpdate(
	ctx context.Context,
	sender string,
	includeNotification bool, // whether to show a tray notification
	m *base.PubSubPayload,
) error {

	var envelope resources.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		return fmt.Errorf(
			"can't unmarshal notification envelope from pubsub data: %w", err)
	}

	var item base.Item
	err = json.Unmarshal(envelope.Payload, &item)
	if err != nil {
		return fmt.Errorf("can't unmarshal item from pubsub data: %w", err)
	}
	// include notifications for persistent items
	var notification *base.FirebaseSimpleNotificationInput
	iconURL := common.DefaultIconPath
	if item.Persistent && includeNotification {
		// also include a notification
		notification = &base.FirebaseSimpleNotificationInput{
			Title:    item.Tagline,
			Body:     item.Summary,
			ImageURL: &iconURL,
		}
	}

	err = n.SendNotificationViaFCM(ctx, item.Users, sender, envelope, notification)
	if err != nil {
		return fmt.Errorf("unable to notify item: %w", err)
	}

	// TODO Send email notifications
	// TODO For urgent (tray), consider whatsapp and sms notifications

	switch sender {
	case itemPublishSender:
		existingLabels, err := n.repository.Labels(ctx, envelope.UID, envelope.Flavour)
		if err != nil {
			return fmt.Errorf("can't fetch existing labels: %w", err)
		}

		if !base.StringSliceContains(existingLabels, item.Label) {
			err = n.repository.SaveLabel(ctx, envelope.UID, envelope.Flavour, item.Label)
			if err != nil {
				return fmt.Errorf("can't save label: %w", err)
			}
		}
	case itemDeleteSender,
		itemResolveSender,
		itemUnresolveSender,
		itemHideSender,
		itemShowSender,
		itemPinSender,
		itemUnpinSender:
		// do nothing...inbox update code will run in the outer scope
	default:
		return fmt.Errorf("unexpected item publish sender: %s", sender)
	}

	err = n.UpdateInbox(ctx, envelope.UID, envelope.Flavour)
	if err != nil {
		return fmt.Errorf("unable to update inbox count: %w", err)
	}

	return nil
}

// UpdateInbox recalculates the inbox count and notifies the client over FCM
func (n NotificationImpl) UpdateInbox(ctx context.Context, uid string, flavour base.Flavour) error {
	err := n.repository.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
	if err != nil {
		return fmt.Errorf("can't update inbox count: %w", err)
	}

	unread, err := n.repository.UnreadPersistentItems(ctx, uid, flavour)
	if err != nil {
		return fmt.Errorf("can't get inbox count: %w", err)
	}

	err = n.NotifyInboxCountUpdate(ctx, uid, flavour, unread)
	if err != nil {
		return fmt.Errorf("can't notify inbox count: %w", err)
	}

	return nil
}

// NotifyNudgeUpdate sends a nudge update notification via FCM
func (n NotificationImpl) NotifyNudgeUpdate(
	ctx context.Context,
	sender string,
	m *base.PubSubPayload,
) error {
	var envelope resources.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		return fmt.Errorf("can't unmarshal notification envelope from pubsub data: %w", err)
	}

	var nudge base.Nudge
	err = json.Unmarshal(envelope.Payload, &nudge)
	if err != nil {
		return fmt.Errorf("can't unmarshal nudge from pubsub data: %w", err)
	}

	err = n.SendNotificationViaFCM(ctx, nudge.Users, sender, envelope, nil)
	if err != nil {
		return fmt.Errorf("unable to notify nudge: %w", err)
	}

	return nil
}

// NotifyInboxCountUpdate sends a message notifying of an update to inbox
// item counts.
func (n NotificationImpl) NotifyInboxCountUpdate(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	count int,
) error {
	notificationEnvelope := resources.NotificationEnvelope{
		UID:     uid,
		Flavour: flavour,
		Payload: []byte(fmt.Sprintf("%d", count)),
		Metadata: map[string]interface{}{
			"sender": inboxCountUpdate,
			"count":  count,
		},
	}

	notifyUIDs := []string{uid}
	err := n.SendNotificationViaFCM(
		ctx, notifyUIDs, feedUpdate, notificationEnvelope, nil)
	if err != nil {
		return fmt.Errorf("unable to notify thin feed: %w", err)
	}

	return nil
}

// GetUserTokens retrieves the user tokens corresponding to the supplied UIDs
func (n NotificationImpl) GetUserTokens(uids []string) ([]string, error) {
	userTokens, err := n.onboarding.GetDeviceTokens(onboarding.UserUIDs{
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

// SendNotificationViaFCM publishes an FCM notification
func (n NotificationImpl) SendNotificationViaFCM(
	ctx context.Context,
	uids []string,
	sender string,
	pl resources.NotificationEnvelope,
	notification *base.FirebaseSimpleNotificationInput,
) error {
	tokens, err := n.GetUserTokens(uids)
	if err != nil {
		return fmt.Errorf("can't get user tokens: %w", err)
	}
	if len(tokens) == 0 {
		return nil
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
	err = n.fcm.Push(ctx, sender, payload)
	if err != nil {
		return fmt.Errorf("can't send element over FCM: %w", err)
	}
	return nil
}
