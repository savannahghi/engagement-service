package rest_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/rest"

	"github.com/imroc/req"
	"github.com/markbates/pkger"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"google.golang.org/api/idtoken"
)

const (
	testHTTPClientTimeout = 180
	intMax                = 9007199254740990
	onboardingService     = "profile"
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
				ctx: context.Background(),
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
	ctx := context.Background()
	router, err := presentation.Router(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, router)

	uid := xid.New().String()
	fl := base.FlavourConsumer
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

	uid := xid.New().String()
	consumer := base.FlavourConsumer
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

			for k, v := range getDefaultHeaders(t, baseURL) {
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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestGetNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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

func TestGetAction(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testAction := getTestAction()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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
	uid := xid.New().String()
	fl := base.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	testItem := getTestItem()

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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testAction := getTestAction()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
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

	headers := getDefaultHeaders(t, baseURL)

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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
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

	headers := getDefaultHeaders(t, baseURL)

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
	uid := xid.New().String()
	fl := base.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
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
	uid := xid.New().String()
	fl := base.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
	nudge := testNudge()

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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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

func TestUnresolveNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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

func TestShowNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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

func TestHideNudge(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testNudge := testNudge()
	err := postElement(
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
	headers := getDefaultHeaders(t, baseURL)

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

func TestPublishAction(t *testing.T) {
	uid := xid.New().String()
	fl := base.FlavourConsumer
	headers := getDefaultHeaders(t, baseURL)
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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestUnresolveFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestPinFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestUnpinFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestHideFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func TestShowFeedItem(t *testing.T) {
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	testItem := getTestItem()
	err := postElement(
		ctx,
		t,
		uid,
		fl,
		&testItem,
		baseURL,
		"publishFeedItem",
	)
	if err != nil {
		t.Errorf("can't post test item: %s", err)
		return
	}
	headers := getDefaultHeaders(t, baseURL)

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

func getInterserviceBearerTokenHeader(t *testing.T, rootDomain string) string {
	isc := getInterserviceClient(t, rootDomain)
	authToken, err := isc.CreateAuthToken()
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

func getDefaultHeaders(t *testing.T, rootDomain string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": getInterserviceBearerTokenHeader(t, rootDomain),
	}
}

func getGraphQLHeaders(t *testing.T) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": getBearerTokenHeader(t),
	}
}

func getBearerTokenHeader(t *testing.T) string {
	ctx := context.Background()
	user, err := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
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

func getInterserviceClient(t *testing.T, rootDomain string) *base.InterServiceClient {
	service := base.ISCService{
		Name:       "feed",
		RootDomain: rootDomain,
	}
	isc, err := base.NewInterserviceClient(service)
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
	fl base.Flavour,
	el base.Element,
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

	headers := getDefaultHeaders(t, baseURL)
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
	fl base.Flavour,
	el base.Element,
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

	headers := getDefaultHeaders(t, baseURL)
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

func getTestItem() base.Item {
	return base.Item{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       base.TextTypePlain,
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
			base.GetYoutubeVideoLink(base.SampleVideoURL, "title", "description", base.LogoURL),
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
			},
		},
		Conversations: []base.Message{
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
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
}

func testNudge() *base.Nudge {
	return &base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: getTestSequenceNumber(),
		Expiry:         time.Now().Add(time.Hour * 24),
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Title:          ksuid.New().String(),
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
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

func getTestSequenceNumber() int {
	return rand.Intn(intMax)
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
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.LogoURL),
		ActionType:     base.ActionTypePrimary,
		Handling:       base.HandlingFullPage,
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

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestGoogleCloudPubSubHandler(t *testing.T) {
	ctx := context.Background()
	b64 := base64.StdEncoding.EncodeToString([]byte(ksuid.New().String()))
	testPush := base.PubSubPayload{
		Subscription: ksuid.New().String(),
		Message: base.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      []byte(b64),
			Attributes: map[string]string{
				"topicID": common.ActionPublishTopic,
			},
		},
	}
	testPushJSON, err := json.Marshal(testPush)
	if err != nil {
		t.Errorf("can't marshal JSON: %s", err)
		return
	}

	idTokenHTTPClient, err := idtoken.NewClient(ctx, base.Aud)
	if err != nil {
		t.Errorf("can't initialize idToken HTTP client: %s", err)
		return
	}

	pubsubURL := fmt.Sprintf("%s%s", baseURL, base.PubSubHandlerPath)
	req, err := http.NewRequest(
		http.MethodPost,
		pubsubURL,
		bytes.NewBuffer(testPushJSON),
	)
	if err != nil {
		t.Errorf("can't initialize request: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

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
			name: "valid pubsub format payload with valid auth",
			args: args{
				r:      req,
				client: idTokenHTTPClient,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "no auth header",
			args: args{
				r:      req,
				client: http.DefaultClient,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		// TODO Bad topic
		// TODO Item publish, good input
		// TODO Item publish, bad input
		// TODO Item resolve, good input
		// TODO Item resolve, bad input
		// TODO Item delete, good input
		// TODO Item delete, bad input
		// TODO Item hide, bad input
		// TODO Item hide, good input
		// TODO Item show, bad input
		// TODO Item show, good input
		// TODO Item pin, bad input
		// TODO Item pin, good input
		// TODO Item unpin, bad input
		// TODO Item unpin, good input
		// TODO Nudge publish, bad input
		// TODO Nudge publish, good input
		// TODO Nudge delete, bad input
		// TODO Nudge delete, good input
		// TODO Nudge resolve, bad input
		// TODO Nudge resolve, good input
		// TODO Nudge unresolve, bad input
		// TODO Nudge unresolve, good input
		// TODO Nudge hide, bad input
		// TODO Nudge hide, good input
		// TODO Nudge show, bad input
		// TODO Nudge show, good input
		// TODO Action publish, bad input
		// TODO Action publish, good input
		// TODO Action delete, bad input
		// TODO Action delete, good input
		// TODO Message post, bad input
		// TODO Message post, good input
		// TODO Message delete, bad input
		// TODO Message delete, good input
		// TODO Incoming event, bad input
		// TODO Incoming event, good input
		// TODO Unknown pubsub topic (fallback)
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
				log.Printf("request JSON: %s", string(testPushJSON))
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
	headers := getDefaultHeaders(t, baseURL)
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
	uploadInput := base.UploadInput{
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
				upload := &base.Upload{}
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
				fetchedUpload := base.Upload{}
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
	fl base.Flavour,
	nudge *base.Nudge,
) error {

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return fmt.Errorf("can't initialize Firebase Repository: %s", err)
	}
	nudge.Status = base.StatusDone
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
	ctx := context.Background()
	uid := xid.New().String()
	fl := base.FlavourConsumer
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize Firebase Repository: %s", err)
	}

	defaultNudges, err := database.SetDefaultNudges(
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
		if nudge.Title == common.AddInsuranceNudgeTitle {
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

	headers := getDefaultHeaders(t, baseURL)

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
					common.AddInsuranceNudgeTitle,
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
}