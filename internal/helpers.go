package internal

import (
	"permissions-storage/internal/pkg/model"
	"slices"
	"strings"
)

//func (p *PermissionsStoragePlugin) GroupNameByID(groupID int) string {
//	for i := range p.groups {
//		if p.groups[i].ID == groupID {
//			return p.groups[i].Name
//		}
//	}
//	return ""
//}
//
//func (p *PermissionsStoragePlugin) GroupIDByName(groupName string) int {
//	for i := range p.groups {
//		if strings.EqualFold(p.groups[i].Name, groupName) {
//			return p.groups[i].ID
//		}
//	}
//	return -1
//}

func (p *PermissionsStoragePlugin) GroupByName(groupName string) *model.Group {
	for i := range p.groups {
		if strings.EqualFold(p.groups[i].Name, groupName) {
			return p.groups[i]
		}
	}
	return nil
}

func (p *PermissionsStoragePlugin) GroupByID(groupID model.GroupID) *model.Group {
	for i := range p.groups {
		if p.groups[i].ID == groupID {
			return p.groups[i]
		}
	}
	return nil
}

func (p *PermissionsStoragePlugin) RemoveGroup(id model.GroupID) bool {
	length := len(p.groups)

	p.groups = slices.DeleteFunc(p.groups, func(g *model.Group) bool {
		return g.ID == id
	})

	return length != len(p.groups)
}
