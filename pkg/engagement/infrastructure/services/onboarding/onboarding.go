package onboarding

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/onboarding")

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
	GetUserProfile(ctx context.Context, uid string) (*base.UserProfile, error)
	IsOptedOut(ctx context.Context, phoneNumber string) (bool, error)
	PhonesWithoutOptOut(ctx context.Context, phones []string) ([]string, error)
}

// NewRemoteProfileService initializes a connection to a remote profile service
// that we will invoke via inter-service communication
func NewRemoteProfileService(
	profileClient *base.InterServiceClient,
) ProfileService {
	return &RemoteProfileService{
		profileClient: profileClient,
	}
}

// RemoteProfileService uses inter-service REST APIs to fetch information
// from a remote profile service
type RemoteProfileService struct {
	profileClient *base.InterServiceClient
}

// NewOnboardingClient initializes a new interservice client for onboarding
func NewOnboardingClient() *base.InterServiceClient {
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
) (*base.UserProfile, error) {
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
	user := base.UserProfile{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error parsing user profile data: %w", err)
	}
	return &user, nil
}

// IsOptedOut checks and returns if a user is opted out or not
func (rps RemoteProfileService) IsOptedOut(
	ctx context.Context,
	phoneNumber string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "IsOptedOut")
	defer span.End()
	payload := map[string]interface{}{
		"phoneNumber": phoneNumber,
	}
	resp, err := rps.profileClient.MakeRequest(
		ctx,
		http.MethodPost,
		isOptedOut,
		payload,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf(
			"unable to make remote profile call with error: %v",
			err,
		)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(
			"failed to get opted out status. Error code: %v",
			resp.StatusCode,
		)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf(
			"error reading profile response body: %w",
			err,
		)
	}

	body := map[string]bool{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("error parsing user profile data: %w", err)
	}

	return body["opted_out"], nil
}

// PhonesWithoutOptOut given a slice of phone numbers, returns numbers that have not opted out
// of our marketing messages programme
func (rps RemoteProfileService) PhonesWithoutOptOut(ctx context.Context, phones []string) ([]string, error) {
	ctx, span := tracer.Start(ctx, "PhonesWithoutOptOut")
	defer span.End()

	var whitelistedNumbers []string
	for _, phone := range phones {
		optedOut, err := rps.IsOptedOut(ctx, phone)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return nil, err
		}
		if !optedOut {
			whitelistedNumbers = append(whitelistedNumbers, phone)
		}
	}

	return whitelistedNumbers, nil
}
