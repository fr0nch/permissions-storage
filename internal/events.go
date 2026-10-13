package internal

import (
	"context"
	"permissions-storage/internal/pkg/model"
	"time"

	permissions "github.com/fr0nch/go-plugify-permissions"
)

const contextTimeout = 15 * time.Second

func (p *PermissionsStoragePlugin) RegisteringEvents() {
	permissions.OnUserRequest_Register(p.OnUserRequest)
	permissions.OnGroupsLoad_Register(p.OnLoadGroups)

	permissions.OnUserPermissionChangeStorage_Register(p.OnUserPermissionChange)
	permissions.OnUserGroupChangeStorage_Register(p.OnUserGroupChange)
	permissions.OnUserCookieChangeStorage_Register(p.OnUserCookieChange)

	permissions.OnPermissionExpiration_Register(p.OnPermissionExpiration)
	permissions.OnGroupExpiration_Register(p.OnGroupExpiration)

	permissions.OnGroupCreateStorage_Register(p.OnGroupCreate)
	permissions.OnGroupDeleteStorage_Register(p.OnGroupDelete)

	permissions.OnSetParentStorage_Register(p.OnGroupSetParent)

	permissions.OnGroupPermissionChangeStorage_Register(p.OnGroupPermissionChange)
	permissions.OnGroupOptionChangeStorage_Register(p.OnGroupOptionChange)
}

func (p *PermissionsStoragePlugin) OnUserRequest(pluginID int64, targetID uint64, username string, offline bool) bool {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting LoadUser for %s[%d], DB not ready: %v", username, targetID, err)
		return false
	}

	user, err := p.storage.LoadUser(ctx, targetID, username)
	if err != nil {
		p.log.Errorf("Error loading user: %v", err)
		return false
	}

	// Doubtful... but OK...
	userGroups := user.Groups[:0]
	for i := range user.Groups {
		group := p.GroupByID(user.Groups[i].ID)
		if group == nil {
			p.log.Warnf("Skipping group %d of user %d: group not loaded", user.Groups[i].ID, targetID)
			continue
		}

		user.Groups[i].Group = group
		userGroups = append(userGroups, user.Groups[i])
	}
	user.Groups = userGroups

	// TODO: Figure out a better way to implement this in the future.
	if permissions.UserExists(targetID) != permissions.PlayerState_NotFound {
		permissions.DeleteUser(p.PluginID, targetID)
	}

	permissions.CreateUser(p.PluginID, targetID, user.Immunity, false, nil)

	p.log.Debugf("Loading user: \n - \tUserID: %d, \n - \tGroups: %v, \n - \tPermissions: %v, \n - \tCookies: %v", user.UserID, user.Groups, user.Permissions, user.Cookies)

	if len(p.groups) > 0 {
		if defaultGroup := p.GroupByID(p.defaultGroupID); defaultGroup != nil {
			p.log.Debugf("Setting user default group.")
			permissions.AddGroup(p.PluginID, targetID, defaultGroup.Name, 0, true)
		}
		for i := range user.Groups {
			var expires int64
			if !user.Groups[i].Expires.IsZero() {
				expires = user.Groups[i].Expires.Unix()
			}

			p.log.Debugf("AddGroup: name='%s', expires=%d(%v)", user.Groups[i].Name, expires, user.Groups[i].Expires)
			permissions.AddGroup(p.PluginID, targetID, user.Groups[i].Name, expires, true)
		}
	} else {
		p.log.Warn("Server groups not found. Skipping user group assignment.")
	}

	for i := range user.Permissions {
		var expires int64
		if !user.Permissions[i].Expires.IsZero() {
			expires = user.Permissions[i].Expires.Unix()
		}

		p.log.Debugf("AddPermission: name='%s', state='%v', expires=%d", user.Permissions[i].Name, user.Permissions[i].State, expires)
		permissions.AddPermission(p.PluginID, targetID, user.Permissions[i].String(), expires, true)
	}

	for key, value := range user.Cookies {
		p.log.Debugf("SetCookie: сookie='%s', value='%s'", key, value)
		permissions.SetCookie(p.PluginID, targetID, key, value, false)
	}

	//state := permissions.PlayerState_Online
	//if offline {
	//	state = permissions.PlayerState_Offline
	//}

	//callback(p.PluginID, targetID, state)
	return true
}

