// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"
	"go-zero-admin/internal/logic/system/helper"
	"go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
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
		l.Errorf("更新接口查询旧记录失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if err := l.svcCtx.SysApiModel.Update(l.ctx, &system.SysApi{
		Id:          req.Id,
		ApiPath:     req.ApiPath,
		Method:      req.Method,
		ApiGroup:    req.Group,
		ApiName:     req.ApiName,
		Description: req.Remark,
	}); err != nil {
		l.Errorf("更新接口失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 判断path或method是否变更
	if oldApi.ApiPath != req.ApiPath || oldApi.Method != req.Method {
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
