// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDosingEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDosingEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDosingEventLogic {
	return &CreateDosingEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDosingEventLogic) CreateDosingEvent(req *types.CreateDosingEventReq) error {
	// todo: add your logic here and delete this line

	return nil
}
