package database

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagement/pkg/engagement/application/common"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/exceptions"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"

	"github.com/savannahghi/engagement/pkg/engagement/repository"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement/pkg/engagement/services/database")

const (
	feedCollectionName           = "feed"
	elementsGroupName            = "elements"
	actionsSubcollectionName     = "actions"
	nudgesSubcollectionName      = "nudges"
	itemsSubcollectionName       = "items"
	messagesSubcollectionName    = "messages"
	incomingEventsCollectionName = "incoming_events"
	outgoingEventsCollectionName = "outgoing_events"
	//AITMarketingMessageName is the name of a Cloud Firestore collection into which AIT
	// callback data will be saved for future analysis
	AITMarketingMessageName = "ait_marketing_sms"

	// NPSResponseCollectionName firestore collection name where nps responses are stored
	NPSResponseCollectionName = "nps_response"

	twilioCallbackCollectionName = "twilio_callbacks"

	twilioVideoCallbackCollectionName = "twilio_video_callbacks"

	notificationCollectionName = "notifications"

	labelsDocID            = "item_labels"
	unreadInboxCountsDocID = "unread_inbox_counts"

	itemsLimit = 1000

	// outgoingEmails represent all the sent emails
	outgoingEmails = "outgoing_emails"
)

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(
	ctx context.Context,
) (repository.Repository, error) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Fatalf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}

	ff := &Repository{
		firestoreClient: fsc,
		mu:              &sync.Mutex{},
	}
	err = ff.checkPreconditions()
	if err != nil {
		log.Fatalf("firebase repository precondition check failed: %s", err)
	}

	return ff, nil
}

// Repository accesses and updates a feed that is stored on Firebase
type Repository struct {
	firestoreClient *firestore.Client
	mu              *sync.Mutex
}

func (fr Repository) checkPreconditions() error {
	if fr.firestoreClient == nil {
		return fmt.Errorf("nil firestore client in feed firebase repository")
	}

	if fr.mu == nil {
		return fmt.Errorf("nil feed repository mutex")
	}

	return nil
}

// GetFeed retrieves a feed by the user's ID and product flavour
//
// Feed items are ordered by:
//
//   1. Timestamp
//   2. Sequence number
//   3. ID (a tie breaker, in the unlikely event that the first two tie)
//
// Having established this ordering, we will implement very lightweight
// pagination using a start and end offset.
//
// This function is intended to be used for initial fetches of the inbox and
// manual refreshes. The clients are expected to implement logic to handle
// streaming inbox updates.
//
// The return parameters are:
//
//  1. A valid feed
//  2. The number of feed items that MATCHED this query (ignoring pagination)
//  3. The number of feed items that are NOT HIDDEN and are PENDING ACTION
//  4. An error, if any
func (fr Repository) GetFeed(
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
	ctx, span := tracer.Start(ctx, "GetFeed")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actions, err := fr.GetActions(ctx, *uid, flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}

	nudges, err := fr.GetNudges(
		ctx,
		*uid,
		flavour,
		status,
		visibility,
		expired,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}

	items, err := fr.GetItems(
		ctx,
		*uid,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get items: %w", err)
	}

	// only add default content if...
	// - the `persistent` filter is set to "BOTH"
	// - all other filters are nil
	noFilters := persistent == feedlib.BooleanFilterBoth && status == nil &&
		visibility == nil &&
		expired == nil &&
		filterParams == nil
	noActions := len(actions) == 0
	noNudges := len(nudges) == 0
	noItems := len(items) == 0
	if noFilters && noActions && noNudges && noItems {
		err = fr.initializeDefaultFeed(ctx, *uid, flavour)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to initialize default feed: %w",
				err,
			)
		}

		// this recursion is potentially dangerous but there's an integration test
		// that exercises this and reduces the risk of infinite recursion
		// we need to do this in order to have confidence that the initialization succeeded
		return fr.GetFeed(
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
	}

	feed := &domain.Feed{
		UID:         *uid,
		Flavour:     flavour,
		Actions:     actions,
		Nudges:      nudges,
		Items:       feedItemsFromCMSFeedTag(ctx, flavour),
		IsAnonymous: isAnonymous,
	}

	return feed, nil
}

func (fr Repository) initializeDefaultFeed(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	ctx, span := tracer.Start(ctx, "initializeDefaultFeed")
	defer span.End()
	fr.mu.Lock() // create default data once

	_, err := SetDefaultActions(ctx, uid, flavour, fr)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to set default actions: %w", err)
	}

	_, err = SetDefaultNudges(ctx, uid, flavour, fr)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to set default nudges: %w", err)
	}

	_, err = SetDefaultItems(ctx, uid, flavour, fr)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to set default items: %w", err)
	}

	fr.mu.Unlock() // release lock on default data
	return nil
}

