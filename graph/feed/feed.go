package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/segmentio/ksuid"
	"github.com/xeipuuv/gojsonschema"
	"gitlab.slade360emr.com/go/base"
)

// feed related constants
const (
	// these topic names are exported because they are also used by sub-packages
	FeedRetrievalTopic     = "feed.get"
	ThinFeedRetrievalTopic = "thinfeed.get"
	ItemRetrievalTopic     = "items.get"
	ItemPublishTopic       = "items.publish"
	ItemDeleteTopic        = "items.delete"
	ItemResolveTopic       = "items.resolve"
	ItemUnresolveTopic     = "items.unresolve"
	ItemHideTopic          = "items.hide"
	ItemShowTopic          = "items.show"
	ItemPinTopic           = "items.pin"
	ItemUnpinTopic         = "items.unpin"
	NudgeRetrievalTopic    = "nudges.get"
	NudgePublishTopic      = "nudges.publish"
	NudgeDeleteTopic       = "nudges.delete"
	NudgeResolveTopic      = "nudges.resolve"
	NudgeUnresolveTopic    = "nudges.unresolve"
	NudgeHideTopic         = "nudges.hide"
	NudgeShowTopic         = "nudges.show"
	ActionRetrievalTopic   = "actions.get"
	ActionPublishTopic     = "actions.publish"
	ActionDeleteTopic      = "actions.delete"
	MessagePostTopic       = "message.post"
	MessageDeleteTopic     = "message.delete"
	IncomingEventTopic     = "incoming.event"

	LogoURL        = "https://assets.healthcloud.co.ke/bewell_logo.png"
	SampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"

	linkSchemaFile       = "link.schema.json"
	messageSchemaFile    = "message.schema.json"
	actionSchemaFile     = "action.schema.json"
	nudgeSchemaFile      = "nudge.schema.json"
	itemSchemaFile       = "item.schema.json"
	feedSchemaFile       = "feed.schema.json"
	contextSchemaFile    = "context.schema.json"
	payloadSchemaFile    = "payload.schema.json"
	eventSchemaFile      = "event.schema.json"
	statusSchemaFile     = "status.schema.json"
	visibilitySchemaFile = "visibility.schema.json"

	fallbackSchemaHost   = "https://schema.healthcloud.co.ke"
	schemaHostEnvVarName = "SCHEMA_HOST"
)

// ErrNilNudge is a marker error for cases when a nudge should have
// been found but was not
var ErrNilNudge = fmt.Errorf("nil nudge")

// ErrNilFeedItem is a sentinel error used to indicate when a feed item
// should have been non nil but was not
var ErrNilFeedItem = fmt.Errorf("nil feed item")

// Element is a building block of a feed e.g a nudge, action, feed item etc
// An element should know how to validate itself against it's JSON schema
type Element interface {
	ValidateAndUnmarshal(b []byte) error
	ValidateAndMarshal() ([]byte, error)
}

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

	// user identifier - who does this feed belong to?
	// this is also the unique identifier for a feed
	UID string `json:"uid" firestore:"uid"`

	// whether this is a consumer or pro feed
	Flavour Flavour `json:"flavour" firestore:"flavour"`

	// what are the global actions available to this user?
	Actions []Action `json:"actions" firestore:"actions"`

	// what does this user's feed contain?
	Items []Item `json:"items" firestore:"items"`

	// what prompts or nudges should this user see?
	Nudges []Nudge `json:"nudges" firestore:"nudges"`
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

	return nil
}

