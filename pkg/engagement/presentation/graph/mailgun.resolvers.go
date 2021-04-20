package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"gitlab.slade360emr.com/go/base"
)

func (r *mutationResolver) SimpleEmail(ctx context.Context, subject string, text string, to []string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	status, err := r.interactor.Mail.SimpleEmail(subject, text, to...)
	if err != nil {
		return "", fmt.Errorf("unable to send an email: %v", err)
	}

	defer base.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"simpleEmail",
		err,
	)

	return status, nil
}
