package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/feed"
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
}

func (fr Repository) checkPreconditions() error {
	if fr.firestoreClient == nil {
		return fmt.Errorf("nil firestore client in feed firebase repository")
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
	uid string,
	flavour feed.Flavour,
	persistent feed.BooleanFilter,
	status *feed.Status,
	visibility *feed.Visibility,
	expired *feed.BooleanFilter,
	filterParams *feed.FilterParams,
) (*feed.Feed, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actions, err := fr.getActions(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}

	nudges, err := fr.getNudges(ctx, uid, flavour, status, visibility)
	if err != nil {
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}

	items, err := fr.getItems(
		ctx,
		uid,
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

	feed := &feed.Feed{
		UID:     uid,
		Flavour: flavour,
		Actions: actions,
		Nudges:  nudges,
		Items:   items,
	}
	return feed, nil
}

// GetFeedItem retrieves and returns a single feed item
func (fr Repository) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	itemID string,
) (*feed.Item, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	itemsCollection := fr.getItemsCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, itemsCollection, itemID, &feed.Item{})
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	item, ok := el.(*feed.Item)
	if !ok {
		return nil, fmt.Errorf("expected an Item, got %T", el)
	}

	messages, err := fr.getMessages(ctx, uid, flavour, itemID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []feed.Message{}
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
	flavour feed.Flavour,
	item *feed.Item,
) (*feed.Item, error) {
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

	messages, err := fr.getMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []feed.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// UpdateFeedItem updates an existing feed item
func (fr Repository) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	item *feed.Item,
) (*feed.Item, error) {
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

	messages, err := fr.getMessages(ctx, uid, flavour, item.ID)
	if err != nil || messages == nil {
		// the thread may not have been initiated yet
		item.Conversations = []feed.Message{}
	} else {
		item.Conversations = messages
	}

	return item, nil
}

// DeleteFeedItem deletes a nudge from a user's feed
func (fr Repository) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
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
	flavour feed.Flavour,
	nudgeID string,
) (*feed.Nudge, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	nudgeCollection := fr.getNudgesCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, nudgeCollection, nudgeID, &feed.Nudge{})
	if err != nil {
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	nudge, ok := el.(*feed.Nudge)
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
	flavour feed.Flavour,
	nudge *feed.Nudge,
) (*feed.Nudge, error) {
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
		true,
	); err != nil {
		return nil, fmt.Errorf("unable to save nudge: %w", err)
	}

	return nudge, nil
}

// UpdateNudge updates an existing nudge e.g to show or hide it
func (fr Repository) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	nudge *feed.Nudge,
) (*feed.Nudge, error) {
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
	flavour feed.Flavour,
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
	flavour feed.Flavour,
	actionID string,
) (*feed.Action, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	actionCollection := fr.getActionsCollection(uid, flavour)
	el, err := fr.getSingleElement(ctx, actionCollection, actionID, &feed.Action{})
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	action, ok := el.(*feed.Action)
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
	flavour feed.Flavour,
	action *feed.Action,
) (*feed.Action, error) {
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
	flavour feed.Flavour,
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
	flavour feed.Flavour,
	itemID string,
	message *feed.Message,
) (*feed.Message, error) {
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
	flavour feed.Flavour,
	itemID string,
) ([]feed.Message, error) {
	if err := fr.checkPreconditions(); err != nil {
		return nil, fmt.Errorf(
			"repository precondition check failed: %w", err)
	}

	messages := []feed.Message{}
	seenMessageIDs := []string{}
	query := fr.getMessagesQuery(uid, flavour, itemID)
	msgDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get messages: %w", err)
	}
	for _, msgDoc := range msgDocs {
		msg := &feed.Message{}
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
	flavour feed.Flavour,
	itemID string,
	messageID string,
) (*feed.Message, error) {
	messageCollection := fr.getMessagesCollection(uid, flavour, itemID)
	el, err := fr.getSingleElement(ctx, messageCollection, messageID, &feed.Message{})
	if err != nil {
		return nil, fmt.Errorf("unable to get message: %w", err)
	}
	if el == nil {
		return nil, nil
	}
	message, ok := el.(*feed.Message)
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
	flavour feed.Flavour,
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
	event *feed.Event,
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
	event *feed.Event,
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
	flavour feed.Flavour,
) *firestore.CollectionRef {
	feedCollection := fr.firestoreClient.Collection(fr.getFeedCollectionName())
	userCollection := feedCollection.Doc(flavour.String()).Collection(uid)
	return userCollection
}

func (fr Repository) getElementCollection(
	uid string,
	flavour feed.Flavour,
	subCollectionName string,
) *firestore.CollectionRef {
	return fr.getUserCollection(
		uid, flavour).Doc(elementsGroupName).Collection(subCollectionName)
}

func (fr Repository) getActionsCollection(
	uid string,
	flavour feed.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, actionsSubcollectionName)
}

func (fr Repository) getNudgesCollection(
	uid string,
	flavour feed.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, nudgesSubcollectionName)
}

func (fr Repository) getItemsCollection(
	uid string,
	flavour feed.Flavour,
) *firestore.CollectionRef {
	return fr.getElementCollection(uid, flavour, itemsSubcollectionName)
}

