package system

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysRoleApiModel = (*customSysRoleApiModel)(nil)

type (
	// SysRoleApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleApiModel.
	SysRoleApiModel interface {
		sysRoleApiModel
		InsertRoleApiTrans(ctx context.Context, tx *gorm.DB, roleApi []SysRoleApi) (int64, error)
		DeleteRoleApiByRoleIdTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
		DeleteRoleApiByApiIdTrans(ctx context.Context, tx *gorm.DB, apiId int64) error
		ListRoleIdsByApiId(ctx context.Context, apiId int64) ([]int64, error)
		ListApiIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	}

	customSysRoleApiModel struct {
		*defaultSysRoleApiModel
		db *gorm.DB
	}
)

// NewSysRoleApiModel returns a model for the database table.
func NewSysRoleApiModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysRoleApiModel {
	return &customSysRoleApiModel{
		defaultSysRoleApiModel: newSysRoleApiModel(conn, c, opts...),
		db:                     db,
	}
}

func (m *customSysRoleApiModel) InsertRoleApiTrans(ctx context.Context, tx *gorm.DB, roleApi []SysRoleApi) (int64, error) {
	result := tx.WithContext(ctx).Table("sys_role_api").Create(&roleApi)
	return result.RowsAffected, result.Error
}

func (m *customSysRoleApiModel) DeleteRoleApiByRoleIdTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
	return tx.WithContext(ctx).Where("role_id = ?", roleId).Delete(&SysRoleApi{}).Error
}

func (m *customSysRoleApiModel) DeleteRoleApiByApiIdTrans(ctx context.Context, tx *gorm.DB, apiId int64) error {
	return tx.WithContext(ctx).Where("api_id = ?", apiId).Delete(&SysRoleApi{}).Error
}

// ListRoleIdsByApiId 根据接口ID查询所有关联的角色ID。
func (m *customSysRoleApiModel) ListRoleIdsByApiId(ctx context.Context, apiId int64) ([]int64, error) {
	var roleIds []int64
	result := m.db.WithContext(ctx).Table("sys_role_api").
		Where("api_id = ?", apiId).
		Pluck("role_id", &roleIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleIds, nil
}

// ListApiIdsByRoleId 根据角色ID查询所有关联的接口ID。
func (m *customSysRoleApiModel) ListApiIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	var apiIds []int64
	result := m.db.WithContext(ctx).Table("sys_role_api").
		Where("role_id = ?", roleId).
		Pluck("api_id", &apiIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return apiIds, nil
}
