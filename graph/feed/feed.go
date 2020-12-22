package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
)

// feed related constants
const (
	// these topic names are exported because they are also used by sub-packages
	ItemPublishTopic    = "items.publish"
	ItemDeleteTopic     = "items.delete"
	ItemResolveTopic    = "items.resolve"
	ItemUnresolveTopic  = "items.unresolve"
	ItemHideTopic       = "items.hide"
	ItemShowTopic       = "items.show"
	ItemPinTopic        = "items.pin"
	ItemUnpinTopic      = "items.unpin"
	NudgePublishTopic   = "nudges.publish"
	NudgeDeleteTopic    = "nudges.delete"
	NudgeResolveTopic   = "nudges.resolve"
	NudgeUnresolveTopic = "nudges.unresolve"
	NudgeHideTopic      = "nudges.hide"
	NudgeShowTopic      = "nudges.show"
	ActionPublishTopic  = "actions.publish"
	ActionDeleteTopic   = "actions.delete"
	MessagePostTopic    = "message.post"
	MessageDeleteTopic  = "message.delete"
	IncomingEventTopic  = "incoming.event"
	ServiceName         = "feed"
	TopicVersion        = "v1"
)

// ErrNilNudge is a marker error for cases when a nudge should have
// been found but was not
var ErrNilNudge = fmt.Errorf("nil nudge")

// ErrNilFeedItem is a sentinel error used to indicate when a feed item
// should have been non nil but was not
var ErrNilFeedItem = fmt.Errorf("nil feed item")

