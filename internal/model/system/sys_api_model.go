package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysApiModel = (*customSysApiModel)(nil)

type (
	// SysApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiModel.
	SysApiModel interface {
		sysApiModel
	}

	customSysApiModel struct {
		*defaultSysApiModel
		db *gorm.DB
	}
)

// NewSysApiModel returns a model for the database table.
func NewSysApiModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysApiModel {
	return &customSysApiModel{
		defaultSysApiModel: newSysApiModel(conn, c, opts...),
		db:                 db,
	}
}
