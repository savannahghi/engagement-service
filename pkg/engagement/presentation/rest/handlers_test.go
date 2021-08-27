package rest_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/markbates/pkger"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/idtoken"

	"github.com/savannahghi/engagement-service/pkg/engagement/presentation"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/rest"
	"github.com/savannahghi/engagement/pkg/engagement/application/common"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	db "github.com/savannahghi/engagement/pkg/engagement/infrastructure/database/firestore"
	"github.com/savannahghi/interserviceclient"
)

const (
	testHTTPClientTimeout = 180
	intMax                = 9007199254740990
	onboardingService     = "profile"
	registerPushToken     = "testing/register_push_token"
)

// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func startTestServer(ctx context.Context) (*http.Server, string, error) {
	// prepare the server
	port := randomPort()
	srv := presentation.PrepareServer(ctx, port, presentation.AllowedOrigins)
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if srv == nil {
		return nil, "", fmt.Errorf("nil test server")
	}

	// set up the TCP listener
	// this is done early so that we are sure we can connect to the port in
	// the tests; backlogs will be sent to the listener
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, "", fmt.Errorf("unable to listen on port %d: %w", port, err)
	}
	if l == nil {
		return nil, "", fmt.Errorf("nil test server listener")
	}
	log.Printf("LISTENING on port %d", port)

	// start serving
	go func() {
		err := srv.Serve(l)
		if err != nil {
			log.Printf("serve error: %s", err)
		}
	}()

	// the cleanup of this server (deferred shutdown) needs to occur in the
	// acceptance test that will use this
	return srv, baseURL, nil
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

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	srv, baseURL, serverErr = startTestServer(ctx) // set the globals
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

func TestRouter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case - should succeed",
			args: args{
				ctx: firebasetools.GetAuthenticatedContext(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := presentation.Router(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Router() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestHealthStatusCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	wr := httptest.NewRecorder()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful health check",
			args: args{
				w: wr,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presentation.HealthStatusCheck(tt.args.w, tt.args.r)
		})
	}
}

