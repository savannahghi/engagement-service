package fcm_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Exit(m.Run())
}

func initializeTestService(ctx context.Context, t *testing.T) (*fcm.Service, error) {
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return nil, err
	}

	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil, err
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil, err
	}
	rps := onboarding.NewRemoteProfileService(profileClient)

	s := fcm.NewService(fr, rps)
	if s == nil {
		t.Errorf("nil FCM service")
		return nil, err
	}
	return s, nil
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	s, err := initializeTestService(ctx, t)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}
	assert.NotNil(t, s)

	tests := []struct {
		name string
		want *fcm.Service
	}{
		{
			name: "good case",
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s
			assert.NotNil(t, got)
		})
	}
}

func TestService_Notifications(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	s, err := initializeTestService(ctx, t)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}
	type args struct {
		ctx               context.Context
		registrationToken string
		newerThan         time.Time
		limit             int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no saved notifications, no error",
			args: args{
				ctx:               ctx,
				registrationToken: uuid.New().String(),
				newerThan:         time.Now(),
				limit:             100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Notifications(tt.args.ctx, tt.args.registrationToken, tt.args.newerThan, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Notifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestService_SendNotification(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	s, err := initializeTestService(ctx, t)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	fakeToken := uuid.New().String()
	imgURL := "https://www.wxpr.org/sites/wxpr/files/styles/medium/public/202007/chipmunk-5401165_1920.jpg"

	type args struct {
		ctx                context.Context
		registrationTokens []string
		data               map[string]string
		notification       *firebasetools.FirebaseSimpleNotificationInput
		android            *firebasetools.FirebaseAndroidConfigInput
		ios                *firebasetools.FirebaseAPNSConfigInput
		web                *firebasetools.FirebaseWebpushConfigInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "non existent token - should fail gracefully",
			args: args{
				ctx:                ctx,
				registrationTokens: []string{fakeToken},
				data: map[string]string{
					"some": "data",
				},
				notification: &firebasetools.FirebaseSimpleNotificationInput{
					Title:    "Test Notification",
					Body:     "From Integration Tests",
					ImageURL: &imgURL,
					Data: map[string]interface{}{
						"more": "data",
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.SendNotification(tt.args.ctx, tt.args.registrationTokens, tt.args.data, tt.args.notification, tt.args.android, tt.args.ios, tt.args.web)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.SendNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendFCMByPhoneOrEmail(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	s, err := initializeTestService(ctx, t)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	phone := interserviceclient.TestUserPhoneNumber
	validEmail := "test@bewell.co.ke"

	type args struct {
		ctx          context.Context
		phoneNumber  *string
		email        *string
		data         map[string]interface{}
		notification firebasetools.FirebaseSimpleNotificationInput
		android      *firebasetools.FirebaseAndroidConfigInput
		ios          *firebasetools.FirebaseAPNSConfigInput
		web          *firebasetools.FirebaseWebpushConfigInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "non existent token using phone - should fail gracefully",
			args: args{
				ctx:         ctx,
				phoneNumber: &phone,
				email:       nil,
				data: map[string]interface{}{
					"more": "data",
				},
				notification: firebasetools.FirebaseSimpleNotificationInput{Title: "Test Notification", Body: "From Integration Tests", Data: map[string]interface{}{"more": "data"}},
				android:      &firebasetools.FirebaseAndroidConfigInput{},
				ios:          &firebasetools.FirebaseAPNSConfigInput{},
				web:          &firebasetools.FirebaseWebpushConfigInput{},
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "non existent token using email - should fail gracefully",
			args: args{
				ctx:         ctx,
				phoneNumber: nil,
				email:       &validEmail,
				data: map[string]interface{}{
					"more": "data",
				},
				notification: firebasetools.FirebaseSimpleNotificationInput{Title: "Test Notification", Body: "From Integration Tests", Data: map[string]interface{}{"more": "data"}},
				android:      &firebasetools.FirebaseAndroidConfigInput{},
				ios:          &firebasetools.FirebaseAPNSConfigInput{},
				web:          &firebasetools.FirebaseWebpushConfigInput{},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.SendFCMByPhoneOrEmail(tt.args.ctx, tt.args.phoneNumber, tt.args.email, tt.args.data, tt.args.notification, tt.args.android, tt.args.ios, tt.args.web)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendFCMByPhoneOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.SendFCMByPhoneOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
