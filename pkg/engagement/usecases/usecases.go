package usecases

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	openSource "github.com/savannahghi/engagement/pkg/engagement/usecases"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	openSource.Usecases
	// New usecases
	// Just an example of how you can add new usecase
	TestFeature(ctx context.Context) (bool, error)
}

//CombinedInteractor is used to combine both our internal and open source
// usecases implementation
type CombinedInteractor struct {
	openSource.Usecases
	InternalInteractor
}

// NewUsecasesInteractor initializes new combined usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Infrastructure) *CombinedInteractor {

	open := openSource.NewUsecasesInteractor(infrastructure)
	internal := NewInternalInteractor(infrastructure)

	impl := &CombinedInteractor{
		open,
		*internal,
	}

	return impl
}

// InternalInteractor ...
type InternalInteractor struct {
	infrastructure infrastructure.Infrastructure
}

// NewInternalInteractor ...
func NewInternalInteractor(infrastructure infrastructure.Infrastructure) *InternalInteractor {
	return &InternalInteractor{infrastructure: infrastructure}
}

// TestFeature ...
func (i InternalInteractor) TestFeature(ctx context.Context) (bool, error) {
	return true, nil
}
