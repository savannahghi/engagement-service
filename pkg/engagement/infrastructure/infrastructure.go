package infrastructure

import openSource "github.com/savannahghi/engagement/pkg/engagement/infrastructure"

type Infrastructure interface {
	openSource.Infrastructure
	// New usecases
}
type CombinedInteractor struct {
	openSource.Infrastructure
	InternalInteractor
}

func NewInfrastructureInteractor() *CombinedInteractor {

	open := openSource.NewInfrastructureInteractor()
	internal := NewInternalInteractor()

	impl := &CombinedInteractor{open, *internal}

	return impl
}

type InternalInteractor struct {
}

func NewInternalInteractor() *InternalInteractor {
	return &InternalInteractor{}
}

func (i InternalInteractor) DummyInfrastructureFunction() {}
