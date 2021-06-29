package mock

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
)

// FakeEngagementRepository is a mock engagement repository
type FakeEngagementRepository struct {
	GetFeedFn func(
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

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	// saving a new feed item
	SaveFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		item *base.Item,
	) (*base.Item, error)

	// updating an existing feed item
	UpdateFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		item *base.Item,
	) (*base.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudgeFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	// saving a new modified nudge
	SaveNudgeFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudge *base.Nudge,
	) (*base.Nudge, error)

	// updating an existing nudge
	UpdateNudgeFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudge *base.Nudge,
	) (*base.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudgeFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetActionFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) (*base.Action, error)

	// saving a new action
	SaveActionFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		action *base.Action,
	) (*base.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteActionFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessageFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		message *base.Message,
	) (*base.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessageFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		messageID string,
	) (*base.Message, error)

	// DeleteMessage deletes a message
	DeleteMessageFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessagesFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) ([]base.Message, error)

	SaveIncomingEventFn func(
		ctx context.Context,
		event *base.Event,
	) error

	SaveOutgoingEventFn func(
		ctx context.Context,
		event *base.Event,
	) error

	GetNudgesFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		status *base.Status,
		visibility *base.Visibility,
		expired *base.BooleanFilter,
	) ([]base.Nudge, error)

	GetActionsFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) ([]base.Action, error)

	GetItemsFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		persistent base.BooleanFilter,
		status *base.Status,
		visibility *base.Visibility,
		expired *base.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]base.Item, error)

	LabelsFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) ([]string, error)

	SaveLabelFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		label string,
	) error

	UnreadPersistentItemsFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCountFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) error

	GetDefaultNudgeByTitleFn func(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		title string,
	) (*base.Nudge, error)

	SaveMarketingMessageFn func(
		ctx context.Context,
		data dto.MarketingSMS,
	) error

	SaveTwilioResponseFn func(
		ctx context.Context,
		data dto.Message,
	) error

	SaveNotificationFn func(
		ctx context.Context,
		firestoreClient *firestore.Client,
		notification dto.SavedNotification,
	) error

	RetrieveNotificationFn func(
		ctx context.Context,
		firestoreClient *firestore.Client,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SaveNPSResponseFn func(
		ctx context.Context,
		response *dto.NPSResponse,
	) error

	UpdateMarketingMessageFn func(
		ctx context.Context,
		phoneNumber string,
		deliveryReport *dto.ATDeliveryReport,
	) (*dto.MarketingSMS, error)

	RetrieveMarketingDataFn func(
		ctx context.Context,
		data *dto.MarketingMessagePayload,
	) ([]*dto.Segment, error)

	UpdateMessageSentStatusFn func(
		ctx context.Context,
		phonenumber string,
	) error

	UpdateUserCRMEmailFn       func(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error
	UpdateUserCRMBewellAwareFn func(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error

	IsOptedOutedFn func(ctx context.Context, phoneNumber string) (bool, error)
}

// GetFeed ...
func (f *FakeEngagementRepository) GetFeed(
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
	return f.GetFeedFn(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
}

// GetFeedItem ...
func (f *FakeEngagementRepository) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	return f.GetFeedItemFn(ctx, uid, flavour, itemID)
}

// SaveFeedItem ...
func (f *FakeEngagementRepository) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	item *base.Item,
) (*base.Item, error) {
	return f.SaveFeedItemFn(ctx, uid, flavour, item)
}

// UpdateFeedItem ...
func (f *FakeEngagementRepository) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	item *base.Item,
) (*base.Item, error) {
	return f.UpdateFeedItemFn(ctx, uid, flavour, item)
}

// DeleteFeedItem permanently deletes a feed item and it's copies
func (f *FakeEngagementRepository) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) error {
	return f.DeleteFeedItemFn(ctx, uid, flavour, itemID)
}

// GetNudge gets THE LATEST VERSION OF a nudge from a feed
func (f *FakeEngagementRepository) GetNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {
	return f.GetNudgeFn(ctx, uid, flavour, nudgeID)
}

// SaveNudge saves a new modified nudge
func (f *FakeEngagementRepository) SaveNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudge *base.Nudge,
) (*base.Nudge, error) {
	return f.SaveNudgeFn(ctx, uid, flavour, nudge)
}

// UpdateNudge updates an existing nudge
func (f *FakeEngagementRepository) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudge *base.Nudge,
) (*base.Nudge, error) {
	return f.UpdateNudgeFn(ctx, uid, flavour, nudge)
}

// DeleteNudge permanently deletes a nudge and it's copies
func (f *FakeEngagementRepository) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) error {
	return f.DeleteNudgeFn(ctx, uid, flavour, nudgeID)
}

// GetAction gets THE LATEST VERSION OF a single action
func (f *FakeEngagementRepository) GetAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) (*base.Action, error) {
	return f.GetActionFn(ctx, uid, flavour, actionID)
}

// SaveAction saves a new action
func (f *FakeEngagementRepository) SaveAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	action *base.Action,
) (*base.Action, error) {
	return f.SaveActionFn(ctx, uid, flavour, action)
}

// DeleteAction permanently deletes an action and it's copies
func (f *FakeEngagementRepository) DeleteAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) error {
	return f.DeleteActionFn(ctx, uid, flavour, actionID)
}