func (p *PermissionsStoragePlugin) OnLoadGroups(pluginID int64) bool {
	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting LoadGroups, DB not ready: %v", err)
		return false
	}

	p.log.Info("Loading groups from database.")
	groups, defaultGroupID, err := p.storage.LoadGroups(ctx)
	if err != nil {
		p.log.Errorf("Error loading groups: %v", err)
		return false
	}

	if groups == nil {
		p.log.Error("Groups not found in storage")
		return false
	}
	p.groups = groups

	p.defaultGroupID = defaultGroupID
	if defaultGroupID == 0 {
		p.log.Warn("Default group ID is not set in storage")
	}

	p.log.Debugf("Loaded Groups: %v", p.groups)

	for i := range p.groups {
		var perms []string
		for perm := range p.groups[i].Permissions {
			perms = append(perms, p.groups[i].Permissions[perm].String())
		}

		result := permissions.CreateGroup(p.PluginID, p.groups[i].Name, perms, p.groups[i].Priority, "", true)
		if result != permissions.Status_Success {
			continue
		}

		for option, value := range p.groups[i].Options {
			p.log.Debugf("SetOptionGroup: %v | %v | %v | %v", p.PluginID, p.groups[i].Name, option, value)
			result = permissions.SetOptionGroup(p.PluginID, p.groups[i].Name, option, value, true)
			if result != permissions.Status_Success {
				continue
			}
		}
	}

	p.log.Debugf("Setting inheritance for groups.\n")

	for i := range p.groups {
		var result permissions.Status

		if p.groups[i].InheritanceID == nil {
			p.log.Warnf("Group '%s' dont have inheritance.\n", p.groups[i].Name)
			continue
		}

		inheritanceGroup := p.GroupByID(*p.groups[i].InheritanceID)
		if inheritanceGroup == nil {
			p.log.Warnf("Group inheritance '%d' dont exist.\n", *p.groups[i].InheritanceID)
			continue
		}

		inheritanceName := p.GroupByID(*p.groups[i].InheritanceID).Name
		result = permissions.SetParent(p.PluginID, p.groups[i].Name, inheritanceName, true)

		p.log.Debugf("Set inheritance ('%s')[%v] for group '%s'.", inheritanceName, p.groups[i].InheritanceID, p.groups[i].Name)

		switch result {
		case permissions.Status_ChildGroupNotFound:
			{
				p.log.Warnf("Failed to set inheritance for group '%s' (child group not found).\n", p.groups[i].Name)
				continue
			}
		case permissions.Status_ParentGroupNotFound:
			{
				p.log.Warnf("Failed to set inheritance for group '%s' (parent group '%s' not found).\n", p.groups[i].Name, inheritanceName)
				continue
			}
		}

		p.log.Debugf("Successfully set ('%s')[%d|%v] inheritance for group '%s'.\n", inheritanceName, *p.groups[i].InheritanceID, p.groups[i].InheritanceID, p.groups[i].Name)
	}

	p.log.Info("Groups loaded successfully.")
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnUserPermissionChange(pluginID int64, action permissions.Action, targetID uint64,
	perm string, oldState permissions.Status, newState permissions.Status, oldTimestamp int64, newTimestamp int64) bool {

	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting UserPermissionChange [ %v | perm: '%s' | timestamp: %d ], DB not ready: %v", targetID, perm, newTimestamp, err)
		return false
	}

	permission := model.UserPermission{
		Permission: model.Permission{
			Name:  perm,
			State: newState == permissions.Status_Allow,
		},
		Expires: time.Unix(newTimestamp, 0),
	}

	switch action {
	case permissions.Action_Add:
		err := p.storage.AddUserPermission(ctx, targetID, &permission, false)
		if err != nil {
			p.log.Errorf("Failed to add permission: %v", err)
			return false
		}
	case permissions.Action_Replace:
		err := p.storage.AddUserPermission(ctx, targetID, &permission, false)
		if err != nil {
			p.log.Errorf("Failed to replace permission: %v", err)
			return false
		}
	case permissions.Action_ReplaceToWC:
		err := p.storage.AddUserPermission(ctx, targetID, &permission, true)
		if err != nil {
			p.log.Errorf("Failed to replace permission with wildcard: %v", err)
			return false
		}
	case permissions.Action_Remove:
		err := p.storage.RemoveUserPermission(ctx, targetID, &permission)
		if err != nil {
			p.log.Errorf("Failed to remove permission: %v", err)
			return false
		}
	}
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnUserGroupChange(pluginID int64, action permissions.Action, targetID uint64,
	group string, oldTimestamp int64, newTimestamp int64) bool {

	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting UserGroupChange [ %v | group: '%s' | timestamp: %d ], DB not ready: %v", targetID, group, newTimestamp, err)
		return false
	}

	userGroup := p.GroupByName(group)
	if userGroup == nil {
		p.log.Errorf("Aborting UserGroupChange [ %v | group: '%s' ], group not found", targetID, group)
		return false
	}

	_group := model.UserGroup{
		Group:   userGroup,
		Expires: time.Unix(newTimestamp, 0),
	}

	switch action {
	case permissions.Action_Add:
		err := p.storage.AddUserGroup(ctx, targetID, &_group)
		if err != nil {
			p.log.Errorf("Failed to add group: %v", err)
			return false
		}
	case permissions.Action_Replace:
		err := p.storage.AddUserGroup(ctx, targetID, &_group)
		if err != nil {
			p.log.Errorf("Failed to replace group: %v", err)
			return false
		}
	case permissions.Action_Remove:
		err := p.storage.RemoveUserGroup(ctx, targetID, &_group)
		if err != nil {
			p.log.Errorf("Failed to remove group: %v", err)
			return false
		}
	}
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnUserCookieChange(pluginID int64, targetID uint64, name string, cookie any) bool {
	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting UserSetCookie [ %v | cookie: '%s' | data: %s ], DB not ready: %v", targetID, name, cookie, err)
		return false
	}

	p.storage.SetCookie(targetID, name, cookie)
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnPermissionExpiration(targetID uint64, perm string, state permissions.Status) {
	if !p.config.Settings.DeleteExpiredRecords {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting PermissionExpiration [ %v | perm: '%s' ], DB not ready: %v", targetID, perm, err)
			return
		}

		err := p.storage.RemoveUserPermission(ctx, targetID, &model.UserPermission{Permission: model.Permission{Name: perm}})
		if err != nil {
			p.log.Errorf("Failed to remove expired permission: %v", err)
			return
		}
	}()
}

