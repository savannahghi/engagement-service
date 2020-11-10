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

const (
	base64PNGSample = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAFNeavDAAAACklEQVQIHWNgAAAAAgABz8g15QAAAABJRU5ErkJggg=="
	intMax          = 9223372036854775807
)

func main() {
	nudge := testNudge()
	data, err := json.MarshalIndent(nudge, "", "    ")
	if err != nil {
		log.Printf("can't marshal nudge to JSON: %s", err)
	}
	fmt.Printf("\n%s\n", data)
}

func testNudge() *feed.Nudge {
	return &feed.Nudge{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          uuid.New().String(),
		Image:          getTestImage(),
		Text:           uuid.New().String(),
		Actions: []feed.Action{
			getTestAction(),
		},
		Users: []string{
			uuid.New().String(),
		},
		Groups: []string{
			uuid.New().String(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func getTestSequenceNumber() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(intMax) // #nosec
}

func getTestImage() feed.Image {
	return feed.Image{
		ID:     uuid.New().String(),
		Base64: base64PNGSample,
	}
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
	}
}

func getTestAction() feed.Action {
	return feed.Action{
		ID:             uuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		ActionType:     feed.ActionTypePrimary,
		Handling:       feed.HandlingFullPage,
		Event:          getTestEvent(),
	}
}