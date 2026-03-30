package permissions

import "github.com/untrustedmodders/go-plugify"

var _ = plugify.Plugin.Loaded

// Generated from permissions

// LoadGroupsCallback - Called when the core requests loading of server groups. This callback is triggered when the system needs to load group definitions associated with a specific plugin. Extensions (e.g., database providers) should subscribe to this event and load the groups into memory.
type LoadGroupsCallback func(pluginID uint64)

// SetParentCallback - Callback invoked when a parent group is set for a child group.
type SetParentCallback func(pluginID uint64, childName string, parentName string)

// SetCookieGroupCallback - Callback invoked when a cookie value is set for a group.
type SetCookieGroupCallback func(pluginID uint64, groupName string, cookieName string, value any)

// UserLoadCallback - Called when a user data load is requested. This callback is triggered by the core when it requires user data to be loaded from an external storage (e.g. database). Extensions can subscribe to this event to perform the actual loading process and initialize the user in memory. This event does NOT guarantee that the user object already exists in memory.
type UserLoadCallback func(pluginID uint64, targetID uint64, username string)

// UserLoadedCallback - Called when a user's data has been fully loaded. This callback is triggered after a storage extension has completed loading and applying the user's persistent data (e.g. groups, permissions, metadata). At this stage, the user is considered fully initialized and ready for normal operation within the system.
type UserLoadedCallback func(pluginID uint64, targetID uint64)

// GroupPermissionCallback - Callback invoked when a permission is added or removed from a group.
type GroupPermissionCallback func(pluginID uint64, action Action, name string, groupName string)

// UserPermissionCallback - Callback invoked when a permission is added or removed for a user.
type UserPermissionCallback func(pluginID uint64, action Action, targetID uint64, perm string, timestamp int64)

// GroupCreateCallback - Callback invoked after a group is successfully created.
type GroupCreateCallback func(pluginID uint64, name string, perms []string, priority int32, parent string)

// GroupDeleteCallback - Callback invoked before a group is deleted.
type GroupDeleteCallback func(pluginID uint64, name string)

// UserSetCookieCallback - Callback invoked when a cookie is set for a user.
type UserSetCookieCallback func(pluginID uint64, targetID uint64, name string, cookie any)

// UserGroupCallback - Callback invoked when a group is added or removed from a user.
type UserGroupCallback func(pluginID uint64, action Action, targetID uint64, group string, timestamp int64)

// UserCreateCallback - Callback invoked after a user is successfully created.
type UserCreateCallback func(pluginID uint64, targetID uint64, immunity int32, groupNames []string, perms []string)

// UserDeleteCallback - Callback invoked before a user is deleted.
type UserDeleteCallback func(pluginID uint64, targetID uint64)

// PermExpirationCallback - Callback invoked when a permission in user has been expired.
type PermExpirationCallback func(targetID uint64, perm string)

// GroupExpirationCallback - Callback invoked when a group in user has been expired.
type GroupExpirationCallback func(targetID uint64, group string)
