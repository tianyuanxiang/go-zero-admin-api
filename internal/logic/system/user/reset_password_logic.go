// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

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

// ResetPassword 管理员重置指定用户的密码。
//
// 业务流程：
//  1. 校验操作者已登录（从 ctx 提取 operatorId）
//  2. 校验目标用户ID合法性
//  3. 拒绝通过该接口修改操作者自身密码（必须走"修改密码"接口验证旧密码）
//  4. 校验目标用户存在且未被软删除
//  5. 对新密码进行 bcrypt 加密
//  6. 更新目标用户的密码字段
//  7. 记录高敏操作审计日志
//
// 参数：
//   - req : 重置密码请求体，Id 来自路径，NewPassword 来自请求体
//
// 返回：
//   - error : 业务错误

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) error {
	// 获取当前用户userId
	operatorId := middleware.GetUserIdFromCtx(l.ctx)
	if operatorId == 0 {
		return xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	// 目标用户ID合法性兜底校验（handler 已校验，此处防御性兜底）
	if req.Id <= 0 {
		return xerr.NewCodeError(xerr.ErrParamInvalid)
	}
	// 禁止通过本接口重置自己的密码，必须引导至"修改密码"功能（需校验旧密码）
	if req.Id == operatorId {
		l.Infof("操作者[%d]尝试通过重置接口修改自身密码，已拦截", operatorId)
		return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "请通过修改密码功能操作")
	}

	// 查询目标用户是否存在且未被软删除
	targetUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("操作者[%d]查询目标用户[%d]失败：%v", operatorId, req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if targetUser.DeletedAt.Valid {
		l.Infof("操作者[%d]尝试重置已删除用户[%d]的密码，已拦截", operatorId, req.Id)
		return xerr.NewCodeErrorMsg(xerr.ErrForbidden, "该用户已被删除")
	}

	// 新密码 bcrypt 加密
	hashedPassword, err := encrypt.HashPassword(req.NewPassword)
	if err != nil {
		l.Errorf("新密码加密失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 更新目标用户的密码（注意此处必须用 req.Id，不能用 operatorId）
	if err := l.svcCtx.SysUserModel.UpdatePassword(l.ctx, req.Id, hashedPassword); err != nil {
		l.Errorf("操作者[%d]重置用户[%d]密码失败：%v", operatorId, req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 高敏操作审计日志（密码重置必须留痕，便于事后追溯）
	l.Infof("操作者[%d]成功重置用户[%d]的密码", operatorId, req.Id)

	return nil
}
