package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement-service/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.ResolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve a Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "resolveFeedItem", err)

	return item, nil
}

func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.UnresolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unresolveFeedItem", err)

	return item, nil
}

func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.PinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to pin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "pinFeedItem", err)

	return item, nil
}

func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.UnpinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unpin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unpinFeedItem", err)

	return item, nil
}

func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.HideFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideFeedItem", err)

	return item, nil
}

func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.interactor.Feed.ShowFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to show Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showFeedItem", err)

	return item, nil
}

func (r *mutationResolver) HideNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.interactor.Feed.HideNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideNudge", err)

	return nudge, nil
}

func (r *mutationResolver) ShowNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.interactor.Feed.ShowNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showNudge", err)

	return nudge, nil
}

func (r *mutationResolver) PostMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, message feedlib.Message) (*feedlib.Message, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	msg, err := r.interactor.Feed.PostMessage(ctx, uid, flavour, itemID, &message)
	if err != nil {
		return nil, fmt.Errorf("unable to post a message: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "postMessage", err)

	return msg, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, messageID string) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.interactor.Feed.DeleteMessage(ctx, uid, flavour, itemID, messageID)
	if err != nil {
		return false, fmt.Errorf("can't delete message: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "deleteMessage", err)

	return true, nil
}

func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour feedlib.Flavour, event feedlib.Event) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.interactor.Feed.ProcessEvent(ctx, uid, flavour, &event)
	if err != nil {
		return false, fmt.Errorf("can't process event: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "processEvent", err)

	return true, nil
}

func (r *queryResolver) GetFeed(ctx context.Context, flavour feedlib.Flavour, isAnonymous bool, persistent feedlib.BooleanFilter, status *feedlib.Status, visibility *feedlib.Visibility, expired *feedlib.BooleanFilter, filterParams *helpers.FilterParams) (*domain.Feed, error) {
	startTime := time.Now()

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
		return nil, fmt.Errorf("can't get Feed: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFeed", err)

	return feed, nil
}

func (r *queryResolver) Labels(ctx context.Context, flavour feedlib.Flavour) ([]string, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	labels, err := r.interactor.Feed.Labels(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get Labels count: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "labels", err)

	return labels, nil
}

func (r *queryResolver) UnreadPersistentItems(ctx context.Context, flavour feedlib.Flavour) (int, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return -1, fmt.Errorf("can't get logged in user UID")
	}
	count, err := r.interactor.Feed.UnreadPersistentItems(ctx, uid, flavour)
	if err != nil {
		return -1, fmt.Errorf("unable to count unread persistent items: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unreadPersistentItems", err)

	return count, nil
}
