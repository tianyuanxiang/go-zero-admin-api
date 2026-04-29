// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRole 根据角色ID查询角色详情。
/* 参数：
//   - roleId : 角色ID
//
// 返回：
//   - *types.RoleItem : 角色详情
//   - error           : 业务错误
*/

func (l *GetRoleLogic) GetRole(roleId int64) (resp *types.RoleItem, err error) {
	role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", roleId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	return &types.RoleItem{
		Id:        role.Id,
		RoleCode:  role.Code,
		RoleName:  role.Name,
		Status:    int(role.Status),
		Sort:      int(role.Sort),
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
	}, err
}
