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

type CreateDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictDataLogic {
	return &CreateDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictDataLogic) CreateDictData(req *types.CreateDictDataReq) error {
	// 创建数据之前先查是否有该字典类型
	dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, int64(req.DictTypeId))
	if err != nil {
		if err == sqlx.ErrNotFound {
			l.Errorf("创建字典数据时，字典类型[%d]不存在", req.DictTypeId)
			return xerr.NewCodeError(xerr.ErrParamInvalid)
		}
		l.Errorf("查询字典类型[%d]失败：%v", req.DictTypeId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictType.DeletedAt.Valid {
		l.Errorf("创建字典数据时，字典类型[%d]已被删除", req.DictTypeId)
		return xerr.NewCodeError(xerr.ErrNotFound)
	}

	_, err = l.svcCtx.SysDictDataModel.Insert(l.ctx, &system.SysDictData{
		TypeId:    int64(req.DictTypeId),
		Label:     req.DictLabel,
		DictValue: req.DictValue,
		Sort:      int64(req.Sort),
		Status:    int64(req.Status),
		Remark:    req.Remark,
	})
	if err != nil {
		l.Errorf("插入字典数据失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return err
}
