package twilio_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/database"
	crmExt "github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/edi"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagement-service/pkg/engagement/repository"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func newTwilioService(ctx context.Context) (*twilio.Service, error) {
	var repo repository.Repository
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate notification service: %w",
			err,
		)
	}
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr, hubspotService)
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}
	mail := mail.NewService(fr)
	edi := edi.NewEdiService(edi.NewEDIClient())
	crmExt := crmExt.NewCrmService(hubspotUsecases, mail)
	sms := sms.NewService(repo, crmExt, ns, edi)

	return twilio.NewService(sms, repo), nil
}

func TestNewService(t *testing.T) {
	srv, err := newTwilioService(context.Background())
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}
	assert.NotNil(t, srv)
	if srv == nil {
		t.Errorf("nil twilio service")
		return
	}
}

func setTwilioCredsToLive() (string, string, error) {
	initialTwilioAuthToken := serverutils.MustGetEnvVar("TWILIO_ACCOUNT_AUTH_TOKEN")
	initialTwilioSID := serverutils.MustGetEnvVar("TWILIO_ACCOUNT_SID")

	liveTwilioAuthToken := serverutils.MustGetEnvVar("TESTING_TWILIO_ACCOUNT_AUTH_TOKEN")
	liveTwilioSID := serverutils.MustGetEnvVar("TESTING_TWILIO_ACCOUNT_SID")

	err := os.Setenv("TWILIO_ACCOUNT_AUTH_TOKEN", liveTwilioAuthToken)
	if err != nil {
		return "", "", fmt.Errorf("unable to set twilio auth token to live: %v", err)
	}
	err = os.Setenv("TWILIO_ACCOUNT_SID", liveTwilioSID)
	if err != nil {
		return "", "", fmt.Errorf("unable to set test twilio auth token to live: %v", err)
	}

	return initialTwilioAuthToken, initialTwilioSID, nil
}

func restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID string) error {
	err := os.Setenv("TWILIO_ACCOUNT_AUTH_TOKEN", initialTwilioAuthToken)
	if err != nil {
		return fmt.Errorf("unable to restore twilio auth token: %v", err)
	}
	err = os.Setenv("TWILIO_ACCOUNT_SID", initialTwilioSID)
	if err != nil {
		return fmt.Errorf("unable to restore twilio sid: %v", err)
	}
	return nil
}

func TestService_Room(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid test case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			room, err := s.Room(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Room() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if room == nil {
				t.Errorf("nil room")
				return
			}

			if tt.wantErr == false {
				if room.Type != "peer-to-peer" {
					t.Errorf("room.Type is not peer to peer")
					return
				}
			}
		})
	}

	// Restore envs after test
	err = restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID)
	if err != nil {
		t.Errorf("unable to restore twilio credentials: %v", err)
		return
	}
}

func TestService_AccessToken(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: firebasetools.GetAuthenticatedContext(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			got, err := s.TwilioAccessToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("nil AccessToken value got")
				return
			}
			if got.JWT == "" {
				t.Errorf("empty access token JWT value got")
				return
			}
			if got.UniqueName == "" {
				t.Errorf("empty access token Unique Name value got")
				return
			}
			if got.SID == "" {
				t.Errorf("empty access token SID value got")
				return
			}
			if got.DateUpdated.IsZero() {
				t.Errorf("empty access token Date Updated value got")
				return
			}
			if got.Status == "" {
				t.Errorf("empty access token Status value got")
				return
			}
			if got.Type == "" {
				t.Errorf("empty access token Type value got")
				return
			}
			if got.MaxParticipants == 0 {
				t.Errorf("empty access token Max Participants value got")
				return
			}
		})
	}

	// Restore envs after test
	err = restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID)
	if err != nil {
		t.Errorf("unable to restore twilio credentials: %v", err)
		return
	}

}

func TestService_SendSMS(t *testing.T) {

	// set test credentials
	initialSmsNumber := serverutils.MustGetEnvVar(twilio.TwilioSMSNumberEnvVarName)
	testSmsNumber := serverutils.MustGetEnvVar("TEST_TWILIO_SMS_NUMBER")
	os.Setenv(twilio.TwilioSMSNumberEnvVarName, testSmsNumber)

	type args struct {
		ctx                              context.Context
		normalizedDestinationPhoneNumber string
		msg                              string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:                              context.Background(),
				normalizedDestinationPhoneNumber: testSmsNumber,
				msg:                              "Test message via Twilio",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			if err := s.SendSMS(tt.args.ctx, tt.args.normalizedDestinationPhoneNumber, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// restore twilio sms phone number
	err := os.Setenv(twilio.TwilioSMSNumberEnvVarName, initialSmsNumber)
	if err != nil {
		t.Errorf("unable to restore twilio sms number envar: %v", err)
	}
}
