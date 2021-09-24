package tests

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"testing"

// 	"github.com/savannahghi/firebasetools"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGraphQLGetFeed(t *testing.T) {
// 	ctx := firebasetools.GetAuthenticatedContext(t)
// 	if ctx == nil {
// 		t.Errorf("nil context")
// 		return
// 	}

// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers := getGraphQLHeaders(t)
// 	gql := map[string]interface{}{}
// 	gql["query"] = `
// 	query getFeed($flavour: Flavour!,$isAnonymous: Boolean!,$playMP4:Boolean,
// 		$persistent: BooleanFilter!,
// 		$status: Status,
// 		  $visibility: Visibility,
// 		  $expired: BooleanFilter){
// 	   getFeed(flavour:$flavour,playMP4:$playMP4, isAnonymous:$isAnonymous,
// 		persistent:$persistent, status:$status,
// 		visibility:$visibility, expired:$expired  ){
// 		id
// 		uid
// 		isAnonymous
// 		actions {
// 				  id
// 				  sequenceNumber
// 				  name
// 				  actionType
// 				  handling
// 			allowAnonymous
// 		  }
// 		  nudges {
// 				  id
// 				  sequenceNumber
// 				  visibility
// 				  status
// 				  title
// 				  text
// 				  actions {
// 					  id
// 					  sequenceNumber
// 					  name
// 					  actionType
// 					  handling
// 					  allowAnonymous
// 				  }
// 				  groups
// 				  users
// 				  links {
// 					  id
// 					  url
// 					  linkType
// 				  }
// 				  notificationChannels
// 		  }

// 		  items {
// 				  id
// 				  sequenceNumber
// 				  expiry
// 				  persistent
// 				  status
// 				  visibility
// 				  icon {
// 					  id
// 					  url
// 					  linkType
// 				  }
// 				  author
// 				  tagline
// 				  label
// 				  timestamp
// 				  summary
// 				  text
// 				  links {
// 					  id
// 					  url
// 					  linkType
// 					}
// 				  actions {
// 					  id
// 					  sequenceNumber
// 					  name
// 					  actionType
// 					  handling
// 				  }
// 				  conversations {
// 					  id
// 					  sequenceNumber
// 					  text
// 					  replyTo
// 					  postedByUID
// 					  postedByName
// 				  }
// 				  users
// 				  groups
// 				  notificationChannels
// 			  }
// 	   }
// 	  }
// 	 `

// 	gql["variables"] = map[string]interface{}{
// 		"flavour":     "CONSUMER",
// 		"playMP4":     true,
// 		"isAnonymous": false,
// 		"persistent":  "BOTH",
// 		"status":      "PENDING",
// 		"visibility":  "SHOW",
// 		"expired":     "FALSE",
// 		"filterParams": map[string]interface{}{
// 			"labels": []string{"a_label", "another_label"},
// 		},
// 	}

// 	validQueryReader, err := mapToJSONReader(gql)
// 	if err != nil {
// 		t.Errorf("unable to get GQL JSON io Reader: %s", err)
// 		return
// 	}

// 	type args struct {
// 		body io.Reader
// 	}
// 	tests := []struct {
// 		name               string
// 		args               args
// 		wantStatus         int
// 		wantErr            bool
// 		wantNonZeroItems   bool
// 		wantNonZeroNudges  bool
// 		wantNonZeroActions bool
// 	}{
// 		{
// 			name: "valid query",
// 			args: args{
// 				body: validQueryReader,
// 			},
// 			wantStatus: 200,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "valid query",
// 			args: args{
// 				body: nil,
// 			},
// 			wantStatus: 400,
// 			wantErr:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r, err := http.NewRequest(
// 				http.MethodPost,
// 				graphQLURL,
// 				tt.args.body,
// 			)
// 			if err != nil {
// 				t.Errorf("unable to compose request: %s", err)
// 				return
// 			}

// 			if r == nil {
// 				t.Errorf("nil request")
// 				return
// 			}

// 			for k, v := range headers {
// 				r.Header.Add(k, v)
// 			}
// 			client := http.DefaultClient
// 			resp, err := client.Do(r)
// 			if err != nil {
// 				t.Errorf("request error: %s", err)
// 				return
// 			}

// 			if resp == nil && !tt.wantErr {
// 				t.Errorf("nil response")
// 				return
// 			}

// 			data, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Errorf("can't read request body: %s", err)
// 				return
// 			}
// 			assert.NotNil(t, data)
// 			if data == nil {
// 				t.Errorf("nil response data")
// 				return
// 			}

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			assert.Equal(t, tt.wantStatus, resp.StatusCode)
// 		})
// 	}
// }

// func TestGraphQLGetFaqsContent(t *testing.T) {
// 	ctx := firebasetools.GetAuthenticatedContext(t)
// 	if ctx == nil {
// 		t.Errorf("nil context")
// 		return
// 	}

// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers := getGraphQLHeaders(t)
// 	gql := map[string]interface{}{}

