package scheduling

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenSource(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTokenSource(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
