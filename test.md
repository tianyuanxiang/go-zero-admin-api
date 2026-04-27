第一步：Model 层需要改造

你的 UpdateRoleTrans 当前只能软删除，不能更新角色字段。需要改方法签名，让它接收一个 SysRole 对象来做通用更新。

文件：internal/model/system/sys_role_model.go

接口定义改为：

// 把原来的 UpdateRoleTrans 签名改掉
UpdateRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) error

实现改为：

// UpdateRoleTrans 在事务中更新角色基本信息（不含关联表）。
func (m *customSysRoleModel) UpdateRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) error {
result := tx.WithContext(ctx).Table("sys_role").
Where("id = ? AND deleted_at IS NULL", role.Id).
Updates(map[string]interface{}{
"name":   role.Name,
"code":   role.Code,
"remark": role.Remark,
"sort":   role.Sort,
})
return result.Error
}

原来那个软删除的功能，建议单独起一个方法名叫 SoftDeleteRoleTrans，别和通用更新混在一起：

// SoftDeleteRoleTrans 在事务中软删除角色。
func (m *customSysRoleModel) SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
result := tx.WithContext(ctx).Table("sys_role").
Where("id = ?", roleId).
Update("deleted_at", sql.NullTime{
Time:  time.Now(),
Valid: true,
})
return result.Error
}

接口定义也要加上：

SysRoleModel interface {
sysRoleModel
FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error)
FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error)
InsertRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) (int64, error)
UpdateRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) error
SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
}

  ---
第二步：UpdateRole 完整代码

文件：internal/logic/system/role/update_role_logic.go

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
"context"
systemmodel "plating/internal/model/system"
"plating/internal/svc"
"plating/internal/types"
casbinpkg "plating/pkg/casbin"
"plating/pkg/xerr"

        "github.com/zeromicro/go-zero/core/logx"
        "github.com/zeromicro/go-zero/core/stores/sqlx"
        "gorm.io/gorm"
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
// 1.检查角色是否存在
existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
if err != nil {
if err == sqlx.ErrNotFound {
return xerr.NewCodeError(xerr.ErrRoleNotFound)
}
l.Errorf("查询角色[%d]失败：%v", req.Id, err)
return xerr.NewCodeError(xerr.ErrInternal)
}

        // 2.如果修改了角色编码，检查新编码是否与其他角色冲突
        oldRoleCode := existRole.Code
        if oldRoleCode != req.RoleCode {
                conflictRole, cErr := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, req.RoleCode)
                if cErr != nil && cErr != sqlx.ErrNotFound {
                        l.Errorf("检查角色编码冲突失败：%v", cErr)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if conflictRole != nil && conflictRole.Id != req.Id {
                        return xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
                }
        }

        // 3.获取接口信息 && 确认接口数量完整
        var apiModels []systemmodel.SysApi
        if len(req.ApiIds) > 0 {
                apiModels, err = l.svcCtx.SysApiModel.ListByIds(l.ctx, req.ApiIds)
                if err != nil {
                        l.Logger.Errorf("查询接口信息失败：%v", err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(apiModels) != len(req.ApiIds) {
                        return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分api接口信息不存在")
                }
        }

        // 4.开启事务
        err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
                // 4.1 更新角色基本信息
                if err := l.svcCtx.SysRoleModel.UpdateRoleTrans(l.ctx, tx, &systemmodel.SysRole{
                        Id:     req.Id,
                        Name:   req.RoleName,
                        Code:   req.RoleCode,
                        Remark: req.Remark,
                }); err != nil {
                        l.Logger.Errorf("更新角色基本信息失败：%v", err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }

                // 4.2 先删后插：更新关联菜单
                if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuTrans(l.ctx, tx, req.Id); err != nil {
                        l.Logger.Errorf("删除角色[%d]旧菜单关联失败：%v", req.Id, err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(req.MenuIds) > 0 {
                        roleMenus := make([]systemmodel.SysRoleMenu, 0, len(req.MenuIds))
                        for _, menuId := range req.MenuIds {
                                roleMenus = append(roleMenus, systemmodel.SysRoleMenu{
                                        RoleId: req.Id,
                                        MenuId: menuId,
                                })
                        }
                        if _, err := l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus); err != nil
                                l.Logger.Errorf("插入角色[%d]新菜单关联失败：%v", req.Id, err)
                                return xerr.NewCodeError(xerr.ErrInternal)
                        }
                }

                // 4.3 先删后插：更新关联接口
                if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiTrans(l.ctx, tx, req.Id); err != nil {
                        l.Logger.Errorf("删除角色[%d]旧接口关联失败：%v", req.Id, err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(req.ApiIds) > 0 {
                        roleApis := make([]systemmodel.SysRoleApi, 0, len(req.ApiIds))
                        for _, apiId := range req.ApiIds {
                                roleApis = append(roleApis, systemmodel.SysRoleApi{
                                        RoleId: req.Id,
                                        ApiId:  apiId,
                                })
                        }
                        if _, err := l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis); err != nil {
                                l.Logger.Errorf("插入角色[%d]新接口关联失败：%v", req.Id, err)
                                return xerr.NewCodeError(xerr.ErrInternal)
                        }
                }

                return nil
        })
        if err != nil {
                l.Logger.Errorf("更新角色事务执行失败: %v", err)
                return xerr.NewCodeError(xerr.ErrInternal)
        }

        // 5.同步Casbin策略（事务成功后执行）
        // 如果角色编码变更了，先清除旧编码的策略
        if oldRoleCode != req.RoleCode {
                if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, oldRoleCode); err != nil {
                        l.Logger.Errorf("清除旧角色编码[%s]Casbin策略失败：%v", oldRoleCode, err)
                }
        }
        // 用新编码全量覆盖策略
        if len(apiModels) > 0 {
                rules := make([][]string, 0, len(apiModels))
                for _, api := range apiModels {
                        rules = append(rules, []string{api.ApiPath, api.Method})
                }
                if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, req.RoleCode, rules); err != nil {
                        l.Logger.Errorf("同步Casbin策略失败（角色编码: %s）：%v", req.RoleCode, err)
                }
        } else {
                // API列表为空，清除该角色所有策略
                if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, req.RoleCode); err != nil {
                        l.Logger.Errorf("清除Casbin策略失败（角色编码: %s）：%v", req.RoleCode, err)
                }
        }

        return nil
}

  ---
