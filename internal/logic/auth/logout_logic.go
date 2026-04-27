// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"go-zero-admin/internal/middleware"
	"go-zero-admin/pkg/jwtx"
	"time"

	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(tokenStr string) error {
	if tokenStr == "" {
		return nil
	}

	// 解析Token获取过期时间（即使Token即将过期也要加入黑名单）
	expireAt := jwtx.GetTokenExpireAt(tokenStr, l.svcCtx.Config.Auth.AccessSecret)
	if expireAt == 0 {
		// Token已过期或无效，无需加入黑名单
		return nil
	}

	// 计算剩余有效时间（秒）
	remainDuration := time.Until(time.Unix(expireAt, 0))
	if remainDuration <= 0 {
		// 已过期，无需处理
		return nil
	}

	// 将Token加入Redis黑名单（键：token:blacklist:<token>，TTL=剩余有效期）
	blacklistKey := fmt.Sprintf("%s%s", middleware.RedisTokenBlacklistPrefix, tokenStr)
	if err := l.svcCtx.RDB.Set(l.ctx, blacklistKey, "1", remainDuration).Err(); err != nil {
		l.Logger.Errorf("Token加入Redis黑名单失败：%v", err)
		// 黑名单写入失败不影响登出响应，日志记录即可
		return nil
	}
	l.Infof("用户[%d]登出成功，Token已加入黑名单", middleware.GetUserIdFromCtx(l.ctx))
	return nil
}
