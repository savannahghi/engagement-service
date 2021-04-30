package common

// feed related constants
const (
	// these topic names are exported because they are also used by sub-packages
	ItemPublishTopic    = "items.publish"
	ItemDeleteTopic     = "items.delete"
	ItemResolveTopic    = "items.resolve"
	ItemUnresolveTopic  = "items.unresolve"
	ItemHideTopic       = "items.hide"
	ItemShowTopic       = "items.show"
	ItemPinTopic        = "items.pin"
	ItemUnpinTopic      = "items.unpin"
	NudgePublishTopic   = "nudges.publish"
	NudgeDeleteTopic    = "nudges.delete"
	NudgeResolveTopic   = "nudges.resolve"
	NudgeUnresolveTopic = "nudges.unresolve"
	NudgeHideTopic      = "nudges.hide"
	NudgeShowTopic      = "nudges.show"
	ActionPublishTopic  = "actions.publish"
	ActionDeleteTopic   = "actions.delete"
	MessagePostTopic    = "message.post"
	MessageDeleteTopic  = "message.delete"
	IncomingEventTopic  = "incoming.event"
	FcmPublishTopic     = "fcm.send_notification"
	SentEmailTopic      = "mails.inbox"

	// DefaultLabel is the label used for welcome content
	DefaultLabel = "WELCOME"

	// StaticBase is the default path at which static assets are hosted
	StaticBase = "https://assets.healthcloud.co.ke"

	// DefaultIconPath is the path to the default Be.Well logo
	DefaultIconPath = StaticBase + "/bewell_logo.png"

	// Default nudges `titles` have been defined as constants as
	// these same titles are used to `resolve` these default nudges
	// once a user completes the action that shall be defined by the nudges

	// AddPrimaryEmailNudgeTitle defines the title for verify email nudge
	// that adds and verifies a primary email address
	AddPrimaryEmailNudgeTitle = "Add Primary Email Address"

	// AddInsuranceNudgeTitle defines the title for add insurance nudge
	AddInsuranceNudgeTitle = "Add Insurance"

	// AddNHIFNudgeTitle defines the title for add NHIF nudge
	AddNHIFNudgeTitle = "Add NHIF"

	// PartnerAccountSetupNudgeTitle defines the title for partner account setup nudge
	PartnerAccountSetupNudgeTitle = "Setup your partner account"

	// HideItemActionName defines the name for the hide action
	HideItemActionName = "HIDE_ITEM"

	// PinItemActionName defines the name for the pin action
	PinItemActionName = "PIN_ITEM"

	// ResolveItemActionName defines the name for the resolve action
	ResolveItemActionName = "RESOLVE_ITEM"

	// ShowItemActionName defines the name for the hide action
	ShowItemActionName = "SHOW_ITEM"

	// UnPinItemActionName defines the name for the pin action
	UnPinItemActionName = "UNPIN_ITEM"

	// UnResolveItemActionName defines the name for the resolve action
	UnResolveItemActionName = "UNRESOLVE_ITEM"
)
