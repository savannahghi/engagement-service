package feed

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// Feed service constants
const (
	ghostCMSAPIEndpoint = "GHOST_CMS_API_ENDPOINT"

	ghostCMSAPIKey = "GHOST_CMS_API_KEY"

	apiRoot = "/ghost/api/v3/content/posts/?"

	includeTags = "&include=tags"

	allowedFeedTagFilter = "&filter=tag:welcome&filter=tag:how-to&filter=tag:what-is&filter=tag:getting-started"

	allowedFAQsTagFilter = "&filter=tag:faqs&filter=tag:how-to"

	allowedLibraryTagFilter = "&filter=tag:diet&filter=tag:health-tips"
)

type requestType int

const (
	feedRequest requestType = iota + 1
	faqsRequest
	libraryRequest
)

// NewService creates a new Feed Service
func NewService() *Service {
	e := base.MustGetEnvVar(ghostCMSAPIEndpoint)
	a := base.MustGetEnvVar(ghostCMSAPIKey)
	return &Service{
		APIEndpoint:  e,
		APIKey:       a,
		PostsAPIRoot: fmt.Sprintf("%v%vkey=%v", e, apiRoot, a),
	}
}

// Service organizes Feed functionality
// APIEndpoint should be of the form https://<name>.ghost.io
type Service struct {
	APIEndpoint  string
	APIKey       string
	PostsAPIRoot string
}

func (s Service) checkPreconditions() {
	if s.APIEndpoint == "" {
		log.Panicf("Ghost API endpoint must be present")
	}

	if s.APIKey == "" {
		log.Panicf("Ghost API key must be present")
	}
}

func (s Service) composeRequest(reqType requestType) *string {
	var urlRequest string
	switch reqType {
	case feedRequest:
		urlRequest = fmt.Sprintf("%v%v%v", s.PostsAPIRoot, includeTags, allowedFeedTagFilter)
	case faqsRequest:
		urlRequest = fmt.Sprintf("%v%v%v", s.PostsAPIRoot, includeTags, allowedFAQsTagFilter)
	case libraryRequest:
		urlRequest = fmt.Sprintf("%v%v%v", s.PostsAPIRoot, includeTags, allowedLibraryTagFilter)
	}

	return &urlRequest
}

// GetFeedContent fetches posts that populate the feed. Since the feed right now is naive,
// we just dump what we get. However, posts for the feed will be of specific tags. Check `allowedFeedTagFilter` above
func (s Service) GetFeedContent(ctx context.Context) ([]*GhostCMSPost, error) {
	s.checkPreconditions()
	url := s.composeRequest(feedRequest)
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create action request with error; %v", err)
	}

	c := &http.Client{Timeout: time.Second * 300}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occured when posting to %v with err %v", url, err)
	}

	defer resp.Body.Close()

	var rr GhostCMSServerResponse

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return nil, fmt.Errorf("failed to decoder response with err %v", err)
	}

	return rr.Posts, nil
}

// GetFaqsContent fetech content of frequently asked question
func (s Service) GetFaqsContent(ctx context.Context) ([]*GhostCMSPost, error) {
	s.checkPreconditions()
	url := s.composeRequest(faqsRequest)
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create action request with error; %v", err)
	}

	c := &http.Client{Timeout: time.Second * 300}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occured when posting to %v with err %v", url, err)
	}

	defer resp.Body.Close()

	var rr GhostCMSServerResponse

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return nil, fmt.Errorf("failed to decoder response with err %v", err)
	}

	return rr.Posts, nil
}

// GetLibraryContent gets library content to be show under libary section of the app
func (s Service) GetLibraryContent(ctx context.Context) ([]*GhostCMSPost, error) {
	s.checkPreconditions()
	url := s.composeRequest(libraryRequest)
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create action request with error; %v", err)
	}

	c := &http.Client{Timeout: time.Second * 300}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occured when posting to %v with err %v", url, err)
	}

	defer resp.Body.Close()

	var rr GhostCMSServerResponse

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return nil, fmt.Errorf("failed to decoder response with err %v", err)
	}

	return rr.Posts, nil
}
