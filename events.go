package main

import (
	"context"
	"permissions-storage/pkg/model"
	"permissions-storage/pkg/permissions"
	"time"
)

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
		user, err := p.storage.LoadUser(context.Background(), targetID, username)
		if err != nil {
			p.log.Errorf("Error loading user: %v", err)
			return
		}

		if user.Name != username {
			user.Name = username
		}

		err = p.storage.UpdateUser(context.Background(), user)
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

				p.log.Debugf("AddGroup: \n\tName: '%s', \n\tExpires: %d(%v)", user.Groups[i].GroupName, expires, user.Groups[i].Expires)
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

			p.log.Debugf("AddPermission: \n\tName: '%s', \n\tExpires: %d", user.Permissions[i].Permission, expires)
			permissions.AddPermission(p.pluginID, targetID, user.Permissions[i].Permission, expires, true)
		}

		for key, value := range user.Cookies {
			p.log.Debugf("SetCookie: \n\tCookie: '%s', \n\tValue: '%s'", key, value)
			permissions.SetCookie(p.pluginID, targetID, key, value, true)
		}

		permissions.LoadedUser(p.pluginID, targetID)
	}()
}

func (p *PermissionsStoragePlugin) OnLoadGroups(pluginID uint64) {
	go func() {
		groups, defaultGroupID, err := p.storage.LoadGroups(context.Background())

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
					p.log.Warnf("Failed to set inheritance for group '%s' (parent group %s not found).\n", p.groups[i].Name, inheritanceName)
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
		permission := model.UserPermission{
			Permission: perm,
			Expires:    time.Unix(timestamp, 0),
		}

		switch action {
		case permissions.Action_Add:
			err := p.storage.AddPermission(context.Background(), targetID, &permission)
			if err != nil {
				p.log.Errorf("Failed to add permission: %v", err)
				return
			}
		case permissions.Action_Remove:
			if !p.config.Settings.DeleteRemovedRecords {
				return
			}

			err := p.storage.RemovePermission(context.Background(), targetID, &permission)
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
		_group := model.UserGroup{
			GroupID: p.GetGroupIDByName(group),
			Expires: time.Unix(timestamp, 0),
		}

		switch action {
		case permissions.Action_Add:
			err := p.storage.AddGroup(context.Background(), targetID, &_group)
			if err != nil {
				p.log.Errorf("Failed to add group: %v", err)
				return
			}
		case permissions.Action_Remove:
			if !p.config.Settings.DeleteRemovedRecords {
				return
			}

			err := p.storage.RemoveGroup(context.Background(), targetID, &_group)
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

	go p.storage.SetCookie(targetID, name, cookie)
}

func (p *PermissionsStoragePlugin) OnPermissionExpiration(targetID uint64, perm string) {
	if !p.config.Settings.DeleteExpiredRecords {
		return
	}

	go func() {
		err := p.storage.RemovePermission(context.Background(), targetID, &model.UserPermission{Permission: perm})
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
		err := p.storage.RemoveGroup(context.Background(), targetID, &model.UserGroup{GroupID: p.GetGroupIDByName(group)})
		if err != nil {
			p.log.Errorf("Failed to remove expired group: %v", err)
			return
		}
	}()
}
