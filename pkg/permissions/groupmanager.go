package permissions

/*
#include "groupmanager.h"
#cgo noescape SetParent
#cgo noescape GetParent
#cgo noescape DumpPermissionsGroup
#cgo noescape GetAllGroups
#cgo noescape HasPermissionGroup
#cgo noescape HasParentGroup
#cgo noescape GetPriorityGroup
#cgo noescape AddPermissionGroup
#cgo noescape RemovePermissionGroup
#cgo noescape GetCookieGroup
#cgo noescape SetCookieGroup
#cgo noescape GetAllCookiesGroup
#cgo noescape CreateGroup
#cgo noescape LoadGroups
#cgo noescape DeleteGroup
#cgo noescape GroupExists
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

// Generated from permissions (group: groupmanager)

// SetParent
//
//	@brief Set parent group for child group
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param childName: Child group name
//	@param parentName: Parent group name to set
//
//	@return Success, ChildGroupNotFound, ParentGroupNotFound
func SetParent(pluginID uint64, childName string, parentName string) Status {
	plugify.Log("permissions::SetParent", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__childName := plugify.ConstructString(childName)
	__parentName := plugify.ConstructString(parentName)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.SetParent(__pluginID, (*C.String)(unsafe.Pointer(&__childName)), (*C.String)(unsafe.Pointer(&__parentName))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__childName)
			plugify.DestroyString(&__parentName)
		},
	}.Do()
	return __retVal
}

// GetParent
//
//	@brief Get parent of requested group
//
//	@param groupName: Group name
//	@param parentName: Parent name
//
//	@return Success, ChildGroupNotFound, ParentGroupNotFound
func GetParent(groupName string, parentName *string) Status {
	plugify.Log("permissions::GetParent", plugify.Trace, 2)
	var __retVal Status
	__groupName := plugify.ConstructString(groupName)
	__parentName := plugify.ConstructString(*parentName)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetParent((*C.String)(unsafe.Pointer(&__groupName)), (*C.String)(unsafe.Pointer(&__parentName))))
			// Unmarshal - Convert native data to managed data.
			*parentName = plugify.GetStringData(&__parentName)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
			plugify.DestroyString(&__parentName)
		},
	}.Do()
	return __retVal
}

// DumpPermissionsGroup
//
//	@brief Get permissions of group
//
//	@param name: Group name
//	@param perms: Permissions
//
//	@return Success, GroupNotFound
func DumpPermissionsGroup(name string, perms *[]string) Status {
	plugify.Log("permissions::DumpPermissionsGroup", plugify.Trace, 2)
	var __retVal Status
	__name := plugify.ConstructString(name)
	__perms := plugify.ConstructVectorString(*perms)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.DumpPermissionsGroup((*C.String)(unsafe.Pointer(&__name)), (*C.Vector)(unsafe.Pointer(&__perms))))
			// Unmarshal - Convert native data to managed data.
			plugify.GetVectorDataStringTo(&__perms, perms)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyVectorString(&__perms)
		},
	}.Do()
	return __retVal
}

// GetAllGroups
//
//	@brief Get all created groups
//
//
//	@return Array of groups
func GetAllGroups() []string {
	plugify.Log("permissions::GetAllGroups", plugify.Trace, 2)
	var __retVal []string
	var __retVal_native plugify.PlgVector
	plugify.Block{
		Try: func() {
			__native := C.GetAllGroups()
			__retVal_native = *(*plugify.PlgVector)(unsafe.Pointer(&__native))
			// Unmarshal - Convert native data to managed data.
			__retVal = plugify.GetVectorDataString(&__retVal_native)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyVectorString(&__retVal_native)
		},
	}.Do()
	return __retVal
}

// HasPermissionGroup
//
//	@brief Check if a group has a specific permission.
//
//	@param name: Group name.
//	@param perm: Permission line.
//
//	@return Allow, Disallow, PermNotFound, GroupNotFound
func HasPermissionGroup(name string, perm string) Status {
	plugify.Log("permissions::HasPermissionGroup", plugify.Trace, 2)
	var __retVal Status
	__name := plugify.ConstructString(name)
	__perm := plugify.ConstructString(perm)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.HasPermissionGroup((*C.String)(unsafe.Pointer(&__name)), (*C.String)(unsafe.Pointer(&__perm))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// HasParentGroup
//
//	@brief Check if parent_name is a parent group for child_name.
//
//	@param childName: Child group name.
//	@param parentName: Parent group name to check.
//
//	@return Allow, Disallow, ChildGroupNotFound, ParentGroupNotFound
func HasParentGroup(childName string, parentName string) Status {
	plugify.Log("permissions::HasParentGroup", plugify.Trace, 2)
	var __retVal Status
	__childName := plugify.ConstructString(childName)
	__parentName := plugify.ConstructString(parentName)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.HasParentGroup((*C.String)(unsafe.Pointer(&__childName)), (*C.String)(unsafe.Pointer(&__parentName))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__childName)
			plugify.DestroyString(&__parentName)
		},
	}.Do()
	return __retVal
}

// GetPriorityGroup
//
//	@brief Get the priority of a group.
//
//	@param groupName: Group name.
//	@param priority: Priority
//
//	@return Success, GroupNotFound
func GetPriorityGroup(groupName string, priority *int32) Status {
	plugify.Log("permissions::GetPriorityGroup", plugify.Trace, 2)
	var __retVal Status
	__groupName := plugify.ConstructString(groupName)
	__priority := C.int32_t(*priority)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetPriorityGroup((*C.String)(unsafe.Pointer(&__groupName)), &__priority))
			// Unmarshal - Convert native data to managed data.
			*priority = int32(__priority)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
		},
	}.Do()
	return __retVal
}

// AddPermissionGroup
//
//	@brief Add a permission to a group.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param name: Group name.
//	@param perm: Permission line.
//
//	@return Success, GroupNotFound
func AddPermissionGroup(pluginID uint64, name string, perm string) Status {
	plugify.Log("permissions::AddPermissionGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__name := plugify.ConstructString(name)
	__perm := plugify.ConstructString(perm)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.AddPermissionGroup(__pluginID, (*C.String)(unsafe.Pointer(&__name)), (*C.String)(unsafe.Pointer(&__perm))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// RemovePermissionGroup
//
//	@brief Remove a permission from a group.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param name: Group name.
//	@param perm: Permission line.
//
//	@return Success, GroupNotFound
func RemovePermissionGroup(pluginID uint64, name string, perm string) Status {
	plugify.Log("permissions::RemovePermissionGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__name := plugify.ConstructString(name)
	__perm := plugify.ConstructString(perm)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.RemovePermissionGroup(__pluginID, (*C.String)(unsafe.Pointer(&__name)), (*C.String)(unsafe.Pointer(&__perm))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyString(&__perm)
		},
	}.Do()
	return __retVal
}

// GetCookieGroup
//
//	@brief Get a cookie value for a group.
//
//	@param groupName: Group name
//	@param cookieName: Cookie name
//	@param value: Cookie value
//
//	@return Success, CookieNotFound, GroupNotFound
func GetCookieGroup(groupName string, cookieName string, value *any) Status {
	plugify.Log("permissions::GetCookieGroup", plugify.Trace, 2)
	var __retVal Status
	__groupName := plugify.ConstructString(groupName)
	__cookieName := plugify.ConstructString(cookieName)
	__value := plugify.ConstructVariant(*value)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetCookieGroup((*C.String)(unsafe.Pointer(&__groupName)), (*C.String)(unsafe.Pointer(&__cookieName)), (*C.Variant)(unsafe.Pointer(&__value))))
			// Unmarshal - Convert native data to managed data.
			*value = plugify.GetVariantData(&__value)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
			plugify.DestroyString(&__cookieName)
			plugify.DestroyVariant(&__value)
		},
	}.Do()
	return __retVal
}

// SetCookieGroup
//
//	@brief Set a cookie value for a group.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param groupName: Group name
//	@param cookieName: Cookie name
//	@param value: Cookie value.
//
//	@return Success, GroupNotFound
func SetCookieGroup(pluginID uint64, groupName string, cookieName string, value any) Status {
	plugify.Log("permissions::SetCookieGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__groupName := plugify.ConstructString(groupName)
	__cookieName := plugify.ConstructString(cookieName)
	__value := plugify.ConstructVariant(value)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.SetCookieGroup(__pluginID, (*C.String)(unsafe.Pointer(&__groupName)), (*C.String)(unsafe.Pointer(&__cookieName)), (*C.Variant)(unsafe.Pointer(&__value))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
			plugify.DestroyString(&__cookieName)
			plugify.DestroyVariant(&__value)
		},
	}.Do()
	return __retVal
}

// GetAllCookiesGroup
//
//	@brief Get all cookies from group.
//
//	@param groupName: Group name
//	@param cookieNames: Array of cookie names
//	@param values: Array of cookie values
//
//	@return Success, GroupNotFound
func GetAllCookiesGroup(groupName string, cookieNames *[]string, values *[]any) Status {
	plugify.Log("permissions::GetAllCookiesGroup", plugify.Trace, 2)
	var __retVal Status
	__groupName := plugify.ConstructString(groupName)
	__cookieNames := plugify.ConstructVectorString(*cookieNames)
	__values := plugify.ConstructVectorVariant(*values)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.GetAllCookiesGroup((*C.String)(unsafe.Pointer(&__groupName)), (*C.Vector)(unsafe.Pointer(&__cookieNames)), (*C.Vector)(unsafe.Pointer(&__values))))
			// Unmarshal - Convert native data to managed data.
			plugify.GetVectorDataStringTo(&__cookieNames, cookieNames)
			plugify.GetVectorDataVariantTo(&__values, values)
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__groupName)
			plugify.DestroyVectorString(&__cookieNames)
			plugify.DestroyVectorVariant(&__values)
		},
	}.Do()
	return __retVal
}

// CreateGroup
//
//	@brief Create a new group.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param name: Group name.
//	@param perms: Array of permission lines.
//	@param priority: Group priority.
//	@param parent: Parent group name.
//
//	@return Success, GroupAlreadyExist, ParentGroupNotFound
func CreateGroup(pluginID uint64, name string, perms []string, priority int32, parent string) Status {
	plugify.Log("permissions::CreateGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__name := plugify.ConstructString(name)
	__perms := plugify.ConstructVectorString(perms)
	__priority := C.int32_t(priority)
	__parent := plugify.ConstructString(parent)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.CreateGroup(__pluginID, (*C.String)(unsafe.Pointer(&__name)), (*C.Vector)(unsafe.Pointer(&__perms)), __priority, (*C.String)(unsafe.Pointer(&__parent))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
			plugify.DestroyVectorString(&__perms)
			plugify.DestroyString(&__parent)
		},
	}.Do()
	return __retVal
}

// LoadGroups
//
//	@brief Dispatches a request to load server groups for a plugin.
//
//	@param pluginID: Identifier of the plugin that calls the method.
func LoadGroups(pluginID uint64) {
	plugify.Log("permissions::LoadGroups", plugify.Trace, 2)
	__pluginID := C.uint64_t(pluginID)
	C.LoadGroups(__pluginID)
}

// DeleteGroup
//
//	@brief Delete a group.
//
//	@param pluginID: Identifier of the plugin that calls the method.
//	@param name: Group name.
//
//	@return Success if deleted; GroupNotFound if group not found.
func DeleteGroup(pluginID uint64, name string) Status {
	plugify.Log("permissions::DeleteGroup", plugify.Trace, 2)
	var __retVal Status
	__pluginID := C.uint64_t(pluginID)
	__name := plugify.ConstructString(name)
	plugify.Block{
		Try: func() {
			__retVal = int32(C.DeleteGroup(__pluginID, (*C.String)(unsafe.Pointer(&__name))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
		},
	}.Do()
	return __retVal
}

// GroupExists
//
//	@brief Check if a group exists.
//
//	@param name: Group name.
//
//	@return True if group exists, false otherwise.
func GroupExists(name string) bool {
	plugify.Log("permissions::GroupExists", plugify.Trace, 2)
	var __retVal bool
	__name := plugify.ConstructString(name)
	plugify.Block{
		Try: func() {
			__retVal = bool(C.GroupExists((*C.String)(unsafe.Pointer(&__name))))
		},
		Finally: func() {
			// Perform cleanup.
			plugify.DestroyString(&__name)
		},
	}.Do()
	return __retVal
}
