package permission

import (
	"gitlab.slade360emr.com/go/base"
)

// FeedView describes view permissions on a feed
var FeedView = base.PermissionInput{
	Resource: "feed_view",
	Action:   "view",
}

// ThinFeedView describes view permissions on a thin feed
var ThinFeedView = base.PermissionInput{
	Resource: "thin_feed_view",
	Action:   "view",
}

// FeedItemView describes view permissions on a feed item
var FeedItemView = base.PermissionInput{
	Resource: "feed_item_view",
	Action:   "view",
}

// NudgeView describes view permissions on a nudge
var NudgeView = base.PermissionInput{
	Resource: "nudge_view",
	Action:   "view",
}

// ActionView describes view permissions on an action
var ActionView = base.PermissionInput{
	Resource: "action_view",
	Action:   "view",
}

// PublishItem describes create permissions on a feed item
var PublishItem = base.PermissionInput{
	Resource: "publish_item",
	Action:   "create",
}

// DeleteItem describes delete permissions on a feed item
var DeleteItem = base.PermissionInput{
	Resource: "delete_item",
	Action:   "delete",
}

// ResolveItem describes the resolve permissions on an item
var ResolveItem = base.PermissionInput{
	Resource: "resolve_item",
	Action:   "resolve",
}

// UnresolveItem describes the unresolve permissions on an item
var UnresolveItem = base.PermissionInput{
	Resource: "unresolve_item",
	Action:   "unresolve",
}

// PinItem describes the pin permissions on an item
var PinItem = base.PermissionInput{
	Resource: "pin_item",
	Action:   "pin",
}

// UnpinItem describes the unpin permissions on an item. To mark a feed item as not persistent
var UnpinItem = base.PermissionInput{
	Resource: "unpin_item",
	Action:   "unpin",
}

// HideItem decribes the hide permissions on an item
var HideItem = base.PermissionInput{
	Resource: "hide_item",
	Action:   "hide",
}

// ShowItem describes the view permissions on an item
var ShowItem = base.PermissionInput{
	Resource: "show_item",
	Action:   "view",
}

// GetLabel describes the view permissions on a label
var GetLabel = base.PermissionInput{
	Resource: "get_label",
	Action:   "view",
}

// CreateLabel describes the create permissions on a label
var CreateLabel = base.PermissionInput{
	Resource: "create_label",
	Action:   "create",
}

// UnreadPersistentItems describes the permissions on a feed item
var UnreadPersistentItems = base.PermissionInput{
	Resource: "unread_persistent_item",
	Action:   "update",
}

// UpdateUnreadPersistentItems describes the permissions on a feed item
var UpdateUnreadPersistentItems = base.PermissionInput{
	Resource: "update_unread_persistent_item",
	Action:   "update",
}

// PostMessage describes the create permissions on a message
var PostMessage = base.PermissionInput{
	Resource: "post_message",
	Action:   "create",
}

// DeleteMessage describes the delete permissions on a message
var DeleteMessage = base.PermissionInput{
	Resource: "delete_message",
	Action:   "delete",
}

// ProcessEvent describes the create permission on processing events
var ProcessEvent = base.PermissionInput{
	Resource: "process_event",
	Action:   "create",
}

// ItemUpdate describes the update permissions on items
var ItemUpdate = base.PermissionInput{
	Resource: "item_update",
	Action:   "update",
}

// SendMessage describes the create permissions on a message
var SendMessage = base.PermissionInput{
	Resource: "send_message",
	Action:   "create",
}

// LoadMarketingData describes the create permissions on a message
var LoadMarketingData = base.PermissionInput{
	Resource: "load_marketing_data",
	Action:   "create",
}
