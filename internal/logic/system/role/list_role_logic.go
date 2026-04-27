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

type ListRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRoleLogic) ListRole(req *types.ListRoleReq) (resp *types.ListRoleResp, err error) {
	roles, count, err := l.svcCtx.SysRoleModel.List(l.ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		l.Errorf("查询角色列表失败 %v\n", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	list := make([]types.RoleItem, 0, len(roles))
	for _, role := range roles {
		list = append(list, types.RoleItem{
			Id:        role.Id,
			RoleName:  role.Name,
			RoleCode:  role.Code,
			Status:    int(role.Status),
			Sort:      int(role.Sort),
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &types.ListRoleResp{
		Total: count,
		List:  list,
	}, err
}
