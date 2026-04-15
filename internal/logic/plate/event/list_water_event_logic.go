// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWaterEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWaterEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWaterEventLogic {
	return &ListWaterEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWaterEventLogic) ListWaterEvent(req *types.ListEventReq) (resp *types.ListWaterEventResp, err error) {
	// todo: add your logic here and delete this line

	return
}
