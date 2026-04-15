package plate

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlateWaterEventsModel = (*customPlateWaterEventsModel)(nil)

type (
	// PlateWaterEventsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlateWaterEventsModel.
	PlateWaterEventsModel interface {
		plateWaterEventsModel
	}

	customPlateWaterEventsModel struct {
		*defaultPlateWaterEventsModel
	}
)

// NewPlateWaterEventsModel returns a model for the database table.
func NewPlateWaterEventsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlateWaterEventsModel {
	return &customPlateWaterEventsModel{
		defaultPlateWaterEventsModel: newPlateWaterEventsModel(conn, c, opts...),
	}
}
