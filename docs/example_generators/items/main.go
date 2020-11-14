package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/feed/graph/feed"
)

const (
	base64PNGSample = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAFNeavDAAAACklEQVQIHWNgAAAAAgABz8g15QAAAABJRU5ErkJggg=="
	base64PDFSample = "JVBERi0xLjUKJbXtrvsKNCAwIG9iago8PCAvTGVuZ3RoIDUgMCBSCiAgIC9GaWx0ZXIgL0ZsYXRlRGVjb2RlCj4+CnN0cmVhbQp4nDNUMABCXUMQpWdkopCcy1XIFcgFADCwBFQKZW5kc3RyZWFtCmVuZG9iago1IDAgb2JqCiAgIDI3CmVuZG9iagozIDAgb2JqCjw8Cj4+CmVuZG9iagoyIDAgb2JqCjw8IC9UeXBlIC9QYWdlICUgMQogICAvUGFyZW50IDEgMCBSCiAgIC9NZWRpYUJveCBbIDAgMCAwLjI0IDAuMjQgXQogICAvQ29udGVudHMgNCAwIFIKICAgL0dyb3VwIDw8CiAgICAgIC9UeXBlIC9Hcm91cAogICAgICAvUyAvVHJhbnNwYXJlbmN5CiAgICAgIC9JIHRydWUKICAgICAgL0NTIC9EZXZpY2VSR0IKICAgPj4KICAgL1Jlc291cmNlcyAzIDAgUgo+PgplbmRvYmoKMSAwIG9iago8PCAvVHlwZSAvUGFnZXMKICAgL0tpZHMgWyAyIDAgUiBdCiAgIC9Db3VudCAxCj4+CmVuZG9iago2IDAgb2JqCjw8IC9Qcm9kdWNlciAoY2Fpcm8gMS4xNi4wIChodHRwczovL2NhaXJvZ3JhcGhpY3Mub3JnKSkKICAgL0NyZWF0aW9uRGF0ZSAoRDoyMDIwMTAzMDA4MDkwOCswMycwMCkKPj4KZW5kb2JqCjcgMCBvYmoKPDwgL1R5cGUgL0NhdGFsb2cKICAgL1BhZ2VzIDEgMCBSCj4+CmVuZG9iagp4cmVmCjAgOAowMDAwMDAwMDAwIDY1NTM1IGYgCjAwMDAwMDAzODEgMDAwMDAgbiAKMDAwMDAwMDE2MSAwMDAwMCBuIAowMDAwMDAwMTQwIDAwMDAwIG4gCjAwMDAwMDAwMTUgMDAwMDAgbiAKMDAwMDAwMDExOSAwMDAwMCBuIAowMDAwMDAwNDQ2IDAwMDAwIG4gCjAwMDAwMDA1NjIgMDAwMDAgbiAKdHJhaWxlcgo8PCAvU2l6ZSA4CiAgIC9Sb290IDcgMCBSCiAgIC9JbmZvIDYgMCBSCj4+CnN0YXJ0eHJlZgo2MTQKJSVFT0YK"
	sampleVideoURL  = "https://www.youtube.com/watch?v=bPiofmZGb8o"
)

func main() {
	item := getTestItem()
	data, err := json.MarshalIndent(item, "", "    ")
	if err != nil {
		log.Printf("can't marshal test item to JSON: %s", err)
	}
	fmt.Printf("\n%s\n", data)
}

func getTestItem() feed.Item {
	return feed.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon: feed.Image{
			ID:     "icon-1",
			Base64: base64PNGSample,
		},
		Author:    "Bot 1",
		Tagline:   "Bot speaks...",
		Label:     "DRUGS",
		Timestamp: time.Now(),
		Summary:   "I am a bot...",
		Text:      "This bot can speak",
		Images: []feed.Image{
			{
				ID:     "img-1",
				Base64: base64PNGSample,
			},
		},
		Videos: []feed.Video{
			{
				ID:  "video-1",
				URL: "https://www.youtube.com/watch?v=bPiofmZGb8o",
			},
		},
		Actions: []feed.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				ActionType:     feed.ActionTypeSecondary,
				Handling:       feed.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
			},
		},
		Conversations: []feed.Message{
			{
				ID:           "msg-2",
				Text:         "hii ni reply",
				ReplyTo:      "msg-1",
				PostedByName: ksuid.New().String(),
				PostedByUID:  ksuid.New().String(),
				Timestamp:    time.Now(),
			},
		},
		Documents: []feed.Document{
			getTestDocument(),
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelFcm,
			feed.ChannelEmail,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func getTestDocument() feed.Document {
	return feed.Document{
		ID:     ksuid.New().String(),
		Base64: base64PDFSample,
	}
}