// Feed manages and serializes the nudges, actions and feed items that a
// specific user should see.
//
// A feed is stored and serialized on a per-user basis. If a feed item is sent
// to a group of users, it should be "expanded" before the user's feed gets
// stored.
type Feed struct {
	repository          Repository
	notificationService NotificationService

	// a string composed by concatenating the UID, a "|" and a flavour
	ID string `json:"id" firestore:"-"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// user identifier - who does this feed belong to?
	// this is also the unique identifier for a feed
	UID string `json:"uid" firestore:"uid"`

	// whether this is a consumer or pro feed
	Flavour base.Flavour `json:"flavour" firestore:"flavour"`

	// what are the global actions available to this user?
	Actions []base.Action `json:"actions" firestore:"actions"`

	// what does this user's feed contain?
	Items []base.Item `json:"items" firestore:"items"`

	// what prompts or nudges should this user see?
	Nudges []base.Nudge `json:"nudges" firestore:"nudges"`

	// indicates whether the user is Anonymous or not
	IsAnonymous *bool `json:"isAnonymous" firestore:"isAnonymous"`
}

func (fe Feed) getID() string {
	return fmt.Sprintf("%s|%s", fe.UID, fe.Flavour.String())
}

func (fe Feed) checkPreconditions() error {
	if fe.repository == nil {
		return fmt.Errorf("feed has a nil repository")
	}

	if fe.notificationService == nil {
		return fmt.Errorf("feed has a nil notification service")
	}

	if fe.UID == "" {
		return fmt.Errorf("feed has a zero valued UID")
	}

	if !fe.Flavour.IsValid() {
		return fmt.Errorf("feed does not have a valid flavour")
	}

	if fe.ID == "" {
		fe.ID = fe.getID()
	}

	if fe.SequenceNumber == 0 {
		fe.SequenceNumber = int(time.Now().Unix())
	}

	return nil
}

// ValidateElement ensures that an element is non nil and valid
func (fe Feed) ValidateElement(el base.Element) error {
	if el == nil {
		return fmt.Errorf("nil element")
	}

	_, err := el.ValidateAndMarshal()
	if err != nil {
		return fmt.Errorf("element failed validation: %w", err)
	}

	return nil
}

// IsEntity marks a feed as an Apollo federation GraphQL entity
func (fe Feed) IsEntity() {}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (fe *Feed) ValidateAndUnmarshal(b []byte) error {
	err := base.ValidateAndUnmarshal(base.FeedSchemaFile, b, fe)
	if err != nil {
		return fmt.Errorf("invalid feed JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against the JSON schema then marshals to JSON
func (fe *Feed) ValidateAndMarshal() ([]byte, error) {
	return base.ValidateAndMarshal(base.FeedSchemaFile, fe)
}

// GetFeedItem retrieves a feed item
func (fe Feed) GetFeedItem(ctx context.Context, itemID string) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve feed item %s: %w", itemID, err)
	}

	if item == nil {
		return nil, nil
	}

	return item, nil
}

// GetNudge retrieves a feed item
func (fe Feed) GetNudge(ctx context.Context, nudgeID string) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve nudge %s: %w", nudgeID, err)
	}

	if nudge == nil {
		return nil, nil
	}

	return nudge, nil
}

// GetAction retrieves a feed item
func (fe Feed) GetAction(
	ctx context.Context,
	actionID string,
) (*base.Action, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	action, err := fe.repository.GetAction(ctx, fe.UID, fe.Flavour, actionID)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve action %s: %w", actionID, err)
	}

	if action == nil {
		return nil, nil
	}

	return action, nil
}

// PublishFeedItem idempotently creates or updates a feed item
func (fe Feed) PublishFeedItem(
	ctx context.Context,
	item *base.Item,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("cano't publish nil feed item")
	}

	if item.SequenceNumber == 0 {
		item.SequenceNumber = int(time.Now().Unix())
	}

	err := fe.ValidateElement(item)
	if err != nil {
		return nil, fmt.Errorf("invalid item: %w", err)
	}

	for _, action := range item.Actions {
		if action.ActionType == base.ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	item, err = fe.repository.SaveFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to publish feed item %s: %w", item.ID, err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemPublishTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify item to channel: %w", err)
	}

	return item, nil
}

// DeleteFeedItem removes a feed item
func (fe Feed) DeleteFeedItem(
	ctx context.Context,
	itemID string,
) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.GetFeedItem(ctx, itemID)
	if err != nil || item == nil {
		// fails to error because it should be safe to retry deletes
		return nil // does not exist, nothing to delete
	}

	err = fe.repository.DeleteFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return fmt.Errorf("unable to delete item: %s", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemDeleteTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return fmt.Errorf("unable to notify item to channel: %w", err)
	}

	return fe.repository.DeleteFeedItem(ctx, fe.UID, fe.Flavour, itemID)
}

// ResolveFeedItem marks a feed item as Done
func (fe Feed) ResolveFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Status = base.StatusDone
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemResolveTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// PinFeedItem marks a feed item as persistent
func (fe Feed) PinFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Persistent = true
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemResolveTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// UnpinFeedItem marks a feed item as not persistent
func (fe Feed) UnpinFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Persistent = false
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to pin feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemPinTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify pinned item to channel: %w", err)
	}

	return item, nil
}

// UnresolveFeedItem marks a feed item as pending
func (fe Feed) UnresolveFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Status = base.StatusPending
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemUnresolveTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify unresolved item to channel: %w", err)
	}

	return item, nil
}

// HideFeedItem hides a feed item from a specific user's feed
func (fe Feed) HideFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Visibility = base.VisibilityHide
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to hide feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemHideTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify hidden item to channel: %w", err)
	}

	return item, nil
}

// ShowFeedItem hides a feed item from a specific user's feed
func (fe Feed) ShowFeedItem(
	ctx context.Context,
	itemID string,
) (*base.Item, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	item, err := fe.repository.GetFeedItem(ctx, fe.UID, fe.Flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, ErrNilFeedItem
	}

	item.Visibility = base.VisibilityShow
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to show feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ItemShowTopic),
		fe.UID,
		fe.Flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify revealed/shown item to channel: %w", err)
	}

	return item, nil
}

// Labels returns the valid labels / filters for this feed
func (fe Feed) Labels(ctx context.Context) ([]string, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	return fe.repository.Labels(ctx, fe.UID, fe.Flavour)
}

// SaveLabel saves the indicated label, if it does not already exist
func (fe Feed) SaveLabel(ctx context.Context, label string) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	return fe.repository.SaveLabel(ctx, fe.UID, fe.Flavour, label)
}

// UnreadPersistentItems returns the number of unread inbox items for this feed
func (fe Feed) UnreadPersistentItems(ctx context.Context) (int, error) {
	if err := fe.checkPreconditions(); err != nil {
		return -1, fmt.Errorf("feed precondition check failed: %w", err)
	}

	return fe.repository.UnreadPersistentItems(ctx, fe.UID, fe.Flavour)
}

// UpdateUnreadPersistentItemsCount updates the number of unread inbox items
func (fe Feed) UpdateUnreadPersistentItemsCount(ctx context.Context) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	return fe.repository.UpdateUnreadPersistentItemsCount(ctx, fe.UID, fe.Flavour)
}

// PublishNudge idempotently creates or updates a nudge
//
// If a nudge with the same ID existed but the sequence number of the new
// nudge is higher, the nudge is replaced.
//
// If a nudge with that ID does not exist, it is inserted at the correct place.
//
// If a nudge with that ID exists, and the existing sequence number is lower,
// it is updated.
//
// If a nudge with that ID and sequence number already exists, the update is
// ignored. This makes the push method idempotent.
//
// If the nudge does not have a sequence number, it is assigned one.
func (fe Feed) PublishNudge(
	ctx context.Context,
	nudge *base.Nudge,
) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	if nudge == nil {
		return nil, fmt.Errorf("can't publish nil nudge")
	}

	if nudge.SequenceNumber == 0 {
		nudge.SequenceNumber = int(time.Now().Unix())
	}

	err := fe.ValidateElement(nudge)
	if err != nil {
		return nil, fmt.Errorf("invalid nudge: %w", err)
	}

	for _, action := range nudge.Actions {
		if action.ActionType == base.ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	nudge, err = fe.repository.SaveNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to publish nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgePublishTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// ResolveNudge marks a feed item as Done
func (fe Feed) ResolveNudge(
	ctx context.Context,
	nudgeID string,
) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, ErrNilNudge
	}

	nudge.Status = base.StatusDone
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgeResolveTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// UnresolveNudge marks a feed item as pending
func (fe Feed) UnresolveNudge(
	ctx context.Context,
	nudgeID string,
) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, ErrNilNudge
	}

	nudge.Status = base.StatusPending
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgeUnresolveTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// HideNudge hides a feed item from a specific user's feed
func (fe Feed) HideNudge(
	ctx context.Context,
	nudgeID string,
) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, ErrNilNudge
	}

	nudge.Visibility = base.VisibilityHide
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgeHideTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}
	return nudge, nil
}

// ShowNudge hides a feed item from a specific user's feed
func (fe Feed) ShowNudge(ctx context.Context, nudgeID string) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, ErrNilNudge
	}

	nudge.Visibility = base.VisibilityShow
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgeShowTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// DeleteNudge removes a nudge
func (fe Feed) DeleteNudge(ctx context.Context, nudgeID string) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.GetNudge(ctx, nudgeID)
	if err != nil || nudge == nil {
		return nil // no error, "re-deleting" a nudge should not cause an error
	}

	err = fe.repository.DeleteNudge(ctx, fe.UID, fe.Flavour, nudgeID)
	if err != nil {
		return fmt.Errorf("can't delete nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(NudgeDeleteTopic),
		fe.UID,
		fe.Flavour,
		nudge,
		map[string]interface{}{
			"nudgeID": nudge.ID,
		},
	); err != nil {
		return fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nil
}

// PublishAction adds/updates an action in a user's feed
//
// If an action with the same ID existed but the sequence number of the new
// nudge is higher, the action is replaced.
//
// If an action with that ID does not exist, it is inserted at the correct
// place.
//
// If an action with that ID exists, and the existing sequence number is lower,
// it is updated.
//
// If an action with that ID and sequence number already exists, the update is
// ignored. This makes the push method idempotent.
//
// If the action does not have a sequence number, it is assigned one.
func (fe Feed) PublishAction(
	ctx context.Context,
	action *base.Action,
) (*base.Action, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	if action == nil {
		return nil, fmt.Errorf("can't publish nil nudge")
	}

	if action.SequenceNumber == 0 {
		action.SequenceNumber = int(time.Now().Unix())
	}

	err := fe.ValidateElement(action)
	if err != nil {
		return nil, fmt.Errorf("invalid action: %w", err)
	}

	action, err = fe.repository.SaveAction(ctx, fe.UID, fe.Flavour, action)
	if err != nil {
		return nil, fmt.Errorf("unable to publish action: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ActionPublishTopic),
		fe.UID,
		fe.Flavour,
		action,
		map[string]interface{}{
			"actionID": action.ID,
		},
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify action to channel: %w", err)
	}

	return action, nil
}

// DeleteAction removes a nudge
func (fe Feed) DeleteAction(ctx context.Context, actionID string) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	action, err := fe.GetAction(ctx, actionID)
	if err != nil || action == nil {
		return nil // no harm "re-deleting" an already deleted action
	}

	err = fe.repository.DeleteAction(ctx, fe.UID, fe.Flavour, actionID)
	if err != nil {
		return fmt.Errorf("unable to delete action: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(ActionDeleteTopic),
		fe.UID,
		fe.Flavour,
		action,
		map[string]interface{}{
			"actionID": action.ID,
		},
	); err != nil {
		return fmt.Errorf("unable to notify action to channel: %w", err)
	}

	return nil
}

// PostMessage updates a feed/thread with a new message OR a reply
func (fe Feed) PostMessage(
	ctx context.Context,
	itemID string,
	message *base.Message,
) (*base.Message, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	if message == nil {
		return nil, fmt.Errorf("can't post nil message")
	}

	if message.ID == "" {
		message.ID = ksuid.New().String()
	}

	if message.SequenceNumber == 0 {
		message.SequenceNumber = int(time.Now().Unix())
	}

	err := fe.ValidateElement(message)
	if err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	msg, err := fe.repository.PostMessage(
		ctx,
		fe.UID,
		fe.Flavour,
		itemID,
		message,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to post a message: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(MessagePostTopic),
		fe.UID,
		fe.Flavour,
		message,
		map[string]interface{}{
			"itemID":    itemID,
			"messageID": message.ID,
		},
	); err != nil {
		return nil, fmt.Errorf("unable to notify message to channel: %w", err)
	}

	return msg, nil
}

// DeleteMessage permanently removes a message
func (fe Feed) DeleteMessage(
	ctx context.Context,
	itemID string,
	messageID string,
) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	message, err := fe.repository.GetMessage(
		ctx,
		fe.UID,
		fe.Flavour,
		itemID,
		messageID,
	)
	if err != nil || message == nil {
		return nil // no harm "re-deleting" an already deleted message
	}

	err = fe.repository.DeleteMessage(
		ctx,
		fe.UID,
		fe.Flavour,
		itemID,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("unable to delete message: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(MessageDeleteTopic),
		fe.UID,
		fe.Flavour,
		message,
		map[string]interface{}{
			"itemID":    itemID,
			"messageID": message.ID,
		},
	); err != nil {
		return fmt.Errorf("unable to notify message delete to channel: %w", err)
	}

	return nil
}

// ProcessEvent publishes an event to an incoming event channel.
//
// Further processing is delegated to listeners to that event channel.
//
// The results of processing an event include but are not limited to:
//
// 	1. Marking feed items as done and notifying their subscribers
//  2. Marking nudges as done and notifying their subscribers
//  3. Updating an audit trail
//  4. Updating (streaming) analytics
func (fe Feed) ProcessEvent(
	ctx context.Context,
	event *base.Event,
) error {
	if err := fe.checkPreconditions(); err != nil {
		return fmt.Errorf("feed precondition check failed: %w", err)
	}

	if event == nil {
		return fmt.Errorf("can't process nil event")
	}

	if !event.Context.Flavour.IsValid() {
		event.Context.Flavour = fe.Flavour
	}

	if event.ID == "" {
		event.ID = ksuid.New().String()
	}

	if event.Context.UserID == "" {
		event.Context.UserID = fe.UID
	}

	err := fe.ValidateElement(event)
	if err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	if event.Context.Flavour != fe.Flavour {
		return fmt.Errorf(
			"the event context flavour (%s) does not match the feed flavour (%s)",
			event.Context.Flavour,
			fe.Flavour,
		)
	}

	err = fe.repository.SaveIncomingEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("can't save incoming event: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		AddPubSubNamespace(IncomingEventTopic),
		fe.UID,
		fe.Flavour,
		event,
		map[string]interface{}{
			"eventID": event.ID,
		},
	); err != nil {
		return fmt.Errorf(
			"unable to publish incoming event to channel: %w", err)
	}

	return nil
}

// FilterParams organizes the parameters needed to filter persistent feed items
type FilterParams struct {
	Labels []string `json:"labels"`
}

// AddPubSubNamespace creates a namespaced topic name
func AddPubSubNamespace(topicName string) string {
	environment := base.GetRunningEnvironment()
	return base.NamespacePubsubIdentifier(
		ServiceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// GetDefaultNudgeByTitle retrieves a default feed nudge
func (fe Feed) GetDefaultNudgeByTitle(ctx context.Context, title string) (*base.Nudge, error) {
	if err := fe.checkPreconditions(); err != nil {
		return nil, fmt.Errorf("feed precondition check failed: %w", err)
	}

	nudge, err := fe.repository.GetDefaultNudgeByTitle(ctx, fe.UID, fe.Flavour, title)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve verify email nudge: %w", err)
	}

	if nudge == nil {
		return nil, fmt.Errorf("can't get the default verify email nudge")
	}

	return nudge, nil
}