func (fr Repository) getMessagesCollection(
	uid string,
	flavour feed.Flavour,
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
	flavour feed.Flavour,
	persistent feed.BooleanFilter,
	status *feed.Status,
	visibility *feed.Visibility,
	expired *feed.BooleanFilter,
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

	switch persistent {
	case feed.BooleanFilterTrue:
		itemsQuery = itemsQuery.Where("persistent", "==", true)
	case feed.BooleanFilterFalse:
		itemsQuery = itemsQuery.Where("persistent", "==", false)
	}

	if status != nil {
		itemsQuery = itemsQuery.Where("status", "==", status)
	}
	if visibility != nil {
		itemsQuery = itemsQuery.Where("visibility", "==", visibility)
	}
	if expired != nil {
		if *expired == feed.BooleanFilterFalse {
			itemsQuery = itemsQuery.Where("expiry", ">=", time.Now())
		}

		if *expired == feed.BooleanFilterTrue {
			itemsQuery = itemsQuery.Where("expiry", "<=", time.Now())
		}
	}
	if filterParams != nil && len(filterParams.Labels) > 0 {
		itemsQuery = itemsQuery.Where("label", "in", filterParams.Labels)
	}

	return &itemsQuery, nil
}

func (fr Repository) getItems(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	persistent feed.BooleanFilter,
	status *feed.Status,
	visibility *feed.Visibility,
	expired *feed.BooleanFilter,
	filterParams *feed.FilterParams,
) ([]feed.Item, error) {
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

	items := []feed.Item{}
	seenItemIDs := []string{}
	itemDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}
	for _, itemDoc := range itemDocs {
		item := &feed.Item{}
		err := itemDoc.DataTo(item)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal item from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenItemIDs, item.ID) {
			messages, err := fr.getMessages(ctx, uid, flavour, item.ID)
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

func (fr Repository) getActionsQuery(
	uid string,
	flavour feed.Flavour,
) *firestore.Query {
	query := fr.getActionsCollection(uid, flavour).Query.OrderBy(
		"id", firestore.Desc).OrderBy("sequenceNumber", firestore.Desc)
	return &query
}

func (fr Repository) getActions(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
) ([]feed.Action, error) {
	actions := []feed.Action{}
	seenActionIDs := []string{}

	query := fr.getActionsQuery(uid, flavour)
	actionDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get actions: %w", err)
	}
	for _, actionDoc := range actionDocs {
		action := &feed.Action{}
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

func (fr Repository) getNudgesQuery(
	uid string,
	flavour feed.Flavour,
	status *feed.Status,
	visibility *feed.Visibility,
) *firestore.Query {
	nudgesQuery := fr.getNudgesCollection(uid, flavour).Query.OrderBy(
		"id", firestore.Desc).OrderBy("sequenceNumber", firestore.Desc)
	if status != nil {
		nudgesQuery = nudgesQuery.Where("status", "==", status)
	}
	if visibility != nil {
		nudgesQuery = nudgesQuery.Where("visibility", "==", visibility)
	}
	return &nudgesQuery
}

func (fr Repository) getMessagesQuery(
	uid string,
	flavour feed.Flavour,
	itemID string,
) *firestore.Query {
	messagesQuery := fr.getMessagesCollection(uid, flavour, itemID).Query.OrderBy(
		"id", firestore.Desc).OrderBy("sequenceNumber", firestore.Desc)
	return &messagesQuery
}

func (fr Repository) getNudges(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	status *feed.Status,
	visibility *feed.Visibility,
) ([]feed.Nudge, error) {
	nudges := []feed.Nudge{}
	seenNudgeIDs := []string{}

	query := fr.getNudgesQuery(
		uid,
		flavour,
		status,
		visibility,
	)
	nudgeDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get nudges: %w", err)
	}
	for _, nudgeDoc := range nudgeDocs {
		nudge := &feed.Nudge{}
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

func (fr Repository) getMessages(
	ctx context.Context,
	uid string,
	flavour feed.Flavour,
	itemID string,
) ([]feed.Message, error) {
	messages := []feed.Message{}
	seenMessageIDs := []string{}

	query := fr.getMessagesQuery(
		uid,
		flavour,
		itemID,
	)
	messageDocs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get messages: %w", err)
	}
	for _, msgDoc := range messageDocs {
		message := &feed.Message{}
		err := msgDoc.DataTo(message)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal message from firebase doc: %w", err)
		}
		if !base.StringSliceContains(seenMessageIDs, message.ID) {
			if message.Timestamp.IsZero() {
				message.Timestamp = time.Now() // backfill after schema change
			}
			messages = append(messages, *message)
			seenMessageIDs = append(seenMessageIDs, message.ID)
		}
	}
	return messages, nil
}

func (fr Repository) getSingleElement(
	ctx context.Context,
	collection *firestore.CollectionRef,
	id string,
	el feed.Element,
) (feed.Element, error) {
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
	el feed.Element,
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

func validateElement(el feed.Element) error {
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
	el feed.Element,
) (feed.Element, error) {
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
