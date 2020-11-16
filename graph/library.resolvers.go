package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"gitlab.slade360emr.com/go/feed/graph/generated"
	"gitlab.slade360emr.com/go/feed/graph/library"
)

func (r *ghostCMSPostResolver) CreatedAt(ctx context.Context, obj *library.GhostCMSPost) (string, error) {
	if obj != nil {
		return obj.CreatedAt.Format(time.RFC3339), nil
	}
	return time.Now().Format(time.RFC3339), nil
}

func (r *ghostCMSPostResolver) PublishedAt(ctx context.Context, obj *library.GhostCMSPost) (string, error) {
	if obj != nil {
		return obj.PublishedAt.Format(time.RFC3339), nil
	}
	return time.Now().Format(time.RFC3339), nil
}

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.libraryService.GetLibraryContent(ctx)
}

func (r *queryResolver) GetFaqsContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.libraryService.GetFaqsContent(ctx)
}

func (r *queryResolver) GetFeedContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.libraryService.GetFeedContent(ctx)
}

// GhostCMSPost returns generated.GhostCMSPostResolver implementation.
func (r *Resolver) GhostCMSPost() generated.GhostCMSPostResolver { return &ghostCMSPostResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type ghostCMSPostResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
