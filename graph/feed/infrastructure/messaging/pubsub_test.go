package messaging_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/messaging"
)

func getProjectNumber(t *testing.T) int {
	projectNumber, err := base.GetEnvVar(base.GoogleProjectNumberEnvVarName)
	if err != nil {
		t.Errorf("project number not found in env var: %s", err)
		return 0
	}

	if projectNumber == "" {
		t.Errorf("nil project number")
		return 0
	}

	projectNumberInt, err := strconv.Atoi(projectNumber)
	if err != nil {
		t.Errorf("non int project number: %s", err)
		return 0
	}

	return projectNumberInt
}

func TestNewPubSubNotificationService(t *testing.T) {
	ctx := context.Background()
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	projectNumber := getProjectNumber(t)
	if projectNumber == 0 {
		t.Errorf("the project number cannot be zero")
		return
	}

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
			got, err := messaging.NewPubSubNotificationService(ctx, projectID, projectNumber)
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

	projectNumber, err := base.GetEnvVar(base.GoogleProjectNumberEnvVarName)
	if err != nil {
		t.Errorf("project number not found in env var: %s", err)
		return
	}

	if projectNumber == "" {
		t.Errorf("nil project number")
		return
	}

	projectNumberInt, err := strconv.Atoi(projectNumber)
	if err != nil {
		t.Errorf("non int project number: %s", err)
		return
	}

	if projectNumberInt == 0 {
		t.Errorf("the project number cannot be zero")
		return
	}

	srv, err := messaging.NewPubSubNotificationService(ctx, projectID, projectNumberInt)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}

	if srv == nil {
		t.Errorf("nil pubsub notification service")
		return
	}

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
					ID:             ksuid.New().String(),
					SequenceNumber: 1,
					Text:           ksuid.New().String(),
					ReplyTo:        ksuid.New().String(),
					PostedByUID:    ksuid.New().String(),
					PostedByName:   ksuid.New().String(),
					Timestamp:      time.Now(),
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
					ID:        ksuid.New().String(),
					Text:      ksuid.New().String(),
					ReplyTo:   ksuid.New().String(),
					Timestamp: time.Now(),
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

func TestPubSubNotificationService_TopicIDs(t *testing.T) {
	ctx := context.Background()
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	projectNumber := getProjectNumber(t)
	if projectNumber == 0 {
		t.Errorf("the project number cannot be zero")
		return
	}
	ps, err := messaging.NewPubSubNotificationService(
		ctx, projectID, projectNumber)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service")
		return
	}

	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ps.TopicIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.TopicIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSubNotificationService_SubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	projectNumber := getProjectNumber(t)
	if projectNumber == 0 {
		t.Errorf("the project number cannot be zero")
		return
	}
	ps, err := messaging.NewPubSubNotificationService(
		ctx, projectID, projectNumber)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service")
		return
	}

	tests := []struct {
		name string
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ps.SubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.SubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
