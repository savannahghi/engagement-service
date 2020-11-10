package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *entityResolver) FindActionByIDAndSequenceNumber(ctx context.Context, id string, sequenceNumber int) (*feed.Action, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindEventByID(ctx context.Context, id string) (*feed.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

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

func (r *entityResolver) FindItemByIDAndSequenceNumber(ctx context.Context, id string, sequenceNumber int) (*feed.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindNudgeByIDAndSequenceNumber(ctx context.Context, id string, sequenceNumber int) (*feed.Nudge, error) {
	panic(fmt.Errorf("not implemented"))
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
