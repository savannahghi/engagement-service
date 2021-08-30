package usecases

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	engagementLib "github.com/savannahghi/engagement/pkg/engagement/usecases"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	engagementLib.Usecases
	// New usecases
	// Just an example of how you can add new usecase
	TestFeature(ctx context.Context) (bool, error)
}

//CombinedInteractor is used to combine both our internal and open source
// usecases implementation
type CombinedInteractor struct {
	engagementLib.Usecases
	InternalInteractor
}

// NewUsecasesInteractor initializes new combined usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Infrastructure) *CombinedInteractor {

	open := engagementLib.NewUsecasesInteractor(infrastructure)
	internal := NewInternalInteractor(infrastructure)

	impl := &CombinedInteractor{
		open,
		*internal,
	}

	return impl
}

// InternalInteractor ...
type InternalInteractor struct {
	engagementLib  engagementLib.Usecases
	infrastructure infrastructure.Infrastructure
}

// NewInternalInteractor ...
func NewInternalInteractor(infrastructure infrastructure.Infrastructure) *InternalInteractor {
	library := engagementLib.NewUsecasesInteractor(infrastructure)

	return &InternalInteractor{
		infrastructure: infrastructure,
		engagementLib:  library,
	}
}

// TestFeature ...
func (i InternalInteractor) TestFeature(ctx context.Context) (bool, error) {
	return true, nil
}
