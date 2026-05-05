// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"go-zero-admin/internal/common"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	casbinpkg "go-zero-admin/pkg/casbin"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) error {
	// 检查角色编码唯一性
	existRole, err := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, req.RoleCode)
	if err != nil && err != sqlx.ErrNotFound {
		l.Errorf("查询角色编码[%s]是否存在失败：%v", req.RoleCode, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if existRole != nil {
		return xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
	}

	// 获取接口信息 && 确认接口数量完整
	var apiModels []systemmodel.SysApi
	if len(req.ApiIds) > 0 {
		apiModels, err = l.svcCtx.SysApiModel.ListByIds(l.ctx, req.ApiIds)
		if err != nil {
			l.Errorf("查询接口信息失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(apiModels) != len(req.ApiIds) {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分api接口信息不存在")
		}
	}

	if len(req.MenuIds) > 0 {
		menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, req.MenuIds)
		if err != nil {
			l.Logger.Errorf("查询菜单信息失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(menus) != len(req.MenuIds) {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分菜单信息不存在或已删除")
		}
	}

	menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, req.MenuIds)
	if err != nil {
		l.Logger.Errorf("查询菜单信息失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if len(menus) != len(req.MenuIds) {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分菜单信息不存在或已删除")
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.插入角色
		roleId, err := l.svcCtx.SysRoleModel.InsertRoleTrans(l.ctx, tx, &systemmodel.SysRole{
			Name:   req.RoleName,
			Code:   req.RoleCode,
			Status: 1,
			Remark: req.Remark,
			Sort:   req.Sort,
		})
		if err != nil {
			l.Errorf("插入角色记录失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}

		// 2.插入关联菜单ID
		if len(req.MenuIds) > 0 {
			// 自动补全菜单祖先链，确保父菜单不会缺失
			req.MenuIds, err = common.CompleteMenuAncestors(l.ctx, l.svcCtx.SysMenuModel, req.MenuIds)
			if err != nil {
				l.Errorf("补全菜单祖先链失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
			roleMenus := make([]systemmodel.SysRoleMenu, 0, len(req.MenuIds))
			for _, menuId := range req.MenuIds {
				roleMenus = append(roleMenus, systemmodel.SysRoleMenu{
					RoleId: roleId,
					MenuId: menuId,
				})
			}
			_, err = l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus)
			if err != nil {
				l.Errorf("插入关联菜单ID失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		// 3.插入关联接口ID
		if len(req.ApiIds) > 0 {
			roleApis := make([]systemmodel.SysRoleApi, 0, len(req.ApiIds))
			for _, apiId := range req.ApiIds {
				roleApis = append(roleApis, systemmodel.SysRoleApi{
					RoleId: roleId,
					ApiId:  apiId,
				})
			}
			_, err = l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis)
			if err != nil {
				l.Errorf("插入关联菜单ID失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		return nil
	})
	if err != nil {
		l.Errorf("创建角色事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 创建casbin规则
	if len(apiModels) > 0 {
		rules := make([][]string, 0, len(apiModels))
		for _, api := range apiModels {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, req.RoleCode, rules); err != nil {
			l.Errorf("同步Casbin策略失败（角色编码: %s）：%v", req.RoleCode, err)
		}
	}

	return nil
}
