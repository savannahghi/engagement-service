package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/messaging"
)

func TestNewPubSubNotificationService(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)

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
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}

	if srv == nil {
		t.Errorf("nil pubsub notification service")
		return
	}

	type args struct {
		channel  string
		uid      string
		flavour  base.Flavour
		el       base.Element
		metadata map[string]interface{}
	}
	tests := []struct {
		name    string
		pubsub  messaging.NotificationService
		args    args
		wantErr bool
	}{
		{
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &base.Message{
					ID:             ksuid.New().String(),
					SequenceNumber: 1,
					Text:           ksuid.New().String(),
					ReplyTo:        ksuid.New().String(),
					PostedByUID:    ksuid.New().String(),
					PostedByName:   ksuid.New().String(),
					Timestamp:      time.Now(),
				},
				uid:      ksuid.New().String(),
				flavour:  base.FlavourConsumer,
				metadata: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:   "invalid message, missing posted by info",
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &base.Message{
					ID:        ksuid.New().String(),
					Text:      ksuid.New().String(),
					ReplyTo:   ksuid.New().String(),
					Timestamp: time.Now(),
				},
				uid:      ksuid.New().String(),
				flavour:  base.FlavourPro,
				metadata: map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pubsub.Notify(
				context.Background(),
				tt.args.channel,
				tt.args.uid,
				tt.args.flavour,
				tt.args.el,
				tt.args.metadata,
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
