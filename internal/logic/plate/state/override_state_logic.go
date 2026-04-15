// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package state

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OverrideStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOverrideStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OverrideStateLogic {
	return &OverrideStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OverrideStateLogic) OverrideState(req *types.OverrideStateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
