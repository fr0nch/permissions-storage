package storage

import (
	"context"
	"permissions-storage/pkg/model"
)

type Storage interface {
	CreateTables(ctx context.Context) error

	LoadGroups(ctx context.Context) (groups []*model.Group, defaultGroupID int, err error)

	LoadUser(ctx context.Context, userID model.UserID, username string) (user *model.User, err error)
	UpdateUser(ctx context.Context, user *model.User) error

	AddPermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error
	RemovePermission(ctx context.Context, userID model.UserID, permission *model.UserPermission) error

	AddGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error
	RemoveGroup(ctx context.Context, userID model.UserID, group *model.UserGroup) error

	SetCookie(userID model.UserID, key string, value any)

	Close()
}