// GetFeedItem retrieves and returns a single feed item
func (fr Repository) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "GetFeedItem")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	itemsCollection := fr.getItemsCollection(uid, flavour)
	el, err := fr.getSingleElement(
		ctx,
		itemsCollection,
		itemID,
		&feedlib.Item{},
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	item, ok := el.(*feedlib.Item)
	if !ok {
		return nil, fmt.Errorf("expected an Item, got %T", el)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, itemID)
	if err != nil || messages == nil {
		helpers.RecordSpanError(span, err)
		// the thread may not have been initiated yet
		item.Conversations = []feedlib.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// SaveFeedItem validates and saves a feed item.
// It's expected to have an ID and sequence number already.
// One suggestion is to use UUIDs for IDs and UTC Unix Epoch seconds for
// sequence numbers.
func (fr Repository) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "SaveFeedItem")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	_, err := item.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("item failed validation: %w", err)
	}

	coll := fr.getItemsCollection(uid, flavour)
	if err := fr.saveElement(
		ctx,
		item,
		item.ID,
		item.SequenceNumber,
		coll,
		true,
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save item: %w", err)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		helpers.RecordSpanError(span, err)
		// the thread may not have been initiated yet
		item.Conversations = []feedlib.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// UpdateFeedItem updates an existing feed item
func (fr Repository) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "UpdateFeedItem")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	_, err := item.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("invalid item: %w", err)
	}

	coll := fr.getItemsCollection(uid, flavour)
	if err := fr.saveElement(
		ctx,
		item,
		item.ID,
		item.SequenceNumber,
		coll,
		false, // not a new item, skip existing checks
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save item: %w", err)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		helpers.RecordSpanError(span, err)
		// the thread may not have been initiated yet
		item.Conversations = []feedlib.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// DeleteFeedItem deletes a nudge from a user's feed
func (fr Repository) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	ctx, span := tracer.Start(ctx, "DeleteFeedItem")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getItemsCollection(uid, flavour).Doc(itemID).Delete(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't delete item: %w", err)
	}

	return nil
}

// GetNudge retrieves a single nudge
func (fr Repository) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "GetNudge")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	nudgeCollection := fr.getNudgesCollection(uid, flavour)
	el, err := fr.getSingleElement(
		ctx,
		nudgeCollection,
		nudgeID,
		&feedlib.Nudge{},
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	nudge, ok := el.(*feedlib.Nudge)
	if !ok {
		return nil, fmt.Errorf("expected an nudge, got %T", el)
	}

	return nudge, nil
}

// SaveNudge saves an updated nudge
// It's expected to have an ID and sequence number already.
// One suggestion is to use UUIDs for IDs and UTC Unix Epoch seconds for
// sequence numbers.
func (fr Repository) SaveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "SaveNudge")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	// find an existing nudge with the same title
	existingNudge, err := fr.GetDefaultNudgeByTitle(
		ctx,
		uid,
		flavour,
		nudge.Title,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		// error means an existing nudge wasn't found
		log.Printf("nudge doesn't exist error: %v", err.Error())
	}

	if existingNudge != nil {
		return nil, fmt.Errorf(
			"cannot save nudge, found an existing nudge with same title",
		)
	}

	coll := fr.getNudgesCollection(uid, flavour)
	if err := fr.saveElement(
		ctx,
		nudge,
		nudge.ID,
		nudge.SequenceNumber,
		coll,
		true, // a new nudge
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}

	return nudge, nil
}

// UpdateNudge updates an existing nudge e.g to show or hide it
func (fr Repository) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "UpdateNudge")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if nudge == nil {
		return nil, fmt.Errorf("nil nudge")
	}

	_, err := nudge.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("nudge failed validation: %w", err)
	}

	coll := fr.getNudgesCollection(uid, flavour)
	if err := fr.saveElement(
		ctx,
		nudge,
		nudge.ID,
		nudge.SequenceNumber,
		coll,
		false, // not a new nudge, should not check for existence
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}

	return nudge, nil
}

// DeleteNudge deletes a nudge from a user's feed
func (fr Repository) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	ctx, span := tracer.Start(ctx, "DeleteNudge")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getNudgesCollection(uid, flavour).Doc(nudgeID).Delete(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't delete nudge: %w", err)
	}

	return nil
}

// GetAction retrieves a single action
func (fr Repository) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "GetAction")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actionCollection := fr.getActionsCollection(uid, flavour)
	el, err := fr.getSingleElement(
		ctx,
		actionCollection,
		actionID,
		&feedlib.Action{},
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	action, ok := el.(*feedlib.Action)
	if !ok {
		return nil, fmt.Errorf("expected an action, got %T", el)
	}

	return action, nil
}

