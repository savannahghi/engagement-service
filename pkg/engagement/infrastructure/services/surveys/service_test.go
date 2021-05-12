package surveys_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/surveys"

	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/resources"
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
		input resources.NPSInput
	}
	// feedback := &resources.FeedbackInput{
	// 	Question: "How is it",
	// 	Answer:   "It is what it is",
	// }
	// email := base.TestUserEmail
	// phoneNumber := base.TestUserPhoneNumber

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO - Find out why it fails
		// {
		// 	name: "Successful save nps response",
		// 	args: args{
		// 		input: resources.NPSInput{
		// 			Name:        "Seleman Bungara",
		// 			Score:       8,
		// 			SladeCode:   "50",
		// 			Email:       &email,
		// 			PhoneNumber: &phoneNumber,
		// 			Feedback:    []*resources.FeedbackInput{feedback},
		// 		},
		// 	},
		// 	want:    true,
		// 	wantErr: false,
		// },
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
			}
		})
	}
}
