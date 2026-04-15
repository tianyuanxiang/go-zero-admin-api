# Casbin 权限管理完全指南

> 本文档结合本项目（plating）的实际代码，从零讲解 Casbin 的核心概念和完整使用方式。
> 读完此文档，你将彻底掌握 Casbin 的 RBAC 权限模型，以及它在 go-zero 项目中的集成方法。

---

## 一、Casbin 是什么？

Casbin 是一个开源的 Go 权限管理框架，核心职责只有一件事：

> **回答一个问题：用户 X 是否有权限对资源 Y 执行操作 Z？**

它不管理用户账号和密码，那是认证（Authentication）的事。
Casbin 只管授权（Authorization）：**谁能做什么**。

### 1.1 认证 vs 授权

先搞清两个概念，否则后面会混淆：

| 概念 | 英文 | 职责 | 本项目中谁负责 |
|------|------|------|---------------|
| 认证 | Authentication | 你是谁？（验证身份） | `AuthMiddleware`（JWT） |
| 授权 | Authorization | 你能做什么？（验证权限） | `CasbinMiddleware`（Casbin） |

**先认证，后授权**。请求进来先过 AuthMiddleware 确认"你是谁"，再过 CasbinMiddleware 确认"你能不能做这件事"。

### 1.2 和传统硬编码鉴权的对比

```go
// 传统做法：每个接口里硬编码权限检查，改一个接口就要改一次代码
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    role := getRole(r)
    if role != "admin" && role != "manager" {
        http.Error(w, "没有权限", 403)
        return
    }
    // 业务逻辑...
}

// Casbin做法：统一中间件拦截，业务代码完全不用管权限
// 只需在数据库里维护"哪个角色能访问哪个接口"的规则即可
// 新增接口？往数据库加一行规则就行，不用改代码
```

---

## 二、核心概念：PERM 模型

Casbin 基于 **PERM 模型**（Policy, Effect, Request, Matchers），四个字母分别对应四个概念：

| 组件 | 含义 | 类比 |
|------|------|------|
| **R**equest（请求） | 待鉴权的三元组：谁、什么资源、什么操作 | 有人刷卡想进门 |
| **P**olicy（策略） | 存储在数据库中的权限规则 | 门禁白名单 |
| **M**atchers（匹配器） | 如何将 Request 和 Policy 进行比较 | 门禁比对规则 |
| **E**ffect（效果） | 多条策略命中时的最终裁决 | 最终放行还是拦截 |

**通俗理解整个流程**：

```
有人刷卡想进A栋（Request）
  → 门禁系统拿着刷卡信息去比对白名单（Policy）
  → 按照预设的比对规则逐条检查（Matchers）
  → 只要白名单里有一条匹配就放行（Effect）
```

---

## 三、本项目的 model.conf 详解

文件位置：`etc/rbac_model.conf`

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

### 逐行解读

#### 3.1 `[request_definition]` -- 请求长什么样

```ini
r = sub, obj, act
```

定义了每次鉴权请求包含三个要素：

| 字段 | 含义 | 本项目中的实际值 | 示例 |
|------|------|-----------------|------|
| `sub` | 主体（Subject），谁在请求 | 角色编码（role.Code） | `"admin"`, `"operator"` |
| `obj` | 对象（Object），访问什么 | HTTP 请求路径 | `"/api/system/user"` |
| `act` | 操作（Action），做什么动作 | HTTP 方法 | `"GET"`, `"POST"`, `"DELETE"` |

**注意**：`sub` 不是用户ID，是**角色编码**。一个用户可以有多个角色，中间件会逐个角色去检查。

#### 3.2 `[policy_definition]` -- 规则长什么样

```ini
p = sub, obj, act
```

数据库 `casbin_rule` 表中每一行存储的就是一条策略，格式和请求一样是三个字段。

**对应到数据库表**：

| 表字段 | 对应 | 含义 | 示例值 |
|--------|------|------|--------|
| `ptype` | 策略类型 | `p` 表示普通策略 | `"p"` |
| `v0` | sub（角色编码） | 哪个角色 | `"admin"` |
| `v1` | obj（路径） | 能访问什么路径 | `"/api/system/user"` |
| `v2` | act（方法） | 用什么HTTP方法 | `"GET"` |

