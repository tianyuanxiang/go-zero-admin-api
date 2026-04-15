// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTankConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTankConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTankConfigLogic {
	return &GetTankConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTankConfigLogic) GetTankConfig() (resp *types.TankConfigItem, err error) {
	// todo: add your logic here and delete this line

	return
}
