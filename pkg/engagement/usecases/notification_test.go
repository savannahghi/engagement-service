package usecases_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/usecases"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
)

const (
	onboardingService = "profile"
	intMax            = 9007199254740990
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

func InitializeTestNewNotification(ctx context.Context) (*usecases.NotificationImpl, error) {
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, err
	}
	fcmNotification, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate push notification service : %w", err)
	}
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	notification := usecases.NewNotification(fr, fcmNotification, onboarding)
	return notification, nil
}

func onboardingISCClient(t *testing.T) *base.InterServiceClient {
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil
	}

	return profileClient
}

func getATestMessage() base.Message {
	return base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTheTestItem(t *testing.T) base.Item {
	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return base.Item{}
	}

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
			token.UID,
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

func getATestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getATestEvent() base.Event {
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

func getATestAction() base.Action {
	return base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		ActionType:     base.ActionTypePrimary,
		Handling:       base.HandlingFullPage,
	}
}

func getTestPubsubPayload(t *testing.T, el base.Element) *base.PubSubPayload {
	elData, err := el.ValidateAndMarshal()
	if err != nil {
		t.Errorf("invalid element: %w", err)
		return nil
	}

	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}

	envelope := resources.NotificationEnvelope{
		UID:     token.UID,
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

func aTestNudge(t *testing.T) *base.Nudge {
	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}
	return &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []base.Action{
			getATestAction(),
		},
		Users: []string{
			token.UID,
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

func TestHandleItemPublish(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemDelete(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemResolve(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnresolve(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemHide(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemShow(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPin(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemPin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnpin(t *testing.T) {
	item := getTheTestItem(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleItemUnpin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnpin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgePublish(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgePublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgePublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeDelete(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgeDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeResolve(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgeResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeUnresolve(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgeUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeHide(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgeHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeShow(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleNudgeShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionPublish(t *testing.T) {
	action := getATestAction()
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleActionPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionDelete(t *testing.T) {
	action := getATestAction()
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleActionDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessagePost(t *testing.T) {
	message := getATestMessage()
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleMessagePost(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessagePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessageDelete(t *testing.T) {
	message := getATestMessage()
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleMessageDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleIncomingEvent(t *testing.T) {
	event := getATestEvent()
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.HandleIncomingEvent(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleIncomingEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotifyItemUpdate(t *testing.T) {
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.NotifyItemUpdate(tt.args.ctx, tt.args.sender, tt.args.includeNotification, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("NotifyItemUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendNotificationViaFCM(t *testing.T) {
	ctx, _, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		ctx          context.Context
		uids         []string
		sender       string
		pl           resources.NotificationEnvelope
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
			if err := notify.SendNotificationViaFCM(tt.args.ctx, tt.args.uids, tt.args.sender, tt.args.pl, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("SendNotificationViaFCM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserTokens(t *testing.T) {
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			got, err := notify.GetUserTokens(tt.args.uids)
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
	ctx, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
		{
			name: "Successful: inbox update notification",
			args: args{
				ctx:     ctx,
				uid:     token.UID,
				flavour: base.FlavourConsumer,
				count:   10,
			},
			wantErr: false,
		},
		{
			name: "Invalid: fail to send a notification message",
			args: args{
				ctx:     context.Background(),
				uid:     "invalid uid",
				flavour: base.FlavourConsumer,
				count:   10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := notify.NotifyInboxCountUpdate(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("NotifyInboxCountUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotifyNudgeUpdate(t *testing.T) {
	ctx := context.Background()
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
			if err := notify.NotifyNudgeUpdate(tt.args.ctx, tt.args.sender, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("NotifyNudgeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateInbox(t *testing.T) {
	ctx, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
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
		{
			name: "Successful inbox update",
			args: args{
				ctx:     ctx,
				uid:     token.UID,
				flavour: base.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Failed inbox update on Consumer",
			args: args{
				ctx:     ctx,
				uid:     "invalid uid",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Failed inbox update on Pro",
			args: args{
				ctx:     ctx,
				uid:     "invalid uid",
				flavour: base.FlavourPro,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := notify.UpdateInbox(tt.args.ctx, tt.args.uid, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInbox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
