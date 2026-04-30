// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/internal/common"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuTreeLogic) GetMenuTree() (resp *types.MenuTreeResp, err error) {

	menus, err := getUserMenus(l.ctx, l.svcCtx, false)
	if err != nil {
		return nil, err
	}

	menuTree := common.BuildMenuTree(menus, 0)

	return &types.MenuTreeResp{List: menuTree}, nil
}

func (l *GetMenuTreeLogic) isAdmin(roleIds []int64) bool {
	for _, roleId := range roleIds {
		role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if err != nil || role == nil || role.Status != 1 {
			continue
		}
		if role.Code == "admin" {
			return true
		}
	}
	return false
}
