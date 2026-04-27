// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"database/sql"
	"go-zero-admin/pkg/xerr"
	"time"

	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(userId int64) error {
	// 检查用户是否存在
	_, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户[%d]失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 软删除用户（设置deleted_at）
	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err = l.svcCtx.SysUserModel.Update(l.ctx, &systemmodel.SysUser{
			Id: userId,
			DeletedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}); err != nil {
			l.Errorf("软删除用户[%d]失败：%v", userId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 清除用户的角色关联
		if err = l.svcCtx.SysUserRoleModel.Delete(l.ctx, userId); err != nil {
			l.Errorf("清除用户[%d]角色关联失败：%v", userId, err)
			// 不影响主流程，记录日志即可
		}
		return nil
	})
	if err != nil {
		l.Logger.Errorf("删除用户事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
