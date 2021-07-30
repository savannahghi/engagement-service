package mock

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/feedlib"
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
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) (*domain.Feed, error)

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) (*feedlib.Item, error)

	// saving a new feed item
	SaveFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// updating an existing feed item
	UpdateFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) (*feedlib.Nudge, error)

	// saving a new modified nudge
	SaveNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// updating an existing nudge
	UpdateNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) (*feedlib.Action, error)

	// saving a new action
	SaveActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		action *feedlib.Action,
	) (*feedlib.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		message *feedlib.Message,
	) (*feedlib.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) (*feedlib.Message, error)

	// DeleteMessage deletes a message
	DeleteMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessagesFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) ([]feedlib.Message, error)

	SaveIncomingEventFn func(
		ctx context.Context,
		event *feedlib.Event,
	) error

	SaveOutgoingEventFn func(
		ctx context.Context,
		event *feedlib.Event,
	) error

	GetNudgesFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
	) ([]feedlib.Nudge, error)

	GetActionsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]feedlib.Action, error)

	GetItemsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]feedlib.Item, error)

	LabelsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]string, error)

	SaveLabelFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		label string,
	) error

	UnreadPersistentItemsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCountFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) error

	GetDefaultNudgeByTitleFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		title string,
	) (*feedlib.Nudge, error)

	SaveMarketingMessageFn func(
		ctx context.Context,
		data dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

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
		data *dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

	UpdateMessageSentStatusFn func(
		ctx context.Context,
		phonenumber string,
		segment string,
	) error

	UpdateUserCRMEmailFn       func(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error
	UpdateUserCRMBewellAwareFn func(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error

	SaveOutgoingEmailsFn          func(ctx context.Context, payload *dto.OutgoingEmailsLog) error
	UpdateMailgunDeliveryStatusFn func(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error)

	GetMarketingSMSByIDFn func(
		ctx context.Context,
		id string,
	) (*dto.MarketingSMS, error)

	GetMarketingSMSByPhoneFn func(
		ctx context.Context,
		phoneNumber string,
	) (*dto.MarketingSMS, error)
}

// GetFeed ...
func (f *FakeEngagementRepository) GetFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) (*domain.Feed, error) {
	return f.GetFeedFn(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
}

// GetFeedItem ...
func (f *FakeEngagementRepository) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return f.GetFeedItemFn(ctx, uid, flavour, itemID)
}

// SaveFeedItem ...
func (f *FakeEngagementRepository) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return f.SaveFeedItemFn(ctx, uid, flavour, item)
}

// UpdateFeedItem ...
func (f *FakeEngagementRepository) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return f.UpdateFeedItemFn(ctx, uid, flavour, item)
}

// DeleteFeedItem permanently deletes a feed item and it's copies
func (f *FakeEngagementRepository) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	return f.DeleteFeedItemFn(ctx, uid, flavour, itemID)
}

// GetNudge gets THE LATEST VERSION OF a nudge from a feed
func (f *FakeEngagementRepository) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return f.GetNudgeFn(ctx, uid, flavour, nudgeID)
}

// SaveNudge saves a new modified nudge
func (f *FakeEngagementRepository) SaveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return f.SaveNudgeFn(ctx, uid, flavour, nudge)
}

// UpdateNudge updates an existing nudge
func (f *FakeEngagementRepository) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return f.UpdateNudgeFn(ctx, uid, flavour, nudge)
}

// DeleteNudge permanently deletes a nudge and it's copies
func (f *FakeEngagementRepository) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	return f.DeleteNudgeFn(ctx, uid, flavour, nudgeID)
}

// GetAction gets THE LATEST VERSION OF a single action
func (f *FakeEngagementRepository) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	return f.GetActionFn(ctx, uid, flavour, actionID)
}

// SaveAction saves a new action
func (f *FakeEngagementRepository) SaveAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	return f.SaveActionFn(ctx, uid, flavour, action)
}

// DeleteAction permanently deletes an action and it's copies
func (f *FakeEngagementRepository) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	return f.DeleteActionFn(ctx, uid, flavour, actionID)
}

// PostMessage posts a message or a reply to a message/thread
func (f *FakeEngagementRepository) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	return f.PostMessageFn(ctx, uid, flavour, itemID, message)
}

