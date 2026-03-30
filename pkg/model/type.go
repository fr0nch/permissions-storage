package model

import (
	"time"
)

type Group struct {
	Permissions   []string
	Name          string
	ID            int
	InheritanceID *int
	Options       map[string]string
	Priority      int32
	Default       bool
}

type Cookie map[string]string

type User struct {
	Groups      []UserGroup
	Permissions []UserPermission
	Name        string
	UserID      UserID
	Cookies     Cookie
	Immunity    int32
}

type UserGroup struct {
	Expires   time.Time
	GroupName string
	GroupID   int
}

type UserPermission struct {
	Expires    time.Time
	Permission string
}

// ===================================
// UserGroups methods
// ===================================

//// AddGroup добавляет группу пользователю
//func (u *User) AddGroup(groupID int, groupName string, expires time.Time, status ActiveStatus) *UserGroup {
//	u.mu.Lock()
//	defer u.mu.Unlock()
//
//	newGroup := UserGroup{
//		GroupID:   groupID,
//		GroupName: groupName,
//		Expires:   expires,
//		Status:    status,
//	}
//
//	u.Groups = append(u.Groups, newGroup)
//	return &u.Groups[len(u.Groups)-1]
//}
//
//// RemoveGroup удаляет группу у пользователя
//func (u *User) RemoveGroup(groupName string) bool {
//	u.mu.Lock()
//	defer u.mu.Unlock()
//
//	length := len(u.Groups)
//
//	u.Groups = slices.DeleteFunc(u.Groups, func(p UserGroup) bool {
//		return p.GroupName == groupName
//	})
//
//	return length != len(u.Permissions)
//}
//
//// GroupExpires возвращает время истечения группы пользователя
//func (u *User) GroupExpires(groupName string) (time.Time, bool) {
//	for i := range u.Groups {
//		if u.Groups[i].GroupName == groupName {
//			return u.Groups[i].Expires, true
//		}
//	}
//
//	return time.Time{}, false
//}
//
//// FindGroupByID возвращает группу пользователя по ID группы
//func (u *User) FindGroupByID(groupID int) *UserGroup {
//	u.mu.RLock()
//	defer u.mu.RUnlock()
//
//	for i := range u.Groups {
//		if u.Groups[i].GroupID == groupID {
//			return &u.Groups[i]
//		}
//	}
//	return nil
//}
//
//// FindGroupByName возвращает группу пользователя по имени группы
//func (u *User) FindGroupByName(groupName string) *UserGroup {
//	u.mu.RLock()
//	defer u.mu.RUnlock()
//
//	for i := range u.Groups {
//		if strings.EqualFold(u.Groups[i].GroupName, groupName) {
//			return &u.Groups[i]
//		}
//	}
//	return nil
//}
//
//// ===================================
//// UserPermissions methods
//// ===================================
//
//// AddPermission добавляет пермишен пользователю
//func (u *User) AddPermission(permission string, expires time.Time, status ActiveStatus) *UserPermission {
//	u.mu.Lock()
//	defer u.mu.Unlock()
//
//	newPermission := UserPermission{
//		Permission: permission,
//		Expires:    expires,
//		Status:     status,
//	}
//
//	u.Permissions = append(u.Permissions, newPermission)
//	return &u.Permissions[len(u.Permissions)-1]
//}
//
//// RemovePermission удаляет пермишен у пользователя
//func (u *User) RemovePermission(permission string) bool {
//	u.mu.Lock()
//	defer u.mu.Unlock()
//
//	length := len(u.Permissions)
//
//	u.Permissions = slices.DeleteFunc(u.Permissions, func(p UserPermission) bool {
//		return p.Permission == permission
//	})
//
//	return length != len(u.Permissions)
//}
//
//// PermissionExpires возвращает время истечения пермишена пользователя
//func (u *User) PermissionExpires(permission string) *time.Time {
//	u.mu.RLock()
//	defer u.mu.RUnlock()
//
//	for i := range u.Permissions {
//		if u.Permissions[i].Permission == permission {
//			return &u.Permissions[i].Expires
//		}
//	}
//
//	return nil
//}
//
//// FindPermission ищет пермишен по его имени
//func (u *User) FindPermission(permission string) *UserPermission {
//	u.mu.RLock()
//	defer u.mu.RUnlock()
//
//	for i := range u.Permissions {
//		if u.Permissions[i].Permission == permission {
//			return &u.Permissions[i]
//		}
//	}
//	return nil
//}
