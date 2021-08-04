// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/surveys"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/uploads"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	"github.com/savannahghi/engagement/pkg/engagement/usecases"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	Feed         usecases.FeedUseCases
	Notification usecases.NotificationUsecases
	Uploads      uploads.ServiceUploads
	Library      library.ServiceLibrary
	SMS          sms.ServiceSMS
	Mail         mail.ServiceMail
	Whatsapp     whatsapp.ServiceWhatsapp
	OTP          otp.ServiceOTP
	Twilio       twilio.ServiceTwilio
	FCM          fcm.ServiceFCM
	Surveys      surveys.ServiceSurveys
	CRM          hubspot.ServiceHubSpotInterface
	CrmExt       crm.ServiceCrm
}

// NewEngagementInteractor returns a new engagement interactor
func NewEngagementInteractor(
	feed usecases.FeedUseCases,
	notification usecases.NotificationUsecases,
	uploads uploads.ServiceUploads,
	library library.ServiceLibrary,
	sms sms.ServiceSMS,
	mail mail.ServiceMail,
	whatsapp whatsapp.ServiceWhatsapp,
	otp otp.ServiceOTP,
	twilio twilio.ServiceTwilio,
	fcm fcm.ServiceFCM,
	surveys surveys.ServiceSurveys,
	CRM hubspot.ServiceHubSpotInterface,
	crmExt crm.ServiceCrm,
) (*Interactor, error) {
	return &Interactor{
		Feed:         feed,
		Notification: notification,
		Uploads:      uploads,
		Library:      library,
		SMS:          sms,
		Mail:         mail,
		Whatsapp:     whatsapp,
		OTP:          otp,
		Twilio:       twilio,
		FCM:          fcm,
		Surveys:      surveys,
		CRM:          CRM,
		CrmExt:       crmExt,
	}, nil
}
