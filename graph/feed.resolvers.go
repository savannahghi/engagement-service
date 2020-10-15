package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/feed/graph/generated"
	"gitlab.slade360emr.com/go/feed/graph/model"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*model.LibraryItem, error) {
	r.checkPreconditions()
	return r.feedService.GetLibraryContent(ctx)
}

func (r *queryResolver) GetFaqs(ctx context.Context) ([]*model.Faq, error) {
	r.checkPreconditions()
	return r.feedService.GetFaqs(ctx)
}

func (r *queryResolver) GetFeedItems(ctx context.Context) ([]*model.FeedItem, error) {
	r.checkPreconditions()
	return r.feedService.GetFeedItems(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
