package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestGraphQLProcessEvent(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ProcessEvent($flavour: Flavour!, $event: EventInput!) {
		processEvent(flavour: $flavour, event: $event)
	}	  
	`

	gql["variables"] = map[string]interface{}{
		"flavour": "CONSUMER",
		"event": map[string]interface{}{
			"name": "TEST_EVENT",
			"context": map[string]string{
				"userID":         "user-1",
				"organizationID": "org-1",
				"locationID":     "location-1",
				"timestamp":      "2020-11-05T03:26:15+00:00",
			},
			"payload": map[string]interface{}{
				"data": map[string]interface{}{
					"some": "stuff",
					"and":  "other stuff",
				},
			},
		},
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
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

func TestGraphQLDeleteMessage(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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
		t.Errorf("can't post item: %s", err)
		return
	}

	testMessage := getTestMessage()
	err = postMessage(ctx, t, uid, fl, &testMessage, baseURL, testItem.ID)
	if err != nil {
		t.Errorf("can't post message: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation DeleteMessage(
		$flavour: Flavour!, 
		$itemID: String!, 
		$messageID: String!
	  ) {
		deleteMessage(flavour: $flavour, itemID: $itemID, messageID: $messageID)
	}
	`
	gql["variables"] = map[string]interface{}{
		"flavour":   fl.String(),
		"itemID":    testItem.ID,
		"messageID": testMessage.ID,
	}
	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLPostMessage(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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
		t.Errorf("can't post item: %s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation PostMessage(
		$flavour: Flavour!, 
		$itemID: String!,
		$message: MsgInput!
	  ) {
		postMessage(flavour: $flavour, itemID: $itemID, message: $message) {
		  id
		  sequenceNumber
		  text
		  replyTo
		  postedByUID
		  postedByName
		}
	}
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
		"message": map[string]string{
			"id":             ksuid.New().String(),
			"text":           ksuid.New().String(),
			"replyTo":        ksuid.New().String(),
			"postedByUID":    ksuid.New().String(),
			"postedByName":   ksuid.New().String(),
			"timestamp":      time.Now().Format(time.RFC3339),
			"sequenceNumber": fmt.Sprintf("%d", time.Now().Unix()),
		},
	}
	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLHideNudge(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation HideNudge($flavour: Flavour!, $nudgeID: String!) {
		hideNudge(flavour: $flavour, nudgeID: $nudgeID) {
		  id
		  sequenceNumber
		  visibility
		  status
		  title
		  text
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  groups
		  users
		  links {
			id
			url
			linkType
		  }
		  notificationChannels
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"nudgeID": testNudge.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLShowNudge(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ShowNudge($flavour: Flavour!, $nudgeID: String!) {
		showNudge(flavour: $flavour, nudgeID: $nudgeID) {
		  id
		  sequenceNumber
		  visibility
		  status
		  title
		  text
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  groups
		  users
		  links {
			id
			url
			linkType
		  }
		  notificationChannels
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"nudgeID": testNudge.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLResolveFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ResolveFeedItem($flavour: Flavour!, $itemID: String!) {
		resolveFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLUnresolveFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation UnresolveFeedItem($flavour: Flavour!, $itemID: String!) {
		unresolveFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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
func TestGraphQLPinFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation PinFeedItem($flavour: Flavour!, $itemID: String!) {
		pinFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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
func TestGraphQLUnpinFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation UnpinFeedItem($flavour: Flavour!, $itemID: String!) {
		unpinFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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
func TestGraphQLHideFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation HideFeedItem($flavour: Flavour!, $itemID: String!) {
		hideFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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
func TestGraphQLShowFeedItem(t *testing.T) {
	ctx, authToken := firebasetools.GetAuthenticatedContextAndToken(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	uid := authToken.UID
	fl := feedlib.FlavourConsumer
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

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	mutation ShowFeedItem($flavour: Flavour!, $itemID: String!) {
		showFeedItem(flavour: $flavour, itemID: $itemID) {
		  id
		  sequenceNumber
		  expiry
		  persistent
		  status
		  visibility
		  author
		  tagline
		  label
		  timestamp
		  summary
		  text
		  users
		  groups
		  notificationChannels
		  links {
			id
			url
			linkType
		  }
		  actions {
			id
			sequenceNumber
			name
			actionType
			handling
		  }
		  conversations {
			id
			sequenceNumber
			text
			replyTo
			postedByUID
			postedByName
		  }
		  icon {
			id
			url
			linkType
		  }
		}
	  }	  
	`
	gql["variables"] = map[string]interface{}{
		"flavour": fl.String(),
		"itemID":  testItem.ID,
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
				r.Header.Add(k, v)
			}
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

func TestGraphQLGetFeed(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	gql := map[string]interface{}{}
	gql["query"] = `
	query getFeed($flavour: Flavour!,$isAnonymous: Boolean!,
		$persistent: BooleanFilter!,
		$status: Status,
		  $visibility: Visibility,
		  $expired: BooleanFilter){
	   getFeed(flavour:$flavour,isAnonymous:$isAnonymous,
		persistent:$persistent, status:$status, 
		visibility:$visibility, expired:$expired  ){
		id
		uid
		isAnonymous
		actions {
				  id
				  sequenceNumber
				  name
				  actionType
				  handling
			allowAnonymous
		  }
		  nudges {
				  id
				  sequenceNumber
				  visibility
				  status
				  title
				  text
				  actions {
					  id
					  sequenceNumber
					  name
					  actionType
					  handling
					  allowAnonymous
				  }
				  groups
				  users
				  links {
					  id
					  url
					  linkType
				  }
				  notificationChannels
		  }
		  
		  items {
				  id
				  sequenceNumber
				  expiry
				  persistent
				  status
				  visibility
				  icon {
					  id
					  url
					  linkType
				  }
				  author
				  tagline
				  label
				  timestamp
				  summary
				  text
				  links {
					  id
					  url
					  linkType
					}
				  actions {
					  id
					  sequenceNumber
					  name
					  actionType
					  handling
				  }
				  conversations {
					  id
					  sequenceNumber
					  text
					  replyTo
					  postedByUID
					  postedByName
				  }
				  users
				  groups
				  notificationChannels
			  }
	   }
	  }	  
	 `

	gql["variables"] = map[string]interface{}{
		"flavour":     "CONSUMER",
		"isAnonymous": false,
		"persistent":  "BOTH",
		"status":      "PENDING",
		"visibility":  "SHOW",
		"expired":     "FALSE",
		"filterParams": map[string]interface{}{
			"labels": []string{"a_label", "another_label"},
		},
	}

	validQueryReader, err := mapToJSONReader(gql)
	if err != nil {
		t.Errorf("unable to get GQL JSON io Reader: %s", err)
		return
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name               string
		args               args
		wantStatus         int
		wantErr            bool
		wantNonZeroItems   bool
		wantNonZeroNudges  bool
		wantNonZeroActions bool
	}{
		{
			name: "valid query",
			args: args{
				body: validQueryReader,
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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

			for k, v := range headers {
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

func TestGraphQLUpload(t *testing.T) {

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)

	graphqlMutation := `mutation upload($input: UploadInput!) {
		upload(input: $input) {
		  id
		  url
		  size
		  hash
		  creation
		  title
		  contentType
		  language
		  base64data
		}
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"title":       "Upload title",
							"contentType": "json",
							"language":    "en",
							"base64data":  "data",
							"filename":    "filename",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query - Using invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"title":       123456,
						"contentType": 123456,
						"language":    123456,
						"base64data":  123456,
						"filename":    123456,
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}

func TestGraphQlGetLibraryContent(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)

	graphqlQuery := `query {
		getLibraryContent {
		  id
		  slug
		  uuid
		  title
		  html
		  excerpt
		  url
		  featureImage
		  readingTime
		  tags {
			id
			name
			slug
			description
			visibility
			url
		  }
		  createdAt
		  publishedAt
		  updatedAt
		  commentID
		}
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"invalid": graphqlQuery,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}

func TestGraphQlGetFaqsContent(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)

	graphqlQuery := `query getFaqsContent($flavour: Flavour!) {
		getFaqsContent(flavour: $flavour) {
		  id
		  slug
		  uuid
		  title
		  html
		  excerpt
		  url
		  featureImage
		  readingTime
		  tags {
			id
			name
			slug
			description
			visibility
			url
		  }
		  createdAt
		  publishedAt
		  updatedAt
		  commentID
		}
	}	  
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query - Consumer",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"flavour": feedlib.FlavourConsumer,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		// TODO: Restore after fixing profile problem
		// {
		// 	name: "valid query - Pro",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlQuery,
		// 			"variables": map[string]interface{}{
		// 				"flavour": feedlib.FlavourPro,
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"invalid": graphqlQuery,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}

func TestGraphQlGetLabels(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)

	graphqlQuery := `query($flavour: Flavour!) {
		labels(flavour: $flavour)
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid - Get labels on pro",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"flavour": "PRO",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid - Get labels on consumer",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"flavour": "CONSUMER",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"invalid": "PRO",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}

func TestGraphQLSimpleEmail(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	testUserMail := "test@bewell.co.ke"

	graphqlMutation := `mutation simpleEmail(
		$subject: String!,
		$text:String!,
		$to:[String!]!,
		){
			simpleEmail(
				subject: $subject,
				text:$text,
				to:$to,
			)
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"subject": "Test Subject",
						"text":    "Hey :)",
						"to":      []string{testUserMail},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query - Using invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"some-invalid": "data",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}

func TestGraphQLSendFCMByPhoneOrEmail(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := getGraphQLHeaders(t)
	testUserPhone := interserviceclient.TestUserPhoneNumber

	testUserEmail := "test@bewell.co.ke"
	imageURL := "https:www.example.com/hey.png"

	graphqlMutation := `mutation sendFCMByPhoneOrEmail(
		$phoneNumber: String
		$email: String
		$data: Map!
		$notification: FirebaseSimpleNotificationInput!
		$android: FirebaseAndroidConfigInput
		$ios: FirebaseAPNSConfigInput
		$web: FirebaseWebpushConfigInput
	  ) {
		sendFCMByPhoneOrEmail(
		  phoneNumber: $phoneNumber
		  email: $email
		  data: $data
		  notification: $notification
		  android: $android
		  ios: $ios
		  web: $web
		)
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Happy case: phone",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phoneNumber": testUserPhone,
						"email":       nil,
						"data": map[string]interface{}{
							"title": "Proof of concept that FCM works",
							"body":  "Testing the hypothesis that FCM works without Pub/Sub :)",
						},
						"notification": map[string]interface{}{
							"title":    "Proof of concept that FCM works",
							"body":     "Testing the hypothesis that FCM works without Pub/Sub :)",
							"imageURL": imageURL,
							"data": map[string]interface{}{
								"key": "value",
							},
						},
						"android": map[string]interface{}{
							"priority": "high",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "Happy case: email",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phoneNumber": nil,
						"email":       testUserEmail,
						"data": map[string]interface{}{
							"title": "Proof of concept that FCM works",
							"body":  "Testing the hypothesis that FCM works without Pub/Sub :)",
						},
						"notification": map[string]interface{}{
							"title":    "Proof of concept that FCM works",
							"body":     "Testing the hypothesis that FCM works without Pub/Sub :)",
							"imageURL": imageURL,
							"data": map[string]interface{}{
								"key": "value",
							},
						},
						"android": map[string]interface{}{
							"priority": "high",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},

		{
			name: "Bad case",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"some": "key",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)

			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				b, _ := httputil.DumpResponse(resp, true)
				t.Errorf("Bad status response returned; %v ", string(b))
				return
			}
		})
	}
}
