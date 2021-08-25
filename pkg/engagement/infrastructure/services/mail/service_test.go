package mail_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement-service/pkg/engagement/repository"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func TestNewService(t *testing.T) {
	ctx := context.Background()

	repo, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("error initializing new firebase repo:%v", err)
		return
	}

	service := mail.NewService(repo)
	tests := []struct {
		name string
		want *mail.Service
	}{
		{
			name: "default case",
			want: service,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mail.NewService(repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			} else {
				got.CheckPreconditions()
			}
		})
	}
}

func TestService_SendEmail(t *testing.T) {
	testUserMail := "test@bewell.co.ke"
	ctx := context.Background()

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("error initializing new firebase repo:%v", err)
		return
	}

	service := mail.NewService(fr)
	tests := []struct {
		name    string
		service *mail.Service
		subject string
		text    string
		to      []string

		expectMsg bool
		expectID  bool
		expectErr bool
	}{
		{
			name:    "valid email",
			service: service,
			subject: "Test Email",
			text:    "Test Email",
			to:      []string{testUserMail},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, id, err := tt.service.SendEmail(ctx, tt.subject, tt.text, nil, tt.to...)
			if tt.expectErr {
				if err == nil {
					t.Errorf("an error was expected")
					return
				}
				if msg != "" && id != "" {
					t.Errorf("expected no message and message ID")
					return
				}
			}
			if !tt.expectErr {
				if err != nil {
					t.Errorf("an error was not expected")
					return
				}
				if msg == "" && id == "" {
					t.Errorf("expected a message and message ID")
					return
				}
			}
		})
	}
}

func TestService_SendInBlue(t *testing.T) {
	ctx := context.Background()
	var repo repository.Repository
	type args struct {
		subject string
		text    string
		to      []string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus string
		wantErr    bool
	}{
		{
			name: "happy case",
			args: args{
				subject: "Test Email",
				text:    "This is a test email",
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mail.NewService(repo)
			if s.SendInBlueEnabled {
				got, _, err := s.SendInBlue(ctx, tt.args.subject, tt.args.text, tt.args.to...)
				if (err != nil) != tt.wantErr {
					t.Errorf("Service.SendInBlue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.wantStatus {
					t.Errorf("Service.SendInBlue() got = %v, want %v", got, tt.wantStatus)
				}
			}
		})
	}
}

func TestService_CheckPreconditions(t *testing.T) {
	var repo repository.Repository
	type fields struct {
		Mg                *mailgun.MailgunImpl
		From              string
		SendInBlueEnabled bool
		SendInBlueAPIKey  string
		Repository        repository.Repository
	}
	tests := []struct {
		name   string
		fields fields
		panics bool
	}{
		{
			name: "invalid: missing mailgin implementation",
			fields: fields{
				From:              "test@example",
				SendInBlueEnabled: false,
				SendInBlueAPIKey:  "key",
				Repository:        repo,
			},
			panics: true,
		},
		{
			name: "invalid: missing from",
			fields: fields{
				Mg:                &mailgun.MailgunImpl{},
				SendInBlueEnabled: false,
				SendInBlueAPIKey:  "key",
				Repository:        repo,
			},
			panics: true,
		},
		{
			name: "invalid: missing sendiblue api key",
			fields: fields{
				Mg:                &mailgun.MailgunImpl{},
				From:              "test@example",
				SendInBlueEnabled: false,
				Repository:        repo,
			},
			panics: true,
		},
		{
			name: "invalid: missing repo",
			fields: fields{
				Mg:                &mailgun.MailgunImpl{},
				From:              "test@example",
				SendInBlueEnabled: false,
				SendInBlueAPIKey:  "key",
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mail.Service{
				Mg:                tt.fields.Mg,
				From:              tt.fields.From,
				SendInBlueEnabled: tt.fields.SendInBlueEnabled,
				SendInBlueAPIKey:  tt.fields.SendInBlueAPIKey,
				Repository:        tt.fields.Repository,
			}
			if tt.panics {
				assert.Panics(t, s.CheckPreconditions)
			}

		})
	}
}