重点说明

1. 角色编码变更的特殊处理（第5步）

Casbin 的 subject 是 roleCode，不是 roleId。如果用户把编码从 operator 改成 op_admin，那 casbin_rule 表里存的还是
operator。所以要：
- 先用旧编码 RemoveAllPoliciesForRole 清掉
- 再用新编码 AddRolePolicies 写入

AddRolePolicies 内部虽然也有清除逻辑，但它清的是新编码的规则——如果编码变了，旧编码的规则它清不到。

2. 关联表的"先删后插"策略

菜单和接口的关联关系采用"全量删除 + 全量插入"，而不是做 diff。原因：
- 逻辑简单，不容易出错
- 关联表数据量小，性能影响可以忽略
- 和你 CreateRole 的写法风格保持一致

3. UpdateRoleReq 里没有 Sort 字段

你的 UpdateRoleReq 定义里没有 Sort，所以更新时我没传 Sort。如果需要支持修改排序，你要在 UpdateRoleReq 里加上 Sort
int64，然后在 UpdateRoleTrans 的 Updates map 里也加上。



✻ Brewed for 1m 55s

❯ 为什么你的代码要"先删后插：更新关联菜单"呢？这样的话前端如果只是更新remark，也要把所有的关联菜单传过来？

● 你说得对，这个设计确实有问题。

核心矛盾在于：前端每次只想改一个字段，却被迫要把所有关联关系都传回来，否则就会被清空。

这取决于你想怎么设计这个接口，有两种思路：

思路一：拆接口（推荐）

把"更新角色基本信息"和"更新角色关联关系"拆成独立的接口：

┌───────────────────────────┬────────────────┬──────────────────────────────────┐
│           接口            │      职责      │               参数               │
├───────────────────────────┼────────────────┼──────────────────────────────────┤
│ PUT /role/:id             │ 只更新基本信息 │ roleName, roleCode, remark, sort │
├───────────────────────────┼────────────────┼──────────────────────────────────┤
│ PUT /role/:id/permissions │ 只更新关联关系 │ menuIds, apiIds                  │
└───────────────────────────┴────────────────┴──────────────────────────────────┘

好处：
- 职责单一，前端改备注不需要传菜单和接口
- Casbin 同步只在 permissions 接口里做
- 前端调用意图明确，不容易误操作

思路二：保持一个接口，用字段判断是否需要更新关联

前端传 null/不传 表示"不改"，传 []（空数组）表示"清空"，传 [1,2,3] 表示"覆盖"。

但这个方案在 Go 里区分"没传"和"传了空数组"比较麻烦，需要用 *[]int64 指针类型，增加复杂度。

  ---
