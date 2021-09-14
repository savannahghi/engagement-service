package authorization

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
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
