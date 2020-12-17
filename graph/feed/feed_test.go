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
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
)

const (
	sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"

	intMax = 9007199254740990
)

func getEmptyJson(t *testing.T) []byte {
	emptyJSONBytes, err := json.Marshal(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, emptyJSONBytes)
	return emptyJSONBytes
}

func getTestItem() base.Item {
	return base.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       base.TextTypePlain,
		Links: []base.Link{
			base.GetYoutubeVideoLink(sampleVideoURL, "title", "description", base.BlankImageURL),
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
				AllowAnonymous: true,
			},
		},
		Conversations: []base.Message{
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
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func TestNewInMemoryFeed(t *testing.T) {
	feeds := getTestFeedAggregate(t)
	assert.NotNil(t, feeds)
}

func TestMessage_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Message{
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
			msg := &base.Message{}
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

	validElement := base.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       base.TextTypeMarkdown,
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
			},
		},
		Conversations: []base.Message{
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
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
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
		Status         base.Status
		Visibility     base.Visibility
		Icon           base.Link
		Author         string
		Tagline        string
		Label          string
		Timestamp      time.Time
		Summary        string
		Text           string
		Links          []base.Link
		Actions        []base.Action
		Conversations  []base.Message
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
			it := &base.Item{
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
				Links:          tt.fields.Links,
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

	validElement := base.Nudge{
		ID:             "nudge-1",
		SequenceNumber: 1,
		Visibility:     base.VisibilityShow,
		Status:         base.StatusPending,
		Title:          "Update your profile!",
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Text: "An up to date profile will help us serve you better!",
		Actions: []base.Action{
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
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
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Visibility     base.Visibility
		Status         base.Status
		Title          string
		Text           string
		Links          []base.Link
		Actions        []base.Action
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
			nu := &base.Nudge{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Visibility:     tt.fields.Visibility,
				Status:         tt.fields.Status,
				Title:          tt.fields.Title,
				Links:          tt.fields.Links,
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

	validElement := base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Name:           "ACTION_NAME",
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		ActionType:     base.ActionTypeSecondary,
		Handling:       base.HandlingFullPage,
		AllowAnonymous: false,
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Name           string
		Icon           base.Link
		ActionType     base.ActionType
		Handling       base.Handling
		AllowAnonymous bool
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
			ac := &base.Action{}
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
	anonymous := false
	emptyJSONBytes := getEmptyJson(t)
	validElement := feed.Feed{
		UID:            "a-uid",
		IsAnonymous:    &anonymous,
		SequenceNumber: int(time.Now().Unix()),
		Flavour:        base.FlavourConsumer,
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
				AllowAnonymous: false,
			},
		},
		Nudges: []base.Nudge{
			{
				ID:             "nudge-1",
				SequenceNumber: 1,
				Visibility:     base.VisibilityShow,
				Status:         base.StatusPending,
				Title:          "Update your profile!",
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []base.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
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
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
				},
			},
		},
		Items: []base.Item{
			{
				ID:             "item-1",
				SequenceNumber: 1,
				Expiry:         time.Now(),
				Persistent:     true,
				Status:         base.StatusPending,
				Visibility:     base.VisibilityShow,
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Author:    "Bot 1",
				Tagline:   "Bot speaks...",
				Label:     "DRUGS",
				Timestamp: time.Now(),
				Summary:   "I am a bot...",
				Text:      "This bot can speak",
				TextType:  base.TextTypeMarkdown,
				Actions: []base.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypeSecondary,
						Handling:       base.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Conversations: []base.Message{
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
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
				},
			},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		UID         string
		IsAnonymous *bool
		Flavour     base.Flavour
		Actions     []base.Action
		Items       []base.Item
		Nudges      []base.Nudge
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
				UID:         tt.fields.UID,
				IsAnonymous: tt.fields.IsAnonymous,
				Flavour:     tt.fields.Flavour,
				Actions:     tt.fields.Actions,
				Items:       tt.fields.Items,
				Nudges:      tt.fields.Nudges,
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
		UID            string
		IsAnonymous    *bool
		SequenceNumber int
		Flavour        base.Flavour
		Actions        []base.Action
		Items          []base.Item
		Nudges         []base.Nudge
	}
	anonymous := false
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid feed",
			fields: fields{
				UID:            "a-uid",
				IsAnonymous:    &anonymous,
				SequenceNumber: int(time.Now().Unix()),
				Flavour:        base.FlavourPro,
				Actions: []base.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypeSecondary,
						Handling:       base.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Nudges: []base.Nudge{
					{
						ID:             "nudge-1",
						SequenceNumber: 1,
						Visibility:     base.VisibilityShow,
						Status:         base.StatusPending,
						Title:          "Update your profile!",
						Links: []base.Link{
							base.GetPNGImageLink(
								base.LogoURL, "title", "description", base.BlankImageURL),
						},
						Text: "Help us serve you better!",
						Actions: []base.Action{
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								Icon: base.GetPNGImageLink(
									base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypePrimary,
								Handling:       base.HandlingInline,
								AllowAnonymous: false,
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
						NotificationChannels: []base.Channel{
							base.ChannelFcm,
							base.ChannelEmail,
							base.ChannelSms,
							base.ChannelWhatsapp,
						},
					},
				},
				Items: []base.Item{
					{
						ID:             "item-1",
						SequenceNumber: 1,
						Expiry:         time.Now(),
						Persistent:     true,
						Status:         base.StatusPending,
						Visibility:     base.VisibilityShow,
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						Links: []base.Link{
							base.GetPNGImageLink(
								base.LogoURL, "title", "description", base.BlankImageURL),
						},
						Author:    "Bot 1",
						Tagline:   "Bot speaks...",
						Label:     "DRUGS",
						Timestamp: time.Now(),
						Summary:   "I am a bot...",
						Text:      "This bot can speak",
						TextType:  base.TextTypeMarkdown,
						Actions: []base.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           "ACTION_NAME",
								Icon: base.GetPNGImageLink(
									base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypeSecondary,
								Handling:       base.HandlingFullPage,
								AllowAnonymous: false,
							},
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								Icon: base.GetPNGImageLink(
									base.LogoURL, "title", "description", base.BlankImageURL),
								ActionType:     base.ActionTypePrimary,
								Handling:       base.HandlingInline,
								AllowAnonymous: false,
							},
						},
						Conversations: []base.Message{
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
						NotificationChannels: []base.Channel{
							base.ChannelFcm,
							base.ChannelEmail,
							base.ChannelSms,
							base.ChannelWhatsapp,
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
				UID:            tt.fields.UID,
				IsAnonymous:    tt.fields.IsAnonymous,
				SequenceNumber: tt.fields.SequenceNumber,
				Flavour:        tt.fields.Flavour,
				Actions:        tt.fields.Actions,
				Items:          tt.fields.Items,
				Nudges:         tt.fields.Nudges,
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
		Icon           base.Link
		ActionType     base.ActionType
		Handling       base.Handling
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
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType: base.ActionTypePrimary,
				Handling:   base.HandlingInline,
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
			ac := &base.Action{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Name:           tt.fields.Name,
				Icon:           tt.fields.Icon,
				ActionType:     tt.fields.ActionType,
				Handling:       tt.fields.Handling,
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
		Visibility           base.Visibility
		Status               base.Status
		Title                string
		Links                []base.Link
		Text                 string
		Actions              []base.Action
		Groups               []string
		Users                []string
		NotificationChannels []base.Channel
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
				Visibility:     base.VisibilityShow,
				Status:         base.StatusPending,
				Title:          "Update your profile!",
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []base.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
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
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
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
			nu := &base.Nudge{
				ID:                   tt.fields.ID,
				SequenceNumber:       tt.fields.SequenceNumber,
				Visibility:           tt.fields.Visibility,
				Status:               tt.fields.Status,
				Title:                tt.fields.Title,
				Links:                tt.fields.Links,
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
		Status               base.Status
		Visibility           base.Visibility
		Icon                 base.Link
		Author               string
		Tagline              string
		Label                string
		Timestamp            time.Time
		Summary              string
		Text                 string
		TextType             base.TextType
		Links                []base.Link
		Actions              []base.Action
		Conversations        []base.Message
		Users                []string
		Groups               []string
		NotificationChannels []base.Channel
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
				Status:         base.StatusPending,
				Visibility:     base.VisibilityShow,
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				Author:    "Bot 1",
				Tagline:   "Bot speaks...",
				Label:     "DRUGS",
				Timestamp: time.Now(),
				Summary:   "I am a bot...",
				Text:      "This bot can speak",
				TextType:  base.TextTypeMarkdown,
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Actions: []base.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypeSecondary,
						Handling:       base.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Conversations: []base.Message{
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
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
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
			it := &base.Item{
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
				TextType:             tt.fields.TextType,
				Links:                tt.fields.Links,
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
			msg := &base.Message{
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

func TestContext_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Context{
		UserID:         "uid-1",
		Flavour:        base.FlavourConsumer,
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
			ct := &base.Context{
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
		Flavour        base.Flavour
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
				Flavour:        base.FlavourConsumer,
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
			ct := &base.Context{
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

	validElement := base.Payload{
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
			pl := &base.Payload{
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
			pl := &base.Payload{
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

	validElement := base.Event{
		ID:   "event-1",
		Name: "THIS_EVENT",
		Context: base.Context{
			UserID:         "user-1",
			Flavour:        base.FlavourConsumer,
			OrganizationID: "org-1",
			LocationID:     "loc-1",
			Timestamp:      time.Now(),
		},
		Payload: base.Payload{
			Data: map[string]interface{}{"a": 1},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		ID      string
		Name    string
		Context base.Context
		Payload base.Payload
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
			ev := &base.Event{
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
		Context base.Context
		Payload base.Payload
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
				Context: base.Context{
					UserID:         "user-1",
					Flavour:        base.FlavourConsumer,
					OrganizationID: "org-1",
					LocationID:     "loc-1",
					Timestamp:      time.Now(),
				},
				Payload: base.Payload{
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
			ev := &base.Event{
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	flavour := base.FlavourConsumer
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	fr, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	testItem := testItem()

	type args struct {
		ctx  context.Context
		item *base.Item
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
				item: &base.Item{},
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantStatus base.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantStatus: base.StatusDone,
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantStatus base.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantStatus: base.StatusPending,
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantVisibility: base.VisibilityHide,
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
	fl := base.FlavourConsumer

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				itemID: item.ID,
			},
			wantVisibility: base.VisibilityShow,
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()

	type args struct {
		ctx   context.Context
		nudge *base.Nudge
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantStatus base.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantStatus: base.StatusDone,
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantStatus base.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantStatus: base.StatusPending,
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: base.VisibilityHide,
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		wantVisibility base.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: base.VisibilityShow,
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()

	type args struct {
		ctx    context.Context
		action *base.Action
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
		message *base.Message
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
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
	fl := base.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	event := getTestEvent()

	type args struct {
		ctx   context.Context
		event *base.Event
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

func getTestMessage() base.Message {
	return base.Message{
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

func testItem() *base.Item {
	return &base.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         getTextExpiry(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon: base.GetPNGImageLink(
			base.LogoURL, "title", "description", base.BlankImageURL),
		Author:    ksuid.New().String(),
		Tagline:   ksuid.New().String(),
		Label:     ksuid.New().String(),
		Timestamp: time.Now(),
		Summary:   ksuid.New().String(),
		Text:      ksuid.New().String(),
		TextType:  base.TextTypeMarkdown,
		Links: []base.Link{
			base.GetPNGImageLink(
				base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Actions: []base.Action{
			getTestAction(),
		},
		Conversations: []base.Message{
			getTestMessage(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []base.Channel{
			base.ChannelEmail,
			base.ChannelFcm,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func getTextExpiry() time.Time {
	return time.Now().Add(time.Hour * 24000)
}

func getTestEvent() base.Event {
	return base.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: base.Context{
			UserID:         ksuid.New().String(),
			Flavour:        base.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() base.Action {
	return base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon: base.GetPNGImageLink(
			base.LogoURL, "title", "description", base.BlankImageURL),
		ActionType: base.ActionTypePrimary,
		Handling:   base.HandlingFullPage,
	}
}

func testNudge() *base.Nudge {
	return &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Text: ksuid.New().String(),
		Actions: []base.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []base.Channel{
			base.ChannelEmail,
			base.ChannelFcm,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func TestLink_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)
	validLink := base.Link{
		ID:          ksuid.New().String(),
		URL:         sampleVideoURL,
		LinkType:    base.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   base.BlankImageURL,
	}
	validLinkJSONBytes, err := json.Marshal(validLink)
	assert.Nil(t, err)
	assert.NotNil(t, validLinkJSONBytes)
	assert.Greater(t, len(validLinkJSONBytes), 3)

	invalidVideoLink := base.Link{
		ID:          ksuid.New().String(),
		URL:         "www.example.com/not_a_youtube_video",
		LinkType:    base.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   base.BlankImageURL,
	}
	invalidLinkJSONBytes, err := json.Marshal(invalidVideoLink)
	assert.Nil(t, err)
	assert.NotNil(t, invalidLinkJSONBytes)
	assert.Greater(t, len(invalidLinkJSONBytes), 3)

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid link",
			args: args{
				b: validLinkJSONBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid link",
			args: args{
				b: invalidLinkJSONBytes,
			},
			wantErr: true,
		},
		{
			name: "empty JSON - invalid",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &base.Link{}
			if err := l.ValidateAndUnmarshal(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Link.ValidateAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLink_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID          string
		URL         string
		Type        base.LinkType
		Title       string
		Description string
		Thumbnail   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid link",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         sampleVideoURL,
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "not a valid URL",
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid YouTube URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_video",
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PNG URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_png",
				Type:        base.LinkTypePngImage,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PDF URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_pdf",
				Type:        base.LinkTypePdfDocument,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &base.Link{
				ID:          tt.fields.ID,
				URL:         tt.fields.URL,
				LinkType:    tt.fields.Type,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Thumbnail:   tt.fields.Thumbnail,
			}
			got, err := l.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.ValidateAndMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestFeed_Labels(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: ctx,
			},
			want:    []string{feed.DefaultLabel},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.Labels(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.Labels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Feed.Labels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeed_SaveLabel(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		ctx   context.Context
		label string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:   ctx,
				label: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.SaveLabel(tt.args.ctx, tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("Feed.SaveLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_UnreadPersistentItems(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: ctx,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fe.UnreadPersistentItems(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.UnreadPersistentItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Feed.UnreadPersistentItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeed_UpdateUnreadPersistentItemsCount(t *testing.T) {
	ctx := context.Background()
	agg := getFeedAggregate(t)
	fl := base.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fe.UpdateUnreadPersistentItemsCount(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Feed.UpdateUnreadPersistentItemsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
