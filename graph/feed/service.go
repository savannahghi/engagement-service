package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/feed/graph/model"
)

// Feed service constants
const (
	requestTimeoutSeconds  = 30
	strapiGraphQlServerURL = "STRAPI_GRAPHQL_SERVER_URL"
)

// NewService creates a new Feed Service
func NewService() *Service {

	cmsServerURL := base.MustGetEnvVar(strapiGraphQlServerURL)
	return &Service{
		cmsServerURL: cmsServerURL,
	}
}

// Service organizes Feed functionality
type Service struct {
	cmsServerURL string
}

func (s Service) checkPreconditions() {
	if s.cmsServerURL == "" {
		log.Panicf("the feed Service has an emty CMS Server URL")
	}
}

func (s Service) newRequest(url string, method string, body []byte) (*http.Response, error) {
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Length", strconv.Itoa(buffer.Len()))

	client := &http.Client{
		Timeout: time.Second * requestTimeoutSeconds,
	}
	return client.Do(request)
}

// GetFaqs Fetches Faqs from CMS
func (s Service) GetFaqs(ctx context.Context) ([]*model.Faq, error) {
	jsonQuery := map[string]string{
		"query": `
            { 
                faqs {
					question,
					answer
                }
            }
        `,
	}
	jsonQueryValue, _ := json.Marshal(jsonQuery)
	response, err := s.newRequest(s.cmsServerURL, "POST", jsonQueryValue)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	// jsonByte, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// err = json.Unmarshal(jsonByte, faqResp)
	// if err != nil {
	// 	panic(err)
	// }

	faqItems := []*model.Faq{}
	return faqItems, nil

}

// GetLibraryContent Fetches Library Content from CMS
func (s Service) GetLibraryContent(ctx context.Context) ([]*model.LibraryItem, error) {
	jsonQuery := map[string]string{
		"query": `
            { 
                libraryItem {
					title,
					description
                }
            }
        `,
	}
	jsonQueryValue, _ := json.Marshal(jsonQuery)
	response, err := s.newRequest(s.cmsServerURL, "POST", jsonQueryValue)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

	libraryContent := make([]*model.LibraryItem, 0)

	return libraryContent, nil

}

// GetFeedItems Fetches Feed Items from CMS
func (s Service) GetFeedItems(ctx context.Context) ([]*model.FeedItem, error) {
	jsonQuery := map[string]string{
		"query": `
            { 
                feedItems {
					title,
					description
                }
            }
        `,
	}
	jsonQueryValue, _ := json.Marshal(jsonQuery)
	response, err := s.newRequest(s.cmsServerURL, "POST", jsonQueryValue)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

	feedItems := make([]*model.FeedItem, 0)

	return feedItems, nil
}
