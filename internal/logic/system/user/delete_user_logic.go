// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"go-zero-admin/internal/middleware"
	"go-zero-admin/internal/svc"
	"go-zero-admin/pkg/constants"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(userId int64) error {

	operatorId := middleware.GetUserIdFromCtx(l.ctx)
	if operatorId == 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if userId == operatorId {
		return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "禁止删除自己")
	}

	// 检查用户是否存在
	targetUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户[%d]失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if targetUser.DeletedAt.Valid {
		l.Errorf("用户[%d]已删除", userId)
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 超级管理员禁止删除：查询目标用户绑定的角色，命中 Code=="admin" 即拒绝
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户[%d]的角色列表失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if len(roleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
		if err != nil {
			l.Errorf("查询用户[%d]的角色明细失败：%v", userId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		for _, role := range roles {
			if role.Code == constants.RoleCodeAdmin {
				l.Infof("尝试删除超级管理员[%d]，已拦截", userId)
				return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "超级管理员禁止删除")
			}
		}
	}
	// 软删除用户（设置deleted_at）
	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err = l.svcCtx.SysUserModel.DeleteUserTrans(l.ctx, tx, userId); err != nil {
			l.Errorf("软删除用户[%d]失败：%v", userId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 清除用户的角色关联
		if err = l.svcCtx.SysUserRoleModel.DeleteByUserIdTrans(l.ctx, tx, userId); err != nil {
			l.Errorf("清除用户[%d]角色关联失败：%v", userId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})
	if err != nil {
		l.Logger.Errorf("删除用户事务执行失败: %v", err)
		return err
	}

	return nil
}
