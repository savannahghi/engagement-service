package feed

import (
	"fmt"
	"io"
	"strconv"
)

// ActionType defines the types for global actions
type ActionType string

// the known action types are constants
const (
	ActionTypePrimary   ActionType = "PRIMARY"
	ActionTypeSecondary ActionType = "SECONDARY"
	ActionTypeOverflow  ActionType = "OVERFLOW"
	ActionTypeFloating  ActionType = "FLOATING"
)

// AllActionType has the known set of action types
var AllActionType = []ActionType{
	ActionTypePrimary,
	ActionTypeSecondary,
	ActionTypeOverflow,
	ActionTypeFloating,
}

// IsValid returns true only for valid action types
func (e ActionType) IsValid() bool {
	switch e {
	case ActionTypePrimary,
		ActionTypeSecondary,
		ActionTypeOverflow,
		ActionTypeFloating:
		return true
	}
	return false
}

func (e ActionType) String() string {
	return string(e)
}

// UnmarshalGQL reads an action type from GQL
func (e *ActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionType", str)
	}
	return nil
}

// MarshalGQL writes an action type to the supplied writer
func (e ActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Handling determines whether an action is handled INLINE or
type Handling string

// known action handling strategies
const (
	HandlingInline   Handling = "INLINE"
	HandlingFullPage Handling = "FULL_PAGE"
)

// AllHandling is the set of all valid handling strategies
var AllHandling = []Handling{
	HandlingInline,
	HandlingFullPage,
}

// IsValid returns true only for valid handling strategies
func (e Handling) IsValid() bool {
	switch e {
	case HandlingInline, HandlingFullPage:
		return true
	}
	return false
}

func (e Handling) String() string {
	return string(e)
}

// UnmarshalGQL reads and validates a handling value from the supplied input
func (e *Handling) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Handling(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Handling", str)
	}
	return nil
}

// MarshalGQL writes the Handling value to the supplied writer
func (e Handling) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Status is the set of known statuses for feed items and nudges
type Status string

// known item and nudge statuses
const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

// AllStatus is the set of known statuses
var AllStatus = []Status{
	StatusPending,
	StatusInProgress,
	StatusDone,
}

// IsValid returns true if a status is valid
func (e Status) IsValid() bool {
	switch e {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

// UnmarshalGQL translates the input value given into a status
func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

// MarshalGQL writes the status to the supplied writer
func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Visibility defines the visibility statuses of feed items
type Visibility string

// known visibility values
const (
	VisibilityShow Visibility = "SHOW"
	VisibilityHide Visibility = "HIDE"
)

// AllVisibility is the set of all known visibility values
var AllVisibility = []Visibility{
	VisibilityShow,
	VisibilityHide,
}

// IsValid returns true if a visibility value is valid
func (e Visibility) IsValid() bool {
	switch e {
	case VisibilityShow, VisibilityHide:
		return true
	}
	return false
}

func (e Visibility) String() string {
	return string(e)
}

// UnmarshalGQL reads and validates a visibility value
func (e *Visibility) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Visibility(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Visibility", str)
	}
	return nil
}

// MarshalGQL writes a visibility value into the supplied writer
func (e Visibility) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Channel represents a notification challen
type Channel string

// known notification channels
const (
	ChannelFcm      Channel = "FCM"
	ChannelEmail    Channel = "EMAIL"
	ChannelSms      Channel = "SMS"
	ChannelWhatsapp Channel = "WHATSAPP"
)

// AllChannel is the set of all supported notification channels
var AllChannel = []Channel{
	ChannelFcm,
	ChannelEmail,
	ChannelSms,
	ChannelWhatsapp,
}

// IsValid returns True only for a valid channel
func (e Channel) IsValid() bool {
	switch e {
	case ChannelFcm, ChannelEmail, ChannelSms, ChannelWhatsapp:
		return true
	}
	return false
}

func (e Channel) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied input into a channel value
func (e *Channel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Channel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Channel", str)
	}
	return nil
}

// MarshalGQL writes the channel to the supplied writer
func (e Channel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Flavour is the flavour of a feed i.e consumer or pro
type Flavour string

// known flavours
const (
	FlavourPro      Flavour = "PRO"
	FlavourConsumer Flavour = "CONSUMER"
)

// AllFlavour is a set of all valid flavours
var AllFlavour = []Flavour{
	FlavourPro,
	FlavourConsumer,
}

// IsValid returns True if a feed is valid
func (e Flavour) IsValid() bool {
	switch e {
	case FlavourPro, FlavourConsumer:
		return true
	}
	return false
}

func (e Flavour) String() string {
	return string(e)
}

// UnmarshalGQL translates and validates the input flavour
func (e *Flavour) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Flavour(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Flavour", str)
	}
	return nil
}

// MarshalGQL writes the flavour to the supplied writer
func (e Flavour) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Keys are the top level keys in a feed
type Keys string

// known feed keys
const (
	KeysActions Keys = "actions"
	KeysNudges  Keys = "nudges"
	KeysItems   Keys = "items"
)

// AllKeys is the set of all valid feed keys
var AllKeys = []Keys{
	KeysActions,
	KeysNudges,
	KeysItems,
}

// IsValid returns true if a feed key is valid
func (e Keys) IsValid() bool {
	switch e {
	case KeysActions, KeysNudges, KeysItems:
		return true
	}
	return false
}

func (e Keys) String() string {
	return string(e)
}

// UnmarshalGQL translates a feed key from a string
func (e *Keys) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Keys(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FeedKeys", str)
	}
	return nil
}

// MarshalGQL writes the feed key to the supplied writer
func (e Keys) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// BooleanFilter defines true/false/both for filtering against bools
type BooleanFilter string

// known boolean filter value
const (
	BooleanFilterTrue  BooleanFilter = "TRUE"
	BooleanFilterFalse BooleanFilter = "FALSE"
	BooleanFilterBoth  BooleanFilter = "BOTH"
)

// IsValid is a set of known boolean filters
var IsValid = []BooleanFilter{
	BooleanFilterTrue,
	BooleanFilterFalse,
	BooleanFilterBoth,
}

// IsValid returns True if the boolean filter value is valid
func (e BooleanFilter) IsValid() bool {
	switch e {
	case BooleanFilterTrue, BooleanFilterFalse, BooleanFilterBoth:
		return true
	}
	return false
}

func (e BooleanFilter) String() string {
	return string(e)
}

// UnmarshalGQL reads the bool value in from input
func (e *BooleanFilter) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BooleanFilter(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BooleanFilter", str)
	}
	return nil
}

// MarshalGQL writes the bool value to the supplied writer
func (e BooleanFilter) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
