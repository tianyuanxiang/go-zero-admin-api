// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"plating/internal/middleware"
	systemmodel "plating/internal/model/system"
	"plating/pkg/xerr"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.UserInfoResp, err error) {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 查询用户基本信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户[%d]失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 查询用户角色
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户[%d]角色失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	roleCodes := make([]string, len(roleIds))
	for _, roleId := range roleIds {
		role, roleErr := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if roleErr != nil || role == nil || role.Status != 1 {
			continue
		}
		roleCodes = append(roleCodes, role.Code)
	}

	// 6. 查询用户的菜单权限（合并所有角色的菜单）
	menuIds, err := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, roleIds)
	if err != nil {
		l.Logger.Errorf("查询用户菜单失败：%v", err)
		menuIds = []int64{}
	}

	menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, menuIds)
	if err != nil {
		l.Logger.Errorf("查询菜单详情失败：%v", err)
		menus = []*systemmodel.SysMenu{}
	}

	// 将菜单列表构建为树形结构
	menuTree := buildMenuTree(menus, 0)

	return &types.UserInfoResp{
		UserInfo: types.UserInfo{
			UserId:   user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Roles:    roleCodes,
			Menus:    menuTree,
		},
	}, nil
}
