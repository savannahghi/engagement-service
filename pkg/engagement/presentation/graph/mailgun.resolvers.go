package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) SimpleEmail(ctx context.Context, subject string, text string, to []string) (string, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	return r.interactor.Mail.SimpleEmail(subject, text, to...)
}
