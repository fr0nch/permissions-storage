package storage

import (
	"context"
	"permissions-storage/internal/pkg/model"
)

type Storage interface {
	User
	Group

	CreateTables(ctx context.Context) error
	WaitReady(ctx context.Context) error
	Close()
}

type User interface {
	LoadUser(ctx context.Context, userID model.UserID, username string) (user *model.User, err error)
	//RemoveUser(ctx context.Context, user *model.User) error
	AddUserPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission, wildcard bool) error
	RemoveUserPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error
	AddUserGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error
	RemoveUserGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error
	SetCookie(ctx context.Context, userID model.UserID, key string, value any)
}

type Group interface {
	LoadGroups(ctx context.Context) (groups []*model.Group, defaultGroupID model.GroupID, err error)
	CreateGroup(ctx context.Context, group *model.Group) error
	DeleteGroup(ctx context.Context, groupID model.GroupID) error
	AddGroupPermission(ctx context.Context, groupID model.GroupID, permission *model.Permission, wildcard bool) error
	RemoveGroupPermission(ctx context.Context, groupID model.GroupID, permission *model.Permission) error
	SetParentGroup(ctx context.Context, groupID model.GroupID, inheritanceID model.GroupID) error
	SetOption(ctx context.Context, groupID model.GroupID, key string, value any) error
}
