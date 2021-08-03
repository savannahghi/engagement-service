package mock

import (
	"context"
	"net/http"
)

// FakeEDIService is an `edi` service mock
type FakeEDIService struct {
	UpdateMessageSentFn func(
		ctx context.Context,
		phoneNumber string,
		segment string,
	) (*http.Response, error)
}

// UpdateMessageSent mocks the update message sent method
func (f *FakeEDIService) UpdateMessageSent(
	ctx context.Context,
	phoneNumber string,
	segment string,
) (*http.Response, error) {
	return f.UpdateMessageSentFn(ctx, phoneNumber, segment)
}