// GetMessage retrieves THE LATEST VERSION OF a message
func (f *FakeEngagementRepository) GetMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) (*feedlib.Message, error) {
	return f.GetMessageFn(ctx, uid, flavour, itemID, messageID)
}

// DeleteMessage deletes a message
func (f *FakeEngagementRepository) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	return f.DeleteMessageFn(ctx, uid, flavour, itemID, messageID)
}

// GetMessages retrieves a message
func (f *FakeEngagementRepository) GetMessages(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) ([]feedlib.Message, error) {
	return f.GetMessagesFn(ctx, uid, flavour, itemID)
}

// SaveIncomingEvent ...
func (f *FakeEngagementRepository) SaveIncomingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return f.SaveIncomingEventFn(ctx, event)
}

// SaveOutgoingEvent ...
func (f *FakeEngagementRepository) SaveOutgoingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return f.SaveOutgoingEventFn(ctx, event)
}

// GetNudges ...
func (f *FakeEngagementRepository) GetNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) ([]feedlib.Nudge, error) {
	return f.GetNudgesFn(ctx, uid, flavour, status, visibility, expired)
}

// GetActions ...
func (f *FakeEngagementRepository) GetActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]feedlib.Action, error) {
	return f.GetActionsFn(ctx, uid, flavour)
}

// GetItems ...
func (f *FakeEngagementRepository) GetItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]feedlib.Item, error) {
	return f.GetItemsFn(ctx, uid, flavour, persistent, status, visibility, expired, filterParams)
}

// Labels ...
func (f *FakeEngagementRepository) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	return f.LabelsFn(ctx, uid, flavour)
}

// SaveLabel ...
func (f *FakeEngagementRepository) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	return f.SaveLabelFn(ctx, uid, flavour, label)
}

// UnreadPersistentItems ...
func (f *FakeEngagementRepository) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	return f.UnreadPersistentItemsFn(ctx, uid, flavour)
}

// UpdateUnreadPersistentItemsCount ...
func (f *FakeEngagementRepository) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return f.UpdateUnreadPersistentItemsCountFn(ctx, uid, flavour)
}

// GetDefaultNudgeByTitle ...
func (f *FakeEngagementRepository) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	return f.GetDefaultNudgeByTitleFn(ctx, uid, flavour, title)
}

// SaveMarketingMessage saves the callback data for future analysis
func (f *FakeEngagementRepository) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
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
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return f.UpdateMarketingMessageFn(ctx, data)
}

// UpdateMessageSentStatus ..
func (f *FakeEngagementRepository) UpdateMessageSentStatus(
	ctx context.Context,
	phonenumber string,
	segment string,
) error {
	return f.UpdateMessageSentStatusFn(ctx, phonenumber, segment)
}

// UpdateUserCRMEmail ..
func (f *FakeEngagementRepository) UpdateUserCRMEmail(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error {
	return f.UpdateUserCRMEmailFn(ctx, phoneNumber, payload)
}

// UpdateUserCRMBewellAware ..
func (f *FakeEngagementRepository) UpdateUserCRMBewellAware(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error {
	return f.UpdateUserCRMBewellAwareFn(ctx, email, payload)
}

// SaveOutgoingEmails ...
func (f *FakeEngagementRepository) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return f.SaveOutgoingEmailsFn(ctx, payload)
}

// UpdateMailgunDeliveryStatus ...
func (f *FakeEngagementRepository) UpdateMailgunDeliveryStatus(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error) {
	return f.UpdateMailgunDeliveryStatusFn(ctx, payload)
}

// GetMarketingSMSByPhone ..
func (f *FakeEngagementRepository) GetMarketingSMSByPhone(ctx context.Context, phoneNumber string) (*dto.MarketingSMS, error) {
	return f.GetMarketingSMSByPhoneFn(ctx, phoneNumber)
}

// GetMarketingSMSByID ..
func (f *FakeEngagementRepository) GetMarketingSMSByID(
	ctx context.Context,
	id string,
) (*dto.MarketingSMS, error) {
	return f.GetMarketingSMSByIDFn(ctx, id)
}
