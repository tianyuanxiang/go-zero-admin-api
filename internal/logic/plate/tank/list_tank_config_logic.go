// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTankConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTankConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTankConfigLogic {
	return &ListTankConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTankConfigLogic) ListTankConfig(req *types.ListTankConfigReq) (resp *types.ListTankConfigResp, err error) {
	// todo: add your logic here and delete this line

	return
}
