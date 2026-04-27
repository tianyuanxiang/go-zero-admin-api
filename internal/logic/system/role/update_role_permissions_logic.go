// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	systemmodel "go-zero-admin/internal/model/system"
	casbinpkg "go-zero-admin/pkg/casbin"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermissionsLogic {
	return &UpdateRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRolePermissionsLogic) UpdateRolePermissions(req *types.UpdateRolePermissionsReq) error {
	// 检查角色是否存在
	existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 1.获取接口信息 && 确认接口数量完整
	var apiModels []systemmodel.SysApi
	if len(req.ApiIds) > 0 {
		apiModels, err = l.svcCtx.SysApiModel.ListByIds(l.ctx, req.ApiIds)
		if err != nil {
			l.Logger.Errorf("查询接口信息失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(apiModels) != len(req.ApiIds) {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分api接口信息不存在")
		}
	}

	// 2.开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 2.1 先删后插：更新关联菜单
		if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuTrans(l.ctx, tx, req.Id); err != nil {
			l.Logger.Errorf("删除角色[%d]旧菜单关联失败：%v", req.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(req.MenuIds) > 0 {
			roleMenus := make([]systemmodel.SysRoleMenu, 0, len(req.MenuIds))
			for _, menuId := range req.MenuIds {
				roleMenus = append(roleMenus, systemmodel.SysRoleMenu{
					RoleId: req.Id,
					MenuId: menuId,
				})
			}
			if _, err := l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus); err != nil {
				l.Logger.Errorf("插入角色[%d]新菜单关联失败：%v", req.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		// 2.2 先删后插：更新关联接口
		if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiTrans(l.ctx, tx, req.Id); err != nil {
			l.Logger.Errorf("删除角色[%d]旧接口关联失败：%v", req.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(req.ApiIds) > 0 {
			roleApis := make([]systemmodel.SysRoleApi, 0, len(req.ApiIds))
			for _, apiId := range req.ApiIds {
				roleApis = append(roleApis, systemmodel.SysRoleApi{
					RoleId: req.Id,
					ApiId:  apiId,
				})
			}
			if _, err := l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis); err != nil {
				l.Logger.Errorf("插入角色[%d]新接口关联失败：%v", req.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}
		return nil
	})

	if err != nil {
		l.Logger.Errorf("更新角色权限事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 3.同步Casbin策略（事务成功后执行）
	// 用新编码全量覆盖策略
	if len(apiModels) > 0 {
		rules := make([][]string, 0, len(apiModels))
		for _, api := range apiModels {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, existRole.Code, rules); err != nil {
			l.Logger.Errorf("同步Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
		}
	} else {
		// API列表为空，清除该角色所有策略
		if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, existRole.Code); err != nil {
			l.Logger.Errorf("清除Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
		}
	}

	return nil
}
