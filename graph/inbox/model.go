package inbox

import (
	"time"
)

// Message structure of a message
type Message struct {
	ID            string         `json:"id" firestore:"id"`
	CreatedAt     time.Time      `json:"created_at" firestore:"createdAt"`
	Body          string         `json:"body" firestore:"body"`
	Read          bool           `json:"read" firestore:"read"`
	ReadAt        time.Time      `json:"read_at" firestore:"readAt"`
	SenderName    string         `json:"sender_name" firestore:"senderName"`
	SenderUID     string         `json:"sender_uid" firestore:"senderUID"`
	RecipientName string         `json:"recipient_name" firestore:"recipientName"`
	RecipientUID  string         `json:"recipient_uid" firestore:"recipientUID"`
	Channel       MessageChannel `json:"channel" firestore:"channel"`
	Tags          []MessageTag   `json:"tags" firestore:"tags"`
}

// IsEntity ...
func (m Message) IsEntity() {}

// MessageChannel channel of a message
type MessageChannel struct {
	ID   string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Slug string `json:"slug" firestore:"slug"`
}

// IsEntity ...
func (m MessageChannel) IsEntity() {}

// MessageTag ...
type MessageTag struct {
	ID   string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Slug string `json:"slug" firestore:"slug"`
}

// IsEntity ...
func (m MessageTag) IsEntity() {}
