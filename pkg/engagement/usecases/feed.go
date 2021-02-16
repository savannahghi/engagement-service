package usecases

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/messaging"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/exceptions"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

// FeedUseCases represents all the profile business logic
type FeedUseCases interface {
	GetFeed(
		ctx context.Context,
		uid *string,
		isAnonymous *bool,
		flavour base.Flavour,
		persistent base.BooleanFilter,
		status *base.Status,
		visibility *base.Visibility,
		expired *base.BooleanFilter,
		filterParams *helpers.FilterParams,
	) (*domain.Feed, error)

	GetThinFeed(
		ctx context.Context,
		uid *string,
		isAnonymous *bool,
		flavour base.Flavour,
	) (*domain.Feed, error)

	GetFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	GetNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	GetAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) (*base.Action, error)

	PublishFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		item *base.Item,
	) (*base.Item, error)

	DeleteFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) error

	ResolveFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	PinFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	UnpinFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	UnresolveFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	HideFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	ShowFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	Labels(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) ([]string, error)

	SaveLabel(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		label string,
	) error

	UnreadPersistentItems(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCount(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) error

	PublishNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudge *base.Nudge,
	) (*base.Nudge, error)

	ResolveNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	UnresolveNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	HideNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	ShowNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	GetDefaultNudgeByTitle(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		title string,
	) (*base.Nudge, error)

	ProcessEvent(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		event *base.Event,
	) error

	DeleteMessage(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		messageID string,

	) error

	PostMessage(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		message *base.Message,
	) (*base.Message, error)

	DeleteAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) error

	PublishAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		action *base.Action,
	) (*base.Action, error)

	DeleteNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) error
}

// FeedUseCaseImpl represents the feed usecase implementation
type FeedUseCaseImpl struct {
	Repository          repository.Repository
	NotificationService messaging.NotificationService
}

// NewFeed initializes a user feed
func NewFeed(
	repository repository.Repository,
	notificationService messaging.NotificationService,
) *FeedUseCaseImpl {
	return &FeedUseCaseImpl{
		Repository:          repository,
		NotificationService: notificationService,
	}
}

// GetFeed retrieves a feed
func (fe FeedUseCaseImpl) GetFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *helpers.FilterParams,
) (*domain.Feed, error) {
	feed, err := fe.Repository.GetFeed(
		ctx,
		uid,
		isAnonymous,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("feed retrieval error: %w", err)
	}

	// set the ID (computed, not stored)
	feed.ID = feed.GetID()
	feed.SequenceNumber = int(time.Now().Unix())

	return feed, nil
}

// GetThinFeed gets a feed with only the UID, flavour and dependencies
// filled in.
//
// It is used for efficient instantiation of feeds by code that does not need
// the full detail.
func (fe FeedUseCaseImpl) GetThinFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour base.Flavour,
) (*domain.Feed, error) {
	feed := &domain.Feed{
		UID:         *uid,
		Flavour:     flavour,
		Actions:     []base.Action{},
		Items:       []base.Item{},
		Nudges:      []base.Nudge{},
		IsAnonymous: isAnonymous,
	}

	// set the ID (computed, not stored)
	feed.ID = feed.GetID()
	feed.SequenceNumber = int(time.Now().Unix())

	return feed, nil
}

// GetFeedItem retrieves a feed item
func (fe FeedUseCaseImpl) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
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
func (fe FeedUseCaseImpl) GetNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {
	nudge, err := fe.Repository.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve nudge %s: %w", nudgeID, err)
	}

	if nudge == nil {
		return nil, nil
	}

	return nudge, nil
}

