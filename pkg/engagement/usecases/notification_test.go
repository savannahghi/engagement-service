package usecases_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/usecases"
)

const (
	onboardingService = "profile"
	intMax            = 9007199254740990
	registerPushToken = "testing/register_push_token"
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
	fcm := fcm.NewService(fr)
	mail := mail.NewService(fr)
	crm := hubspot.NewHubSpotService()
	notification := usecases.NewNotification(fr, fcmNotification, onboarding, fcm, mail, crm)
	return notification, nil
}

func onboardingISCClient(t *testing.T) *interserviceclient.InterServiceClient {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil
	}

	return profileClient
}

func RegisterPushToken(
	ctx context.Context,
	t *testing.T,
	UID string,
	onboardingClient *interserviceclient.InterServiceClient,
) (bool, error) {
	token := "random"
	if onboardingClient == nil {
		return false, fmt.Errorf("nil ISC client")
	}

	payload := map[string]interface{}{
		"pushTokens": token,
		"uid":        UID,
	}
	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		registerPushToken,
		payload,
	)
	if err != nil {
		return false, fmt.Errorf("unable to make a request to register push token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("expected a StatusOK (200) status code but instead got %v", resp.StatusCode)
	}

	return true, nil
}

func getATestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTheTestItem(t *testing.T) feedlib.Item {
	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return feedlib.Item{}
	}

	return feedlib.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feedlib.TextTypePlain,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
			feedlib.GetYoutubeVideoLink(feedlib.SampleVideoURL, "title", "description", feedlib.LogoURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
			},
		},
		Conversations: []feedlib.Message{
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
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelFcm,
			feedlib.ChannelEmail,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func getATestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getATestEvent() feedlib.Event {
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

func getATestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		ActionType:     feedlib.ActionTypePrimary,
		Handling:       feedlib.HandlingFullPage,
	}
}

func getTestPubsubPayload(t *testing.T, el feedlib.Element) *pubsubtools.PubSubPayload {
	elData, err := el.ValidateAndMarshal()
	if err != nil {
		t.Errorf("invalid element: %w", err)
		return nil
	}

	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}

	envelope := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
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

	return &pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      data,
			Attributes: map[string]string{
				"topicID": ksuid.New().String(),
			},
		},
	}
}

