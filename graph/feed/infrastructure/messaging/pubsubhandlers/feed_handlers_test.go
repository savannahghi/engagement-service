package pubsubhandlers_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging/pubsubhandlers"
)

// TODO Get test FCM token from env
// TODO Document the senders, they are very clear
// TODO Document item publish: item is notified + *tray notification*
// TODO Document item delete: item is notified, no tray notification
// TODO Document item resolve: item is notified, no tray notification
// TODO Document item unresolve: item is notified, no tray notification
// TODO Document item hide: item is notified, no tray notification
// TODO Document item show: item is notified, no tray notification
// TODO Document item pin: item is notified, no tray notification
// TODO Document item unpin: item is notified, no tray notification
// TODO Document nudge publish: nudge is notified
// TODO Document nudge delete: nudge is notified
// TODO Document nudge resolve: nudge is notified
// TODO Document nudge unresolve: nudge is notified
// TODO Document nudge hide: nudge is notified
// TODO Document show nudge: nudge is notified
// TODO Document action publish: thin feed is notified
// TODO Document action delete: thin feed is notified
// TODO Document message post: message notified
// TODO Document message delete: message notified

const (
	intMax = 9007199254740990
)

func getTestPubsubPayload(t *testing.T, el base.Element) *base.PubSubPayload {
	elData, err := el.ValidateAndMarshal()
	if err != nil {
		t.Errorf("invalid element: %w", err)
		return nil
	}

	envelope := feed.NotificationEnvelope{
		UID:     ksuid.New().String(),
		Flavour: base.FlavourConsumer,
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

func getTestItem() base.Item {
	return base.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       base.TextTypePlain,
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
			base.GetYoutubeVideoLink(base.SampleVideoURL, "title", "description", base.LogoURL),
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
			},
		},
		Conversations: []base.Message{
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
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func testNudge() *base.Nudge {
	return &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []base.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []base.Channel{
			base.ChannelEmail,
			base.ChannelFcm,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getTestEvent() base.Event {
	return base.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: base.Context{
			UserID:         ksuid.New().String(),
			Flavour:        base.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() base.Action {
	return base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		ActionType:     base.ActionTypePrimary,
		Handling:       base.HandlingFullPage,
	}
}

func getTestMessage() base.Message {
	return base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func TestNotifyItemUpdate(t *testing.T) {
	type args struct {
		ctx                 context.Context
		sender              string
		includeNotification bool
		m                   *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO Wrong notification data in payload
		// TODO Item can't be unmarshalled from envelope contents
		// TODO Valid persistent item
		// TODO Valid non-persistent item
		// TODO Default: unknown sender
		// TODO Item publish sender...new labels
		// TODO Item publish sender, existing/default label
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.NotifyItemUpdate(tt.args.ctx, tt.args.sender, tt.args.includeNotification, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("NotifyItemUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendNotificationViaFCM(t *testing.T) {
	type args struct {
		ctx          context.Context
		uids         []string
		sender       string
		pl           feed.NotificationEnvelope
		notification *base.FirebaseSimpleNotificationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO UIDs with no tokens...error
		// TODO UIDs with no tokens...blank
		// TODO UIDs with tokens...valid push
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.SendNotificationViaFCM(tt.args.ctx, tt.args.uids, tt.args.sender, tt.args.pl, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("SendNotificationViaFCM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetThinFeed(t *testing.T) {
	type args struct {
		ctx     context.Context
		uid     string
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *feed.Feed
		wantErr bool
	}{
		// TODO: Initialize new feed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pubsubhandlers.GetThinFeed(tt.args.ctx, tt.args.uid, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThinFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetThinFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserTokens(t *testing.T) {
	type args struct {
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO UIDs with tokens
		// TODO UIDs with no tokens
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pubsubhandlers.GetUserTokens(tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifyInboxCountUpdate(t *testing.T) {
	type args struct {
		ctx     context.Context
		uid     string
		flavour base.Flavour
		count   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO Successful inbox update notification
		// TODO Unsuccessful inbox update notification
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.NotifyInboxCountUpdate(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("NotifyInboxCountUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotifyNudgeUpdate(t *testing.T) {
	type args struct {
		ctx    context.Context
		sender string
		m      *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO Successful nudge update notification
		// TODO Unsuccessful nudge update notification
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.NotifyNudgeUpdate(tt.args.ctx, tt.args.sender, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("NotifyNudgeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateInbox(t *testing.T) {
	type args struct {
		ctx     context.Context
		uid     string
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO Successful inbox update
		// TODO Failed inbox update
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pubsubhandlers.UpdateInbox(tt.args.ctx, tt.args.uid, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInbox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
