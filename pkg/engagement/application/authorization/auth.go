package authorization

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
)

var (
	enforcer *casbin.Enforcer
)

// this function helps to initialize the global variable `enforcer` that cannot be initialized in the global context.

func init() {
	initEnforcer()
}

func initEnforcer() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	conf := filepath.Join(basepath, "/rbac_model.conf")
	dataFile := filepath.Join(basepath, "/data/rbac_policy.csv")
	e, err := casbin.NewEnforcer(conf, dataFile)
	if err != nil {
		log.Panicf("unable to initialize and enforce permissions %v", err)
	}
	enforcer = e
}

// CheckPemissions is used to check whether the permissions of a subject are set
func CheckPemissions(subject string, input profileutils.PermissionInput) (bool, error) {

	ok, err := enforcer.Enforce(subject, input.Resource, input.Action)
	if err != nil {
		return false, fmt.Errorf("unable to check permissions %w", err)
	}
	if ok {
		return true, nil
	}
	return false, nil
}

// CheckAuthorization is used to check the user permissions
func CheckAuthorization(subject string, permission profileutils.PermissionInput) (bool, error) {
	isAuthorized, err := CheckPemissions(subject, permission)
	if err != nil {
		return false, fmt.Errorf("internal server error: can't authorize user: %w", err)
	}

	if !isAuthorized {
		return false, nil
	}

	return true, nil
}

// IsAuthorized checks if the subject identified by their email has permission to access the
// specified resource
// currently only known internal anonymous users and external API Integrations emails are checked, internal and default logged in users
// have access by default.
// for subjects identified by their phone number normalize the phone and omit the first (+) character
func IsAuthorized(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
	if user.PhoneNumber != "" && converterandformatter.StringSliceContains(profileutils.AuthorizedPhones, user.PhoneNumber) {
		return CheckAuthorization(user.PhoneNumber[1:], permission)
	}
	if user.Email != "" && converterandformatter.StringSliceContains(profileutils.AuthorizedEmails, user.Email) {
		return CheckAuthorization(user.Email, permission)

	}
	return true, nil
}

// GetInterserviceBearerTokenHeader gets the Authorization header
func GetInterserviceBearerTokenHeader(ctx context.Context) (string, error) {
	service := interserviceclient.ISCService{} // name and domain not necessary for our use case
	isc, err := interserviceclient.NewInterserviceClient(service)
	if err != nil {
		return "", fmt.Errorf("can't initialize interservice client: %w", err)
	}

	authToken, err := isc.CreateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("can't get auth token: %w", err)
	}
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader, nil
}
