package onboarding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// specific endpoint paths for ISC
const (
	profileEmails       = "internal/contactdetails/emails/"
	profilePhoneNumbers = "internal/contactdetails/phonenumbers/"
	profileTokens       = "internal/contactdetails/tokens/"
)

// UserUIDs is used to serialize user UIDs for inter-service calls to the
// profile service
type UserUIDs struct {
	UIDs []string `json:"uids"`
}

// ProfileService defines the interactions with the profile service
type ProfileService interface {
	GetEmailAddresses(uids UserUIDs) (map[string][]string, error)
	GetPhoneNumbers(uids UserUIDs) (map[string][]string, error)
	GetDeviceTokens(uid UserUIDs) (map[string][]string, error)
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

func (rps RemoteProfileService) callProfileService(
	uids UserUIDs, path string,
) (map[string][]string, error) {
	resp, err := rps.profileClient.MakeRequest(http.MethodPost, path, uids)
	if err != nil {
		return nil, fmt.Errorf("error calling profile service: %w", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading profile response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"error status code after calling profile service, got status %d and data `%s`",
			resp.StatusCode, string(data),
		)
	}

	var contacts map[string][]string
	err = json.Unmarshal(data, &contacts)
	if err != nil {
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
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(uids, profileEmails)
}

// GetPhoneNumbers gets the specified users' phone numbers from the
// staging / testing / prod profile service
func (rps RemoteProfileService) GetPhoneNumbers(
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(uids, profilePhoneNumbers)
}

// GetDeviceTokens gets the specified users' FCM push tokens from the
// staging / testing / prod profile service
func (rps RemoteProfileService) GetDeviceTokens(
	uids UserUIDs,
) (map[string][]string, error) {
	return rps.callProfileService(uids, profileTokens)
}
