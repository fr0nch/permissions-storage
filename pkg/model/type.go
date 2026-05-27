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
	State      bool
}
