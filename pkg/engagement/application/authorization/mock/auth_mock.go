package mock

import "github.com/savannahghi/profileutils"

// FakeAuth ...
type FakeAuth struct {
	// CheckPemissionsFn ...
	CheckPemissionsFn func(subject string, input profileutils.PermissionInput) (bool, error)
	// CheckAuthorizationFn ...
	CheckAuthorizationFn func(subject string, permission profileutils.PermissionInput) (bool, error)
	// IsAuthorizedFn ...
	IsAuthorizedFn func(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error)
}

// CheckPemissions is a mock version of the original function
func (a *FakeAuth) CheckPemissions(subject string, input profileutils.PermissionInput) (bool, error) {
	return a.CheckPemissionsFn(subject, input)
}

// CheckAuthorization is a mock version of the original function
func (a *FakeAuth) CheckAuthorization(subject string, permission profileutils.PermissionInput) (bool, error) {
	return a.CheckAuthorizationFn(subject, permission)
}

// IsAuthorized is a mock version of the original function
func (a *FakeAuth) IsAuthorized(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
	return a.IsAuthorizedFn(user, permission)
}
