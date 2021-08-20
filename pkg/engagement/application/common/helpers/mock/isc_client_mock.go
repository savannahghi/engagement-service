package mock

import "github.com/savannahghi/interserviceclient"

// FakeIscClient ...
type FakeIscClient struct {
	// InitializeInterServiceClientFn ...
	InitializeInterServiceClientFn func(serviceName string) *interserviceclient.InterServiceClient
}

// InitializeInterServiceClient is a mock version of the original function
func (c *FakeIscClient) InitializeInterServiceClient(
	serviceName string,
) *interserviceclient.InterServiceClient {
	return c.InitializeInterServiceClientFn(serviceName)
}
