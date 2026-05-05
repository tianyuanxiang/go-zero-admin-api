// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) error {
	oldMenu, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("查询菜单[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if oldMenu.DeletedAt.Valid {
		return xerr.NewCodeError(xerr.ErrMenuNotFound)
	}

	// 更新菜单基本信息
	updates := make(map[string]interface{})
	if req.ParentId != nil {
		updates["parent_id"] = *req.ParentId
	}
	if req.MenuName != nil {
		updates["name"] = *req.MenuName
	}
	if req.Path != nil {
		updates["menu_path"] = *req.Path
	}
	if req.Component != nil {
		updates["component"] = *req.Component
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.MenuType != nil {
		updates["menu_type"] = *req.MenuType
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Perms != nil {
		updates["permission"] = *req.Perms
	}
	// 可见性和可用性判断
	if req.Visible != nil {
		updates["visible"] = *req.Visible
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 1.开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.1 停用级联子孙
		if req.Status != nil && *req.Status == 0 && oldMenu.Status == 1 {
			if err := l.cascadeStatus(req.Id, 0, tx); err != nil {
				l.Errorf("批量设置级联子孙关闭状态失败%v", err)
				return err
			}
		}
		// 1.2 开启级联子孙
		if req.Status != nil && *req.Status == 1 && oldMenu.Status == 0 {
			if err := l.cascadeStatus(req.Id, 1, tx); err != nil {
				l.Errorf("批量设置级联子孙开启状态失败%v", err)
				return err
			}
		}

		// 2.1 隐藏级联子孙
		if req.Visible != nil && *req.Visible == 0 && oldMenu.Visible == 1 && oldMenu.MenuType != 2 {
			if err := l.cascadeVisible(req.Id, 0, tx); err != nil {
				l.Errorf("批量隐藏级联子孙失败%v", err)
				return err
			}
		}
		// 2.2 可见级联子孙
		if req.Visible != nil && *req.Visible == 1 && oldMenu.Visible == 0 && oldMenu.MenuType != 2 {
			if err := l.cascadeVisible(req.Id, 1, tx); err != nil {
				l.Errorf("批量可见级联子孙失败%v", err)
				return err
			}
		}

		if err = l.svcCtx.SysMenuModel.UpdateMenuTrans(l.ctx, tx, req.Id, updates); err != nil {
			l.Errorf("更新菜单[%d]失败：%v", req.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})

	if err != nil {
		l.Logger.Errorf("更新菜单事务执行失败: %v", err)
		return err
	}

	return nil
}
