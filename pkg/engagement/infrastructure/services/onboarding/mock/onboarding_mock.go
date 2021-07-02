package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"
)

// FakeServiceOnboarding is an "onboarding" service mock
type FakeServiceOnboarding struct {
	GetEmailAddressesFn func(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error)
	GetPhoneNumbersFn   func(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error)
	GetDeviceTokensFn   func(ctx context.Context, uid onboarding.UserUIDs) (map[string][]string, error)
	GetUserProfileFn    func(ctx context.Context, uid string) (*base.UserProfile, error)

	// GetDeviceTokensFn   func(uids UserUIDsFn) (map[string][]string, error)
}

// GetEmailAddresses ...
func (f *FakeServiceOnboarding) GetEmailAddresses(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetEmailAddressesFn(ctx, uids)
}

// GetPhoneNumbers ...
func (f *FakeServiceOnboarding) GetPhoneNumbers(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetPhoneNumbersFn(ctx, uids)
}

// GetDeviceTokens ...
func (f *FakeServiceOnboarding) GetDeviceTokens(ctx context.Context, uid onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetDeviceTokensFn(ctx, uid)
}

// GetUserProfile ...
func (f *FakeServiceOnboarding) GetUserProfile(ctx context.Context, uid string) (*base.UserProfile, error) {
	return f.GetUserProfileFn(ctx, uid)
}
