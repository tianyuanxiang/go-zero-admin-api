// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.CreateDictTypeReq) error {
	exist, err := l.svcCtx.SysDictTypeModel.FindOneByCode(l.ctx, req.DictCode)
	if err != nil || err != sqlx.ErrNotFound {
		l.Errorf("查询字典类型编码[%s]失败：%v", req.DictCode, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if exist != nil {
		return xerr.NewCodeErrorMsg(xerr.ErrDuplicate, "字典类型编码已存在")
	}

	_, err = l.svcCtx.SysDictTypeModel.Insert(l.ctx, &systemmodel.SysDictType{
		Name:   req.DictName,
		Code:   req.DictCode,
		Status: int64(req.Status),
		Remark: req.Remark,
	})
	if err != nil {
		l.Errorf("插入字典类型失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return nil
}
