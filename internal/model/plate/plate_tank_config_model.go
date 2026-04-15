package plate

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlateTankConfigModel = (*customPlateTankConfigModel)(nil)

type (
	// PlateTankConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlateTankConfigModel.
	PlateTankConfigModel interface {
		plateTankConfigModel
	}

	customPlateTankConfigModel struct {
		*defaultPlateTankConfigModel
	}
)

// NewPlateTankConfigModel returns a model for the database table.
func NewPlateTankConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlateTankConfigModel {
	return &customPlateTankConfigModel{
		defaultPlateTankConfigModel: newPlateTankConfigModel(conn, c, opts...),
	}
}
