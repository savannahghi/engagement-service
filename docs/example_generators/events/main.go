package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/feed/graph/feed"
)

func main() {
	event := getTestEvent()
	data, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		log.Printf("can't marshal event to JSON: %s", err)
	}
	fmt.Printf("\n%s\n", data)
}

func getTestEvent() feed.Event {
	return feed.Event{
		ID:   uuid.New().String(),
		Name: "TEST_EVENT",
		Context: feed.Context{
			UserID:         uuid.New().String(),
			Flavour:        feed.FlavourConsumer,
			OrganizationID: uuid.New().String(),
			LocationID:     uuid.New().String(),
			Timestamp:      time.Now(),
		},
		Payload: feed.Payload{
			Data: map[string]interface{}{
				"a":       "key",
				"another": "key",
			},
		},
	}
}
