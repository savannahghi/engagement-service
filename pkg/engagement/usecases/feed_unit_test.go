package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
)

func TestPublishFeedItem(t *testing.T) {
	ctx := context.Background()

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
		flavour base.Flavour
		item    *base.Item
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
				flavour: base.FlavourConsumer,
				item:    testItem,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_save_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				item:    testItem,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_a_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				item:    testItem,
			},
			wantErr: true,
		},
		{
			name: "invalid:use_a_nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				item:    nil,
			},
			wantErr: true,
		},
		{
			name: "invalid:use_an_invalid_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				item:    &base.Item{},
			},
			wantErr: true,
		},
		{
			name: "invalid:use_an_invalid_action_type",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_save_feed_item" {
				fakeEngagement.SaveFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("unable to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_a_notification" {
				fakeEngagement.SaveFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send a notification")
				}

				if tt.name == "invalid:use_a_nil_item" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour base.Flavour,
						item *base.Item,
					) (*base.Item, error) {
						return nil, fmt.Errorf("can't publish nil feed item")
					}
				}

				if tt.name == "invalid:use_an_invalid_item" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour base.Flavour,
						item *base.Item,
					) (*base.Item, error) {
						return nil, fmt.Errorf("unable to publish feed item")
					}
				}

				if tt.name == "invalid:use_an_invalid_action_type" {
					fakeEngagement.SaveFeedItemFn = func(
						ctx context.Context,
						uid string,
						flavour base.Flavour,
						item *base.Item,
					) (*base.Item, error) {
						return &base.Item{
							ID: uuid.New().String(),
							Actions: []base.Action{
								{
									ID:         ksuid.New().String(),
									Name:       "TEST_ACTION",
									ActionType: base.ActionTypeFloating,
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
	ctx := context.Background()

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
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  "",
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_delete_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) error {
					return nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, fmt.Errorf("unable to retrieve feed item")
				}
			}

			if tt.name == "invalid:fail_to_delete_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) error {
					return fmt.Errorf("failed to delete feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.DeleteFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) error {
					return nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
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
	ctx := context.Background()

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
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ResolveItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
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
	ctx := context.Background()

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
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
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
	ctx := context.Background()

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
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.UnPinItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to publish feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.PinItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
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
	ctx := context.Background()

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
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:nil_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				// itemID:  "",
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.UnResolveItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "invalid:nil_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
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
	ctx := context.Background()

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
		flavour base.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "valid:successfully_hide_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantVisibility: base.VisibilityHide,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.HideItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
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
	ctx := context.Background()

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
		flavour base.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "valid:successfully_show_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantVisibility: base.VisibilityShow,
			wantErr:        false,
		},
		{
			name: "invalid:fail_to_get_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_feed_item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
				itemID:  testItem.ID,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_send_notification",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
					metadata map[string]interface{},
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to get feed item")
				}
			}

			if tt.name == "invalid:fail_to_update_feed_item" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return nil, fmt.Errorf("failed to update feed item")
				}
			}

			if tt.name == "invalid:fail_to_send_notification" {
				fakeEngagement.GetFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					itemID string,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           common.ShowItemActionName,
								Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
						},
					}, nil
				}

				fakeEngagement.UpdateFeedItemFn = func(
					ctx context.Context,
					uid string,
					flavour base.Flavour,
					item *base.Item,
				) (*base.Item, error) {
					return &base.Item{
						ID: uuid.New().String(),
					}, nil
				}

				fakeMessaging.NotifyFn = func(
					ctx context.Context,
					topicID string,
					uid string,
					flavour base.Flavour,
					payload base.Element,
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
	ctx := context.Background()

	i, err := InitializeFakeEngagementInteractor()
	if err != nil {
		t.Errorf("failed to initialize the fake engagement interactor: %v", err)
		return
	}
	uid := ksuid.New().String()

	type args struct {
		ctx     context.Context
		uid     string
		flavour base.Flavour
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
				flavour: base.FlavourConsumer,
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
					flavour base.Flavour,
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
