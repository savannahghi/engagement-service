package mock

import (
	"context"
	"net/url"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

// FakeServiceTwilio defines the interaction with the twilio service
type FakeServiceTwilio struct {
	MakeTwilioRequestFn func(
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error

	RoomFn func(ctx context.Context) (*dto.Room, error)

	TwilioAccessTokenFn func(ctx context.Context) (*dto.AccessToken, error)

	SendSMSFn func(ctx context.Context, to string, msg string) error

	SaveTwilioVideoCallbackStatusFn func(
		ctx context.Context,
		data dto.CallbackData,
	) error
}

// MakeTwilioRequest ...
func (f *FakeServiceTwilio) MakeTwilioRequest(
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	return f.MakeTwilioRequestFn(method, urlPath, content, target)
}

// Room ...
func (f *FakeServiceTwilio) Room(ctx context.Context) (*dto.Room, error) {
	return f.RoomFn(ctx)
}

// TwilioAccessToken ...
func (f *FakeServiceTwilio) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return f.TwilioAccessTokenFn(ctx)
}

// SendSMS ...
func (f *FakeServiceTwilio) SendSMS(ctx context.Context, to string, msg string) error {
	return f.SendSMSFn(ctx, to, msg)
}

// SaveTwilioVideoCallbackStatus ..
func (f *FakeServiceTwilio) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return f.SaveTwilioVideoCallbackStatusFn(ctx, data)
}
