package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/dto"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/repository"
)

type MarketingDataUseCases interface {
	GetMarketingData(ctx context.Context, data *dto.MarketingMessagePayload) ([]*dto.Segment, error)
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
