package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
)

func (r *mutationResolver) Send(ctx context.Context, to string, message string) (*resources.SendMessageResponse, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	return r.interactor.SMS.Send(to, message)
}

func (r *mutationResolver) SendToMany(ctx context.Context, message string, to []string) (*resources.SendMessageResponse, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	return r.interactor.SMS.SendToMany(message, to)
}
