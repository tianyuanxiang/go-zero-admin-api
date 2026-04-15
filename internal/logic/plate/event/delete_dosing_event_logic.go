// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type DeleteDosingEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDosingEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDosingEventLogic {
	return &DeleteDosingEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDosingEventLogic) DeleteDosingEvent() error {
	// todo: add your logic here and delete this line

	return nil
}
