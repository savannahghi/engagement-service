// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/mail"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/sms"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/uploads"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	Feed         usecases.FeedUseCases
	Notification usecases.NotificationUsecases
	Uploads      uploads.ServiceUploads
	Library      library.ServiceLibrary
	SMS          sms.ServiceSMS
	Mail         mail.ServiceMail
}

// NewEngagementInteractor returns a new engagement interactor
func NewEngagementInteractor(
	feed usecases.FeedUseCases,
	notification usecases.NotificationUsecases,
	uploads uploads.ServiceUploads,
	library library.ServiceLibrary,
	sms sms.ServiceSMS,
	mail mail.ServiceMail,
) (*Interactor, error) {
	return &Interactor{
		Feed:         feed,
		Notification: notification,
		Uploads:      uploads,
		Library:      library,
		SMS:          sms,
		Mail:         mail,
	}, nil
}
