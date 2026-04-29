package common

import (
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/types"
)

// buildMenuTree 将扁平菜单列表递归构建为树形结构。
func BuildMenuTree(menus []*systemmodel.SysMenu, parentId int64) []types.MenuItem {
	result := make([]types.MenuItem, 0)
	for _, m := range menus {
		if m.ParentId != parentId {
			continue
		}
		menuNode := types.MenuItem{
			Id:        m.Id,
			ParentId:  m.ParentId,
			MenuName:  m.Name,
			Path:      m.MenuPath,
			Component: m.Component,
			Icon:      m.Icon,
			MenuType:  m.MenuType,
			Perms:     m.Permission,
			Sort:      m.Sort,
		}
		children := BuildMenuTree(menus, m.Id)
		if len(children) > 0 {
			menuNode.Children = children
		}
		result = append(result, menuNode)
	}
	return result
}
