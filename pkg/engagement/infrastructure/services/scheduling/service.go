package scheduling

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/engagement-service/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement-service/pkg/engagement/services/scheduling")

// calendar app constants
const (
	DefaultCalendarEmail                   = "be.well@healthcloud.co.ke"
	DefaultCalendarDisplayName             = "Be.Well"
	DefaultCalendarTimezone                = "Africa/Nairobi"
	DefaultCalendarDescription             = "Auto-created calendar"
	DefaultCalendarVisibility              = "private"
	DefaultCalendarLocked                  = false
	DefaultCalendarGuestsCanSeeOtherGuests = true
	DefaultGuestsCanModify                 = true
	DefaultGuestsCanInviteOthers           = true
	DefaultCalendarEmailReminderMinutes    = 30
	DefaultCalendarPopupReminderMinutes    = 10
	DefaultCoordinates                     = "-1.288018, 36.783740"
)

// NewService initializes a Google Calendar service
func NewService() *Service {
	ctx := context.Background()
	ts, err := GetTokenSource(ctx)
	if err != nil {
		log.Panicf("unable to initialize token source for the Google Calendar API: %s", err)
	}

	// the Google calendar API has a rate limiter that makes it unwise to do actual
	// API calls in automated tests. It's also a b***** to set up due to the
	// immaturity of the Go library at the time of writing and the complicated
	// OAUth 2 "delegated permissions" model. To cap off a bad situation, the
	// automatically generated "official" Go client is difficult to mock.
	var gcalService *calendar.Service
	if serverutils.IsRunningTests() {
		mockClient := MockGCALHTTPClient()
		gcalService, err = calendar.NewService(ctx, option.WithHTTPClient(mockClient))
		if err != nil {

			log.Panicf("unable to create mock Google Calendar service: %s", err)
		}
	} else {
		gcalService, err = calendar.NewService(ctx, option.WithTokenSource(ts))
		if err != nil {

			log.Panicf("unable to create Google Calendar service: %s", err)
		}
	}

	fc := &firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {

		log.Printf("unable to initialize Firebase app: %s", err)
	}
	firestore, err := fa.Firestore(ctx)
	if err != nil {

		log.Printf("unable to initialize Firestore: %s", err)
	}

	service := &Service{
		gcalService: gcalService,
		firestore:   firestore,
	}
	return service
}

// Service offers methods to interact with the Google Calendar API
type Service struct {
	gcalService *calendar.Service
	firestore   *firestore.Client
}

func (s Service) checkPreconditions() {
	if s.gcalService == nil {
		log.Panicf("scheduling service: nil Google calendar service")
	}

	if s.firestore == nil {
		log.Panicf("scheduling service: nil firestore")
	}
}

// CreateCalendar creates a calendar e.g for a specific provider or practitioner.
//
// At the time of writing, the global per user rate limit was 500 queries/100 seconds/user.
//
// Also: If you create more than 60 new calendars in a short period, your
// calendar might go into read-only mode for several hours.
func (s Service) CreateCalendar(ctx context.Context, name string) (*calendar.Calendar, error) {
	_, span := tracer.Start(ctx, "CreateCalendar")
	defer span.End()
	s.checkPreconditions()

	inp := &calendar.Calendar{
		Summary:     name,
		TimeZone:    DefaultCalendarTimezone,
		Description: DefaultCalendarDescription,
		ConferenceProperties: &calendar.ConferenceProperties{
			AllowedConferenceSolutionTypes: []string{"hangoutsMeet"},
		},
	}
	insertCall := s.gcalService.Calendars.Insert(inp)
	calendar, err := insertCall.Do()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to insert calendar: %w", err)
	}
	return calendar, nil
}

