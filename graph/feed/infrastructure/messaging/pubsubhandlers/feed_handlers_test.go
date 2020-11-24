package pubsubhandlers_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging/pubsubhandlers"
)

const (
	intMax = 9007199254740990
)

func getTestPubsubPayload(t *testing.T, el feed.Element) *base.PubSubPayload {
	elData, err := el.ValidateAndMarshal()
	if err != nil {
		t.Errorf("invalid element: %w", err)
		return nil
	}

	envelope := feed.NotificationEnvelope{
		UID:     ksuid.New().String(),
		Flavour: feed.FlavourConsumer,
		Payload: elData,
		Metadata: map[string]interface{}{
			ksuid.New().String(): ksuid.New().String(),
		},
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		t.Errorf("can't marshal envelope data: %w", err)
		return nil
	}

	return &base.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: base.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      data,
			Attributes: map[string]string{
				"topicID": ksuid.New().String(),
			},
		},
	}
}

func TestHandleFeedRetrieval(t *testing.T) {
	thinFeed := getTestThinFeed(t)
	if thinFeed == nil {
		return
	}

	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, thinFeed),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleFeedRetrieval(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleFeedRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleThinFeedRetrieval(t *testing.T) {
	thinFeed := getTestThinFeed(t)
	if thinFeed == nil {
		return
	}

	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, thinFeed),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleThinFeedRetrieval(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleThinFeedRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemRetrieval(t *testing.T) {
	item := getTestItem()

	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemRetrieval(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPublish(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemDelete(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemResolve(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnresolve(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemHide(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemShow(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPin(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemPin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnpin(t *testing.T) {
	item := getTestItem()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleItemUnpin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnpin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeRetrieval(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeRetrieval(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgePublish(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgePublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgePublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeDelete(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeResolve(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeUnresolve(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeHide(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeShow(t *testing.T) {
	nudge := testNudge()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleNudgeShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionRetrieval(t *testing.T) {
	action := getTestAction()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &action),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionRetrieval(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionPublish(t *testing.T) {
	action := getTestAction()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &action),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionDelete(t *testing.T) {
	action := getTestAction()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &action),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleActionDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessagePost(t *testing.T) {
	message := getTestMessage()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &message),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleMessagePost(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessagePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessageDelete(t *testing.T) {
	message := getTestMessage()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &message),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleMessageDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleIncomingEvent(t *testing.T) {
	event := getTestEvent()
	ctx := context.Background()
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &event),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.HandleIncomingEvent(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleIncomingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getTestThinFeed(t *testing.T) *feed.Feed {
	ctx := context.Background()
	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	repository, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate new firebase repository: %v", err)
		return nil
	}

	notificationService, err := messaging.NewMockNotificationService()
	if err != nil {
		t.Errorf("can't instantiate notification service")
		return nil
	}

	agg, err := feed.NewCollection(repository, notificationService)
	if err != nil {
		t.Errorf("can't instantiate feed collection: %w", err)
		return nil
	}
	thinFeed, err := agg.GetThinFeed(ctx, uid, flavour)
	if err != nil {
		t.Errorf("can't instantiate thin feed: %w", err)
		return nil
	}
	return thinFeed
}

func getTestItem() feed.Item {
	return feed.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feed.TextTypePlain,
		Links: []feed.Link{
			feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
			feed.GetYoutubeVideoLink(feed.SampleVideoURL, "title", "description", feed.LogoURL),
		},
		Actions: []feed.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
				ActionType:     feed.ActionTypeSecondary,
				Handling:       feed.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
			},
		},
		Conversations: []feed.Message{
			{
				ID:             "msg-2",
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
				SequenceNumber: int(time.Now().Unix()),
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

func testNudge() *feed.Nudge {
	return &feed.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feed.Link{
			feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
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

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
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
		Icon:           feed.GetPNGImageLink(feed.LogoURL, "title", "description", feed.LogoURL),
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
	}
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
