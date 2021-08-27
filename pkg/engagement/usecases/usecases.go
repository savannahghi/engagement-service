package usecases

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	openSource "github.com/savannahghi/engagement/pkg/engagement/usecases"
)

type Usecases interface {
	openSource.Usecases
	// New usecases
	TestFeature(ctx context.Context) (bool, error)
}

type CombinedInteractor struct {
	openSource.Usecases
	InternalInteractor
}

func NewUsecasesInteractor(infrastructure infrastructure.Infrastructure) *CombinedInteractor {

	open := openSource.NewUsecasesInteractor(infrastructure)
	internal := NewInternalInteractor(infrastructure)

	impl := &CombinedInteractor{
		open,
		*internal,
	}

	return impl
}

type InternalInteractor struct {
	infrastructure infrastructure.Infrastructure
}

func NewInternalInteractor(infrastructure infrastructure.Infrastructure) *InternalInteractor {
	return &InternalInteractor{infrastructure: infrastructure}
}

func (i InternalInteractor) TestFeature(ctx context.Context) (bool, error) {
	return true, nil
}
