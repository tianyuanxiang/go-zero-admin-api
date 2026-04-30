// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	casbinpkg "go-zero-admin/pkg/casbin"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) error {
	// 检查角色是否存在
	existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	if existRole.DeletedAt.Valid {
		return xerr.NewCodeError(xerr.ErrRoleNotFound)
	}

	// 如果修改了角色编码，检查新编码是否与其他角色冲突
	if req.RoleCode != nil && existRole.Code != *req.RoleCode {
		conflictRole, cErr := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, *req.RoleCode)
		if cErr != nil && cErr != sqlx.ErrNotFound {
			l.Errorf("检查角色编码冲突失败：%v", cErr)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if conflictRole != nil && conflictRole.Id != req.Id {
			return xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
		}
	}

	// 更新角色基本信息
	updates := make(map[string]interface{})
	if req.RoleName != nil {
		updates["name"] = *req.RoleName
	}
	if req.RoleCode != nil {
		updates["code"] = *req.RoleCode
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}

	if len(updates) == 0 {
		return nil
	}

	if err := l.svcCtx.SysRoleModel.UpdateRoleTrans(l.ctx, req.Id, updates); err != nil {
		l.Logger.Errorf("更新角色基本信息失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 如果编码从 operator 改成 op_admin，Casbin 里存的还是
	// operator，策略就失效了。所以如果允许修改角色编码，需要在事务后面加 Casbin 迁移
	if req.RoleCode != nil && existRole.Code != *req.RoleCode {
		oldPolicies, err := casbinpkg.GetRolePolicies(l.svcCtx.Enforcer, existRole.Code)
		if err != nil {
			l.Errorf("查询角色旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
		} else if len(oldPolicies) > 0 {
			// 提取 [path, method] 部分
			rules := make([][]string, 0, len(oldPolicies))
			for _, p := range oldPolicies {
				if len(p) >= 3 {
					rules = append(rules, []string{p[1], p[2]})
				}
			}
			// 删除旧编码策略
			if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, existRole.Code); err != nil {
				l.Errorf("清除角色旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
			}
			// 写入新编码策略
			if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, *req.RoleCode, rules); err != nil {
				l.Errorf("迁移casbin策略到新角色编码[%s]失败: %v", req.RoleCode, err)
			}
		}
	}

	return err
}
