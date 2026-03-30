#pragma once

#include "shared.h"

extern int32_t (*__permissions_DumpPermissions)(uint64_t, Vector*);

static int32_t DumpPermissions(uint64_t targetID, Vector* perms) {
	return __permissions_DumpPermissions(targetID, perms);
}

extern int32_t (*__permissions_CanAffectUser)(uint64_t, uint64_t);

static int32_t CanAffectUser(uint64_t actorID, uint64_t targetID) {
	return __permissions_CanAffectUser(actorID, targetID);
}

extern int32_t (*__permissions_SetImmunity)(uint64_t, int32_t);

static int32_t SetImmunity(uint64_t targetID, int32_t immunity) {
	return __permissions_SetImmunity(targetID, immunity);
}

extern int32_t (*__permissions_HasPermission)(uint64_t, String*);

static int32_t HasPermission(uint64_t targetID, String* perm) {
	return __permissions_HasPermission(targetID, perm);
}

extern int32_t (*__permissions_HasGroup)(uint64_t, String*);

static int32_t HasGroup(uint64_t targetID, String* groupName) {
	return __permissions_HasGroup(targetID, groupName);
}

extern int32_t (*__permissions_GetUserGroups)(uint64_t, Vector*);

static int32_t GetUserGroups(uint64_t targetID, Vector* outGroups) {
	return __permissions_GetUserGroups(targetID, outGroups);
}

extern int32_t (*__permissions_GetImmunity)(uint64_t, int32_t*);

static int32_t GetImmunity(uint64_t targetID, int32_t* immunity) {
	return __permissions_GetImmunity(targetID, immunity);
}

extern int32_t (*__permissions_AddPermission)(uint64_t, uint64_t, String*, int64_t, bool);

static int32_t AddPermission(uint64_t pluginID, uint64_t targetID, String* perm, int64_t timestamp, bool dontBroadcast) {
	return __permissions_AddPermission(pluginID, targetID, perm, timestamp, dontBroadcast);
}

extern int32_t (*__permissions_RemovePermission)(uint64_t, uint64_t, String*);

static int32_t RemovePermission(uint64_t pluginID, uint64_t targetID, String* perm) {
	return __permissions_RemovePermission(pluginID, targetID, perm);
}

extern int32_t (*__permissions_AddGroup)(uint64_t, uint64_t, String*, int64_t, bool);

static int32_t AddGroup(uint64_t pluginID, uint64_t targetID, String* groupName, int64_t timestamp, bool dontBroadcast) {
	return __permissions_AddGroup(pluginID, targetID, groupName, timestamp, dontBroadcast);
}

extern int32_t (*__permissions_RemoveGroup)(uint64_t, uint64_t, String*);

static int32_t RemoveGroup(uint64_t pluginID, uint64_t targetID, String* groupName) {
	return __permissions_RemoveGroup(pluginID, targetID, groupName);
}

extern int32_t (*__permissions_GetCookie)(uint64_t, String*, Variant*);

static int32_t GetCookie(uint64_t targetID, String* name, Variant* value) {
	return __permissions_GetCookie(targetID, name, value);
}

extern int32_t (*__permissions_SetCookie)(uint64_t, uint64_t, String*, Variant*, bool);

static int32_t SetCookie(uint64_t pluginID, uint64_t targetID, String* name, Variant* cookie, bool dontBroadcast) {
	return __permissions_SetCookie(pluginID, targetID, name, cookie, dontBroadcast);
}

extern int32_t (*__permissions_GetAllCookies)(uint64_t, Vector*, Vector*);

static int32_t GetAllCookies(uint64_t targetID, Vector* names, Vector* values) {
	return __permissions_GetAllCookies(targetID, names, values);
}

extern int32_t (*__permissions_CreateUser)(uint64_t, uint64_t, int32_t, Vector*, Vector*);

static int32_t CreateUser(uint64_t pluginID, uint64_t targetID, int32_t immunity, Vector* groupsList, Vector* permsList) {
	return __permissions_CreateUser(pluginID, targetID, immunity, groupsList, permsList);
}

extern int32_t (*__permissions_DeleteUser)(uint64_t, uint64_t);

static int32_t DeleteUser(uint64_t pluginID, uint64_t targetID) {
	return __permissions_DeleteUser(pluginID, targetID);
}

extern void (*__permissions_LoadUser)(uint64_t, uint64_t, String*);

static void LoadUser(uint64_t pluginID, uint64_t targetID, String* username) {
	__permissions_LoadUser(pluginID, targetID, username);
}

extern void (*__permissions_LoadedUser)(uint64_t, uint64_t);

static void LoadedUser(uint64_t pluginID, uint64_t targetID) {
	__permissions_LoadedUser(pluginID, targetID);
}

extern bool (*__permissions_UserExists)(uint64_t);

static bool UserExists(uint64_t targetID) {
	return __permissions_UserExists(targetID);
}