func (p *PermissionsStoragePlugin) OnGroupExpiration(targetID uint64, group string) {
	if !p.config.Settings.DeleteExpiredRecords {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting GroupExpiration [ %v | group: '%s' ], DB not ready: %v", targetID, group, err)
			return
		}

		expiredGroup := p.GroupByName(group)
		if expiredGroup == nil {
			p.log.Errorf("Aborting GroupExpiration [ %v | group: '%s' ], group not found", targetID, group)
			return
		}

		err := p.storage.RemoveUserGroup(ctx, targetID, &model.UserGroup{Group: expiredGroup})
		if err != nil {
			p.log.Errorf("Failed to remove expired group: %v", err)
			return
		}
	}()
}

func (p *PermissionsStoragePlugin) OnGroupCreate(pluginID int64, name string, perms []string, priority int32, parent string) bool {
	p.log.Debugf("OnGroupCreated called!")

	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting CreateGroup, DB not ready: %v", err)
		return false
	}

	p.log.Debugf("perms: '%v'", perms)
	p.log.Debugf("perms count: %v", len(perms))

	var _perms []model.Permission
	for i := range perms {
		if len(perms[i]) > 0 {
			_perms = append(_perms, model.NewPermission(perms[i]))
		}
	}

	p.log.Debugf("_perms: '%v'", _perms)
	p.log.Debugf("_perms count: %v", len(_perms))

	var parentID *model.GroupID

	parentGroup := p.GroupByName(parent)
	if parentGroup != nil {
		parentID = &parentGroup.ID
	}

	group := model.Group{
		Name:          name,
		Priority:      priority,
		InheritanceID: parentID,
		Permissions:   _perms,
	}

	err := p.storage.CreateGroup(ctx, &group)
	if err != nil {
		p.log.Errorf("Failed to create group: %v", err)
		return false
	}

	p.groups = append(p.groups, &group)
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnGroupDelete(pluginID int64, name string) bool {
	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting DeleteGroup, DB not ready: %v", err)
		return false
	}

	group := p.GroupByName(name)
	if group == nil {
		p.log.Errorf("Aborting DeleteGroup [ group: '%s' ], group not found", name)
		return false
	}

	var groupID = group.ID
	err := p.storage.DeleteGroup(ctx, groupID)
	if err != nil {
		p.log.Errorf("Failed to delete group: %v", err)
		return false
	}

	if !p.RemoveGroup(groupID) {
		p.log.Errorf("failed to remove group '%v' with ID %v: group not found", name, groupID)
		return false
	}
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnGroupSetParent(pluginID int64, childName string, parentName string) bool {
	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting SetParent [ child: '%s' | parent: '%s' ], DB not ready: %v", childName, parentName, err)
		return false
	}

	group := p.GroupByName(childName)
	if group == nil {
		p.log.Errorf("Aborting SetParent [ child: '%s' | parent: '%s' ], child group not found", childName, parentName)
		return false
	}

	parentGroup := p.GroupByName(parentName)
	if parentGroup == nil {
		p.log.Errorf("Aborting SetParent [ child: '%s' | parent: '%s' ], parent group not found", childName, parentName)
		return false
	}

	err := p.storage.SetParentGroup(ctx, group.ID, parentGroup.ID)
	if err != nil {
		p.log.Errorf("Failed to set group parent: %v", err)
		return false
	}

	group.InheritanceID = &parentGroup.ID
	//}()

	return true
}

