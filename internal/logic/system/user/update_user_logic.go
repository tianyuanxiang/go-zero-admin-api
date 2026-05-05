// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) error {
	// 检查用户是否存在
	sysUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户失败，用户ID[%d]：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if sysUser.DeletedAt.Valid {
		l.Errorf("用户[%d]已删除", req.Id)
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 校验用户名
	if req.Username != "" {
		userByName, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, req.Username)
		if err != nil {
			l.Errorf("更新用户时校验用户名[%s]失败：%v", req.Username, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}

		if userByName != nil && userByName.Id != req.Id {
			l.Errorf("用户[%s]已存在", req.Username)
			return xerr.NewCodeError(xerr.ErrParamInvalid)
		}
	}
	// 更新字典类型信息
	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	// 是否需要触达角色绑定逻辑
	needUpdateRoles := req.RoleIds != nil

	if len(updates) == 0 && !needUpdateRoles {
		l.Error("更新字段为空")
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 角色合法性校验
	if needUpdateRoles && len(req.RoleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, req.RoleIds)
		if err != nil {
			l.Errorf("角色合法性校验失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(roles) != len(req.RoleIds) { // 严格匹配，避免传无效ID被忽略
			return xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.SysUserModel.UpdateTrans(l.ctx, tx, req.Id, updates); err != nil {
			l.Errorf("更新用户[%d]失败：%v", req.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if needUpdateRoles {
			if err := l.svcCtx.SysUserRoleModel.AssignRolesTrans(l.ctx, tx, req.Id, req.RoleIds); err != nil {
				l.Errorf("更新用户[%d]分配角色失败：%v", req.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}

		}
		return nil
	})

	if err != nil {
		l.Logger.Errorf("更新用户事务执行失败: %v", err)
		return err
	}
	return nil
}
