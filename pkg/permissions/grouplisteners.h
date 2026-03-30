#pragma once

#include "shared.h"

extern int32_t (*__permissions_OnLoadGroups_Register)(void*);

static int32_t OnLoadGroups_Register(void* callback) {
	return __permissions_OnLoadGroups_Register(callback);
}

extern int32_t (*__permissions_OnLoadGroups_Unregister)(void*);

static int32_t OnLoadGroups_Unregister(void* callback) {
	return __permissions_OnLoadGroups_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupSetParent_Register)(void*);

static int32_t OnGroupSetParent_Register(void* callback) {
	return __permissions_OnGroupSetParent_Register(callback);
}

extern int32_t (*__permissions_OnGroupSetParent_Unregister)(void*);

static int32_t OnGroupSetParent_Unregister(void* callback) {
	return __permissions_OnGroupSetParent_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupSetCookie_Register)(void*);

static int32_t OnGroupSetCookie_Register(void* callback) {
	return __permissions_OnGroupSetCookie_Register(callback);
}

extern int32_t (*__permissions_OnGroupSetCookie_Unregister)(void*);

static int32_t OnGroupSetCookie_Unregister(void* callback) {
	return __permissions_OnGroupSetCookie_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupPermissionChange_Register)(void*);

static int32_t OnGroupPermissionChange_Register(void* callback) {
	return __permissions_OnGroupPermissionChange_Register(callback);
}

extern int32_t (*__permissions_OnGroupPermissionChange_Unregister)(void*);

static int32_t OnGroupPermissionChange_Unregister(void* callback) {
	return __permissions_OnGroupPermissionChange_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupCreate_Register)(void*);

static int32_t OnGroupCreate_Register(void* callback) {
	return __permissions_OnGroupCreate_Register(callback);
}

extern int32_t (*__permissions_OnGroupCreate_Unregister)(void*);

static int32_t OnGroupCreate_Unregister(void* callback) {
	return __permissions_OnGroupCreate_Unregister(callback);
}

extern int32_t (*__permissions_OnGroupDelete_Register)(void*);

static int32_t OnGroupDelete_Register(void* callback) {
	return __permissions_OnGroupDelete_Register(callback);
}

extern int32_t (*__permissions_OnGroupDelete_Unregister)(void*);

static int32_t OnGroupDelete_Unregister(void* callback) {
	return __permissions_OnGroupDelete_Unregister(callback);
}

