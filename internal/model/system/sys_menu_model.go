package system

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysMenuModel = (*customSysMenuModel)(nil)

type (
	SysMenuModel interface {
		sysMenuModel
		ListByIds(ctx context.Context, ids []int64) ([]*SysMenu, error)
	}

	customSysMenuModel struct {
		*defaultSysMenuModel
		db *gorm.DB
	}
)

func NewSysMenuModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysMenuModel {
	return &customSysMenuModel{
		defaultSysMenuModel: newSysMenuModel(conn, c, opts...),
		db:                  db,
	}
}

// ListByIds 根据菜单ID列表批量查询未软删除的菜单，按父级和排序字段升序返回。
func (m *customSysMenuModel) ListByIds(ctx context.Context, ids []int64) ([]*SysMenu, error) {
	if len(ids) == 0 {
		return []*SysMenu{}, nil
	}
	var menus []*SysMenu
	result := m.db.WithContext(ctx).Table("sys_menu").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Order("parent_id ASC, sort ASC").
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}
