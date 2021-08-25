package helpers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const ActionPublishTopic = "actions.publish"

var tracer = otel.Tracer("example.com")

func TestAddPubSubNamespace(t *testing.T) {
	type args struct {
		topicName string
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
	}{
		{
			name: "happy case: valid action publish topic provided",
			args: args{
				topicName: ActionPublishTopic,
			},
			wantValue: true,
		},
		{
			name:      "sad case: empty action publish topic provided",
			args:      args{},
			wantValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.AddPubSubNamespace(tt.args.topicName)
			if tt.wantValue {
				assert.NotNil(t, got)
				assert.NotEmpty(t, got)
				assert.Contains(t, got, tt.args.topicName)
			}
			if !tt.wantValue {
				assert.Contains(t, got, "--")
			}
		})
	}
}

func TestValidateElement(t *testing.T) {
	type args struct {
		el feedlib.Element
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid feed lib element provided",
			args: args{
				el: &feedlib.Message{
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
			name:    "sad case: nil element provided",
			args:    args{},
			wantErr: true,
		},
		{
			name: "sad case: missing field in element",
			args: args{
				el: &feedlib.Message{
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
			if err := helpers.ValidateElement(tt.args.el); (err != nil) != tt.wantErr {
				t.Errorf("ValidateElement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecordSpanError(t *testing.T) {
	ctx := context.Background()
	_, span := tracer.Start(ctx, "SetDefaultActions")
	type args struct {
		span trace.Span
		err  error
	}
	tests := []struct {
		name   string
		args   args
		panics bool
	}{
		{
			name: "happy case: valid error and span params provided",
			args: args{
				span: span,
				err:  fmt.Errorf("test error"),
			},
			panics: false,
		},
		{
			name: "sad case: no span params",
			args: args{
				err: fmt.Errorf("test error"),
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcRecordSpanError := func() { helpers.RecordSpanError(tt.args.span, tt.args.err) }
			if tt.panics {
				assert.Panics(t, fcRecordSpanError)
			}
		})
	}
}

func TestEpochTimetoStandardTime(t *testing.T) {
	testTime := time.Now()
	testTimeString := testTime.String()

	type args struct {
		timeOfDelivery string
	}

	tests := []struct {
		name   string
		args   args
		panics bool
	}{
		{
			name: "happy case: valid time of delivery provided",
			args: args{
				timeOfDelivery: testTimeString,
			},
			panics: false,
		},
		{
			name: "sad case: invalid time of delivery provided",
			args: args{
				timeOfDelivery: "",
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcEpochTimetoStandardTime := func() { helpers.EpochTimetoStandardTime(tt.args.timeOfDelivery) }
			if tt.panics {
				assert.Panics(t, fcEpochTimetoStandardTime)
			}
			if !tt.panics {
				got := helpers.EpochTimetoStandardTime(tt.args.timeOfDelivery)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestStripMailGunIDSpecialCharacters(t *testing.T) {
	testMessageID := "12>dd<gg"

	type args struct {
		messageID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string with symbol characters provided",
			args: args{
				messageID: testMessageID,
			},
			want: "12ddgg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.StripMailGunIDSpecialCharacters(tt.args.messageID)
			if got != tt.want {
				t.Errorf("StripMailGunIDSpecialCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}