// SaveAction saves a modified action
// It's expected to have an ID and sequence number already.
// One suggestion is to use UUIDs for IDs and UTC Unix Epoch seconds for
// sequence numbers.
func (fr Repository) SaveAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "SaveAction")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if action == nil {
		return nil, fmt.Errorf("nil action")
	}

	_, err := action.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("action failed validation: %w", err)
	}

	coll := fr.getActionsCollection(uid, flavour)
	if err := fr.saveElement(
		ctx,
		action,
		action.ID,
		action.SequenceNumber,
		coll,
		true,
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save action: %w", err)
	}

	return action, nil
}

// DeleteAction deletes an action from a user's feed
func (fr Repository) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	ctx, span := tracer.Start(ctx, "DeleteAction")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getActionsCollection(uid, flavour).Doc(actionID).Delete(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't delete action: %w", err)
	}

	return nil
}

// PostMessage adds a message or reply to an item's thread
func (fr Repository) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "PostMessage")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	coll := fr.getMessagesCollection(uid, flavour, itemID)
	if err := fr.saveElement(
		ctx,
		message,
		message.ID,
		message.SequenceNumber,
		coll,
		true,
	); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save message: %w", err)
	}

	return message, nil
}

// GetMessages gets the conversation thread for a single item
func (fr Repository) GetMessages(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) ([]feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "GetMessages")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	messages := []feedlib.Message{}
	seenMessageIDs := []string{}
	query := fr.getMessagesQuery(uid, flavour, itemID)
	msgDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get messages: %w", err)
	}
	for _, msgDoc := range msgDocs {
		msg := &feedlib.Message{}
		err := msgDoc.DataTo(msg)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to unmarshal message from firebase doc: %w", err)
		}
		if !converterandformatter.StringSliceContains(seenMessageIDs, msg.ID) {
			messages = append(messages, *msg)
			if msg.Timestamp.IsZero() {
				msg.Timestamp = time.Now() // backwards compat after schema change
			}
			seenMessageIDs = append(seenMessageIDs, msg.ID)
		}
	}
	return messages, nil
}

// GetMessage retrieves a message
func (fr Repository) GetMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) (*feedlib.Message, error) {
	ctx, span := tracer.Start(ctx, "GetMessage")
	defer span.End()
	messageCollection := fr.getMessagesCollection(uid, flavour, itemID)
	el, err := fr.getSingleElement(
		ctx,
		messageCollection,
		messageID,
		&feedlib.Message{},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get message: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	message, ok := el.(*feedlib.Message)
	if !ok {
		return nil, fmt.Errorf("expected a message, got %T", el)
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now() // backwards compat after schema change
	}
	return message, nil
}

// DeleteMessage removes a specific message
func (fr Repository) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	ctx, span := tracer.Start(ctx, "DeleteMessage")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getMessagesCollection(uid, flavour, itemID).
		Doc(messageID).
		Delete(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't delete message: %w", err)
	}

	return nil
}

// SaveIncomingEvent saves events that have been received from clients
// before they are processed further
func (fr Repository) SaveIncomingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	ctx, span := tracer.Start(ctx, "SaveIncomingEvent")
	defer span.End()
	if event == nil {
		return fmt.Errorf("nil event")
	}

	_, err := event.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("event failed validation: %w", err)
	}

	collectionName := firebasetools.SuffixCollection(incomingEventsCollectionName)
	coll := fr.firestoreClient.Collection(collectionName)
	doc := coll.Doc(event.ID)
	_, err = doc.Set(ctx, event)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to save event: %w", err)
	}
	return nil
}

// SaveOutgoingEvent saves events that are to be sent to clients
// before they are sent
func (fr Repository) SaveOutgoingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	ctx, span := tracer.Start(ctx, "SaveOutgoingEvent")
	defer span.End()
	if event == nil {
		return fmt.Errorf("nil event")
	}

	_, err := event.ValidateAndMarshal()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("event failed validation: %w", err)
	}

	collectionName := firebasetools.SuffixCollection(outgoingEventsCollectionName)
	coll := fr.firestoreClient.Collection(collectionName)
	doc := coll.Doc(event.ID)
	_, err = doc.Set(ctx, event)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to save event: %w", err)
	}
	return nil
}

func (fr Repository) getFeedCollectionName() string {
	suffixed := firebasetools.SuffixCollection(feedCollectionName)
	return suffixed
}

func (fr Repository) getMaretingSMSCollectionName() string {
	suffixed := firebasetools.SuffixCollection(AITMarketingMessageName)
	return suffixed
}

func (fr Repository) getNotificationCollectionName() string {
	suffixed := firebasetools.SuffixCollection(notificationCollectionName)
	return suffixed
}

