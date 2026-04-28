// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByTypeLogic {
	return &GetDictDataByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictDataByTypeLogic) GetDictDataByType(dictTypeId int64) (resp *types.ListDictDataResp, err error) {

	DictData, count, err := l.svcCtx.SysDictDataModel.ListByDictTypeId(l.ctx, dictTypeId)
	if err != nil {
		l.Errorf("根据字典类型查询字典数据失败 %v\n", err)
		return nil, err
	}

	list := make([]types.DictDataItem, 0, len(DictData))
	for _, item := range DictData {
		list = append(list, types.DictDataItem{
			Id:        int(item.Id),
			DictType:  int(item.TypeId),
			DictLabel: item.Label,
			DictValue: item.DictValue,
			Sort:      int(item.Sort),
			Status:    int(item.Status),
			Remark:    item.Remark,
		})
	}
	return &types.ListDictDataResp{
		Total: count,
		List:  list,
	}, err
}
