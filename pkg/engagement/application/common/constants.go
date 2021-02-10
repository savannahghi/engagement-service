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

	// DefaultLabel is the label used for welcome content
	DefaultLabel = "WELCOME"

	// StaticBase is the default path at which static assets are hosted
	StaticBase = "https://assets.healthcloud.co.ke"

	// DefaultIconPath is the path to the default Be.Well logo
	DefaultIconPath = StaticBase + "/bewell_logo.png"

	// AddPrimaryEmailNudgeTitle defines the title for verify email nudge
	// that adds and verifies a primary email address
	AddPrimaryEmailNudgeTitle = "Add Primary Email Address"

	// AddInsuranceNudgeTitle defines the title for add insurance nudge
	AddInsuranceNudgeTitle = "Add Insurance"

	// AddNHIFNudgeTitle defines the title for add NHIF nudge
	AddNHIFNudgeTitle = "Add NHIF"

	// PartnerAccountSetupNudgeTitle defines the title for partner account setup nudge
	PartnerAccountSetupNudgeTitle = "Setup your partner account"
)