func (fr Repository) getNPSResponseCollectionName() string {
	suffixed := firebasetools.SuffixCollection(NPSResponseCollectionName)
	return suffixed
}

func (fr Repository) getTwilioCallbackCollectionName() string {
	suffixed := firebasetools.SuffixCollection(twilioCallbackCollectionName)
	return suffixed
}

func (fr Repository) getOutgoingEmailsCollectionName() string {
	return firebasetools.SuffixCollection(outgoingEmails)
}

func (fr Repository) getUserCollection(
	uid string,
	flavour feedlib.Flavour,
) *firestore.CollectionRef {
	feedCollection := fr.firestoreClient.Collection(fr.getFeedCollectionName())
	userCollection := feedCollection.Doc(flavour.String()).Collection(uid)
	return userCollection
}

func (fr Repository) getElementCollection(
	uid string,
	flavour feedlib.Flavour,
	subCollectionName string,
) *firestore.CollectionRef {
	return fr.getUserCollection(
		uid,
		flavour,
	).Doc(elementsGroupName).Collection(subCollectionName)
}

func (fr Repository) getActionsCollection(
	uid string,
	flavour feedlib.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, actionsSubcollectionName)
}

func (fr Repository) getNudgesCollection(
	uid string,
	flavour feedlib.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, nudgesSubcollectionName)
}

func (fr Repository) getItemsCollection(
	uid string,
	flavour feedlib.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, itemsSubcollectionName)
}

func (fr Repository) getMessagesCollection(
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) *firestore.CollectionRef {
	itemsColl := fr.getElementCollection(uid, flavour, itemsSubcollectionName)
	messagesColl := itemsColl.Doc(itemID).Collection(messagesSubcollectionName)
	return messagesColl
}

func (fr Repository) getTwilioVideoCallbackCollectionName() string {
	suffixed := firebasetools.SuffixCollection(twilioVideoCallbackCollectionName)
	return suffixed
}

func (fr Repository) elementExists(
	ctx context.Context,
	collection *firestore.CollectionRef,
	id string,
	sequenceNumber int,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "elementExists")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	query := collection.Where(
		"id", "==", id,
	).Where(
		"sequenceNumber", "==", sequenceNumber,
	).LimitToLast(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to fetch firestore docs: %w", err)
	}
	return len(docs) > 0, nil
}
func (fr Repository) getItemsQuery(
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) (*firestore.Query, error) {
	itemsQuery := fr.getItemsCollection(
		uid, flavour,
	).Query.OrderBy(
		"expiry", firestore.Desc,
	).OrderBy(
		"id", firestore.Desc,
	).OrderBy(
		"sequenceNumber", firestore.Desc,
	).Limit(itemsLimit)

	if status == nil {
		itemsQuery = itemsQuery.Where(
			"status",
			"==",
			feedlib.StatusPending,
		)
	}

	if visibility == nil {
		itemsQuery = itemsQuery.Where(
			"visibility",
			"==",
			feedlib.VisibilityShow,
		)
	}

	if expired == nil {
		itemsQuery = itemsQuery.Where(
			"expiry",
			">=",
			time.Now(),
		)
	}

	switch persistent {
	case feedlib.BooleanFilterTrue:
		itemsQuery = itemsQuery.Where("persistent", "==", true)
	case feedlib.BooleanFilterFalse:
		itemsQuery = itemsQuery.Where("persistent", "==", false)

		// feedlib.BooleanFilterBoth is "passed": no filters added to the query
	}

	if status != nil {
		itemsQuery = itemsQuery.Where("status", "==", status)
	}
	if visibility != nil {
		itemsQuery = itemsQuery.Where("visibility", "==", visibility)
	}
	if expired != nil {
		if *expired == feedlib.BooleanFilterFalse {
			itemsQuery = itemsQuery.Where("expiry", ">=", time.Now())
		}

		if *expired == feedlib.BooleanFilterTrue {
			itemsQuery = itemsQuery.Where("expiry", "<=", time.Now())
		}
	}
	if filterParams != nil && len(filterParams.Labels) > 0 {
		itemsQuery = itemsQuery.Where("label", "in", filterParams.Labels)
	}

	return &itemsQuery, nil
}

// GetItems fetches feed items
func (fr Repository) GetItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]feedlib.Item, error) {
	ctx, span := tracer.Start(ctx, "GetItems")
	defer span.End()
	query, err := fr.getItemsQuery(
		uid,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to compose items query: %w", err)
	}

	items := []feedlib.Item{}
	seenItemIDs := []string{}
	itemDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	for _, itemDoc := range itemDocs {
		item := &feedlib.Item{}
		err := itemDoc.DataTo(item)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to unmarshal item from firebase doc: %w", err)
		}
		if !converterandformatter.StringSliceContains(seenItemIDs, item.ID) {
			messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
			if err != nil {
				helpers.RecordSpanError(span, err)
				return nil, fmt.Errorf("can't get feed item messages: %w", err)
			}
			item.Conversations = messages
			items = append(items, *item)
			seenItemIDs = append(seenItemIDs, item.ID)
		}
	}
	return items, nil
}