// GetAction retrieves a feed item
func (fe FeedUseCaseImpl) GetAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) (*base.Action, error) {
	action, err := fe.Repository.GetAction(ctx, uid, flavour, actionID)
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
func (fe FeedUseCaseImpl) PublishFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	item *base.Item,
) (*base.Item, error) {
	if item == nil {
		return nil, fmt.Errorf("cano't publish nil feed item")
	}

	if item.SequenceNumber == 0 {
		item.SequenceNumber = int(time.Now().Unix())
	}

	err := helpers.ValidateElement(item)
	if err != nil {
		return nil, fmt.Errorf("invalid item: %w", err)
	}

	for _, action := range item.Actions {
		if action.ActionType == base.ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	item, err = fe.Repository.SaveFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to publish feed item %s: %w", item.ID, err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemPublishTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) error {
	item, err := fe.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil || item == nil {
		// fails to error because it should be safe to retry deletes
		return nil // does not exist, nothing to delete
	}

	err = fe.Repository.DeleteFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return fmt.Errorf("unable to delete item: %s", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemDeleteTopic),
		uid,
		flavour,
		item,
		map[string]interface{}{
			"itemID": item.ID,
		},
	); err != nil {
		return fmt.Errorf("unable to notify item to channel: %w", err)
	}

	return fe.Repository.DeleteFeedItem(ctx, uid, flavour, itemID)
}

// ResolveFeedItem marks a feed item as Done
func (fe FeedUseCaseImpl) ResolveFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {

	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Status = base.StatusDone
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemResolveTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) PinFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Persistent = true
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemResolveTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) UnpinFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Persistent = false
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to pin feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemPinTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) UnresolveFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Status = base.StatusPending
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemUnresolveTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) HideFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Visibility = base.VisibilityHide
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to hide feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemHideTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) ShowFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {

	item, err := fe.Repository.GetFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get feed item with ID %s", itemID)
	}

	if item == nil {
		return nil, exceptions.ErrNilFeedItem
	}

	item.Visibility = base.VisibilityShow
	item.SequenceNumber = item.SequenceNumber + 1

	item, err = fe.Repository.UpdateFeedItem(ctx, uid, flavour, item)
	if err != nil {
		return nil, fmt.Errorf("unable to show feed item: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ItemShowTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) Labels(
	ctx context.Context,
	uid string, flavour base.Flavour,
) ([]string, error) {
	return fe.Repository.Labels(ctx, uid, flavour)
}

// SaveLabel saves the indicated label, if it does not already exist
func (fe FeedUseCaseImpl) SaveLabel(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	label string,
) error {

	return fe.Repository.SaveLabel(ctx, uid, flavour, label)
}

// UnreadPersistentItems returns the number of unread inbox items for this feed
func (fe FeedUseCaseImpl) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) (int, error) {
	return fe.Repository.UnreadPersistentItems(ctx, uid, flavour)
}

