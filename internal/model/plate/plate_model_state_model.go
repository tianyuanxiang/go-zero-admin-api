package plate

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlateModelStateModel = (*customPlateModelStateModel)(nil)

type (
	// PlateModelStateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlateModelStateModel.
	PlateModelStateModel interface {
		plateModelStateModel
	}

	customPlateModelStateModel struct {
		*defaultPlateModelStateModel
	}
)

// NewPlateModelStateModel returns a model for the database table.
func NewPlateModelStateModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlateModelStateModel {
	return &customPlateModelStateModel{
		defaultPlateModelStateModel: newPlateModelStateModel(conn, c, opts...),
	}
}
