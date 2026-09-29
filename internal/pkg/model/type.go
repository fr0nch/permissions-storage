package model

import (
	"strings"
	"time"
)

type GroupID int

type ID interface {
	UserID | GroupID
}

type Option map[string]string
type Cookie map[string]string

type Permission struct {
	Name  string
	State bool
}

func NewPermission(s string) Permission {
	return Permission{
		Name:  strings.TrimPrefix(s, "-"),
		State: !strings.HasPrefix(s, "-"),
	}
}

func (p *Permission) String() string {
	if !p.State {
		return "-" + p.Name
	}
	return p.Name
}

type Group struct {
	Permissions   []Permission
	Name          string
	ID            GroupID
	InheritanceID *GroupID
	Options       Option
	Priority      int32
	Default       bool
}

type User struct {
	Permissions []UserPermission
	Groups      []UserGroup
	Name        string
	UserID      UserID
	Cookies     Cookie
	Immunity    int32
}

type UserPermission struct {
	Permission
	Expires time.Time
}

type UserGroup struct {
	*Group
	Expires time.Time
}
