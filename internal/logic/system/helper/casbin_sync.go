package helper

import (
	"context"
	"go-zero-admin/internal/svc"
	casbinpkg "go-zero-admin/pkg/casbin"

	"github.com/zeromicro/go-zero/core/logx"
)

// RebuildCasbinByRoleIds 根据角色ID列表，重建每个角色的Casbin策略。
//
// 适用场景：API的 path/method 变更、API被删除后，需要同步更新所有关联角色的Casbin规则。
// 该方法只记录错误日志，不返回错误（降级处理，不影响主流程）。
func RebuildCasbinByRoleIds(ctx context.Context, svcCtx *svc.ServiceContext, roleIds []int64) {
	if len(roleIds) == 0 {
		return
	}

	logger := logx.WithContext(ctx)

	// 1.批量查询角色信息（拿roleCode）
	roles, err := svcCtx.SysRoleModel.FindByIds(ctx, roleIds)
	if err != nil {
		logger.Errorf("重建Casbin策略-查询角色信息失败：%v", err)
		return
	}

	// 2.逐个角色重建策略
	for _, role := range roles {
		// 2.1 查询该角色绑定的所有接口ID
		apiIds, err := svcCtx.SysRoleApiModel.ListApiIdsByRoleId(ctx, role.Id)
		if err != nil {
			logger.Errorf("重建Casbin策略-查询角色[%s]绑定接口失败：%v", role.Code, err)
			continue
		}

		// 2.2 该角色没有绑定任何接口，清空策略
		if len(apiIds) == 0 {
			if err := casbinpkg.RemoveAllPoliciesForRole(svcCtx.Enforcer, role.Code); err != nil {
				logger.Errorf("重建Casbin策略-清除角色[%s]策略失败：%v", role.Code, err)
			}
			continue
		}

		// 2.3 查询接口详情（拿path和method）
		apis, err := svcCtx.SysApiModel.ListByIds(ctx, apiIds)
		if err != nil {
			logger.Errorf("重建Casbin策略-查询接口列表失败：%v", err)
			continue
		}

		// 2.4 构建规则并全量覆盖
		rules := make([][]string, 0, len(apis))
		for _, api := range apis {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(svcCtx.Enforcer, role.Code, rules); err != nil {
			logger.Errorf("重建Casbin策略-角色[%s]写入失败：%v", role.Code, err)
		}
	}
}
