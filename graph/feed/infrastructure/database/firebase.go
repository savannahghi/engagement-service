package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	feedCollectionName           = "feed"
	elementsGroupName            = "elements"
	actionsSubcollectionName     = "actions"
	nudgesSubcollectionName      = "nudges"
	itemsSubcollectionName       = "items"
	messagesSubcollectionName    = "messages"
	incomingEventsCollectionName = "incoming_events"
	outgoingEventsCollectionName = "outgoing_events"

	labelsDocID            = "item_labels"
	unreadInboxCountsDocID = "unread_inbox_counts"

	itemsLimit = 1000
)

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(ctx context.Context) (feed.Repository, error) {
	fc := base.FirebaseClient{}
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
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *feed.FilterParams,
) (*feed.Feed, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actions, err := fr.GetActions(ctx, *uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}

	nudges, err := fr.GetNudges(ctx, *uid, flavour, status, visibility, expired)
	if err != nil {
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
		return nil, fmt.Errorf("unable to get items: %w", err)
	}

	// only add default content if...
	// - the `persistent` filter is set to "BOTH"
	// - all other filters are nil
	noFilters := persistent == base.BooleanFilterBoth && status == nil && visibility == nil && expired == nil && filterParams == nil
	noActions := len(actions) == 0
	noNudges := len(nudges) == 0
	noItems := len(items) == 0
	if noFilters && noActions && noNudges && noItems {
		err = fr.initializeDefaultFeed(ctx, *uid, flavour)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize default feed: %w", err)
		}

		// this recursion is potentially dangerous but there's an integration test
		// that exercises this and reduces the risk of infinite recursion
		// we need to do this in order to have confidence that the initialization succeeded
		return fr.GetFeed(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
	}

	feed := &feed.Feed{
		UID:         *uid,
		Flavour:     flavour,
		Actions:     actions,
		Nudges:      nudges,
		Items:       items,
		IsAnonymous: isAnonymous,
	}

	return feed, nil
}

func (fr Repository) initializeDefaultFeed(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) error {
	fr.mu.Lock() // create default data once

	_, err := feed.SetDefaultActions(ctx, uid, flavour, fr)
	if err != nil {
		return fmt.Errorf("unable to set default actions: %w", err)
	}

	_, err = feed.SetDefaultNudges(ctx, uid, flavour, fr)
	if err != nil {
		return fmt.Errorf("unable to set default nudges: %w", err)
	}

	_, err = feed.SetDefaultItems(ctx, uid, flavour, fr)
	if err != nil {
		return fmt.Errorf("unable to set default items: %w", err)
	}

	fr.mu.Unlock() // release lock on default data
	return nil
}

// GetFeedItem retrieves and returns a single feed item
func (fr Repository) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) (*base.Item, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	itemsCollection := fr.getItemsCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, itemsCollection, itemID, &base.Item{})
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	item, ok := el.(*base.Item)
	if !ok {
		return nil, fmt.Errorf("expected an Item, got %T", el)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, itemID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []base.Message{}
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
	flavour base.Flavour,
	item *base.Item,
) (*base.Item, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	_, err := item.ValidateAndMarshal()
	if err != nil {
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
		return nil, fmt.Errorf("unable to save item: %w", err)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []base.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// UpdateFeedItem updates an existing feed item
func (fr Repository) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	item *base.Item,
) (*base.Item, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if item == nil {
		return nil, fmt.Errorf("nil item")
	}

	_, err := item.ValidateAndMarshal()
	if err != nil {
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
		return nil, fmt.Errorf("unable to save item: %w", err)
	}

	messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []base.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// DeleteFeedItem deletes a nudge from a user's feed
func (fr Repository) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) error {
	if err := fr.checkPreconditions(); err != nil {
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getItemsCollection(uid, flavour).Doc(itemID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("can't delete item: %w", err)
	}

	return nil
}

// GetNudge retrieves a single nudge
func (fr Repository) GetNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) (*base.Nudge, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	nudgeCollection := fr.getNudgesCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, nudgeCollection, nudgeID, &base.Nudge{})
	if err != nil {
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	nudge, ok := el.(*base.Nudge)
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
	flavour base.Flavour,
	nudge *base.Nudge,
) (*base.Nudge, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
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
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}

	return nudge, nil
}

// UpdateNudge updates an existing nudge e.g to show or hide it
func (fr Repository) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudge *base.Nudge,
) (*base.Nudge, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if nudge == nil {
		return nil, fmt.Errorf("nil nudge")
	}

	_, err := nudge.ValidateAndMarshal()
	if err != nil {
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
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}

	return nudge, nil
}

