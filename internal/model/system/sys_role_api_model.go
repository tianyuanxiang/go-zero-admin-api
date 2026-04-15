package system

import (
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
