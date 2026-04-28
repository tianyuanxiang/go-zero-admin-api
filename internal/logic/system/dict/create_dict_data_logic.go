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
	_, err := l.svcCtx.SysDictDataModel.Insert(l.ctx, &system.SysDictData{
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
