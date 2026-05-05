package menu

import (
	"context"
	"go-zero-admin/internal/common"
	"go-zero-admin/internal/middleware"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/pkg/constants"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

func getUserMenus(ctx context.Context, svcCtx *svc.ServiceContext, filterHidden bool) ([]*systemmodel.SysMenu, error) {
	userId := middleware.GetUserIdFromCtx(ctx)
	if userId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	roleIds, err := svcCtx.SysUserRoleModel.GetRoleIdsByUserId(ctx, userId)
	if err != nil {
		logx.WithContext(ctx).Errorf("查询用户[%d]角色失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if len(roleIds) == 0 {
		logx.Errorf("用户角色ID为空")
		return nil, nil
	}

	var menus []*systemmodel.SysMenu
	if isAdmin(ctx, svcCtx, roleIds) {
		menus, err = svcCtx.SysMenuModel.ListAll(ctx)
	} else {
		menuIds, menuErr := svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(ctx, roleIds)
		if menuErr != nil {
			logx.WithContext(ctx).Errorf("查询角色菜单权限失败：%v", menuErr)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		menuIds, err = common.CompleteMenuAncestors(ctx, svcCtx.SysMenuModel, menuIds)
		if err != nil {
			logx.WithContext(ctx).Errorf("补齐菜单祖先链失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		menus, err = svcCtx.SysMenuModel.ListByIds(ctx, menuIds)
	}
	if err != nil {
		logx.WithContext(ctx).Errorf("查询菜单失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if filterHidden {
		visible := make([]*systemmodel.SysMenu, 0, len(menus))
		for _, m := range menus {
			if m.Visible == 1 {
				visible = append(visible, m)
			}
		}
		return visible, nil
	}
	return menus, nil
}

func isAdmin(ctx context.Context, svcCtx *svc.ServiceContext, roleIds []int64) bool {
	for _, roleId := range roleIds {
		role, err := svcCtx.SysRoleModel.FindOneByRoleId(ctx, roleId)
		if err != nil || role == nil || role.Status != 1 {
			continue
		}
		if role.Code == constants.RoleCodeAdmin {
			return true
		}
	}
	return false
}

func (l *UpdateMenuLogic) cascadeStatus(parentId, status int64, tx *gorm.DB) error {
	menus, err := l.svcCtx.SysMenuModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部菜单失败：%v", err)
		return err
	}

	menuIds := collectChildren(menus, parentId)
	if len(menuIds) == 0 {
		return nil
	}
	// 批量更新
	if err := l.svcCtx.SysMenuModel.BatchUpdateMenuStatus(l.ctx, tx, menuIds, int(status)); err != nil {
		return err
	}
	return nil
}

func (l *UpdateMenuLogic) cascadeVisible(parentId, visible int64, tx *gorm.DB) error {
	menus, err := l.svcCtx.SysMenuModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部菜单失败：%v", err)
		return err
	}

	menuIds := collectChildren(menus, parentId)
	if len(menuIds) == 0 {
		return nil
	}
	if err := l.svcCtx.SysMenuModel.BatchUpdateMenuVisible(l.ctx, tx, menuIds, int(visible)); err != nil {
		return err
	}
	return nil
}

func collectChildren(menus []*systemmodel.SysMenu, parentId int64) []int64 {
	var ids []int64
	for _, m := range menus {
		if m.ParentId == parentId && !m.DeletedAt.Valid {
			ids = append(ids, m.Id)
			ids = append(ids, collectChildren(menus, m.Id)...)
		}
	}
	return ids
}
