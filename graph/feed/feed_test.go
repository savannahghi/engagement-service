package feed_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/feed/graph/feed"
	db "gitlab.slade360emr.com/go/feed/graph/feed/infrastructure/database"
)

const (
	base64PNGSample = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAFNeavDAAAACklEQVQIHWNgAAAAAgABz8g15QAAAABJRU5ErkJggg=="
	base64PDFSample = "JVBERi0xLjUKJbXtrvsKNCAwIG9iago8PCAvTGVuZ3RoIDUgMCBSCiAgIC9GaWx0ZXIgL0ZsYXRlRGVjb2RlCj4+CnN0cmVhbQp4nDNUMABCXUMQpWdkopCcy1XIFcgFADCwBFQKZW5kc3RyZWFtCmVuZG9iago1IDAgb2JqCiAgIDI3CmVuZG9iagozIDAgb2JqCjw8Cj4+CmVuZG9iagoyIDAgb2JqCjw8IC9UeXBlIC9QYWdlICUgMQogICAvUGFyZW50IDEgMCBSCiAgIC9NZWRpYUJveCBbIDAgMCAwLjI0IDAuMjQgXQogICAvQ29udGVudHMgNCAwIFIKICAgL0dyb3VwIDw8CiAgICAgIC9UeXBlIC9Hcm91cAogICAgICAvUyAvVHJhbnNwYXJlbmN5CiAgICAgIC9JIHRydWUKICAgICAgL0NTIC9EZXZpY2VSR0IKICAgPj4KICAgL1Jlc291cmNlcyAzIDAgUgo+PgplbmRvYmoKMSAwIG9iago8PCAvVHlwZSAvUGFnZXMKICAgL0tpZHMgWyAyIDAgUiBdCiAgIC9Db3VudCAxCj4+CmVuZG9iago2IDAgb2JqCjw8IC9Qcm9kdWNlciAoY2Fpcm8gMS4xNi4wIChodHRwczovL2NhaXJvZ3JhcGhpY3Mub3JnKSkKICAgL0NyZWF0aW9uRGF0ZSAoRDoyMDIwMTAzMDA4MDkwOCswMycwMCkKPj4KZW5kb2JqCjcgMCBvYmoKPDwgL1R5cGUgL0NhdGFsb2cKICAgL1BhZ2VzIDEgMCBSCj4+CmVuZG9iagp4cmVmCjAgOAowMDAwMDAwMDAwIDY1NTM1IGYgCjAwMDAwMDAzODEgMDAwMDAgbiAKMDAwMDAwMDE2MSAwMDAwMCBuIAowMDAwMDAwMTQwIDAwMDAwIG4gCjAwMDAwMDAwMTUgMDAwMDAgbiAKMDAwMDAwMDExOSAwMDAwMCBuIAowMDAwMDAwNDQ2IDAwMDAwIG4gCjAwMDAwMDA1NjIgMDAwMDAgbiAKdHJhaWxlcgo8PCAvU2l6ZSA4CiAgIC9Sb290IDcgMCBSCiAgIC9JbmZvIDYgMCBSCj4+CnN0YXJ0eHJlZgo2MTQKJSVFT0YK"
	sampleVideoURL  = "https://www.youtube.com/watch?v=bPiofmZGb8o"

	intMax = 9007199254740990
)

func getEmptyJson(t *testing.T) []byte {
	emptyJSONBytes, err := json.Marshal(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, emptyJSONBytes)
	return emptyJSONBytes
}

func getTestItem() feed.Item {
	return feed.Item{
		ID:             "item-1",
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
				Event: feed.Event{
					ID:   "event-1",
					Name: "THIS_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourPro,
						OrganizationID: "org-1",
						LocationID:     "loc-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": 1},
					},
				},
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
				Event: feed.Event{
					ID:   "event-1",
					Name: "AN_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourConsumer,
						LocationID:     "location-1",
						OrganizationID: "organization-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": "1"},
					},
				},
			},
		},
		Conversations: []feed.Message{
			{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
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

func TestNewInMemoryFeed(t *testing.T) {
	feeds := getTestFeedAggregate(t)
	assert.NotNil(t, feeds)
}

func TestVideo_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)
	validVideo := feed.Video{
		ID:  ksuid.New().String(),
		URL: "https://www.youtube.com/watch?v=mlv36Yxy3Wk",
	}
	validVideoJSONBytes, err := json.Marshal(validVideo)
	assert.Nil(t, err)
	assert.NotNil(t, validVideoJSONBytes)
	assert.Greater(t, len(validVideoJSONBytes), 3)

	type fields struct {
		ID  string
		URL string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validVideoJSONBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vi := &feed.Video{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			if err := vi.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Video.UnmarshalJSON() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if !tt.wantErr {
				assert.NotZero(t, vi.ID)
				assert.NotZero(t, vi.URL)
			}
		})
	}
}

