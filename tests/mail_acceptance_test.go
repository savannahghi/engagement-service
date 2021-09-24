package tests

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httputil"
// 	"testing"
// 	"time"
// )

// func TestGraphQLSimpleEmail(t *testing.T) {
// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers := getGraphQLHeaders(t)
// 	testUserMail := "test@bewell.co.ke"

// 	graphqlMutation := `mutation simpleEmail(
// 		$subject: String!,
// 		$text:String!,
// 		$to:[String!]!,
// 		){
// 			simpleEmail(
// 				subject: $subject,
// 				text:$text,
// 				to:$to,
// 			)
// 	  }
// 	`

// 	type args struct {
// 		query map[string]interface{}
// 	}

// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantStatus int
// 		wantErr    bool
// 	}{
// 		{
// 			name: "valid query",
// 			args: args{
// 				query: map[string]interface{}{
// 					"query": graphqlMutation,
// 					"variables": map[string]interface{}{
// 						"subject": "Test Subject",
// 						"text":    "Hey :)",
// 						"to":      []string{testUserMail},
// 					},
// 				},
// 			},
// 			wantStatus: http.StatusOK,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "invalid query - Using invalid payload",
// 			args: args{
// 				query: map[string]interface{}{
// 					"query": graphqlMutation,
// 					"variables": map[string]interface{}{
// 						"some-invalid": "data",
// 					},
// 				},
// 			},
// 			wantStatus: http.StatusUnprocessableEntity,
// 			wantErr:    true,
// 		},

// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			body, err := mapToJSONReader(tt.args.query)

// 			if err != nil {
// 				t.Errorf("unable to get GQL JSON io Reader: %s", err)
// 				return
// 			}

// 			r, err := http.NewRequest(
// 				http.MethodPost,
// 				graphQLURL,
// 				body,
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
// 			client := http.Client{
// 				Timeout: time.Second * testHTTPClientTimeout,
// 			}
// 			resp, err := client.Do(r)
// 			if err != nil {
// 				t.Errorf("request error: %s", err)
// 				return
// 			}

// 			dataResponse, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Errorf("can't read request body: %s", err)
// 				return
// 			}
// 			if dataResponse == nil {
// 				t.Errorf("nil response data")
// 				return
// 			}

// 			data := map[string]interface{}{}
// 			err = json.Unmarshal(dataResponse, &data)
// 			if err != nil {
// 				t.Errorf("bad data returned")
// 				return
// 			}

// 			if tt.wantErr {
// 				errMsg, ok := data["errors"]
// 				if !ok {
// 					t.Errorf("GraphQL error: %s", errMsg)
// 					return
// 				}
// 			}

// 			if !tt.wantErr {
// 				_, ok := data["errors"]
// 				if ok {
// 					t.Errorf("error not expected")
// 					return
// 				}
// 			}
// 			if tt.wantStatus != resp.StatusCode {
// 				b, _ := httputil.DumpResponse(resp, true)
// 				t.Errorf("Bad status response returned; %v ", string(b))
// 				return
// 			}
// 		})
// 	}
// }
