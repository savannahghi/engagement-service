package helpers

import (
	"fmt"

	"gitlab.slade360emr.com/go/base"
)

const (
	// ServiceName ...
	ServiceName = "engagement"
	// TopicVersion ...
	TopicVersion = "v1"
)

// AddPubSubNamespace creates a namespaced topic name
func AddPubSubNamespace(topicName string) string {
	environment := base.GetRunningEnvironment()
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