// DeleteNudge deletes a nudge from a user's feed
func (fr Repository) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	nudgeID string,
) error {
	if err := fr.checkPreconditions(); err != nil {
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getNudgesCollection(uid, flavour).Doc(nudgeID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("can't delete nudge: %w", err)
	}

	return nil
}

// GetAction retrieves a single action
func (fr Repository) GetAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) (*base.Action, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actionCollection := fr.getActionsCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, actionCollection, actionID, &base.Action{})
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	action, ok := el.(*base.Action)
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
	flavour base.Flavour,
	action *base.Action,
) (*base.Action, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	if action == nil {
		return nil, fmt.Errorf("nil action")
	}

	_, err := action.ValidateAndMarshal()
	if err != nil {
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
		return nil, fmt.Errorf("unable to save action: %w", err)
	}

	return action, nil
}

// DeleteAction deletes an action from a user's feed
func (fr Repository) DeleteAction(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	actionID string,
) error {
	if err := fr.checkPreconditions(); err != nil {
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getActionsCollection(uid, flavour).Doc(actionID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("can't delete action: %w", err)
	}

	return nil
}

// PostMessage adds a message or reply to an item's thread
func (fr Repository) PostMessage(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
	message *base.Message,
) (*base.Message, error) {
	if err := fr.checkPreconditions(); err != nil {
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
		return nil, fmt.Errorf("unable to save message: %w", err)
	}

	return message, nil
}

// GetMessages gets the conversation thread for a single item
func (fr Repository) GetMessages(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	itemID string,
) ([]base.Message, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	messages := []base.Message{}
	seenMessageIDs := []string{}
	query := fr.getMessagesQuery(uid, flavour, itemID)
	msgDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get messages: %w", err)
	}
	for _, msgDoc := range msgDocs {
		msg := &base.Message{}
		err := msgDoc.DataTo(msg)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal message from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenMessageIDs, msg.ID) {
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
	flavour base.Flavour,
	itemID string,
	messageID string,
) (*base.Message, error) {
	messageCollection := fr.getMessagesCollection(uid, flavour, itemID)
	el, err := fr.getSingleElement(ctx, messageCollection, messageID, &base.Message{})
	if err != nil {
		return nil, fmt.Errorf("unable to get message: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	message, ok := el.(*base.Message)
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
	flavour base.Flavour,
	itemID string,
	messageID string,
) error {
	if err := fr.checkPreconditions(); err != nil {
		return fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	_, err := fr.getMessagesCollection(uid, flavour, itemID).Doc(messageID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}

	return nil
}

// SaveIncomingEvent saves events that have been received from clients
// before they are processed further
func (fr Repository) SaveIncomingEvent(
	ctx context.Context,
	event *base.Event,
) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	_, err := event.ValidateAndMarshal()
	if err != nil {
		return fmt.Errorf("event failed validation: %w", err)
	}

	collectionName := base.SuffixCollection(incomingEventsCollectionName)
	coll := fr.firestoreClient.Collection(collectionName)
	doc := coll.Doc(event.ID)
	_, err = doc.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("unable to save event: %w", err)
	}
	return nil
}

// SaveOutgoingEvent saves events that are to be sent to clients
// before they are sent
func (fr Repository) SaveOutgoingEvent(
	ctx context.Context,
	event *base.Event,
) error {
	if event == nil {
		return fmt.Errorf("nil event")
	}

	_, err := event.ValidateAndMarshal()
	if err != nil {
		return fmt.Errorf("event failed validation: %w", err)
	}

	collectionName := base.SuffixCollection(outgoingEventsCollectionName)
	coll := fr.firestoreClient.Collection(collectionName)
	doc := coll.Doc(event.ID)
	_, err = doc.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("unable to save event: %w", err)
	}
	return nil
}

func (fr Repository) getFeedCollectionName() string {
	suffixed := base.SuffixCollection(feedCollectionName)
	return suffixed
}

func (fr Repository) getUserCollection(
	uid string,
	flavour base.Flavour,
) *firestore.CollectionRef {
	feedCollection := fr.firestoreClient.Collection(fr.getFeedCollectionName())
	userCollection := feedCollection.Doc(flavour.String()).Collection(uid)
	return userCollection
}

func (fr Repository) getElementCollection(
	uid string,
	flavour base.Flavour,
	subCollectionName string,
) *firestore.CollectionRef {
	return fr.getUserCollection(
		uid,
		flavour,
	).Doc(elementsGroupName).Collection(subCollectionName)
}

func (fr Repository) getActionsCollection(
	uid string,
	flavour base.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, actionsSubcollectionName)
}

