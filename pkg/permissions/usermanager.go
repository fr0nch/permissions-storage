package permissions

/*
#include "usermanager.h"
#cgo noescape DumpPermissions
#cgo noescape CanAffectUser
#cgo noescape SetImmunity
#cgo noescape HasPermission
#cgo noescape HasGroup
#cgo noescape GetUserGroups
#cgo noescape GetImmunity
#cgo noescape AddPermission
#cgo noescape RemovePermission
#cgo noescape AddGroup
#cgo noescape RemoveGroup
#cgo noescape GetCookie
#cgo noescape SetCookie
#cgo noescape GetAllCookies
#cgo noescape CreateUser
#cgo noescape DeleteUser
#cgo noescape LoadUser
#cgo noescape LoadedUser
#cgo noescape UserExists
*/
import "C"
import (
	"errors"
	"github.com/untrustedmodders/go-plugify"
	"reflect"
	"runtime"
	"unsafe"
)

var _ = errors.New("")
var _ = reflect.TypeOf(0)
var _ = runtime.GOOS
var _ = unsafe.Sizeof(0)
var _ = plugify.Plugin.Loaded

// Generated from permissions (group: usermanager)

// DumpPermissions
//
//	@brief Get permissions of user
//
//	@param targetID: Player ID
//	@param perms: Permissions
//
//	@return Success, TargetUserNotFound
func DumpPermissions(targetID uint64, perms *[]string) Status {
	plugify.Log("permissions::DumpPermissions", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__perms := plugify.ConstructVectorString(*perms)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.DumpPermissions(__targetID, (*C.Vector)(unsafe.Pointer(&__perms))))
			// Unmarshal - Convert native data to managed data.
			plugify.GetVectorDataStringTo(&__perms, perms)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyVectorString(&__perms)
		},
	}.Do()
	return __retVal
}

// CanAffectUser
//
//	@brief Check players immunity or groups priority
//
//	@param actorID: Player performing the action
//	@param targetID: Player receiving the action
//
//	@return Allow, Disallow, ActorUserNotFound, or TargetUserNotFound
func CanAffectUser(actorID uint64, targetID uint64) Status {
	plugify.Log("permissions::CanAffectUser", plugify.Trace, 2)
	var __retVal Status
	__actorID := C.uint64_t(actorID)
	__targetID := C.uint64_t(targetID)
	__retVal = int32(C.CanAffectUser(__actorID, __targetID))
	return __retVal
}

