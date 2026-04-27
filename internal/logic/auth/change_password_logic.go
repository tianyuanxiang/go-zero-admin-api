// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"go-zero-admin/internal/middleware"
	"go-zero-admin/pkg/encrypt"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdatePassword 修改当前用户的登录密码。

/*
	业务流程：

//  1. 获取当前用户信息
//  2. 验证旧密码是否正确
//  3. 对新密码进行bcrypt加密
//  4. 更新数据库中的密码字段
//
// 参数：
//   - req : 修改密码请求体（旧密码 + 新密码）
//
// 返回：
//   - error : 业务错误
*/
func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)

	if userId == 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 1. 查询用户信息（需要获取当前密码Hash）
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Logger.Errorf("查询用户[%d]失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 2. 验证旧密码
	if !encrypt.CheckPassword(req.OldPassword, user.Password) {
		return xerr.NewCodeError(xerr.ErrOldPasswordWrong)
	}

	// 3. 对新密码进行bcrypt加密
	newHashedPassword, err := encrypt.HashPassword(req.NewPassword)
	if err != nil {
		l.Logger.Errorf("新密码加密失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 4. 更新密码
	if err = l.svcCtx.SysUserModel.UpdatePassword(l.ctx, userId, newHashedPassword); err != nil {
		l.Logger.Errorf("更新用户[%d]密码失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