func TestMessage_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Text:           "some message text",
		PostedByName:   ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		Timestamp:      time.Now(),
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &feed.Message{}
			if err := msg.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Message.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestItem_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Item{
		ID:             "item-1",
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
				Event: feed.Event{
					ID:   "event-1",
					Name: "THIS_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourConsumer,
						OrganizationID: "org-1",
						LocationID:     "loc-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": 1},
					},
				},
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
				Event: feed.Event{
					ID:   "event-1",
					Name: "AN_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourPro,
						LocationID:     "location-1",
						OrganizationID: "organization-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": "1"},
					},
				},
			},
		},
		Conversations: []feed.Message{
			{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
			},
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		Documents: []feed.Document{
			getTestDocument(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelFcm,
			feed.ChannelEmail,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Expiry         time.Time
		Persistent     bool
		Status         feed.Status
		Visibility     feed.Visibility
		Icon           feed.Image
		Author         string
		Tagline        string
		Label          string
		Timestamp      time.Time
		Summary        string
		Text           string
		Images         []feed.Image
		Videos         []feed.Video
		Actions        []feed.Action
		Conversations  []feed.Message
		Users          []string
		Groups         []string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &feed.Item{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Expiry:         tt.fields.Expiry,
				Persistent:     tt.fields.Persistent,
				Status:         tt.fields.Status,
				Visibility:     tt.fields.Visibility,
				Icon:           tt.fields.Icon,
				Author:         tt.fields.Author,
				Tagline:        tt.fields.Tagline,
				Label:          tt.fields.Label,
				Timestamp:      tt.fields.Timestamp,
				Summary:        tt.fields.Summary,
				Text:           tt.fields.Text,
				Images:         tt.fields.Images,
				Videos:         tt.fields.Videos,
				Actions:        tt.fields.Actions,
				Conversations:  tt.fields.Conversations,
				Users:          tt.fields.Users,
				Groups:         tt.fields.Groups,
			}
			if err := it.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Item.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestNudge_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Nudge{
		ID:             "nudge-1",
		SequenceNumber: 1,
		Visibility:     feed.VisibilityShow,
		Status:         feed.StatusPending,
		Title:          "Update your profile!",
		Image: feed.Image{
			ID:     "image-1",
			Base64: base64PNGSample,
		},
		Text: "An up to date profile will help us serve you better!",
		Actions: []feed.Action{
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
				Event: feed.Event{
					ID:   "event-1",
					Name: "AN_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourConsumer,
						LocationID:     "location-1",
						OrganizationID: "organization-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": "1"},
					},
				},
			},
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelFcm,
			feed.ChannelEmail,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Visibility     feed.Visibility
		Status         feed.Status
		Title          string
		Image          feed.Image
		Text           string
		Actions        []feed.Action
		Groups         []string
		Users          []string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &feed.Nudge{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Visibility:     tt.fields.Visibility,
				Status:         tt.fields.Status,
				Title:          tt.fields.Title,
				Image:          tt.fields.Image,
				Text:           tt.fields.Text,
				Actions:        tt.fields.Actions,
				Groups:         tt.fields.Groups,
				Users:          tt.fields.Users,
			}
			if err := nu.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Nudge.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestAction_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Name:           "ACTION_NAME",
		ActionType:     feed.ActionTypeSecondary,
		Handling:       feed.HandlingFullPage,
		Event: feed.Event{
			ID:   "event-1",
			Name: "THIS_EVENT",
			Context: feed.Context{
				UserID:         "user-1",
				Flavour:        feed.FlavourConsumer,
				OrganizationID: "org-1",
				LocationID:     "loc-1",
				Timestamp:      time.Now(),
			},
			Payload: feed.Payload{
				Data: map[string]interface{}{"a": 1},
			},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Name           string
		ActionType     feed.ActionType
		Handling       feed.Handling
		Event          feed.Event
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &feed.Action{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Name:           tt.fields.Name,
				ActionType:     tt.fields.ActionType,
				Handling:       tt.fields.Handling,
				Event:          tt.fields.Event,
			}
			if err := ac.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Action.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFeed_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)
	validElement := feed.Feed{
		UID:     "a-uid",
		Flavour: feed.FlavourConsumer,
		Actions: []feed.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				ActionType:     feed.ActionTypeSecondary,
				Handling:       feed.HandlingFullPage,
				Event: feed.Event{
					ID:   "event-1",
					Name: "THIS_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourConsumer,
						OrganizationID: "org-1",
						LocationID:     "loc-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": 1},
					},
				},
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
				Event: feed.Event{
					ID:   "event-1",
					Name: "AN_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourPro,
						LocationID:     "location-1",
						OrganizationID: "organization-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": "1"},
					},
				},
			},
		},
		Nudges: []feed.Nudge{
			{
				ID:             "nudge-1",
				SequenceNumber: 1,
				Visibility:     feed.VisibilityShow,
				Status:         feed.StatusPending,
				Title:          "Update your profile!",
				Image: feed.Image{
					ID:     "image-1",
					Base64: base64PNGSample,
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []feed.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						ActionType:     feed.ActionTypePrimary,
						Handling:       feed.HandlingInline,
						Event: feed.Event{
							ID:   "event-1",
							Name: "AN_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								LocationID:     "location-1",
								OrganizationID: "organization-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": "1"},
							},
						},
					},
				},
				Groups: []string{
					"group-1",
					"group-2",
				},
				Users: []string{
					"user-1",
					"user-2",
				},
				NotificationChannels: []feed.Channel{
					feed.ChannelFcm,
					feed.ChannelEmail,
					feed.ChannelSms,
					feed.ChannelWhatsapp,
				},
			},
		},
		Items: []feed.Item{
			{
				ID:             "item-1",
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
						Event: feed.Event{
							ID:   "event-1",
							Name: "THIS_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								OrganizationID: "org-1",
								LocationID:     "loc-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": 1},
							},
						},
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						ActionType:     feed.ActionTypePrimary,
						Handling:       feed.HandlingInline,
						Event: feed.Event{
							ID:   "event-1",
							Name: "AN_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourPro,
								LocationID:     "location-1",
								OrganizationID: "organization-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": "1"},
							},
						},
					},
				},
				Conversations: []feed.Message{
					{
						ID:             "msg-2",
						SequenceNumber: 1,
						Text:           "hii ni reply",
						ReplyTo:        "msg-1",
						PostedByName:   ksuid.New().String(),
						PostedByUID:    ksuid.New().String(),
						Timestamp:      time.Now(),
					},
				},
				Users: []string{
					"user-1",
					"user-2",
				},
				Groups: []string{
					"group-1",
					"group-2",
				},
				Documents: []feed.Document{
					getTestDocument(),
				},
				NotificationChannels: []feed.Channel{
					feed.ChannelFcm,
					feed.ChannelEmail,
					feed.ChannelSms,
					feed.ChannelWhatsapp,
				},
			},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		UID     string
		Flavour feed.Flavour
		Actions []feed.Action
		Items   []feed.Item
		Nudges  []feed.Nudge
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := &feed.Feed{
				UID:     tt.fields.UID,
				Flavour: tt.fields.Flavour,
				Actions: tt.fields.Actions,
				Items:   tt.fields.Items,
				Nudges:  tt.fields.Nudges,
			}
			if err := fe.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Feed.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFeed_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		UID     string
		Flavour feed.Flavour
		Actions []feed.Action
		Items   []feed.Item
		Nudges  []feed.Nudge
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid feed",
			fields: fields{
				UID:     "a-uid",
				Flavour: feed.FlavourPro,
				Actions: []feed.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						ActionType:     feed.ActionTypeSecondary,
						Handling:       feed.HandlingFullPage,
						Event: feed.Event{
							ID:   "event-1",
							Name: "THIS_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								OrganizationID: "org-1",
								LocationID:     "loc-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": 1},
							},
						},
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						ActionType:     feed.ActionTypePrimary,
						Handling:       feed.HandlingInline,
						Event: feed.Event{
							ID:   "event-1",
							Name: "AN_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								LocationID:     "location-1",
								OrganizationID: "organization-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": "1"},
							},
						},
					},
				},
				Nudges: []feed.Nudge{
					{
						ID:             "nudge-1",
						SequenceNumber: 1,
						Visibility:     feed.VisibilityShow,
						Status:         feed.StatusPending,
						Title:          "Update your profile!",
						Image: feed.Image{
							ID:     "image-1",
							Base64: base64PNGSample,
						},
						Text: "Help us serve you better!",
						Actions: []feed.Action{
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								ActionType:     feed.ActionTypePrimary,
								Handling:       feed.HandlingInline,
								Event: feed.Event{
									ID:   "event-1",
									Name: "AN_EVENT",
									Context: feed.Context{
										UserID:         "user-1",
										Flavour:        feed.FlavourConsumer,
										LocationID:     "location-1",
										OrganizationID: "organization-1",
										Timestamp:      time.Now(),
									},
									Payload: feed.Payload{
										Data: map[string]interface{}{"a": "1"},
									},
								},
							},
						},
						Groups: []string{
							"group-1",
							"group-2",
						},
						Users: []string{
							"user-1",
							"user-2",
						},
						NotificationChannels: []feed.Channel{
							feed.ChannelFcm,
							feed.ChannelEmail,
							feed.ChannelSms,
							feed.ChannelWhatsapp,
						},
					},
				},
				Items: []feed.Item{
					{
						ID:             "item-1",
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
								URL: sampleVideoURL,
							},
						},
						Actions: []feed.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           "ACTION_NAME",
								ActionType:     feed.ActionTypeSecondary,
								Handling:       feed.HandlingFullPage,
								Event: feed.Event{
									ID:   "event-1",
									Name: "THIS_EVENT",
									Context: feed.Context{
										UserID:         "user-1",
										Flavour:        feed.FlavourConsumer,
										OrganizationID: "org-1",
										LocationID:     "loc-1",
										Timestamp:      time.Now(),
									},
									Payload: feed.Payload{
										Data: map[string]interface{}{"a": 1},
									},
								},
							},
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								ActionType:     feed.ActionTypePrimary,
								Handling:       feed.HandlingInline,
								Event: feed.Event{
									ID:   "event-1",
									Name: "AN_EVENT",
									Context: feed.Context{
										UserID:         "user-1",
										Flavour:        feed.FlavourPro,
										LocationID:     "location-1",
										OrganizationID: "organization-1",
										Timestamp:      time.Now(),
									},
									Payload: feed.Payload{
										Data: map[string]interface{}{"a": "1"},
									},
								},
							},
						},
						Conversations: []feed.Message{
							{
								ID:             "msg-2",
								SequenceNumber: 1,
								Text:           "hii ni reply",
								ReplyTo:        "msg-1",
								PostedByName:   ksuid.New().String(),
								PostedByUID:    ksuid.New().String(),
								Timestamp:      time.Now(),
							},
						},
						Users: []string{
							"user-1",
							"user-2",
						},
						Groups: []string{
							"group-1",
							"group-2",
						},
						Documents: []feed.Document{
							getTestDocument(),
						},
						NotificationChannels: []feed.Channel{
							feed.ChannelFcm,
							feed.ChannelEmail,
							feed.ChannelSms,
							feed.ChannelWhatsapp,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := &feed.Feed{
				UID:     tt.fields.UID,
				Flavour: tt.fields.Flavour,
				Actions: tt.fields.Actions,
				Items:   tt.fields.Items,
				Nudges:  tt.fields.Nudges,
			}
			got, err := fe.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Feed.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestAction_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID             string
		SequenceNumber int
		Name           string
		ActionType     feed.ActionType
		Handling       feed.Handling
		Event          feed.Event
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid action",
			fields: fields{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				ActionType:     feed.ActionTypePrimary,
				Handling:       feed.HandlingInline,
				Event: feed.Event{
					ID:   "event-1",
					Name: "AN_EVENT",
					Context: feed.Context{
						UserID:         "user-1",
						Flavour:        feed.FlavourConsumer,
						LocationID:     "location-1",
						OrganizationID: "organization-1",
						Timestamp:      time.Now(),
					},
					Payload: feed.Payload{
						Data: map[string]interface{}{"a": "1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &feed.Action{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Name:           tt.fields.Name,
				ActionType:     tt.fields.ActionType,
				Handling:       tt.fields.Handling,
				Event:          tt.fields.Event,
			}
			got, err := ac.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Action.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestNudge_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID                   string
		SequenceNumber       int
		Visibility           feed.Visibility
		Status               feed.Status
		Title                string
		Image                feed.Image
		Text                 string
		Actions              []feed.Action
		Groups               []string
		Users                []string
		NotificationChannels []feed.Channel
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case - valid nudge",
			fields: fields{
				ID:             "nudge-1",
				SequenceNumber: 1,
				Visibility:     feed.VisibilityShow,
				Status:         feed.StatusPending,
				Title:          "Update your profile!",
				Image: feed.Image{
					ID:     "image-1",
					Base64: base64PNGSample,
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []feed.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						ActionType:     feed.ActionTypePrimary,
						Handling:       feed.HandlingInline,
						Event: feed.Event{
							ID:   "event-1",
							Name: "AN_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								LocationID:     "location-1",
								OrganizationID: "organization-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": "1"},
							},
						},
					},
				},
				Groups: []string{
					"group-1",
					"group-2",
				},
				Users: []string{
					"user-1",
					"user-2",
				},
				NotificationChannels: []feed.Channel{
					feed.ChannelFcm,
					feed.ChannelEmail,
					feed.ChannelSms,
					feed.ChannelWhatsapp,
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &feed.Nudge{
				ID:                   tt.fields.ID,
				SequenceNumber:       tt.fields.SequenceNumber,
				Visibility:           tt.fields.Visibility,
				Status:               tt.fields.Status,
				Title:                tt.fields.Title,
				Image:                tt.fields.Image,
				Text:                 tt.fields.Text,
				Actions:              tt.fields.Actions,
				Groups:               tt.fields.Groups,
				Users:                tt.fields.Users,
				NotificationChannels: tt.fields.NotificationChannels,
			}
			got, err := nu.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Nudge.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestItem_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID                   string
		SequenceNumber       int
		Expiry               time.Time
		Persistent           bool
		Status               feed.Status
		Visibility           feed.Visibility
		Icon                 feed.Image
		Author               string
		Tagline              string
		Label                string
		Timestamp            time.Time
		Summary              string
		Text                 string
		Images               []feed.Image
		Documents            []feed.Document
		Videos               []feed.Video
		Actions              []feed.Action
		Conversations        []feed.Message
		Users                []string
		Groups               []string
		NotificationChannels []feed.Channel
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case - valid item",
			fields: fields{
				ID:             "item-1",
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
				Documents: []feed.Document{
					getTestDocument(),
				},
				Videos: []feed.Video{
					{
						ID:  "video-1",
						URL: sampleVideoURL,
					},
				},
				Actions: []feed.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						ActionType:     feed.ActionTypeSecondary,
						Handling:       feed.HandlingFullPage,
						Event: feed.Event{
							ID:   "event-1",
							Name: "THIS_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourConsumer,
								OrganizationID: "org-1",
								LocationID:     "loc-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": 1},
							},
						},
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						ActionType:     feed.ActionTypePrimary,
						Handling:       feed.HandlingInline,
						Event: feed.Event{
							ID:   "event-1",
							Name: "AN_EVENT",
							Context: feed.Context{
								UserID:         "user-1",
								Flavour:        feed.FlavourPro,
								LocationID:     "location-1",
								OrganizationID: "organization-1",
								Timestamp:      time.Now(),
							},
							Payload: feed.Payload{
								Data: map[string]interface{}{"a": "1"},
							},
						},
					},
				},
				Conversations: []feed.Message{
					{
						ID:             "msg-2",
						SequenceNumber: 1,
						Text:           "hii ni reply",
						ReplyTo:        "msg-1",
						PostedByName:   ksuid.New().String(),
						PostedByUID:    ksuid.New().String(),
						Timestamp:      time.Now(),
					},
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
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &feed.Item{
				ID:                   tt.fields.ID,
				SequenceNumber:       tt.fields.SequenceNumber,
				Expiry:               tt.fields.Expiry,
				Persistent:           tt.fields.Persistent,
				Status:               tt.fields.Status,
				Visibility:           tt.fields.Visibility,
				Icon:                 tt.fields.Icon,
				Author:               tt.fields.Author,
				Tagline:              tt.fields.Tagline,
				Label:                tt.fields.Label,
				Timestamp:            tt.fields.Timestamp,
				Summary:              tt.fields.Summary,
				Text:                 tt.fields.Text,
				Images:               tt.fields.Images,
				Documents:            tt.fields.Documents,
				Videos:               tt.fields.Videos,
				Actions:              tt.fields.Actions,
				Conversations:        tt.fields.Conversations,
				Users:                tt.fields.Users,
				Groups:               tt.fields.Groups,
				NotificationChannels: tt.fields.NotificationChannels,
			}
			got, err := it.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Item.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestMessage_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID             string
		SequenceNumber int
		Text           string
		ReplyTo        string
		PostedByName   string
		PostedByUID    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "this is a message",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &feed.Message{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Text:           tt.fields.Text,
				ReplyTo:        tt.fields.ReplyTo,
				PostedByName:   tt.fields.PostedByName,
				PostedByUID:    tt.fields.PostedByUID,
				Timestamp:      time.Now(),
			}
			got, err := msg.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Message.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestVideo_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID  string
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:  "video-1",
				URL: "https://y.yb/12345",
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vi := &feed.Video{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			got, err := vi.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Video.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestImage_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Image{
		ID:     ksuid.New().String(),
		Base64: base64PNGSample,
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID     string
		Base64 string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &feed.Image{
				ID:     tt.fields.ID,
				Base64: tt.fields.Base64,
			}
			if err := im.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Image.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestImage_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID     string
		Base64 string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:     "image-1",
				Base64: base64PNGSample,
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &feed.Image{
				ID:     tt.fields.ID,
				Base64: tt.fields.Base64,
			}
			got, err := im.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Image.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestContext_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Context{
		UserID:         "uid-1",
		Flavour:        feed.FlavourConsumer,
		OrganizationID: "org-1",
		LocationID:     "loc-1",
		Timestamp:      time.Now(),
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		UserID         string
		OrganizationID string
		LocationID     string
		Timestamp      time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &feed.Context{
				UserID:         tt.fields.UserID,
				OrganizationID: tt.fields.OrganizationID,
				LocationID:     tt.fields.LocationID,
				Timestamp:      tt.fields.Timestamp,
			}
			if err := ct.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Context.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestContext_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		UserID         string
		Flavour        feed.Flavour
		OrganizationID string
		LocationID     string
		Timestamp      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				UserID:         "uid-1",
				Flavour:        feed.FlavourConsumer,
				OrganizationID: "org-1",
				LocationID:     "loc-1",
				Timestamp:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &feed.Context{
				UserID:         tt.fields.UserID,
				Flavour:        tt.fields.Flavour,
				OrganizationID: tt.fields.OrganizationID,
				LocationID:     tt.fields.LocationID,
				Timestamp:      tt.fields.Timestamp,
			}
			got, err := ct.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Context.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestPayload_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Payload{
		Data: map[string]interface{}{"a": 1},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		Data map[string]interface{}
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &feed.Payload{
				Data: tt.fields.Data,
			}
			if err := pl.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Payload.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestPayload_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		Data map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				Data: map[string]interface{}{"a": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &feed.Payload{
				Data: tt.fields.Data,
			}
			got, err := pl.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Payload.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestEvent_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feed.Event{
		ID:   "event-1",
		Name: "THIS_EVENT",
		Context: feed.Context{
			UserID:         "user-1",
			Flavour:        feed.FlavourConsumer,
			OrganizationID: "org-1",
			LocationID:     "loc-1",
			Timestamp:      time.Now(),
		},
		Payload: feed.Payload{
			Data: map[string]interface{}{"a": 1},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		ID      string
		Name    string
		Context feed.Context
		Payload feed.Payload
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &feed.Event{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Context: tt.fields.Context,
				Payload: tt.fields.Payload,
			}
			if err := ev.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Event.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestEvent_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID      string
		Name    string
		Context feed.Context
		Payload feed.Payload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:   "event-1",
				Name: "THIS_EVENT",
				Context: feed.Context{
					UserID:         "user-1",
					Flavour:        feed.FlavourConsumer,
					OrganizationID: "org-1",
					LocationID:     "loc-1",
					Timestamp:      time.Now(),
				},
				Payload: feed.Payload{
					Data: map[string]interface{}{"a": 1},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &feed.Event{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Context: tt.fields.Context,
				Payload: tt.fields.Payload,
			}
			got, err := ev.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Event.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestFeed_GetItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "valid case - item that exists",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid case - item that does not exist",
			args: args{
				ctx:    ctx,
				itemID: ksuid.New().String(),
			},
			wantErr: false,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.GetFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_GetNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	flavour := feed.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := fr.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "valid case - nudge that exists",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid case - nudge that does not exist",
			args: args{
				ctx:     ctx,
				nudgeID: ksuid.New().String(),
			},
			wantErr: false,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.GetNudge(tt.args.ctx, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_GetAction(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	action := getTestAction()
	savedAction, err := fr.SaveAction(ctx, uid, fl, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		ctx      context.Context
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "valid case - action that exists",
			args: args{
				ctx:      ctx,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
		{
			name: "invalid case - action that does not exist",
			args: args{
				ctx:      ctx,
				actionID: ksuid.New().String(),
			},
			wantErr: false,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.GetAction(tt.args.ctx, tt.args.actionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_PublishFeedItem(t *testing.T) {

	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	type args struct {
		ctx  context.Context
		item *feed.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful, valid item",
			args: args{
				ctx:  ctx,
				item: testItem,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful, invalid item, will fail validation",
			args: args{
				ctx:  ctx,
				item: &feed.Item{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.PublishFeedItem(tt.args.ctx, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PublishFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_DeleteFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case, should delete",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent ID, should not fail",
			args: args{
				ctx:    ctx,
				itemID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.DeleteFeedItem(tt.args.ctx, tt.args.itemID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteFeedItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_ResolveFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feed.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantStatus: feed.StatusDone,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.ResolveFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.ResolveFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Status, tt.wantStatus) {
				t.Errorf("Feed.ResolveFeedItem() = %v, want %v", got.Status, tt.wantStatus)
			}
		})
	}
}

func TestFeed_UnresolveFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feed.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantStatus: feed.StatusPending,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.UnresolveFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.UnresolveFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Status, tt.wantStatus) {
				t.Errorf("Feed.ResolveFeedItem() = %v, want %v", got.Status, tt.wantStatus)
			}
		})
	}
}

func TestFeed_PinFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name           string
		args           args
		wantPersistent bool
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantPersistent: true,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.PinFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PinFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Persistent, tt.wantPersistent) {
				t.Errorf("Feed.PinFeedItem() = %v, want %v", got.Persistent, tt.wantPersistent)
			}
		})
	}
}

func TestFeed_UnpinFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name           string
		args           args
		wantPersistent bool
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantPersistent: false,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.UnpinFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.UnpinFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Persistent, tt.wantPersistent) {
				t.Errorf("Feed.PinFeedItem() = %v, want %v", got.Persistent, tt.wantPersistent)
			}
		})
	}
}

func TestFeed_HideFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feed.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantVisibility: feed.VisibilityHide,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.HideFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.HideFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Visibility, tt.wantVisibility) {
				t.Errorf("Feed.HideFeedItem() = %v, want %v", got.Visibility, tt.wantVisibility)
			}
		})
	}
}

