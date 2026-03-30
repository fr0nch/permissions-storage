#include "shared.h"

PLUGIFY_EXPORT int32_t (*__permissions_SetParent)(uint64_t, String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetParent)(String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_DumpPermissionsGroup)(String*, Vector*) = NULL;


PLUGIFY_EXPORT Vector (*__permissions_GetAllGroups)() = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_HasPermissionGroup)(String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_HasParentGroup)(String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetPriorityGroup)(String*, int32_t*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_AddPermissionGroup)(uint64_t, String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_RemovePermissionGroup)(uint64_t, String*, String*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetCookieGroup)(String*, String*, Variant*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_SetCookieGroup)(uint64_t, String*, String*, Variant*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_GetAllCookiesGroup)(String*, Vector*, Vector*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_CreateGroup)(uint64_t, String*, Vector*, int32_t, String*) = NULL;


PLUGIFY_EXPORT void (*__permissions_LoadGroups)(uint64_t) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_DeleteGroup)(uint64_t, String*) = NULL;


PLUGIFY_EXPORT bool (*__permissions_GroupExists)(String*) = NULL;


