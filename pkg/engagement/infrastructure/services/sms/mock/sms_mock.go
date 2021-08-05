package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
)

//FakeSMSService ...
type FakeSMSService struct {
	SendToManyFn func(
		ctx context.Context,
		message string,
		to []string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	SendFn func(
		ctx context.Context,
		to, message string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	SendMarketingSMSFn func(
		ctx context.Context,
		to []string,
		message string,
		from enumutils.SenderID,
		segment string,
	) (*dto.SendMessageResponse, error)

	SaveMarketingMessageFn func(
		ctx context.Context,
		data dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

	UpdateMarketingMessageFn func(
		ctx context.Context,
		data *dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

	GetMarketingSMSByPhoneFn func(
		ctx context.Context,
		phoneNumber string,
	) (*dto.MarketingSMS, error)
}

//SendToMany mock function
func (f *FakeSMSService) SendToMany(ctx context.Context, message string, to []string, from enumutils.SenderID) (*dto.SendMessageResponse, error) {
	return f.SendToManyFn(ctx, message, to, from)
}

//Send mock function
func (f *FakeSMSService) Send(ctx context.Context, to, message string, from enumutils.SenderID) (*dto.SendMessageResponse, error) {
	return f.SendFn(ctx, to, message, from)
}

//SendMarketingSMS mock function
func (f *FakeSMSService) SendMarketingSMS(ctx context.Context, to []string, message string, from enumutils.SenderID, segment string) (*dto.SendMessageResponse, error) {
	return f.SendMarketingSMSFn(ctx, to, message, from, segment)
}

//SaveMarketingMessage mock function
func (f *FakeSMSService) SaveMarketingMessage(ctx context.Context, data dto.MarketingSMS) (*dto.MarketingSMS, error) {
	return f.SaveMarketingMessageFn(ctx, data)
}

//UpdateMarketingMessage mock function
func (f *FakeSMSService) UpdateMarketingMessage(ctx context.Context, data *dto.MarketingSMS) (*dto.MarketingSMS, error) {
	return f.UpdateMarketingMessageFn(ctx, data)
}

//GetMarketingSMSByPhone mock function
func (f *FakeSMSService) GetMarketingSMSByPhone(ctx context.Context, phoneNumber string) (*dto.MarketingSMS, error) {
	return f.GetMarketingSMSByPhoneFn(ctx, phoneNumber)
}
