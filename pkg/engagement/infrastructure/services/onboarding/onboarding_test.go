package onboarding_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
)

func initializeTestService(t *testing.T) (onboarding.ProfileService, context.Context, *auth.Token, error) {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil, nil, nil, err
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil, nil, nil, err
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		profileClient,
	)
	if err != nil {
		t.Errorf("can't get phone number user: %v", err)
		return nil, nil, nil, err
	}
	return rps, ctx, token, nil
}

func TestNewRemoteProfileService(t *testing.T) {
	rps, _, _, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
		return
	}
	if rps == nil {
		t.Errorf("got back nil remote profile service")
		return
	}
}

func TestRemoteProfileService_GetEmailAddresses(t *testing.T) {
	rps, ctx, token, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
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
	rps, ctx, token, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
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
	rps, ctx, token, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
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
	rps, ctx, _, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
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

func TestRemoteProfileService_GetUserProfileByPhoneOrEmail(t *testing.T) {
	rps, ctx, _, err := initializeTestService(t)
	if err != nil {
		t.Errorf("an error occurred %v", err)
		return
	}

	validPhone := interserviceclient.TestUserPhoneNumber
	invalidPhone := "+2547+"
	invalidEmail := "test"
	type args struct {
		ctx     context.Context
		payload *dto.RetrieveUserProfileInput
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "Happy case:phone",
			args: args{
				ctx: ctx,
				payload: &dto.RetrieveUserProfileInput{
					PhoneNumber: &validPhone,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case:phone",
			args: args{
				ctx: context.Background(),
				payload: &dto.RetrieveUserProfileInput{
					PhoneNumber: &invalidPhone,
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Sad case:email",
			args: args{
				ctx: context.Background(),
				payload: &dto.RetrieveUserProfileInput{
					EmailAddress: &invalidEmail,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rps.GetUserProfileByPhoneOrEmail(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