// FreeBusy returns the free/busy slots of the indicated calendar between
// the indicated time
func (s Service) FreeBusy(
	ctx context.Context,
	calendarIDs []string,
	start time.Time,
	end time.Time,
) (map[string][]*BusySlot, error) {
	_, span := tracer.Start(ctx, "FreeBusy")
	defer span.End()
	s.checkPreconditions()

	items := []*calendar.FreeBusyRequestItem{}
	for _, calendarID := range calendarIDs {
		item := &calendar.FreeBusyRequestItem{
			Id: calendarID,
		}
		items = append(items, item)
	}
	req := &calendar.FreeBusyRequest{
		TimeMin:  start.Format(time.RFC3339),
		TimeMax:  end.Format(time.RFC3339),
		TimeZone: DefaultCalendarTimezone,
		Items:    items,
	}
	call := s.gcalService.Freebusy.Query(req)
	resp, err := call.Do()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get free/busy for calendar(s): %w", err)
	}
	output := make(map[string][]*BusySlot)
	for calendarID, calFreeBusy := range resp.Calendars {
		for _, busyPeriod := range calFreeBusy.Busy {
			timeFormat := "2006-01-02T15:04:05-07:00"
			start, err := time.Parse(timeFormat, busyPeriod.Start)
			if err == nil {
				return nil, fmt.Errorf(
					"can't parse Google calendar busy period start %s into time: %w", busyPeriod.Start, err)
			}
			end, err := time.Parse(timeFormat, busyPeriod.End)
			if err == nil {
				return nil, fmt.Errorf(
					"can't parse Google calendar busy period end %s into time: %w", busyPeriod.End, err)
			}
			busySlot := &BusySlot{
				Start: start,
				End:   end,
			}
			output[calendarID] = append(output[calendarID], busySlot)
		}
	}
	return output, nil
}

// AddEvent adds an event to a calendar.
func (s Service) AddEvent(
	ctx context.Context,
	title string,
	description string,
	start time.Time,
	end time.Time,
	attendeeEmails []string,
	privateExtendedPropertiesContent map[string]string,
) (*calendar.Event, error) {
	_, span := tracer.Start(ctx, "SendOTPToEmail")
	defer span.End()
	s.checkPreconditions()

	if len(attendeeEmails) == 0 {
		return nil, fmt.Errorf("an event can't be created with no attendeeEmails")
	}

	calendarID := "primary"
	calGet := s.gcalService.Calendars.Get(calendarID)
	cal, err := calGet.Do()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get Google Calendar with ID %s: %w", calendarID, err)
	}

	locked := DefaultCalendarLocked
	guestsCanSeeOtherGuests := DefaultCalendarGuestsCanSeeOtherGuests
	guestsCanModify := DefaultGuestsCanModify
	guestsCanInviteOthers := DefaultGuestsCanInviteOthers
	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: DefaultCalendarTimezone,
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: DefaultCalendarTimezone,
		},
		GuestsCanInviteOthers:   &guestsCanInviteOthers,
		GuestsCanModify:         guestsCanModify,
		GuestsCanSeeOtherGuests: &guestsCanSeeOtherGuests,
		Locked:                  locked,
		Visibility:              DefaultCalendarVisibility,
		Organizer: &calendar.EventOrganizer{
			Email:       DefaultCalendarEmail,
			DisplayName: DefaultCalendarDisplayName,
		},
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{
					Method:  "email",
					Minutes: DefaultCalendarEmailReminderMinutes,
				},
				{
					Method:  "popup",
					Minutes: DefaultCalendarPopupReminderMinutes,
				},
			},
			ForceSendFields: []string{"UseDefault"},
		},
	}

	if attendeeEmails != nil {
		attendees := []*calendar.EventAttendee{}
		for _, email := range attendeeEmails {
			attendees = append(
				attendees,
				&calendar.EventAttendee{
					Email: email,
				},
			)
		}
		event.Attendees = attendees
	}

	if privateExtendedPropertiesContent != nil {
		extendedProperties := &calendar.EventExtendedProperties{
			Private: privateExtendedPropertiesContent,
		}
		event.ExtendedProperties = extendedProperties
	}

	calInsertCall := s.gcalService.Events.Insert(cal.Id, event)
	calInsertCall.SendNotifications(true)
	event, err = calInsertCall.Do()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't create calendar event: %w", err)
	}
	return event, nil
}
