// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"net/http"
	"plating/internal/config"
	"plating/internal/middleware"
	"plating/pkg/casbin"
	"plating/pkg/orm"
	pkgsqlx "plating/pkg/sqlx"
	"time"

	platemodel "plating/internal/model/plate"
	systemmodel "plating/internal/model/system"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	// DB go-zero sqlx数据库连接（用于业务Model层查询）
	Orm *gorm.DB
	// RDB Redis客户端（用于Token黑名单、缓存等）
	RDB *redis.Client
	// Enforcer Casbin权限执行器（用于API接口鉴权）
	Enforcer *casbinv2.Enforcer

	// JWT认证中间件实例
	AuthMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

	// Casbin鉴权中间件实例（需要Enforcer和DB连接用于动态查询用户角色）
	CasbinMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

	// --- 系统管理 Model 层 ---
	// SysUserModel 系统用户数据访问
	SysUserModel systemmodel.SysUserModel
	// SysRoleModel 系统角色数据访问
	SysRoleModel systemmodel.SysRoleModel
	// SysUserRoleModel 用户角色关联数据访问
	SysUserRoleModel systemmodel.SysUserRoleModel
	// SysMenuModel 系统菜单数据访问
	SysMenuModel systemmodel.SysMenuModel
	// SysRoleMenuModel 角色菜单关联数据访问
	SysRoleMenuModel systemmodel.SysRoleMenuModel
	// SysApiModel 系统接口权限数据访问
	SysApiModel systemmodel.SysApiModel
	// SysRoleApiModel 角色接口关联数据访问
	SysRoleApiModel systemmodel.SysRoleApiModel
	// SysDictTypeModel 字典类型数据访问
	SysDictTypeModel systemmodel.SysDictTypeModel
	// SysDictDataModel 字典数据访问
	SysDictDataModel systemmodel.SysDictDataModel
	// SysLoginLogModel 登录日志数据访问
	SysLoginLogModel systemmodel.SysLoginLogModel
	// SysOperLogModel 操作日志数据访问
	SysOperLogModel systemmodel.SysOperLogModel
	// SysFileModel 文件记录数据访问
	SysFileModel systemmodel.SysFileModel

	// --- 槽液分析业务 Model 层 ---

	// PlateTankConfigModel 槽体配置数据访问
	PlateTankConfigModel platemodel.PlateTankConfigModel
	// PlateModelStateModel 模型状态数据访问
	PlateModelStateModel platemodel.PlateModelStateModel
	// PlateProductionEventModel 生产事件数据访问
	PlateProductionEventModel platemodel.PlateProductionEventsModel
	// DosingEventModel 加药事件数据访问
	PlateDosingEventModel platemodel.PlateDosingEventsModel
	// PlateWaterEventModel 补水事件数据访问
	PlateWaterEventModel platemodel.PlateWaterEventsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 数据库连接（go-zero sqlx风格），并用超时包装器包裹。
	// 包装器使所有 Model 层的 DB 操作使用独立的 10s 超时上下文，
	// 避免被 HTTP 请求上下文的短超时（RestConf.Timeout 默认 3s）截断。
	rawConn := sqlx.NewMysql(c.DB.DataSource)
	conn := pkgsqlx.NewTimeoutConn(rawConn, 30*time.Second)
	logx.Info("sqlx MySQL连接初始化成功")

	// 初始化gorm
	db := orm.NewMysql(&orm.Config{
		DSN:         c.DB.DataSource,
		Active:      20,
		Idle:        10,
		IdleTimeout: time.Hour * 24,
	})

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:        c.BizRedis.Host,
		Password:    c.BizRedis.Pass,
		DB:          c.BizRedis.DB,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})

	// 验证Redis连接是否正常
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis连接失败：%v", err))
	}
	logx.Info("RedisBiz 连接初始化成功")

	// 初始化 Casbin 执行器（使用gorm-adapter连接MySQL）
	enforcer, err := casbin.NewCasbin(db, c.CasbinModelPath)
	if err != nil {
		panic(fmt.Sprintf("Casbin初始化失败：%v", err))
	}
	logx.Info("Casbin执行器初始化成功")

	return &ServiceContext{
		Config:   c,
		Orm:      db,
		RDB:      rdb,
		Enforcer: enforcer,

		AuthMiddleware:   middleware.AuthMiddleware(c),
		CasbinMiddleware: middleware.CasbinMiddleware(enforcer, conn, c.CacheRedis, db),
		// 系统管理 Model 层初始化
		SysUserModel:     systemmodel.NewSysUserModel(conn, c.CacheRedis, db),
		SysRoleModel:     systemmodel.NewSysRoleModel(conn, c.CacheRedis, db),
		SysUserRoleModel: systemmodel.NewSysUserRoleModel(conn, c.CacheRedis, db),
		SysMenuModel:     systemmodel.NewSysMenuModel(conn, c.CacheRedis, db),
		SysRoleMenuModel: systemmodel.NewSysRoleMenuModel(conn, c.CacheRedis, db),
		SysApiModel:      systemmodel.NewSysApiModel(conn, c.CacheRedis, db),
		SysRoleApiModel:  systemmodel.NewSysRoleApiModel(conn, c.CacheRedis, db),
		SysDictTypeModel: systemmodel.NewSysDictTypeModel(conn, c.CacheRedis, db),
		SysDictDataModel: systemmodel.NewSysDictDataModel(conn, c.CacheRedis, db),
		SysLoginLogModel: systemmodel.NewSysLoginLogModel(conn, c.CacheRedis),
		SysOperLogModel:  systemmodel.NewSysOperLogModel(conn, c.CacheRedis),
		SysFileModel:     systemmodel.NewSysFileModel(conn, c.CacheRedis),

		// 槽液分析业务 Model 层初始化
		PlateTankConfigModel:      platemodel.NewPlateTankConfigModel(conn, c.CacheRedis),
		PlateModelStateModel:      platemodel.NewPlateModelStateModel(conn, c.CacheRedis),
		PlateProductionEventModel: platemodel.NewPlateProductionEventsModel(conn, c.CacheRedis),
		PlateDosingEventModel:     platemodel.NewPlateDosingEventsModel(conn, c.CacheRedis),
		PlateWaterEventModel:      platemodel.NewPlateWaterEventsModel(conn, c.CacheRedis),
	}
}