// UpdateUnreadPersistentItemsCount updates the number of unread inbox items
func (fe FeedUseCaseImpl) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) error {
	return fe.Repository.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
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
func (fe FeedUseCaseImpl) PublishNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudge *base.Nudge,
) (*base.Nudge, error) {

	if nudge == nil {
		return nil, fmt.Errorf("can't publish nil nudge")
	}

	if nudge.SequenceNumber == 0 {
		nudge.SequenceNumber = int(time.Now().Unix())
	}

	err := helpers.ValidateElement(nudge)
	if err != nil {
		return nil, fmt.Errorf("invalid nudge: %w", err)
	}

	for _, action := range nudge.Actions {
		if action.ActionType == base.ActionTypeFloating {
			return nil, fmt.Errorf("floating actions are only allowed at the global level")
		}
	}

	nudge, err = fe.Repository.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to publish nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgePublishTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) ResolveNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {

	nudge, err := fe.Repository.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, exceptions.ErrNilNudge
	}

	nudge.Status = base.StatusDone
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.Repository.UpdateNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgeResolveTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) UnresolveNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {

	nudge, err := fe.Repository.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, exceptions.ErrNilNudge
	}

	nudge.Status = base.StatusPending
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.Repository.UpdateNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgeUnresolveTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) HideNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {

	nudge, err := fe.Repository.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, exceptions.ErrNilNudge
	}

	nudge.Visibility = base.VisibilityHide
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.Repository.UpdateNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgeHideTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) ShowNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {

	nudge, err := fe.Repository.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge with ID %s", nudgeID)
	}

	if nudge == nil {
		return nil, exceptions.ErrNilNudge
	}

	nudge.Visibility = base.VisibilityShow
	nudge.SequenceNumber = nudge.SequenceNumber + 1

	nudge, err = fe.Repository.UpdateNudge(ctx, uid, flavour, nudge)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgeShowTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) error {

	nudge, err := fe.GetNudge(ctx, uid, flavour, nudgeID)
	if err != nil || nudge == nil {
		return nil // no error, "re-deleting" a nudge should not cause an error
	}

	err = fe.Repository.DeleteNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return fmt.Errorf("can't delete nudge: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.NudgeDeleteTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) PublishAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	action *base.Action,
) (*base.Action, error) {

	if action == nil {
		return nil, fmt.Errorf("can't publish nil nudge")
	}

	if action.SequenceNumber == 0 {
		action.SequenceNumber = int(time.Now().Unix())
	}

	err := helpers.ValidateElement(action)
	if err != nil {
		return nil, fmt.Errorf("invalid action: %w", err)
	}

	action, err = fe.Repository.SaveAction(ctx, uid, flavour, action)
	if err != nil {
		return nil, fmt.Errorf("unable to publish action: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ActionPublishTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) DeleteAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) error {

	action, err := fe.GetAction(ctx, uid, flavour, actionID)
	if err != nil || action == nil {
		return nil // no harm "re-deleting" an already deleted action
	}

	err = fe.Repository.DeleteAction(ctx, uid, flavour, actionID)
	if err != nil {
		return fmt.Errorf("unable to delete action: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.ActionDeleteTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) PostMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	message *base.Message,
) (*base.Message, error) {

	if message == nil {
		return nil, fmt.Errorf("can't post nil message")
	}

	if message.ID == "" {
		message.ID = ksuid.New().String()
	}

	if message.SequenceNumber == 0 {
		message.SequenceNumber = int(time.Now().Unix())
	}

	err := helpers.ValidateElement(message)
	if err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	msg, err := fe.Repository.PostMessage(
		ctx,
		uid,
		flavour,
		itemID,
		message,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to post a message: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.MessagePostTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	messageID string,
) error {

	message, err := fe.Repository.GetMessage(
		ctx,
		uid,
		flavour,
		itemID,
		messageID,
	)
	if err != nil || message == nil {
		return nil // no harm "re-deleting" an already deleted message
	}

	err = fe.Repository.DeleteMessage(
		ctx,
		uid,
		flavour,
		itemID,
		messageID,
	)
	if err != nil {
		return fmt.Errorf("unable to delete message: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.MessageDeleteTopic),
		uid,
		flavour,
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
func (fe FeedUseCaseImpl) ProcessEvent(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	event *base.Event,
) error {

	if event == nil {
		return fmt.Errorf("can't process nil event")
	}

	if !event.Context.Flavour.IsValid() {
		event.Context.Flavour = flavour
	}

	if event.ID == "" {
		event.ID = ksuid.New().String()
	}

	if event.Context.UserID == "" {
		event.Context.UserID = uid
	}

	err := helpers.ValidateElement(event)
	if err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	if event.Context.Flavour != flavour {
		return fmt.Errorf(
			"the event context flavour (%s) does not match the feed flavour (%s)",
			event.Context.Flavour,
			flavour,
		)
	}

	err = fe.Repository.SaveIncomingEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("can't save incoming event: %w", err)
	}

	if err := fe.NotificationService.Notify(
		ctx,
		helpers.AddPubSubNamespace(common.IncomingEventTopic),
		uid,
		flavour,
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

// GetDefaultNudgeByTitle retrieves a default feed nudge
func (fe FeedUseCaseImpl) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	title string,
) (*base.Nudge, error) {

	nudge, err := fe.Repository.GetDefaultNudgeByTitle(ctx, uid, flavour, title)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve verify email nudge: %w", err)
	}

	if nudge == nil {
		return nil, fmt.Errorf("can't get the default verify email nudge")
	}

	return nudge, nil
}