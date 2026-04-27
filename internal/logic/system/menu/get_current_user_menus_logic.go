// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/internal/common"
	"go-zero-admin/internal/middleware"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserMenusLogic {
	return &GetCurrentUserMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserMenusLogic) GetCurrentUserMenus() (resp *types.MenuTreeResp, err error) {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
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

	// 查询菜单权限
	menuIds, err := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, roleIds)
	if err != nil {
		l.Errorf("查询用户菜单失败：%v", err)
		menuIds = []int64{}
	}

	menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, menuIds)
	if err != nil {
		l.Errorf("查询菜单详情失败：%v", err)
		menus = []*systemmodel.SysMenu{}
	}

	menuTree := common.BuildMenuTree(menus, 0)

	return &types.MenuTreeResp{
		List: menuTree,
	}, nil
}
