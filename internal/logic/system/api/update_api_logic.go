// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"
	"go-zero-admin/internal/logic/system/helper"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApiLogic) UpdateApi(req *types.UpdateApiReq) error {
	oldApi, err := l.svcCtx.SysApiModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("更新api接口查询旧记录失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if oldApi.DeletedAt.Valid {
		return xerr.NewCodeError(xerr.ErrNotFound)
	}
	// 更新菜单基本信息
	updates := make(map[string]interface{})
	if req.ApiPath != nil {
		updates["api_path"] = *req.ApiPath
	}
	if req.Method != nil {
		updates["method"] = *req.Method
	}
	if req.Group != nil {
		updates["api_group"] = *req.Group
	}
	if req.Remark != nil {
		updates["description"] = *req.Remark
	}
	if req.ApiName != nil {
		updates["api_name"] = *req.ApiName
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	if err = l.svcCtx.SysApiModel.UpdateApi(l.ctx, req.Id, updates); err != nil {
		l.Errorf("更新api[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 判断path或method是否变更
	if oldApi.ApiPath != *req.ApiPath || oldApi.Method != *req.Method {
		// 同步Casbin
		roleIds, err := l.svcCtx.SysRoleApiModel.ListRoleIdsByApiId(l.ctx, req.Id)
		if err != nil {
			l.Logger.Errorf("查询API[%d]关联角色失败：%v", req.Id, err)
		} else {
			helper.RebuildCasbinByRoleIds(l.ctx, l.svcCtx, roleIds)
		}
	}

	return err
}
