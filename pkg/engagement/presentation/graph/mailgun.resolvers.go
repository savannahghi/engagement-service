package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
)

func (r *mutationResolver) TestFeature(ctx context.Context) (bool, error) {
	return r.usecases.TestFeature(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
