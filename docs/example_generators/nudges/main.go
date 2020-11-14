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

const (
	base64PNGSample = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAFNeavDAAAACklEQVQIHWNgAAAAAgABz8g15QAAAABJRU5ErkJggg=="
	intMax          = 9007199254740990
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
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Title:          ksuid.New().String(),
		Image:          getTestImage(),
		Text:           ksuid.New().String(),
		Actions: []feed.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
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
		ID:     ksuid.New().String(),
		Base64: base64PNGSample,
	}
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
