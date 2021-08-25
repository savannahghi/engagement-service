package crm_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/crm"
	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/mail"
	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

const (
	phone    = "+254701223334"
	email    = "%s.users@bewell.co.ke"
	newEmail = "test@bewell.co.ke"
)

func newHubspotUsecases() *hubspotUsecases.HubSpot {
	ctx := context.Background()
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		log.Panic("failed to initialize hubspot crm repository: %w", err)
	}
	return hubspotUsecases.NewHubSpotUsecases(hubspotfr, hubspotService)
}

func newServiceCrm(ctx context.Context) *crm.Hubspot {
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		log.Panic(
			"can't instantiate firebase repository in resolver: %w",
			err,
		)
	}
	mail := mail.NewService(fr)
	hubspotUsecases := newHubspotUsecases()
	return crm.NewCrmService(hubspotUsecases, mail)
}

func testContact() (*hubspotDomain.CRMContact, error) {
	ctx := context.Background()
	phone := "+254701223334"
	hubspotUsecases := newHubspotUsecases()
	c := &hubspotDomain.CRMContact{
		Properties: hubspotDomain.ContactProperties{
			Phone: phone,
			Email: fmt.Sprintf(email, phone),
		},
	}
	contact, err := hubspotUsecases.CreateHubSpotContact(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to create tests contact: %w", err)
	}

	// Sleep to ensure the created contact propagates
	time.Sleep(5 * time.Second)

	return contact, nil
}

func TestHubspot_CollectEmails(t *testing.T) {
	ctx := context.Background()
	g := newServiceCrm(ctx)

	contact, err := testContact()
	if err != nil {
		t.Errorf("failed to create a test contact: %w", err)
		return
	}
	if contact == nil {
		t.Errorf("nil contact created")
		return
	}
	type args struct {
		ctx         context.Context
		email       string
		phonenumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "collect email happy case",
			args: args{
				ctx:         ctx,
				email:       newEmail,
				phonenumber: contact.Properties.Phone,
			},
			wantErr: false,
		},
		{
			name: "collect email sad case",
			args: args{
				ctx:         ctx,
				email:       newEmail,
				phonenumber: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact, err := g.CollectEmails(tt.args.ctx, tt.args.email, tt.args.phonenumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hubspot.CollectEmails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "collect email happy case" {
				if contact == nil {
					t.Errorf("expected a contact")
					return
				}
				if contact.Properties.Email != newEmail {
					t.Errorf("expected a contact email to be collected")
					return
				}
			}
			if tt.name == "collect email sad case" {
				if contact != nil {
					t.Errorf("did not expect a contact")
					return
				}
			}
		})
	}
}

func TestHubspot_BeWellAware(t *testing.T) {
	ctx := context.Background()
	g := newServiceCrm(ctx)

	contact, err := testContact()
	if err != nil {
		t.Errorf("failed to create a test contact: %w", err)
		return
	}
	if contact == nil {
		t.Errorf("nil contact created")
		return
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "bewell aware happy case",
			args: args{
				ctx:   ctx,
				email: newEmail,
			},
			wantErr: false,
		},
		{
			name: "bewell aware sad case",
			args: args{
				ctx:   ctx,
				email: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact, err := g.BeWellAware(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hubspot.BeWellAware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "bewell aware happy case" {
				if contact == nil {
					t.Errorf("expected a contact")
					return
				}
				if contact.Properties.BeWellAware != hubspotDomain.GeneralOptionTypeYes {
					t.Errorf("expected a contact to be bewell aware")
					return
				}
			}
			if tt.name == "collect email sad case" {
				if contact != nil {
					t.Errorf("did not expect a contact")
					return
				}
			}
		})
	}
}
