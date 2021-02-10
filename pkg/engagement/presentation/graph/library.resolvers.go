package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.interactor.Library.GetLibraryContent(ctx)
}

func (r *queryResolver) GetFaqsContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	return r.interactor.Library.GetFaqsContent(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
