package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
)

func TestPublishFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}

	uid := ksuid.New().String()
	testItem := testItem()
	type args struct {
		ctx     context.Context
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
			name: "valid:publish_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    testItem,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_save_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    testItem,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_a_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    testItem,
			},
			wantErr: true,
		},
		{
			name: "invalid:use_a_nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    nil,
			},
			wantErr: true,
		},
		{
			name: "invalid:use_an_invalid_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    &feedlib.Item{},
			},
			wantErr: true,
		},
		{
			name: "invalid:use_an_invalid_action_type",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				item:    testItem,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:publish_feed_item" {
				fakeEngagement.SaveFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_save_feed_item" {
				fakeEngagement.SaveFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("unable to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_a_notification" {
				fakeEngagement.SaveFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send a notification")
				}

				if tt.name == "invalid:use_a_nil_item" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour feedlib.Flavour,
						item *feedlib.Item,
					) (*feedlib.Item, error) {
						return nil, fmt.Errorf("can't publish nil feed item")
					}
				}

				if tt.name == "invalid:use_an_invalid_item" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour feedlib.Flavour,
						item *feedlib.Item,
					) (*feedlib.Item, error) {
						return nil, fmt.Errorf("unable to publish feed item")
					}
				}

				if tt.name == "invalid:use_an_invalid_action_type" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour feedlib.Flavour,
						item *feedlib.Item,
					) (*feedlib.Item, error) {
						return &feedlib.Item{
							ID: uuid.New().String(),
							Actions: []feedlib.Action{
								{
									ID:         ksuid.New().String(),
									Name:       "TEST_ACTION",
									ActionType: feedlib.ActionTypeFloating,
								},
							},
						}, fmt.Errorf("floating actions are only allowed at the global level")
					}
				}
			}

			got, err := i.Feed.PublishFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.PublishFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}

		})
	}

}

func TestDeleteFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := testItem()

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
	}{
		{
			name: "valid:delete_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  "",
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_delete_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:delete_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) error {
					return nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("unable to retrieve feed item")
				}
			}

			if tt.name == "invalid:fail_to_delete_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) error {
					return fmt.Errorf("failed to delete feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) error {
					return nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send a notification")
				}
			}

			err := i.Feed.DeleteFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestDeleteFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestResolveFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

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
	}{
		{
			name: "valid:successfully_resolve_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_resolve_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ResolveItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			got, err := i.Feed.ResolveFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestResolveFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}

}

func TestPinFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantPersistent bool
		wantErr        bool
	}{
		{
			name: "valid:successfully_pin_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: true,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  "",
			},
			wantPersistent: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_pin_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to publish nil item")
				}
			}

			got, err := i.Feed.PinFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.PinFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestUnpinFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantPersistent bool
		wantErr        bool
	}{
		{
			name: "valid:successfully_unpin_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: true,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantPersistent: false,
			wantErr:        true,
		},
		{
			name: "invalid:nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  "",
			},
			wantPersistent: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_unpin_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.UnPinItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to publish nil item")
				}
			}

			got, err := i.Feed.UnpinFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.UnpinFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestUnresolveFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

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
	}{
		{
			name: "valid:successfully_unresolve_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_unresolve_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.UnResolveItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("nil item")
				}
			}

			got, err := i.Feed.UnresolveFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestUnresolveFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestHideFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "valid:successfully_hide_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantVisibility: feedlib.VisibilityHide,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_hide_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.HideItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			got, err := i.Feed.HideFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestHideFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestShowFeedItem(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	testItem := getTestItem()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "valid:successfully_show_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantVisibility: feedlib.VisibilityShow,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_show_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					itemID string,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					item *feedlib.Item,
				) (*feedlib.Item, error) {
					return &feedlib.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			got, err := i.Feed.ShowFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestShowFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestLabels(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

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
			name: "valid:successfully_return_valid_labels",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
			},
			want:    []string{common.DefaultLabel},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_return_valid_labels" {
				fakeEngagement.LabelsFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
				) ([]string, error) {
					return []string{common.DefaultLabel}, nil
				}
			}
			got, err := i.Feed.Labels(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}

}

