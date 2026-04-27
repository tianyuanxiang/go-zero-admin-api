// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(menuId int64) error {

	_, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, menuId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("查询菜单[%d]失败：%v", menuId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 检查是否有子菜单
	hasChildren, err := l.svcCtx.SysMenuModel.HasChildren(l.ctx, menuId)
	if err != nil {
		l.Errorf("检查菜单[%d]的子节点失败：%v", menuId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if hasChildren {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "该菜单存在子菜单，请先删除子菜单")
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.软删除菜单
		if err := l.svcCtx.SysMenuModel.SoftDeleteTrans(l.ctx, tx, menuId); err != nil {
			l.Errorf("删除菜单[%d]失败：%v\n", menuId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 2.删除关联关系
		if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuByMenuIdTrans(l.ctx, tx, menuId); err != nil {
			l.Errorf("删除角色菜单关联关系失败：%v\n", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return err
	})
	if err != nil {
		l.Logger.Errorf("创建删除菜单事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return err
}
