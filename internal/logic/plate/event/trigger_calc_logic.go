// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerCalcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerCalcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerCalcLogic {
	return &TriggerCalcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerCalcLogic) TriggerCalc(req *types.TriggerCalcReq) (resp *types.TriggerCalcResp, err error) {
	// todo: add your logic here and delete this line

	return
}