// ValidateElement ensures that an element is non nil and valid
func (fe Feed) ValidateElement(el Element) error {
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
	err := validateAndUnmarshal(feedSchemaFile, b, fe)
	if err != nil {
		return fmt.Errorf("invalid feed JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against the JSON schema then marshals to JSON
func (fe *Feed) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(feedSchemaFile, fe)
}

// GetFeedItem retrieves a feed item
func (fe Feed) GetFeedItem(ctx context.Context, itemID string) (*Item, error) {
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

	if err := fe.notificationService.Notify(
		ctx, ItemRetrievalTopic, item); err != nil {
		return nil, fmt.Errorf("unable to notify item to channel: %w", err)
	}

	return item, nil
}

// GetNudge retrieves a feed item
func (fe Feed) GetNudge(ctx context.Context, nudgeID string) (*Nudge, error) {
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

	if err := fe.notificationService.Notify(
		ctx, NudgeRetrievalTopic, nudge); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// GetAction retrieves a feed item
func (fe Feed) GetAction(
	ctx context.Context,
	actionID string,
) (*Action, error) {
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

	if err := fe.notificationService.Notify(
		ctx, ActionRetrievalTopic, action); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return action, nil
}

// PublishFeedItem idempotently creates or updates a feed item
func (fe Feed) PublishFeedItem(
	ctx context.Context,
	item *Item,
) (*Item, error) {
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
		if action.ActionType == ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	item, err = fe.repository.SaveFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to publish feed item %s: %w", item.ID, err)
	}

	if err := fe.notificationService.Notify(
		ctx, ItemPublishTopic, item); err != nil {
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
		ctx, ItemDeleteTopic, item); err != nil {
		return fmt.Errorf("unable to notify item to channel: %w", err)
	}

	return fe.repository.DeleteFeedItem(ctx, fe.UID, fe.Flavour, itemID)
}

// ResolveFeedItem marks a feed item as Done
func (fe Feed) ResolveFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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

	item.Status = StatusDone
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, ItemResolveTopic, item); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// PinFeedItem marks a feed item as persistent
func (fe Feed) PinFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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
		ctx, ItemResolveTopic, item); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// UnpinFeedItem marks a feed item as not persistent
func (fe Feed) UnpinFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, ItemResolveTopic, item); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// UnresolveFeedItem marks a feed item as pending
func (fe Feed) UnresolveFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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

	item.Status = StatusPending
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		ItemUnresolveTopic,
		item,
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// HideFeedItem hides a feed item from a specific user's feed
func (fe Feed) HideFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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

	item.Visibility = VisibilityHide
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, ItemHideTopic, item); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
}

// ShowFeedItem hides a feed item from a specific user's feed
func (fe Feed) ShowFeedItem(
	ctx context.Context,
	itemID string,
) (*Item, error) {
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

	item.Visibility = VisibilityShow
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.repository.UpdateFeedItem(ctx, fe.UID, fe.Flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx,
		ItemShowTopic,
		item,
	); err != nil {
		return nil, fmt.Errorf(
			"unable to notify resolved item to channel: %w", err)
	}

	return item, nil
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
	nudge *Nudge,
) (*Nudge, error) {
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
		if action.ActionType == ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	nudge, err = fe.repository.SaveNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to publish nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, NudgePublishTopic, nudge); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// ResolveNudge marks a feed item as Done
func (fe Feed) ResolveNudge(
	ctx context.Context,
	nudgeID string,
) (*Nudge, error) {
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

	nudge.Status = StatusDone
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, NudgeResolveTopic, nudge); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// UnresolveNudge marks a feed item as pending
func (fe Feed) UnresolveNudge(
	ctx context.Context,
	nudgeID string,
) (*Nudge, error) {
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

	nudge.Status = StatusPending
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, NudgeUnresolveTopic, nudge); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}

	return nudge, nil
}

// HideNudge hides a feed item from a specific user's feed
func (fe Feed) HideNudge(
	ctx context.Context,
	nudgeID string,
) (*Nudge, error) {
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

	nudge.Visibility = VisibilityHide
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, NudgeHideTopic, nudge); err != nil {
		return nil, fmt.Errorf("unable to notify nudge to channel: %w", err)
	}
	return nudge, nil
}

// ShowNudge hides a feed item from a specific user's feed
func (fe Feed) ShowNudge(ctx context.Context, nudgeID string) (*Nudge, error) {
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

	nudge.Visibility = VisibilityShow
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.repository.UpdateNudge(ctx, fe.UID, fe.Flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %w", err)
	}

	if err := fe.notificationService.Notify(
		ctx, NudgeShowTopic, nudge); err != nil {
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
		NudgeDeleteTopic,
		nudge,
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
	action *Action,
) (*Action, error) {
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
		ctx, ActionPublishTopic, action); err != nil {
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
		ActionDeleteTopic,
		action,
	); err != nil {
		return fmt.Errorf("unable to notify action to channel: %w", err)
	}

	return nil
}

// PostMessage updates a feed/thread with a new message OR a reply
func (fe Feed) PostMessage(
	ctx context.Context,
	itemID string,
	message *Message,
) (*Message, error) {
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
		MessagePostTopic,
		message,
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
		MessageDeleteTopic,
		message,
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
	event *Event,
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
		IncomingEventTopic,
		event,
	); err != nil {
		return fmt.Errorf(
			"unable to publish incoming event to channel: %w", err)
	}

	return nil
}

