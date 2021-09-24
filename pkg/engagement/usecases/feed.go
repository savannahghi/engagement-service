package usecases

import (
	libInfra "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	libFeed "github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
)

// FeedUsecases represent logic required to make Feed
type FeedUsecases interface {
	libFeed.Usecases
}

// FeedImpl represents the Feed usecase implementation
type FeedImpl struct {
	LibInfrastructure libInfra.Interactor
	LibUsecases       libFeed.Usecases
}

// NewFeed initializes a Feed usecase
func NewFeed(
	libInfra libInfra.Interactor,
	libUsecases libFeed.Usecases,
) *FeedImpl {
	return &FeedImpl{
		LibInfrastructure: libInfra,
		LibUsecases:       libUsecases,
	}
}
