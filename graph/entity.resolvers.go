package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *entityResolver) FindFeedByUIDAndFlavour(ctx context.Context, uid string, flavour feed.Flavour) (*feed.Feed, error) {
	r.checkPreconditions()

	agg, err := feed.NewCollection(r.repository, r.notificationService)
	if err != nil {
		return nil, fmt.Errorf("can't initialize feed aggregate")
	}

	feed, err := agg.GetFeed(
		ctx,
		uid,
		flavour,
		feed.BooleanFilterBoth,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return feed, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
