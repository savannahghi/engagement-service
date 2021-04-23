package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *entityResolver) FindAccessTokenByJwt(ctx context.Context, jwt string) (*resources.AccessToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindDummyByID(ctx context.Context, id *string) (*resources.Dummy, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindFeedByID(ctx context.Context, id string) (*domain.Feed, error) {
	r.checkPreconditions()

	components := strings.Split(id, "|")
	if len(components) != 2 {
		return nil, fmt.Errorf(
			"expected `id` to be a string with exactly two parts separated by a | i.e the uid and flavour as `uid|flavour`")
	}

	uid := components[0]
	flavour := base.Flavour(components[1])
	if !flavour.IsValid() {
		return nil, fmt.Errorf("%s is not a valid flavour", flavour)
	}
	anonymous := false
	return r.interactor.Feed.GetFeed(ctx,
		&uid,
		&anonymous,
		flavour,
		base.BooleanFilterBoth,
		nil,
		nil,
		nil,
		nil,
	)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
