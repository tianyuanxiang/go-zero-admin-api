// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type ClearOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearOperLogLogic {
	return &ClearOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearOperLogLogic) ClearOperLog() error {
	// todo: add your logic here and delete this line

	return nil
}
