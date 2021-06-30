package usecases

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

type MarketingDataUseCases interface {
	GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error)
	UpdateUserCRMEmail(ctx context.Context, email string, phonenumber string) error
	BeWellAware(ctx context.Context, email string) error
}

// MarketingDataImpl represents the marketing usecase implementation
type MarketingDataImpl struct {
	repository repository.Repository
}

// NewMarketing initialises a marketing usecase
func NewMarketing(
	repository repository.Repository,
) *MarketingDataImpl {
	return &MarketingDataImpl{
		repository: repository,
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
