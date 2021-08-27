package infrastructure

import openSource "github.com/savannahghi/engagement/pkg/engagement/infrastructure"

// Infrastructure is an interface that combines of all infrastructure
type Infrastructure interface {
	openSource.Infrastructure
}

//CombinedInteractor is used to combine both our internal and open source
// infrastructure implementation
type CombinedInteractor struct {
	openSource.Infrastructure
	InternalInteractor
}

// NewInfrastructureInteractor initializes new combined interactor
func NewInfrastructureInteractor() *CombinedInteractor {
	// open is open source interactor that both internal and open source services depend on
	open := openSource.NewInfrastructureInteractor()
	// internal is internal interactor that our internal services depends on
	internal := NewInternalInteractor()

	impl := &CombinedInteractor{open, *internal}

	return impl
}

// InternalInteractor ...
type InternalInteractor struct {
}

// NewInternalInteractor initializes new internal interactor
func NewInternalInteractor() *InternalInteractor {
	return &InternalInteractor{}
}

// DummyInfrastructureFunction ...
func (i InternalInteractor) DummyInfrastructureFunction() {}
