package db_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
)

const (
	intMax         = 9007199254740990
	sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"
)

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getTestMessage() feed.Message {
	return feed.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestEvent() feed.Event {
	return feed.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: feed.Context{
			UserID:         ksuid.New().String(),
			Flavour:        feed.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() feed.Action {
	return feed.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.BlankImageURL),
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
	}
}

func testNudge() *feed.Nudge {
	return &feed.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feed.Link{
			feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.BlankImageURL),
		},
		Text: ksuid.New().String(),
		Actions: []feed.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func TestNewFirebaseRepository(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case - should succeed",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.NewFirebaseRepository(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewFirebaseRepository() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_GetFeed(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	status := feed.StatusPending
	visibility := feed.VisibilityHide
	expired := feed.BooleanFilterFalse

	type args struct {
		uid          string
		flavour      feed.Flavour
		persistent   feed.BooleanFilter
		status       *feed.Status
		visibility   *feed.Visibility
		expired      *feed.BooleanFilter
		filterParams *feed.FilterParams
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		wantInitialization bool
	}{
		{
			name: "no filter params",
			args: args{
				uid:        uid,
				flavour:    flavour,
				persistent: feed.BooleanFilterBoth,
			},
			wantErr:            false,
			wantInitialization: true,
		},
		{
			name: "with filter params",
			args: args{
				uid:        uid,
				flavour:    flavour,
				persistent: feed.BooleanFilterFalse,
				status:     &status,
				visibility: &visibility,
				expired:    &expired,
				filterParams: &feed.FilterParams{
					Labels: []string{ksuid.New().String()},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialFeed, err := fr.GetFeed(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.persistent,
				tt.args.status,
				tt.args.visibility,
				tt.args.expired,
				tt.args.filterParams,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetFeed() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && initialFeed == nil {
				t.Errorf("nil feed")
				return
			}

			if tt.wantInitialization {
				// re-fetch, ensure it does not change in counts
				initialNudges := len(initialFeed.Nudges)
				initialItems := len(initialFeed.Items)
				initialActions := len(initialFeed.Actions)

				if initialActions < 1 {
					t.Errorf("zero initial actions")
				}

				if initialItems < 1 {
					t.Errorf("zero initial items")
				}

				if initialNudges < 1 {
					t.Errorf("zero initial nudges")
				}

				for range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
					refetchedFeed, err := fr.GetFeed(
						ctx,
						tt.args.uid,
						tt.args.flavour,
						tt.args.persistent,
						tt.args.status,
						tt.args.visibility,
						tt.args.expired,
						tt.args.filterParams,
					)
					if err != nil {
						t.Errorf("error when refetching feed: %s", err)
						return
					}
					if refetchedFeed == nil {
						t.Errorf("nil refetched feed")
						return
					}

					refetchedNudges := len(refetchedFeed.Nudges)
					refetchedItems := len(refetchedFeed.Items)
					refetchedActions := len(refetchedFeed.Actions)

					if refetchedActions != initialActions {
						t.Errorf("initially got %d actions, refetched and got %d", initialActions, refetchedActions)
					}

					if refetchedNudges != initialNudges {
						t.Errorf("initially got %d nudges, refetched and got %d", initialActions, refetchedActions)
					}

					if refetchedItems != initialItems {
						t.Errorf("initially got %d items, refetched and got %d", initialItems, refetchedItems)
					}
				}
			}
		})
	}
}

func TestFirebaseRepository_GetFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "get back saved feed items",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent feed item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				itemID:  ksuid.New().String(),
			},
			wantErr: false,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetFeedItem(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantNil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.itemID, got.ID)
			}
		})
	}
}

func TestFirebaseRepository_SaveFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	proItem := getTestItem()
	consumerItem := getTestItem()

	type args struct {
		uid     string
		flavour feed.Flavour
		item    *feed.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pro item",
			args: args{
				uid:     ksuid.New().String(),
				flavour: feed.FlavourConsumer,
				item:    &proItem,
			},
			wantErr: false,
		},
		{
			name: "consumer item",
			args: args{
				uid:     ksuid.New().String(),
				flavour: feed.FlavourPro,
				item:    &consumerItem,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveFeedItem(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.item,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)

				bs, err := got.ValidateAndMarshal()
				assert.Nil(t, err)
				assert.NotNil(t, bs)
			}
		})
	}
}