// Action represents the global and non-global actions that a user can see/do
type Action struct {
	// A unique identifier for each action
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// A friendly name for the action; rich text with Unicode, can have emoji
	Name string `json:"name" firestore:"name"`

	// Action types are: primary, secondary, overflow and floating
	// Primary actions get dominant visual treatment;
	// secondary actions less so;
	// overflow actions are hidden;
	// floating actions are material FABs
	ActionType ActionType `json:"actionType" firestore:"actionType"`

	// How the action should be handled e.g inline or full page.
	// This is a hint for frontend logic.
	Handling Handling `json:"handling" firestore:"handling"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ac *Action) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(actionSchemaFile, b, ac)
	if err != nil {
		return fmt.Errorf("invalid action JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ac *Action) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(actionSchemaFile, ac)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (ac Action) IsEntity() {}

// Event An event indicating that this action was triggered
type Event struct {
	// A unique identifier for each action
	ID string `json:"id" firestore:"id"`

	// An event name - two upper case words separated by an underscore
	Name string `json:"name" firestore:"name"`

	// Technical metadata - when/where/why/who/what/how etc
	Context Context `json:"context,omitempty" firestore:"context,omitempty"`

	// The actual 'business data' carried by the event
	Payload Payload `json:"payload,omitempty" firestore:"payload,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ev *Event) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(eventSchemaFile, b, ev)
	if err != nil {
		return fmt.Errorf("invalid event JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ev *Event) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(eventSchemaFile, ev)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (ev Event) IsEntity() {}

// Context identifies when/where/why/who/what/how an event occured.
type Context struct {
	// the system or human user that created this event
	UserID string `json:"userID" firestore:"userID"`

	// the flavour of the feed/app that originated this event
	Flavour Flavour `json:"flavour" firestore:"flavour"`

	// the client (organization) that this user belongs to
	OrganizationID string `json:"organizationID" firestore:"organizationID"`

	// the location (e.g branch) from which the event was sent
	LocationID string `json:"locationID" firestore:"locationID"`

	// when this event was sent
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (ct *Context) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(contextSchemaFile, b, ct)
	if err != nil {
		return fmt.Errorf("invalid context JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (ct *Context) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(contextSchemaFile, ct)
}

// Payload carries the actual 'business data' carried by the event.
// It varies from event to event.
type Payload struct {
	Data map[string]interface{} `json:"data" firestore:"data"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (pl *Payload) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(payloadSchemaFile, b, pl)
	if err != nil {
		return fmt.Errorf("invalid payload JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (pl *Payload) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(payloadSchemaFile, pl)
}

// Nudge represents a "prompt" for a user e.g to set a PIN
type Nudge struct {
	// A unique identifier for each nudge
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// Visibility determines if a nudge should be visible or not
	Visibility Visibility `json:"visibility" firestore:"visibility"`

	// whether the nudge is done (acted on) or pending
	Status Status `json:"status" firestore:"status"`

	// When this nudge should be expired/removed, automatically. RFC3339.
	Expiry time.Time `json:"expiry" firestore:"expiry"`

	// the title (lead line) of the nudge
	Title string `json:"title" firestore:"title"`

	// the text/copy of the nudge
	Text string `json:"text" firestore:"text"`

	// an illustrative image for the nudge
	Links []Link `json:"links" firestore:"links"`

	// actions to include on the nudge
	Actions []Action `json:"actions" firestore:"actions"`

	// Identifiers of all the users that got this message
	Users []string `json:"users,omitempty" firestore:"users,omitempty"`

	// Identifiers of all the groups that got this message
	Groups []string `json:"groups,omitempty" firestore:"groups,omitempty"`

	// How the user should be notified of this new item, if at all
	NotificationChannels []Channel `json:"notificationChannels,omitempty" firestore:"notificationChannels,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (nu *Nudge) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(nudgeSchemaFile, b, nu)
	if err != nil {
		return fmt.Errorf("invalid nudge JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal verifies against JSON schema then marshals to JSON
func (nu *Nudge) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(nudgeSchemaFile, nu)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (nu Nudge) IsEntity() {}

// Item is a single item in a feed or in an inbox
type Item struct {
	// A unique identifier for each feed item
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// When this feed item should be expired/removed, automatically. RFC3339.
	Expiry time.Time `json:"expiry" firestore:"expiry"`

	// If a feed item is persistent, it also goes to the inbox
	// AND triggers a push notification.
	// Pinning a feed item makes it persistent.
	Persistent bool `json:"persistent" firestore:"persistent"`

	// Whether the task under a feed item is completed, pending etc
	Status Status `json:"status" firestore:"status"`

	// Whether the feed item is to be shown or hidden
	Visibility Visibility `json:"visibility" firestore:"visibility"`

	// A link to a PNG image that would serve as an avatar
	Icon Link `json:"icon" firestore:"icon"`

	// The person - real or robot - that generated this feed item. Rich text.
	Author string `json:"author" firestore:"author"`

	// An OPTIONAL second title line. Rich text.
	Tagline string `json:"tagline" firestore:"tagline"`

	// A label e.g for the queue that this item belongs to
	Label string `json:"label" firestore:"label"`

	// When this feed item was created. RFC3339.
	// This is used to calculate the feed item's age for display.
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`

	// An OPTIONAL summary line. Rich text.
	Summary string `json:"summary" firestore:"summary"`

	// Rich text that can include any unicode e.g emoji
	Text string `json:"text" firestore:"text"`

	// TextType determines how the frontend will render the text
	TextType TextType `json:"textType" firestore:"textType"`

	// an illustrative image for the item
	Links []Link `json:"links" firestore:"links"`

	// Actions are the primary, secondary and overflow actions associated
	// with a feed item
	Actions []Action `json:"actions,omitempty" firestore:"actions,omitempty"`

	// Conversations are messages and replies around a feed item
	Conversations []Message `json:"conversations,omitempty" firestore:"conversations,omitempty"`

	// Identifiers of all the users that got this message
	Users []string `json:"users,omitempty" firestore:"users,omitempty"`

	// Identifiers of all the groups that got this message
	Groups []string `json:"groups,omitempty" firestore:"groups,omitempty"`

	// How the user should be notified of this new item, if at all
	NotificationChannels []Channel `json:"notificationChannels,omitempty" firestore:"notificationChannels,omitempty"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (it *Item) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(itemSchemaFile, b, it)
	if err != nil {
		return fmt.Errorf("invalid item JSON: %w", err)
	}

	if it.Icon.LinkType != LinkTypePngImage {
		return fmt.Errorf("an icon must be a PNG image")
	}

	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (it *Item) ValidateAndMarshal() ([]byte, error) {
	if it.Icon.LinkType != LinkTypePngImage {
		return nil, fmt.Errorf("an icon must be a PNG image")
	}

	return validateAndMarshal(itemSchemaFile, it)
}

// IsEntity marks this as an Apollo federation GraphQL entity
func (it Item) IsEntity() {}

// Message is a message in a thread of conversations attached to a feed item
type Message struct {
	// A unique identifier for each message on the thread
	ID string `json:"id" firestore:"id"`

	// A higher sequence number means that it came later
	SequenceNumber int `json:"sequenceNumber" firestore:"sequenceNumber"`

	// Rich text that can include any unicode e.g emoji
	Text string `json:"text" firestore:"text"`

	// The unique ID of any message that this one is replying to - a thread
	ReplyTo string `json:"replyTo" firestore:"replyTo"`

	// The UID of the user that posted the message
	PostedByUID string `json:"postedByUID" firestore:"postedByUID"`

	// The UID of the user that posted the message
	PostedByName string `json:"postedByName" firestore:"postedByName"`

	// when this message was sent
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (msg *Message) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(messageSchemaFile, b, msg)
	if err != nil {
		return fmt.Errorf("invalid message JSON: %w", err)
	}
	return nil
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (msg *Message) ValidateAndMarshal() ([]byte, error) {
	return validateAndMarshal(messageSchemaFile, msg)
}

// Link holds references to media that is part of the feed.
// The URL should embed authentication details.
// The treatment will depend on the specified asset type.
type Link struct {
	// A unique identifier for each feed item
	ID string `json:"id" firestore:"id"`

	// A URL at which the video can be accessed.
	// For a private video, the URL should include authentication information.
	URL string `json:"url" firestore:"url"`

	// LinkType of link
	LinkType LinkType `json:"linkType" firestore:"linkType"`
}

func (l *Link) validateLinkType() error {
	if !govalidator.IsURL(l.URL) {
		return fmt.Errorf("%s is not a valid URL", l.URL)
	}
	switch l.LinkType {
	case LinkTypePdfDocument:
		if !strings.Contains(l.URL, ".png") {
			return fmt.Errorf("%s does not end with .pdf", l.URL)
		}
	case LinkTypePngImage:
		if !strings.Contains(l.URL, ".png") {
			return fmt.Errorf("%s does not end with .png", l.URL)
		}
	case LinkTypeYoutubeVideo:
		if !strings.Contains(l.URL, "youtube.com") {
			return fmt.Errorf("%s is not a youtube.com URL", l.URL)
		}
	}
	return nil
}

// ValidateAndUnmarshal checks that the input data is valid as per the
// relevant JSON schema and unmarshals it if it is
func (l *Link) ValidateAndUnmarshal(b []byte) error {
	err := validateAndUnmarshal(linkSchemaFile, b, l)
	if err != nil {
		return fmt.Errorf("invalid video JSON: %w", err)
	}

	return l.validateLinkType()
}

// ValidateAndMarshal validates against JSON schema then marshals to JSON
func (l *Link) ValidateAndMarshal() ([]byte, error) {
	err := l.validateLinkType()
	if err != nil {
		return nil, fmt.Errorf("can't marshal invalid link: %w", err)
	}
	return validateAndMarshal(linkSchemaFile, l)
}

// FilterParams organizes the parameters needed to filter persistent feed items
type FilterParams struct {
	Labels []string `json:"labels"`
}

func validateAgainstSchema(sch string, b []byte) error {
	schemaURL := fmt.Sprintf("%s/%s", getSchemaURL(), sch)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)
	documentLoader := gojsonschema.NewStringLoader(string(b))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf(
			"failed to validate `%s` against %s, got %#v: %w",
			string(b),
			sch,
			result,
			err,
		)
	}
	if !result.Valid() {
		errMsgs := []string{}
		for _, vErr := range result.Errors() {
			errType := vErr.Type()
			val := vErr.Value()
			context := vErr.Context().String()
			field := vErr.Field()
			desc := vErr.Description()
			descFormat := vErr.DescriptionFormat()
			details := vErr.Details()
			errMsg := fmt.Sprintf(
				"errType: %s\nval: %s\ncontext: %s\nfield: %s\ndesc: %s\ndescFormat: %s\ndetails: %s\n",
				errType,
				val,
				context,
				field,
				desc,
				descFormat,
				details,
			)
			errMsgs = append(errMsgs, errMsg)
		}
		return fmt.Errorf(
			"the result of validating `%s` against %s is not valid: %#v",
			string(b),
			sch,
			errMsgs,
		)
	}
	return nil
}

func validateAndUnmarshal(sch string, b []byte, el Element) error {
	err := validateAgainstSchema(sch, b)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	err = json.Unmarshal(b, el)
	if err != nil {
		return fmt.Errorf("can't unmarshal JSON to struct: %w", err)
	}
	return nil
}

func validateAndMarshal(sch string, el Element) ([]byte, error) {
	bs, err := json.Marshal(el)
	if err != nil {
		return nil, fmt.Errorf("can't marshal %T to JSON: %w", el, err)
	}
	err = validateAgainstSchema(sch, bs)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return bs, nil
}

// getSchemaURL serves JSON schema from this server and only falls back to a
// remote schema host when the local server cannot serve the JSON schema files.
// This has been done so as to reduce the impact of the network and DNS on the
// schema validation process - a critical path activity.
func getSchemaURL() string {
	schemaHost, err := base.GetEnvVar(schemaHostEnvVarName)
	if err != nil {
		log.Printf("can't get env var `%s`: %s", schemaHostEnvVarName, err)
	}

	client := http.Client{
		Timeout: time.Second * 1, // aggressive timeout
	}
	req, err := http.NewRequest(http.MethodGet, schemaHost, nil)
	if err != nil {
		log.Printf("can't create request to local schema URL: %s", err)
	}
	if err == nil {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error accessing schema URL: %s", err)
		}
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				log.Printf("schema URL error status code: %s", resp.Status)
			}
			if resp.StatusCode == http.StatusOK {
				return schemaHost // we want this case to be the most common
			}
		}
	}

	// fall back to an externally hosted schema
	return fallbackSchemaHost
}

// GetPNGImageLink returns an initialized PNG image link.
//
// It is used in testing and default data generation.
func GetPNGImageLink(url string) Link {
	return Link{
		ID:       ksuid.New().String(),
		URL:      url,
		LinkType: LinkTypePngImage,
	}
}

// GetYoutubeVideoLink returns an initialized YouTube video link.
//
// It is used in testing and default data generation.
func GetYoutubeVideoLink(url string) Link {
	return Link{
		ID:       ksuid.New().String(),
		URL:      url,
		LinkType: LinkTypeYoutubeVideo,
	}
}

// GetPDFDocumentLink returns an initialized PDF document link.
//
// It is used in testing and default data generation.
func GetPDFDocumentLink(url string) Link {
	return Link{
		ID:       ksuid.New().String(),
		URL:      url,
		LinkType: LinkTypePdfDocument,
	}
}