// GetActions retrieves the actions that a single feed has
func (fr Repository) GetActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]feedlib.Action, error) {
	ctx, span := tracer.Start(ctx, "GetActions")
	defer span.End()
	actions := []feedlib.Action{}
	seenActionIDs := []string{}

	query := fr.getActionsQuery(uid, flavour)
	actionDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	for _, actionDoc := range actionDocs {
		action := &feedlib.Action{}
		err := actionDoc.DataTo(action)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to unmarshal action from firebase doc: %w", err)
		}
		if !converterandformatter.StringSliceContains(
			seenActionIDs,
			action.ID,
		) {
			actions = append(actions, *action)
			seenActionIDs = append(seenActionIDs, action.ID)
		}
	}
	return actions, nil
}

// Labels retrieves the labels
func (fr Repository) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	ctx, span := tracer.Start(ctx, "Labels")
	defer span.End()
	labelDoc := fr.getUserCollection(uid, flavour).Doc(labelsDocID)
	lDoc, err := labelDoc.Get(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		if status.Code(err) != codes.NotFound {
			return nil, fmt.Errorf("error fetching labels collection: %w", err)
		}
		// create it if not found
		defaultLabel := map[string][]string{
			"labels": {common.DefaultLabel},
		}
		_, err := labelDoc.Set(ctx, defaultLabel)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf("can't set default label: %w", err)
		}

		// recurse after saving initial label
		return fr.Labels(ctx, uid, flavour)
	}

	var labels map[string][]string
	err = lDoc.DataTo(&labels)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"can't unmarshal labels from Firestore doc to list: %w", err)
	}

	lls, prs := labels["labels"]
	if !prs {
		return nil, fmt.Errorf("no `labels` key in labels doc")
	}

	return lls, nil
}

// SaveLabel saves the indicated label, if it does not already exist
func (fr Repository) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	ctx, span := tracer.Start(ctx, "SaveLabel")
	defer span.End()
	labels, err := fr.Labels(ctx, uid, flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to retrieve labels: %w", err)
	}

	if !converterandformatter.StringSliceContains(labels, label) {
		labelDoc := fr.getUserCollection(uid, flavour).Doc(labelsDocID)
		l := map[string][]string{
			"labels": {label},
		}
		_, err := labelDoc.Set(ctx, l)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return fmt.Errorf("can't save label: %w", err)
		}
	}

	return nil
}

// UnreadPersistentItems fetches unread persistent items
func (fr Repository) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	ctx, span := tracer.Start(ctx, "UnreadPersistentItems")
	defer span.End()
	// unreadInboxCountsDocID
	unreadDoc := fr.getUserCollection(uid, flavour).Doc(unreadInboxCountsDocID)
	uDoc, err := unreadDoc.Get(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		if status.Code(err) != codes.NotFound {
			return -1, fmt.Errorf(
				"error fetching unread docs collection: %w",
				err,
			)
		}
		// create it if not found
		defaultCount := map[string]int{
			"count": 0,
		}
		_, err := unreadDoc.Set(ctx, defaultCount)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return -1, fmt.Errorf("can't set default unread count: %w", err)
		}

		// recurse after saving initial unread count
		return fr.UnreadPersistentItems(ctx, uid, flavour)
	}

	var counts map[string]int
	err = uDoc.DataTo(&counts)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return -1, fmt.Errorf(
			"can't unmarshal unread counts from Firestore doc to list: %w",
			err,
		)
	}

	count, present := counts["count"]
	if !present {
		return -1, fmt.Errorf("`count` key not present in unread counts doc")
	}

	return count, nil
}

// UpdateUnreadPersistentItemsCount updates the unread inbox count
func (fr Repository) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	ctx, span := tracer.Start(ctx, "UpdateUnreadPersistentItemsCount")
	defer span.End()
	unreadDoc := fr.getUserCollection(uid, flavour).Doc(unreadInboxCountsDocID)

	persistentItemsQ, err := fr.getItemsQuery(
		uid, flavour, feedlib.BooleanFilterTrue, nil, nil, nil, nil)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't compose persistent items query: %w", err)
	}

	// (todo) nn
	// getting a count via iteration is VERY expensive because each firestore
	// doc iterated over is charged as a read.
	// this logic MUST be replaced soon (Dec 2020)
	persistentCount := 0
	it := persistentItemsQ.Documents(ctx)
	for {
		_, err = it.Next()
		if err != nil && err != iterator.Done {
			return fmt.Errorf("error iterating over persistent items: %w", err)
		}
		persistentCount++
		if err == iterator.Done {
			break
		}
	}

	count := map[string]int{
		"count": persistentCount,
	}
	_, err = unreadDoc.Set(ctx, count)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't set unread count: %w", err)
	}

	return nil
}

