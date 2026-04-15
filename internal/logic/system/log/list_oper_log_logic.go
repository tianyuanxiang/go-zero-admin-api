// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperLogLogic {
	return &ListOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOperLogLogic) ListOperLog(req *types.ListOperLogReq) (resp *types.ListOperLogResp, err error) {
	// todo: add your logic here and delete this line

	return
}
