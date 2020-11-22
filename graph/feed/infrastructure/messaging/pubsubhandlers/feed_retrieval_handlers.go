package pubsubhandlers

import (
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/messaging"
)

// HandlePubsubPayload defines the signature of a function that handles
// payloads received from Google Cloud Pubsub
type HandlePubsubPayload func(m *messaging.PubSubPayload) error

// HandleFeedRetrieval responds to feed retrieval messages
func HandleFeedRetrieval(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("feed retrieval data: \n%s\n", string(m.Message.Data))
	log.Printf("feed retrieval subscription: %s", m.Subscription)
	log.Printf("feed retrieval message ID: %s", m.Message.MessageID)
	log.Printf("feed retrieval message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleThinFeedRetrieval responds to thin feed retrieval messages
func HandleThinFeedRetrieval(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("thin feed retrieval data: \n%s\n", string(m.Message.Data))
	log.Printf("thin feed retrieval subscription: %s", m.Subscription)
	log.Printf("thin feed retrieval message ID: %s", m.Message.MessageID)
	log.Printf("thin feed retrieval message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemRetrieval responds to item retrieval messages
func HandleItemRetrieval(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item retrieval data: \n%s\n", string(m.Message.Data))
	log.Printf("item retrieval subscription: %s", m.Subscription)
	log.Printf("item retrieval message ID: %s", m.Message.MessageID)
	log.Printf("item retrieval message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemPublish responds to item publish messages
func HandleItemPublish(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item publish data: \n%s\n", string(m.Message.Data))
	log.Printf("item publish subscription: %s", m.Subscription)
	log.Printf("item publish message ID: %s", m.Message.MessageID)
	log.Printf("item publish message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemDelete responds to item delete messages
func HandleItemDelete(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item delete data: \n%s\n", string(m.Message.Data))
	log.Printf("item delete subscription: %s", m.Subscription)
	log.Printf("item delete message ID: %s", m.Message.MessageID)
	log.Printf("item delete message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemResolve responds to item resolve messages
func HandleItemResolve(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item resolve data: \n%s\n", string(m.Message.Data))
	log.Printf("item resolve subscription: %s", m.Subscription)
	log.Printf("item resolve message ID: %s", m.Message.MessageID)
	log.Printf("item resolve message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemUnresolve responds to item unresolve messages
func HandleItemUnresolve(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item resolve data: \n%s\n", string(m.Message.Data))
	log.Printf("item resolve subscription: %s", m.Subscription)
	log.Printf("item resolve message ID: %s", m.Message.MessageID)
	log.Printf("item resolve message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemHide responds to item hide messages
func HandleItemHide(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item hide data: \n%s\n", string(m.Message.Data))
	log.Printf("item hide subscription: %s", m.Subscription)
	log.Printf("item hide message ID: %s", m.Message.MessageID)
	log.Printf("item hide message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemShow responds to item show messages
func HandleItemShow(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item show data: \n%s\n", string(m.Message.Data))
	log.Printf("item show subscription: %s", m.Subscription)
	log.Printf("item show message ID: %s", m.Message.MessageID)
	log.Printf("item show message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemPin responds to item pin messages
func HandleItemPin(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item pin data: \n%s\n", string(m.Message.Data))
	log.Printf("item pin subscription: %s", m.Subscription)
	log.Printf("item pin message ID: %s", m.Message.MessageID)
	log.Printf("item pin message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleItemUnpin responds to item unpin messages
func HandleItemUnpin(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("item unpin data: \n%s\n", string(m.Message.Data))
	log.Printf("item unpin subscription: %s", m.Subscription)
	log.Printf("item unpin message ID: %s", m.Message.MessageID)
	log.Printf("item unpin message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeRetrieval responds to nudge retrieval messages
func HandleNudgeRetrieval(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge retrieval data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge retrieval subscription: %s", m.Subscription)
	log.Printf("nudge retrieval message ID: %s", m.Message.MessageID)
	log.Printf("nudge retrieval message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgePublish responds to nudge publish messages
func HandleNudgePublish(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge publish data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge publish subscription: %s", m.Subscription)
	log.Printf("nudge publish message ID: %s", m.Message.MessageID)
	log.Printf("nudge publish message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeDelete responds to nudge delete messages
func HandleNudgeDelete(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge delete data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge delete subscription: %s", m.Subscription)
	log.Printf("nudge delete message ID: %s", m.Message.MessageID)
	log.Printf("nudge delete message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeResolve responds to nudge resolve messages
func HandleNudgeResolve(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge resolve data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge resolve subscription: %s", m.Subscription)
	log.Printf("nudge resolve message ID: %s", m.Message.MessageID)
	log.Printf("nudge resolve message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeUnresolve responds to nudge unresolve messages
func HandleNudgeUnresolve(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge unresolve data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge unresolve subscription: %s", m.Subscription)
	log.Printf("nudge unresolve message ID: %s", m.Message.MessageID)
	log.Printf("nudge unresolve message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeHide responds to nudge hide messages
func HandleNudgeHide(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge hide data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge hide subscription: %s", m.Subscription)
	log.Printf("nudge hide message ID: %s", m.Message.MessageID)
	log.Printf("nudge hide message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleNudgeShow responds to nudge hide messages
func HandleNudgeShow(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("nudge show data: \n%s\n", string(m.Message.Data))
	log.Printf("nudge show subscription: %s", m.Subscription)
	log.Printf("nudge show message ID: %s", m.Message.MessageID)
	log.Printf("nudge show message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleActionRetrieval responds to action retrieval messages
func HandleActionRetrieval(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("action retrieval data: \n%s\n", string(m.Message.Data))
	log.Printf("action retrieval subscription: %s", m.Subscription)
	log.Printf("action retrieval message ID: %s", m.Message.MessageID)
	log.Printf("action retrieval message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleActionPublish responds to action publish messages
func HandleActionPublish(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("action publish data: \n%s\n", string(m.Message.Data))
	log.Printf("action publish subscription: %s", m.Subscription)
	log.Printf("action publish message ID: %s", m.Message.MessageID)
	log.Printf("action publish message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleActionDelete responds to action publish messages
func HandleActionDelete(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("action delete data: \n%s\n", string(m.Message.Data))
	log.Printf("action delete subscription: %s", m.Subscription)
	log.Printf("action delete message ID: %s", m.Message.MessageID)
	log.Printf("action delete message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleMessagePost responds to message post pubsub messages
func HandleMessagePost(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("message post data: \n%s\n", string(m.Message.Data))
	log.Printf("message post subscription: %s", m.Subscription)
	log.Printf("message post message ID: %s", m.Message.MessageID)
	log.Printf("message post message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleMessageDelete responds to message delete pubsub messages
func HandleMessageDelete(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("message delete data: \n%s\n", string(m.Message.Data))
	log.Printf("message delete subscription: %s", m.Subscription)
	log.Printf("message delete message ID: %s", m.Message.MessageID)
	log.Printf("message delete message attributes: %#v", m.Message.Attributes)
	return nil
}

// HandleIncomingEvent responds to message delete pubsub messages
func HandleIncomingEvent(m *messaging.PubSubPayload) error {
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	log.Printf("incoming event data: \n%s\n", string(m.Message.Data))
	log.Printf("incoming event subscription: %s", m.Subscription)
	log.Printf("incoming event message ID: %s", m.Message.MessageID)
	log.Printf("incoming event message attributes: %#v", m.Message.Attributes)
	return nil
}
