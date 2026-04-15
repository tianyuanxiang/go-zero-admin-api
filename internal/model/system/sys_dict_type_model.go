package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysDictTypeModel = (*customSysDictTypeModel)(nil)

type (
	// SysDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictTypeModel.
	SysDictTypeModel interface {
		sysDictTypeModel
	}

	customSysDictTypeModel struct {
		*defaultSysDictTypeModel
		db *gorm.DB
	}
)

// NewSysDictTypeModel returns a model for the database table.
func NewSysDictTypeModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysDictTypeModel {
	return &customSysDictTypeModel{
		defaultSysDictTypeModel: newSysDictTypeModel(conn, c, opts...),
		db:                      db,
	}
}
