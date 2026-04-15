// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApiLogic) ListApi(req *types.ListApiReq) (resp *types.ListApiResp, err error) {
	// todo: add your logic here and delete this line

	return
}
