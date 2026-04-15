// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type DeleteProductionEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductionEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductionEventLogic {
	return &DeleteProductionEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductionEventLogic) DeleteProductionEvent() error {
	// todo: add your logic here and delete this line

	return nil
}
