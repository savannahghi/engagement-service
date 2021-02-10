package resources

import "gitlab.slade360emr.com/go/base"

// NotificationEnvelope is used to "wrap" elements with context and metadata
// before they are sent as notifications.
//
// This context and metadata allows the recipients of the notifications to
// process them intelligently.
type NotificationEnvelope struct {
	UID      string                 `json:"uid"`
	Flavour  base.Flavour           `json:"flavour"`
	Payload  []byte                 `json:"payload"`
	Metadata map[string]interface{} `json:"metadata"`
}
