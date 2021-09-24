package usecases_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleItemPublish(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Library: %v", err)
	}

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
			name: "Sad Case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy Case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemDelete(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.HandleItemDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemResolve(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.HandleItemResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnresolve(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sd case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemHide(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemShow(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemPin(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemPin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleItemUnpin(t *testing.T) {
	item := getTheTestItem(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &item),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleItemUnpin(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleItemUnpin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgePublish(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgePublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgePublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeDelete(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgeDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeResolve(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgeResolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeResolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeUnresolve(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgeUnresolve(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeUnresolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeHide(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgeHide(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeHide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleNudgeShow(t *testing.T) {
	nudge := aTestNudge(t)
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, nudge),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleNudgeShow(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleNudgeShow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionPublish(t *testing.T) {
	action := getATestAction()
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &action),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleActionPublish(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActionDelete(t *testing.T) {
	action := getATestAction()
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &action),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleActionDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleActionDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessagePost(t *testing.T) {
	message := getATestMessage()
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &message),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleMessagePost(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessagePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleMessageDelete(t *testing.T) {
	message := getATestMessage()
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &message),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleMessageDelete(ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HandleMessageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleIncomingEvent(t *testing.T) {
	event := getATestEvent()
	ctx := firebasetools.GetAuthenticatedContext(t)
	n, _, err := InitializeTestNewNotification(ctx)
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
			name: "Sad case: nil payload",
			args: args{
				m: nil,
			},
			wantErr: true,
		},
		{
			name: "Happy case: non nil payload",
			args: args{
				m: getTestPubsubPayload(t, &event),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := n.LibUsecases.HandleIncomingEvent(ctx, tt.args.m); (err != nil) != tt.wantErr {
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
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.NotifyItemUpdate(tt.args.ctx, tt.args.sender, tt.args.includeNotification, tt.args.m); (err != nil) != tt.wantErr {
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
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.SendNotificationViaFCM(tt.args.ctx, tt.args.uids, tt.args.sender, tt.args.pl, tt.args.notification); (err != nil) != tt.wantErr {
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

	n, _, err := InitializeTestNewNotification(ctx)
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
			tokens, err := n.LibUsecases.GetUserTokens(tt.args.ctx, tt.args.uids)
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
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.NotifyInboxCountUpdate(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.count); (err != nil) != tt.wantErr {
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
	n, _, err := InitializeTestNewNotification(ctx)
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
			err := n.LibUsecases.NotifyNudgeUpdate(tt.args.ctx, tt.args.sender, tt.args.m)
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
	n, _, err := InitializeTestNewNotification(ctx)
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
			if err := n.LibUsecases.UpdateInbox(tt.args.ctx, tt.args.uid, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInbox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotificationImpl_SendNotificationEmail(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	n, _, err := InitializeTestNewNotification(ctx)
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
			err := n.LibUsecases.SendNotificationEmail(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationImpl.SendNotificationEmail() error = %v, wantErr %v", err, tt.wantErr)
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
