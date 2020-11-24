package feed

import (
	"context"
)

// Repository defines methods for persistence and retrieval of feeds
type Repository interface {
	// getting a feed...create a default feed if it does not exist
	// return: feed, matching count, total count, optional error
	GetFeed(
		ctx context.Context,
		uid string,
		flavour Flavour,
		persistent BooleanFilter,
		status *Status,
		visibility *Visibility,
		expired *BooleanFilter,
		filterParams *FilterParams,
	) (*Feed, error)

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItem(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
	) (*Item, error)

	// saving a new feed item
	SaveFeedItem(
		ctx context.Context,
		uid string,
		flavour Flavour,
		item *Item,
	) (*Item, error)

	// updating an existing feed item
	UpdateFeedItem(
		ctx context.Context,
		uid string,
		flavour Flavour,
		item *Item,
	) (*Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItem(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudge(
		ctx context.Context,
		uid string,
		flavour Flavour,
		nudgeID string,
	) (*Nudge, error)

	// saving a new modified nudge
	SaveNudge(
		ctx context.Context,
		uid string,
		flavour Flavour,
		nudge *Nudge,
	) (*Nudge, error)

	// updating an existing nudge
	UpdateNudge(
		ctx context.Context,
		uid string,
		flavour Flavour,
		nudge *Nudge,
	) (*Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudge(
		ctx context.Context,
		uid string,
		flavour Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetAction(
		ctx context.Context,
		uid string,
		flavour Flavour,
		actionID string,
	) (*Action, error)

	// saving a new action
	SaveAction(
		ctx context.Context,
		uid string,
		flavour Flavour,
		action *Action,
	) (*Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteAction(
		ctx context.Context,
		uid string,
		flavour Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessage(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
		message *Message,
	) (*Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessage(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
		messageID string,
	) (*Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessages(
		ctx context.Context,
		uid string,
		flavour Flavour,
		itemID string,
	) ([]Message, error)

	SaveIncomingEvent(
		ctx context.Context,
		event *Event,
	) error

	SaveOutgoingEvent(
		ctx context.Context,
		event *Event,
	) error

	GetNudges(
		ctx context.Context,
		uid string,
		flavour Flavour,
		status *Status,
		visibility *Visibility,
	) ([]Nudge, error)

	GetActions(
		ctx context.Context,
		uid string,
		flavour Flavour,
	) ([]Action, error)

	GetItems(
		ctx context.Context,
		uid string,
		flavour Flavour,
		persistent BooleanFilter,
		status *Status,
		visibility *Visibility,
		expired *BooleanFilter,
		filterParams *FilterParams,
	) ([]Item, error)

	Labels(
		ctx context.Context,
		uid string,
		flavour Flavour,
	) ([]string, error)

	SaveLabel(
		ctx context.Context,
		uid string,
		flavour Flavour,
		label string,
	) error

	UnreadPersistentItems(
		ctx context.Context,
		uid string,
		flavour Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCount(
		ctx context.Context,
		uid string,
		flavour Flavour,
	) error
}
