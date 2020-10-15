package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *entityResolver) FindGhostCMSPostByID(ctx context.Context, id *string) (*feed.GhostCMSPost, error) {
	return nil, nil
}

func (r *entityResolver) FindGhostCMSTagByID(ctx context.Context, id *string) (*feed.GhostCMSTag, error) {
	return nil, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
