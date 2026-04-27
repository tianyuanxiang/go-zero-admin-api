// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
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
