package main

import (
	"strings"
)

func (p *PermissionsStoragePlugin) GetGroupNameByID(groupID int) string {
	for i := range p.groups {
		if p.groups[i].ID == groupID {
			return p.groups[i].Name
		}
	}

	return ""
}

func (p *PermissionsStoragePlugin) GetGroupIDByName(groupName string) int {
	for i := range p.groups {
		if strings.EqualFold(p.groups[i].Name, groupName) {
			return p.groups[i].ID
		}
	}

	return -1
}
