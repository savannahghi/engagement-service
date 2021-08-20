package whatsapp_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	"github.com/stretchr/testify/assert"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	m.Run()
}

func newService() *whatsapp.Service {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		log.Panicf(
			"can't instantiate firebase repository in resolver: %v",
			err,
		)
	}
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		log.Panicf("failed to initialize hubspot crm repository: %v", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr, hubspotService)

	mail := mail.NewService(fr)
	crmExt := crm.NewCrmService(hubspotUsecases, mail)
	return whatsapp.NewService(fr, crmExt)
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "valid initialization",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newService()
			assert.NotNil(t, got)

			// should not panic
			got.CheckPreconditions()
		})
	}
}

func TestService_checkPreconditions(t *testing.T) {
	type fields struct {
		baseURL          string
		accountSID       string
		accountAuthToken string
		sender           string
		httpClient       *http.Client
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name:      "http client not set",
			wantPanic: true,
		},
		{
			name: "base URL not set",
			fields: fields{
				httpClient: http.DefaultClient,
			},
			wantPanic: true,
		},
		{
			name: "account SID not set",
			fields: fields{
				httpClient: http.DefaultClient,
				baseURL:    "http://example.com",
			},
			wantPanic: true,
		},
		{
			name: "account auth token not set",
			fields: fields{
				httpClient: http.DefaultClient,
				baseURL:    "http://example.com",
				accountSID: "something",
			},
			wantPanic: true,
		},
		{
			name: "sender not set",
			fields: fields{
				httpClient:       http.DefaultClient,
				baseURL:          "http://example.com",
				accountSID:       "something",
				accountAuthToken: "something else",
			},
			wantPanic: true,
		},
		{
			name: "all set",
			fields: fields{
				httpClient:       http.DefaultClient,
				baseURL:          "http://example.com",
				accountSID:       "something",
				accountAuthToken: "something else",
				sender:           "+1765432345",
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := whatsapp.Service{
				BaseURL:          tt.fields.baseURL,
				AccountSID:       tt.fields.accountSID,
				AccountAuthToken: tt.fields.accountAuthToken,
				Sender:           tt.fields.sender,
				HTTPClient:       tt.fields.httpClient,
			}

			if tt.wantPanic {
				assert.Panics(t, func() {
					s.CheckPreconditions()
				})
			} else {
				s.CheckPreconditions()
			}
		})
	}
}

func TestService_PhoneNumberVerificationCode(t *testing.T) {
	type args struct {
		ctx              context.Context
		to               string
		code             string
		marketingMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid number",
			args: args{
				ctx: context.Background(),
				to:  "this is not a valid number",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "valid number",
			args: args{
				ctx:              context.Background(),
				to:               "+25423002959",
				code:             "345",
				marketingMessage: "This is a test",
			},
			want: false,
			// TODO - investigate why an error is returned
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService()
			got, err := s.PhoneNumberVerificationCode(tt.args.ctx, tt.args.to, tt.args.code, tt.args.marketingMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PhoneNumberVerificationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.PhoneNumberVerificationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ReceiveInboundMessages(t *testing.T) {
	ctx := context.Background()
	message := &dto.TwilioMessage{
		AccountSID:       uuid.New().String(),
		From:             gofakeit.Phone(),
		To:               gofakeit.Phone(),
		Body:             gofakeit.Sentence(10),
		NumMedia:         "1",
		NumSegments:      "1",
		APIVersion:       uuid.New().String(),
		ProfileName:      gofakeit.FirstName(),
		SmsMessageSID:    uuid.New().String(),
		SmsSid:           uuid.New().String(),
		SmsStatus:        "delivered",
		WaID:             uuid.New().String(),
		MediaContentType: "image/jpeg",
		MediaURL:         gofakeit.URL(),
		TimeReceived:     time.Now(),
	}
	type args struct {
		ctx     context.Context
		message *dto.TwilioMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy inboud WA message :)",
			args: args{
				ctx:     ctx,
				message: message,
			},
			wantErr: false,
		},
		{
			name: "sad inboud WA message :)",
			args: args{
				ctx:     ctx,
				message: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newService()
			msg, err := s.ReceiveInboundMessages(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ReceiveInboundMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && msg == nil {
				t.Errorf("expected inbound message to be returned")
				return
			}
			if tt.wantErr && msg != nil {
				t.Errorf("did not expect an inbound message to be returned")
				return
			}
		})
	}
}