你倾向哪种？ 我按你选的方案给你出具体代码。

❯ 拆接口（推荐）

● 先看一下当前的 API 定义和路由结构。

Searched for 2 patterns, read 3 files (ctrl+o to expand)

● 再看看 handler 里的 response 引入方式，确认 import 风格。

Read 1 file (ctrl+o to expand)

● 好了，所有上下文都齐了。下面是完整的改造方案，你需要改动的文件逐个列出来。

  ---
需要改动的文件清单

┌──────┬─────────────────────────────────────────────────────────────────┬─────────────────────────────────────┐
│ 序号 │                              文件                               │                动作                 │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 1    │ api/desc/system/role.api                                        │ 修改 UpdateRoleReq，新增            │
│      │                                                                 │ UpdateRolePermissionsReq 和路由     │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 2    │ internal/types/system_role.go                                   │ 重新生成或手动改（跟着 api 文件走） │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 3    │ internal/handler/system/role/update_role_handler.go             │ 改：去掉菜单/接口逻辑               │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 4    │ internal/handler/system/role/update_role_permissions_handler.go │ 新增                                │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 5    │ internal/logic/system/role/update_role_logic.go                 │ 改：只更新基本信息                  │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 6    │ internal/logic/system/role/update_role_permissions_logic.go     │ 新增                                │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 7    │ internal/model/system/sys_role_model.go                         │ 改：拆分 UpdateRoleTrans            │
├──────┼─────────────────────────────────────────────────────────────────┼─────────────────────────────────────┤
│ 8    │ 路由注册文件                                                    │ 注册新 handler（goctl               │
│      │                                                                 │ 生成或手动加）                      │
└──────┴─────────────────────────────────────────────────────────────────┴─────────────────────────────────────┘

  ---
1. API 定义

文件：api/desc/system/role.api

UpdateRoleReq 改为只保留基本信息字段，新增 UpdateRolePermissionsReq：

// 更新角色基本信息请求体
type UpdateRoleReq {
Id       int64  `path:"id"`               // 角色ID
RoleCode string `json:"roleCode"`         // 角色编码
RoleName string `json:"roleName"`         // 角色名称
Sort     int64  `json:"sort"`             // 排序值
Remark   string `json:"remark,optional"`  // 备注
}

// 更新角色权限（菜单+接口）请求体
type UpdateRolePermissionsReq {
Id      int64   `path:"id"`               // 角色ID
MenuIds []int64 `json:"menuIds,optional"` // 关联菜单ID列表
ApiIds  []int64 `json:"apiIds,optional"`  // 关联接口ID列表
}

路由部分新增一条：

// 更新角色基本信息
@handler UpdateRoleHandler
put /role/:id (UpdateRoleReq)

// 更新角色权限（菜单+接口绑定）
@handler UpdateRolePermissionsHandler
put /role/:id/permissions (UpdateRolePermissionsReq)

  ---
2. Model 层改造

文件：internal/model/system/sys_role_model.go

接口定义：

SysRoleModel interface {
sysRoleModel
FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error)
FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error)
InsertRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) (int64, error)
UpdateRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) error
SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
}

实现部分，把原来的 UpdateRoleTrans 改成真正的更新，软删除独立出去：

// UpdateRoleTrans 在事务中更新角色基本信息。
func (m *customSysRoleModel) UpdateRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) error {
result := tx.WithContext(ctx).Table("sys_role").
Where("id = ? AND deleted_at IS NULL", role.Id).
Updates(map[string]interface{}{
"name":   role.Name,
"code":   role.Code,
"remark": role.Remark,
"sort":   role.Sort,
})
return result.Error
}

// SoftDeleteRoleTrans 在事务中软删除角色。
func (m *customSysRoleModel) SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
result := tx.WithContext(ctx).Table("sys_role").
Where("id = ?", roleId).
Update("deleted_at", sql.NullTime{
Time:  time.Now(),
Valid: true,
})
return result.Error
}

  ---
3. UpdateRole Logic — 只管基本信息

