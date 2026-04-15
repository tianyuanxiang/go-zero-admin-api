// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"plating/internal/middleware"
	"plating/pkg/encrypt"
	"plating/pkg/xerr"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) error {
	// 获取当前用户userId
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.Id == userId {
		l.Logger.Infof("该接口禁止修改自己的密码!")
		return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "请通过修改密码功能操作")
	}

	userInfo, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Logger.Errorf("查询用户[%d]失败：%v", userId, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if userInfo.DeletedAt.Valid {
		l.Logger.Infof("禁止修改已删除用户的密码!")
		return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "该用户已被删除")
	}

	// 密码加密
	hashedPassword, err := encrypt.HashPassword(req.NewPassword)
	if err != nil {
		l.Logger.Errorf("新密码加密失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 修改密码
	if err := l.svcCtx.SysUserModel.UpdatePassword(l.ctx, userId, hashedPassword); err != nil {
		l.Logger.Errorf("修改密码失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