func (fr Repository) getNudgesQuery(
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) *firestore.Query {
	nudgesQuery := fr.getNudgesCollection(
		uid, flavour,
	).Query.OrderBy(
		"expiry", firestore.Desc,
	).OrderBy(
		"id", firestore.Desc,
	).OrderBy(
		"sequenceNumber", firestore.Desc,
	)

	if status == nil {
		nudgesQuery = nudgesQuery.Where(
			"status",
			"==",
			feedlib.StatusPending,
		)
	}

	if visibility == nil {
		nudgesQuery = nudgesQuery.Where(
			"visibility",
			"==",
			feedlib.VisibilityShow,
		)
	}

	if expired == nil {
		nudgesQuery = nudgesQuery.Where(
			"expiry",
			">=",
			time.Now(),
		)
	}

	if status != nil {
		nudgesQuery = nudgesQuery.Where("status", "==", status)
	}
	if visibility != nil {
		nudgesQuery = nudgesQuery.Where("visibility", "==", visibility)
	}
	if expired != nil {
		if *expired == feedlib.BooleanFilterFalse {
			nudgesQuery = nudgesQuery.Where("expiry", ">=", time.Now())
		}

		if *expired == feedlib.BooleanFilterTrue {
			nudgesQuery = nudgesQuery.Where("expiry", "<=", time.Now())
		}
	}
	return &nudgesQuery
}

func (fr Repository) getMessagesQuery(
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) *firestore.Query {
	messagesQuery := fr.getMessagesCollection(uid, flavour, itemID).
		Query.OrderBy(
		"id", firestore.Desc).
		OrderBy("sequenceNumber", firestore.Desc)
	return &messagesQuery
}

// GetNudges fetches nudges from the database
func (fr Repository) GetNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) ([]feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "GetNudges")
	defer span.End()
	nudges := []feedlib.Nudge{}
	seenNudgeIDs := []string{}

	query := fr.getNudgesQuery(
		uid,
		flavour,
		status,
		visibility,
		expired,
	)
	nudgeDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	for _, nudgeDoc := range nudgeDocs {
		nudge := &feedlib.Nudge{}
		err := nudgeDoc.DataTo(nudge)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to unmarshal nudge from firebase doc: %w", err)
		}
		if !converterandformatter.StringSliceContains(seenNudgeIDs, nudge.ID) {
			nudges = append(nudges, *nudge)
			seenNudgeIDs = append(seenNudgeIDs, nudge.ID)
		}
	}
	return nudges, nil
}

func (fr Repository) getActionsQuery(
	uid string,
	flavour feedlib.Flavour,
) *firestore.Query {
	query := fr.getActionsCollection(uid, flavour).Query.OrderBy(
		"id", firestore.Desc,
	).OrderBy("sequenceNumber", firestore.Desc)
	return &query
}

func (fr Repository) getSingleElement(
	ctx context.Context,
	collection *firestore.CollectionRef,
	id string,
	el feedlib.Element,
) (feedlib.Element, error) {
	ctx, span := tracer.Start(ctx, "getSingleElement")
	defer span.End()
	query := orderAndLimitBySequence(collection.Where(
		"id", "==", id,
	))

	docs, err := fetchQueryDocs(ctx, query, true)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"unable to get element with ID %s: %w", id, err)
	}

	el, err = docToElement(docs[0], el)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"unable to unmarshal feed item from doc snapshot: %w", err)
	}

	return el, nil
}

func (fr Repository) saveElement(
	ctx context.Context,
	el feedlib.Element,
	id string,
	sequenceNumber int,
	coll *firestore.CollectionRef,
	isNewElement bool,
) error {
	ctx, span := tracer.Start(ctx, "saveElement")
	defer span.End()
	if err := validateElement(el); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("%T failed validation: %w", el, err)
	}

	if isNewElement {
		exists, err := fr.elementExists(ctx, coll, id, sequenceNumber)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return fmt.Errorf("can't determine if item exists: %w", err)
		}

		if exists {
			return fmt.Errorf(
				"an element with the same ID and sequence number exists")
		}
	}

	doc := coll.Doc(id)
	_, err := doc.Set(ctx, el)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to save item: %w", err)
	}

	return nil
}