#### 3.3 `[role_definition]` -- 角色继承

```ini
g = _, _
```

`g` 定义角色之间的继承关系。`_, _` 表示支持两个参数：子角色和父角色。

**本项目当前没有使用角色继承**。如果未来需要"经理继承操作员的所有权限"，可以在 `casbin_rule` 表中添加：

| ptype | v0 | v1 |
|-------|----|----|
| g | manager | operator |

这样 manager 自动拥有 operator 的所有权限，不用重复配置。

#### 3.4 `[policy_effect]` -- 最终裁决规则

```ini
e = some(where (p.eft == allow))
```

含义：**只要存在任意一条匹配的策略，就允许访问**。
这是"白名单"模式 -- 默认拒绝，只有明确授权的才放行。

#### 3.5 `[matchers]` -- 匹配规则（最核心）

```ini
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

这一行决定了"请求"和"规则"怎么比对。拆开看：

| 表达式 | 含义 | 解释 |
|--------|------|------|
| `g(r.sub, p.sub)` | 角色匹配 | 请求者的角色 == 策略中的角色（含继承关系） |
| `keyMatch2(r.obj, p.obj)` | **路径通配匹配** | 请求路径能匹配策略中的路径模式 |
| `r.act == p.act` | 方法精确匹配 | HTTP方法完全一致 |
| `p.act == "*"` | 方法通配 | 策略中方法为`*`时匹配任意方法 |

**三个条件用 `&&` 连接，必须同时满足才算匹配**。

##### keyMatch2 路径匹配详解

`keyMatch2` 是 Casbin 内置的路径匹配函数，支持 `:param` 风格的参数通配：

| 策略中的路径（p.obj） | 请求路径（r.obj） | 是否匹配 | 说明 |
|----------------------|-------------------|---------|------|
| `/api/system/user` | `/api/system/user` | 匹配 | 精确匹配 |
| `/api/system/user/:id` | `/api/system/user/123` | 匹配 | `:id` 通配任意值 |
| `/api/system/user/:id` | `/api/system/user/456/detail` | 不匹配 | 多了一层路径 |
| `/api/*` | `/api/system/user` | 匹配 | `*` 通配所有子路径 |
| `/api/*` | `/api/plating/event/dosing` | 匹配 | `*` 通配所有子路径 |

##### 方法通配详解

`(r.act == p.act || p.act == "*")` 的含义：

| 策略中的方法（p.act） | 请求方法（r.act） | 是否匹配 |
|---------------------|------------------|---------|
| `GET` | `GET` | 匹配 |
| `GET` | `POST` | 不匹配 |
| `*` | `GET` | 匹配 |
| `*` | `POST` | 匹配 |
| `*` | 任意方法 | 都匹配 |

---

## 四、本项目的数据库表结构

### 4.1 casbin_rule 表（Casbin 自动管理）

文件位置：`schema/plating.sql`

```sql
CREATE TABLE `casbin_rule` (
  `id`    bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ptype` varchar(100) NULL DEFAULT NULL,  -- 策略类型: p=普通策略, g=角色继承
  `v0`    varchar(100) NULL DEFAULT NULL,  -- 对应 sub（角色编码）
  `v1`    varchar(100) NULL DEFAULT NULL,  -- 对应 obj（路径）
  `v2`    varchar(100) NULL DEFAULT NULL,  -- 对应 act（方法）
  `v3`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用）
  `v4`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用）
  `v5`    varchar(100) NULL DEFAULT NULL,  -- 扩展字段（本项目未用）
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
) COMMENT = 'Casbin权限规则表';
```

**初始数据**（超级管理员拥有所有权限）：

```sql
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');
```

这一条规则的含义：`admin` 角色可以用**任意HTTP方法**访问 `/api/` 下的**所有路径**。

### 4.2 相关联的业务表

Casbin 只管 `casbin_rule` 这一张表，但鉴权流程中还会查询以下表：

| 表名 | 作用 | 关键字段 |
|------|------|---------|
| `sys_user` | 用户表 | id, username, status |
| `sys_role` | 角色表 | id, code, name, status |
| `sys_user_role` | 用户-角色关联表 | user_id, role_id |
| `sys_api` | API接口注册表 | path, method, group, description |
| `sys_role_api` | 角色-API关联表 | role_id, api_id |

**这些表之间的关系**：

```
sys_user  ←──  sys_user_role  ──→  sys_role
                                      │
                                      ↓
                                sys_role_api  ──→  sys_api
                                      │
                                      ↓
                                casbin_rule（Casbin自动读写）
```

通俗来说：
- `sys_user_role` 记录"用户有哪些角色"
- `sys_role_api` 记录"角色能访问哪些接口"（业务层面的关联关系）
- `casbin_rule` 记录"角色能访问哪些接口"（Casbin 层面的策略，和 sys_role_api 是同步的）

为什么 `sys_role_api` 和 `casbin_rule` 看起来重复？因为 `sys_role_api` 是给前端页面展示用的（"这个角色勾选了哪些接口"），`casbin_rule` 是给 Casbin 引擎做鉴权用的。两者通过 "为角色分配API" 的接口保持同步。

---

## 五、完整鉴权链条：从请求到放行

下面用一个真实的例子，把整条链路串起来。

### 场景：操作员张三要录入一条加药事件

```
POST /api/plating/event/dosing
Authorization: Bearer eyJhbGci...（张三的JWT令牌）
Content-Type: application/json

{"tankId": "T001", "eventTime": "2026-04-15 10:30:00", ...}
```

### 第一步：路由匹配

go-zero 根据 `routes.go` 找到这个路由，它属于带有 `[AuthMiddleware, CasbinMiddleware]` 的路由组：

```go
// routes.go 第66-123行
server.AddRoutes(
    rest.WithMiddlewares(
        []rest.Middleware{serverCtx.AuthMiddleware, serverCtx.CasbinMiddleware},
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/event/dosing",
                Handler: plateevent.CreateDosingEventHandler(serverCtx),
            },
            // ...
        }...,
    ),
    rest.WithPrefix("/api/plating"),
)
```

**中间件按顺序执行**：先 AuthMiddleware，再 CasbinMiddleware，最后才到 Handler。

### 第二步：AuthMiddleware -- 确认"你是谁"

文件位置：`internal/middleware/auth_middleware.go`

```
1. 从请求头提取 Authorization: Bearer eyJhbGci...
2. 解析JWT令牌，提取 userId=5, username="zhangsan"
3. 将 userId 和 username 写入请求的 Context 中
4. 放行，交给下一个中间件
```

如果JWT无效或过期，直接返回 401，不会走到 CasbinMiddleware。

### 第三步：CasbinMiddleware -- 确认"你能不能做"

文件位置：`internal/middleware/casbin_middleware.go`

```
1. 从 Context 读取 userId=5
   （调用 GetUserIdFromCtx(r.Context())）

2. 查 sys_user_role 表：userId=5 的角色有哪些？
   → 得到 roleIds = [2]（张三只有一个角色）

3. 查 sys_role 表：roleId=2 的角色编码是什么？
   → 得到 role.Code = "operator"

4. 调用 Casbin 鉴权：
   enforcer.Enforce("operator", "/api/plating/event/dosing", "POST")

5. Casbin 在内存中匹配 casbin_rule 表的规则：
   找到一条：("p", "operator", "/api/plating/event/dosing", "POST")
   
   用 matchers 比对：
   - g("operator", "operator") → 角色匹配
   - keyMatch2("/api/plating/event/dosing", "/api/plating/event/dosing") → 路径匹配
   - "POST" == "POST" → 方法匹配
   
   三个条件都满足 → 返回 true

6. hasPermission = true → 放行，交给 Handler 处理业务逻辑
```

**如果张三没有这个权限呢？**

```
Casbin 匹配不到任何规则 → 返回 false
→ 中间件返回 403 Forbidden，请求到此结束，不会进入 Handler
```

### 流程图总结

```
客户端请求
    │
    ▼
路由匹配（routes.go）
    │
    ▼
AuthMiddleware（JWT认证）
    │ 失败 → 401 Unauthorized
    ▼ 成功
CasbinMiddleware（权限校验）
    │
    ├─ 1. 从Context取userId
    ├─ 2. 查sys_user_role → 得到roleIds
    ├─ 3. 查sys_role → 得到role.Code
    ├─ 4. enforcer.Enforce(role.Code, path, method)
    │     └─ 在内存中匹配casbin_rule的规则
    │
    │ 无权限 → 403 Forbidden
    ▼ 有权限
Handler（业务逻辑）
    │
    ▼
返回响应
```

---

## 六、哪些路由受 Casbin 保护？

从 `routes.go` 可以看出，路由分为三个安全等级：

### 6.1 完全公开（无任何中间件）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/auth/login` | POST | 登录 |
| `/api/auth/refresh` | POST | 刷新令牌 |

### 6.2 仅需登录（只有 AuthMiddleware，无 Casbin）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/auth/changePassword` | POST | 修改自己的密码 |
| `/api/auth/logout` | POST | 登出 |
| `/api/auth/userInfo` | GET | 获取当前用户信息 |

**这些接口只要登录了就能访问，不检查角色权限**。因为它们是每个用户都应该有的基础操作。

### 6.3 需要登录 + 角色授权（AuthMiddleware + CasbinMiddleware）

以下所有接口都受 Casbin 保护，必须在 `casbin_rule` 表中有对应规则才能访问：

**系统管理 -- 用户**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/user` | POST | 创建用户 |
| `/api/system/user` | GET | 用户列表 |
| `/api/system/user/:id` | PUT | 更新用户 |
| `/api/system/user/:id` | DELETE | 删除用户 |
| `/api/system/user/:id` | GET | 用户详情 |
| `/api/system/user/:id/reset-password` | POST | 重置用户密码 |

**系统管理 -- 角色**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/role` | POST | 创建角色 |
| `/api/system/role` | GET | 角色列表 |
| `/api/system/role/:id` | PUT | 更新角色 |
| `/api/system/role/:id` | DELETE | 删除角色 |
| `/api/system/role/:id` | GET | 角色详情 |
| `/api/system/role/all` | GET | 所有角色（不分页） |

**系统管理 -- 菜单**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/menu` | POST | 创建菜单 |
| `/api/system/menu/:id` | PUT | 更新菜单 |
| `/api/system/menu/:id` | DELETE | 删除菜单 |
| `/api/system/menu/current` | GET | 当前用户的菜单 |
| `/api/system/menu/tree` | GET | 菜单树 |

**系统管理 -- 接口**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/api` | POST | 注册接口 |
| `/api/system/api` | GET | 接口列表 |
| `/api/system/api/:id` | PUT | 更新接口 |
| `/api/system/api/:id` | DELETE | 删除接口 |
| `/api/system/api/all` | GET | 所有接口（不分页） |

**系统管理 -- 字典**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/dict/type` | POST | 创建字典类型 |
| `/api/system/dict/type` | GET | 字典类型列表 |
| `/api/system/dict/type/:id` | PUT | 更新字典类型 |
| `/api/system/dict/type/:id` | DELETE | 删除字典类型 |
| `/api/system/dict/data` | POST | 创建字典数据 |
| `/api/system/dict/data/:dictType` | GET | 按类型查字典数据 |
| `/api/system/dict/data/:id` | PUT | 更新字典数据 |
| `/api/system/dict/data/:id` | DELETE | 删除字典数据 |

**系统管理 -- 文件**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/file` | GET | 文件列表 |
| `/api/system/file/:id` | DELETE | 删除文件 |
| `/api/system/file/upload` | POST | 上传文件 |

**系统管理 -- 日志**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/system/log/login` | GET | 登录日志列表 |
| `/api/system/log/login/clear` | DELETE | 清空登录日志 |
| `/api/system/log/oper` | GET | 操作日志列表 |
| `/api/system/log/oper/clear` | DELETE | 清空操作日志 |

**业务功能 -- 槽液事件**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/event/dosing` | POST | 录入加药事件 |
| `/api/plating/event/dosing` | GET | 加药事件列表 |
| `/api/plating/event/dosing/:id` | DELETE | 删除加药事件 |
| `/api/plating/event/production` | POST | 录入生产事件 |
| `/api/plating/event/production` | GET | 生产事件列表 |
| `/api/plating/event/production/:id` | DELETE | 删除生产事件 |
| `/api/plating/event/water` | POST | 录入换水事件 |
| `/api/plating/event/water` | GET | 换水事件列表 |
| `/api/plating/event/water/:id` | DELETE | 删除换水事件 |
| `/api/plating/event/trigger-calc` | POST | 触发计算 |

**业务功能 -- 槽体状态**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/state/:tankId` | GET | 最新状态 |
| `/api/plating/state/export` | GET | 导出报表 |
| `/api/plating/state/override` | POST | 覆盖状态 |
| `/api/plating/state/trend` | GET | 趋势数据 |

**业务功能 -- 槽体配置**

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/plating/tank` | POST | 创建槽体 |
| `/api/plating/tank` | GET | 槽体列表 |
| `/api/plating/tank/:tankId` | PUT | 更新槽体 |
| `/api/plating/tank/:tankId` | DELETE | 删除槽体 |
| `/api/plating/tank/:tankId` | GET | 槽体详情 |
| `/api/plating/tank/init-state` | POST | 初始化模型状态 |

---

## 七、超级管理员是怎么拥有所有权限的？

初始化SQL中只给 admin 角色写了一条规则：

```sql
INSERT INTO `casbin_rule` VALUES (1, 'p', 'admin', '/api/*', '*', '', '', '');
```

翻译一下：

| 字段 | 值 | 含义 |
|------|-----|------|
| ptype | `p` | 这是一条策略 |
| v0 (sub) | `admin` | admin 角色 |
| v1 (obj) | `/api/*` | 所有 `/api/` 下的路径 |
| v2 (act) | `*` | 所有HTTP方法 |

当 admin 角色访问任何接口时：

```
enforcer.Enforce("admin", "/api/system/user/5/reset-password", "POST")

匹配过程：
1. g("admin", "admin") → 角色匹配
2. keyMatch2("/api/system/user/5/reset-password", "/api/*") → 路径匹配（*通配所有子路径）
3. "POST" == "*" → 不等... 但 p.act == "*" → 方法通配匹配

三个条件都满足 → 返回 true → 放行
```

**所以 admin 不需要一条条配置，一条通配规则搞定一切。**

---

## 八、本项目的 Casbin 初始化代码

### 8.1 Casbin 引擎初始化

文件位置：`pkg/casbin/casbin.go`

```go
func NewCasbin(db *gorm.DB, modelPath string) (*casbinv2.Enforcer, error) {
    // 1. 使用 gorm-adapter 连接 MySQL
    //    Casbin 规则自动读写 casbin_rule 表
    //    如果表不存在，gorm-adapter 会自动建表
    adapter, err := gormadapter.NewAdapterByDB(db)
    
    // 2. 加载模型配置文件（etc/rbac_model.conf）
    //    定义了请求格式、策略格式、匹配规则
    enforcer, err := casbinv2.NewEnforcer(modelPath, adapter)
    
    // 3. 从数据库加载所有策略到内存
    //    鉴权时直接在内存中匹配，不查数据库，性能极高
    enforcer.LoadPolicy()
    
    return enforcer, nil
}
```

**关键理解**：Casbin 把数据库中的规则**一次性加载到内存**，后续所有 `Enforce()` 调用都在内存中比对，不会频繁查数据库。这就是它性能高的原因。但代价是：**修改了数据库中的规则后，必须调用 `LoadPolicy()` 重新加载，否则内存中还是旧规则。**

### 8.2 注入到 ServiceContext

文件位置：`internal/svc/service_context.go`

```go
type ServiceContext struct {
    // ... 其他字段
    Enforcer         *casbinv2.Enforcer                        // Casbin执行器
    AuthMiddleware   func(handlerFunc http.HandlerFunc) http.HandlerFunc
    CasbinMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc
    // ...
}

// 初始化时
enforcer, err := casbin.NewCasbin(db, c.CasbinModelPath)
// ...
return &ServiceContext{
    Enforcer:         enforcer,
    AuthMiddleware:   middleware.AuthMiddleware(c),
    CasbinMiddleware: middleware.CasbinMiddleware(enforcer, conn, c.CacheRedis, db),
    // ...
}
```

**配置文件**（`etc/plating-api.yaml`）：

```yaml
CasbinModelPath: "etc/rbac_model.conf"
```

### 8.3 鉴权中间件的完整代码

文件位置：`internal/middleware/casbin_middleware.go`

```go
func CasbinMiddleware(enforcer *casbinv2.Enforcer, conn sqlx.SqlConn,
    c cache.CacheConf, db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
    
    // 初始化需要查询的Model
    userRoleModel := systemmodel.NewSysUserRoleModel(conn, c, db)
    roleModel := systemmodel.NewSysRoleModel(conn, c, db)

    return func(next http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 第一步：从Context获取当前用户ID（AuthMiddleware写入的）
            userId := GetUserIdFromCtx(r.Context())
            if userId == 0 {
                response.FailUnauthorized(w, r)
                return
            }

            reqPath := r.URL.Path   // 如：/api/system/user
            reqMethod := r.Method   // 如：GET

            // 第二步：查询该用户的所有角色ID
            roleIds, err := userRoleModel.GetRoleIdsByUserId(r.Context(), userId)
            // ...（错误处理）

            if len(roleIds) == 0 {
                // 用户没有任何角色，直接拒绝
                response.FailForbidden(w, r)
                return
            }

            // 第三步：逐个角色检查权限（OR逻辑：任意一个角色有权限就放行）
            hasPermission := false
            for _, roleId := range roleIds {
                role, err := roleModel.FindOneByRoleId(r.Context(), roleId)
                if err != nil || role == nil {
                    continue
                }

                // 用角色编码（code）调用Casbin鉴权
                allowed, err := enforcer.Enforce(role.Code, reqPath, reqMethod)
                if err != nil {
                    continue
                }

                if allowed {
                    hasPermission = true
                    break  // 有一个角色有权限就够了，不用继续检查
                }
            }

            if !hasPermission {
                response.FailForbidden(w, r)
                return
            }

            // 有权限，放行到下一个Handler
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 九、casbin_rule 表数据配置实战

### 9.1 给"操作员"角色配置权限

假设你的系统有两个角色：

| id | code | name |
|----|------|------|
| 1 | admin | 超级管理员 |
| 2 | operator | 操作员 |

admin 已经有通配规则了。现在需要给 operator 配置"能录入和查看加药事件、能查看槽体状态"的权限。

在 `casbin_rule` 表中插入：

```sql
-- 操作员可以录入和查看加药事件
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/event/dosing', 'POST'),
('p', 'operator', '/api/plating/event/dosing', 'GET');

-- 操作员可以查看槽体状态和趋势
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/plating/state/:tankId', 'GET'),
('p', 'operator', '/api/plating/state/trend', 'GET');

-- 操作员可以查看自己的菜单
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'operator', '/api/system/menu/current', 'GET');
```

**注意**：路径中的 `:tankId` 和 `:id` 是 `keyMatch2` 的参数通配语法，能匹配实际请求中的 `/api/plating/state/T001`。

### 9.2 完整的 casbin_rule 表示例

| id | ptype | v0 | v1 | v2 |
|----|-------|------------|------------------------------------------|--------|
| 1 | p | admin | /api/* | * |
| 2 | p | operator | /api/plating/event/dosing | POST |
| 3 | p | operator | /api/plating/event/dosing | GET |
| 4 | p | operator | /api/plating/event/water | POST |
| 5 | p | operator | /api/plating/event/water | GET |
| 6 | p | operator | /api/plating/event/production | POST |
| 7 | p | operator | /api/plating/event/production | GET |
| 8 | p | operator | /api/plating/state/:tankId | GET |
| 9 | p | operator | /api/plating/state/trend | GET |
| 10 | p | operator | /api/system/menu/current | GET |
| 11 | p | viewer | /api/plating/state/:tankId | GET |
| 12 | p | viewer | /api/plating/state/trend | GET |
| 13 | p | viewer | /api/plating/state/export | GET |
| 14 | p | viewer | /api/system/menu/current | GET |

**规则设计思路**：
- `admin`：一条通配搞定，拥有所有权限
- `operator`：只能操作业务功能（录入事件、查看状态），不能管理系统
- `viewer`：只读权限，只能查看状态和导出报表

---

## 十、本项目中完整的权限配置流程（通过API操作）

### 场景：给"操作员"角色配置"可以录入加药事件"的权限

**第一步：在系统中注册 API 接口**（只需做一次）

```
POST /api/system/api
{
    "path": "/api/plating/event/dosing",
    "method": "POST",
    "group": "槽液事件",
    "description": "录入加药事件"
}
```

接口信息写入 `sys_api` 表，得到 id=10。

**第二步：确认"操作员"角色存在**

```
POST /api/system/role
{
    "name": "操作员",
    "code": "operator",
    "status": 1
}
```

角色编码 `operator` 将作为 Casbin 中的 sub。

**第三步：为角色绑定 API 权限**

这一步会同时操作两张表：

1. 在 `sys_role_api` 表中记录"角色2绑定了接口10"（给前端展示用）
2. 在 `casbin_rule` 表中写入策略 `("p", "operator", "/api/plating/event/dosing", "POST")`（给Casbin鉴权用）
3. 调用 `enforcer.LoadPolicy()` 刷新内存中的策略

**第四步：为用户分配"操作员"角色**

将 userId 和 roleId 的关系写入 `sys_user_role` 表。

**第五步：用户发起请求，自动鉴权**

```
POST /api/plating/event/dosing
Authorization: Bearer eyJhbGci...
```

CasbinMiddleware 自动完成鉴权，业务代码无需关心权限。

---

## 十一、常用 Casbin API

本项目在 `pkg/casbin/casbin.go` 中封装了常用操作：

### 11.1 添加策略

```go
// 单条添加：允许 operator 角色 POST /api/plating/event/dosing
casbin.AddPolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// 批量添加（全量覆盖模式：先删旧的，再写新的）
rules := [][]string{
    {"/api/plating/event/dosing", "POST"},
    {"/api/plating/event/dosing", "GET"},
    {"/api/plating/state/:tankId", "GET"},
}
casbin.AddRolePolicies(enforcer, "operator", rules)
```

### 11.2 删除策略

```go
// 删除单条
casbin.RemovePolicyForRole(enforcer, "operator", "/api/plating/event/dosing", "POST")

// 删除某角色的所有策略（删除角色时用）
casbin.RemoveAllPoliciesForRole(enforcer, "operator")
```

### 11.3 查询策略

```go
// 查询某角色的所有权限
policies, _ := casbin.GetRolePolicies(enforcer, "operator")
// 返回: [["operator", "/api/plating/event/dosing", "POST"], ...]

// 获取所有策略（调试用）
allPolicies := enforcer.GetPolicy()
```

### 11.4 检查权限

```go
// 检查 operator 是否有权限 POST /api/plating/event/dosing
allowed, _ := casbin.CheckPermission(enforcer, "operator", "/api/plating/event/dosing", "POST")
```

### 11.5 重新加载策略（重要）

```go
// 每次修改casbin_rule表后必须调用，否则内存中的策略不会更新
casbin.ReloadPolicy(enforcer)
```

---

## 十二、常见问题排查

### Q1：添加了策略但 Enforce 返回 false？

**最常见原因**：没有调用 `enforcer.LoadPolicy()` 重新加载。

Casbin 在内存中匹配，数据库写入后必须重新加载才生效。每次通过 API 修改策略后，代码中已经自动调用了 `LoadPolicy()`。但如果你直接操作数据库（比如手动 INSERT），必须重启服务或手动触发加载。

### Q2：路径带参数时匹配不上？

本项目的 matchers 使用 `keyMatch2`，支持 `:param` 风格通配。

**正确的策略写法**：

```sql
-- 策略中用 :id 占位
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'operator', '/api/system/user/:id', 'GET');
```

这样 `/api/system/user/123`、`/api/system/user/456` 都能匹配。

**常见错误**：策略中写了精确的ID，如 `/api/system/user/123`，这样只能匹配用户123。

### Q3：多个角色时如何判定？

用户可以同时拥有多个角色（比如既是 operator 又是 viewer）。中间件会**逐个角色检查**，只要有一个角色有权限就放行（OR 逻辑）。

```
用户角色: [operator, viewer]

检查 operator → enforcer.Enforce("operator", path, method) → false
检查 viewer  → enforcer.Enforce("viewer", path, method) → true
→ 有权限，放行
```

### Q4：角色被删除后，casbin_rule 中的规则会残留吗？

删除角色时，应该同步清理 Casbin 策略：

```go
// 删除该角色在 casbin_rule 表中的所有规则
casbin.RemoveAllPoliciesForRole(enforcer, roleCode)
casbin.ReloadPolicy(enforcer)
```

### Q5：如何调试当前内存中有哪些策略？

```go
// 打印所有策略
policies := enforcer.GetPolicy()
for _, p := range policies {
    fmt.Printf("角色:%s  路径:%s  方法:%s\n", p[0], p[1], p[2])
}
```

### Q6：为什么我直接在数据库加了规则但不生效？

因为 Casbin 在内存中做匹配。你往数据库插了数据，但内存还是旧的。两种解法：
1. 重启服务（服务启动时会 `LoadPolicy()`）
2. 通过代码调用 `enforcer.LoadPolicy()` 重新加载

**建议通过系统的API接口来管理权限，API内部会自动刷新策略。**

### Q7：路径中 `*` 和 `:param` 的区别？

| 模式 | 示例 | 匹配范围 |
|------|------|---------|
| `:param` | `/api/system/user/:id` | 只匹配一层：`/api/system/user/123` |
| `*` | `/api/*` | 匹配所有子路径：`/api/system/user/123/detail` 也能匹配 |

所以 admin 的 `/api/*` 能匹配所有 API 接口。

---

## 十三、本项目 Casbin 相关文件一览

| 文件 | 作用 |
|------|------|
| `etc/rbac_model.conf` | RBAC 模型定义文件，定义请求格式、策略格式、匹配规则 |
| `etc/plating-api.yaml` | 配置文件，`CasbinModelPath` 指定模型文件路径 |
| `pkg/casbin/casbin.go` | Casbin 初始化 + 策略操作封装函数 |
| `internal/middleware/casbin_middleware.go` | HTTP 鉴权中间件，拦截请求并执行 Enforce |
| `internal/middleware/auth_middleware.go` | JWT 认证中间件，为 Casbin 提供 userId |
| `internal/svc/service_context.go` | 初始化 Enforcer 并注入中间件 |
| `internal/handler/routes.go` | 路由注册，决定哪些路由受 Casbin 保护 |
| `internal/model/system/sys_role_model.go` | 角色Model，提供按ID查角色编码 |
| `internal/model/system/sys_user_role_model.go` | 用户角色关联Model，提供按用户ID查角色列表 |
| `internal/model/system/sys_role_api_model.go` | 角色API关联Model，记录角色绑定了哪些接口 |
| `internal/model/system/sys_api_model.go` | API接口Model，存储所有可配置的接口信息 |
| `schema/plating.sql` | 数据库表结构，包含 casbin_rule 表定义和初始数据 |

---

## 十四、快速备忘录

```
鉴权的三元组
(角色code, 请求路径, HTTP方法)  →  Enforce()  →  true/false

增删策略
enforcer.AddPolicy(role, path, method)           // 添加单条
enforcer.AddPolicies(rules)                       // 批量添加
enforcer.RemovePolicy(role, path, method)         // 删除单条
enforcer.RemoveFilteredPolicy(0, role)            // 删除角色的所有规则

查询
enforcer.GetPolicy()                              // 所有规则
enforcer.GetFilteredPolicy(0, role)               // 某角色的规则
enforcer.Enforce(role, path, method)              // 执行鉴权

刷新（每次改完casbin_rule后必须调用）
enforcer.LoadPolicy()

路径匹配规则（本项目使用 keyMatch2）
/api/system/user      精确匹配 /api/system/user
/api/system/user/:id  通配匹配 /api/system/user/123（一层）
/api/*                通配匹配 /api/ 下所有路径（多层）

方法匹配
GET/POST/PUT/DELETE   精确匹配
*                     匹配任意方法
```
