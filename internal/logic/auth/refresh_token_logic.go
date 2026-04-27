// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"go-zero-admin/pkg/jwtx"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	// 1. 解析刷新令牌
	claims, err := jwtx.ParseToken(req.RefreshToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		l.Logger.Infof("刷新令牌验证失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrTokenExpired)
	}

	// 2. 确认令牌类型必须是 RefreshToken
	if claims.TokenType != jwtx.TokenTypeRefresh {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrTokenInvalid, "请使用刷新令牌换取新的访问令牌")
	}

	// 3. 验证用户状态（防止用户被禁用后仍能刷新token）
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, claims.UserId)
	if err != nil {
		l.Logger.Errorf("查询用户[%d]失败：%v", claims.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	if user.Status != 1 {
		return nil, xerr.NewCodeError(xerr.ErrAccountDisabled)
	}

	// 4. 生成新的访问令牌
	newAccessToken, err := jwtx.GenerateToken(user.Id, user.Username, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Errorf("生成新AccessToken失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &types.RefreshTokenResp{
		AccessToken: newAccessToken,
		ExpiresIn:   l.svcCtx.Config.Auth.AccessExpire,
	}, nil
}
