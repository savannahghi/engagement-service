package onboarding

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement/pkg/engagement/services/onboarding")

// specific endpoint paths for ISC
const (
	profileEmails       = "internal/contactdetails/emails/"
	profilePhoneNumbers = "internal/contactdetails/phonenumbers/"
	profileTokens       = "internal/contactdetails/tokens/"
	userProfile         = "internal/user_profile"
	isOptedOut          = "internal/is_opted_out"

	onboardingService = "profile"
)

// UserUIDs is used to serialize user UIDs for inter-service calls to the
// profile service
type UserUIDs struct {
	UIDs []string `json:"uids"`
}

// ProfileService defines the interactions with the profile service
type ProfileService interface {
	GetEmailAddresses(
		ctx context.Context,
		uids UserUIDs,
	) (map[string][]string, error)
	GetPhoneNumbers(
		ctx context.Context,
		uids UserUIDs,
	) (map[string][]string, error)
	GetDeviceTokens(
		ctx context.Context,
		uid UserUIDs,
	) (map[string][]string, error)
	GetUserProfile(ctx context.Context, uid string) (*profileutils.UserProfile, error)
}

// NewRemoteProfileService initializes a connection to a remote profile service
// that we will invoke via inter-service communication
func NewRemoteProfileService(
	profileClient *interserviceclient.InterServiceClient,
) ProfileService {
	return &RemoteProfileService{
		profileClient: profileClient,
	}
}

// RemoteProfileService uses inter-service REST APIs to fetch information
// from a remote profile service
type RemoteProfileService struct {
	profileClient *interserviceclient.InterServiceClient
}

// NewOnboardingClient initializes a new interservice client for onboarding
func NewOnboardingClient() *interserviceclient.InterServiceClient {
	return helpers.InitializeInterServiceClient(onboardingService)
}

func (rps RemoteProfileService) callProfileService(
	ctx context.Context,
	uids UserUIDs, path string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "callProfileService")
	defer span.End()
	resp, err := rps.profileClient.MakeRequest(
		ctx,
		http.MethodPost,
		path,
		uids,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error calling profile service: %w", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error reading profile response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error status code after calling profile service, got status %d and data `%s`",
			resp.StatusCode,
			string(data),
		)
	}

	var contacts map[string][]string
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf(
			"can't unmarshal profile response data \n(\n%s\n)\n: %w",
			string(data),
			err,
		)
	}

	return contacts, nil
}

// GetEmailAddresses gets the specified users' email addresses from the
// staging / testing / prod profile service
func (rps RemoteProfileService) GetEmailAddresses(
	ctx context.Context,
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(ctx, uids, profileEmails)
}

// GetPhoneNumbers gets the specified users' phone numbers from the
// staging / testing / prod profile service
func (rps RemoteProfileService) GetPhoneNumbers(
	ctx context.Context,
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(ctx, uids, profilePhoneNumbers)
}

// GetDeviceTokens gets the specified users' FCM push tokens from the
// staging / testing / prod profile service
func (rps RemoteProfileService) GetDeviceTokens(
	ctx context.Context,
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(ctx, uids, profileTokens)
}

// GetUserProfile gets the specified users' profile from the onboarding service
func (rps RemoteProfileService) GetUserProfile(
	ctx context.Context,
	uid string,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfile")
	defer span.End()
	uidPayoload := dto.UIDPayload{
		UID: &uid,
	}
	resp, err := rps.profileClient.MakeRequest(
		ctx,
		http.MethodPost,
		userProfile,
		uidPayoload,
	)

	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error calling profile service: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"user profile not found. Error code: %v",
			resp.StatusCode,
		)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error reading profile response body: %w", err)
	}
	user := profileutils.UserProfile{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error parsing user profile data: %w", err)
	}
	return &user, nil
}
