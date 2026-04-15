// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitModelStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitModelStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitModelStateLogic {
	return &InitModelStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitModelStateLogic) InitModelState(req *types.InitModelStateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
