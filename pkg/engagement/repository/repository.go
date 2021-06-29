package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"

	"gitlab.slade360emr.com/go/base"
)

// Repository defines methods for persistence and retrieval of feeds
type Repository interface {
	// getting a feed...create a default feed if it does not exist
	// return: feed, matching count, total count, optional error
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

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) (*base.Item, error)

	// saving a new feed item
	SaveFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		item *base.Item,
	) (*base.Item, error)

	// updating an existing feed item
	UpdateFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		item *base.Item,
	) (*base.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItem(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) (*base.Nudge, error)

	// saving a new modified nudge
	SaveNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudge *base.Nudge,
	) (*base.Nudge, error)

	// updating an existing nudge
	UpdateNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudge *base.Nudge,
	) (*base.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudge(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) (*base.Action, error)

	// saving a new action
	SaveAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		action *base.Action,
	) (*base.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteAction(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessage(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		message *base.Message,
	) (*base.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessage(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		messageID string,
	) (*base.Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessages(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		itemID string,
	) ([]base.Message, error)

	SaveIncomingEvent(
		ctx context.Context,
		event *base.Event,
	) error

	SaveOutgoingEvent(
		ctx context.Context,
		event *base.Event,
	) error

	GetNudges(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		status *base.Status,
		visibility *base.Visibility,
		expired *base.BooleanFilter,
	) ([]base.Nudge, error)

	GetActions(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
	) ([]base.Action, error)

	GetItems(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		persistent base.BooleanFilter,
		status *base.Status,
		visibility *base.Visibility,
		expired *base.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]base.Item, error)

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

	GetDefaultNudgeByTitle(
		ctx context.Context,
		uid string,
		flavour base.Flavour,
		title string,
	) (*base.Nudge, error)

	SaveMarketingMessage(
		ctx context.Context,
		data dto.MarketingSMS,
	) error

	UpdateMarketingMessage(
		ctx context.Context,
		phoneNumber string,
		deliveryReport *dto.ATDeliveryReport,
	) (*dto.MarketingSMS, error)

	SaveTwilioResponse(
		ctx context.Context,
		data dto.Message,
	) error

	SaveNotification(
		ctx context.Context,
		firestoreClient *firestore.Client,
		notification dto.SavedNotification,
	) error

	RetrieveNotification(
		ctx context.Context,
		firestoreClient *firestore.Client,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SaveNPSResponse(
		ctx context.Context,
		response *dto.NPSResponse,
	) error

	RetrieveMarketingData(
		ctx context.Context,
		data *dto.MarketingMessagePayload,
	) ([]*dto.Segment, error)

	UpdateMessageSentStatus(
		ctx context.Context,
		phonenumber string,
	) error

	UpdateUserCRMEmail(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error
	UpdateUserCRMBewellAware(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error
}
