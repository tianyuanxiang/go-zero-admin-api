// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByTypeLogic {
	return &GetDictDataByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictDataByTypeLogic) GetDictDataByType() (resp *types.ListDictDataResp, err error) {
	// todo: add your logic here and delete this line

	return
}
