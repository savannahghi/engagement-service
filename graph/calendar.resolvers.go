package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/engagement/graph/generated"
	"gitlab.slade360emr.com/go/engagement/graph/model"
	calendar "google.golang.org/api/calendar/v3"
)

func (r *calendarEventResolver) Attachments(ctx context.Context, obj *calendar.Event) ([]*model.EventAttachment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *calendarEventResolver) Attendees(ctx context.Context, obj *calendar.Event) ([]*model.EventAttendee, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *calendarEventResolver) OriginalStartTime(ctx context.Context, obj *calendar.Event) (*model.EventDateTime, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *calendarEventResolver) Start(ctx context.Context, obj *calendar.Event) (*model.EventDateTime, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *calendarEventResolver) End(ctx context.Context, obj *calendar.Event) (*model.EventDateTime, error) {
	panic(fmt.Errorf("not implemented"))
}

// CalendarEvent returns generated.CalendarEventResolver implementation.
func (r *Resolver) CalendarEvent() generated.CalendarEventResolver { return &calendarEventResolver{r} }

type calendarEventResolver struct{ *Resolver }
