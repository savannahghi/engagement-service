package scheduling

import "time"

// BusySlot is a representation of an occupied start and end time
type BusySlot struct {
	Start time.Time
	End   time.Time
}
