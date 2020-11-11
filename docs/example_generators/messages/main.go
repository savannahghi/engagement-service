package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/feed/graph/feed"
)

const intMax = 9007199254740990

func main() {
	msg := getTestMessage()
	data, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		log.Printf("can't marshal message to JSON: %s", err)
	}
	fmt.Printf("\n%s\n", data)
}

func getTestMessage() feed.Message {
	return feed.Message{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           uuid.New().String(),
		ReplyTo:        uuid.New().String(),
		PostedByUID:    uuid.New().String(),
		PostedByName:   uuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestSequenceNumber() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(intMax) // #nosec
}
