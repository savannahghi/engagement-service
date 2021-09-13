package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
)

func (r *mutationResolver) SendNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendFCMByPhoneOrEmail(ctx context.Context, phoneNumber *string, email *string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) HideNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ShowNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PostMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, message feedlib.Message) (*feedlib.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, messageID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour feedlib.Flavour, event feedlib.Event) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SimpleEmail(ctx context.Context, subject string, text string, to []string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyOtp(ctx context.Context, msisdn string, otp string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PhoneNumberVerificationCode(ctx context.Context, to string, code string, marketingMessage string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notifications(ctx context.Context, registrationToken string, newerThan time.Time, limit int) ([]*dto.SavedNotification, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetFeed(ctx context.Context, flavour feedlib.Flavour, isAnonymous bool, persistent feedlib.BooleanFilter, status *feedlib.Status, visibility *feedlib.Visibility, expired *feedlib.BooleanFilter, filterParams *helpers.FilterParams) (*domain.Feed, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Labels(ctx context.Context, flavour feedlib.Flavour) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) UnreadPersistentItems(ctx context.Context, flavour feedlib.Flavour) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GenerateOtp(ctx context.Context, msisdn string, appID *string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GenerateAndEmailOtp(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GenerateRetryOtp(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListNPSResponse(ctx context.Context) ([]*dto.NPSResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
