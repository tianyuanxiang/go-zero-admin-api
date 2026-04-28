// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-zero-admin/internal/svc"
)

type DeleteDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictTypeLogic) DeleteDictType(dictTypeId int64) error {
	_, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, dictTypeId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Logger.Errorf("删除字典类型时查询dictTypeId[%d]是否存在 失败:%v\n", dictTypeId, err)
		return err
	}
	// 先删除该类型下所有字典数据
	if err = l.svcCtx.SysDictDataModel.SoftDeleteDictDataByDictTypeId(l.ctx, dictTypeId); err != nil {
		l.Errorf("删除字典类型[%d]下字典数据失败：%v", dictTypeId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	// 再删除该类型
	if err := l.svcCtx.SysDictTypeModel.SoftDeleteDictType(l.ctx, dictTypeId); err != nil {
		l.Logger.Errorf("删除字典类型失败:%v\n", err)
		return err
	}

	return err
}
