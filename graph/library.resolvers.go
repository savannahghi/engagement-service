package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/engagement/graph/generated"
	"gitlab.slade360emr.com/go/engagement/graph/library"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.libraryService.GetLibraryContent(ctx)
}

func (r *queryResolver) GetFaqsContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.libraryService.GetFaqsContent(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