文件：internal/logic/system/role/update_role_logic.go

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
"context"
systemmodel "plating/internal/model/system"
"plating/internal/svc"
"plating/internal/types"
"plating/pkg/xerr"

        "github.com/zeromicro/go-zero/core/logx"
        "github.com/zeromicro/go-zero/core/stores/sqlx"
        "gorm.io/gorm"
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
// 1.检查角色是否存在
existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
if err != nil {
if err == sqlx.ErrNotFound {
return xerr.NewCodeError(xerr.ErrRoleNotFound)
}
l.Errorf("查询角色[%d]失败：%v", req.Id, err)
return xerr.NewCodeError(xerr.ErrInternal)
}

        // 2.如果修改了角色编码，检查新编码是否与其他角色冲突
        if existRole.Code != req.RoleCode {
                conflictRole, cErr := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, req.RoleCode)
                if cErr != nil && cErr != sqlx.ErrNotFound {
                        l.Errorf("检查角色编码冲突失败：%v", cErr)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if conflictRole != nil && conflictRole.Id != req.Id {
                        return xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
                }
        }

        // 3.更新角色基本信息
        err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
                if err := l.svcCtx.SysRoleModel.UpdateRoleTrans(l.ctx, tx, &systemmodel.SysRole{
                        Id:     req.Id,
                        Name:   req.RoleName,
                        Code:   req.RoleCode,
                        Remark: req.Remark,
                        Sort:   req.Sort,
                }); err != nil {
                        l.Logger.Errorf("更新角色基本信息失败：%v", err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                return nil
        })
        if err != nil {
                l.Logger.Errorf("更新角色事务执行失败: %v", err)
                return xerr.NewCodeError(xerr.ErrInternal)
        }

        return nil
}

注意：这里不涉及 Casbin，因为只改了基本信息，没动接口绑定关系。

但有一个边界情况——角色编码变更。如果编码从 operator 改成 op_admin，Casbin 里存的还是
operator，策略就失效了。所以如果你允许修改角色编码，需要在事务后面加 Casbin 迁移：

// 4.如果角色编码变更，迁移Casbin策略到新编码
if existRole.Code != req.RoleCode {
oldPolicies, err := casbinpkg.GetRolePolicies(l.svcCtx.Enforcer, existRole.Code)
if err != nil {
l.Logger.Errorf("查询旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
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
l.Logger.Errorf("清除旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
}
// 写入新编码策略
if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, req.RoleCode, rules); err != nil {
l.Logger.Errorf("迁移Casbin策略到新编码[%s]失败：%v", req.RoleCode, err)
}
}
}

如果你不允许修改角色编码（把 RoleCode 从 UpdateRoleReq 里去掉），那这段就不需要，整个 UpdateRole 就完全不碰
Casbin。

  ---
4. UpdateRolePermissions Logic — 只管权限绑定

新文件：internal/logic/system/role/update_role_permissions_logic.go

package role

import (
"context"
systemmodel "plating/internal/model/system"
"plating/internal/svc"
"plating/internal/types"
casbinpkg "plating/pkg/casbin"
"plating/pkg/xerr"

        "github.com/zeromicro/go-zero/core/logx"
        "github.com/zeromicro/go-zero/core/stores/sqlx"
        "gorm.io/gorm"
)

type UpdateRolePermissionsLogic struct {
logx.Logger
ctx    context.Context
svcCtx *svc.ServiceContext
}

func NewUpdateRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermissionsLogic {
return &UpdateRolePermissionsLogic{
Logger: logx.WithContext(ctx),
ctx:    ctx,
svcCtx: svcCtx,
}
}

