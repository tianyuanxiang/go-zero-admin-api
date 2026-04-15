package plate

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlateProductionEventsModel = (*customPlateProductionEventsModel)(nil)

type (
	// PlateProductionEventsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlateProductionEventsModel.
	PlateProductionEventsModel interface {
		plateProductionEventsModel
	}

	customPlateProductionEventsModel struct {
		*defaultPlateProductionEventsModel
	}
)

// NewPlateProductionEventsModel returns a model for the database table.
func NewPlateProductionEventsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlateProductionEventsModel {
	return &customPlateProductionEventsModel{
		defaultPlateProductionEventsModel: newPlateProductionEventsModel(conn, c, opts...),
	}
}
