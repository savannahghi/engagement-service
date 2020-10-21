package messaging_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/messaging"
)

func TestNewPubSubNotificationService(t *testing.T) {
	ctx := context.Background()
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := messaging.NewPubSubNotificationService(ctx, projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewPubSubNotificationService() error = %v, wantErr %v",
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

func TestPubSubNotificationService_Notify(t *testing.T) {
	ctx := context.Background()
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	type args struct {
		channel string
		el      feed.Element
	}
	tests := []struct {
		name    string
		pubsub  feed.NotificationService
		args    args
		wantErr bool
	}{
		{
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &feed.Message{
					ID:           uuid.New().String(),
					Text:         uuid.New().String(),
					ReplyTo:      uuid.New().String(),
					PostedByUID:  uuid.New().String(),
					PostedByName: uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name:   "invalid message, missing posted by info",
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &feed.Message{
					ID:      uuid.New().String(),
					Text:    uuid.New().String(),
					ReplyTo: uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pubsub.Notify(
				context.Background(),
				tt.args.channel,
				tt.args.el,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"PubSubNotificationService.Notify() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
