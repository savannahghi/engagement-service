package scheduling

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// const calendarTestEmail = firebasetools.TestUserEmail

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func TestNewService(t *testing.T) {
	assert.NotPanics(t, func() {
		srv := NewService()
		srv.checkPreconditions()
	})
}

func TestService_CreateCalendar(t *testing.T) {
	ctx := context.Background()
	srv := NewService()
	assert.NotNil(t, srv)
	srv.checkPreconditions()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{  
		// TODO: restore after resolving issue
		// {
		// 	name: "successful test",
		// 	args: args{
		// 		name: "Mock calendar for testing",
		// 	},
		// 	wantErr: false, 
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.CreateCalendar(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateCalendar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.args.name, got.Summary)

			deleteCall := s.gcalService.Calendars.Delete(got.Id)
			err = deleteCall.Do()
			assert.Nil(t, err)
		})
	}
}

func TestService_AddEvent(t *testing.T) {
	ctx := context.Background()
	srv := NewService()
	assert.NotNil(t, srv)
	srv.checkPreconditions()

	cal, err := srv.CreateCalendar(ctx, "integration test calendar")
	assert.Nil(t, err)
	assert.NotNil(t, cal)

	defer func() {
		// clean up
		deleteCall := srv.gcalService.Calendars.Delete(cal.Id)
		err = deleteCall.Do()
		assert.Nil(t, err)
	}()

	// start := time.Now().Add(time.Hour)
	// end := start.Add(time.Hour)
	type args struct {
		title              string
		description        string
		start              time.Time
		end                time.Time
		attendeeEmails     []string
		extendedProperties map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: restore after resolving issue
		// {
		// 	name: "valid case",
		// 	args: args{
		// 		title:          "Integration Test Event",
		// 		description:    "Integration Test Event Description",
		// 		start:          start,
		// 		end:            end,
		// 		attendeeEmails: []string{calendarTestEmail},
		// 		extendedProperties: map[string]string{
		// 			"a": "1",
		// 			"b": "2",
		// 		},
		// 	},
		// 	wantErr: false,  
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.AddEvent(
				ctx,
				tt.args.title,
				tt.args.description,
				tt.args.start,
				tt.args.end,
				tt.args.attendeeEmails,
				tt.args.extendedProperties,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
