// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

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

func (l *ListAllRoleLogic) ListAllRole() (resp *types.ListAllResp, err error) {
	roles, err := l.svcCtx.SysRoleModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部角色失败: %v\n", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	list := make([]types.RoleOption, 0, len(roles))
	for _, role := range roles {
		list = append(list, types.RoleOption{
			Id:       role.Id,
			RoleName: role.Name,
			RoleCode: role.Code,
		})
	}

	return &types.ListAllResp{
		ListAll: list,
	}, err
}
