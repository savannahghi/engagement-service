package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/feedlib"
	"gitlab.slade360emr.com/go/apiclient"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
)

// Repository defines methods for persistence and retrieval of feeds
type Repository interface {
	// getting a feed...create a default feed if it does not exist
	// return: feed, matching count, total count, optional error
	GetFeed(
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
	GetFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) (*feedlib.Item, error)

	// saving a new feed item
	SaveFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// updating an existing feed item
	UpdateFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) (*feedlib.Nudge, error)

	// saving a new modified nudge
	SaveNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// updating an existing nudge
	UpdateNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) (*feedlib.Action, error)

	// saving a new action
	SaveAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		action *feedlib.Action,
	) (*feedlib.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		message *feedlib.Message,
	) (*feedlib.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) (*feedlib.Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessages(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) ([]feedlib.Message, error)

	SaveIncomingEvent(
		ctx context.Context,
		event *feedlib.Event,
	) error

	SaveOutgoingEvent(
		ctx context.Context,
		event *feedlib.Event,
	) error

	GetNudges(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
	) ([]feedlib.Nudge, error)

	GetActions(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]feedlib.Action, error)

	GetItems(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]feedlib.Item, error)

	Labels(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]string, error)

	SaveLabel(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		label string,
	) error

	UnreadPersistentItems(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCount(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) error

	GetDefaultNudgeByTitle(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		title string,
	) (*feedlib.Nudge, error)

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
	) ([]*apiclient.Segment, error)

	UpdateMessageSentStatus(
		ctx context.Context,
		phonenumber string,
		segment string,
	) error

	UpdateUserCRMEmail(ctx context.Context, phoneNumber string, payload *dto.UpdateContactPSMessage) error
	UpdateUserCRMBewellAware(ctx context.Context, email string, payload *dto.UpdateContactPSMessage) error

	LoadMarketingData(ctx context.Context, data apiclient.Segment) (int, error)

	RollBackMarketingData(ctx context.Context, data apiclient.Segment) error
	SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error
	UpdateMailgunDeliveryStatus(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error)

	GetSladerDataByPhone(
		ctx context.Context,
		phonenumber string,
	) (*apiclient.Segment, error)
}
