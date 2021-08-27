package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	generated1 "github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	ghostCMSPost, err := r.usecases.GetLibraryContent(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get library content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getLibraryContent", err)

	return ghostCMSPost, nil
}

func (r *queryResolver) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	faqs, err := r.usecases.GetFaqsContent(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get FAQs content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFaqsContent", err)

	return faqs, nil
}

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
