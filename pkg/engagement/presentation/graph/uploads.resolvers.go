package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
)

func (r *mutationResolver) Upload(ctx context.Context, input base.UploadInput) (*base.Upload, error) {
	startTime := time.Now()

	r.checkPreconditions()
	upload, err := r.interactor.Uploads.Upload(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to upload: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"upload",
		err,
	)

	return upload, nil
}

func (r *queryResolver) FindUploadByID(ctx context.Context, id string) (*base.Upload, error) {
	startTime := time.Now()

	r.checkPreconditions()
	upload, err := r.interactor.Uploads.FindUploadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to find upload by ID: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"findUploadByID",
		err,
	)

	return upload, nil
}
