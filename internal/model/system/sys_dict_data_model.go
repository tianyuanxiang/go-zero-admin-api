package system

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysDictDataModel = (*customSysDictDataModel)(nil)

type (
	// SysDictDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictDataModel.
	SysDictDataModel interface {
		sysDictDataModel
	}

	customSysDictDataModel struct {
		*defaultSysDictDataModel
		db *gorm.DB
	}
)

// NewSysDictDataModel returns a model for the database table.
func NewSysDictDataModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysDictDataModel {
	return &customSysDictDataModel{
		defaultSysDictDataModel: newSysDictDataModel(conn, c, opts...),
		db:                      db,
	}
}
