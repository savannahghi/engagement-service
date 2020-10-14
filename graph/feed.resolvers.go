package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/feed/graph/generated"
	"gitlab.slade360emr.com/go/feed/graph/model"
)

func (r *queryResolver) ContentItems(ctx context.Context) ([]*model.ContentItem, error) {
	r.checkPreconditions()
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
