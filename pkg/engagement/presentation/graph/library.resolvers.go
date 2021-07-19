package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/presentation/graph/generated"
)

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*library.GhostCMSPost, error) {
	startTime := time.Now()

	ghostCMSPost, err := r.interactor.Library.GetLibraryContent(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get library content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getLibraryContent", err)

	return ghostCMSPost, nil
}

func (r *queryResolver) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*library.GhostCMSPost, error) {
	startTime := time.Now()

	faqs, err := r.interactor.Library.GetFaqsContent(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get FAQs content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFaqsContent", err)

	return faqs, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