func TestFirebaseRepository_DeleteFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		uid     string
		flavour feed.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing feed item",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing feed item",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteFeedItem(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteFeedItem() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_GetNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feed.Flavour
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - retrieves successfully",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: savedNudge.ID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudgeID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_SaveNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	type args struct {
		uid     string
		flavour feed.Flavour
		nudge   *feed.Nudge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - save nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudge:   nudge,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudge,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_DeleteNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feed.Flavour
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudgeID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteNudge(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudgeID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteNudge() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_GetAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourPro

	action := getTestAction()
	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feed.Flavour
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - will save",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.actionID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.GetAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_SaveAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourPro
	action := getTestAction()

	type args struct {
		uid     string
		flavour feed.Flavour
		action  *feed.Action
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case, should save",
			args: args{
				uid:     uid,
				flavour: flavour,
				action:  &action,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SaveAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.action,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.SaveAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_DeleteAction(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	action := getTestAction()

	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feed.Flavour
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing action",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
		{
			name: "non existing action",
			args: args{
				uid:      uid,
				flavour:  flavour,
				actionID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteAction(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.actionID,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.DeleteAction() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFirebaseRepository_PostMessage(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()

	type args struct {
		uid     string
		flavour feed.Flavour
		itemID  string
		message *feed.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case - should save",
			args: args{
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
				message: &message,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.PostMessage(
				ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.message,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FirebaseRepository.PostMessage() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFirebaseRepository_UpdateFeedItem(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	item.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		item    *feed.Item
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "valid case - will update",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				item:    item,
			},
			wantText: "updated",
			wantErr:  false,
		},
		{
			name: "error case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				item:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UpdateFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.UpdateFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantText, got.Text)
			}
		})
	}
}

func TestFirebaseRepository_UpdateNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	savedNudge.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feed.Flavour
		nudge   *feed.Nudge
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "valid case - update an existing nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				nudge:   savedNudge,
			},
			wantText: "updated",
			wantErr:  false,
		},
		{
			name: "nil nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				nudge:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UpdateNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudge)
			if (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.UpdateNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantText, got.Text)
			}
		})
	}
}

func TestFirebaseRepository_DeleteMessage(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()
	savedMessage, err := fr.PostMessage(ctx, uid, flavour, item.ID, &message)
	assert.Nil(t, err)
	assert.NotNil(t, savedMessage)

	type args struct {
		ctx       context.Context
		uid       string
		flavour   feed.Flavour
		itemID    string
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing message, should delete",
			args: args{
				ctx:       ctx,
				uid:       uid,
				flavour:   flavour,
				itemID:    item.ID,
				messageID: savedMessage.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent message, should not error",
			args: args{
				ctx:       ctx,
				uid:       uid,
				flavour:   flavour,
				itemID:    ksuid.New().String(),
				messageID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.DeleteMessage(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.messageID,
			); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFirebaseRepository_SaveIncomingEvent(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *feed.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:   ctx,
				event: &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:   ctx,
				event: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveIncomingEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.SaveIncomingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFirebaseRepository_SaveOutgoingEvent(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *feed.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:   ctx,
				event: &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:   ctx,
				event: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveOutgoingEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("FirebaseRepository.SaveOutgoingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestItem() feed.Item {
	return feed.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon: feed.GetPNGImageLink(
			feed.LogoURL, "title", "description", feed.BlankImageURL),
		Author:    "Bot 1",
		Tagline:   "Bot speaks...",
		Label:     "DRUGS",
		Timestamp: time.Now(),
		Summary:   "I am a bot...",
		Text:      "This bot can speak",
		TextType:  feed.TextTypePlain,
		Links: []feed.Link{
			feed.GetYoutubeVideoLink(
				sampleVideoURL, "title", "description", feed.BlankImageURL),
		},
		Actions: []feed.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon: feed.GetPNGImageLink(
					feed.LogoURL, "title", "description", feed.BlankImageURL),
				ActionType: feed.ActionTypeSecondary,
				Handling:   feed.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon: feed.GetPNGImageLink(
					feed.LogoURL, "title", "description", feed.BlankImageURL),
				ActionType: feed.ActionTypePrimary,
				Handling:   feed.HandlingInline,
			},
		},
		Conversations: []feed.Message{
			{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
			},
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelFcm,
			feed.ChannelEmail,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}
