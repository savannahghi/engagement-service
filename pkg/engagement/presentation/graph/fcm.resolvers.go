package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) SendNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notificationData, err := converterandformatter.MapInterfaceToMapString(data)
	if err != nil {
		return false, err
	}

	sent, err := r.usecases.SendNotification(
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

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sendNotification",
		err,
	)

	return sent, nil
}

func (r *mutationResolver) SendFCMByPhoneOrEmail(ctx context.Context, phoneNumber *string, email *string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	sent, err := r.usecases.SendFCMByPhoneOrEmail(
		ctx,
		phoneNumber,
		email,
		data,
		notification,
		android,
		ios,
		web,
	)
	if err != nil {
		return false, fmt.Errorf("failed to send an FCM notification by email or phone : %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"SendFCMByPhoneOrEmail",
		err,
	)

	return sent, nil
}

func (r *queryResolver) Notifications(ctx context.Context, registrationToken string, newerThan time.Time, limit int) ([]*dto.SavedNotification, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notification, err := r.usecases.Notifications(ctx, registrationToken, newerThan, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
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
