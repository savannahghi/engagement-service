package mock

import (
	"time"

	"github.com/savannahghi/feedlib"
	"go.opentelemetry.io/otel/trace"
)

// FakeHelpers ...
type FakeHelpers struct {
	// AddPubSubNameSpaceFn ...
	AddPubSubNameSpaceFn func(topicName string) string
	// ValidateElementFn ...
	ValidateElementFn func(el feedlib.Element) error
	// RecordSpanErrorFn ...
	RecordSpanErrorFn func(span trace.Span, err error)
	// EpochTimeToStandardTimeFn ....
	EpochTimeToStandardTimeFn func(timeOfDelivery string) time.Time
	// StripMailGunIDSpecialCharactersFn ...
	StripMailGunIDSpecialCharactersFn func(messageID string) string
}

// AddPubSubNameSpace is a mock version of the original function
func (f *FakeHelpers) AddPubSubNameSpace(
	topicName string,
) string {
	return f.AddPubSubNameSpaceFn(topicName)
}

// ValidateElement is a mock version of the original function
func (f *FakeHelpers) ValidateElement(
	el feedlib.Element,
) error {
	return f.ValidateElementFn(el)
}

// RecordSpanError is a mock version of the original function
func (f *FakeHelpers) RecordSpanError(
	span trace.Span,
	err error,
) {
}

// EpochTimeToStandardTime is a mock version of the original function
func (f *FakeHelpers) EpochTimeToStandardTime(
	timeOfDelivery string,
) time.Time {
	return f.EpochTimeToStandardTimeFn(timeOfDelivery)
}

// StripMailGunIDSpecialCharacters is a mock version of the original function
func (f *FakeHelpers) StripMailGunIDSpecialCharacters(
	messageID string,
) string {
	return f.StripMailGunIDSpecialCharactersFn(messageID)
}
