package usecases

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/authorization/permission"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

const (
	campaignDataFileName = "campaign.dataset.csv"
)

type MarketingDataUseCases interface {
	GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error)
	UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error
	BeWellAware(ctx context.Context, email string) error
	LoadCampaignDataset(ctx context.Context, phone string) error
}

// MarketingDataImpl represents the marketing usecase implementation
type MarketingDataImpl struct {
	repository repository.Repository
	hubspot    hubspot.ServiceHubSpotInterface
}

// NewMarketing initialises a marketing usecase
func NewMarketing(
	repository repository.Repository, hubspot hubspot.ServiceHubSpotInterface,
) *MarketingDataImpl {
	return &MarketingDataImpl{
		repository: repository,
		hubspot:    hubspot,
	}
}

func (m MarketingDataImpl) GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error) {
	segmentData, err := m.repository.RetrieveMarketingData(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the marketing data")
	}

	return segmentData, nil
}

// UpdateUserCRMEmail updates a user CRM contact with the supplied email
func (m MarketingDataImpl) UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error {
	CRMContactProperties := CRMDomain.ContactProperties{
		Email: email,
	}

	if err := m.repository.UpdateUserCRMEmail(ctx, phonenumber, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
		Phone:      phonenumber,
	}); err != nil {
		return fmt.Errorf("failed to create CRM staging payload %v", err)
	}
	return nil

}

//BeWellAware toggles the user identified by the provided email= as bewell-aware on the CRM
func (m MarketingDataImpl) BeWellAware(ctx context.Context, email string) error {
	CRMContactProperties := CRMDomain.ContactProperties{
		BeWellAware: CRMDomain.GeneralOptionTypeYes,
	}

	logrus.Printf("bewellAware payload: %v with email: %v", CRMContactProperties, email)

	if err := m.repository.UpdateUserCRMBewellAware(ctx, email, &dto.UpdateContactPSMessage{
		Properties: CRMContactProperties,
	}); err != nil {
		return fmt.Errorf("failed to create CRM staging payload %v", err)
	}
	return nil
}

// LoadCampaignDataset publishes the campaign dataset into firestore and CRM
// this call is idempotent,hence its safeguards against duplicates on both firstore and CRM
// this function expects that the CSV containing the dataset has been preprocessed and contains
// attributes expected by firestore and CRM
func (m MarketingDataImpl) LoadCampaignDataset(ctx context.Context, phone string) error {
	if p := base.StringSliceContains(base.AuthorizedPhones, phone); !p {
		return fmt.Errorf("not authorized to access this resource")
	}
	isAuthorized, err := authorization.IsAuthorized(&base.UserInfo{
		PhoneNumber: phone,
	}, permission.LoadMarketingData)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, campaignDataFileName)

	csvFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening the CSV file: %w", err)
	}
	defer csvFile.Close()

	csvContent, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read data from the CSV file :%w", err)
	}

	csvPayload := csvContent[1:]

	count := 0

	for _, line := range csvPayload {
		data := dto.Segment{
			BeWellEnrolled:        line[0],
			OptOut:                line[1],
			BeWellAware:           line[2],
			BeWellPersona:         line[3],
			HasWellnessCard:       line[4],
			HasCover:              line[5],
			Payor:                 line[6],
			FirstChannelOfContact: line[7],
			InitialSegment:        line[8],
			HasVirtualCard:        line[9],
			Email:                 line[10],
			PhoneNumber:           line[11],
			FirstName:             line[12],
			LastName:              line[13],
			Wing:                  line[14],
			MessageSent:           line[15],
			IsSynced:              line[16],
			TimeSynced:            line[17],
		}

		logrus.Printf("Loading date of email %v", data.Email)

		// publish to firestore
		if err := m.repository.LoadMarketingData(ctx, data); err != nil {
			return fmt.Errorf("%v", err)
		}

		res, err := m.hubspot.SearchContactByPhone(line[11])
		if err != nil {
			_ = m.repository.RollBackMarketingData(ctx, data)
			continue
		}

		// contact already exists. Nothing to do here
		if len(res.Results) >= 1 {
			continue
		}

		convertor := func(d string) CRMDomain.GeneralOptionType {
			if d == "YES" {
				return CRMDomain.GeneralOptionTypeYes
			}
			return CRMDomain.GeneralOptionTypeNo
		}

		convertorPersona := func(d string) CRMDomain.Persona {
			if d == "ALICE" {
				return CRMDomain.PersonaAlice
			}

			if d == "JUMA" {
				return CRMDomain.PersonaJuma
			}

			if d == "BOB" {
				return CRMDomain.PersonaBob
			}

			if d == "ANDREW" {
				return CRMDomain.PersonaAndrew
			}
			return CRMDomain.PersonaSlader
		}

		convertorPayor := func(d string) CRMDomain.Payor {
			if d == "RESOLUTION" || d == "RESOLUTION INSURANCE" {
				return CRMDomain.PayorResolution
			}

			if d == "APA" || d == "APA INSURANCE" {
				return CRMDomain.PayorApa
			}

			if d == "JUBILEE" || d == "JUBILEE INSURANCE" {
				return CRMDomain.PayorJubilee
			}

			if d == "BRITAM" || d == "BRITAM INSURANCE" {
				return CRMDomain.PayorBritam
			}

			if d == "MADISON" || d == "MADISON INSURANCE" {
				return CRMDomain.PayorMadison
			}

			return CRMDomain.PayorJubilee
		}

		convertorChannel := func(d string) CRMDomain.ChannelOfContact {
			if d == "APP" {
				return CRMDomain.ChannelOfContactApp
			}
			if d == "USSD" {
				return CRMDomain.ChannelOfContactUssd
			}
			return CRMDomain.ChannelOfContactShortcode
		}

		if _, err := m.hubspot.CreateContact(CRMDomain.CRMContact{
			Properties: CRMDomain.ContactProperties{
				BeWellEnrolled:        convertor(line[0]),
				OptOut:                convertor(line[1]),
				BeWellAware:           convertor(line[2]),
				BeWellPersona:         convertorPersona(line[3]),
				HasSladeID:            convertor(line[4]),
				HasCover:              convertor(line[5]),
				Payor:                 convertorPayor(line[6]),
				FirstChannelOfContact: convertorChannel(line[7]),
				HasVirtualCard:        convertor(line[9]),
				Email:                 line[10],
				Phone:                 line[11],
				FirstName:             line[12],
				LastName:              line[13],
			},
		}); err != nil {
			_ = m.repository.RollBackMarketingData(ctx, data)
			continue
		}

		count++
	}

	logrus.Printf("Records successfully created on firestore and CRM %v", count)

	return nil
}
