package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
	"github.com/savannahghi/serverutils"
)

func (r *queryResolver) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	accessToken, err := r.interactor.Twilio.TwilioAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to generate access token: %w", err)
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"twilioAccessToken",
		err,
	)

	return accessToken, nil
}
