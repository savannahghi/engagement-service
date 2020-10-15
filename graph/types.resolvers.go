package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *ghostCMSPostResolver) CreatedAt(ctx context.Context, obj *feed.GhostCMSPost) (*string, error) {
	return nil, nil
}

func (r *ghostCMSPostResolver) PublishedAt(ctx context.Context, obj *feed.GhostCMSPost) (*string, error) {
	return nil, nil
}

// GhostCMSPost returns generated.GhostCMSPostResolver implementation.
func (r *Resolver) GhostCMSPost() generated.GhostCMSPostResolver { return &ghostCMSPostResolver{r} }

type ghostCMSPostResolver struct{ *Resolver }
