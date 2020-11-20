package library

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// Library service constants
const (
	ghostCMSAPIEndpoint     = "GHOST_CMS_API_ENDPOINT"
	ghostCMSAPIKey          = "GHOST_CMS_API_KEY"
	apiRoot                 = "/ghost/api/v3/content/posts/?"
	includeTags             = "&include=tags"
	includeAuthors          = "&include=authors"
	formats                 = "&formats=html,plaintext"
	allowedFeedTagFilter    = "&filter=tag:welcome&filter=tag:how-to&filter=tag:what-is&filter=tag:getting-started"
	allowedFAQsTagFilter    = "&filter=tag:faqs&filter=tag:how-to"
	allowedLibraryTagFilter = "&filter=tag:diet&filter=tag:health-tips"
)

type requestType int

const (
	feedRequest requestType = iota + 1
	faqsRequest
	libraryRequest
)

// NewService creates a new library Service
func NewService() *Service {
	e := base.MustGetEnvVar(ghostCMSAPIEndpoint)
	a := base.MustGetEnvVar(ghostCMSAPIKey)
	srv := &Service{
		APIEndpoint:  e,
		APIKey:       a,
		PostsAPIRoot: fmt.Sprintf("%v%vkey=%v", e, apiRoot, a),
	}
	srv.checkPreconditions()
	return srv
}

// Service organizes library functionality
// APIEndpoint should be of the form https://<name>.ghost.io
type Service struct {
	APIEndpoint  string
	APIKey       string
	PostsAPIRoot string
}

func (s Service) checkPreconditions() {
	if s.APIEndpoint == "" {
		log.Panicf("Ghost API endpoint must be set")
	}

	if s.APIKey == "" {
		log.Panicf("Ghost API key must be set")
	}

	if s.PostsAPIRoot == "" {
		log.Panicf("Ghost Post API root must be set")
	}
}

func (s Service) composeRequest(reqType requestType) string {
	var urlRequest string
	switch reqType {
	case feedRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedFeedTagFilter,
			includeAuthors,
			formats,
		)
	case faqsRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedFAQsTagFilter,
			includeAuthors,
			formats,
		)
	case libraryRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedLibraryTagFilter,
			includeAuthors,
			formats,
		)
	}
	return urlRequest
}

func (s Service) getCMSPosts(ctx context.Context, requestType requestType) ([]*GhostCMSPost, error) {
	s.checkPreconditions()
	url := s.composeRequest(feedRequest)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

// GetFeedContent fetches posts that should be added to the feed.
func (s Service) GetFeedContent(ctx context.Context) ([]*GhostCMSPost, error) {
	return s.getCMSPosts(ctx, feedRequest)
}

// GetFaqsContent fetches posts tagged as FAQs.
func (s Service) GetFaqsContent(ctx context.Context) ([]*GhostCMSPost, error) {
	return s.getCMSPosts(ctx, faqsRequest)
}

// GetLibraryContent gets library content to be show under libary section of the app.
func (s Service) GetLibraryContent(ctx context.Context) ([]*GhostCMSPost, error) {
	return s.getCMSPosts(ctx, libraryRequest)
}