func aTestNudge(t *testing.T) *feedlib.Nudge {
	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			getATestAction(),
		},
		Users: []string{
			token.UID,
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

func TestHandleItemPublish(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		m *pubsubtools.PubSubPayload
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
	item := getTheTestItem(t)
	invalidItem := getTheTestItem(t)
	invalidItem.Label = "new-label"
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))

	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		ctx                 context.Context
		sender              string
		includeNotification bool
		m                   *pubsubtools.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:use a valid persistent item",
			args: args{
				ctx:                 ctx,
				sender:              "ITEM_PUBLISHED",
				includeNotification: true,
				m:                   getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
		{
			name: "invalid:unknown sender",
			args: args{
				ctx:                 ctx,
				sender:              "unknown sender",
				includeNotification: true,
				m:                   getTestPubsubPayload(t, &item),
			},
			wantErr: true,
		},
		{
			name: "Invalid: Item can't be unmarshalled from envelope contents",
			args: args{
				ctx:                 ctx,
				sender:              "ITEM_PUBLISHED",
				includeNotification: true,
				m: &pubsubtools.PubSubPayload{
					Subscription: ksuid.New().String(),
					Message: pubsubtools.PubSubMessage{
						MessageID: ksuid.New().String(),
						Data:      nil,
						Attributes: map[string]string{
							"topicID": common.ActionPublishTopic,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid: Wrong notification data in payload",
			args: args{
				ctx:                 ctx,
				sender:              "ITEM_PUBLISHED",
				includeNotification: true,
				m: &pubsubtools.PubSubPayload{
					Subscription: ksuid.New().String(),
					Message: pubsubtools.PubSubMessage{
						MessageID: "invalid id",
						Data:      []byte("data"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid: use a new label",
			args: args{
				ctx:                 ctx,
				sender:              "ITEM_PUBLISHED",
				includeNotification: true,
				m:                   getTestPubsubPayload(t, &invalidItem),
			},
			wantErr: false,
		},
		// TODO Valid non-persistent item
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
	_ = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingISCClient(t))
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))

	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}
	imageurl := "some-image-url"
	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		ctx          context.Context
		uids         []string
		sender       string
		pl           dto.NotificationEnvelope
		notification *firebasetools.FirebaseSimpleNotificationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: UIDs with tokens",
			args: args{
				ctx:    ctx,
				uids:   []string{token.UID},
				sender: "Be.Well",
				pl: dto.NotificationEnvelope{
					UID:     token.UID,
					Flavour: feedlib.FlavourConsumer,
					Payload: []byte("payload"),
					Metadata: map[string]interface{}{
						"language": "en",
					},
				},
				notification: &firebasetools.FirebaseSimpleNotificationInput{
					Title:    "Scheduled visit",
					Body:     "Your visit has been scheduled for Thursday morning",
					ImageURL: &imageurl,
					Data: map[string]interface{}{
						"patient": "outpatient",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: UIDs with blank tokens",
			args: args{
				ctx:    ctx,
				uids:   []string{""},
				sender: "Be.Well",
				pl: dto.NotificationEnvelope{
					UID:     "",
					Flavour: feedlib.FlavourConsumer,
					Payload: []byte("payload"),
					Metadata: map[string]interface{}{
						"language": "en",
					},
				},
				notification: &firebasetools.FirebaseSimpleNotificationInput{
					Title:    "Scheduled visit",
					Body:     "Your visit has been scheduled for Thursday morning",
					ImageURL: &imageurl,
					Data: map[string]interface{}{
						"patient": "outpatient",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: UIDs with no tokens",
			args: args{
				ctx:    context.Background(),
				uids:   []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
				sender: "Be.Well",
				pl: dto.NotificationEnvelope{
					UID:     "",
					Flavour: feedlib.FlavourConsumer,
					Payload: []byte("payload"),
					Metadata: map[string]interface{}{
						"language": "en",
					},
				},
				notification: &firebasetools.FirebaseSimpleNotificationInput{
					Title:    "Scheduled visit",
					Body:     "Your visit has been scheduled for Thursday morning",
					ImageURL: &imageurl,
					Data: map[string]interface{}{
						"patient": "outpatient",
					},
				},
			},
			wantErr: true,
		},
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
	// clean up
	_ = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingISCClient(t))
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))

	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}

	notify, err := InitializeTestNewNotification(ctx)
	assert.Nil(t, err)
	type args struct {
		ctx  context.Context
		uids []string
	}
	tokens := token.UID
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get user push tokens",
			args: args{
				ctx: ctx,
				uids: []string{
					tokens,
				},
			},
			wantErr: false,
		},

		{
			name: "sad case: get user push tokens",
			args: args{
				ctx: ctx,
				uids: []string{
					"invalid_uid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := notify.GetUserTokens(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(tokens) == 0 {
					t.Errorf("expected user tokens")
					return
				}
			}
			if tt.wantErr {
				if len(tokens) != 0 {
					t.Errorf("did not expected user tokens")
					return
				}
			}
		})
	}
}

func TestNotifyInboxCountUpdate(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
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
		flavour feedlib.Flavour
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
				flavour: feedlib.FlavourConsumer,
				count:   10,
			},
			wantErr: false,
		},
		{
			name: "Invalid: fail to send a notification message",
			args: args{
				ctx:     context.Background(),
				uid:     "invalid uid",
				flavour: feedlib.FlavourConsumer,
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
	nudge := aTestNudge(t)
	action := getATestAction()
	ctx, _, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
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
		ctx    context.Context
		sender string
		m      *pubsubtools.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful nudge update notification",
			args: args{
				ctx:    ctx,
				sender: "NUDGE_PUBLISHED",
				m:      getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
		{
			name: "Fail to notify nudge update - Use invalid payload",
			args: args{
				ctx:    ctx,
				sender: "Be Well",
				m:      getTestPubsubPayload(t, &action),
			},
			wantErr: true,
		},
		{
			name: "Fail to notify nudge update - unknown sender",
			args: args{
				ctx:    ctx,
				sender: "Be Well",
				m:      getTestPubsubPayload(t, nudge),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notify.NotifyNudgeUpdate(tt.args.ctx, tt.args.sender, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotifyNudgeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateInbox(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
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
		flavour feedlib.Flavour
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
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Failed inbox update on Consumer",
			args: args{
				ctx:     ctx,
				uid:     "invalid uid",
				flavour: feedlib.FlavourConsumer,
			},
			// TODO: Restore after the milestone @mathenge
			wantErr: false,
		},
		{
			name: "Failed inbox update on Pro",
			args: args{
				ctx:     ctx,
				uid:     "invalid uid",
				flavour: feedlib.FlavourPro,
			},
			// TODO: Restore after the milestone @mathenge
			wantErr: false,
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

func TestNotificationImpl_SendEmail(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	notify, err := InitializeTestNewNotification(ctx)
	if err != nil {
		t.Errorf("failed to initialize new service")
		return
	}

	body := map[string]interface{}{
		"to":      []string{"automated.test.user.bewell-app-ci@healthcloud.co.ke"},
		"text":    "This is a test message",
		"subject": "Test Subject",
	}

	payloadData, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal data: %v", err)
		return
	}

	pubSubMessage := pubsubtools.PubSubMessage{
		Data: payloadData,
		Attributes: map[string]string{
			"topicID": "mails.send",
		},
	}

	payload := &pubsubtools.PubSubPayload{
		Message: pubSubMessage,
	}
	type args struct {
		ctx context.Context
		m   *pubsubtools.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Send email to a valid email",
			args: args{
				ctx: ctx,
				m:   payload,
			},
			wantErr: false,
		},
		{
			name: "invalid case - missing payload",
			args: args{
				ctx: context.Background(),
				m:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notify.SendEmail(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationImpl.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got: %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got: %v", err)
					return
				}
			}
		})
	}
}
