package system

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysRoleMenuModel = (*customSysRoleMenuModel)(nil)

type (
	SysRoleMenuModel interface {
		sysRoleMenuModel
		GetMenuIdsByRoleIds(ctx context.Context, roleIds []int64) ([]int64, error)
	}

	customSysRoleMenuModel struct {
		*defaultSysRoleMenuModel
		db *gorm.DB
	}
)

func NewSysRoleMenuModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysRoleMenuModel {
	return &customSysRoleMenuModel{
		defaultSysRoleMenuModel: newSysRoleMenuModel(conn, c, opts...),
		db:                      db,
	}
}

// GetMenuIdsByRoleIds 根据多个角色ID查询去重后的菜单ID列表。
func (m *customSysRoleMenuModel) GetMenuIdsByRoleIds(ctx context.Context, roleIds []int64) ([]int64, error) {
	if len(roleIds) == 0 {
		return []int64{}, nil
	}
	var menuIds []int64
	result := m.db.WithContext(ctx).Table("sys_role_menu").
		Where("role_id IN ?", roleIds).
		Distinct("menu_id").
		Pluck("menu_id", &menuIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return menuIds, nil
}
