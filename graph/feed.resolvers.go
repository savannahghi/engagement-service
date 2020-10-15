package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*feed.GhostCMSPost, error) {
	return r.feedService.GetLibraryContent(ctx)
}

func (r *queryResolver) GetFaqsContent(ctx context.Context) ([]*feed.GhostCMSPost, error) {
	return r.feedService.GetFaqsContent(ctx)
}

func (r *queryResolver) GetFeedContent(ctx context.Context) ([]*feed.GhostCMSPost, error) {
	return r.feedService.GetFeedContent(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
