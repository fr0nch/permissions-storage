#pragma once

#include "shared.h"

extern int32_t (*__permissions_SetParent)(uint64_t, String*, String*);

static int32_t SetParent(uint64_t pluginID, String* childName, String* parentName) {
	return __permissions_SetParent(pluginID, childName, parentName);
}

extern int32_t (*__permissions_GetParent)(String*, String*);

static int32_t GetParent(String* groupName, String* parentName) {
	return __permissions_GetParent(groupName, parentName);
}

extern int32_t (*__permissions_DumpPermissionsGroup)(String*, Vector*);

static int32_t DumpPermissionsGroup(String* name, Vector* perms) {
	return __permissions_DumpPermissionsGroup(name, perms);
}

extern Vector (*__permissions_GetAllGroups)();

static Vector GetAllGroups() {
	return __permissions_GetAllGroups();
}

extern int32_t (*__permissions_HasPermissionGroup)(String*, String*);

static int32_t HasPermissionGroup(String* name, String* perm) {
	return __permissions_HasPermissionGroup(name, perm);
}

extern int32_t (*__permissions_HasParentGroup)(String*, String*);

static int32_t HasParentGroup(String* childName, String* parentName) {
	return __permissions_HasParentGroup(childName, parentName);
}

extern int32_t (*__permissions_GetPriorityGroup)(String*, int32_t*);

static int32_t GetPriorityGroup(String* groupName, int32_t* priority) {
	return __permissions_GetPriorityGroup(groupName, priority);
}

extern int32_t (*__permissions_AddPermissionGroup)(uint64_t, String*, String*);

static int32_t AddPermissionGroup(uint64_t pluginID, String* name, String* perm) {
	return __permissions_AddPermissionGroup(pluginID, name, perm);
}

extern int32_t (*__permissions_RemovePermissionGroup)(uint64_t, String*, String*);

static int32_t RemovePermissionGroup(uint64_t pluginID, String* name, String* perm) {
	return __permissions_RemovePermissionGroup(pluginID, name, perm);
}

extern int32_t (*__permissions_GetCookieGroup)(String*, String*, Variant*);

static int32_t GetCookieGroup(String* groupName, String* cookieName, Variant* value) {
	return __permissions_GetCookieGroup(groupName, cookieName, value);
}

extern int32_t (*__permissions_SetCookieGroup)(uint64_t, String*, String*, Variant*);

static int32_t SetCookieGroup(uint64_t pluginID, String* groupName, String* cookieName, Variant* value) {
	return __permissions_SetCookieGroup(pluginID, groupName, cookieName, value);
}

extern int32_t (*__permissions_GetAllCookiesGroup)(String*, Vector*, Vector*);

static int32_t GetAllCookiesGroup(String* groupName, Vector* cookieNames, Vector* values) {
	return __permissions_GetAllCookiesGroup(groupName, cookieNames, values);
}

extern int32_t (*__permissions_CreateGroup)(uint64_t, String*, Vector*, int32_t, String*);

static int32_t CreateGroup(uint64_t pluginID, String* name, Vector* perms, int32_t priority, String* parent) {
	return __permissions_CreateGroup(pluginID, name, perms, priority, parent);
}

extern void (*__permissions_LoadGroups)(uint64_t);

static void LoadGroups(uint64_t pluginID) {
	__permissions_LoadGroups(pluginID);
}

extern int32_t (*__permissions_DeleteGroup)(uint64_t, String*);

static int32_t DeleteGroup(uint64_t pluginID, String* name) {
	return __permissions_DeleteGroup(pluginID, name);
}

extern bool (*__permissions_GroupExists)(String*);

static bool GroupExists(String* name) {
	return __permissions_GroupExists(name);
}

