// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"plating/internal/svc"
)

type DeleteTankConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTankConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTankConfigLogic {
	return &DeleteTankConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTankConfigLogic) DeleteTankConfig() error {
	// todo: add your logic here and delete this line

	return nil
}