func TestFeed_ShowFeedItem(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feed.FlavourConsumer

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	item, err := fr.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx    context.Context
		itemID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feed.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantVisibility: feed.VisibilityShow,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.ShowFeedItem(tt.args.ctx, tt.args.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.ShowFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Visibility, tt.wantVisibility) {
				t.Errorf("Feed.HideFeedItem() = %v, want %v", got.Visibility, tt.wantVisibility)
			}
		})
	}
}

func TestFeed_PublishNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()

	type args struct {
		ctx   context.Context
		nudge *feed.Nudge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				ctx:   ctx,
				nudge: nudge,
			},
			wantErr: false,
		},
		{
			name: "failure case - nil nudge",
			args: args{
				ctx:   ctx,
				nudge: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.PublishNudge(tt.args.ctx, tt.args.nudge)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PublishNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_DeleteNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := fe.PublishNudge(ctx, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
		{
			name: "another success case - should not fail even if IDs dont exist",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.DeleteNudge(tt.args.ctx, tt.args.nudgeID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteNudge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_ResolveNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := fe.PublishNudge(ctx, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feed.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantStatus: feed.StatusDone,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.ResolveNudge(tt.args.ctx, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.ResolveNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Status, tt.wantStatus) {
				t.Errorf("Feed.ResolveNudge() = %v, want %v", got.Status, tt.wantStatus)
			}
		})
	}
}

func TestFeed_UnresolveNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := fe.PublishNudge(ctx, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feed.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantStatus: feed.StatusPending,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.UnresolveNudge(tt.args.ctx, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.UnresolveNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Status, tt.wantStatus) {
				t.Errorf("Feed.ResolveNudge() = %v, want %v", got.Status, tt.wantStatus)
			}
		})
	}
}

