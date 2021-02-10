package exceptions

import "fmt"

// ErrNilNudge is a marker error for cases when a nudge should have
// been found but was not
var ErrNilNudge = fmt.Errorf("nil nudge")

// ErrNilFeedItem is a sentinel error used to indicate when a feed item
// should have been non nil but was not
var ErrNilFeedItem = fmt.Errorf("nil feed item")
