// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDosingEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDosingEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDosingEventLogic {
	return &ListDosingEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDosingEventLogic) ListDosingEvent(req *types.ListEventReq) (resp *types.ListDosingEventResp, err error) {
	// todo: add your logic here and delete this line

	return
}