func (l *UpdateRolePermissionsLogic) UpdateRolePermissions(req *types.UpdateRolePermissionsReq) error {
// 1.检查角色是否存在
existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
if err != nil {
if err == sqlx.ErrNotFound {
return xerr.NewCodeError(xerr.ErrRoleNotFound)
}
l.Errorf("查询角色[%d]失败：%v", req.Id, err)
return xerr.NewCodeError(xerr.ErrInternal)
}

        // 2.获取接口信息并确认数量完整
        var apiModels []systemmodel.SysApi
        if len(req.ApiIds) > 0 {
                apiModels, err = l.svcCtx.SysApiModel.ListByIds(l.ctx, req.ApiIds)
                if err != nil {
                        l.Logger.Errorf("查询接口信息失败：%v", err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(apiModels) != len(req.ApiIds) {
                        return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分api接口信息不存在")
                }
        }

        // 3.开启事务：先删后插更新关联关系
        err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
                // 3.1 更新关联菜单
                if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuTrans(l.ctx, tx, req.Id); err != nil {
                        l.Logger.Errorf("删除角色[%d]旧菜单关联失败：%v", req.Id, err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(req.MenuIds) > 0 {
                        roleMenus := make([]systemmodel.SysRoleMenu, 0, len(req.MenuIds))
                        for _, menuId := range req.MenuIds {
                                roleMenus = append(roleMenus, systemmodel.SysRoleMenu{
                                        RoleId: req.Id,
                                        MenuId: menuId,
                                })
                        }
                        if _, err := l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus); err != nil
                                l.Logger.Errorf("插入角色[%d]新菜单关联失败：%v", req.Id, err)
                                return xerr.NewCodeError(xerr.ErrInternal)
                        }
                }

                // 3.2 更新关联接口
                if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiTrans(l.ctx, tx, req.Id); err != nil {
                        l.Logger.Errorf("删除角色[%d]旧接口关联失败：%v", req.Id, err)
                        return xerr.NewCodeError(xerr.ErrInternal)
                }
                if len(req.ApiIds) > 0 {
                        roleApis := make([]systemmodel.SysRoleApi, 0, len(req.ApiIds))
                        for _, apiId := range req.ApiIds {
                                roleApis = append(roleApis, systemmodel.SysRoleApi{
                                        RoleId: req.Id,
                                        ApiId:  apiId,
                                })
                        }
                        if _, err := l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis); err != nil {
                                l.Logger.Errorf("插入角色[%d]新接口关联失败：%v", req.Id, err)
                                return xerr.NewCodeError(xerr.ErrInternal)
                        }
                }

                return nil
        })
        if err != nil {
                l.Logger.Errorf("更新角色权限事务执行失败: %v", err)
                return xerr.NewCodeError(xerr.ErrInternal)
        }

        // 4.同步Casbin策略（事务成功后执行）
        if len(apiModels) > 0 {
                rules := make([][]string, 0, len(apiModels))
                for _, api := range apiModels {
                        rules = append(rules, []string{api.ApiPath, api.Method})
                }
                if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, existRole.Code, rules); err != nil {
                        l.Logger.Errorf("同步Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
                }
        } else {
                // 菜单和接口全部清空，同步清除Casbin策略
                if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, existRole.Code); err != nil {
                        l.Logger.Errorf("清除Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
                }
        }

        return nil
}

  ---
5. UpdateRolePermissions Handler

新文件：internal/handler/system/role/update_role_permissions_handler.go

package role

import (
"net/http"
"plating/pkg/response"
"plating/pkg/xerr"

        "plating/internal/logic/system/role"
        "plating/internal/svc"
        "plating/internal/types"

        "github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
var req types.UpdateRolePermissionsReq
if err := httpx.Parse(r, &req); err != nil {
response.FailWithMsg(w, r, err.Error())
return
}

                l := role.NewUpdateRolePermissionsLogic(r.Context(), svcCtx)
                err := l.UpdateRolePermissions(&req)
                if err != nil {
                        if codeErr, ok := err.(*xerr.CodeError); ok {
                                response.Fail(w, r, codeErr.Code, codeErr.Msg)
                                return
                        }
                        response.FailInternal(w, r)
                        return
                } else {
                        response.OK(w, r)
                }
        }
}

  ---
6. 路由注册

你需要找到路由注册的地方，把新 handler 挂上去。看一下你的路由文件在哪。



Searched for 2 patterns, read 1 file (ctrl+o to expand)

● 在 routes.go 第 381 行之后（UpdateRoleHandler 那条之后），加一条路由：

{
Method:  http.MethodPut,
Path:    "/role/:id/permissions",
Handler: systemrole.UpdateRolePermissionsHandler(serverCtx),
},

  ---
总结：改动全景

api/desc/system/role.api
- UpdateRoleReq 去掉 MenuIds、ApiIds，加上 Sort
- 新增 UpdateRolePermissionsReq
- 新增路由 PUT /role/:id/permissions

internal/types/system_role.go
- 同步改 UpdateRoleReq
- 新增 UpdateRolePermissionsReq

internal/model/system/sys_role_model.go
- UpdateRoleTrans 改成真正更新字段
- 新增 SoftDeleteRoleTrans（给后面删除角色用）

internal/logic/system/role/update_role_logic.go
- 只更新基本信息，不碰关联表和Casbin
- 编码变更时迁移Casbin策略

internal/logic/system/role/update_role_permissions_logic.go  (新文件)
- 先删后插菜单和接口关联
- 事务后同步Casbin策略

internal/handler/system/role/update_role_permissions_handler.go  (新文件)
- 新handler

internal/handler/routes.go
- 新增路由注册

代码都给你了，按这个顺序改就行。types 那个文件如果你用 goctl api go 生成的话，改完 .api 文件重新生成即可，handler
和 logic 需要手动建。