package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *ghostCMSPostResolver) ReadingTime(ctx context.Context, obj *feed.GhostCMSPost) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
