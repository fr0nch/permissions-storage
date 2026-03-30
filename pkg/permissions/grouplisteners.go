package permissions

/*
#include "grouplisteners.h"
#cgo noescape OnLoadGroups_Register
#cgo noescape OnLoadGroups_Unregister
#cgo noescape OnGroupSetParent_Register
#cgo noescape OnGroupSetParent_Unregister
#cgo noescape OnGroupSetCookie_Register
#cgo noescape OnGroupSetCookie_Unregister
#cgo noescape OnGroupPermissionChange_Register
#cgo noescape OnGroupPermissionChange_Unregister
#cgo noescape OnGroupCreate_Register
#cgo noescape OnGroupCreate_Unregister
#cgo noescape OnGroupDelete_Register
#cgo noescape OnGroupDelete_Unregister
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

// Generated from permissions (group: grouplisteners)

// OnLoadGroups_Register
//
//	@brief Register listener on LoadGroups event.
//
//	@param callback: Function callback.
func OnLoadGroups_Register(callback LoadGroupsCallback) Status {
	plugify.Log("permissions::OnLoadGroups_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadGroups_Register(__callback))
	return __retVal
}

// OnLoadGroups_Unregister
//
//	@brief Unregister listener on LoadGroups event.
//
//	@param callback: Function callback.
func OnLoadGroups_Unregister(callback LoadGroupsCallback) Status {
	plugify.Log("permissions::OnLoadGroups_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnLoadGroups_Unregister(__callback))
	return __retVal
}

// OnGroupSetParent_Register
//
//	@brief Register listener on group parent changing
//
//	@param callback: Function callback.
func OnGroupSetParent_Register(callback SetParentCallback) Status {
	plugify.Log("permissions::OnGroupSetParent_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupSetParent_Register(__callback))
	return __retVal
}

// OnGroupSetParent_Unregister
//
//	@brief Unregister listener on group parent changing
//
//	@param callback: Function callback.
func OnGroupSetParent_Unregister(callback SetParentCallback) Status {
	plugify.Log("permissions::OnGroupSetParent_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupSetParent_Unregister(__callback))
	return __retVal
}

// OnGroupSetCookie_Register
//
//	@brief Register listener on group cookie sets
//
//	@param callback: Function callback.
func OnGroupSetCookie_Register(callback SetCookieGroupCallback) Status {
	plugify.Log("permissions::OnGroupSetCookie_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupSetCookie_Register(__callback))
	return __retVal
}

// OnGroupSetCookie_Unregister
//
//	@brief Unregister listener on group cookie sets
//
//	@param callback: Function callback.
func OnGroupSetCookie_Unregister(callback SetCookieGroupCallback) Status {
	plugify.Log("permissions::OnGroupSetCookie_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupSetCookie_Unregister(__callback))
	return __retVal
}

// OnGroupPermissionChange_Register
//
//	@brief Register listener on group permissions add/remove
//
//	@param callback: Function callback.
func OnGroupPermissionChange_Register(callback GroupPermissionCallback) Status {
	plugify.Log("permissions::OnGroupPermissionChange_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupPermissionChange_Register(__callback))
	return __retVal
}

// OnGroupPermissionChange_Unregister
//
//	@brief Unregister listener on group permissions add/remove
//
//	@param callback: Function callback.
func OnGroupPermissionChange_Unregister(callback GroupPermissionCallback) Status {
	plugify.Log("permissions::OnGroupPermissionChange_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupPermissionChange_Unregister(__callback))
	return __retVal
}

// OnGroupCreate_Register
//
//	@brief Register listener on group creation
//
//	@param callback: Function callback.
func OnGroupCreate_Register(callback GroupCreateCallback) Status {
	plugify.Log("permissions::OnGroupCreate_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupCreate_Register(__callback))
	return __retVal
}

// OnGroupCreate_Unregister
//
//	@brief Unregister listener on group creation
//
//	@param callback: Listener
func OnGroupCreate_Unregister(callback GroupCreateCallback) Status {
	plugify.Log("permissions::OnGroupCreate_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupCreate_Unregister(__callback))
	return __retVal
}

// OnGroupDelete_Register
//
//	@brief Register listener on group deletion
//
//	@param callback: Listener
func OnGroupDelete_Register(callback GroupDeleteCallback) Status {
	plugify.Log("permissions::OnGroupDelete_Register", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupDelete_Register(__callback))
	return __retVal
}

// OnGroupDelete_Unregister
//
//	@brief Unregister listener on group deletion
//
//	@param callback: Listener
func OnGroupDelete_Unregister(callback GroupDeleteCallback) Status {
	plugify.Log("permissions::OnGroupDelete_Unregister", plugify.Trace, 2)
	var __retVal Status
	__callback := plugify.GetFunctionPointerForDelegate(callback)
	__retVal = int32(C.OnGroupDelete_Unregister(__callback))
	return __retVal
}
