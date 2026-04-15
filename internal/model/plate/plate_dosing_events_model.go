package plate

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PlateDosingEventsModel = (*customPlateDosingEventsModel)(nil)

type (
	// PlateDosingEventsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPlateDosingEventsModel.
	PlateDosingEventsModel interface {
		plateDosingEventsModel
	}

	customPlateDosingEventsModel struct {
		*defaultPlateDosingEventsModel
	}
)

// NewPlateDosingEventsModel returns a model for the database table.
func NewPlateDosingEventsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PlateDosingEventsModel {
	return &customPlateDosingEventsModel{
		defaultPlateDosingEventsModel: newPlateDosingEventsModel(conn, c, opts...),
	}
}
