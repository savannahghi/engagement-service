package feed

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// Feed service constants
const (
	apiRoot = "/ghost/api/v3/content/posts/?"

	ghostCMSAPIEndpoint = "GHOST_CMS_API_ENDPOINT"

	ghostCMSAPIKey = "GHOST_CMS_API_KEY"

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
// APIEndpoint should be of the form https://bewell.ghost.io
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
	_, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatalf("Failed to create action request with error; %v", err)
		// fail silently
		return nil, nil
	}

	// continue

	return nil, nil
}
