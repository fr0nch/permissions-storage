package permissions

/*
#include "userlisteners.h"
#cgo noescape OnLoadUser_Register
#cgo noescape OnLoadUser_Unregister
#cgo noescape OnLoadedUser_Register
#cgo noescape OnLoadedUser_Unregister
#cgo noescape OnUserPermissionChange_Register
#cgo noescape OnUserPermissionChange_Unregister
#cgo noescape OnUserSetCookie_Register
#cgo noescape OnUserSetCookie_Unregister
#cgo noescape OnUserGroupChange_Register
#cgo noescape OnUserGroupChange_Unregister
#cgo noescape OnUserCreate_Register
#cgo noescape OnUserCreate_Unregister
#cgo noescape OnUserDelete_Register
#cgo noescape OnUserDelete_Unregister
#cgo noescape OnPermissionExpirationCallback_Register
#cgo noescape OnPermissionExpirationCallback_Unregister
#cgo noescape OnGroupExpirationCallback_Register
#cgo noescape OnGroupExpirationCallback_Unregister
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

// Generated from permissions (group: userlisteners)

// OnLoadUser_Register
//
//	@brief Register listener on LoadUser event.
//
//	@param callback: Function callback.
func OnLoadUser_Register(callback UserLoadCallback) Status {
	plugify.Log("permissions::OnLoadUser_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadUser_Register(__callback))
	return __retVal
}

// OnLoadUser_Unregister
//
//	@brief Unregister listener on LoadUser event.
//
//	@param callback: Function callback.
func OnLoadUser_Unregister(callback UserLoadCallback) Status {
	plugify.Log("permissions::OnLoadUser_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadUser_Unregister(__callback))
	return __retVal
}

// OnLoadedUser_Register
//
//	@brief Register listener on LoadedUser event.
//
//	@param callback: Function callback.
func OnLoadedUser_Register(callback UserLoadedCallback) Status {
	plugify.Log("permissions::OnLoadedUser_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadedUser_Register(__callback))
	return __retVal
}

// OnLoadedUser_Unregister
//
//	@brief Unregister listener on LoadedUser event.
//
//	@param callback: Function callback.
func OnLoadedUser_Unregister(callback UserLoadedCallback) Status {
	plugify.Log("permissions::OnLoadedUser_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadedUser_Unregister(__callback))
	return __retVal
}

// OnUserPermissionChange_Register
//
//	@brief Register listener on user permission add/remove
//
//	@param callback: Function callback.
func OnUserPermissionChange_Register(callback UserPermissionCallback) Status {
	plugify.Log("permissions::OnUserPermissionChange_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserPermissionChange_Register(__callback))
	return __retVal
}

// OnUserPermissionChange_Unregister
//
//	@brief Unregister listener on user permission add/remove
//
//	@param callback: Function callback.
func OnUserPermissionChange_Unregister(callback UserPermissionCallback) Status {
	plugify.Log("permissions::OnUserPermissionChange_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserPermissionChange_Unregister(__callback))
	return __retVal
}

// OnUserSetCookie_Register
//
//	@brief Register listener on user cookie sets
//
//	@param callback: Function callback.
func OnUserSetCookie_Register(callback UserSetCookieCallback) Status {
	plugify.Log("permissions::OnUserSetCookie_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserSetCookie_Register(__callback))
	return __retVal
}

// OnUserSetCookie_Unregister
//
//	@brief Register listener on user cookie sets
//
//	@param callback: Function callback.
func OnUserSetCookie_Unregister(callback UserSetCookieCallback) Status {
	plugify.Log("permissions::OnUserSetCookie_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserSetCookie_Unregister(__callback))
	return __retVal
}

// OnUserGroupChange_Register
//
//	@brief Register listener on user groups changing
//
//	@param callback: Function callback.
func OnUserGroupChange_Register(callback UserGroupCallback) Status {
	plugify.Log("permissions::OnUserGroupChange_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserGroupChange_Register(__callback))
	return __retVal
}

// OnUserGroupChange_Unregister
//
//	@brief Unregister listener on user groups changing
//
//	@param callback: Function callback.
func OnUserGroupChange_Unregister(callback UserGroupCallback) Status {
	plugify.Log("permissions::OnUserGroupChange_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserGroupChange_Unregister(__callback))
	return __retVal
}

// OnUserCreate_Register
//
//	@brief Register listener on user creation
//
//	@param callback: Function callback.
func OnUserCreate_Register(callback UserCreateCallback) Status {
	plugify.Log("permissions::OnUserCreate_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserCreate_Register(__callback))
	return __retVal
}

// OnUserCreate_Unregister
//
//	@brief Unregister listener on user creation
//
//	@param callback: Function callback.
func OnUserCreate_Unregister(callback UserCreateCallback) Status {
	plugify.Log("permissions::OnUserCreate_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserCreate_Unregister(__callback))
	return __retVal
}

// OnUserDelete_Register
//
//	@brief Register listener on user deletion
//
//	@param callback: Function callback.
func OnUserDelete_Register(callback UserDeleteCallback) Status {
	plugify.Log("permissions::OnUserDelete_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserDelete_Register(__callback))
	return __retVal
}

// OnUserDelete_Unregister
//
//	@brief Unregister listener on user deletion
//
//	@param callback: Function callback.
func OnUserDelete_Unregister(callback UserDeleteCallback) Status {
	plugify.Log("permissions::OnUserDelete_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnUserDelete_Unregister(__callback))
	return __retVal
}

// OnPermissionExpirationCallback_Register
//
//	@brief Register listener on user permission expiration
//
//	@param callback: Function callback.
func OnPermissionExpirationCallback_Register(callback PermExpirationCallback) Status {
	plugify.Log("permissions::OnPermissionExpirationCallback_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnPermissionExpirationCallback_Register(__callback))
	return __retVal
}

// OnPermissionExpirationCallback_Unregister
//
//	@brief Unregister listener on user permission expiration
//
//	@param callback: Function callback.
func OnPermissionExpirationCallback_Unregister(callback PermExpirationCallback) Status {
	plugify.Log("permissions::OnPermissionExpirationCallback_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnPermissionExpirationCallback_Unregister(__callback))
	return __retVal
}

// OnGroupExpirationCallback_Register
//
//	@brief Register listener on user permission expiration
//
//	@param callback: Function callback.
func OnGroupExpirationCallback_Register(callback GroupExpirationCallback) Status {
	plugify.Log("permissions::OnGroupExpirationCallback_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupExpirationCallback_Register(__callback))
	return __retVal
}

// OnGroupExpirationCallback_Unregister
//
//	@brief Unregister listener on user group expiration
//
//	@param callback: Function callback.
func OnGroupExpirationCallback_Unregister(callback GroupExpirationCallback) Status {
	plugify.Log("permissions::OnGroupExpirationCallback_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupExpirationCallback_Unregister(__callback))
	return __retVal
}
