// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllRoleLogic {
	return &ListAllRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllRoleLogic) ListAllRole() (resp *types.ListRoleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
