package usecases

import (
	libRepository "github.com/savannahghi/engagementcore/pkg/engagement/repository"
	libFeed "github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
)

// FeedUsecases represents all the profile business logic
type FeedUsecases interface {
	libFeed.Usecases
}

// FeedImpl represents the feed usecase implementation
type FeedImpl struct {
	repository libRepository.Repository
}

// NewFeed initializes a user feed
func NewFeed(
	repository libRepository.Repository,
) *FeedImpl {
	return &FeedImpl{
		repository: repository,
	}
}
