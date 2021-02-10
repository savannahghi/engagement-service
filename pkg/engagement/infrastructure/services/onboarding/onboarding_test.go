package onboarding_test

import (
	"testing"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"

	"gitlab.slade360emr.com/go/base"
)

func TestNewRemoteProfileService(t *testing.T) {
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
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
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "email inter-service API call to profile",
			args: args{
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetEmailAddresses(tt.args.uids)
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
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "phone inter-service API call to profile",
			args: args{
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetPhoneNumbers(tt.args.uids)
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
	deps, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return
	}

	profileClient, err := base.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	_, token, err := base.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return
	}

	type args struct {
		uids onboarding.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "tokens inter-service API call to profile",
			args: args{
				uids: onboarding.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetDeviceTokens(tt.args.uids)
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
