package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *mutationResolver) SendNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification base.FirebaseSimpleNotificationInput, android *base.FirebaseAndroidConfigInput, ios *base.FirebaseAPNSConfigInput, web *base.FirebaseWebpushConfigInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notificationData, err := base.MapInterfaceToMapString(data)
	if err != nil {
		return false, err
	}

	sent, err := r.interactor.FCM.SendNotification(
		ctx,
		registrationTokens,
		notificationData,
		&notification,
		android,
		ios,
		web,
	)
	if err != nil {
		return false, fmt.Errorf("failed to send a notification : %w", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sendNotification",
		err,
	)

	return sent, nil
}

func (r *queryResolver) Notifications(ctx context.Context, registrationToken string, newerThan time.Time, limit int) ([]*resources.SavedNotification, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notification, err := r.interactor.FCM.Notifications(ctx, registrationToken, newerThan, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %w", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"notifications",
		err,
	)

	return notification, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
