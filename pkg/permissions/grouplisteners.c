#include "shared.h"

PLUGIFY_EXPORT int32_t (*__permissions_OnLoadGroups_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnLoadGroups_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupSetParent_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupSetParent_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupSetCookie_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupSetCookie_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupPermissionChange_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupPermissionChange_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupCreate_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupCreate_Unregister)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupDelete_Register)(void*) = NULL;


PLUGIFY_EXPORT int32_t (*__permissions_OnGroupDelete_Unregister)(void*) = NULL;


