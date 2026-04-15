// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package state

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLatestStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestStateLogic {
	return &GetLatestStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLatestStateLogic) GetLatestState() (resp *types.ModelStateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
