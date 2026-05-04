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
	dictData, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictData.DeletedAt.Valid {
		l.Errorf("字典数据[%d]已删除", req.Id)
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}
	// 更新字典类型信息
	updates := make(map[string]interface{})
	if req.DictTypeId != 0 {
		updates["type_id"] = req.DictTypeId
	}
	if req.DictLabel != "" {
		updates["label"] = req.DictLabel
	}
	if req.DictValue != "" {
		updates["dict_value"] = req.DictValue
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
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

	if err = l.svcCtx.SysDictDataModel.UpdateDictData(l.ctx, int64(req.Id), updates); err != nil {
		l.Errorf("更新DictData[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return err
}