func TestSaveLabel(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		label   string
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid:successfully_save_valid_label",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				label:   ksuid.New().String(),
			},
			want:    []string{common.DefaultLabel},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_save_valid_label" {
				fakeEngagement.SaveLabelFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					label string,
				) error {
					return nil
				}
			}
			err := i.Feed.SaveLabel(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestSaveLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestUnreadPersistentItems(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

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
			name: "valid:successfully_return_no_of_unread_inboxItems",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_return_no_of_unread_inboxItems" {
				fakeEngagement.UnreadPersistentItemsFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
				) (int, error) {
					return 0, nil
				}
			}
			_, err := i.Feed.UnreadPersistentItems(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestUnreadPersistentItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestUpdateUnreadPersistentItemsCount(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

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
			name: "valid:successfully_update_no_of_unread_inboxItems",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_update_no_of_unread_inboxItems" {
				fakeEngagement.UpdateUnreadPersistentItemsCountFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
				) error {
					return nil
				}
			}
			err := i.Feed.UpdateUnreadPersistentItemsCount(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeedUseCaseImpl.TestUpdateUnreadPersistentItemsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestPublishNudge(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()
	nudge := testNudge()
	type args struct {
		ctx     context.Context
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
			name: "valid:successfully_publish_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge:   nudge,
			},
			wantErr: false,
		},
		{
			name: "invalid:nil_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge:   nudge,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge:   nudge,
			},
			wantErr: true,
		},
		{
			name: "invalid:invalid_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge:   &feedlib.Nudge{},
			},
			wantErr: true,
		},
		{
			name: "invalid:invalid_nudgeActionType",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudge: &feedlib.Nudge{
					ID: uuid.New().String(),
					Actions: []feedlib.Action{
						{
							ActionType: feedlib.ActionTypeFloating,
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_publish_nudge" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						SequenceNumber: 0,
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:nil_nudge" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("can't publish nil nudge")
				}
			}

			if tt.name == "invalid:fail_to_save_nudge" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("unable to publish nudge")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						SequenceNumber: 0,
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:invalid_nudge" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("invalid nudge")
				}
			}
			if tt.name == "invalid:invalid_nudgeActionType" {
				fakeEngagement.SaveNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID: uuid.New().String(),
						Actions: []feedlib.Action{
							{
								ActionType: feedlib.ActionTypeFloating,
							},
						},
					}, fmt.Errorf("invalid nudge")
				}
			}
			got, err := i.Feed.PublishNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudge)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PublishNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil item response returned")
					return
				}
			}
		})
	}
}

func TestResolveNudge(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()
	nudge := testNudge()
	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "valid:successfully_resolve_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantStatus: feedlib.StatusDone,
			wantErr:    false,
		},
		{
			name: "invalid:fail_to_get_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_resolve_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.ResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("unable to get nudge")
				}
			}

			if tt.name == "invalid:fail_to_update_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.ResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("unable to resolve nudge")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.ResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("unable to send notification")
				}
			}
			got, err := i.Feed.ResolveNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.ResolveNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil nudge response returned")
					return
				}
			}
		})
	}
}

func TestUnresolveNudge(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()
	nudge := testNudge()
	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "valid:successfully_Unresolve_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantStatus: feedlib.StatusDone,
			wantErr:    false,
		},
		{
			name: "invalid:fail_to_get_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: feedlib.FlavourConsumer,
				nudgeID: nudge.ID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_Unresolve_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.UnResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("unable to get nudge")
				}
			}

			if tt.name == "invalid:fail_to_update_nudge" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.UnResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return nil, fmt.Errorf("unable to Unresolve nudge")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudgeID string,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID:             uuid.New().String(),
						Status:         feedlib.StatusDone,
						SequenceNumber: 1,
						Actions: []feedlib.Action{
							{
								Name: common.UnResolveItemActionName,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateNudgeFn = func(
					ctx context.Context,
					uid string,
					flavour feedlib.Flavour,
					nudge *feedlib.Nudge,
				) (*feedlib.Nudge, error) {
					return &feedlib.Nudge{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour feedlib.Flavour,
					payload feedlib.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("unable to send notification")
				}
			}
			got, err := i.Feed.UnresolveNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.UnresolveNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil nudge response returned")
					return
				}
			}
		})
	}
}
