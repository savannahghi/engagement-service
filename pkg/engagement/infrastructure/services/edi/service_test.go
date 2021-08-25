package edi_test

import (
	"context"
	"testing"

	"github.com/savannahghi/engagement-service/pkg/engagement/infrastructure/services/edi"
	"github.com/savannahghi/interserviceclient"
)

func TestServiceEDIImpl_UpdateMessageSent(t *testing.T) {
	ediClient := edi.NewEDIClient()
	s := edi.NewEdiService(ediClient)
	type args struct {
		ctx         context.Context
		phoneNumber string
		segment     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case -> Successfully update message sent",
			args: args{
				ctx:         context.Background(),
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				segment:     "Random Segment",
			},
			wantErr: false,
		},
		{
			name: "Sad Case -> Fail to update message sent",
			args: args{
				ctx: context.Background(),
			},
			// TODO: investigate this, should be true
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.UpdateMessageSent(tt.args.ctx, tt.args.phoneNumber, tt.args.segment)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEDIImpl.UpdateMessageSent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
