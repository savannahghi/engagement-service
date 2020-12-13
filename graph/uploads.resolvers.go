package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/engagement/graph/uploads"
)

func (r *mutationResolver) Upload(ctx context.Context, input *uploads.UploadInput) (*uploads.Upload, error) {
	r.checkPreconditions()
	return r.uploadService.Upload(ctx, input)
}
