// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
	"go-zero-admin/internal/svc"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictDataLogic) DeleteDictData(dictDataId int64) error {
	_, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, dictDataId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return l.svcCtx.SysDictDataModel.SoftDeleteDictDataByDictDataId(l.ctx, dictDataId)
}
