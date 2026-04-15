package system

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysRoleModel = (*customSysRoleModel)(nil)

type (
	SysRoleModel interface {
		sysRoleModel
		FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error)
		FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error)
	}

	customSysRoleModel struct {
		*defaultSysRoleModel
		db *gorm.DB
	}
)

func NewSysRoleModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: newSysRoleModel(conn, c, opts...),
		db:                  db,
	}
}

// FindOneByRoleId 按角色ID查询未软删除的角色。
func (m *customSysRoleModel) FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error) {
	var role SysRole
	result := m.db.WithContext(ctx).Table("sys_role").
		Where("id = ? AND deleted_at IS NULL", roleId).
		First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &role, nil
}

// FindByIds 批量查询多个角色信息
func (m *customSysRoleModel) FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error) {
	var roles []*SysRole
	result := m.db.WithContext(ctx).Table("sys_role").
		Where("id IN ? AND deleted_at IS NULL", roleIds).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
