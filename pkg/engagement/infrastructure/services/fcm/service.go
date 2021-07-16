package fcm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/fcm")

// Service provides methods for sending Firebase Cloud Messaging notifications
type Service struct {
	fcmClient       *messaging.Client
	firestoreClient *firestore.Client
	Repository      repository.Repository
}

func initializeFirestoreClient() (*firestore.Client, error) {
	ctx := context.Background()
	fc := &firebasetools.FirebaseClient{}
	app, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app: %s", err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firestore client: %v", err)
	}

	return firestoreClient, nil
}

func initializeFCMClient() (*messaging.Client, error) {
	ctx := context.Background()
	fc := &firebasetools.FirebaseClient{}
	app, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app: %s", err)
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Messaging client: %v", err)
	}

	return fcmClient, nil
}

// ServiceFCM defines all interactions with the FCM service
type ServiceFCM interface {
	SendNotification(
		ctx context.Context,
		registrationTokens []string,
		data map[string]string,
		notification *firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)

	Notifications(
		ctx context.Context,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)
}

// NewService initializes a service to interact with Firebase Cloud Messaging
func NewService(repository repository.Repository) *Service {
	fcmClient, err := initializeFCMClient()
	if err != nil {
		log.Panicf("error getting Messaging client: %v\n", err)
	}

	firestoreClient, err := initializeFirestoreClient()
	if err != nil {
		log.Panicf("error getting Firestore client: %v\n", err)
	}

	srv := &Service{
		fcmClient:       fcmClient,
		firestoreClient: firestoreClient,
		Repository:      repository,
	}
	srv.checkPreconditions()
	return srv
}

func (s Service) checkPreconditions() {
	if s.fcmClient == nil {
		log.Panicf("nil messaging client in FCM service")
	}
	if s.firestoreClient == nil {
		log.Panicf("nil firestore client in FCM service")
	}
}

// SendNotification sends a data message to the specified registration tokens.
//
// It returns:
//
//  - a list of registration tokens for which message sending failed
//  - an error, if no message sending occured
//
// Notification messages can also be accompanied by custom `data`.
//
// For data messages, the following keys should be avoided:
//
//  - reserved words: "from", "notification" and "message_type"
//  - any word starting with "gcm" or "google"
//
// Messages that are time sensitive (e.g video calls) should be sent with
// `HIGH_PRIORITY`. Their time to live should also be limited (or the expiry)
// set on iOS. For Android, there is a `TTL` key in `messaging.AndroidConfig`.
// For iOS, the `apns-expiration` header should be set to a specific timestamp
// e.g `"apns-expiration":"1604750400"`. For web, there's a `TTL` header that
// is also a number of seconds e.g. `"TTL":"4500"`.
//
// For Android, priority is set via the `messaging.AndroidConfig` `priority`
// key to either "normal" or "high". It should be set to "high" only for urgent
// notification e.g video call notifications. For web, it is set via the
// `Urgency` header e.g "Urgency": "high". For iOS, the "apns-priority" header
// is used, with "5" for normal/low and "10" to mean urgent/high.
//
// The callers of this method should implement retries and exponential backoff,
// if necessary.
func (s Service) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "SendNotification")
	defer span.End()
	s.checkPreconditions()

	if registrationTokens == nil {
		return false, fmt.Errorf("can't send FCM notifications to nil registration tokens")
	}

	message := &messaging.MulticastMessage{Tokens: registrationTokens}

	if data != nil {
		err := ValidateFCMData(data)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return false, err
		}
		message.Data = data
	}
	if notification != nil {
		message.Notification = &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		}
		if notification.ImageURL != nil {
			message.Notification.ImageURL = *notification.ImageURL
		}
	}
	if android != nil {
		message.Android = &messaging.AndroidConfig{
			Priority: android.Priority,
			Data:     base.ConvertInterfaceMap(android.Data),
		}
		if android.CollapseKey != nil {
			message.Android.CollapseKey = *android.CollapseKey
		}
		if android.RestrictedPackageName != nil {
			message.Android.RestrictedPackageName = *android.RestrictedPackageName
		}
	}
	if web != nil {
		message.Webpush = &messaging.WebpushConfig{
			Headers: base.ConvertInterfaceMap(web.Headers),
			Data:    base.ConvertInterfaceMap(web.Data),
		}
	}
	if ios != nil {
		message.APNS = &messaging.APNSConfig{
			Headers: base.ConvertInterfaceMap(web.Headers),
		}
	}

	batchResp, err := s.fcmClient.SendMulticast(ctx, message)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to send FCM messages: %w", err)
	}

	var errorMessages []string
	for idx, resp := range batchResp.Responses {
		if !resp.Success {
			// The order of responses corresponds to the order of the registration tokens.
			msg := fmt.Sprintf(
				"fcm: failed to send message to %s: %v",
				registrationTokens[idx],
				resp.Error,
			)
			errorMessages = append(errorMessages, msg)
		}
		if notification != nil {
			savedNotification := dto.SavedNotification{
				ID:                uuid.New().String(),
				RegistrationToken: registrationTokens[idx],
				MessageID:         resp.MessageID,
				Timestamp:         time.Now(),
			}
			if notification != nil {
				savedNotification.Notification = &dto.FirebaseSimpleNotification{
					Title:    notification.Title,
					Body:     notification.Body,
					ImageURL: notification.ImageURL,
				}
			}
			if data != nil {
				savedNotification.Data = base.ConvertStringMap(data)
			}
			if android != nil {
				savedNotification.AndroidConfig = &dto.FirebaseAndroidConfig{
					CollapseKey: android.CollapseKey,
					Priority:    android.Priority,
					Data:        android.Data,
				}
			}
			if web != nil {
				savedNotification.WebpushConfig = &dto.FirebaseWebpushConfig{
					Headers: web.Headers,
					Data:    web.Data,
				}
			}
			if ios != nil {
				savedNotification.APNSConfig = &dto.FirebaseAPNSConfig{
					Headers: ios.Headers,
				}
			}
			err = s.Repository.SaveNotification(ctx, s.firestoreClient, savedNotification)
			if err != nil {
				helpers.RecordSpanError(span, err)
				log.Printf("unable to save notification: %v", err)
			}
		}
	}
	if len(errorMessages) > 0 {
		return false, fmt.Errorf(strings.Join(errorMessages, "; "))
	}
	return true, nil
}

// Notifications is used to query a user's priorities
func (s Service) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	s.checkPreconditions()
	return s.Repository.RetrieveNotification(ctx, s.firestoreClient, registrationToken, newerThan, limit)
}
