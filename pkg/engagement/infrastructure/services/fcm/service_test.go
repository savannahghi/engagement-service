package fcm_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/fcm"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Exit(m.Run())
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}
	sampleService := fcm.NewService(fr)
	assert.NotNil(t, sampleService)

	tests := []struct {
		name string
		want *fcm.Service
	}{
		{
			name: "good case",
			want: sampleService,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fcm.NewService(fr)
			assert.NotNil(t, got)
		})
	}
}

func TestService_Notifications(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}
	s := fcm.NewService(fr)
	if s == nil {
		t.Errorf("nil FCM service")
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
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}
	fakeToken := uuid.New().String()
	imgURL := "https://www.wxpr.org/sites/wxpr/files/styles/medium/public/202007/chipmunk-5401165_1920.jpg"
	s := fcm.NewService(fr)
	if s == nil {
		t.Errorf("nil FCM service")
		return
	}
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
