package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/feed/graph/inbox"
)

func (r *queryResolver) GetUserMessages(ctx context.Context) ([]*inbox.Message, error) {
	return r.inboxService.GetUserMessages(ctx)
}

func (r *queryResolver) SendWelcomeMessage(ctx context.Context) (*bool, error) {
	return r.inboxService.SendWelcomeMessageToUser(ctx)
}