// 	gql["query"] = `
// 	query GetFAQContent($flavour: Flavour!) {
// 		getFaqsContent(flavour: $flavour){
// 		  id,
// 		  slug
// 		  uuid,
// 		  title,
// 		  html,
// 		  excerpt,
// 		  url,
// 		  featureImage,
// 		  readingTime,
// 		  tags {
// 			id,
// 			name,
// 			slug,
// 			description,
// 			visibility,
// 			url
// 		  },
// 		  createdAt,
// 		  updatedAt,
// 		  commentID
// 		}
// 	  }
// 	`

// 	gql["variables"] = map[string]interface{}{
// 		"flavour": "CONSUMER",
// 	}

// 	validQueryReader, err := mapToJSONReader(gql)
// 	if err != nil {
// 		t.Errorf("unable to get GQL JSON io Reader: %s", err)
// 		return
// 	}

// 	type args struct {
// 		body io.Reader
// 	}

// 	tests := []struct {
// 		name               string
// 		args               args
// 		wantStatus         int
// 		wantErr            bool
// 		wantNonZeroItems   bool
// 		wantNonZeroNudges  bool
// 		wantNonZeroActions bool
// 	}{
// 		{
// 			name: "valid query",
// 			args: args{
// 				body: validQueryReader,
// 			},
// 			wantStatus: 200,
// 			wantErr:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r, err := http.NewRequest(
// 				http.MethodPost,
// 				graphQLURL,
// 				tt.args.body,
// 			)
// 			if err != nil {
// 				t.Errorf("unable to compose request: %s", err)
// 				return
// 			}

// 			if r == nil {
// 				t.Errorf("nil request")
// 				return
// 			}

// 			for k, v := range headers {
// 				r.Header.Add(k, v)
// 			}
// 			client := http.DefaultClient
// 			resp, err := client.Do(r)
// 			if err != nil {
// 				t.Errorf("request error: %s", err)
// 				return
// 			}

// 			if resp == nil && !tt.wantErr {
// 				t.Errorf("nil response")
// 				return
// 			}

// 			data, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Errorf("can't read request body: %s", err)
// 				return
// 			}
// 			assert.NotNil(t, data)
// 			if data == nil {
// 				t.Errorf("nil response data")
// 				return
// 			}

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			assert.Equal(t, tt.wantStatus, resp.StatusCode)
// 		})
// 	}

// }

// func TestGraphQLGetLibraryContent(t *testing.T) {
// 	ctx := firebasetools.GetAuthenticatedContext(t)
// 	if ctx == nil {
// 		t.Errorf("nil context")
// 		return
// 	}

// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers := getGraphQLHeaders(t)
// 	gql := map[string]interface{}{}

// 	gql["query"] = `
// 	query GetLibraryContent {
// 		getLibraryContent{
// 		  id,
// 		  slug
// 		  uuid,
// 		  title,
// 		  html,
// 		  excerpt,
// 		  url,
// 		  featureImage,
// 		  readingTime,
// 		  tags {
// 			id,
// 			name,
// 			slug,
// 			description,
// 			visibility,
// 			url
// 		  },
// 		  createdAt,
// 		  updatedAt,
// 		  commentID
// 		}
// 	  }
// 	`
// 	validQueryReader, err := mapToJSONReader(gql)
// 	if err != nil {
// 		t.Errorf("unable to get GQL JSON io Reader: %s", err)
// 		return
// 	}

// 	type args struct {
// 		body io.Reader
// 	}

// 	tests := []struct {
// 		name               string
// 		args               args
// 		wantStatus         int
// 		wantErr            bool
// 		wantNonZeroItems   bool
// 		wantNonZeroNudges  bool
// 		wantNonZeroActions bool
// 	}{
// 		{
// 			name: "valid query",
// 			args: args{
// 				body: validQueryReader,
// 			},
// 			wantStatus: 200,
// 			wantErr:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r, err := http.NewRequest(
// 				http.MethodPost,
// 				graphQLURL,
// 				tt.args.body,
// 			)
// 			if err != nil {
// 				t.Errorf("unable to compose request: %s", err)
// 				return
// 			}

// 			if r == nil {
// 				t.Errorf("nil request")
// 				return
// 			}

// 			for k, v := range headers {
// 				r.Header.Add(k, v)
// 			}
// 			client := http.DefaultClient
// 			resp, err := client.Do(r)
// 			if err != nil {
// 				t.Errorf("request error: %s", err)
// 				return
// 			}

// 			if resp == nil && !tt.wantErr {
// 				t.Errorf("nil response")
// 				return
// 			}

// 			data, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Errorf("can't read request body: %s", err)
// 				return
// 			}
// 			assert.NotNil(t, data)
// 			if data == nil {
// 				t.Errorf("nil response data")
// 				return
// 			}

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			assert.Equal(t, tt.wantStatus, resp.StatusCode)
// 		})
// 	}

// }
