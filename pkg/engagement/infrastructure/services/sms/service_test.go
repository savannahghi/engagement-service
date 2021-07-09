package sms_test

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/sms"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func TestSendToMany(t *testing.T) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}
	crm := hubspot.NewHubSpotService()
	service := sms.NewService(fr, crm)

	type args struct {
		message string
		to      []string
		sender  base.SenderID
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "valid:successfully send to many using BeWell",
			args: args{
				message: "This is a test",
				to:      []string{"+254711223344", "+254700990099"},
				sender:  base.SenderIDBewell,
			},
			wantErr: false,
		},
		{
			name: "valid:successfully send to many using Slade260",
			args: args{
				message: "This is a test",
				to:      []string{"+254711223344", "+254700990099"},
				sender:  base.SenderIDSLADE360,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SendToMany(tt.args.message, tt.args.to, base.SenderIDBewell)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendToMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestSend(t *testing.T) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository: %w", err)
		return
	}
	crm := hubspot.NewHubSpotService()
	service := sms.NewService(fr, crm)

	type args struct {
		to      string
		message string
		sender  base.SenderID
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "valid:successfully send",
			args: args{
				message: "This is a test",
				to:      "+254711223344",
				sender:  base.SenderIDSLADE360,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail to send",
			args: args{
				message: "",
				to:      "+",
				sender:  base.SenderIDSLADE360,
			},
			wantErr: true,
		},
		{
			name: "send from an unknown sender",
			args: args{
				message: "This is a test",
				to:      "+254711223344",
				sender:  "na-kitambi-utaezana",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Send(tt.args.to, tt.args.message, tt.args.sender)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestService_SendMarketingSMS(t *testing.T) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository: %w", err)
		return
	}
	crm := hubspot.NewHubSpotService()
	s := sms.NewService(fr, crm)
	type args struct {
		ctx     context.Context
		to      []string
		message string
		from    base.SenderID
		segment string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happily send a marketing SMS :)",
			args: args{
				ctx:     context.Background(),
				to:      []string{gofakeit.Phone()},
				message: gofakeit.HipsterSentence(10),
				from:    base.SenderIDBewell,
				segment: gofakeit.UUID(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smsResp, err := s.SendMarketingSMS(tt.args.ctx, tt.args.to, tt.args.message, tt.args.from, tt.args.segment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendMarketingSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if smsResp == nil {
				t.Errorf("expected an sms response to be returned")
				return
			}
		})
	}
}
