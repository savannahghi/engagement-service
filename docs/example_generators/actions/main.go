package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/feed/graph/feed"
)

const intMax = 9007199254740990

func main() {
	action := getTestAction()
	data, err := json.MarshalIndent(action, "", "    ")
	if err != nil {
		log.Printf("can't marshal action to JSON: %s", err)
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
	}
}

func getTestAction() feed.Action {
	return feed.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
		Event:          getTestEvent(),
	}
}

func getTestSequenceNumber() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(intMax) // #nosec
}
