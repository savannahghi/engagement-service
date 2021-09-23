package usecases_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/usecases"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	libInfra "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	libRepository "github.com/savannahghi/engagementcore/pkg/engagement/repository"
	libNotification "github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/segmentio/ksuid"
)

const (
	// sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"
	onboardingService = "profile"
	intMax            = 9007199254740990
	registerPushToken = "testing/register_push_token"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("ENVIRONMENT", "staging")
	os.Exit(m.Run())
}

func InitializeTestNewNotification(ctx context.Context) (*usecases.NotificationImpl, libRepository.Repository, error) {
	var repo libRepository.Repository
	infra := libInfra.NewInteractor()
	libUsc := libNotification.NewNotification(infra)
	lib := usecases.NewNotification(repo, libUsc)
	return lib, repo, nil
}

func onboardingISCClient(t *testing.T) *interserviceclient.InterServiceClient {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil
	}

	return profileClient
}

func RegisterPushToken(
	ctx context.Context,
	t *testing.T,
	UID string,
	onboardingClient *interserviceclient.InterServiceClient,
) (bool, error) {
	token := "random"
	if onboardingClient == nil {
		return false, fmt.Errorf("nil ISC client")
	}

	payload := map[string]interface{}{
		"pushTokens": token,
		"uid":        UID,
	}
	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		registerPushToken,
		payload,
	)
	if err != nil {
		return false, fmt.Errorf("unable to make a request to register push token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("expected a StatusOK (200) status code but instead got %v", resp.StatusCode)
	}

	return true, nil
}

func getATestSequenceNumber() int {
	return rand.Intn(intMax)
}

func getATestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTheTestItem(t *testing.T) feedlib.Item {
	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return feedlib.Item{}
	}

	return feedlib.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feedlib.TextTypePlain,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
			feedlib.GetYoutubeVideoLink(feedlib.SampleVideoURL, "title", "description", feedlib.LogoURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingInline,
			},
		},
		Conversations: []feedlib.Message{
			{
				ID:             "msg-2",
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
				SequenceNumber: int(time.Now().Unix()),
			},
		},
		Users: []string{
			token.UID,
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

func getATestEvent() feedlib.Event {
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

func getATestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		ActionType:     feedlib.ActionTypePrimary,
		Handling:       feedlib.HandlingFullPage,
	}
}

func getTestPubsubPayload(t *testing.T, el feedlib.Element) *pubsubtools.PubSubPayload {
	elData, err := el.ValidateAndMarshal()
	if err != nil {
		t.Errorf("invalid element: %v", err)
		return nil
	}

	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}

	envelope := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: elData,
		Metadata: map[string]interface{}{
			ksuid.New().String(): ksuid.New().String(),
		},
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		t.Errorf("can't marshal envelope data: %v", err)
		return nil
	}

	return &pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      data,
			Attributes: map[string]string{
				"topicID": ksuid.New().String(),
			},
		},
	}
}

func aTestNudge(t *testing.T) *feedlib.Nudge {
	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getATestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			getATestAction(),
		},
		Users: []string{
			token.UID,
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
