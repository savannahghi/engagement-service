package mock

import "gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"

// FakeServiceOnboarding is an "onboarding" service mock
type FakeServiceOnboarding struct {
	GetEmailAddressesFn func(uids onboarding.UserUIDs) (map[string][]string, error)
	GetPhoneNumbersFn   func(uids onboarding.UserUIDs) (map[string][]string, error)
	GetDeviceTokensFn   func(uid onboarding.UserUIDs) (map[string][]string, error)

	// GetDeviceTokensFn   func(uids UserUIDsFn) (map[string][]string, error)
}

// GetEmailAddresses ...
func (f *FakeServiceOnboarding) GetEmailAddresses(uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetEmailAddressesFn(uids)
}

// GetPhoneNumbers ...
func (f *FakeServiceOnboarding) GetPhoneNumbers(uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetPhoneNumbersFn(uids)
}

// GetDeviceTokens ...
func (f *FakeServiceOnboarding) GetDeviceTokens(uid onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetDeviceTokensFn(uid)
}