func (p *PermissionsStoragePlugin) OnGroupPermissionChange(pluginID int64, action permissions.Action, groupName string,
	perm string, oldState permissions.Status, newState permissions.Status) bool {

	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting GroupPermissionChange [ group: '%s' | perm: '%s' ], DB not ready: %v", groupName, perm, err)
		return false
	}

	group := p.GroupByName(groupName)
	if group == nil {
		p.log.Errorf("Aborting GroupPermissionChange [ group: '%s' | perm: '%s' ], group not found", groupName, perm)
		return false
	}

	permission := model.Permission{
		Name:  perm,
		State: newState == permissions.Status_Allow,
	}

	switch action {
	case permissions.Action_Add:
		err := p.storage.AddGroupPermission(ctx, group.ID, &permission, false)
		if err != nil {
			p.log.Errorf("Failed to add permission: %v", err)
			return false
		}
	case permissions.Action_Replace:
		err := p.storage.AddGroupPermission(ctx, group.ID, &permission, false)
		if err != nil {
			p.log.Errorf("Failed to replace permission: %v", err)
			return false
		}
	case permissions.Action_ReplaceToWC:
		err := p.storage.AddGroupPermission(ctx, group.ID, &permission, true)
		if err != nil {
			p.log.Errorf("Failed to replace permission with wildcard: %v", err)
			return false
		}
	case permissions.Action_Remove:
		err := p.storage.RemoveGroupPermission(ctx, group.ID, &permission)
		if err != nil {
			p.log.Errorf("Failed to remove permission: %v", err)
			return false
		}
	}

	//}()
	return true
}

func (p *PermissionsStoragePlugin) OnGroupOptionChange(pluginID int64, groupName string, optionName string, value any) bool {
	if pluginID == p.PluginID {
		return true
	}

	//go func() {
	ctx, cancel := context.WithTimeout(p.Plugin.Context(), contextTimeout)
	defer cancel()

	if err := p.storage.WaitReady(ctx); err != nil {
		p.log.Errorf("Aborting GroupOptionChange [ group: '%s' | option: '%s' ], DB not ready: %v", groupName, optionName, err)
		return false
	}

	group := p.GroupByName(groupName)
	if group == nil {
		p.log.Errorf("Aborting GroupOptionChange [ group: '%s' | option: '%s' ], group not found", groupName, optionName)
		return false
	}

	err := p.storage.SetOption(ctx, group.ID, optionName, value)
	if err != nil {
		p.log.Errorf("Failed to set group option: %v", err)
		return false
	}
	//}()

	return true
}