func TestFeed_HideNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := fe.PublishNudge(ctx, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feed.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: feed.VisibilityHide,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.HideNudge(tt.args.ctx, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.HideNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Visibility, tt.wantVisibility) {
				t.Errorf("Feed.HideNudge() = %v, want %v", got.Visibility, tt.wantVisibility)
			}
		})
	}
}

func TestFeed_ShowNudge(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := fe.PublishNudge(ctx, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		nudgeID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feed.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: feed.VisibilityShow,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.ShowNudge(tt.args.ctx, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.ShowNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Visibility, tt.wantVisibility) {
				t.Errorf("Feed.HideNudge() = %v, want %v", got.Visibility, tt.wantVisibility)
			}
		})
	}
}

func TestFeed_PublishAction(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()

	type args struct {
		ctx    context.Context
		action *feed.Action
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:    ctx,
				action: &action,
			},
			wantErr: false,
		},
		{
			name: "invalid case - nil input",
			args: args{
				ctx:    ctx,
				action: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.PublishAction(tt.args.ctx, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PublishAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_DeleteAction(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()
	savedAction, err := fe.PublishAction(ctx, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		ctx      context.Context
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "action that exists",
			args: args{
				ctx:      ctx,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
		{
			name: "action that does not exist - should not error",
			args: args{
				ctx:      ctx,
				actionID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.DeleteAction(tt.args.ctx, tt.args.actionID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_PostMessage(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := fe.PublishFeedItem(ctx, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()

	type args struct {
		ctx     context.Context
		itemID  string
		message *feed.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful post",
			args: args{
				ctx:     ctx,
				itemID:  item.ID,
				message: &message,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful post - non existent item",
			args: args{
				ctx:     ctx,
				itemID:  ksuid.New().String(),
				message: &message,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful post - nil item",
			args: args{
				ctx:     ctx,
				itemID:  item.ID,
				message: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.PostMessage(tt.args.ctx, tt.args.itemID, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.PostMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_DeleteMessage(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()

	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := getTestItem()
	item, err := fe.PublishFeedItem(ctx, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()
	savedMessage, err := fe.PostMessage(ctx, item.ID, &message)
	assert.Nil(t, err)
	assert.NotNil(t, savedMessage)

	type args struct {
		ctx       context.Context
		itemID    string
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success - message that exists",
			args: args{
				ctx:       ctx,
				itemID:    item.ID,
				messageID: savedMessage.ID,
			},
			wantErr: false,
		},
		{
			name: "success - message that does not exist",
			args: args{
				ctx:       ctx,
				itemID:    item.ID,
				messageID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fe.DeleteMessage(tt.args.ctx, tt.args.itemID, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFeed_ProcessEvent(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := feed.FlavourConsumer
	uid := ksuid.New().String()
	fe, err := agg.GetThinFeed(ctx, uid, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *feed.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:   ctx,
				event: &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:   ctx,
				event: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.ProcessEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Feed.ProcessEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocument_ValidateAndUnmarshal(t *testing.T) {
	samplePDFBytes := []byte{
		123, 34, 105, 100, 34, 58, 34, 55, 49, 55, 51, 51, 51, 55, 50, 45, 50, 53, 50, 51, 45, 52, 51,
		100, 52, 45, 56, 102, 57, 97, 45, 53, 53, 53, 100, 51, 102, 52, 55, 54, 97, 52, 97, 34, 44, 34,
		98, 97, 115, 101, 54, 52, 34, 58, 34, 74, 86, 66, 69, 82, 105, 48, 120, 76, 106, 85, 75, 74, 98,
		88, 116, 114, 118, 115, 75, 78, 67, 65, 119, 73, 71, 57, 105, 97, 103, 111, 56, 80, 67, 65,
		118, 84, 71, 86, 117, 90, 51, 82, 111, 73, 68, 85, 103, 77, 67, 66, 83, 67, 105, 65, 103, 73,
		67, 57, 71, 97, 87, 120, 48, 90, 88, 73, 103, 76, 48, 90, 115, 89, 88, 82, 108, 82, 71, 86, 106,
		98, 50, 82, 108, 67, 106, 52, 43, 67, 110, 78, 48, 99, 109, 86, 104, 98, 81, 112, 52, 110, 68,
		78, 85, 77, 65, 66, 67, 88, 85, 77, 81, 112, 87, 100, 107, 111, 112, 67, 99, 121, 49, 88, 73, 70,
		99, 103, 70, 65, 68, 67, 119, 66, 70, 81, 75, 90, 87, 53, 107, 99, 51, 82, 121, 90, 87, 70, 116,
		67, 109, 86, 117, 90, 71, 57, 105, 97, 103, 111, 49, 73, 68, 65, 103, 98, 50, 74, 113, 67, 105,
		65, 103, 73, 68, 73, 51, 67, 109, 86, 117, 90, 71, 57, 105, 97, 103, 111, 122, 73, 68, 65, 103,
		98, 50, 74, 113, 67, 106, 119, 56, 67, 106, 52, 43, 67, 109, 86, 117, 90, 71, 57, 105, 97, 103,
		111, 121, 73, 68, 65, 103, 98, 50, 74, 113, 67, 106, 119, 56, 73, 67, 57, 85, 101, 88, 66, 108,
		73, 67, 57, 81, 89, 87, 100, 108, 73, 67, 85, 103, 77, 81, 111, 103, 73, 67, 65, 118, 85, 71, 70,
		121, 90, 87, 53, 48, 73, 68, 69, 103, 77, 67, 66, 83, 67, 105, 65, 103, 73, 67, 57, 78, 90, 87,
		82, 112, 89, 85, 74, 118, 101, 67, 66, 98, 73, 68, 65, 103, 77, 67, 65, 119, 76, 106, 73, 48,
		73, 68, 65, 117, 77, 106, 81, 103, 88, 81, 111, 103, 73, 67, 65, 118, 81, 50, 57, 117, 100, 71,
		86, 117, 100, 72, 77, 103, 78, 67, 65, 119, 73, 70, 73, 75, 73, 67, 65, 103, 76, 48, 100, 121,
		98, 51, 86, 119, 73, 68, 119, 56, 67, 105, 65, 103, 73, 67, 65, 103, 73, 67, 57, 85, 101, 88, 66,
		108, 73, 67, 57, 72, 99, 109, 57, 49, 99, 65, 111, 103, 73, 67, 65, 103, 73, 67, 65, 118, 85,
		121, 65, 118, 86, 72, 74, 104, 98, 110, 78, 119, 89, 88, 74, 108, 98, 109, 78, 53, 67, 105, 65,
		103, 73, 67, 65, 103, 73, 67, 57, 74, 73, 72, 82, 121, 100, 87, 85, 75, 73, 67, 65, 103, 73, 67,
		65, 103, 76, 48, 78, 84, 73, 67, 57, 69, 90, 88, 90, 112, 89, 50, 86, 83, 82, 48, 73, 75, 73, 67,
		65, 103, 80, 106, 52, 75, 73, 67, 65, 103, 76, 49, 74, 108, 99, 50, 57, 49, 99, 109, 78, 108, 99,
		121, 65, 122, 73, 68, 65, 103, 85, 103, 111, 43, 80, 103, 112, 108, 98, 109, 82, 118, 89, 109,
		111, 75, 77, 83, 65, 119, 73, 71, 57, 105, 97, 103, 111, 56, 80, 67, 65, 118, 86, 72, 108, 119,
		90, 83, 65, 118, 85, 71, 70, 110, 90, 88, 77, 75, 73, 67, 65, 103, 76, 48, 116, 112, 90, 72, 77,
		103, 87, 121, 65, 121, 73, 68, 65, 103, 85, 105, 66, 100, 67, 105, 65, 103, 73, 67, 57, 68, 98,
		51, 86, 117, 100, 67, 65, 120, 67, 106, 52, 43, 67, 109, 86, 117, 90, 71, 57, 105, 97, 103, 111,
		50, 73, 68, 65, 103, 98, 50, 74, 113, 67, 106, 119, 56, 73, 67, 57, 81, 99, 109, 57, 107, 100,
		87, 78, 108, 99, 105, 65, 111, 89, 50, 70, 112, 99, 109, 56, 103, 77, 83, 52, 120, 78, 105, 52,
		119, 73, 67, 104, 111, 100, 72, 82, 119, 99, 122, 111, 118, 76, 50, 78, 104, 97, 88, 74, 118,
		90, 51, 74, 104, 99, 71, 104, 112, 89, 51, 77, 117, 98, 51, 74, 110, 75, 83, 107, 75, 73, 67,
		65, 103, 76, 48, 78, 121, 90, 87, 70, 48, 97, 87, 57, 117, 82, 71, 70, 48, 90, 83, 65, 111, 82,
		68, 111, 121, 77, 68, 73, 119, 77, 84, 65, 122, 77, 68, 65, 52, 77, 68, 107, 119, 79, 67, 115,
		119, 77, 121, 99, 119, 77, 67, 107, 75, 80, 106, 52, 75, 90, 87, 53, 107, 98, 50, 74, 113, 67,
		106, 99, 103, 77, 67, 66, 118, 89, 109, 111, 75, 80, 68, 119, 103, 76, 49, 82, 53, 99, 71, 85,
		103, 76, 48, 78, 104, 100, 71, 70, 115, 98, 50, 99, 75, 73, 67, 65, 103, 76, 49, 66, 104, 90,
		50, 86, 122, 73, 68, 69, 103, 77, 67, 66, 83, 67, 106, 52, 43, 67, 109, 86, 117, 90, 71, 57,
		105, 97, 103, 112, 52, 99, 109, 86, 109, 67, 106, 65, 103, 79, 65, 111, 119, 77, 68, 65, 119,
		77, 68, 65, 119, 77, 68, 65, 119, 73, 68, 89, 49, 78, 84, 77, 49, 73, 71, 89, 103, 67, 106, 65,
		119, 77, 68, 65, 119, 77, 68, 65, 122, 79, 68, 69, 103, 77, 68, 65, 119, 77, 68, 65, 103, 98,
		105, 65, 75, 77, 68, 65, 119, 77, 68, 65, 119, 77, 68, 69, 50, 77, 83, 65, 119, 77, 68, 65, 119,
		77, 67, 66, 117, 73, 65, 111, 119, 77, 68, 65, 119, 77, 68, 65, 119, 77, 84, 81, 119, 73, 68,
		65, 119, 77, 68, 65, 119, 73, 71, 52, 103, 67, 106, 65, 119, 77, 68, 65, 119, 77, 68, 65, 119,
		77, 84, 85, 103, 77, 68, 65, 119, 77, 68, 65, 103, 98, 105, 65, 75, 77, 68, 65, 119, 77, 68, 65,
		119, 77, 68, 69, 120, 79, 83, 65, 119, 77, 68, 65, 119, 77, 67, 66, 117, 73, 65, 111, 119, 77,
		68, 65, 119, 77, 68, 65, 119, 78, 68, 81, 50, 73, 68, 65, 119, 77, 68, 65, 119, 73, 71, 52, 103,
		67, 106, 65, 119, 77, 68, 65, 119, 77, 68, 65, 49, 78, 106, 73, 103, 77, 68, 65, 119, 77, 68,
		65, 103, 98, 105, 65, 75, 100, 72, 74, 104, 97, 87, 120, 108, 99, 103, 111, 56, 80, 67, 65,
		118, 85, 50, 108, 54, 90, 83, 65, 52, 67, 105, 65, 103, 73, 67, 57, 83, 98, 50, 57, 48, 73, 68,
		99, 103, 77, 67, 66, 83, 67, 105, 65, 103, 73, 67, 57, 74, 98, 109, 90, 118, 73, 68, 89, 103,
		77, 67, 66, 83, 67, 106, 52, 43, 67, 110, 78, 48, 89, 88, 74, 48, 101, 72, 74, 108, 90, 103,
		111, 50, 77, 84, 81, 75, 74, 83, 86, 70, 84, 48, 89, 75, 34, 125,
	}

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid PDF document",
			args: args{
				b: samplePDFBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid PDF document",
			args: args{
				b: []byte{1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &feed.Document{}
			if err := doc.ValidateAndUnmarshal(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Document.ValidateAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocument_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID     string
		Base64 string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid PDF document",
			fields: fields{
				ID:     ksuid.New().String(),
				Base64: base64PDFSample,
			},
			wantErr: false,
		},
		{
			name: "invalid PDF document",
			fields: fields{
				ID:     ksuid.New().String(),
				Base64: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &feed.Document{
				ID:     tt.fields.ID,
				Base64: tt.fields.Base64,
			}
			got, err := doc.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Document.ValidateAndMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func getTestDocument() feed.Document {
	return feed.Document{
		ID:     ksuid.New().String(),
		Base64: base64PDFSample,
	}
}

func getTestMessage() feed.Message {
	return feed.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func testItem() *feed.Item {
	return &feed.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         getTextExpiry(),
		Persistent:     true,
		Status:         feed.StatusPending,
		Visibility:     feed.VisibilityShow,
		Icon:           getTestImage(),
		Author:         ksuid.New().String(),
		Tagline:        ksuid.New().String(),
		Label:          ksuid.New().String(),
		Timestamp:      time.Now(),
		Summary:        ksuid.New().String(),
		Text:           ksuid.New().String(),
		Images: []feed.Image{
			getTestImage(),
		},
		Videos: []feed.Video{
			getTestVideo(),
		},
		Actions: []feed.Action{
			getTestAction(),
		},
		Conversations: []feed.Message{
			getTestMessage(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		Documents: []feed.Document{
			getTestDocument(),
		},
		NotificationChannels: []feed.Channel{
			feed.ChannelEmail,
			feed.ChannelFcm,
			feed.ChannelSms,
			feed.ChannelWhatsapp,
		},
	}
}

func getTextExpiry() time.Time {
	return time.Now().Add(time.Hour * 24000)
}

func getTestImage() feed.Image {
	return feed.Image{
		ID:     ksuid.New().String(),
		Base64: base64PNGSample,
	}
}

func getTestVideo() feed.Video {
	return feed.Video{
		ID:  ksuid.New().String(),
		URL: sampleVideoURL,
	}
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
