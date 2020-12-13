//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"context"
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	"gitlab.slade360emr.com/go/engagement/graph/library"
	"gitlab.slade360emr.com/go/engagement/graph/scheduling"
	"gitlab.slade360emr.com/go/engagement/graph/uploads"
)

// NewResolver sets up the dependencies needed for WhatsApp query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	fr feed.Repository,
	ns feed.NotificationService,
) (*Resolver, error) {
	return &Resolver{
		repository:          fr,
		notificationService: ns,
		libraryService:      library.NewService(),
		uploadService:       uploads.NewService(),
		schedulingService:   scheduling.NewService(),
	}, nil
}

// Resolver serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	repository          feed.Repository
	notificationService feed.NotificationService
	libraryService      *library.Service
	uploadService       *uploads.Service
	schedulingService   *scheduling.Service
}

func (r Resolver) checkPreconditions() {
	if r.repository == nil {
		log.Panicf("nil repository in resolver")
	}

	if r.notificationService == nil {
		log.Panicf("nil notification service in resolver")
	}

	if r.libraryService == nil {
		log.Panicf("nil library service in resolver")
	}

	if r.uploadService == nil {
		log.Panicf("nil upload service in resolver")
	}

	if r.schedulingService == nil {
		log.Panicf("nil scheduling service in resolver")
	}
}

func (r Resolver) getLoggedInUserUID(ctx context.Context) (string, error) {
	authToken, err := base.GetUserTokenFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("auth token not found in context: %w", err)
	}
	return authToken.UID, nil
}

func (r Resolver) getThinFeed(ctx context.Context, flavour feed.Flavour) (*feed.Feed, error) {
	r.checkPreconditions()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}

	agg, err := feed.NewCollection(r.repository, r.notificationService)
	if err != nil {
		return nil, fmt.Errorf("can't initialize feed aggregate")
	}

	anonymous := false

	thinFeed, err := agg.GetThinFeed(ctx, &uid, &anonymous, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed, nil
}
