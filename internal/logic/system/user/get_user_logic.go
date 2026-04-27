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
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUser 根据用户ID查询用户详情（含角色列表）。
func (l *GetUserLogic) GetUser(userId int64) (resp *types.UserItem, err error) {
	// 查询用户基本信息
	userInfo, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Logger.Errorf("查询用户[%d]失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 查询用户角色
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户[%d]的角色失败：%v", userId, err)
		roleIds = []int64{}
	}

	// 4. 批量查询所有角色信息
	roleCodes := make([]string, 0, len(roleIds)) // roleId -> code
	if len(roleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
		if err != nil {
			l.Logger.Errorf("批量查询角色信息失败: %v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		for _, role := range roles {
			roleCodes = append(roleCodes, role.Code)
		}
	}

	return &types.UserItem{
		Id:        userInfo.Id,
		Username:  userInfo.Username,
		Nickname:  userInfo.Nickname,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		Status:    int(userInfo.Status),
		Remark:    userInfo.Remark,
		Roles:     roleCodes,
		CreatedAt: userInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
