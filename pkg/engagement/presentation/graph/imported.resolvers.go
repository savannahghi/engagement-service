package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) SendNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notificationData, err := converterandformatter.MapInterfaceToMapString(data)
	if err != nil {
		return false, err
	}

	sent, err := r.usecases.SendNotification(
		ctx,
		registrationTokens,
		notificationData,
		&notification,
		android,
		ios,
		web,
	)
	if err != nil {
		return false, fmt.Errorf("failed to send a notification : %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sendNotification",
		err,
	)

	return sent, nil
}

func (r *mutationResolver) SendFCMByPhoneOrEmail(ctx context.Context, phoneNumber *string, email *string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	sent, err := r.usecases.SendFCMByPhoneOrEmail(
		ctx,
		phoneNumber,
		email,
		data,
		notification,
		android,
		ios,
		web,
	)
	if err != nil {
		return false, fmt.Errorf("failed to send an FCM notification by email or phone : %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"SendFCMByPhoneOrEmail",
		err,
	)

	return sent, nil
}

func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.ResolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve a Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "resolveFeedItem", err)

	return item, nil
}

func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.UnresolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unresolveFeedItem", err)

	return item, nil
}

func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.PinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to pin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "pinFeedItem", err)

	return item, nil
}

func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.UnpinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unpin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unpinFeedItem", err)

	return item, nil
}

func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.HideFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideFeedItem", err)

	return item, nil
}

func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.ShowFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to show Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showFeedItem", err)

	return item, nil
}

func (r *mutationResolver) HideNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.usecases.HideNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideNudge", err)

	return nudge, nil
}

func (r *mutationResolver) ShowNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.usecases.ShowNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showNudge", err)

	return nudge, nil
}