// PostMessage posts a message or a reply to a message/thread
func (f *FakeEngagementRepository) PostMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	message *base.Message,
) (*base.Message, error) {
	return f.PostMessageFn(ctx, uid, flavour, itemID, message)
}

// GetMessage retrieves THE LATEST VERSION OF a message
func (f *FakeEngagementRepository) GetMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	messageID string,
) (*base.Message, error) {
	return f.GetMessageFn(ctx, uid, flavour, itemID, messageID)
}

// DeleteMessage deletes a message
func (f *FakeEngagementRepository) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	messageID string,
) error {
	return f.DeleteMessageFn(ctx, uid, flavour, itemID, messageID)
}

// GetMessages retrieves a message
func (f *FakeEngagementRepository) GetMessages(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) ([]base.Message, error) {
	return f.GetMessagesFn(ctx, uid, flavour, itemID)
}

// SaveIncomingEvent ...
func (f *FakeEngagementRepository) SaveIncomingEvent(
	ctx context.Context,
	event *base.Event,
) error {
	return f.SaveIncomingEventFn(ctx, event)
}

// SaveOutgoingEvent ...
func (f *FakeEngagementRepository) SaveOutgoingEvent(
	ctx context.Context,
	event *base.Event,
) error {
	return f.SaveOutgoingEventFn(ctx, event)
}

// GetNudges ...
func (f *FakeEngagementRepository) GetNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
) ([]base.Nudge, error) {
	return f.GetNudgesFn(ctx, uid, flavour, status, visibility, expired)
}

// GetActions ...
func (f *FakeEngagementRepository) GetActions(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) ([]base.Action, error) {
	return f.GetActionsFn(ctx, uid, flavour)
}

// GetItems ...
func (f *FakeEngagementRepository) GetItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]base.Item, error) {
	return f.GetItemsFn(ctx, uid, flavour, persistent, status, visibility, expired, filterParams)
}

// Labels ...
func (f *FakeEngagementRepository) Labels(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) ([]string, error) {
	return f.LabelsFn(ctx, uid, flavour)
}

// SaveLabel ...
func (f *FakeEngagementRepository) SaveLabel(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	label string,
) error {
	return f.SaveLabelFn(ctx, uid, flavour, label)
}

// UnreadPersistentItems ...
func (f *FakeEngagementRepository) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) (int, error) {
	return f.UnreadPersistentItemsFn(ctx, uid, flavour)
}

// UpdateUnreadPersistentItemsCount ...
func (f *FakeEngagementRepository) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) error {
	return f.UpdateUnreadPersistentItemsCountFn(ctx, uid, flavour)
}

// GetDefaultNudgeByTitle ...
func (f *FakeEngagementRepository) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	title string,
) (*base.Nudge, error) {
	return f.GetDefaultNudgeByTitleFn(ctx, uid, flavour, title)
}

// SaveMarketingMessage saves the callback data for future analysis
func (f *FakeEngagementRepository) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) error {
	return f.SaveMarketingMessageFn(ctx, data)
}

// SaveTwilioResponse saves the callback data for future analysis
func (f *FakeEngagementRepository) SaveTwilioResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return f.SaveTwilioResponseFn(ctx, data)
}

// SaveNotification saves a notification
func (f *FakeEngagementRepository) SaveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	notification dto.SavedNotification,
) error {
	return f.SaveNotificationFn(ctx, firestoreClient, notification)
}

// RetrieveNotification retrieves a notification
func (f *FakeEngagementRepository) RetrieveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.RetrieveNotificationFn(ctx, firestoreClient, registrationToken, newerThan, limit)
}

// SaveNPSResponse saves a NPS response
func (f *FakeEngagementRepository) SaveNPSResponse(
	ctx context.Context,
	response *dto.NPSResponse,
) error {
	return f.SaveNPSResponseFn(ctx, response)
}

// UpdateMarketingMessage ..
func (f *FakeEngagementRepository) UpdateMarketingMessage(
	ctx context.Context,
	phoneNumber string,
	deliveryReport *dto.ATDeliveryReport,
) (*dto.MarketingSMS, error) {
	return f.UpdateMarketingMessageFn(ctx, phoneNumber, deliveryReport)
}

// RetrieveMarketingData ..
func (f *FakeEngagementRepository) RetrieveMarketingData(
	ctx context.Context,
	data *dto.MarketingMessagePayload,
) ([]*dto.Segment, error) {
	return f.RetrieveMarketingDataFn(ctx, data)
}

// UpdateMessageSentStatus ..
func (f *FakeEngagementRepository) UpdateMessageSentStatus(
	ctx context.Context,
	phonenumber string,
) error {
	return f.UpdateMessageSentStatusFn(ctx, phonenumber)
}

// UpdateUserCRMEmail ..
func (f *FakeEngagementRepository) UpdateUserCRMEmail(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error {
	return f.UpdateUserCRMEmailFn(ctx, phoneNumber, payload)
}

// UpdateUserCRMBewellAware ..
func (f *FakeEngagementRepository) UpdateUserCRMBewellAware(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error {
	return f.UpdateUserCRMBewellAwareFn(ctx, email, payload)
}

// IsOptedOuted ..
func (f *FakeEngagementRepository) IsOptedOuted(ctx context.Context, phoneNumber string) (bool, error) {
	return f.IsOptedOutedFn(ctx, phoneNumber)
}
