package helpers_test

import (
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	mock "github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers/mock"
	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
)

var fakeIscClient mock.FakeIscClient

func Test_InitializeInterServiceClient(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
	}{
		{
			name: "happy case: valid interservice client",
			args: args{
				serviceName: "profile",
			},
			wantValue: true,
		},
		{
			name:      "sad case: missing interservice client",
			args:      args{},
			wantValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case: valid interservice client" {
				fakeIscClient.InitializeInterServiceClientFn = func(
					serviceName string,
				) *interserviceclient.InterServiceClient {
					return &interserviceclient.InterServiceClient{
						Name:              "test-service",
						RequestRootDomain: "http://example.com",
					}
				}
				got := fakeIscClient.InitializeInterServiceClientFn(tt.args.serviceName)
				assert.NotEmpty(t, got.Name)
				assert.NotEmpty(t, got.RequestRootDomain)
			}

			if tt.name == "sad case: missing interservice client" {
				fakeIscClient.InitializeInterServiceClientFn = func(
					serviceName string,
				) *interserviceclient.InterServiceClient {
					return &interserviceclient.InterServiceClient{
						Name:              "",
						RequestRootDomain: "",
					}
				}
				got := fakeIscClient.InitializeInterServiceClientFn(tt.args.serviceName)
				assert.Empty(t, got.Name)
				assert.Empty(t, got.RequestRootDomain)
			}

			got := helpers.InitializeInterServiceClient(tt.args.serviceName)
			if tt.wantValue {
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.Name)
				assert.NotEmpty(t, got.RequestRootDomain)
			}
			if !tt.wantValue {
				assert.Empty(t, got.Name)
				assert.Empty(t, got.RequestRootDomain)
			}
		})
	}
}
