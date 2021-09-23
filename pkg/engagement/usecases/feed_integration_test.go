package usecases_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/usecases"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	libInfra "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	libFeed "github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"

	IntMax = 9007199254740990
)

func InitializeTestNewFeed(ctx context.Context) (*usecases.FeedImpl, libInfra.Interactor, error) {
	infra := libInfra.NewInteractor()
	libUsc := libFeed.NewFeed(infra)
	lib := usecases.NewFeed(infra, libUsc)
	return lib, infra, nil
}

func getTestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getInvalidTestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             "",
		SequenceNumber: 0,
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestSequenceNumber() int {
	return rand.Intn(IntMax)
}

func testItem() *feedlib.Item {
	return &feedlib.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         getTextExpiry(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon: feedlib.GetPNGImageLink(
			feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		Author:    ksuid.New().String(),
		Tagline:   ksuid.New().String(),
		Label:     ksuid.New().String(),
		Timestamp: time.Now(),
		Summary:   ksuid.New().String(),
		Text:      ksuid.New().String(),
		TextType:  feedlib.TextTypeMarkdown,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(
				feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		},
		Actions: []feedlib.Action{
			getTestAction(),
		},
		Conversations: []feedlib.Message{
			getTestMessage(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelEmail,
			feedlib.ChannelFcm,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func getTextExpiry() time.Time {
	return time.Now().Add(time.Hour * 24000)
}

func getTestEvent() feedlib.Event {
	return feedlib.Event{
		ID:   ksuid.New().String(),
		Name: "TEST_EVENT",
		Context: feedlib.Context{
			UserID:         ksuid.New().String(),
			Flavour:        feedlib.FlavourConsumer,
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getInvalidTestEvent() feedlib.Event {
	return feedlib.Event{
		ID:   "",
		Name: "TEST_EVENT",
		Context: feedlib.Context{
			UserID:         "",
			Flavour:        "",
			OrganizationID: ksuid.New().String(),
			LocationID:     ksuid.New().String(),
			Timestamp:      time.Now(),
		},
	}
}

func getTestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon: feedlib.GetPNGImageLink(
			feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		ActionType: feedlib.ActionTypePrimary,
		Handling:   feedlib.HandlingFullPage,
	}
}

func testNudge() *feedlib.Nudge {
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			getTestAction(),
		},
		Users: []string{
			ksuid.New().String(),
		},
		Groups: []string{
			ksuid.New().String(),
		},
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelEmail,
			feedlib.ChannelFcm,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func getEmptyJson(t *testing.T) []byte {
	emptyJSONBytes, err := json.Marshal(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, emptyJSONBytes)
	return emptyJSONBytes
}

func getTestItem() feedlib.Item {
	return feedlib.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feedlib.TextTypePlain,
		Links: []feedlib.Link{
			feedlib.GetYoutubeVideoLink(sampleVideoURL, "title", "description", feedlib.BlankImageURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
				AllowAnonymous: false,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
				AllowAnonymous: true,
			},
		},
		Conversations: []feedlib.Message{
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
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelFcm,
			feedlib.ChannelEmail,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
}

func TestMessage_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := feedlib.Message{
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
			msg := &feedlib.Message{}
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

	validElement := feedlib.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feedlib.TextTypeMarkdown,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
			},
		},
		Conversations: []feedlib.Message{
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
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelFcm,
			feedlib.ChannelEmail,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
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
		Status         feedlib.Status
		Visibility     feedlib.Visibility
		Icon           feedlib.Link
		Author         string
		Tagline        string
		Label          string
		Timestamp      time.Time
		Summary        string
		Text           string
		Links          []feedlib.Link
		Actions        []feedlib.Action
		Conversations  []feedlib.Message
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
			it := &feedlib.Item{
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

	validElement := feedlib.Nudge{
		ID:             "nudge-1",
		SequenceNumber: 1,
		Visibility:     feedlib.VisibilityShow,
		Status:         feedlib.StatusPending,
		Title:          "Update your profile!",
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		},
		Text: "An up to date profile will help us serve you better!",
		Actions: []feedlib.Action{
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
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
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelFcm,
			feedlib.ChannelEmail,
			feedlib.ChannelSms,
			feedlib.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Visibility     feedlib.Visibility
		Status         feedlib.Status
		Title          string
		Text           string
		Links          []feedlib.Link
		Actions        []feedlib.Action
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
			nu := &feedlib.Nudge{
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

	validElement := feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Name:           "ACTION_NAME",
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
		ActionType:     feedlib.ActionTypeSecondary,
		Handling:       feedlib.HandlingFullPage,
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
		Icon           feedlib.Link
		ActionType     feedlib.ActionType
		Handling       feedlib.Handling
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
			ac := &feedlib.Action{}
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
	validElement := domain.Feed{
		UID:            "a-uid",
		IsAnonymous:    &anonymous,
		SequenceNumber: int(time.Now().Unix()),
		Flavour:        feedlib.FlavourConsumer,
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon: feedlib.GetPNGImageLink(
					feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
				AllowAnonymous: false,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon: feedlib.GetPNGImageLink(
					feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
				AllowAnonymous: false,
			},
		},
		Nudges: []feedlib.Nudge{
			{
				ID:             "nudge-1",
				SequenceNumber: 1,
				Visibility:     feedlib.VisibilityShow,
				Status:         feedlib.StatusPending,
				Title:          "Update your profile!",
				Links: []feedlib.Link{
					feedlib.GetPNGImageLink(
						feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []feedlib.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypePrimary,
						Handling:       feedlib.HandlingInline,
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
				NotificationChannels: []feedlib.Channel{
					feedlib.ChannelFcm,
					feedlib.ChannelEmail,
					feedlib.ChannelSms,
					feedlib.ChannelWhatsapp,
				},
			},
		},
		Items: []feedlib.Item{
			{
				ID:             "item-1",
				SequenceNumber: 1,
				Expiry:         time.Now(),
				Persistent:     true,
				Status:         feedlib.StatusPending,
				Visibility:     feedlib.VisibilityShow,
				Icon: feedlib.GetPNGImageLink(
					feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				Links: []feedlib.Link{
					feedlib.GetPNGImageLink(
						feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				},
				Author:    "Bot 1",
				Tagline:   "Bot speaks...",
				Label:     "DRUGS",
				Timestamp: time.Now(),
				Summary:   "I am a bot...",
				Text:      "This bot can speak",
				TextType:  feedlib.TextTypeMarkdown,
				Actions: []feedlib.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypeSecondary,
						Handling:       feedlib.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypePrimary,
						Handling:       feedlib.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Conversations: []feedlib.Message{
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
				NotificationChannels: []feedlib.Channel{
					feedlib.ChannelFcm,
					feedlib.ChannelEmail,
					feedlib.ChannelSms,
					feedlib.ChannelWhatsapp,
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
		Flavour     feedlib.Flavour
		Actions     []feedlib.Action
		Items       []feedlib.Item
		Nudges      []feedlib.Nudge
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
			fe := &domain.Feed{
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
		Flavour        feedlib.Flavour
		Actions        []feedlib.Action
		Items          []feedlib.Item
		Nudges         []feedlib.Nudge
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
				Flavour:        feedlib.FlavourPro,
				Actions: []feedlib.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypeSecondary,
						Handling:       feedlib.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypePrimary,
						Handling:       feedlib.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Nudges: []feedlib.Nudge{
					{
						ID:             "nudge-1",
						SequenceNumber: 1,
						Visibility:     feedlib.VisibilityShow,
						Status:         feedlib.StatusPending,
						Title:          "Update your profile!",
						Links: []feedlib.Link{
							feedlib.GetPNGImageLink(
								feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						},
						Text: "Help us serve you better!",
						Actions: []feedlib.Action{
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								Icon: feedlib.GetPNGImageLink(
									feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
								ActionType:     feedlib.ActionTypePrimary,
								Handling:       feedlib.HandlingInline,
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
						NotificationChannels: []feedlib.Channel{
							feedlib.ChannelFcm,
							feedlib.ChannelEmail,
							feedlib.ChannelSms,
							feedlib.ChannelWhatsapp,
						},
					},
				},
				Items: []feedlib.Item{
					{
						ID:             "item-1",
						SequenceNumber: 1,
						Expiry:         time.Now(),
						Persistent:     true,
						Status:         feedlib.StatusPending,
						Visibility:     feedlib.VisibilityShow,
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						Links: []feedlib.Link{
							feedlib.GetPNGImageLink(
								feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						},
						Author:    "Bot 1",
						Tagline:   "Bot speaks...",
						Label:     "DRUGS",
						Timestamp: time.Now(),
						Summary:   "I am a bot...",
						Text:      "This bot can speak",
						TextType:  feedlib.TextTypeMarkdown,
						Actions: []feedlib.Action{
							{
								ID:             ksuid.New().String(),
								SequenceNumber: 1,
								Name:           "ACTION_NAME",
								Icon: feedlib.GetPNGImageLink(
									feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
								ActionType:     feedlib.ActionTypeSecondary,
								Handling:       feedlib.HandlingFullPage,
								AllowAnonymous: false,
							},
							{
								ID:             "action-1",
								SequenceNumber: 1,
								Name:           "First action",
								Icon: feedlib.GetPNGImageLink(
									feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
								ActionType:     feedlib.ActionTypePrimary,
								Handling:       feedlib.HandlingInline,
								AllowAnonymous: false,
							},
						},
						Conversations: []feedlib.Message{
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
						NotificationChannels: []feedlib.Channel{
							feedlib.ChannelFcm,
							feedlib.ChannelEmail,
							feedlib.ChannelSms,
							feedlib.ChannelWhatsapp,
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
			fe := &domain.Feed{
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
		Icon           feedlib.Link
		ActionType     feedlib.ActionType
		Handling       feedlib.Handling
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
				Icon: feedlib.GetPNGImageLink(
					feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				ActionType: feedlib.ActionTypePrimary,
				Handling:   feedlib.HandlingInline,
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
			ac := &feedlib.Action{
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
		Visibility           feedlib.Visibility
		Status               feedlib.Status
		Title                string
		Links                []feedlib.Link
		Text                 string
		Actions              []feedlib.Action
		Groups               []string
		Users                []string
		NotificationChannels []feedlib.Channel
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
				Visibility:     feedlib.VisibilityShow,
				Status:         feedlib.StatusPending,
				Title:          "Update your profile!",
				Links: []feedlib.Link{
					feedlib.GetPNGImageLink(
						feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []feedlib.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypePrimary,
						Handling:       feedlib.HandlingInline,
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
				NotificationChannels: []feedlib.Channel{
					feedlib.ChannelFcm,
					feedlib.ChannelEmail,
					feedlib.ChannelSms,
					feedlib.ChannelWhatsapp,
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
			nu := &feedlib.Nudge{
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
		Status               feedlib.Status
		Visibility           feedlib.Visibility
		Icon                 feedlib.Link
		Author               string
		Tagline              string
		Label                string
		Timestamp            time.Time
		Summary              string
		Text                 string
		TextType             feedlib.TextType
		Links                []feedlib.Link
		Actions              []feedlib.Action
		Conversations        []feedlib.Message
		Users                []string
		Groups               []string
		NotificationChannels []feedlib.Channel
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
				Status:         feedlib.StatusPending,
				Visibility:     feedlib.VisibilityShow,
				Icon: feedlib.GetPNGImageLink(
					feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				Author:    "Bot 1",
				Tagline:   "Bot speaks...",
				Label:     "DRUGS",
				Timestamp: time.Now(),
				Summary:   "I am a bot...",
				Text:      "This bot can speak",
				TextType:  feedlib.TextTypeMarkdown,
				Links: []feedlib.Link{
					feedlib.GetPNGImageLink(
						feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
				},
				Actions: []feedlib.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypeSecondary,
						Handling:       feedlib.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: feedlib.GetPNGImageLink(
							feedlib.LogoURL, "title", "description", feedlib.BlankImageURL),
						ActionType:     feedlib.ActionTypePrimary,
						Handling:       feedlib.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Conversations: []feedlib.Message{
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
				NotificationChannels: []feedlib.Channel{
					feedlib.ChannelFcm,
					feedlib.ChannelEmail,
					feedlib.ChannelSms,
					feedlib.ChannelWhatsapp,
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
			it := &feedlib.Item{
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
			msg := &feedlib.Message{
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

	validElement := feedlib.Context{
		UserID:         "uid-1",
		Flavour:        feedlib.FlavourConsumer,
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
			ct := &feedlib.Context{
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
		Flavour        feedlib.Flavour
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
				Flavour:        feedlib.FlavourConsumer,
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
			ct := &feedlib.Context{
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

	validElement := feedlib.Payload{
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
			pl := &feedlib.Payload{
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
			pl := &feedlib.Payload{
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

	validElement := feedlib.Event{
		ID:   "event-1",
		Name: "THIS_EVENT",
		Context: feedlib.Context{
			UserID:         "user-1",
			Flavour:        feedlib.FlavourConsumer,
			OrganizationID: "org-1",
			LocationID:     "loc-1",
			Timestamp:      time.Now(),
		},
		Payload: feedlib.Payload{
			Data: map[string]interface{}{"a": 1},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		ID      string
		Name    string
		Context feedlib.Context
		Payload feedlib.Payload
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
			ev := &feedlib.Event{
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
		Context feedlib.Context
		Payload feedlib.Payload
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
				Context: feedlib.Context{
					UserID:         "user-1",
					Flavour:        feedlib.FlavourConsumer,
					OrganizationID: "org-1",
					LocationID:     "loc-1",
					Timestamp:      time.Now(),
				},
				Payload: feedlib.Payload{
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
			ev := &feedlib.Event{
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *feedlib.Item
	}{
		{
			name: "valid case - item that exists",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantErr: false,
			want:    item,
		},
		{
			name: "invalid case - item that does not exist",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.GetFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && err != nil {
				assert.NotNil(t, got.Expiry)
				expirytext := getTextExpiry()
				got.Expiry = expirytext
				tt.want.Expiry = expirytext
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Feed.GetFeedItem() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestFeed_GetNudge(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	flavour := feedlib.FlavourConsumer
	nudge := testNudge()

	savedNudge, err := infra.SaveNudge(ctx, uid, flavour, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *feedlib.Nudge
	}{
		{
			name: "valid case - nudge that exists",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
			want:    savedNudge,
		},
		{
			name: "invalid case - nudge that does not exist",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.GetNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Feed.GetNudge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeedUseCaseImpl_GetDefaultNudgeByTitle(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	if err != nil {
		t.Errorf("failed to initialize a new feed")
		return
	}
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	nudge := testNudge()
	savedNudge, err := infra.SaveNudge(ctx, uid, fl, nudge)
	if err != nil {
		t.Errorf("failed to save nudge")
		return
	}
	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		title   string
	}
	tests := []struct {
		name    string
		args    args
		want    *feedlib.Nudge
		wantErr bool
	}{
		{
			name: "valid case - Existing nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				title:   savedNudge.Title,
			},
			want:    savedNudge,
			wantErr: false,
		},
		{
			name: "invalid case - Non-existent nudge",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				title:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.GetDefaultNudgeByTitle(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.title,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FeedUseCaseImpl.GetDefaultNudgeByTitle() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FeedUseCaseImpl.GetDefaultNudgeByTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestFeed_GetAction(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()
	savedAction, err := infra.SaveAction(ctx, uid, fl, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		ctx      context.Context
		uid      string
		flavour  feedlib.Flavour
		actionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *feedlib.Action
	}{
		{
			name: "valid case - action that exists",
			args: args{
				ctx:      ctx,
				uid:      uid,
				flavour:  fl,
				actionID: savedAction.ID,
			},
			wantErr: false,
			want:    savedAction,
		},
		{
			name: "invalid case - action that does not exist",
			args: args{
				ctx:      ctx,
				uid:      uid,
				flavour:  fl,
				actionID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.GetAction(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.actionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.GetAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Feed.GetAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeed_PublishFeedItem(t *testing.T) {

	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		item    *feedlib.Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful, valid item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				item:    testItem,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful, invalid item, will fail validation",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				item:    &feedlib.Item{},
			},
			wantErr: true,
		},
		{
			name: "invalid case: nil item, will fail to publish feed item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				item:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.PublishFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.item)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case, should delete",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantErr: false,
		},
		{
			name: "non existent ID, should not fail",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.DeleteFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteFeedItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_ResolveFeedItem(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantStatus: feedlib.StatusDone,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.ResolveFeedItem(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantStatus: feedlib.StatusPending,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.UnresolveFeedItem(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
			)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
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
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantPersistent: true,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.PinFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
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
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantPersistent: false,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.UnpinFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: flavour,
				itemID:  item.ID,
			},
			wantVisibility: feedlib.VisibilityHide,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.HideFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, infra, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer

	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := infra.SaveFeedItem(ctx, uid, flavour, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
			},
			wantVisibility: feedlib.VisibilityShow,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.ShowFeedItem(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.itemID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudge   *feedlib.Nudge
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
				uid:     uid,
				flavour: fl,
				nudge:   nudge,
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
			got, err := agg.LibUsecases.PublishNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudge)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := agg.LibUsecases.PublishNudge(ctx, uid, fl, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
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
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
		{
			name: "another success case - should not fail even if IDs dont exist",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.DeleteNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteNudge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_ResolveNudge(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := agg.LibUsecases.PublishNudge(ctx, uid, fl, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantStatus: feedlib.StatusDone,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.ResolveNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := agg.LibUsecases.PublishNudge(ctx, uid, fl, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus feedlib.Status
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantStatus: feedlib.StatusPending,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.UnresolveNudge(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.nudgeID,
			)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := agg.LibUsecases.PublishNudge(ctx, uid, fl, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: feedlib.VisibilityHide,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.HideNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	nudge := testNudge()
	savedNudge, err := agg.LibUsecases.PublishNudge(ctx, uid, fl, nudge)
	assert.Nil(t, err)
	assert.NotNil(t, savedNudge)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		nudgeID string
	}
	tests := []struct {
		name           string
		args           args
		wantVisibility feedlib.Visibility
		wantErr        bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				nudgeID: savedNudge.ID,
			},
			wantVisibility: feedlib.VisibilityShow,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.ShowNudge(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.nudgeID)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		action  *feedlib.Action
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				action:  &action,
			},
			wantErr: false,
		},
		{
			name: "invalid case - nil input",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				action:  nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.PublishAction(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.action)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	action := getTestAction()
	savedAction, err := agg.LibUsecases.PublishAction(ctx, uid, fl, &action)
	assert.Nil(t, err)
	assert.NotNil(t, savedAction)

	type args struct {
		ctx      context.Context
		uid      string
		flavour  feedlib.Flavour
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
				uid:      uid,
				flavour:  fl,
				actionID: savedAction.ID,
			},
			wantErr: false,
		},
		{
			name: "action that does not exist - should not error",
			args: args{
				ctx:      ctx,
				uid:      uid,
				flavour:  fl,
				actionID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.DeleteAction(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.actionID); (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_PostMessage(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := testItem()

	item, err := agg.LibUsecases.PublishFeedItem(ctx, uid, fl, testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()
	invalidMessage := getInvalidTestMessage()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		itemID  string
		message *feedlib.Message
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
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
				message: &message,
			},
			wantErr: false,
		},
		{
			name: "successful post -missing message ID and sequence number set to 0",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
				message: &invalidMessage,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful post - non existent item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  ksuid.New().String(),
				message: &message,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful post - nil item",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				itemID:  item.ID,
				message: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.PostMessage(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.message,
			)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	testItem := getTestItem()
	item, err := agg.LibUsecases.PublishFeedItem(ctx, uid, fl, &testItem)
	assert.NotNil(t, item)
	assert.Nil(t, err)

	message := getTestMessage()
	savedMessage, err := agg.LibUsecases.PostMessage(ctx, uid, fl, item.ID, &message)
	assert.Nil(t, err)
	assert.NotNil(t, savedMessage)

	type args struct {
		ctx       context.Context
		uid       string
		flavour   feedlib.Flavour
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
				uid:       uid,
				flavour:   fl,
				itemID:    item.ID,
				messageID: savedMessage.ID,
			},
			wantErr: false,
		},
		{
			name: "success - message that does not exist",
			args: args{
				ctx:       ctx,
				uid:       uid,
				flavour:   fl,
				itemID:    item.ID,
				messageID: ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := agg.LibUsecases.DeleteMessage(
				tt.args.ctx,
				tt.args.uid,
				tt.args.flavour,
				tt.args.itemID,
				tt.args.messageID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFeed_ProcessEvent(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()
	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	event := getTestEvent()
	invalidEvent := getInvalidTestEvent()

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		event   *feedlib.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				event:   &event,
			},
			wantErr: false,
		},
		{
			name: "invalid event",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				event:   nil,
			},
			wantErr: true,
		},
		{
			name: "invalid flavour",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: "invalid flavour",
				event:   &event,
			},
			wantErr: true,
		},
		{
			name: "event with missing details and invalid flavour",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				event:   &invalidEvent,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.ProcessEvent(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Feed.ProcessEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLink_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)
	validLink := feedlib.Link{
		ID:          ksuid.New().String(),
		URL:         sampleVideoURL,
		LinkType:    feedlib.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   feedlib.BlankImageURL,
	}
	validLinkJSONBytes, err := json.Marshal(validLink)
	assert.Nil(t, err)
	assert.NotNil(t, validLinkJSONBytes)
	assert.Greater(t, len(validLinkJSONBytes), 3)

	invalidVideoLink := feedlib.Link{
		ID:          ksuid.New().String(),
		URL:         "www.example.com/not_a_youtube_video",
		LinkType:    feedlib.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   feedlib.BlankImageURL,
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
			l := &feedlib.Link{}
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
		Type        feedlib.LinkType
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
				Type:        feedlib.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   feedlib.BlankImageURL,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "not a valid URL",
				Type:        feedlib.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   feedlib.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid YouTube URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_video",
				Type:        feedlib.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   feedlib.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PNG URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_png",
				Type:        feedlib.LinkTypePngImage,
				Title:       "title",
				Description: "description",
				Thumbnail:   feedlib.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PDF URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_pdf",
				Type:        feedlib.LinkTypePdfDocument,
				Title:       "title",
				Description: "description",
				Thumbnail:   feedlib.BlankImageURL,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &feedlib.Link{
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		uid     string
		flavour feedlib.Flavour
		ctx     context.Context
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
				uid:     uid,
				flavour: fl,
				ctx:     ctx,
			},
			want:    []string{common.DefaultLabel},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.Labels(tt.args.ctx, tt.args.uid, tt.args.flavour)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		ctx     context.Context
		uid     string
		flavour feedlib.Flavour
		label   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx:     ctx,
				uid:     uid,
				flavour: fl,
				label:   ksuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.SaveLabel(tt.args.ctx, tt.args.uid, tt.args.flavour, tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("Feed.SaveLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_UnreadPersistentItems(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		uid     string
		flavour feedlib.Flavour
		ctx     context.Context
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
				uid:     uid,
				flavour: fl,
				ctx:     ctx,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.LibUsecases.UnreadPersistentItems(tt.args.ctx, tt.args.uid, tt.args.flavour)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	agg, _, err := InitializeTestNewFeed(ctx)
	assert.Nil(t, err)
	fl := feedlib.FlavourConsumer
	uid := ksuid.New().String()

	anonymous := false
	fe, err := agg.LibUsecases.GetThinFeed(ctx, &uid, &anonymous, fl)
	assert.Nil(t, err)
	assert.NotNil(t, fe)

	type args struct {
		uid     string
		flavour feedlib.Flavour
		ctx     context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				uid:     uid,
				flavour: fl,
				ctx:     ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := agg.LibUsecases.UpdateUnreadPersistentItemsCount(tt.args.ctx, tt.args.uid, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf(
					"Feed.UpdateUnreadPersistentItemsCount() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
