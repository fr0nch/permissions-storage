package main

import (
	"context"
	"permissions-storage/pkg/model"
	"permissions-storage/pkg/permissions"
	"time"
)

const contextTimeout = 15 * time.Second

func (p *PermissionsStoragePlugin) RegisteringEvents() {
	permissions.OnLoadUser_Register(p.OnLoadUser)
	permissions.OnLoadGroups_Register(p.OnLoadGroups)

	permissions.OnUserPermissionChange_Register(p.OnUserPermissionChange)
	permissions.OnUserGroupChange_Register(p.OnUserGroupChange)
	permissions.OnUserSetCookie_Register(p.OnUserSetCookie)

	permissions.OnPermissionExpirationCallback_Register(p.OnPermissionExpiration)
	permissions.OnGroupExpirationCallback_Register(p.OnGroupExpiration)
}

func (p *PermissionsStoragePlugin) OnLoadUser(pluginID uint64, targetID uint64, username string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting LoadUser for %s[%d], DB not ready: %v", username, targetID, err)
			return
		}

		user, err := p.storage.LoadUser(ctx, targetID, username)
		if err != nil {
			p.log.Errorf("Error loading user: %v", err)
			return
		}

		if user.Name != username {
			user.Name = username
		}

		err = p.storage.UpdateUser(ctx, user)
		if err != nil {
			p.log.Errorf("Error updating user: %v", err)
		}

		permissions.CreateUser(p.pluginID, targetID, user.Immunity, nil, nil)

		p.log.Debugf("Loading user: \n - \tUserID: %d, \n - \tGroups: %v, \n - \tPermissions: %v, \n - \tCookies: %v", user.UserID, user.Groups, user.Permissions, user.Cookies)

		if len(p.groups) > 0 {
			if p.defaultGroupID > 0 {
				p.log.Debugf("Setting user default group.")
				permissions.AddGroup(p.pluginID, targetID, p.GetGroupNameByID(p.defaultGroupID), 0, true)
			}

			for i := range user.Groups {
				var expires int64
				if !user.Groups[i].Expires.IsZero() {
					expires = user.Groups[i].Expires.Unix()
				}

				p.log.Debugf("AddGroup: name='%s', expires=%d(%v)", user.Groups[i].GroupName, expires, user.Groups[i].Expires)
				permissions.AddGroup(p.pluginID, targetID, user.Groups[i].GroupName, expires, true)
			}
		} else {
			p.log.Warn("Server groups not found. Skipping user group assignment.")
		}

		for i := range user.Permissions {
			var expires int64
			if !user.Permissions[i].Expires.IsZero() {
				expires = user.Permissions[i].Expires.Unix()
			}

			p.log.Debugf("AddPermission: name='%s', expires=%d", user.Permissions[i].Permission, expires)
			permissions.AddPermission(p.pluginID, targetID, user.Permissions[i].Permission, expires, true)
		}

		for key, value := range user.Cookies {
			p.log.Debugf("SetCookie: сookie='%s', value='%s'", key, value)
			permissions.SetCookie(p.pluginID, targetID, key, value, true)
		}

		permissions.LoadedUser(p.pluginID, targetID)
	}()
}

func (p *PermissionsStoragePlugin) OnLoadGroups(pluginID uint64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting LoadGroups, DB not ready: %v", err)
			return
		}

		groups, defaultGroupID, err := p.storage.LoadGroups(ctx)

		p.log.Info("Loading groups from database.")

		if err != nil {
			p.log.Errorf("Error loading groups: %v", err)
			return
		}

		if defaultGroupID == 0 {
			p.log.Warn("Default group ID is not set in storage")
		}

		p.defaultGroupID = defaultGroupID

		if groups == nil {
			p.log.Error("No groups found in storage")
			return
		}

		p.groups = groups

		p.log.Debugf("Loaded Groups: %v", p.groups)

		for i := range p.groups {
			result := permissions.CreateGroup(p.pluginID, p.groups[i].Name, p.groups[i].Permissions, p.groups[i].Priority, "")

			if result != permissions.Status_Success {
				continue
			}

			for cookie, value := range p.groups[i].Options {
				result = permissions.SetCookieGroup(p.pluginID, p.groups[i].Name, cookie, value)

				if result != permissions.Status_Success {
					continue
				}
			}

		}

		p.log.Info("Setting inheritance for groups.\n")

		for i := range p.groups {
			var result permissions.Status

			if p.groups[i].InheritanceID == nil {
				p.log.Warnf("Group '%s' dont have inheritance.\n", p.groups[i].Name)
				continue
			}

			inheritanceName := p.GetGroupNameByID(*p.groups[i].InheritanceID)

			result = permissions.SetParent(p.pluginID, p.groups[i].Name, inheritanceName)

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

			p.log.Debugf("Successfully set ('%s')[%v] inheritance for group '%s'.\n", inheritanceName, p.groups[i].InheritanceID, p.groups[i].Name)
		}

		p.log.Info("Groups loaded successfully.")
	}()
}

