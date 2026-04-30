// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/pkg/xerr"

	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) error {

	// 确认父菜单是否存在
	if req.ParentId > 0 {
		menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, []int64{req.ParentId})
		if err != nil {
			l.Errorf("查询父菜单[%d]失败：%v", req.ParentId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(menus) == 0 {
			return xerr.NewCodeErrorMsg(xerr.ErrMenuNotFound, "父菜单不存在或已被删除")
		}
		// 额外校验：父菜单不能是按钮类型
		if menus[0].MenuType == 2 {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "按钮类型的菜单不能作为父菜单")
		}
	}
	_, err := l.svcCtx.SysMenuModel.Insert(l.ctx, &systemmodel.SysMenu{
		ParentId:   req.ParentId,
		Name:       req.MenuName,
		MenuType:   req.MenuType,
		MenuPath:   req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       int64(req.Sort),
		Permission: req.Perms,
		Status:     int64(req.Status),
		Remark:     req.Remark,
	})
	if err != nil {
		l.Errorf("插入菜单记录失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return nil
}
