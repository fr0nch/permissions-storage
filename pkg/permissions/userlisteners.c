#include "shared.h"

PLUGIFY_EXPORT int32_t (*__permissions_OnLoadUser_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnLoadUser_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnLoadedUser_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnLoadedUser_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserPermissionChange_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserPermissionChange_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserSetCookie_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserSetCookie_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserGroupChange_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserGroupChange_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserCreate_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserCreate_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserDelete_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnUserDelete_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnPermissionExpirationCallback_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnPermissionExpirationCallback_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupExpirationCallback_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupExpirationCallback_Unregister)(void*) = NULL;


