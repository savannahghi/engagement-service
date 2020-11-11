package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
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
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: feed.Context{
			UserID:         ksuid.New().String(),
			Flavour:        feed.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
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
