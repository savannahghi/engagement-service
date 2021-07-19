package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/imroc/req"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation"
)

const (
	testHTTPClientTimeout = 180
	intMax                = 9007199254740990
)

// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	// run the tests
	code := m.Run()

	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func GetTestSequenceNumber() int {
	return rand.Intn(intMax)
}

func GetTestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: GetTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		ActionType:     feedlib.ActionTypePrimary,
		Handling:       feedlib.HandlingFullPage,
	}
}

func getTestMessage() feedlib.Message {
	return feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: GetTestSequenceNumber(),
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
}

func getTestItem() feedlib.Item {
	return feedlib.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       feedlib.TextTypePlain,
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
			feedlib.GetYoutubeVideoLink(base.SampleVideoURL, "title", "description", base.LogoURL),
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
				ActionType:     feedlib.ActionTypeSecondary,
				Handling:       feedlib.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
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

func GetInterserviceClient(t *testing.T, rootDomain string) *base.InterServiceClient {
	service := base.ISCService{
		Name:       "feed",
		RootDomain: rootDomain,
	}
	isc, err := base.NewInterserviceClient(service)
	assert.Nil(t, err)
	assert.NotNil(t, isc)
	return isc
}

func GetInterserviceBearerTokenHeader(ctx context.Context, t *testing.T, rootDomain string) string {
	isc := GetInterserviceClient(t, rootDomain)
	authToken, err := isc.CreateAuthToken(ctx)
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

func GetDefaultHeaders(ctx context.Context, t *testing.T, rootDomain string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": GetInterserviceBearerTokenHeader(ctx, t, rootDomain),
	}
}

func getGraphQLHeaders(t *testing.T) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": GetBearerTokenHeader(t),
	}
}

func GetBearerTokenHeader(t *testing.T) string {
	ctx := context.Background()
	user, err := base.GetOrCreateFirebaseUser(ctx, converterandformatter.TestUserEmail)
	if err != nil {
		t.Errorf("can't get or create firebase user: %s", err)
		return ""
	}

	if user == nil {
		t.Errorf("nil firebase user")
		return ""
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, user.UID)
	if err != nil {
		t.Errorf("can't create custom token: %s", err)
		return ""
	}

	if customToken == "" {
		t.Errorf("blank custom token: %s", err)
		return ""
	}

	idTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		t.Errorf("can't authenticate custom token: %s", err)
		return ""
	}
	if idTokens == nil {
		t.Errorf("nil idTokens")
		return ""
	}

	return fmt.Sprintf("Bearer %s", idTokens.IDToken)
}

func testNudge() *feedlib.Nudge {
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: GetTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			GetTestAction(),
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

func postElement(
	ctx context.Context,
	t *testing.T,
	uid string,
	fl feedlib.Flavour,
	el feedlib.Element,
	baseURL string,
	routeName string,
) error {
	router, err := presentation.Router(ctx)
	if err != nil {
		t.Errorf("can't set up router: %s", err)
		return err
	}

	params := []string{
		"uid", uid,
		"flavour", fl.String(),
		"isAnonymous", "false",
	}

	route := router.Get(routeName)
	if route == nil {
		return fmt.Errorf(
			"there's no registered route with the name `%s`", routeName)
	}
	path, err := router.Get(routeName).URL(params...)
	if err != nil {
		t.Errorf("can't get URL: %s", err)
		return err
	}
	url := fmt.Sprintf("%s%s", baseURL, path.String())

	data, err := json.Marshal(el)
	if err != nil {
		t.Errorf("can't marshal nudge to JSON: %s", err)
		return err
	}
	payload := bytes.NewBuffer(data)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		t.Errorf("error when creating request to post `%v` to %s: %s", payload, url, err)
		return err
	}
	if r == nil {
		t.Errorf("nil request when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	headers := GetDefaultHeaders(ctx, t, baseURL)
	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	resp, err := client.Do(r)
	if resp == nil {
		t.Errorf("nil response: %s", err)
		return err
	}

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("error when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	assert.NotNil(t, data)
	if data == nil {
		t.Errorf("nil response when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
		return fmt.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
	}

	return nil
}

func postMessage(
	ctx context.Context,
	t *testing.T,
	uid string,
	fl feedlib.Flavour,
	el feedlib.Element,
	baseURL string,
	itemID string,
) error {
	router, err := presentation.Router(ctx)
	if err != nil {
		t.Errorf("can't set up router: %s", err)
		return err
	}

	params := []string{
		"uid", uid,
		"flavour", fl.String(),
		"itemID", itemID,
		"isAnonymous", "false",
	}

	path, err := router.Get("postMessage").URL(params...)
	if err != nil {
		t.Errorf("can't get URL: %s", err)
		return err
	}
	url := fmt.Sprintf("%s%s", baseURL, path.String())

	data, err := json.Marshal(el)
	if err != nil {
		t.Errorf("can't marshal nudge to JSON: %s", err)
		return err
	}
	payload := bytes.NewBuffer(data)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		t.Errorf("error when creating request to post `%v` to %s: %s", payload, url, err)
		return err
	}
	if r == nil {
		t.Errorf("nil request when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	headers := GetDefaultHeaders(ctx, t, baseURL)
	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := http.DefaultClient
	resp, err := client.Do(r)
	if resp == nil {
		t.Errorf("nil response: %s", err)
		return err
	}

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("error when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	assert.NotNil(t, data)
	if data == nil {
		t.Errorf("nil response when posting `%v` to %s: %s", payload, url, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
		return fmt.Errorf("error status code `%d` and data `%s`", resp.StatusCode, data)
	}

	return nil
}
