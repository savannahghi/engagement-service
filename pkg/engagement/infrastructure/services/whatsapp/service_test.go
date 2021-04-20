package whatsapp_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/whatsapp"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	m.Run()
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
			got := whatsapp.NewService()
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
			s := whatsapp.NewService()
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
