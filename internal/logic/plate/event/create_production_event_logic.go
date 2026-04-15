// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductionEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductionEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductionEventLogic {
	return &CreateProductionEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductionEventLogic) CreateProductionEvent(req *types.CreateProductionEventReq) error {
	// todo: add your logic here and delete this line

	return nil
}
