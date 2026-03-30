#include "shared.h"

PLUGIFY_EXPORT int32_t (*__permissions_DumpPermissions)(uint64_t, Vector*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_CanAffectUser)(uint64_t, uint64_t) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_SetImmunity)(uint64_t, int32_t) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_HasPermission)(uint64_t, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_HasGroup)(uint64_t, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetUserGroups)(uint64_t, Vector*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetImmunity)(uint64_t, int32_t*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_AddPermission)(uint64_t, uint64_t, String*, int64_t, bool) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_RemovePermission)(uint64_t, uint64_t, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_AddGroup)(uint64_t, uint64_t, String*, int64_t, bool) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_RemoveGroup)(uint64_t, uint64_t, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetCookie)(uint64_t, String*, Variant*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_SetCookie)(uint64_t, uint64_t, String*, Variant*, bool) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetAllCookies)(uint64_t, Vector*, Vector*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_CreateUser)(uint64_t, uint64_t, int32_t, Vector*, Vector*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_DeleteUser)(uint64_t, uint64_t) = NULL;


PLUGIFY_EXPORT void (*__permissions_LoadUser)(uint64_t, uint64_t, String*) = NULL;


PLUGIFY_EXPORT void (*__permissions_LoadedUser)(uint64_t, uint64_t) = NULL;


PLUGIFY_EXPORT bool (*__permissions_UserExists)(uint64_t) = NULL;


