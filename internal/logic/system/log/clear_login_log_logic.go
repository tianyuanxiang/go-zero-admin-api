// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type ClearLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearLoginLogLogic {
	return &ClearLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearLoginLogLogic) ClearLoginLog() error {
	// todo: add your logic here and delete this line

	return nil
}
