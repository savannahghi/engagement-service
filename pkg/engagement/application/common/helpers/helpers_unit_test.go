package helpers_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	mock "github.com/savannahghi/engagement/pkg/engagement/application/common/helpers/mock"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

var fakeHelpers mock.FakeHelpers

func Test_AddPubSubNamespace(t *testing.T) {
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
			if tt.name == "happy case: valid action publish topic provided" {
				fakeHelpers.AddPubSubNameSpaceFn = func(topicName string) string {
					return "service-namespace-env-version"
				}
				got := fakeHelpers.AddPubSubNameSpaceFn(tt.args.topicName)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got)
				assert.Contains(t, got, "namespace")
			}
			if tt.name == "sad case: empty action publish topic provided" {
				fakeHelpers.AddPubSubNameSpaceFn = func(topicName string) string {
					return "service--env-version"
				}
				got := fakeHelpers.AddPubSubNameSpaceFn(tt.args.topicName)
				assert.Contains(t, got, "--")
			}
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

func Test_ValidateElement(t *testing.T) {
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
			if tt.name == "happy case: valid feed lib element provided" {
				fakeHelpers.ValidateElementFn = func(el feedlib.Element) error {
					return nil
				}
				err := helpers.ValidateElement(tt.args.el)
				assert.Nil(t, err)
			}
			if tt.name == "sad case: nil element provided" {
				fakeHelpers.ValidateElementFn = func(el feedlib.Element) error {
					return fmt.Errorf("test error")
				}
				err := helpers.ValidateElement(tt.args.el)
				assert.NotNil(t, err)
			}
			if tt.name == "sad case: missing field in element" {
				fakeHelpers.ValidateElementFn = func(el feedlib.Element) error {
					return fmt.Errorf("test error")
				}
				err := helpers.ValidateElement(tt.args.el)
				assert.NotNil(t, err)
			}
			err := helpers.ValidateElement(tt.args.el)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateElement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RecordSpanError(t *testing.T) {
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
			if tt.name == "happy case: valid error and span params provided" {
				fakeHelpers.RecordSpanErrorFn = func(span trace.Span, err error) {}
				fcRecordSpanError := func() {
					fakeHelpers.RecordSpanErrorFn(tt.args.span, tt.args.err)
				}
				assert.NotPanics(t, fcRecordSpanError)

			}
			if tt.name == "sad case: no span params" {
				fakeHelpers.RecordSpanErrorFn = func(span trace.Span, err error) {
					log.Panic("paniced")
				}
				fcRecordSpanError := func() {
					fakeHelpers.RecordSpanErrorFn(tt.args.span, tt.args.err)
				}
				assert.Panics(t, fcRecordSpanError)

			}
			fcRecordSpanError := func() { helpers.RecordSpanError(tt.args.span, tt.args.err) }
			if tt.panics {
				assert.Panics(t, fcRecordSpanError)
			}
			if !tt.panics {
				assert.NotPanics(t, fcRecordSpanError)
			}
		})
	}
}

func Test_EpochTimetoStandardTime(t *testing.T) {
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
			if tt.name == "happy case: valid time of delivery provided" {
				fakeHelpers.EpochTimeToStandardTimeFn = func(timeOfDelivery string) time.Time {
					return testTime
				}
				got := fakeHelpers.EpochTimeToStandardTimeFn(tt.args.timeOfDelivery)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got)
			}
			if tt.name == "sad case: invalid time of delivery provided" {
				fakeHelpers.EpochTimeToStandardTimeFn = func(timeOfDelivery string) time.Time {
					log.Panic("paniced")
					return time.Now()
				}
				fcEpochTimetoStandardTime := func() { fakeHelpers.EpochTimeToStandardTimeFn(tt.args.timeOfDelivery) }
				assert.Panics(t, fcEpochTimetoStandardTime)
			}

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

func Test_StripMailGunIDSpecialCharacters(t *testing.T) {
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
			name: "happy case: string with symbol characters provided",
			args: args{
				messageID: testMessageID,
			},
			want: "12ddgg",
		},
	}
	for _, tt := range tests {
		if tt.name == "happy case: string with symbol characters provided" {
			fakeHelpers.StripMailGunIDSpecialCharactersFn = func(messageID string) string {
				return "12ddgg"
			}
			if got := fakeHelpers.StripMailGunIDSpecialCharactersFn(tt.args.messageID); got != tt.want {
				t.Errorf("StripMailGunIDSpecialCharacters() = %v, want %v", got, tt.want)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.StripMailGunIDSpecialCharacters(tt.args.messageID); got != tt.want {
				t.Errorf("StripMailGunIDSpecialCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}
