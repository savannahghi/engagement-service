package authorization_test

import (
	"log"
	"testing"

	"github.com/savannahghi/engagement/pkg/engagement/application/authorization"
	"github.com/savannahghi/engagement/pkg/engagement/application/authorization/mock"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/stretchr/testify/assert"
)

var fakeAuthorization mock.FakeAuth

func TestCheckPemissions(t *testing.T) {
	type args struct {
		subject string
		input   profileutils.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: permission is set and subject has permission",
			args: args{
				subject: "254711223344",
				input: profileutils.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid: unknown subject with unknown resource",
			args: args{
				subject: "mail@example.com",
				input: profileutils.PermissionInput{
					Resource: "unknown_resource",
					Action:   "edit",
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "sad case: missing subject and input",
			args:    args{},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: permission is set and subject has permission" {
				fakeAuthorization.CheckPemissionsFn = func(subject string, input profileutils.PermissionInput) (bool, error) {
					return true, nil
				}
				got, err := fakeAuthorization.CheckPemissionsFn(tt.args.subject, tt.args.input)
				assert.True(t, got)
				assert.Nil(t, err)
			}
			if tt.name == "valid: unknown subject with unknown resource" {
				fakeAuthorization.CheckPemissionsFn = func(subject string, input profileutils.PermissionInput) (bool, error) {
					return false, nil
				}
				got, err := fakeAuthorization.CheckPemissionsFn(tt.args.subject, tt.args.input)
				assert.False(t, got)
				assert.Nil(t, err)
			}
			if tt.name == "sad case: missing subject and input" {
				fakeAuthorization.CheckPemissionsFn = func(subject string, input profileutils.PermissionInput) (bool, error) {
					return false, nil
				}
				got, err := fakeAuthorization.CheckPemissionsFn(tt.args.subject, tt.args.input)
				assert.False(t, got)
				assert.Nil(t, err)
			}

			got, err := authorization.CheckPemissions(tt.args.subject, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPemissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPemissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAuthorization(t *testing.T) {
	type args struct {
		subject    string
		permission profileutils.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: permission is set and subject has permission",
			args: args{
				subject: "254711223344",
				permission: profileutils.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid: unknown subject with unknown resource",
			args: args{
				subject: "mail@example.com",
				permission: profileutils.PermissionInput{
					Resource: "unknown_resource",
					Action:   "edit",
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "sad case: missing args, subject and input",
			args:    args{},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: permission is set and subject has permission" {
				fakeAuthorization.CheckAuthorizationFn = func(subject string, permission profileutils.PermissionInput) (bool, error) {
					return true, nil
				}
				got, err := fakeAuthorization.CheckAuthorizationFn(tt.args.subject, tt.args.permission)
				assert.True(t, got)
				assert.Nil(t, err)
			}
			if tt.name == "valid: unknown subject with unknown resource" {
				fakeAuthorization.CheckAuthorizationFn = func(subject string, permission profileutils.PermissionInput) (bool, error) {
					return false, nil
				}
				got, err := fakeAuthorization.CheckAuthorizationFn(tt.args.subject, tt.args.permission)
				assert.False(t, got)
				assert.Nil(t, err)
			}
			if tt.name == "sad case: missing args, subject and input" {
				fakeAuthorization.CheckAuthorizationFn = func(subject string, permission profileutils.PermissionInput) (bool, error) {
					return false, nil
				}
				got, err := fakeAuthorization.CheckAuthorizationFn(tt.args.subject, tt.args.permission)
				assert.False(t, got)
				assert.Nil(t, err)
			}

			got, err := authorization.CheckAuthorization(tt.args.subject, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAuthorization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAuthorization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	type args struct {
		user       *profileutils.UserInfo
		permission profileutils.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
		panics  bool
	}{
		{
			name: "valid: permission is set and subject has permission",
			args: args{
				user: &profileutils.UserInfo{
					DisplayName: "test",
					Email:       "test@example.com",
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
				},
				permission: profileutils.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case: missing user",
			args: args{
				permission: profileutils.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			panics: true,
		},
		{
			name:   "sad case: missing args, user and permission",
			args:   args{},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid: permission is set and subject has permission" {
				fakeAuthorization.IsAuthorizedFn = func(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
					return true, nil
				}
				got, err := fakeAuthorization.IsAuthorizedFn(tt.args.user, tt.args.permission)
				if (err != nil) != tt.wantErr {
					t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
				}
			}
			if tt.name == "sad case: missing user" {
				fakeAuthorization.IsAuthorizedFn = func(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
					log.Panic("panic")
					return false, nil
				}
				fcIsAuthorized := func() { _, _ = fakeAuthorization.IsAuthorizedFn(tt.args.user, tt.args.permission) }
				assert.Panics(t, fcIsAuthorized)
			}
			if tt.name == "sad case: missing args, user and permission" {
				fakeAuthorization.IsAuthorizedFn = func(user *profileutils.UserInfo, permission profileutils.PermissionInput) (bool, error) {
					log.Panic("panic")
					return false, nil
				}
				fcIsAuthorized := func() { _, _ = fakeAuthorization.IsAuthorizedFn(tt.args.user, tt.args.permission) }
				assert.Panics(t, fcIsAuthorized)
			}

			if !tt.panics {
				got, err := authorization.IsAuthorized(tt.args.user, tt.args.permission)
				if (err != nil) != tt.wantErr {
					t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
				}
			}
			if tt.panics {
				fcIsAuthorized := func() { _, _ = authorization.IsAuthorized(tt.args.user, tt.args.permission) }
				assert.Panics(t, fcIsAuthorized)
			}
		})
	}
}
