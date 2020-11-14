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

func getTestAction() feed.Action {
	return feed.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
	}
}

func getTestSequenceNumber() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(intMax) // #nosec
}