func (p *PermissionsStoragePlugin) OnUserPermissionChange(pluginID uint64, action permissions.Action, targetID uint64, perm string, timestamp int64) {
	if pluginID == p.pluginID {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting UserPermissionChange [ %v | perm: '%s' | timestamp: %d ], DB not ready: %v", targetID, perm, timestamp, err)
			return
		}

		permission := model.UserPermission{
			Permission: perm,
			Expires:    time.Unix(timestamp, 0),
		}

		switch action {
		case permissions.Action_Add /*, permissions.Action_Set*/ :
			err := p.storage.AddPermission(ctx, targetID, &permission)
			if err != nil {
				p.log.Errorf("Failed to add permission: %v", err)
				return
			}
		case permissions.Action_Remove:
			if !p.config.Settings.DeleteRemovedRecords {
				return
			}

			err := p.storage.RemovePermission(ctx, targetID, &permission)
			if err != nil {
				p.log.Errorf("Failed to remove permission: %v", err)
				return
			}
		}
	}()

}

func (p *PermissionsStoragePlugin) OnUserGroupChange(pluginID uint64, action permissions.Action, targetID uint64, group string, timestamp int64) {
	if pluginID == p.pluginID {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting UserGroupChange [ %v | group: '%s' | timestamp: %d ], DB not ready: %v", targetID, group, timestamp, err)
			return
		}

		_group := model.UserGroup{
			GroupID: p.GetGroupIDByName(group),
			Expires: time.Unix(timestamp, 0),
		}

		switch action {
		case permissions.Action_Add /*, permissions.Action_Set*/ :
			err := p.storage.AddGroup(ctx, targetID, &_group)
			if err != nil {
				p.log.Errorf("Failed to add group: %v", err)
				return
			}
		case permissions.Action_Remove:
			if !p.config.Settings.DeleteRemovedRecords {
				return
			}

			err := p.storage.RemoveGroup(ctx, targetID, &_group)
			if err != nil {
				p.log.Errorf("Failed to remove group: %v", err)
				return
			}
		}
	}()
}

func (p *PermissionsStoragePlugin) OnUserSetCookie(pluginID uint64, targetID uint64, name string, cookie any) {
	if pluginID == p.pluginID {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting UserSetCookie [ %v | cookie: '%s' | data: %s ], DB not ready: %v", targetID, name, cookie, err)
			return
		}

		p.storage.SetCookie(targetID, name, cookie)
	}()
}

func (p *PermissionsStoragePlugin) OnPermissionExpiration(targetID uint64, perm string) {
	if !p.config.Settings.DeleteExpiredRecords {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting PermissionExpiration [ %v | perm: '%s' ], DB not ready: %v", targetID, perm, err)
			return
		}

		err := p.storage.RemovePermission(ctx, targetID, &model.UserPermission{Permission: perm})
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
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		if err := p.storage.WaitReady(ctx); err != nil {
			p.log.Errorf("Aborting GroupExpiration [ %v | group: '%s' ], DB not ready: %v", targetID, group, err)
			return
		}

		err := p.storage.RemoveGroup(ctx, targetID, &model.UserGroup{GroupID: p.GetGroupIDByName(group)})
		if err != nil {
			p.log.Errorf("Failed to remove expired group: %v", err)
			return
		}
	}()
}
