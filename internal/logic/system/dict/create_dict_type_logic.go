// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.CreateDictTypeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
