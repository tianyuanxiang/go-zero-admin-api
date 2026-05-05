// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/internal/common"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/constants"

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
	menus, err := getUserMenus(l.ctx, l.svcCtx, true)
	if err != nil {
		return nil, err
	}

	menuTree := common.BuildMenuTree(menus, 0)

	return &types.MenuTreeResp{
		List: menuTree,
	}, nil
}

func (l *GetCurrentUserMenusLogic) isAdmin(roleIds []int64) bool {
	for _, roleId := range roleIds {
		role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if err != nil || role == nil || role.Status != 1 {
			continue
		}
		if role.Code == constants.RoleCodeAdmin {
			return true
		}
	}
	return false
}
