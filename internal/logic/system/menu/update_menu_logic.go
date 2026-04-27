// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	_, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("查询菜单[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if err = l.svcCtx.SysMenuModel.Update(l.ctx, &systemmodel.SysMenu{
		Id:         req.Id,
		ParentId:   req.ParentId,
		Name:       req.MenuName,
		MenuPath:   req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		MenuType:   int64(req.MenuType),
		Sort:       int64(req.Sort),
		Status:     int64(req.Status),
		Permission: req.Perms,
		// Visible:    req.Visible,
		Remark: req.Remark,
	}); err != nil {
		l.Errorf("更新菜单[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
