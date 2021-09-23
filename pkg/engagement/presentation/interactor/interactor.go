package interactor

import (
	"github.com/savannahghi/engagement-service/pkg/engagement/usecases"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	sharelib "github.com/savannahghi/engagementcore/pkg/engagement/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra     infrastructure.Interactor
	OpenSourceUsecases  sharelib.Interactor
	UsecaseNotification usecases.NotificationUsecases
}

// NewEngagementInteractor returns a new engagement interactor
func NewEngagementInteractor(
	openSourceInfra infrastructure.Interactor,
	openSourceUsecases sharelib.Interactor,
	notification usecases.NotificationUsecases,

) (*Interactor, error) {
	return &Interactor{
		OpenSourceInfra:     openSourceInfra,
		OpenSourceUsecases:  openSourceUsecases,
		UsecaseNotification: notification,
	}, nil
}
