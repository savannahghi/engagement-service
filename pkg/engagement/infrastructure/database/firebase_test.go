package database_test

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/apiclient"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	db "gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
)

const (
	intMax         = 9007199254740990
	sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"
)

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getTestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestEvent() feedlib.Event {
	return feedlib.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: feedlib.Context{
			UserID:         ksuid.New().String(),
			Flavour:        feedlib.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		ActionType:     feedlib.ActionTypePrimary,
		Handling:       feedlib.HandlingFullPage,
	}
}

func testNudge() *feedlib.Nudge {
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelEmail,
			feedlib.ChannelFcm,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func composeMarketingDataPayload(initialSegment, wing, phoneNumber, email string) apiclient.Segment {
	return apiclient.Segment{
		Properties: domain.ContactProperties{
			BeWellEnrolled:        "NO",
			OptOut:                "NO",
			BeWellAware:           "NO",
			BeWellPersona:         "SLADER",
			HasWellnessCard:       "YES",
			HasCover:              "YES",
			Payor:                 "Jubilee Insuarance Kenya",
			FirstChannelOfContact: "SMS",
			InitialSegment:        initialSegment,
			HasVirtualCard:        "NO",
			Email:                 email,
			Phone:                 phoneNumber,
			FirstName:             gofakeit.FirstName(),
			LastName:              gofakeit.LastName(),
		},

		Wing:        wing,
		MessageSent: "FALSE",
		IsSynced:    "FALSE",
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
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer
	status := feedlib.StatusPending
	visibility := feedlib.VisibilityHide
	expired := feedlib.BooleanFilterFalse

	type args struct {
		uid          string
		isAnonymous  bool
		flavour      feedlib.Flavour
		persistent   feedlib.BooleanFilter
		status       *feedlib.Status
		visibility   *feedlib.Visibility
		expired      *feedlib.BooleanFilter
		filterParams *helpers.FilterParams
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
				uid:         uid,
				isAnonymous: false,
				flavour:     flavour,
				persistent:  feedlib.BooleanFilterBoth,
			},
			wantErr:            false,
			wantInitialization: true,
		},
		{
			name: "with filter params",
			args: args{
				uid:         uid,
				isAnonymous: false,
				flavour:     flavour,
				persistent:  feedlib.BooleanFilterFalse,
				status:      &status,
				visibility:  &visibility,
				expired:     &expired,
				filterParams: &helpers.FilterParams{
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
				&tt.args.uid,
				&tt.args.isAnonymous,
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
						&tt.args.uid,
						&tt.args.isAnonymous,
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
					persistentItemCount := len(refetchedFeed.Items)
					refetchedActions := len(refetchedFeed.Actions)

					if refetchedActions != initialActions {
						t.Errorf("initially got %d actions, refetched and got %d", initialActions, refetchedActions)
					}

					if refetchedNudges != initialNudges {
						t.Errorf("initially got %d nudges, refetched and got %d", initialActions, refetchedActions)
					}

					if persistentItemCount != initialItems {
						t.Errorf("initially got %d items, refetched and got %d", initialItems, persistentItemCount)
					}

					// filter by 'persistent=TRUE'
					persistentFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterTrue,
						nil,
						nil,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the persistent=TRUE filter: %s", err)
						return
					}
					if persistentFeed == nil {
						t.Errorf("nil feed when fetching with the persistent=TRUE filter")
						return
					}
					if len(persistentFeed.Items) < 1 {
						t.Errorf("expected at least one persistent feed item, got none")
						return
					}

					// filter by persistent=FALSE
					nonPersistentFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterFalse,
						nil,
						nil,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the persistent=FALSE filter: %s", err)
						return
					}
					if nonPersistentFeed == nil {
						t.Errorf("nil feed when fetching with the persistent=FALSE filter")
						return
					}
					if len(nonPersistentFeed.Items) < 1 {
						t.Errorf("expected at least one non-persistent feed item, got none")
						return
					}

					// filter by persistent=BOTH
					bothPersistentFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						nil,
						nil,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the persistent=BOTH filter: %s", err)
						return
					}
					if bothPersistentFeed == nil {
						t.Errorf("nil feed when fetching with the persistent=BOTH filter")
						return
					}
					if len(bothPersistentFeed.Items) < 1 {
						t.Errorf("expected at least one persistent=BOTH feed item, got none")
						return
					}

					// filter by visibility=SHOW
					show := feedlib.VisibilityShow
					hiddenFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						nil,
						&show,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the visibility=SHOW filter: %s", err)
						return
					}
					if hiddenFeed == nil {
						t.Errorf("nil feed when fetching with the visibility=SHOW filter")
						return
					}
					if len(hiddenFeed.Items) < 1 {
						t.Errorf("expected at least one visibiity=SHOW feed item, got none")
						return
					}

					// filter by visibility=HIDE
					hide := feedlib.VisibilityHide
					visibilityHideFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						nil,
						&hide,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the visibility=HIDE filter: %s", err)
						return
					}
					if visibilityHideFeed == nil {
						t.Errorf("nil feed when fetching with the visibility=HIDE filter")
						return
					}

					for _, item := range visibilityHideFeed.Items {
						if item.Visibility == feedlib.VisibilityHide {
							t.Errorf("unexpectedly found > 0 visibiity=HIDE feed items")
							return
						}
					}

					// filter by status pending
					pending := feedlib.StatusPending
					pendingFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the status=PENDING filter: %s", err)
						return
					}
					if pendingFeed == nil {
						t.Errorf("nil feed when fetching with the status=PENDING filter")
						return
					}
					if len(pendingFeed.Items) < 1 {
						t.Errorf("expected at least one status=PENDING feed item, got none")
						return
					}

					// filter by status done
					done := feedlib.StatusDone
					doneFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&done,
						&show,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the status=DONE filter: %s", err)
						return
					}
					if doneFeed == nil {
						t.Errorf("nil feed when fetching with the status=DONE filter")
						return
					}

					for _, item := range doneFeed.Items {
						if item.Status == feedlib.StatusDone {
							t.Errorf("expected no status=DONE feed item")
							return
						}
					}

					// filter for in progress feed items
					inProgress := feedlib.StatusInProgress
					inProgressFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&inProgress,
						&show,
						nil,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the status=IN_PROGRESS filter: %s", err)
						return
					}
					if inProgressFeed == nil {
						t.Errorf("nil feed when fetching with the status=IN_PROGRESS filter")
						return
					}

					for _, item := range inProgressFeed.Items {
						if item.Status == feedlib.StatusInProgress {
							t.Errorf("expected no status=IN PROGRESS feed item")
							return
						}
					}

					// filter by expired=BOTH
					both := feedlib.BooleanFilterBoth
					expiredBothFeed, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						&both,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the expired=BOTH filter: %s", err)
						return
					}
					if expiredBothFeed == nil {
						t.Errorf("nil feed when fetching with the expired=BOTH filter")
						return
					}
					if len(expiredBothFeed.Items) < 1 {
						t.Errorf("expected at least one expired=BOTH feed item, got none")
						return
					}

					// filter by expired=FALSE
					falseVal := feedlib.BooleanFilterFalse
					unexpiredFilter, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						&falseVal,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the expired=FALSE filter: %s", err)
						return
					}
					if unexpiredFilter == nil {
						t.Errorf("nil feed when fetching with the expired=FALSE filter")
						return
					}
					if len(unexpiredFilter.Items) < 1 {
						t.Errorf("expected at least one expired=FALSE feed item, got none")
						return
					}

					// filter by expired=TRUE
					trueVal := feedlib.BooleanFilterTrue
					expiredFilter, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						&trueVal,
						nil,
					)
					if err != nil {
						t.Errorf("error when fetching feed with the expired=TRUE filter: %s", err)
						return
					}
					if expiredFilter == nil {
						t.Errorf("nil feed when fetching with the expired=TRUE filter")
						return
					}

					for _, item := range expiredFilter.Items {
						if item.Expiry == time.Now() {
							t.Errorf("did not expect any expired=TRUE feed item")
							return
						}
					}

					// filter by welcome label
					welcomeLabelFilter, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						&falseVal,
						&helpers.FilterParams{
							Labels: []string{common.DefaultLabel},
						},
					)
					if err != nil {
						t.Errorf("error when fetching feed with the welcome label filter: %s", err)
						return
					}
					if welcomeLabelFilter == nil {
						t.Errorf("nil feed when fetching with the welcome label filter")
						return
					}
					if len(welcomeLabelFilter.Items) < 1 {
						t.Errorf("expected at least one feed item with the welcome label, got none")
						return
					}

					// filter by non existent welcome label
					nonExistentLabelFilter, err := fr.GetFeed(
						ctx,
						&tt.args.uid,
						&tt.args.isAnonymous,
						tt.args.flavour,
						feedlib.BooleanFilterBoth,
						&pending,
						&show,
						&falseVal,
						&helpers.FilterParams{
							Labels: []string{ksuid.New().String()},
						},
					)
					if err != nil {
						t.Errorf("error when fetching feed a non-existent label filter: %s", err)
						return
					}
					if nonExistentLabelFilter == nil {
						t.Errorf("nil feed when fetching with a non existent label filter")
						return
					}
					if len(nonExistentLabelFilter.Items) < 1 {
						t.Errorf("expected to find only ghost items")
						return
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
	flavour := feedlib.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
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
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
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
		flavour feedlib.Flavour
		item    *feedlib.Item
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
				flavour: feedlib.FlavourConsumer,
				item:    &proItem,
			},
			wantErr: false,
		},
		{
			name: "consumer item",
			args: args{
				uid:     ksuid.New().String(),
				flavour: feedlib.FlavourPro,
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
	flavour := feedlib.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		uid     string
		flavour feedlib.Flavour
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
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feedlib.Flavour
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
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	nudge2 := testNudge()
	fr.SaveNudge(ctx, uid, flavour, nudge2)

	type args struct {
		uid     string
		flavour feedlib.Flavour
		nudge   *feedlib.Nudge
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
		{
			name: "invalid case - save duplicate nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				nudge:   nudge2,
			},
			wantErr: true,
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
	// Teardown
	err = fr.DeleteNudge(ctx, uid, flavour, nudge2.ID)
	if err != nil {
		t.Errorf("teardown failed")
		return
	}
}

func TestFirebaseRepository_DeleteNudge(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		uid     string
		flavour feedlib.Flavour
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
	flavour := feedlib.FlavourPro

	action := getTestAction()
	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feedlib.Flavour
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
	flavour := feedlib.FlavourPro
	action := getTestAction()

	type args struct {
		uid     string
		flavour feedlib.Flavour
		action  *feedlib.Action
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
	flavour := feedlib.FlavourConsumer
	action := getTestAction()

	savedAction, err := fr.SaveAction(ctx, uid, flavour, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		uid      string
		flavour  feedlib.Flavour
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
	flavour := feedlib.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()

	type args struct {
		uid     string
		flavour feedlib.Flavour
		itemID  string
		message *feedlib.Message
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
	flavour := feedlib.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	item.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		item    *feedlib.Item
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
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	savedNudge.Text = "updated"

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudge   *feedlib.Nudge
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
	flavour := feedlib.FlavourConsumer

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
		flavour   feedlib.Flavour
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
		event *feedlib.Event
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
		event *feedlib.Event
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

func getTestItem() feedlib.Item {
	return feedlib.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon: feedlib.GetPNGImageLink(
			base.LogoURL, "title", "description", base.BlankImageURL),
		Author:    "Bot 1",
		Tagline:   "Bot speaks...",
		Label:     "DRUGS",
		Timestamp: time.Now(),
		Summary:   "I am a bot...",
		Text:      "This bot can speak",
		TextType:  feedlib.TextTypePlain,
		Links: []feedlib.Link{
			feedlib.GetYoutubeVideoLink(
				sampleVideoURL, "title", "description", base.BlankImageURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon: feedlib.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType: feedlib.ActionTypeSecondary,
				Handling:   feedlib.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon: feedlib.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType: feedlib.ActionTypePrimary,
				Handling:   feedlib.HandlingInline,
			},
		},
		Conversations: []feedlib.Message{
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
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelFcm,
			feedlib.ChannelEmail,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func TestRepository_Labels(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "default labels",
			args: args{
				ctx:     ctx,
				uid:     ksuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    []string{common.DefaultLabel},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.Labels(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Labels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Labels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_SaveLabel(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		label   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save label successfully",
			args: args{
				ctx:     ctx,
				uid:     ksuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				label:   ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveLabel(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UnreadPersistentItems(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "default - user with no persistent count",
			args: args{
				ctx:     ctx,
				uid:     ksuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UnreadPersistentItems(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UnreadPersistentItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.UnreadPersistentItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdateUnreadPersistentItemsCount(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default - user with no persistent count",
			args: args{
				ctx:     ctx,
				uid:     ksuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateUnreadPersistentItemsCount(tt.args.ctx, tt.args.uid, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateUnreadPersistentItemsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetDefaultNudgeByTitle(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		t.Errorf("can't save the nudge %v:", err)
		return
	}
	if savedNudge == nil {
		t.Errorf("nil saved nudge")
		return
	}

	type args struct {
		uid     string
		flavour feedlib.Flavour
		title   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case - get a nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				title:   nudge.Title,
			},
			wantErr: false,
		},
		{
			name: "sad case - get a non existent nudge",
			args: args{
				uid:     uid,
				flavour: flavour,
				title:   "non existent title",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nudge, err := fr.GetDefaultNudgeByTitle(ctx, tt.args.uid, tt.args.flavour, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetDefaultNudgeByTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && nudge == nil {
				t.Errorf("expected to get a nudge")
				return
			}
		})
	}
}

func TestRepository_GetNudges(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	if err != nil {
		t.Errorf("can't save the nudge %v:", err)
		return
	}
	if savedNudge == nil {
		t.Errorf("nil saved nudge")
		return
	}

	pending := feedlib.StatusPending
	show := feedlib.VisibilityShow

	type args struct {
		ctx        context.Context
		uid        string
		flavour    feedlib.Flavour
		status     *feedlib.Status
		visibility *feedlib.Visibility
		expired    *feedlib.BooleanFilter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: default logic",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "happy case: filters provided",
			args: args{
				ctx:        ctx,
				uid:        uid,
				flavour:    flavour,
				status:     &pending,
				visibility: &show,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nudges, err := fr.GetNudges(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.status,
				tt.args.visibility,
				tt.args.expired,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetNudges() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && nudges != nil {
				t.Errorf("nudge was not expected since an error occured: %v", err)
				return
			}

			if !tt.wantErr && nudges == nil {
				t.Errorf("nudge was expected since no error occured: %v", err)
				return
			}
		})
	}
}

func TestRepository_GetItems(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := getTestItem()
	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	item, err := fr.SaveFeedItem(ctx, uid, flavour, &testItem)
	if err != nil {
		t.Errorf("can't save the item %v:", err)
		return
	}
	if item == nil {
		t.Errorf("nil saved item")
		return
	}

	pending := feedlib.StatusPending
	show := feedlib.VisibilityShow

	type args struct {
		ctx          context.Context
		uid          string
		flavour      feedlib.Flavour
		persistent   feedlib.BooleanFilter
		status       *feedlib.Status
		visibility   *feedlib.Visibility
		expired      *feedlib.BooleanFilter
		filterParams *helpers.FilterParams
	}
	tests := []struct {
		name    string
		args    args
		want    []feedlib.Item
		wantErr bool
	}{
		{
			name: "happy case: default logic",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "happy case: filters provided",
			args: args{
				ctx:        ctx,
				uid:        uid,
				flavour:    flavour,
				status:     &pending,
				visibility: &show,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := fr.GetItems(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.persistent,
				tt.args.status,
				tt.args.visibility,
				tt.args.expired,
				tt.args.filterParams,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetNudges() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && items != nil {
				t.Errorf("nudge was not expected since an error occured: %v", err)
				return
			}

			if !tt.wantErr && items == nil {
				t.Errorf("nudge was expected since no error occured: %v", err)
				return
			}
		})
	}
}

func TestRepository_SaveNPSResponse(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	type args struct {
		ctx      context.Context
		response *dto.NPSResponse
	}

	feedback := &dto.Feedback{
		Question: "How is it",
		Answer:   "It is what it is",
	}
	email := converterandformatter.TestUserEmail
	phoneNumber := base.TestUserPhoneNumber

	response := &dto.NPSResponse{
		ID:        uuid.New().String(),
		Name:      "Test User",
		Score:     8,
		SladeCode: "123456723",
		Email:     &email,
		MSISDN:    &phoneNumber,
		Feedback:  []dto.Feedback{*feedback},
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:      ctx,
				response: response,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				response: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.SaveNPSResponse(tt.args.ctx, tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveNPSResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_RetrieveMarketingData(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase repository: %w", err)
		return
	}
	if fr == nil {
		t.Errorf("Firebase repository is nil: %v:", err)
		return
	}

	// Setup test data
	segment := composeMarketingDataPayload(
		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
		fmt.Sprintf("WING %s", ksuid.New().String()),
		gofakeit.PhoneFormatted(),
		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
	)
	_, err = fr.LoadMarketingData(ctx, segment)
	if err != nil {
		t.Errorf("Error loading marketing data: %v:", err)
		return
	}

	payload := dto.MarketingMessagePayload{
		InitialSegment: segment.Properties.InitialSegment,
		Wing:           segment.Wing,
	}

	payload2 := dto.MarketingMessagePayload{
		InitialSegment: fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
		Wing:           fmt.Sprintf("WING %s", ksuid.New().String()),
	}

	type args struct {
		ctx  context.Context
		data *dto.MarketingMessagePayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Fetch all data from Initial data",
			args: args{
				ctx:  ctx,
				data: &payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case - Missing segment ",
			args: args{
				ctx:  ctx,
				data: &payload2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fr.RetrieveMarketingData(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.RetrieveMarketingData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestRepository_UpdateMessageSentStatus(t *testing.T) {
// 	assert := assert.New(t)
// 	ctx := context.Background()
// 	repository, err := db.NewFirebaseRepository(ctx)
// 	if !assert.Nilf(err, "Error initializing Firebase repository: %s", err) {
// 		return
// 	}
// 	if !assert.NotNil(repository, "nil Firebase repository") {
// 		return
// 	}

// 	// Setup test data
// 	segment := composeMarketingDataPayload(
// 		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
// 		fmt.Sprintf("WING %s", ksuid.New().String()),
// 		gofakeit.PhoneFormatted(),
// 		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
// 	)
// 	_, err = repository.LoadMarketingData(ctx, segment)
// 	if !assert.Nilf(err, "Error loading marketing data: %s", err) {
// 		return
// 	}

// 	payload1 := dto.MarketingMessagePayload{
// 		InitialSegment: segment.Properties.InitialSegment,
// 		Wing:           segment.Wing,
// 	}
// 	payload2 := dto.MarketingMessagePayload{
// 		InitialSegment: fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
// 		Wing:           fmt.Sprintf("WING %s", ksuid.New().String()),
// 	}

// 	// Make sure that segment was loaded in the repository
// 	segments, err := repository.RetrieveMarketingData(ctx, &payload1)
// 	if !assert.Nilf(err, "Error, unable to retrieve loaded marketing data: %s", err) {
// 		return
// 	}
// 	if !assert.Equalf(len(segments), 1, "Error, expected exactly 1 segment with wing '%s'", err) {
// 		return
// 	}

// 	tests := []struct {
// 		name    string
// 		payload dto.MarketingMessagePayload
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Update message sent status of an existing segment",
// 			payload: payload1,
// 			wantErr: true, // TODO: fix and make false
// 		},
// 		{
// 			name:    "Update message sent status of an non-existing segment",
// 			payload: payload2,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err = repository.UpdateMessageSentStatus(ctx, segment.Properties.Phone, segment.Properties.InitialSegment)
// 			assert.False(!tt.wantErr && err != nil, "Error, unable to update message sent status: %s", err)
// 		})
// 	}

// 	// Teardown test data
// 	err = repository.RollBackMarketingData(ctx, segment)
// 	if !assert.Nilf(err, "Error, unable to roll back market data: %s", err) {
// 		return
// 	}
// }

func TestRepository_LoadMarketingData(t *testing.T) {
	ctx := context.Background()
	repository, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to initialize Firebase repository: %s", err)
		return
	}
	if repository == nil {
		t.Errorf("nil Firebase repository returned")
		return
	}

	// Create a new test segment
	marketingData := composeMarketingDataPayload(
		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
		fmt.Sprintf("WING %s", ksuid.New().String()),
		gofakeit.PhoneFormatted(),
		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
	)

	tests := []struct {
		name          string
		marketingData apiclient.Segment
		wantStatus    int
	}{
		{
			name:          "Load a new segment",
			marketingData: marketingData,
			wantStatus:    -1,
		},
		{
			name:          "Load an existing segment",
			marketingData: marketingData,
			wantStatus:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := repository.LoadMarketingData(ctx, marketingData)
			if err != nil {
				t.Errorf("failed to load new marketing data: %s", err)
				return
			}
			if status != tt.wantStatus {
				t.Errorf("expected load status %v, but got %v instead", tt.wantStatus, status)
				return
			}
		})
	}

	// Cleanup test data
	err = repository.RollBackMarketingData(ctx, marketingData)
	if err != nil {
		t.Errorf("failed to clean test data: %s", err)
		return
	}
}

func TestRepository_RollBackMarketingData(t *testing.T) {
	ctx := context.Background()
	repository, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to initialize Firebase repository: %s", err)
		return
	}
	if repository == nil {
		t.Errorf("nil Firebase repository returned")
		return
	}

	// Setup test data
	marketingData1 := composeMarketingDataPayload(
		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
		fmt.Sprintf("WING %s", ksuid.New().String()),
		gofakeit.PhoneFormatted(),
		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
	)
	marketingData2 := composeMarketingDataPayload(
		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
		fmt.Sprintf("WING %s", ksuid.New().String()),
		gofakeit.PhoneFormatted(),
		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
	)
	_, err = repository.LoadMarketingData(ctx, marketingData1)
	if err != nil {
		t.Errorf("failed to setup test data: %s", err)
	}

	tests := []struct {
		name          string
		marketingData apiclient.Segment
		wantErr       bool
	}{
		{
			name:          "Rollback an existing segment",
			marketingData: marketingData1,
			wantErr:       false,
		},
		{
			name:          "Rollback a non existing segment",
			marketingData: marketingData2,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repository.RollBackMarketingData(ctx, marketingData1)
			if !tt.wantErr && err != nil {
				t.Errorf("failed to rollback marketing data: %s", err)
				return
			}
		})
	}
}

func TestService_SaveOutgoingEmails(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("an error ocurred")
	}

	to := "kathurima@healthcloud.co.ke"
	subject := "Test subject"
	text := "Hello test"
	messageID := "123456"

	outgoingEmail := &dto.OutgoingEmailsLog{
		UUID:        uuid.NewString(),
		To:          []string{to},
		From:        mail.MailGunFromEnvVarName,
		Subject:     subject,
		Text:        text,
		MessageID:   messageID,
		EmailSentOn: time.Now(),
	}

	type args struct {
		ctx     context.Context
		payload *dto.OutgoingEmailsLog
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				payload: outgoingEmail,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				payload: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.SaveOutgoingEmails(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Service.MailgunDeliveryWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestRepository_UpdateMailgunDeliveryStatus(t *testing.T) {
// 	ctx := context.Background()
// 	fr, err := db.NewFirebaseRepository(ctx)
// 	if err != nil {
// 		t.Errorf("an error ocurred")
// 	}

// 	type args struct {
// 		ctx     context.Context
// 		payload *dto.MailgunEvent
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *dto.OutgoingEmailsLog
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				payload: &dto.MailgunEvent{
// 					EventName:   "delivered",
// 					DeliveredOn: "123456789.12456",
// 					MessageID:   "20210715172955.1.63EC29EF167F09B9@sandboxb30d61fba25641a9983c3b3a3c84abde.mailgun.org",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 				payload: &dto.MailgunEvent{
// 					EventName:   "delivered",
// 					DeliveredOn: "123456789.12456",
// 					MessageID:   "",
// 				},
// 			},
// 			want:    &dto.OutgoingEmailsLog{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_ = helpers.EpochTimetoStandardTime("123456789.12456")

// 			got, err := fr.UpdateMailgunDeliveryStatus(tt.args.ctx, tt.args.payload)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Repository.UpdateMailgunDeliveryStatus() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected to get a nudge")
// 				return
// 			}
// 		})
// 	}
// }

// func TestRepository_RetrieveSingleSladerData(t *testing.T) {
// 	ctx := context.Background()
// 	repository, err := db.NewFirebaseRepository(ctx)
// 	if err != nil {
// 		t.Errorf("failed to initialize Firebase repository: %s", err)
// 		return
// 	}
// 	if repository == nil {
// 		t.Errorf("nil Firebase repository returned")
// 		return
// 	}

// 	// Compose the Test user payload
// 	marketingData := composeMarketingDataPayload(
// 		fmt.Sprintf("SIL Segment %s", ksuid.New().String()),
// 		fmt.Sprintf("WING %s", ksuid.New().String()),
// 		base.TestUserPhoneNumber,
// 		fmt.Sprintf("test-%s@savannah.com", ksuid.New().String()),
// 	)

// 	// Create the test user
// 	_, err = repository.LoadMarketingData(ctx, marketingData)
// 	if err != nil {
// 		t.Errorf("failed to setup test data: %s", err)
// 	}

// 	type args struct {
// 		ctx         context.Context
// 		phonenumber string
// 	}

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy Case: Successfully retrieve a slader data",
// 			args: args{
// 				ctx:         ctx,
// 				phonenumber: base.TestUserPhoneNumber,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Happy Case: Attempt to retrieve non-existent data",
// 			args: args{
// 				ctx:         ctx,
// 				phonenumber: "",
// 			},
// 			// This should not return an error
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := repository.GetSladerDataByPhone(tt.args.ctx, tt.args.phonenumber)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Repository.RetrieveSingleSladerData() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// 	// Tear down the created user
// 	repository.RollBackMarketingData(ctx, marketingData)
// }