func validateElement(el feedlib.Element) error {
	if el == nil {
		return fmt.Errorf("failed validation: nil element")
	}

	_, err := el.ValidateAndMarshal()
	return err
}

func orderAndLimitBySequence(query firestore.Query) firestore.Query {
	return query.OrderBy(
		"sequenceNumber", firestore.Desc,
	).LimitToLast(1)
}

func fetchQueryDocs(
	ctx context.Context,
	query firestore.Query,
	requireAtLeastOne bool,
) ([]*firestore.DocumentSnapshot, error) {
	ctx, span := tracer.Start(ctx, "fetchQueryDocs")
	defer span.End()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"unable to fetch documents: %w", err)
	}
	if requireAtLeastOne && len(docs) == 0 {
		return nil, fmt.Errorf("expected at least one matching document")
	}
	return docs, nil
}

func docToElement(
	doc *firestore.DocumentSnapshot,
	el feedlib.Element,
) (feedlib.Element, error) {
	if el == nil {
		return nil, fmt.Errorf("nil element")
	}
	if doc == nil {
		return nil, fmt.Errorf("nil firestore document snapshot")
	}
	if !isPointer(el) {
		return nil, fmt.Errorf("%T is not a pointer", el)
	}
	err := doc.DataTo(el)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal feed item from doc snapshot: %w", err)
	}
	return el, err
}

func isPointer(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Ptr
}

// GetDefaultNudgeByTitle returns a default nudge given its title
func (fr Repository) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	ctx, span := tracer.Start(ctx, "GetDefaultNudgeByTitle")
	defer span.End()
	collection := fr.getNudgesCollection(uid, flavour)
	query := collection.Where("title", "==", title)
	nudgeDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get nudge: %w", err)
	}

	if len(nudgeDocs) == 0 {
		return nil, exceptions.ErrNilNudge
	}

	var nudge *feedlib.Nudge
	for _, nudgeDoc := range nudgeDocs {
		nudgeData := &feedlib.Nudge{}
		err = nudgeDoc.DataTo(nudgeData)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"unable to unmarshal nudge from firebase doc: %w", err)
		}
		nudge = nudgeData
	}
	return nudge, nil
}

// SaveMarketingMessage saves SMS data for future analysis
func (fr Repository) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	ctx, span := tracer.Start(ctx, "SaveMarketingMessage")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("repository precondition check failed: %w", err)
	}

	collectionName := fr.getMaretingSMSCollectionName()
	_, _, err := fr.firestoreClient.Collection(collectionName).Add(ctx, data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to save callback response")
	}

	return &data, nil
}

func getMarketingSMSData(
	ctx context.Context,
	atleastOne bool,
	query firestore.Query,
) (*dto.MarketingSMS, error) {
	ctx, span := tracer.Start(ctx, "getMarketingSMSData")
	defer span.End()
	docs, err := fetchQueryDocs(ctx, query, atleastOne)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	var marketingSMS dto.MarketingSMS
	err = docs[0].DataTo(&marketingSMS)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"unable to unmarshal marketing SMS from doc snapshot: %w", err)
	}

	return &marketingSMS, nil
}

// GetMarketingSMSByID returns a message given its id
func (fr Repository) GetMarketingSMSByID(
	ctx context.Context,
	id string,
) (*dto.MarketingSMS, error) {
	ctx, span := tracer.Start(ctx, "GetMarketingSMSByID")
	defer span.End()
	query := fr.firestoreClient.Collection(fr.getMaretingSMSCollectionName()).
		Where("ID", "==", id)
	return getMarketingSMSData(ctx, true, query)
}

// GetMarketingSMSByPhone returns the latest message given a phone number
func (fr Repository) GetMarketingSMSByPhone(
	ctx context.Context,
	phoneNumber string,
) (*dto.MarketingSMS, error) {
	ctx, span := tracer.Start(ctx, "GetMarketingSMSByPhone")
	defer span.End()
	query := fr.firestoreClient.Collection(fr.getMaretingSMSCollectionName()).
		Where("PhoneNumber", "==", phoneNumber).
		OrderBy("MessageSentTimeStamp", firestore.Desc)
	return getMarketingSMSData(ctx, true, query)
}

// UpdateMarketingMessage updates marketing SMS data
func (fr Repository) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	ctx, span := tracer.Start(ctx, "UpdateMarketingMessage")
	defer span.End()
	query := fr.firestoreClient.Collection(fr.getMaretingSMSCollectionName()).
		Where("ID", "==", data.ID)

	docs, err := fetchQueryDocs(ctx, query, true)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	doc := fr.firestoreClient.Collection(fr.getMaretingSMSCollectionName()).
		Doc(docs[0].Ref.ID)
	if _, err = doc.Set(ctx, data); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	return fr.GetMarketingSMSByID(ctx, data.ID)
}

