package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

func (r *mutationResolver) Upload(ctx context.Context, input base.UploadInput) (*base.Upload, error) {
	r.checkPreconditions()
	return r.uploadService.Upload(ctx, input)
}

func (r *queryResolver) FindUploadByID(ctx context.Context, id string) (*base.Upload, error) {
	r.checkPreconditions()
	return r.uploadService.FindUploadByID(ctx, id)
}
