package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
)

func (r *mutationResolver) RecordNPSResponse(ctx context.Context, input resources.NPSInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	response, err := r.interactor.Surveys.RecordNPSResponse(ctx, input)
	if err != nil {
		return false, fmt.Errorf("failed to record nps response")
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"recordNPSResponse",
		err,
	)

	return response, nil

}

func (r *queryResolver) ListNPSResponse(ctx context.Context) ([]*resources.NPSResponse, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	return nil, nil
}
