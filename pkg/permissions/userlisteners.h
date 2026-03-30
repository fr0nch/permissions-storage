#pragma once

#include "shared.h"

extern int32_t (*__permissions_OnLoadUser_Register)(void*);

static int32_t OnLoadUser_Register(void* callback) {
	return __permissions_OnLoadUser_Register(callback);
}

extern int32_t (*__permissions_OnLoadUser_Unregister)(void*);

static int32_t OnLoadUser_Unregister(void* callback) {
	return __permissions_OnLoadUser_Unregister(callback);
}

extern int32_t (*__permissions_OnLoadedUser_Register)(void*);

static int32_t OnLoadedUser_Register(void* callback) {
	return __permissions_OnLoadedUser_Register(callback);
}

extern int32_t (*__permissions_OnLoadedUser_Unregister)(void*);

static int32_t OnLoadedUser_Unregister(void* callback) {
	return __permissions_OnLoadedUser_Unregister(callback);
}

extern int32_t (*__permissions_OnUserPermissionChange_Register)(void*);

static int32_t OnUserPermissionChange_Register(void* callback) {
	return __permissions_OnUserPermissionChange_Register(callback);
}

extern int32_t (*__permissions_OnUserPermissionChange_Unregister)(void*);

static int32_t OnUserPermissionChange_Unregister(void* callback) {
	return __permissions_OnUserPermissionChange_Unregister(callback);
}

extern int32_t (*__permissions_OnUserSetCookie_Register)(void*);

static int32_t OnUserSetCookie_Register(void* callback) {
	return __permissions_OnUserSetCookie_Register(callback);
}

extern int32_t (*__permissions_OnUserSetCookie_Unregister)(void*);

static int32_t OnUserSetCookie_Unregister(void* callback) {
	return __permissions_OnUserSetCookie_Unregister(callback);
}

extern int32_t (*__permissions_OnUserGroupChange_Register)(void*);

static int32_t OnUserGroupChange_Register(void* callback) {
	return __permissions_OnUserGroupChange_Register(callback);
}

extern int32_t (*__permissions_OnUserGroupChange_Unregister)(void*);

static int32_t OnUserGroupChange_Unregister(void* callback) {
	return __permissions_OnUserGroupChange_Unregister(callback);
}

extern int32_t (*__permissions_OnUserCreate_Register)(void*);

static int32_t OnUserCreate_Register(void* callback) {
	return __permissions_OnUserCreate_Register(callback);
}

extern int32_t (*__permissions_OnUserCreate_Unregister)(void*);

static int32_t OnUserCreate_Unregister(void* callback) {
	return __permissions_OnUserCreate_Unregister(callback);
}

extern int32_t (*__permissions_OnUserDelete_Register)(void*);

static int32_t OnUserDelete_Register(void* callback) {
	return __permissions_OnUserDelete_Register(callback);
}

extern int32_t (*__permissions_OnUserDelete_Unregister)(void*);

static int32_t OnUserDelete_Unregister(void* callback) {
	return __permissions_OnUserDelete_Unregister(callback);
}

extern int32_t (*__permissions_OnPermissionExpirationCallback_Register)(void*);

static int32_t OnPermissionExpirationCallback_Register(void* callback) {
	return __permissions_OnPermissionExpirationCallback_Register(callback);
}

extern int32_t (*__permissions_OnPermissionExpirationCallback_Unregister)(void*);

static int32_t OnPermissionExpirationCallback_Unregister(void* callback) {
	return __permissions_OnPermissionExpirationCallback_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupExpirationCallback_Register)(void*);

static int32_t OnGroupExpirationCallback_Register(void* callback) {
	return __permissions_OnGroupExpirationCallback_Register(callback);
}

extern int32_t (*__permissions_OnGroupExpirationCallback_Unregister)(void*);

static int32_t OnGroupExpirationCallback_Unregister(void* callback) {
	return __permissions_OnGroupExpirationCallback_Unregister(callback);
}