// SaveTwilioResponse saves the callback data
func (fr Repository) SaveTwilioResponse(
	ctx context.Context,
	data dto.Message,
) error {
	ctx, span := tracer.Start(ctx, "SaveTwilioResponse")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("repository precondition check failed: %w", err)
	}

	collectionName := fr.getTwilioCallbackCollectionName()
	_, _, err := fr.firestoreClient.Collection(collectionName).Add(ctx, data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to save callback response")
	}

	return nil
}

// SaveNotification saves a notification
func (fr Repository) SaveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	notification dto.SavedNotification,
) error {
	ctx, span := tracer.Start(ctx, "SaveNotification")
	defer span.End()
	collectionName := fr.getNotificationCollectionName()
	_, _, err := firestoreClient.Collection(collectionName).
		Add(ctx, notification)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't save notification: %w", err)
	}
	return nil
}

// RetrieveNotification retrieves a notification
func (fr Repository) RetrieveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	ctx, span := tracer.Start(ctx, "RetrieveNotification")
	defer span.End()
	collectionName := fr.getNotificationCollectionName()

	docs, err := firestoreClient.Collection(
		collectionName,
	).Where(
		"RegistrationToken", "==", registrationToken,
	).Where(
		"Timestamp", ">=", newerThan,
	).Limit(limit).Documents(ctx).GetAll()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to retrieve notifications: %w", err)
	}
	notifications := []*dto.SavedNotification{}
	for _, doc := range docs {
		var notification dto.SavedNotification
		err = doc.DataTo(&notification)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, fmt.Errorf(
				"error unmarshalling saved notification: %w",
				err,
			)
		}
		notifications = append(notifications, &notification)
	}
	return notifications, nil
}

// SaveNPSResponse  stores the nps responses
func (fr Repository) SaveNPSResponse(
	ctx context.Context,
	response *dto.NPSResponse,
) error {
	ctx, span := tracer.Start(ctx, "SaveNPSResponse")
	defer span.End()
	collection := fr.getNPSResponseCollectionName()
	_, _, err := fr.firestoreClient.Collection(collection).Add(ctx, response)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't save nps response: %w", err)
	}
	return nil
}

// SaveOutgoingEmails saves all the outgoing emails
func (fr Repository) SaveOutgoingEmails(
	ctx context.Context,
	payload *dto.OutgoingEmailsLog,
) error {
	ctx, span := tracer.Start(ctx, "SaveOutgoingEmails")
	defer span.End()

	collection := fr.getOutgoingEmailsCollectionName()
	_, _, err := fr.firestoreClient.Collection(collection).Add(ctx, payload)
	if err != nil {
		return fmt.Errorf("unable to save ougoing email logs")
	}
	return nil
}

// UpdateMailgunDeliveryStatus updates the status and delivery time of the sent email message
func (fr Repository) UpdateMailgunDeliveryStatus(
	ctx context.Context,
	payload *dto.MailgunEvent,
) (*dto.OutgoingEmailsLog, error) {
	ctx, span := tracer.Start(ctx, "UpdateMailgunDeliveryStatus")
	defer span.End()

	collection := fr.getOutgoingEmailsCollectionName()

	query := fr.firestoreClient.Collection(collection).
		Where("messageID", "==", payload.MessageID)

	docs, err := fetchQueryDocs(ctx, query, true)
	if err != nil {
		return nil, err
	}

	var emailLogData dto.OutgoingEmailsLog
	err = docs[0].DataTo(&emailLogData)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to marshal data from docs snapshot: %w",
			err,
		)
	}

	timeOfDelivery := helpers.EpochTimetoStandardTime(payload.DeliveredOn)

	eventOutput := dto.MailgunEventOutput{
		EventName:   payload.EventName,
		DeliveredOn: timeOfDelivery,
	}

	emailLogData.Event = &eventOutput

	doc := fr.firestoreClient.Collection(fr.getOutgoingEmailsCollectionName()).
		Doc(docs[0].Ref.ID)
	if _, err = doc.Set(ctx, emailLogData); err != nil {
		logrus.Print("an error occurred while updating data: ", err)
		return nil, fmt.Errorf("unable to update data")
	}

	return &emailLogData, nil
}

// SaveTwilioVideoCallbackStatus saves the callback data
func (fr Repository) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	ctx, span := tracer.Start(ctx, "SaveTwilioVideoCallbackStatus")
	defer span.End()
	if err := fr.checkPreconditions(); err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("repository precondition check failed: %w", err)
	}

	collectionName := fr.getTwilioVideoCallbackCollectionName()
	_, _, err := fr.firestoreClient.Collection(collectionName).Add(ctx, data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to save callback response")
	}

	return nil
}
