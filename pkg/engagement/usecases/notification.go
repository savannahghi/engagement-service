package usecases

import (
	libRepository "github.com/savannahghi/engagementcore/pkg/engagement/repository"
	libNotification "github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
)

// NotificationUsecases represent logic required to make notification
type NotificationUsecases interface {
	libNotification.NotificationUsecases
}

// NotificationImpl represents the notification usecase implementation
type NotificationImpl struct {
	LibRepository libRepository.Repository
	LibUsecases   libNotification.NotificationUsecases
}

// NewNotification initializes a notification usecase
func NewNotification(
	libRepository libRepository.Repository,
	libUsecases libNotification.NotificationUsecases,
) *NotificationImpl {
	return &NotificationImpl{
		LibRepository: libRepository,
		LibUsecases:   libUsecases,
	}
}
