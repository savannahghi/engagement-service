package mock

import (
	"context"
	"net/url"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

//FakeTwilioService ...
type FakeTwilioService struct {
	MakeTwilioRequestFn func(
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error

	RoomFn func(ctx context.Context) (*dto.Room, error)

	TwilioAccessTokenFn func(ctx context.Context) (*dto.AccessToken, error)

	SendSMSFn func(ctx context.Context, to string, msg string) error
}

//MakeTwilioRequest mock function
func (f *FakeTwilioService) MakeTwilioRequest(method string, urlPath string, content url.Values, target interface{}) error {
	return f.MakeTwilioRequestFn(method, urlPath, content, target)
}

//Room mock function
func (f *FakeTwilioService) Room(ctx context.Context) (*dto.Room, error) {
	return f.RoomFn(ctx)
}

//TwilioAccessToken mock function
func (f *FakeTwilioService) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return f.TwilioAccessTokenFn(ctx)
}

//SendSMS mock function
func (f *FakeTwilioService) SendSMS(ctx context.Context, to string, msg string) error {
	return f.SendSMSFn(ctx, to, msg)
}
