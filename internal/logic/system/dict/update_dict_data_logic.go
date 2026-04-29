// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
	"go-zero-admin/internal/model/system"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictDataLogic) UpdateDictData(req *types.UpdateDictDataReq) error {
	_, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return l.svcCtx.SysDictDataModel.Update(l.ctx, &system.SysDictData{
		Id:        int64(req.Id),
		TypeId:    int64(req.DictTypeId),
		Label:     req.DictLabel,
		DictValue: req.DictValue,
		Sort:      int64(req.Sort),
		Status:    int64(req.Status),
		Remark:    req.Remark,
	})
}
