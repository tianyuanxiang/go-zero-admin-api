// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
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
	dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictType.DeletedAt.Valid {
		l.Errorf("字典类型[%d]已删除", req.Id)
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 更新字典类型信息
	updates := make(map[string]interface{})
	if req.DictCode != nil {
		updates["code"] = *req.DictCode
	}
	if req.DictName != nil {
		updates["name"] = *req.DictName
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	if err = l.svcCtx.SysDictTypeModel.UpdateDictType(l.ctx, int64(req.Id), updates); err != nil {
		l.Errorf("更新DictType[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return err
}
