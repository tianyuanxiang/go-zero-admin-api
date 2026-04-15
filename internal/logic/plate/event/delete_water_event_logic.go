// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type DeleteWaterEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWaterEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWaterEventLogic {
	return &DeleteWaterEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWaterEventLogic) DeleteWaterEvent() error {
	// todo: add your logic here and delete this line

	return nil
}