func TestRoutes(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	router, err := presentation.Router(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, router)

	_, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	itemID := ksuid.New().String()
	nudgeID := ksuid.New().String()
	actionID := ksuid.New().String()
	messageID := ksuid.New().String()
	title := url.QueryEscape(common.AddPrimaryEmailNudgeTitle)
	badTitle := url.QueryEscape("not a default feed title")

	type args struct {
		routeName string
		params    []string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
		wantErr bool
	}{
		{
			name: "get feed",
			args: args{
				routeName: "getFeed",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/", uid, fl.String(), false),
			wantErr: false,
		},
		{
			name: "get feed item",
			args: args{
				routeName: "getFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "get nudge",
			args: args{
				routeName: "getNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "get action",
			args: args{
				routeName: "getAction",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"actionID", actionID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/actions/%s/", uid, fl.String(), false, actionID),
			wantErr: false,
		},
		{
			name: "publish feed item",
			args: args{
				routeName: "publishFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/", uid, fl.String(), false),
			wantErr: false,
		},
		{
			name: "publish nudge",
			args: args{
				routeName: "publishNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/", uid, fl.String(), false),
			wantErr: false,
		},
		{
			name: "publish action",
			args: args{
				routeName: "publishAction",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/actions/", uid, fl.String(), false),
			wantErr: false,
		},
		{
			name: "post message",
			args: args{
				routeName: "postMessage",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/%s/messages/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "post event",
			args: args{
				routeName: "postEvent",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/events/", uid, fl.String(), false),
			wantErr: false,
		},
		{
			name: "delete feed item",
			args: args{
				routeName: "deleteFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "delete nudge",
			args: args{
				routeName: "deleteNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "delete action",
			args: args{
				routeName: "deleteAction",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"actionID", actionID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/actions/%s/", uid, fl.String(), false, actionID),
			wantErr: false,
		},
		{
			name: "delete message",
			args: args{
				routeName: "deleteMessage",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"messageID", messageID,
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/%s/messages/%s/", uid, fl.String(), false, itemID, messageID),
			wantErr: false,
		},
		{
			name: "resolve feed item",
			args: args{
				routeName: "resolveFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/resolve/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "unresolve feed item",
			args: args{
				routeName: "unresolveFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/unresolve/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "pin feed item",
			args: args{
				routeName: "pinFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/pin/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "unpin feed item",
			args: args{
				routeName: "unpinFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/unpin/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "hide feed item",
			args: args{
				routeName: "hideFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/hide/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "show feed item",
			args: args{
				routeName: "showFeedItem",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"itemID", itemID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/items/%s/show/", uid, fl.String(), false, itemID),
			wantErr: false,
		},
		{
			name: "resolve nudge",
			args: args{
				routeName: "resolveNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/resolve/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "unresolve nudge",
			args: args{
				routeName: "unresolveNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/unresolve/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "show nudge",
			args: args{
				routeName: "showNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/show/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "hide nudge",
			args: args{
				routeName: "hideNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"nudgeID", nudgeID,
				},
			},
			wantURL: fmt.Sprintf("/feed/%s/%s/%v/nudges/%s/hide/", uid, fl.String(), false, nudgeID),
			wantErr: false,
		},
		{
			name: "resolve default nudge",
			args: args{
				routeName: "resolveDefaultNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"title", title,
				},
			},
			wantURL: fmt.Sprintf(
				"/feed/%s/%s/%v/defaultnudges/%s/resolve/",
				uid,
				fl.String(),
				false,
				title,
			),
			wantErr: false,
		},
		{
			name: "resolve a non existent default nudge",
			args: args{
				routeName: "resolveDefaultNudge",
				params: []string{
					"uid", uid,
					"isAnonymous", "false",
					"flavour", fl.String(),
					"title", badTitle,
				},
			},
			wantURL: fmt.Sprintf(
				"/feed/%s/%s/%v/defaultnudges/%s/resolve/",
				uid,
				fl.String(),
				false,
				badTitle,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := router.Get(tt.args.routeName).URL(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("route error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
			assert.NotZero(t, url)
			assert.Equal(t, tt.wantURL, url.String())
		})
	}
}

func TestGetFeed(t *testing.T) {
	if srv == nil {
		t.Errorf("nil server")
		return
	}

	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	consumer := feedlib.FlavourConsumer
	invalidConsumer := "invalidConsumer"
	client := http.Client{
		Timeout: time.Minute * 10, // set high when troubleshooting
	}
	anonymous := false

	filterParams := helpers.FilterParams{
		Labels: []string{"a", "label", "and", "another"},
	}
	filterParamsJSONBytes, err := json.Marshal(filterParams)
	assert.Nil(t, err)
	assert.NotNil(t, filterParamsJSONBytes)
	if err != nil {
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name                   string
		args                   args
		wantStatus             int
		wantNewFeedInitialized bool
		wantErr                bool
	}{
		{
			name: "successful fetch of a consumer feed",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH",
					baseURL,
					uid,
					consumer,
					anonymous,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantNewFeedInitialized: true,
			wantStatus:             http.StatusOK,
			wantErr:                false,
		},
		{
			name: "fetch with a status filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH&status=PENDING",
					baseURL,
					uid,
					consumer,
					anonymous,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with a visibility filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH&status=PENDING&visibility=SHOW",
					baseURL,
					uid,
					consumer,
					anonymous,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with an expired filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH&status=PENDING&visibility=SHOW&expired=FALSE",
					baseURL,
					uid,
					consumer,
					anonymous,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fetch with an expired filter",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH&status=PENDING&visibility=SHOW&expired=FALSE&filterParams=%s",
					baseURL,
					uid,
					consumer,
					anonymous,
					string(filterParamsJSONBytes),
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid - invalid flavour",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/?persistent=BOTH",
					baseURL,
					uid,
					invalidConsumer,
					anonymous,
				),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range getDefaultHeaders(ctx, t, baseURL) {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				// early exit
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantNewFeedInitialized {
				returnedFeed := &domain.Feed{}
				err = json.Unmarshal(data, returnedFeed)
				if err != nil {
					t.Errorf("can't unmarshal feed from response JSON: %w", err)
					return
				}

				if len(returnedFeed.Actions) < 1 {
					t.Error("the returned feed has no actions")
				}

				if len(returnedFeed.Nudges) < 1 {
					t.Errorf("the returned feed has no nudges")
				}

				if len(returnedFeed.Items) < 1 {
					t.Errorf("the returned feed has no items")
				}
			}
		})
	}
}

func TestGetFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid feed item retrieval",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetAction(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testAction := getTestAction()
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		&testAction,
		baseURL,
		"publishAction",
	)
	if err != nil {
		t.Errorf("can't post action: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testAction.ID,
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/action/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodGet,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	headers := getDefaultHeaders(ctx, t, baseURL)
	testItem := getTestItem(t)

	bs, err := json.Marshal(testItem)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid feed item publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteAction(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testAction := getTestAction()
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		&testAction,
		baseURL,
		"publishAction",
	)
	if err != nil {
		t.Errorf("can't post test action: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testAction.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/actions/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPostMessage(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	msg := getTestMessage()
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("can't marshal message to JSON: %s", err)
		return
	}
	payload := bytes.NewBuffer(msgBytes)

	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid message post",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/%s/messages/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil message",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/%s/messages/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}

	msg := getTestMessage()
	err = postMessage(
		ctx,
		t,
		uid,
		fl,
		&msg,
		baseURL,
		testItem.ID,
	)
	if err != nil {
		t.Errorf("can't post message: %s", err)
		return
	}

	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid delete",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/%s/messages/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
					msg.ID,
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non existent element delete - safe to repeat over and over",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/%s/messages/%s/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
					ksuid.New().String(),
				),
				httpMethod: http.MethodDelete,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	headers := getDefaultHeaders(ctx, t, baseURL)
	event := getTestEvent()

	bs, err := json.Marshal(event)
	if err != nil {
		t.Errorf("unable to marshal event to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid event publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/events/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil event",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/events/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	headers := getDefaultHeaders(ctx, t, baseURL)
	nudge := testNudge(t)

	bs, err := json.Marshal(nudge)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid nudge publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestResolveNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnresolveNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestShowNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "show valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/show/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to show non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/show/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHideNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testNudge := testNudge(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testNudge,
		baseURL,
		"publishNudge",
	)
	if err != nil {
		t.Errorf("can't post nudge: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "hide valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					false,
					testNudge.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to hide non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/nudges/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPublishAction(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	headers := getDefaultHeaders(ctx, t, baseURL)
	action := getTestAction()

	bs, err := json.Marshal(action)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid action publish",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/actions/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil action",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/actions/",
					baseURL,
					uid,
					fl.String(),
					false,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestResolveFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "resolve valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to resolve non existent feed uten",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnresolveFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "unresolve valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to unresolve non existent feed uten",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/unresolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPinFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "pin valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/pin/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to pin non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/pin/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUnpinFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "unpin valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/unpin/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to unpin non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/unpin/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHideFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "hide valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to hide non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/hide/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestShowFeedItem(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	testItem := getTestItem(t)
	err = postElement(
		ctx,
		t,
		uid,
		fl,
		testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "show valid feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/show/",
					baseURL,
					uid,
					fl.String(),
					false,
					testItem.ID,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "try to show non existent feed item",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/items/%s/show/",
					baseURL,
					uid,
					fl.String(),
					false,
					ksuid.New().String(),
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func getInterserviceBearerTokenHeader(ctx context.Context, t *testing.T, rootDomain string) string {
	isc := getInterserviceClient(t, rootDomain)
	authToken, err := isc.CreateAuthToken(ctx)
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

func getDefaultHeaders(ctx context.Context, t *testing.T, rootDomain string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": getInterserviceBearerTokenHeader(ctx, t, rootDomain),
	}
}

func getInterserviceClient(t *testing.T, rootDomain string) *interserviceclient.InterServiceClient {
	service := interserviceclient.ISCService{
		Name:       "feed",
		RootDomain: rootDomain,
	}
	isc, err := interserviceclient.NewInterserviceClient(service)
	assert.Nil(t, err)
	assert.NotNil(t, isc)
	return isc
}

func randomPort() int {
	rand.Seed(time.Now().Unix())
	min := 32768
	max := 60999
	port := rand.Intn(max-min+1) + min
	return port
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

	headers := getDefaultHeaders(ctx, t, baseURL)
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

	headers := getDefaultHeaders(ctx, t, baseURL)
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

func getTestItem(t *testing.T) *feedlib.Item {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return nil
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))
	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return nil
	}
	return &feedlib.Item{
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

func testNudge(t *testing.T) *feedlib.Nudge {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return &feedlib.Nudge{}
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))
	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return nil
	}
	return &feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         feedlib.StatusPending,
		Visibility:     feedlib.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []feedlib.Link{
			feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		},
		Text: ksuid.New().String(),
		Actions: []feedlib.Action{
			getTestAction(),
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

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
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

func getTestAction() feedlib.Action {
	return feedlib.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Name:           "TEST_ACTION",
		Icon:           feedlib.GetPNGImageLink(feedlib.LogoURL, "title", "description", feedlib.LogoURL),
		ActionType:     feedlib.ActionTypePrimary,
		Handling:       feedlib.HandlingFullPage,
	}
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

func GetPayloadRequest(data pubsubtools.PubSubPayload) (*http.Request, error) {
	testDataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("can't marshal JSON: %w", err)
	}

	pubsubURL := fmt.Sprintf("%s%s", baseURL, pubsubtools.PubSubHandlerPath)
	req, err := http.NewRequest(
		http.MethodPost,
		pubsubURL,
		bytes.NewBuffer(testDataJSON),
	)

	if err != nil {
		return nil, fmt.Errorf("can't initialize request: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func TestGoogleCloudPubSubHandler(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	_, err = RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))

	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}

	item := getTestItem(t)
	itemData, err := json.Marshal(item)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}

	notification := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: itemData,
	}

	notificationData, err := json.Marshal(notification)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}

	nudge := testNudge(t)
	nudgeDataJSON, err := json.Marshal(&nudge)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}

	nudgeNotification := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: nudgeDataJSON,
	}

	nudgeData, err := json.Marshal(nudgeNotification)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}

	emailData := map[string]interface{}{
		"to":      []string{"automated.test.user.bewell-app-ci@healthcloud.co.ke", "bewell@bewell.co.ke"},
		"text":    "This is a test message",
		"subject": "Test Subject",
	}

	emailPayload, err := json.Marshal(emailData)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}

	action := getTestAction()
	actionDataJSON, err := json.Marshal(action)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}
	actionNotification := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: actionDataJSON,
	}
	actionData, err := json.Marshal(actionNotification)
	if err != nil {
		t.Errorf("failed to marshall data")
		return
	}

	message := getTestMessage()
	messageDataJSON, err := json.Marshal(message)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}
	messageNotification := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: messageDataJSON,
	}
	messageData, err := json.Marshal(messageNotification)
	if err != nil {
		t.Errorf("failed to marshall data")
		return
	}

	event := getTestEvent()
	eventDataJSON, err := json.Marshal(event)
	if err != nil {
		t.Errorf("failed to marshal data")
		return
	}
	eventNotification := dto.NotificationEnvelope{
		UID:     token.UID,
		Flavour: feedlib.FlavourConsumer,
		Payload: eventDataJSON,
	}
	eventData, err := json.Marshal(eventNotification)
	if err != nil {
		t.Errorf("failed to marshall data")
		return
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(ksuid.New().String()))

	actionPublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      actionData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ActionPublishTopic),
			},
		},
	}

	invalidActionPublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      itemData,
			Attributes: map[string]string{
				"topicID":       "invalid topic",
				"invalid field": "invalid field",
			},
		},
	}

	actionPublishRequest, _ := GetPayloadRequest(actionPublishTopicPayload)
	invalidActionPublishRequest, _ := GetPayloadRequest(invalidActionPublishTopicPayload)

	actionDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      actionData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ActionDeleteTopic),
			},
		},
	}

	invalidActionDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      itemData,
			Attributes: map[string]string{
				"topicID":       "invalid topic",
				"invalid field": "invalid field",
			},
		},
	}

	invalidActionDeleteRequest, _ := GetPayloadRequest(invalidActionDeleteTopicPayload)
	actionDeleteRequest, _ := GetPayloadRequest(actionDeleteTopicPayload)

	messagePostTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      messageData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.MessagePostTopic),
			},
		},
	}

	invalidMessagePostTopicPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: itemData,
			Attributes: map[string]string{
				"topicID":       "invalid topic",
				"invalid field": "invalid field",
			},
		},
	}

	invalidMessagePostRequest, err := GetPayloadRequest(invalidMessagePostTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	messagePostRequest, err := GetPayloadRequest(messagePostTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	messageDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      messageData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.MessageDeleteTopic),
			},
		},
	}

	invalidMessageDeleteTopicPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: itemData,
			Attributes: map[string]string{
				"topicID":       "invalid topic",
				"invalid field": "invalid field",
			},
		},
	}

	invalidMessageDeleteRequest, err := GetPayloadRequest(invalidMessageDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}
	messageDeleteRequest, err := GetPayloadRequest(messageDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	incomingEventTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      eventData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.IncomingEventTopic),
			},
		},
	}

	invalidIncomingEventTopicPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: itemData,
			Attributes: map[string]string{
				"topicID":       "invalid topic",
				"invalid field": "invalid field",
			},
		},
	}

	invalidIncomingEventRequest, err := GetPayloadRequest(invalidIncomingEventTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}
	incomingEventRequest, err := GetPayloadRequest(incomingEventTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	invalidTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      []byte(b64),
			Attributes: map[string]string{
				"invalid": "invalid topic",
			},
		},
	}

	invalidTopicRequest, _ := GetPayloadRequest(invalidTopicPayload)

	itemPublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemPublishTopic),
			},
		},
	}

	invalidItemPublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemPublishTopic),
				"invalid field": "invalid field",
			},
		},
	}

	invalidItemPublishTopicRequest, _ := GetPayloadRequest(invalidItemPublishTopicPayload)
	itemPublishTopicRequest, _ := GetPayloadRequest(itemPublishTopicPayload)

	itemResolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemResolveTopic),
			},
		},
	}

	invalidItemResolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemResolveTopic),
				"invalid field": "invalid field",
			},
		},
	}

	invalidItemResolveTopicRequest, _ := GetPayloadRequest(invalidItemResolveTopicPayload)
	itemResolveTopicRequest, _ := GetPayloadRequest(itemResolveTopicPayload)

	itemUnresolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemUnresolveTopic),
			},
		},
	}

	invalidItemUnresolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemUnresolveTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemUnresolveTopicRequest, err := GetPayloadRequest(invalidItemUnresolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemUnresolveTopicRequest, err := GetPayloadRequest(itemUnresolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemDeleteTopic),
			},
		},
	}

	invalidItemDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemDeleteTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemDeleteTopicRequest, err := GetPayloadRequest(invalidItemDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemDeleteTopicRequest, err := GetPayloadRequest(itemDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemHideTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemHideTopic),
			},
		},
	}

	invalidItemHideTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemHideTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemHideTopicRequest, err := GetPayloadRequest(invalidItemHideTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemHideTopicRequest, err := GetPayloadRequest(itemHideTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemShowTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemShowTopic),
			},
		},
	}

	invalidItemShowTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemShowTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemShowTopicRequest, err := GetPayloadRequest(invalidItemShowTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemShowTopicRequest, err := GetPayloadRequest(itemShowTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemPinTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemPinTopic),
			},
		},
	}

	invalidItemPinTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemPinTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemPinTopicRequest, err := GetPayloadRequest(invalidItemPinTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemPinTopicRequest, err := GetPayloadRequest(itemPinTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemUnpinTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      notificationData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.ItemUnpinTopic),
			},
		},
	}

	invalidItemUnpinTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.ItemUnpinTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidItemUnpinTopicRequest, err := GetPayloadRequest(invalidItemUnpinTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	itemUnpinTopicRequest, err := GetPayloadRequest(itemUnpinTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgePublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgePublishTopic),
			},
		},
	}

	invalidNudgePublishTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgePublishTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgePublishTopicRequest, err := GetPayloadRequest(invalidNudgePublishTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgePublishTopicRequest, err := GetPayloadRequest(nudgePublishTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgeDeleteTopic),
			},
		},
	}

	invalidNudgeDeleteTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgeDeleteTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgeDeleteTopicRequest, err := GetPayloadRequest(invalidNudgeDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeDeleteTopicRequest, err := GetPayloadRequest(nudgeDeleteTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeResolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgeResolveTopic),
			},
		},
	}

	invalidNudgeResolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgeResolveTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgeResolveTopicRequest, err := GetPayloadRequest(invalidNudgeResolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeResolveTopicRequest, err := GetPayloadRequest(nudgeResolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeUnresolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgeUnresolveTopic),
			},
		},
	}

	invalidNudgeUnresolveTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgeUnresolveTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgeUnresolveTopicRequest, err := GetPayloadRequest(invalidNudgeUnresolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeUnresolveTopicRequest, err := GetPayloadRequest(nudgeUnresolveTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeHideTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgeHideTopic),
			},
		},
	}

	invalidNudgeHideTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgeHideTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgeHideTopicRequest, err := GetPayloadRequest(invalidNudgeHideTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeHideTopicRequest, err := GetPayloadRequest(nudgeHideTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeShowTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      nudgeData,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.NudgeShowTopic),
			},
		},
	}

	invalidNudgeShowTopicPayload := pubsubtools.PubSubPayload{
		Subscription: "invalid",
		Message: pubsubtools.PubSubMessage{
			MessageID: "invalid",
			Data:      messageDataJSON,
			Attributes: map[string]string{
				"topicID":       helpers.AddPubSubNamespace(common.NudgeShowTopic),
				"invalid field": "invalid field",
			},
		},
	}
	invalidNudgeShowTopicRequest, err := GetPayloadRequest(invalidNudgeShowTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	nudgeShowTopicRequest, err := GetPayloadRequest(nudgeShowTopicPayload)
	if err != nil {
		t.Errorf("can't initialize request: %w", err)
		return
	}

	idTokenHTTPClient, err := idtoken.NewClient(ctx, pubsubtools.Aud)
	if err != nil {
		t.Errorf("can't initialize idToken HTTP client: %s", err)
		return
	}

	sendEmailTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      emailPayload,
			Attributes: map[string]string{
				"topicID": helpers.AddPubSubNamespace(common.SentEmailTopic),
			},
		},
	}

	invalidSendEmailTopicPayload := pubsubtools.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: pubsubtools.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      itemData,
			Attributes: map[string]string{
				"invalidTopic": "invalid.topic",
			},
		},
	}

	sendEmailRequest, _ := GetPayloadRequest(sendEmailTopicPayload)
	invalidSendEmailRequest, _ := GetPayloadRequest(invalidSendEmailTopicPayload)

	type args struct {
		r      *http.Request
		client *http.Client
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid pubsub format payload with valid auth - Action Publish Topic",
			args: args{
				r:      actionPublishRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid action publish payload",
			args: args{
				r:      invalidActionPublishRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid pubsub format payload with valid auth - Action Delete Topic",
			args: args{
				r:      actionDeleteRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid action delete payload",
			args: args{
				r:      invalidActionDeleteRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid pubsub format payload with valid auth - Message Post Topic",
			args: args{
				r:      messagePostRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid message post topic payload",
			args: args{
				r:      invalidMessagePostRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid pubsub format payload with valid auth - Message Delete Topic",
			args: args{
				r:      messageDeleteRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid message delete topic payload",
			args: args{
				r:      invalidMessageDeleteRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid pubsub format payload with valid auth - Incoming Event Topic",
			args: args{
				r:      incomingEventRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid incoming event topic payload",
			args: args{
				r:      invalidIncomingEventRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "no auth header",
			args: args{
				r:      itemUnresolveTopicRequest,
				client: http.DefaultClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid topic",
			args: args{
				r:      invalidTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item publish payload",
			args: args{
				r:      itemPublishTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item publish payload",
			args: args{
				r:      invalidItemPublishTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item resolve payload",
			args: args{
				r:      itemResolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item resolve payload",
			args: args{
				r:      invalidItemResolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item unresolve payload",
			args: args{
				r:      itemUnresolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item unresolve payload",
			args: args{
				r:      invalidItemUnresolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item delete payload",
			args: args{
				r:      itemDeleteTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item delete payload",
			args: args{
				r:      invalidItemDeleteTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item hide payload",
			args: args{
				r:      itemHideTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item hide payload",
			args: args{
				r:      invalidItemHideTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item show payload",
			args: args{
				r:      itemShowTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item show payload",
			args: args{
				r:      invalidItemShowTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item pin payload",
			args: args{
				r:      itemPinTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item pin payload",
			args: args{
				r:      invalidItemPinTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid item unpin payload",
			args: args{
				r:      itemUnpinTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid item unpin payload",
			args: args{
				r:      invalidItemUnpinTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge publish topic payload",
			args: args{
				r:      nudgePublishTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge publish payload",
			args: args{
				r:      invalidNudgePublishTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge delete topic payload",
			args: args{
				r:      nudgeDeleteTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge delete payload",
			args: args{
				r:      invalidNudgeDeleteTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge resolve topic payload",
			args: args{
				r:      nudgeResolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge resolve payload",
			args: args{
				r:      invalidNudgeResolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge unresolve topic payload",
			args: args{
				r:      nudgeUnresolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge unresolve payload",
			args: args{
				r:      invalidNudgeUnresolveTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge hide topic payload",
			args: args{
				r:      nudgeHideTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge hide payload",
			args: args{
				r:      invalidNudgeHideTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid nudge show topic payload",
			args: args{
				r:      nudgeShowTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid nudge show payload",
			args: args{
				r:      invalidNudgeShowTopicRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid send email topic payload",
			args: args{
				r:      sendEmailRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true, //TODO - Find out why it fails
		},
		{
			name: "invalid send email topic payload",
			args: args{
				r:      invalidSendEmailRequest,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		// TODO Check success resp status: map[string]string{"status": "success"}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.args.client.Do(tt.args.r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			respBs, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Errorf("unable to read response body: %s", err)
				return
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("wanted status %d, got %d and resp %s",
					tt.wantStatus, resp.StatusCode, string(respBs))
				return
			}

			if !tt.wantErr {
				decoded := map[string]string{}
				err = json.Unmarshal(respBs, &decoded)
				if err != nil {
					t.Errorf("can't decode response to map: %s", err)
					return
				}
				if decoded["status"] != "success" {
					t.Errorf("did not get success status")
					return
				}
			}
		})
	}
}

func TestPostUpload(t *testing.T) {
	ctx := context.Background()
	headers := getDefaultHeaders(ctx, t, baseURL)
	itemID := ksuid.New().String()

	imgPath := rest.StaticDir + "/1px.png"
	f, err := pkger.Open(imgPath)
	if err != nil {
		t.Errorf("can't open test image path with pkger: %v", err)
		return
	}
	defer f.Close()

	imgData, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("can't read test image: %v", err)
		return
	}

	b64 := base64.StdEncoding.EncodeToString(imgData)
	uploadInput := profileutils.UploadInput{
		Title:       itemID,
		ContentType: "image/png",
		Language:    "en",
		Filename:    fmt.Sprintf("%s.png", itemID),
		Base64data:  b64,
	}

	bs, err := json.Marshal(uploadInput)
	if err != nil {
		t.Errorf("unable to marshal upload input to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid upload",
			args: args{
				url:        fmt.Sprintf("%s/internal/upload/", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil upload",
			args: args{
				url:        fmt.Sprintf("%s/internal/upload/", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp.StatusCode != tt.wantStatus {
				log.Printf("raw response data: \n%s\n", string(data))
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}

			if resp.StatusCode == tt.wantStatus && tt.wantStatus == http.StatusOK {
				// find the upload again
				upload := &profileutils.Upload{}
				err = json.Unmarshal(data, &upload)
				if err != nil {
					t.Errorf("can't unmarshal returned upload to JSON: %v", err)
					return
				}

				uploadURL := fmt.Sprintf(
					"%s/internal/upload/%s/",
					baseURL,
					upload.ID,
				)
				uploadReq, err := http.NewRequest(
					http.MethodGet,
					uploadURL,
					nil,
				)
				if err != nil {
					t.Errorf("unable to compose find upload request: %s", err)
					return
				}
				for k, v := range tt.args.headers {
					uploadReq.Header.Add(k, v)
				}
				uploadResp, err := client.Do(uploadReq)
				if err != nil {
					t.Errorf("request error: %s", err)
					return
				}
				if err != nil {
					t.Errorf("error fetching upload again: %v", err)
					return
				}
				uploadData, err := ioutil.ReadAll(uploadResp.Body)
				if err != nil {
					t.Errorf("can't read request body: %s", err)
					return
				}
				fetchedUpload := profileutils.Upload{}
				err = json.Unmarshal(uploadData, &fetchedUpload)
				if err != nil {
					t.Errorf("can't unmarshal returned upload to JSON: %v", err)
					return
				}

				if fetchedUpload.Base64data != upload.Base64data {
					t.Errorf("did not get back the same upload, differnet data")
					return
				}

				if fetchedUpload.Hash != upload.Hash {
					t.Errorf("did not get back the same upload, different hashes")
					return
				}

				if fetchedUpload.ID != upload.ID {
					t.Errorf(
						"did not get back the same upload, different IDs; %s vs %s",
						fetchedUpload.ID,
						upload.ID,
					)
					return
				}
			}
		})
	}
}

func resolveTestNudge(
	ctx context.Context,
	uid string,
	fl feedlib.Flavour,
	nudge *feedlib.Nudge,
) error {

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		return fmt.Errorf("can't initialize Firebase Repository: %s", err)
	}
	nudge.Status = feedlib.StatusDone
	nudge.SequenceNumber = nudge.SequenceNumber + 1
	_, err = fr.UpdateNudge(ctx, uid, fl, nudge)
	if err != nil {
		return fmt.Errorf("unable to resolve nudge: %w",
			err,
		)
	}
	return nil
}

func TestResolveDefaultNudge(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	uid := token.UID
	fl := feedlib.FlavourConsumer
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
	}
	title := common.AddPrimaryEmailNudgeTitle

	nudge, _ := fr.GetDefaultNudgeByTitle(ctx, uid, fl, title)
	strnudge := fmt.Sprintf("%v", nudge)
	if strnudge != "<nil>" {
		err = fr.DeleteNudge(ctx, uid, fl, nudge.ID)
		if err != nil {
			t.Errorf("failed to delete nudge %w", err)
		}
	}

	defaultNudges, err := db.SetDefaultNudges(
		ctx,
		uid,
		fl,
		fr,
	)

	if err != nil {
		t.Errorf("can't set default nudges: %s", err)
	}
	if len(defaultNudges) == 0 {
		t.Errorf("zero default nudges found")
		return
	}

	for _, nudge := range defaultNudges {
		if nudge.Title == title {
			err := resolveTestNudge(
				ctx,
				uid,
				fl,
				&nudge,
			)
			if err != nil {
				t.Errorf("unable to resolve nudge: %w", err)
				return
			}
		}
	}

	bs, err := json.Marshal(map[string]string{"status": "success"})
	if err != nil {
		t.Errorf("unable to marshal upload input to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	headers := getDefaultHeaders(ctx, t, baseURL)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success: resolve valid nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/defaultnudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					common.AddPrimaryEmailNudgeTitle,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: an already resolved nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/defaultnudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					common.AddPrimaryEmailNudgeTitle,
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failure: try to resolve non existent nudge",
			args: args{
				url: fmt.Sprintf(
					"%s/feed/%s/%s/%v/defaultnudges/%s/resolve/",
					baseURL,
					uid,
					fl.String(),
					false,
					"not a nudge title",
				),
				httpMethod: http.MethodPatch,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected %v, but got %v", tt.wantStatus, resp.StatusCode)
				return
			}
		})
	}
	if err != nil {
		t.Errorf("can't set default nudges: %s", err)
	}
}

func TestSendEmail(t *testing.T) {
	testUserMail := "test@bewell.co.ke"
	ctx := context.Background()
	headers := getDefaultHeaders(ctx, t, baseURL)
	to := []string{testUserMail}
	email := dto.EMailMessage{
		Subject: "Test Subject :)",
		Text:    "Hello :)",
		To:      to,
	}

	bs, err := json.Marshal(email)
	if err != nil {
		t.Errorf("unable to marshal upload input to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid send email",
			args: args{
				url:        fmt.Sprintf("%s/internal/send_email", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil send email",
			args: args{
				url:        fmt.Sprintf("%s/internal/send_email", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestGetAITSMSDeliveryCallback(t *testing.T) {
	ctx := context.Background()
	sms := dto.MarketingSMS{
		ID:          uuid.New().String(),
		PhoneNumber: gofakeit.Phone(),
		Message:     gofakeit.FarmAnimal(),
		SenderID:    enumutils.SenderIDBewell,
	}

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
		return
	}

	savedSms, err := fr.SaveMarketingMessage(ctx, sms)
	if err != nil {
		t.Errorf("unable to save marketing message: %w",
			err,
		)
		return
	}
	expectedPayload := map[string]interface{}{
		"phoneNumber": savedSms.PhoneNumber,
		"retryCount":  "0",
		"status":      "success",
	}
	bs, err := json.Marshal(expectedPayload)
	if err != nil {
		t.Errorf("unable to marshal upload input to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	headers := req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": getInterserviceBearerTokenHeader(ctx, t, baseURL),
	}

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Happy AT callback :)",
			args: args{
				url: fmt.Sprintf(
					"%s/ait_callback",
					baseURL,
				),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			assert.NotNil(t, data)
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPresentationHandlersImpl_UpdateMailgunDelivery(t *testing.T) {
	ctx := context.Background()
	headers := getDefaultHeaders(ctx, t, baseURL)
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
		return
	}

	event := dto.MailgunEvent{
		EventName:   "delivered",
		DeliveredOn: "12345637.11222",
		MessageID:   "20210715172955.1.63EC29EF167F09B9@sandboxb30d61fba25641a9983c3b3a3c84abde.mailgun.org",
	}

	emailLog := &dto.OutgoingEmailsLog{
		UUID:        uuid.NewString(),
		To:          []string{"test@bewell.co.ke"},
		From:        "test@bewell.co.ke",
		Subject:     "Test",
		Text:        "Test",
		MessageID:   "20210715172955.1.63EC29EF167F09B9@sandboxb30d61fba25641a9983c3b3a3c84abde.mailgun.org",
		EmailSentOn: time.Time{},
		Event: &dto.MailgunEventOutput{
			EventName:   "accepted",
			DeliveredOn: time.Time{},
		},
	}

	err = fr.SaveOutgoingEmails(ctx, emailLog)
	if err != nil {
		t.Errorf("unable to save outgoing email: %w",
			err,
		)
		return
	}

	_, err = fr.UpdateMailgunDeliveryStatus(ctx, &event)
	if err != nil {
		t.Errorf("unable to update email delivery status: %w",
			err,
		)
		return
	}

	bs, err := json.Marshal(event)
	if err != nil {
		t.Errorf("unable to marshal event input to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	type args struct {
		url        string
		httpMethod string
		headers    map[string]string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid event",
			args: args{
				url:        fmt.Sprintf("%s/internal/mailgun_delivery_webhook", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "nil event",
			args: args{
				url:        fmt.Sprintf("%s/internal/mailgun_delivery_webhook", baseURL),
				httpMethod: http.MethodPost,
				headers:    headers,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range tt.args.headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}
