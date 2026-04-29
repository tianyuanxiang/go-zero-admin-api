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

type UpdateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictTypeLogic) UpdateDictType(req *types.UpdateDictTypeReq) error {
	_, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return l.svcCtx.SysDictTypeModel.Update(l.ctx, &system.SysDictType{
		Id:     int64(req.Id),
		Name:   req.DictName,
		Code:   req.DictType,
		Status: int64(req.Status),
		Remark: req.Remark,
	})
}
