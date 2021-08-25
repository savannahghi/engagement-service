package sms_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/database"
	crmExt "github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/edi"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func newHubSpotService(ctx context.Context) *hubspotUsecases.HubSpot {
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		log.Panic("failed to initialize hubspot crm repository: %w", err)
	}
	return hubspotUsecases.NewHubSpotUsecases(hubspotfr, hubspotService)
}

func newTestSMSService() (*sms.Service, error) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	hubspotUsecases := newHubSpotService(ctx)
	mail := mail.NewService(fr)
	crmExt := crmExt.NewCrmService(hubspotUsecases, mail)
	edi := edi.NewEdiService(edi.NewEDIClient())
	ps, err := messaging.NewPubSubNotificationService(
		ctx,
		serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate notification service in resolver: %w",
			err,
		)
	}
	return sms.NewService(fr, crmExt, ps, edi), nil
}

func TestSendToMany(t *testing.T) {
	ctx := context.Background()
	service, err := newTestSMSService()
	if err != nil {
		t.Errorf("unable to initialize test service with error %v", err)
		return
	}

	type args struct {
		message string
		to      []string
		sender  enumutils.SenderID
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
				sender:  enumutils.SenderIDBewell,
			},
			wantErr: false,
		},
		{
			name: "valid:successfully send to many using Slade260",
			args: args{
				message: "This is a test",
				to:      []string{"+254711223344", "+254700990099"},
				sender:  enumutils.SenderIDSLADE360,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SendToMany(ctx, tt.args.message, tt.args.to, enumutils.SenderIDBewell)
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
	service, err := newTestSMSService()
	if err != nil {
		t.Errorf("unable to initialize test service with error %v", err)
		return
	}

	type args struct {
		to      string
		message string
		sender  enumutils.SenderID
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
				sender:  enumutils.SenderIDSLADE360,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail to send",
			args: args{
				message: "",
				to:      "+",
				sender:  enumutils.SenderIDSLADE360,
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
			got, err := service.Send(ctx, tt.args.to, tt.args.message, tt.args.sender)
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
	ctx := firebasetools.GetAuthenticatedContext(t)
	s, err := newTestSMSService()
	if err != nil {
		t.Errorf("unable to initialize test service with error %v", err)
		return
	}
	to := []string{"+254711223344", "+254700990099"}

	for _, phone := range to {
		hs := newHubSpotService(ctx)
		_, err := hs.CreateHubSpotContact(ctx, &hubspotDomain.CRMContact{Properties: hubspotDomain.ContactProperties{Phone: phone}})
		if err != nil {
			t.Errorf("failed to create test contact: %w", err)
			return
		}
	}

	type args struct {
		ctx     context.Context
		to      []string
		message string
		from    enumutils.SenderID
		segment string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "Happily send a marketing SMS :)",
			args: args{
				ctx:     ctx,
				to:      to,
				message: gofakeit.HipsterSentence(10),
				from:    enumutils.SenderIDBewell,
				segment: "WING A",
			},
			wantErr: false,
		},
		{
			name: "Sad Case missing message :(",
			args: args{
				ctx:     ctx,
				to:      to,
				message: "",
				from:    enumutils.SenderIDBewell,
				segment: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case missing sender :(",
			args: args{
				ctx:     ctx,
				to:      to,
				message: gofakeit.HipsterSentence(10),
				from:    "",
				segment: "WING A",
			},
			wantErr: true,
		},
		{
			name: "Sad Case invalid sender :(",
			args: args{
				ctx:     ctx,
				to:      to,
				message: gofakeit.HipsterSentence(10),
				from:    "invalid",
				segment: "WING A",
			},
			wantErr: true,
		},
		{
			name: "Sad Case missing recipient :(",
			args: args{
				ctx:     ctx,
				to:      []string{""},
				message: gofakeit.HipsterSentence(10),
				from:    enumutils.SenderIDBewell,
				segment: "WING A",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SendMarketingSMS(tt.args.ctx, tt.args.to, tt.args.message, tt.args.from, tt.args.segment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendMarketingSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
