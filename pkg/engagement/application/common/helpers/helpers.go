package helpers

import (
	"fmt"

	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/base"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	// ServiceName ...
	ServiceName = "engagement"
	// TopicVersion ...
	TopicVersion = "v1"
)

// AddPubSubNamespace creates a namespaced topic name
func AddPubSubNamespace(topicName string) string {
	environment := serverutils.GetRunningEnvironment()
	return base.NamespacePubsubIdentifier(
		ServiceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// ValidateElement ensures that an element is non nil and valid
func ValidateElement(el base.Element) error {
	if el == nil {
		return fmt.Errorf("nil element")
	}

	_, err := el.ValidateAndMarshal()
	if err != nil {
		return fmt.Errorf("element failed validation: %w", err)
	}

	return nil
}

// RecordSpanError is a helper function to capture errors in a span
func RecordSpanError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