func (fr Repository) getNudgesCollection(
	uid string,
	flavour base.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, nudgesSubcollectionName)
}

func (fr Repository) getItemsCollection(
	uid string,
	flavour base.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, itemsSubcollectionName)
}

func (fr Repository) getMessagesCollection(
	uid string,
	flavour base.Flavour,
	itemID string,
) *firestore.CollectionRef {
	itemsColl := fr.getElementCollection(uid, flavour, itemsSubcollectionName)
	messagesColl := itemsColl.Doc(itemID).Collection(messagesSubcollectionName)
	return messagesColl
}

func (fr Repository) elementExists(
	ctx context.Context,
	collection *firestore.CollectionRef,
	id string,
	sequenceNumber int,
) (bool, error) {
	if err := fr.checkPreconditions(); err != nil {
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
		return false, fmt.Errorf("unable to fetch firestore docs: %w", err)
	}
	return len(docs) > 0, nil
}
func (fr Repository) getItemsQuery(
	uid string,
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *feed.FilterParams,
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
			base.StatusPending,
		)
	}

	if visibility == nil {
		itemsQuery = itemsQuery.Where(
			"visibility",
			"==",
			base.VisibilityShow,
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
	case base.BooleanFilterTrue:
		itemsQuery = itemsQuery.Where("persistent", "==", true)
	case base.BooleanFilterFalse:
		itemsQuery = itemsQuery.Where("persistent", "==", false)

		// base.BooleanFilterBoth is "passed": no filters added to the query
	}

	if status != nil {
		itemsQuery = itemsQuery.Where("status", "==", status)
	}
	if visibility != nil {
		itemsQuery = itemsQuery.Where("visibility", "==", visibility)
	}
	if expired != nil {
		if *expired == base.BooleanFilterFalse {
			itemsQuery = itemsQuery.Where("expiry", ">=", time.Now())
		}

		if *expired == base.BooleanFilterTrue {
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
	flavour base.Flavour,
	persistent base.BooleanFilter,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
	filterParams *feed.FilterParams,
) ([]base.Item, error) {
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
		return nil, fmt.Errorf("unable to compose items query: %w", err)
	}

	items := []base.Item{}
	seenItemIDs := []string{}
	itemDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	for _, itemDoc := range itemDocs {
		item := &base.Item{}
		err := itemDoc.DataTo(item)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal item from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenItemIDs, item.ID) {
			messages, err := fr.GetMessages(ctx, uid, flavour, item.ID)
			if err != nil {
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
	flavour base.Flavour,
) ([]base.Action, error) {
	actions := []base.Action{}
	seenActionIDs := []string{}

	query := fr.getActionsQuery(uid, flavour)
	actionDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	for _, actionDoc := range actionDocs {
		action := &base.Action{}
		err := actionDoc.DataTo(action)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal action from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenActionIDs, action.ID) {
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
	flavour base.Flavour,
) ([]string, error) {
	labelDoc := fr.getUserCollection(uid, flavour).Doc(labelsDocID)
	lDoc, err := labelDoc.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, fmt.Errorf("error fetching labels collection: %w", err)
		}
		// create it if not found
		defaultLabel := map[string][]string{
			"labels": {feed.DefaultLabel},
		}
		_, err := labelDoc.Set(ctx, defaultLabel)
		if err != nil {
			return nil, fmt.Errorf("can't set default label: %w", err)
		}

		// recurse after saving initial label
		return fr.Labels(ctx, uid, flavour)
	}

	var labels map[string][]string
	err = lDoc.DataTo(&labels)
	if err != nil {
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
	flavour base.Flavour,
	label string,
) error {
	labels, err := fr.Labels(ctx, uid, flavour)
	if err != nil {
		return fmt.Errorf("unable to retrieve labels: %w", err)
	}

	if !base.StringSliceContains(labels, label) {
		labelDoc := fr.getUserCollection(uid, flavour).Doc(labelsDocID)
		l := map[string][]string{
			"labels": {label},
		}
		_, err := labelDoc.Set(ctx, l)
		if err != nil {
			return fmt.Errorf("can't save label: %w", err)
		}
	}

	return nil
}

// UnreadPersistentItems fetches unread persistent items
func (fr Repository) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
) (int, error) {
	// unreadInboxCountsDocID
	unreadDoc := fr.getUserCollection(uid, flavour).Doc(unreadInboxCountsDocID)
	uDoc, err := unreadDoc.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return -1, fmt.Errorf("error fetching unread docs collection: %w", err)
		}
		// create it if not found
		defaultCount := map[string]int{
			"count": 0,
		}
		_, err := unreadDoc.Set(ctx, defaultCount)
		if err != nil {
			return -1, fmt.Errorf("can't set default unread count: %w", err)
		}

		// recurse after saving initial unread count
		return fr.UnreadPersistentItems(ctx, uid, flavour)
	}

	var counts map[string]int
	err = uDoc.DataTo(&counts)
	if err != nil {
		return -1, fmt.Errorf(
			"can't unmarshal unread counts from Firestore doc to list: %w", err)
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
	flavour base.Flavour,
) error {
	unreadDoc := fr.getUserCollection(uid, flavour).Doc(unreadInboxCountsDocID)

	persistentItemsQ, err := fr.getItemsQuery(
		uid, flavour, base.BooleanFilterTrue, nil, nil, nil, nil)
	if err != nil {
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
		return fmt.Errorf("can't set unread count: %w", err)
	}

	return nil
}

func (fr Repository) getNudgesQuery(
	uid string,
	flavour base.Flavour,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
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
			base.StatusPending,
		)
	}

	if visibility == nil {
		nudgesQuery = nudgesQuery.Where(
			"visibility",
			"==",
			base.VisibilityShow,
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
		if *expired == base.BooleanFilterFalse {
			nudgesQuery = nudgesQuery.Where("expiry", ">=", time.Now())
		}

		if *expired == base.BooleanFilterTrue {
			nudgesQuery = nudgesQuery.Where("expiry", "<=", time.Now())
		}
	}
	return &nudgesQuery
}

