package database

import (
	"context"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/stretchr/testify/assert"
)

func Test_partnerAccountSetupNudge(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	nudge, err := partnerAccountSetupNudge(ctx, emptyUID, consumer, fr)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)

	nudge, err = partnerAccountSetupNudge(ctx, emptyUID, pro, fr)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)
}

func Test_verifyEmailNudge(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	nudge, err := verifyEmailNudge(ctx, emptyUID, consumer, fr)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)

	nudge, err = verifyEmailNudge(ctx, emptyUID, pro, fr)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)
}

func Test_createNudge(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	title := "test title"
	text := "text"
	imageURL := "http://example.com/image.png"
	imageTitle := "test image"
	imageDescription := "test image description"
	actions := []feedlib.Action{}

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	notificationBody := feedlib.NotificationBody{}

	nudge, err := createNudge(
		ctx,
		emptyUID,
		consumer,
		title,
		text,
		imageURL,
		imageTitle,
		imageDescription,
		actions,
		fr,
		notificationBody,
	)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)

	nudge, err = createNudge(ctx,
		emptyUID,
		pro,
		title,
		text,
		imageURL,
		imageTitle,
		imageDescription,
		actions,
		fr,
		notificationBody,
	)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)
}

func Test_createGlobalAction(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	allowAnonymous := false
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	name := "test title"
	actionType := feedlib.ActionTypePrimary
	handling := feedlib.HandlingInline
	iconLink := "http://example.com/image.png"
	iconTitle := "test image"
	iconDescription := "test image description"

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	nudge, err := createGlobalAction(
		ctx,
		emptyUID,
		allowAnonymous,
		consumer,
		name,
		actionType,
		handling,
		iconLink,
		iconTitle,
		iconDescription,
		fr,
	)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)

	nudge, err = createGlobalAction(
		ctx,
		emptyUID,
		allowAnonymous,
		pro,
		name,
		actionType,
		handling,
		iconLink,
		iconTitle,
		iconDescription,
		fr,
	)
	assert.Empty(t, nudge)
	assert.NotNil(t, err)
}

func Test_createFeedItem(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	itemID := "test"
	author := "test"
	tagline := "test"
	label := "test"
	iconImageURL := "test"
	iconTitle := "test"
	iconDescription := "test"
	summary := "test"
	text := "test"
	links := []feedlib.Link{}
	actions := []feedlib.Action{}
	conversations := []feedlib.Message{}
	persistent := false

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	feed, err := createFeedItem(
		ctx,
		emptyUID,
		consumer,
		itemID,
		author,
		tagline,
		label,
		iconImageURL,
		iconTitle,
		iconDescription,
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		fr,
	)
	assert.Empty(t, feed)
	assert.NotNil(t, err)

	feed, err = createFeedItem(
		ctx,
		emptyUID,
		pro,
		itemID,
		author,
		tagline,
		label,
		iconImageURL,
		iconTitle,
		iconDescription,
		summary,
		text,
		links,
		actions,
		conversations,
		persistent,
		fr,
	)
	assert.Empty(t, feed)
	assert.NotNil(t, err)
}

func Test_simpleConsumerWelcome(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	welcome, err := simpleConsumerWelcome(ctx, emptyUID, consumer, fr)
	assert.Empty(t, welcome)
	assert.NotNil(t, err)

	welcome, err = simpleConsumerWelcome(ctx, emptyUID, pro, fr)
	assert.Empty(t, welcome)
	assert.NotNil(t, err)
}

func Test_simpleProWelcome(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	welcome, err := simpleProWelcome(ctx, emptyUID, consumer, fr)
	assert.Empty(t, welcome)
	assert.NotNil(t, err)

	welcome, err = simpleProWelcome(ctx, emptyUID, pro, fr)
	assert.Empty(t, welcome)
	assert.NotNil(t, err)
}

func Test_getMessage(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	itemID := "test"
	text := "test"
	replyTo := feedlib.Message{}
	postedByName := "text"

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	message, err := getMessage(
		ctx,
		emptyUID,
		consumer,
		itemID,
		text,
		&replyTo,
		postedByName,
		fr,
	)
	assert.Empty(t, message)
	assert.NotNil(t, err)

	message, err = getMessage(
		ctx,
		emptyUID,
		pro,
		itemID,
		text,
		&replyTo,
		postedByName,
		fr,
	)
	assert.Empty(t, message)
	assert.NotNil(t, err)
}

func Test_getConsumerWelcomeThread(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	itemID := "test"

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	message, err := getConsumerWelcomeThread(
		ctx,
		emptyUID,
		consumer,
		itemID,
		fr,
	)
	assert.Empty(t, message)
	assert.NotNil(t, err)

	message, err = getConsumerWelcomeThread(
		ctx,
		emptyUID,
		pro,
		itemID,
		fr)
	assert.Empty(t, message)
	assert.NotNil(t, err)
}

func Test_getProWelcomeThread(t *testing.T) {
	ctx := context.Background()
	emptyUID := ""
	consumer := feedlib.FlavourConsumer
	pro := feedlib.FlavourPro
	itemID := "test"

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("unable to create FirebaseRepository: %v", err)
	}

	message, err := getProWelcomeThread(
		ctx,
		emptyUID,
		consumer,
		itemID,
		fr,
	)
	assert.Empty(t, message)
	assert.NotNil(t, err)

	message, err = getProWelcomeThread(
		ctx,
		emptyUID,
		pro,
		itemID,
		fr,
	)
	assert.Empty(t, message)
	assert.NotNil(t, err)
}
