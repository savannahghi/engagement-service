package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *entityResolver) FindFeedByID(ctx context.Context, id string) (*feed.Feed, error) {
	r.checkPreconditions()

	components := strings.Split(id, "|")
	if len(components) != 2 {
		return nil, fmt.Errorf(
			"expected `id` to be a string with exactly two parts separated by a | i.e the uid and flavour as `uid|flavour`")
	}

	uid := components[0]
	flavour := feed.Flavour(components[1])
	if !flavour.IsValid() {
		return nil, fmt.Errorf("%s is not a valid flavour", flavour)
	}

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
