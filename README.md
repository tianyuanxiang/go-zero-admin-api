# go-zero-admin Platform

电镀管理平台后端 API 服务，基于 Go 语言 + go-zero 微服务框架构建。

## 技术栈

- **语言**: Go 1.24
- **框架**: [go-zero](https://go-zero.dev/) (REST API)
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT (Access Token + Refresh Token)
- **权限**: Casbin RBAC
- **Excel**: excelize

## 项目结构

```
go-zero-admin/
├── api/                    # API 定义文件
│   ├── desc/               # 分模块 API 描述
│   │   ├── auth.api        # 认证模块
│   │   ├── plate/          # 电镀业务模块
│   │   └── system/         # 系统管理模块
│   └── go-zero-admin.api         # API 主入口
├── etc/                    # 配置文件
│   ├── go-zero-admin-api.yaml    # 服务配置 (gitignore, 需自行创建)
│   ├── go-zero-admin-api.yaml.example  # 配置示例
│   └── rbac_model.conf     # Casbin RBAC 模型
├── internal/               # 私有业务代码
│   ├── common/             # 公共定义
│   ├── config/             # 配置结构体
│   ├── handler/            # HTTP Handler
│   │   ├── auth/           # 认证 (登录/登出/改密/刷新Token)
│   │   ├── plate/          # 电镀业务 (槽体/事件/状态)
│   │   └── system/         # 系统管理 (用户/角色/菜单/字典/日志/文件/API)
│   ├── logic/              # 业务逻辑层
│   ├── middleware/          # 中间件 (鉴权等)
│   ├── model/              # 数据模型
│   ├── svc/                # 服务上下文 (依赖注入)
│   └── types/              # 请求/响应类型
├── pkg/                    # 可复用公共包
│   ├── casbin/             # Casbin 权限工具
│   ├── encrypt/            # 加密工具
│   ├── excel/              # Excel 导入导出
│   ├── jwtx/               # JWT 工具
│   ├── orm/                # ORM 工具
│   ├── response/           # 统一响应
│   ├── sqlx/               # SQL 工具
│   ├── upload/             # 文件上传
│   └── xerr/               # 自定义错误
├── schema/                 # 数据库脚本
│   ├── init.sql            # 初始化数据
│   └── go-zero-admin.sql         # 表结构
├── docx/                   # 项目文档
├── go-zero-admin.go              # 程序入口
├── go.mod
└── go.sum
```

## 功能模块

### 认证模块 (Auth)
- 用户登录 / 登出
- JWT Token 签发与刷新
- 密码修改
- 获取当前用户信息

### 系统管理 (System)
- 用户管理 (CRUD)
- 角色管理 (CRUD + 权限分配)
- 菜单管理 (树形结构)
- 字典管理 (数据字典)
- 操作日志
- 文件管理 (上传/下载)
- API 管理 (接口权限)

### 电镀业务 (Plate)
- 槽体管理 (Tank)
- 事件管理 (Event)
- 状态管理 (State)

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 1. 克隆项目

```bash
git clone https://github.com/tianyuanxiang/go-zero-admin.git
cd go-zero-admin
```

### 2. 初始化数据库

```bash
# 导入表结构
mysql -u root -p plating < schema/plating.sql

# 导入初始化数据
mysql -u root -p plating < schema/init.sql
```

### 3. 配置文件

```bash
# 复制配置示例并修改
cp etc/plating-api.yaml.example etc/plating-api.yaml
# 编辑 etc/plating-api.yaml, 填写实际的数据库和 Redis 连接信息
```

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 启动服务

```bash
go run plating.go
# 或指定配置文件
go run plating.go -f etc/plating-api.yaml
```

服务启动后监听 `http://0.0.0.0:8888`

## 开发规范

### 分支管理

| 分支 | 用途 |
|------|------|
| `main` | 主分支，稳定版本，仅通过 PR 合并 |
| `dev` | 开发分支，日常开发在此提交 |

### 工作流程

```
1. 在 dev 分支开发和提交代码
2. 功能完成并测试通过后，创建 PR 合并到 main
3. main 分支始终保持可部署状态
```

## License

Private
