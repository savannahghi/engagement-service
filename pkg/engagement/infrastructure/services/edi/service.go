package edi

import (
	"context"
	"net/http"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/interserviceclient"
	"gitlab.slade360emr.com/go/apiclient"
)

const (
	// UpdateMessageSent ISC endpoint to update the `message sent` value
	UpdateMessageSent = "internal/update_message_sent"

	ediService = "edi"
)

// ServiceEdi defines the business logic required to interact with EDI
type ServiceEdi interface {
	UpdateMessageSent(
		ctx context.Context,
		phoneNumber string,
		segment string,
	) (*http.Response, error)
}

// ServiceEDIImpl uses inter-service REST APIs to fetch information
// from a remote edi service
type ServiceEDIImpl struct {
	EdiExt *interserviceclient.InterServiceClient
}

// NewEDIClient initializes a new interservice client for edi
func NewEDIClient() *interserviceclient.InterServiceClient {
	return helpers.InitializeInterServiceClient(ediService)
}

// NewEdiService returns a new instance of edi implementations
func NewEdiService(
	edi *interserviceclient.InterServiceClient,
) ServiceEdi {
	return &ServiceEDIImpl{
		EdiExt: edi,
	}
}

// UpdateMessageSent calls the `EDI` service to update the value of message sent
func (s *ServiceEDIImpl) UpdateMessageSent(
	ctx context.Context,
	phoneNumber string,
	segment string,
) (*http.Response, error) {
	payload := apiclient.UpdateMessageSentInput{
		PhoneNumber: phoneNumber,
		Segment:     segment,
	}
	return s.EdiExt.MakeRequest(
		ctx,
		http.MethodPost,
		UpdateMessageSent,
		payload,
	)
}