// SetImmunity
//
//	@brief Set the immunity level of a user.
//
//	@param targetID: Player ID.
//	@param immunity: Immunity
//
//	@return Success, TargetUserNotFound
func SetImmunity(targetID uint64, immunity int32) Status {
	plugify.Log("permissions::SetImmunity", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__immunity := C.int32_t(immunity)
	__retVal = int32(C.SetImmunity(__targetID, __immunity))
	return __retVal
}

// HasPermission
//
//	@brief Check if a user has a specific permission.
//
//	@param targetID: Player ID.
//	@param perm: Permission line.
//
//	@return Allow, Disallow, PermNotFound, TargetUserNotFound
func HasPermission(targetID uint64, perm string) Status {
	plugify.Log("permissions::HasPermission", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__perm := plugify.ConstructString(perm)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.HasPermission(__targetID, (*C.String)(unsafe.Pointer(&__perm))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// HasGroup
//
//	@brief Check if a user belongs to a specific group (directly or via parent groups).
//
//	@param targetID: Player ID.
//	@param groupName: Group name.
//
//	@return PermanentGroup, TemporalGroup, GroupNotDefined, TargetUserNotFound, GroupNotFound
func HasGroup(targetID uint64, groupName string) Status {
	plugify.Log("permissions::HasGroup", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__groupName := plugify.ConstructString(groupName)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.HasGroup(__targetID, (*C.String)(unsafe.Pointer(&__groupName))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
		},
	}.Do()
	return __retVal
}

// GetUserGroups
//
//	@brief Get user groups.
//
//	@param targetID: Player ID.
//	@param outGroups: Groups
//
//	@return Success, TargetUserNotFound
func GetUserGroups(targetID uint64, outGroups *[]string) Status {
	plugify.Log("permissions::GetUserGroups", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__outGroups := plugify.ConstructVectorString(*outGroups)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetUserGroups(__targetID, (*C.Vector)(unsafe.Pointer(&__outGroups))))
			// Unmarshal - Convert native data to managed data.
			plugify.GetVectorDataStringTo(&__outGroups, outGroups)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyVectorString(&__outGroups)
		},
	}.Do()
	return __retVal
}

// GetImmunity
//
//	@brief Get the immunity level of a user.
//
//	@param targetID: Player ID.
//	@param immunity: Immunity
//
//	@return Success, TargetUserNotFound
func GetImmunity(targetID uint64, immunity *int32) Status {
	plugify.Log("permissions::GetImmunity", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__immunity := C.int32_t(*immunity)
	__retVal = int32(C.GetImmunity(__targetID, &__immunity))
	// Unmarshal - Convert native data to managed data.
	*immunity = int32(__immunity)
	return __retVal
}

// AddPermission
//
//	@brief Add a permission to a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param perm: Permission line.
//	@param timestamp: Permission duration.
//	@param dontBroadcast: If set to `true`, suppresses dispatching of the permission change event to registered UserPermission listeners. The permission is still applied internally.
//
//	@return Success, TargetUserNotFound, PermAlreadyGranted
func AddPermission(pluginID uint64, targetID uint64, perm string, timestamp int64, dontBroadcast bool) Status {
	plugify.Log("permissions::AddPermission", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__perm := plugify.ConstructString(perm)
	__timestamp := C.int64_t(timestamp)
	__dontBroadcast := C.bool(dontBroadcast)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.AddPermission(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__perm)), __timestamp, __dontBroadcast))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// RemovePermission
//
//	@brief Remove a permission from a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param perm: Permission line.
//
//	@return Success, TargetUserNotFound
func RemovePermission(pluginID uint64, targetID uint64, perm string) Status {
	plugify.Log("permissions::RemovePermission", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__perm := plugify.ConstructString(perm)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.RemovePermission(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__perm))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// AddGroup
//
//	@brief Add a group to a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param groupName: Group name.
//	@param timestamp: Group duration.
//	@param dontBroadcast: If set to `true`, suppresses dispatching of the group change event to registered UserGroup listeners. The group is still applied internally.
//
//	@return Success, TargetUserNotFound, GroupNotFound, GroupAlreadyExist
func AddGroup(pluginID uint64, targetID uint64, groupName string, timestamp int64, dontBroadcast bool) Status {
	plugify.Log("permissions::AddGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__groupName := plugify.ConstructString(groupName)
	__timestamp := C.int64_t(timestamp)
	__dontBroadcast := C.bool(dontBroadcast)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.AddGroup(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__groupName)), __timestamp, __dontBroadcast))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
		},
	}.Do()
	return __retVal
}

// RemoveGroup
//
//	@brief Remove a group from a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param groupName: Group name.
//
//	@return Success, TargetUserNotFound, ChildGroupNotFound, ParentGroupNotFound
func RemoveGroup(pluginID uint64, targetID uint64, groupName string) Status {
	plugify.Log("permissions::RemoveGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__groupName := plugify.ConstructString(groupName)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.RemoveGroup(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__groupName))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
		},
	}.Do()
	return __retVal
}

// GetCookie
//
//	@brief Get a cookie value for a user.
//
//	@param targetID: Player ID.
//	@param name: Cookie name.
//	@param value: Cookie value.
//
//	@return Success, TargetUserNotFound, CookieNotFound
func GetCookie(targetID uint64, name string, value *any) Status {
	plugify.Log("permissions::GetCookie", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__name := plugify.ConstructString(name)
	__value := plugify.ConstructVariant(*value)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetCookie(__targetID, (*C.String)(unsafe.Pointer(&__name)), (*C.Variant)(unsafe.Pointer(&__value))))
			// Unmarshal - Convert native data to managed data.
			*value = plugify.GetVariantData(&__value)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyVariant(&__value)
		},
	}.Do()
	return __retVal
}

// SetCookie
//
//	@brief Set a cookie value for a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param name: Cookie name.
//	@param cookie: Cookie value.
//	@param dontBroadcast: If set to `true`, suppresses dispatching of the cookie change event to registered UserSetCookie listeners. The cookie is still applied internally.
//
//	@return Success, TargetUserNotFound
func SetCookie(pluginID uint64, targetID uint64, name string, cookie any, dontBroadcast bool) Status {
	plugify.Log("permissions::SetCookie", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__name := plugify.ConstructString(name)
	__cookie := plugify.ConstructVariant(cookie)
	__dontBroadcast := C.bool(dontBroadcast)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.SetCookie(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__name)), (*C.Variant)(unsafe.Pointer(&__cookie)), __dontBroadcast))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyVariant(&__cookie)
		},
	}.Do()
	return __retVal
}

