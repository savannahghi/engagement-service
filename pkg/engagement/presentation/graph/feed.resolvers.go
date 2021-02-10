package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.ResolveFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.UnresolveFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.PinFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.UnpinFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.HideFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour base.Flavour, itemID string) (*base.Item, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.ShowFeedItem(ctx, uid, flavour, itemID)
}

func (r *mutationResolver) HideNudge(ctx context.Context, flavour base.Flavour, nudgeID string) (*base.Nudge, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.HideNudge(ctx, uid, flavour, nudgeID)
}

func (r *mutationResolver) ShowNudge(ctx context.Context, flavour base.Flavour, nudgeID string) (*base.Nudge, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.ShowNudge(ctx, uid, flavour, nudgeID)
}

func (r *mutationResolver) PostMessage(ctx context.Context, flavour base.Flavour, itemID string, message base.Message) (*base.Message, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.PostMessage(ctx, uid, flavour, itemID, &message)
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour base.Flavour, itemID string, messageID string) (bool, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.interactor.Feed.DeleteMessage(ctx, uid, flavour, itemID, messageID)
	if err != nil {
		return false, fmt.Errorf("can't delete message: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour base.Flavour, event base.Event) (bool, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.interactor.Feed.ProcessEvent(ctx, uid, flavour, &event)
	if err != nil {
		return false, fmt.Errorf("can't process event: %w", err)
	}

	return true, nil
}

func (r *queryResolver) GetFeed(ctx context.Context, flavour base.Flavour, isAnonymous bool, persistent base.BooleanFilter, status *base.Status, visibility *base.Visibility, expired *base.BooleanFilter, filterParams *helpers.FilterParams) (*domain.Feed, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	feed, err := r.interactor.Feed.GetFeed(
		ctx,
		&uid,
		&isAnonymous,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("can't process event: %w", err)
	}
	return feed, nil
}

func (r *queryResolver) Labels(ctx context.Context, flavour base.Flavour) ([]string, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.Labels(ctx, uid, flavour)
}

func (r *queryResolver) UnreadPersistentItems(ctx context.Context, flavour base.Flavour) (int, error) {
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return -1, fmt.Errorf("can't get logged in user UID")
	}
	return r.interactor.Feed.UnreadPersistentItems(ctx, uid, flavour)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
