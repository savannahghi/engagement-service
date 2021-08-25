package mock

import (
	"context"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
)

// FakeServiceSMS defines the interactions with the mock sms service
type FakeServiceSMS struct {
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

// SendToMany ...
func (f *FakeServiceSMS) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendToManyFn(ctx, message, to, from)
}

// Send ...
func (f *FakeServiceSMS) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendFn(ctx, to, message, from)
}

// SendMarketingSMS ...
func (f *FakeServiceSMS) SendMarketingSMS(
	ctx context.Context,
	to []string,
	message string,
	from enumutils.SenderID,
	segment string,
) (*dto.SendMessageResponse, error) {
	return f.SendMarketingSMSFn(ctx, to, message, from, segment)
}

// SaveMarketingMessage ...
func (f *FakeServiceSMS) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return f.SaveMarketingMessageFn(ctx, data)
}

// UpdateMarketingMessage ...
func (f *FakeServiceSMS) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return f.UpdateMarketingMessageFn(ctx, data)
}

// GetMarketingSMSByPhone ...
func (f *FakeServiceSMS) GetMarketingSMSByPhone(
	ctx context.Context,
	phoneNumber string,
) (*dto.MarketingSMS, error) {
	return f.GetMarketingSMSByPhoneFn(ctx, phoneNumber)
}
