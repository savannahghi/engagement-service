package scheduling

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
)

// RoundTripFunc is the signature of a round trip function that can be
// assigned to a HTTP Client Transport
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is a mock round trip implementation that is used for testing.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// MockGCALHTTPClient is a mock HTTP client that always returns the same
// valid (mock) Google Calendar API CreateCalendar response
func MockGCALHTTPClient() *http.Client {
	mockCalendar := &calendar.Calendar{
		Description: "Mock calendar for testing",
		Id:          uuid.New().String(),
		Etag:        uuid.New().String(),
		Kind:        "calendar#calendar",
		Location:    "Mawinguni",
		Summary:     "Mock calendar for testing",
		TimeZone:    "Africa/Nairobi",
	}
	mockCalendarJSON, err := json.Marshal(mockCalendar)
	if err != nil {
		log.Panicf("unable to marshal mock calendar JSON: %v", err)
	}
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(mockCalendarJSON)),

			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	return client
}
