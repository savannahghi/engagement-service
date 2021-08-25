package library_test

import (
	"context"
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
)

const (
	onboardingService = "profile"
)

func TestNewLibraryService(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	srv := library.NewLibraryService(onboarding)
	if srv == nil {
		t.Errorf("nil library service")
	}
}

func TestService_GetFeedContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	ctx := firebasetools.GetAuthenticatedContext(t)

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				flavour: feedlib.FlavourConsumer,
				ctx:     ctx,
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFeedContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFeedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}

func TestService_GetFaqsContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:retrieved_user_rbac_faq",
			args: args{
				ctx:     ctx,
				flavour: "PRO",
			},
			wantErr: false,
		},
		{
			name: "valid:retrieved_consumer_faq",
			args: args{
				ctx:     ctx,
				flavour: "CONSUMER",
			},
			wantErr: false,
		},
		{
			name: "invalid:pass_invalid_flavor",
			args: args{
				ctx:     context.Background(),
				flavour: "INVALID",
			},
			wantErr: true,
		},
		{
			name: "invalid:failed_to_get_logged_in_user",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
		{
			name: "invalid:failed_to_get_user_profile",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
		{
			name: "ivalid:failed_to_get_user_faqs",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetFaqsContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFaqsContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_GetLibraryContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	ctx := firebasetools.GetAuthenticatedContext(t)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx,
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetLibraryContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetLibraryContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}
