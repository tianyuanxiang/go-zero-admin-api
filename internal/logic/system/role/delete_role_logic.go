// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"go-zero-admin/internal/svc"
	casbinpkg "go-zero-admin/pkg/casbin"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(roleId int64) error {
	// 1. 检查角色是否存在
	result, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", roleId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 软删除角色记录
		if err = l.svcCtx.SysRoleModel.SoftDeleteRoleTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("软删除角色[%d]失败：%v", roleId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}

		// 2.清除角色与用户的关联
		if err = l.svcCtx.SysUserRoleModel.DeleteUserRoleTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]用户关联失败：%v", roleId, err)
		}

		// 3. 清除角色与菜单的关联
		if err = l.svcCtx.SysRoleMenuModel.DeleteRoleMenuTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]菜单关联失败：%v", roleId, err)
		}

		// 4.清除角色与接口的关联
		if err = l.svcCtx.SysRoleApiModel.DeleteRoleApiTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]接口关联失败：%v", roleId, err)
		}

		return err
	})
	if err != nil {
		l.Logger.Errorf("删除角色事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 同步casbin
	if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, result.Code); err != nil {
		l.Logger.Errorf("删除Casbin策略失败（角色编码: %s）：%v", result.Code, err)
	}

	return nil
}