func (fr Repository) getMessagesQuery(
	uid string,
	flavour base.Flavour,
	itemID string,
) *firestore.Query {
	messagesQuery := fr.getMessagesCollection(uid, flavour, itemID).Query.OrderBy(
		"id", firestore.Desc).OrderBy("sequenceNumber", firestore.Desc)
	return &messagesQuery
}

// GetNudges fetches nudges from the database
func (fr Repository) GetNudges(
	ctx context.Context,
	uid string,
	flavour base.Flavour,
	status *base.Status,
	visibility *base.Visibility,
	expired *base.BooleanFilter,
) ([]base.Nudge, error) {
	nudges := []base.Nudge{}
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
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	for _, nudgeDoc := range nudgeDocs {
		nudge := &base.Nudge{}
		err := nudgeDoc.DataTo(nudge)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal nudge from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenNudgeIDs, nudge.ID) {
			nudges = append(nudges, *nudge)
			seenNudgeIDs = append(seenNudgeIDs, nudge.ID)
		}
	}
	return nudges, nil
}

func (fr Repository) getActionsQuery(
	uid string,
	flavour base.Flavour,
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
	el base.Element,
) (base.Element, error) {
	query := orderAndLimitBySequence(collection.Where(
		"id", "==", id,
	))

	docs, err := fetchQueryDocs(ctx, query, true)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get element with ID %s: %w", id, err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	el, err = docToElement(docs[0], el)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal feed item from doc snapshot: %w", err)
	}

	return el, nil
}

func (fr Repository) saveElement(
	ctx context.Context,
	el base.Element,
	id string,
	sequenceNumber int,
	coll *firestore.CollectionRef,
	isNewElement bool,
) error {
	if err := validateElement(el); err != nil {
		return fmt.Errorf("%T failed validation: %w", el, err)
	}

	if isNewElement {
		exists, err := fr.elementExists(ctx, coll, id, sequenceNumber)
		if err != nil {
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
		return fmt.Errorf("unable to save item: %w", err)
	}

	return nil
}

func validateElement(el base.Element) error {
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
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to fetch documents: %w", err)
	}
	if requireAtLeastOne && len(docs) < 0 {
		return nil, fmt.Errorf("expected at least one matching document")
	}
	return docs, nil
}

func docToElement(
	doc *firestore.DocumentSnapshot,
	el base.Element,
) (base.Element, error) {
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
	flavour base.Flavour,
	title string,
) (*base.Nudge, error) {
	collection := fr.getNudgesCollection(uid, flavour)
	query := collection.Where("title", "==", title)
	nudgeDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get nudge: %w", err)
	}

	if len(nudgeDocs) == 0 {
		return nil, feed.ErrNilNudge
	}

	var nudge *base.Nudge
	for _, nudgeDoc := range nudgeDocs {
		nudgeData := &base.Nudge{}
		err = nudgeDoc.DataTo(nudgeData)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal nudge from firebase doc: %w", err)
		}
		nudge = nudgeData
	}
	return nudge, nil
}
