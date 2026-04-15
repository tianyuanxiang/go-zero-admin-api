// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWaterEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWaterEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWaterEventLogic {
	return &CreateWaterEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWaterEventLogic) CreateWaterEvent(req *types.CreateWaterEventReq) error {
	// todo: add your logic here and delete this line

	return nil
}