func (r *mutationResolver) PostMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, message feedlib.Message) (*feedlib.Message, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	msg, err := r.usecases.PostMessage(ctx, uid, flavour, itemID, &message)
	if err != nil {
		return nil, fmt.Errorf("unable to post a message: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "postMessage", err)

	return msg, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, messageID string) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.usecases.DeleteMessage(ctx, uid, flavour, itemID, messageID)
	if err != nil {
		return false, fmt.Errorf("can't delete message: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "deleteMessage", err)

	return true, nil
}

func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour feedlib.Flavour, event feedlib.Event) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.usecases.ProcessEvent(ctx, uid, flavour, &event)
	if err != nil {
		return false, fmt.Errorf("can't process event: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "processEvent", err)

	return true, nil
}

func (r *mutationResolver) SimpleEmail(ctx context.Context, subject string, text string, to []string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	status, err := r.usecases.SimpleEmail(ctx, subject, text, nil, to)
	if err != nil {
		return "", fmt.Errorf("unable to send an email: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"simpleEmail",
		err,
	)

	return status, nil
}

func (r *mutationResolver) VerifyOtp(ctx context.Context, msisdn string, otp string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verifyOTP, err := r.usecases.VerifyOtp(ctx, msisdn, otp)
	if err != nil {
		return false, fmt.Errorf("failed to check for the validity of the supplied OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"verifyOTP",
		err,
	)

	return verifyOTP, nil
}

func (r *mutationResolver) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verifyEmailOTP, err := r.usecases.VerifyEmailOtp(ctx, email, otp)
	if err != nil {
		return false, fmt.Errorf("failed to check for the validity of the supplied OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"verifyEmailOTP",
		err,
	)

	return verifyEmailOTP, nil
}

func (r *mutationResolver) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	smsResponse, err := r.usecases.Send(ctx, to, message)
	if err != nil {
		return nil, fmt.Errorf("unable send SMS: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"send",
		err,
	)

	return smsResponse, nil
}

func (r *mutationResolver) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	smsResponse, err := r.usecases.SendToMany(ctx, message, to)
	if err != nil {
		return nil, fmt.Errorf("unable to send SMS to many: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"sendToMany",
		err,
	)

	return smsResponse, nil
}

func (r *mutationResolver) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	response, err := r.usecases.RecordNPSResponse(ctx, input)
	if err != nil {
		return false, fmt.Errorf("failed to record nps response")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"recordNPSResponse",
		err,
	)

	return response, nil
}

func (r *mutationResolver) Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error) {
	startTime := time.Now()

	r.checkPreconditions()
	upload, err := r.usecases.Upload(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to upload: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"upload",
		err,
	)

	return upload, nil
}

func (r *mutationResolver) PhoneNumberVerificationCode(ctx context.Context, to string, code string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verificationCode, err := r.usecases.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
	if err != nil {
		return false, fmt.Errorf("failed to send a verification code: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"phoneNumberVerificationCode",
		err,
	)

	return verificationCode, nil
}

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	ghostCMSPost, err := r.usecases.GetLibraryContent(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get library content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getLibraryContent", err)

	return ghostCMSPost, nil
}

func (r *queryResolver) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	faqs, err := r.usecases.GetFaqsContent(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get FAQs content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFaqsContent", err)

	return faqs, nil
}

func (r *queryResolver) Notifications(ctx context.Context, registrationToken string, newerThan time.Time, limit int) ([]*dto.SavedNotification, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	notification, err := r.usecases.Notifications(ctx, registrationToken, newerThan, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"notifications",
		err,
	)

	return notification, nil
}

func (r *queryResolver) GetFeed(ctx context.Context, flavour feedlib.Flavour, isAnonymous bool, persistent feedlib.BooleanFilter, status *feedlib.Status, visibility *feedlib.Visibility, expired *feedlib.BooleanFilter, filterParams *helpers.FilterParams) (*domain.Feed, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	feed, err := r.usecases.GetFeed(
		ctx,
		&uid,
		&isAnonymous,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get Feed: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFeed", err)

	return feed, nil
}

func (r *queryResolver) Labels(ctx context.Context, flavour feedlib.Flavour) ([]string, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	labels, err := r.usecases.Labels(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get Labels count: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "labels", err)

	return labels, nil
}

func (r *queryResolver) UnreadPersistentItems(ctx context.Context, flavour feedlib.Flavour) (int, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return -1, fmt.Errorf("can't get logged in user UID")
	}
	count, err := r.usecases.UnreadPersistentItems(ctx, uid, flavour)
	if err != nil {
		return -1, fmt.Errorf("unable to count unread persistent items: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unreadPersistentItems", err)

	return count, nil
}

func (r *queryResolver) GenerateOtp(ctx context.Context, msisdn string, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.usecases.GenerateAndSendOTP(ctx, msisdn, appID)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) GenerateAndEmailOtp(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.usecases.SendOTPToEmail(ctx, msisdn, email, appID)
	if err != nil {
		return "", fmt.Errorf("failed to send the generated OTP to the provided email address")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateAndEmailOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) GenerateRetryOtp(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.usecases.GenerateRetryOTP(ctx, msisdn, retryStep, appID)
	if err != nil {
		return "", fmt.Errorf("failed to generate fallback OTPs")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"generateRetryOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	otp, err := r.usecases.EmailVerificationOtp(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate an OTP to the supplied email for verification")
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"emailVerificationOTP",
		err,
	)

	return otp, nil
}

func (r *queryResolver) ListNPSResponse(ctx context.Context) ([]*dto.NPSResponse, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	return nil, nil
}

func (r *queryResolver) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	accessToken, err := r.usecases.TwilioAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to generate access token: %w", err)
	}
	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"twilioAccessToken",
		err,
	)

	return accessToken, nil
}

func (r *queryResolver) FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error) {
	startTime := time.Now()

	r.checkPreconditions()
	upload, err := r.usecases.FindUploadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to find upload by ID: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"findUploadByID",
		err,
	)

	return upload, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
