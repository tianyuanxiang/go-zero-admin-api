// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductionEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductionEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductionEventLogic {
	return &ListProductionEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductionEventLogic) ListProductionEvent(req *types.ListEventReq) (resp *types.ListProductionEventResp, err error) {
	// todo: add your logic here and delete this line

	return
}