// GetAllCookies
//
//	@brief Get all cookies from user.
//
//	@param targetID: Player ID.
//	@param names: Array of cookie names
//	@param values: Array of cookie values
//
//	@return Success, TargetUserNotFound
func GetAllCookies(targetID uint64, names *[]string, values *[]any) Status {
	plugify.Log("permissions::GetAllCookies", plugify.Trace, 2)
	var __retVal Status
	__targetID := C.uint64_t(targetID)
	__names := plugify.ConstructVectorString(*names)
	__values := plugify.ConstructVectorVariant(*values)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetAllCookies(__targetID, (*C.Vector)(unsafe.Pointer(&__names)), (*C.Vector)(unsafe.Pointer(&__values))))
			// Unmarshal - Convert native data to managed data.
			plugify.GetVectorDataStringTo(&__names, names)
			plugify.GetVectorDataVariantTo(&__values, values)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyVectorString(&__names)
			plugify.DestroyVectorVariant(&__values)
		},
	}.Do()
	return __retVal
}

// CreateUser
//
//	@brief Create a new user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//	@param immunity: User immunity (set -1 to return highest group priority).
//	@param groupsList: Array of groups to inherit ("group timestamp").
//	@param permsList: Array of permissions (perm.iss.ion timestamp) or (perm.iss.ion).
//
//	@return Success, UserAlreadyExist, GroupNotFound, ChildGroupNotFound
func CreateUser(pluginID uint64, targetID uint64, immunity int32, groupsList []string, permsList []string) Status {
	plugify.Log("permissions::CreateUser", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__immunity := C.int32_t(immunity)
	__groupsList := plugify.ConstructVectorString(groupsList)
	__permsList := plugify.ConstructVectorString(permsList)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.CreateUser(__pluginID, __targetID, __immunity, (*C.Vector)(unsafe.Pointer(&__groupsList)), (*C.Vector)(unsafe.Pointer(&__permsList))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyVectorString(&__groupsList)
			plugify.DestroyVectorString(&__permsList)
		},
	}.Do()
	return __retVal
}

// DeleteUser
//
//	@brief Delete a user.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: Player ID.
//
//	@return Success, TargetUserNotFound
func DeleteUser(pluginID uint64, targetID uint64) Status {
	plugify.Log("permissions::DeleteUser", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__retVal = int32(C.DeleteUser(__pluginID, __targetID))
	return __retVal
}

// LoadUser
//
//	@brief Dispatches a request to load user data. Notifies all registered listeners that the specified user's data should be loaded from an external storage provider. This function does not perform any storage operations by itself. It only broadcasts the load request event.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param targetID: PlayerID of the user whose data should be loaded.
//	@param username: The user's current username. Intended for synchronizing the username with external storage (e.g. updating an existing record or setting it during initial user creation).
func LoadUser(pluginID uint64, targetID uint64, username string) {
	plugify.Log("permissions::LoadUser", plugify.Trace, 2)
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	__username := plugify.ConstructString(username)
	plugify.Block{
		Try: func() {
			C.LoadUser(__pluginID, __targetID, (*C.String)(unsafe.Pointer(&__username)))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__username)
		},
	}.Do()
}

// LoadedUser
//
//	@brief Dispatches a user-loaded event. Invoked by a storage provider to indicate that the requested user data loading process has completed successfully. After this call, the user is considered fully initialized.
//
//	@param pluginID: Identifier of the storage plugin reporting completion.
//	@param targetID: PlayerID of the loaded user.
func LoadedUser(pluginID uint64, targetID uint64) {
	plugify.Log("permissions::LoadedUser", plugify.Trace, 2)
	__pluginID := C.uint64_t(pluginID)
	__targetID := C.uint64_t(targetID)
	C.LoadedUser(__pluginID, __targetID)
}

// UserExists
//
//	@brief Check if a user exists.
//
//	@param targetID: Player ID.
//
//	@return True if user exists, false otherwise.
func UserExists(targetID uint64) bool {
	plugify.Log("permissions::UserExists", plugify.Trace, 2)
	var __retVal bool
	__targetID := C.uint64_t(targetID)
	__retVal = bool(C.UserExists(__targetID))
	return __retVal
}
