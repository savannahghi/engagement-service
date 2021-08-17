package surveys_test

import (
	"context"
	"testing"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/surveys"	
	"github.com/savannahghi/firebasetools" 
	"github.com/savannahghi/interserviceclient"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
)

func TestService_RecordNPSResponse(t *testing.T) {
	ctx := context.Background()
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}
	service := surveys.NewService(fr)
	type args struct {
		ctx   context.Context
		input dto.NPSInput
	}
	feedback := &dto.FeedbackInput{
		Question: "How is it",
		Answer:   "It is what it is",
	}
	email := firebasetools.TestUserEmail
	phoneNumber := interserviceclient.TestUserPhoneNumber

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Successfully save NPS response",
			args: args{
				ctx:   context.Background(),
				input: dto.NPSInput{
					Name:        "Seleman Bungara",
					Score:       8,
					SladeCode:   "50",
					Email:       &email,
					PhoneNumber: &phoneNumber,
					Feedback:    []*dto.FeedbackInput{feedback},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.RecordNPSResponse(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RecordNPSResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.RecordNPSResponse() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
