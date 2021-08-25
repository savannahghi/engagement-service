package helpers_test

import (
	"os"
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	originalEnvironment := os.Getenv("ENVIRONMENT")

	os.Setenv("ENVIRONMENT", "staging")

	code := m.Run()

	os.Setenv("ENVIRONMENT", originalEnvironment)

	os.Exit(code)
}

func TestInitializeInterServiceClient(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
	}{
		{
			name: "happy case",
			args: args{
				serviceName: "profile",
			},
			wantValue: true,
		},
		{
			name:      "sad case",
			args:      args{},
			wantValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.InitializeInterServiceClient(tt.args.serviceName)
			if tt.wantValue {
				assert.NotNil(t, got)
				assert.NotEmpty(t, got)
			}
			if tt.wantValue && got == nil {
				t.Errorf("expected value, got %v", got)
			}

		})
	}
}
