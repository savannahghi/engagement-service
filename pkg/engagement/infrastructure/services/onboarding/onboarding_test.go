package onboarding_test

import (
	"context"
	"testing"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"
)

func TestNewRemoteProfileService(t *testing.T) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}

	if profileClient == nil {
		t.Errorf("nil profile client")
		return
	}
	if profileClient.RequestRootDomain == "" {
		t.Errorf("blank request root domain")
		return
	}

	rps := onboarding.NewRemoteProfileService(profileClient)
	if rps == nil {
		t.Errorf("got back nil remote profile service")
		return
	}
}

func TestRemoteProfileService_GetEmailAddresses(t *testing.T) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get email addresses inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get email addresses inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{},
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetEmailAddresses(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetEmailAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetPhoneNumbers(t *testing.T) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get phone numbers inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get phone numbers inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{}, // empty UID list
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetPhoneNumbers(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetDeviceTokens(t *testing.T) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get device tokens inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get device tokens inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboarding.UserUIDs{
					UIDs: []string{},
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetDeviceTokens(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetDeviceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetUserProfile(t *testing.T) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	ctx, _, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}
	type args struct {
		ctx context.Context
		uid string
	}

	UID, err := firebasetools.GetLoggedInUserUID(ctx)
	if err != nil {
		t.Errorf("can't get logged in user: %v", err)
		return
	}
	invalidUID := "9VwnREOH8GdSfaxH69J6MvCu1gp9"

	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name:    "happy case: got user profile",
			args:    args{ctx: ctx, uid: UID},
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "sad case: unable to get user profile",
			args:    args{ctx: ctx, uid: invalidUID},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetUserProfile(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil profile data")
					return
				}
			}
		})
	}
}
