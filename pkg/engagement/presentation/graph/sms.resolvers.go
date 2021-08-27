package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	smsResponse, err := r.usecases.Send(ctx, to, message)
	if err != nil {
		return nil, fmt.Errorf("unable send SMS: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"send",
		err,
	)

	return smsResponse, nil
}

func (r *mutationResolver) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	smsResponse, err := r.usecases.SendToMany(ctx, message, to)
	if err != nil {
		return nil, fmt.Errorf("unable to send SMS to many: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sendToMany",
		err,
	)

	return smsResponse, nil
}
