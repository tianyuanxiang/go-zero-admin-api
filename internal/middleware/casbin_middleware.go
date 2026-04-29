// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"go-zero-admin/pkg/response"
	"net/http"

	systemmodel "go-zero-admin/internal/model/system"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

// CasbinMiddleware Casbin RBAC接口鉴权中间件。
//
// 从请求Context中获取当前用户ID，查询其拥有的角色，
// 然后通过Casbin Enforcer验证该角色是否有权限访问当前接口。
//
// 鉴权失败（用户无角色、角色无权限）时直接返回 403 响应。
//
// 参数：
//   - enforcer      : Casbin执行器
//   - db            : 数据库连接（用于查询用户角色）
//
// 返回：
//   - func(http.Handler) http.Handler : 标准中间件函数
func CasbinMiddleware(enforcer *casbinv2.Enforcer, conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	userRoleModel := systemmodel.NewSysUserRoleModel(conn, c, db)
	roleModel := systemmodel.NewSysRoleModel(conn, c, db)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logx.Infof("\033[31m%s\033[0m", "***************************** casbinMiddleware start **********************")
			// 从Context获取当前用户ID（由AuthMiddleware写入）
			userId := GetUserIdFromCtx(r.Context())
			if userId == 0 {
				response.FailUnauthorized(w, r)
				return
			}

			// 获取当前请求的路径和方法
			reqPath := r.URL.Path
			reqMethod := r.Method

			// 查询用户拥有的所有角色ID
			roleIds, err := userRoleModel.GetRoleIdsByUserId(r.Context(), userId)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("查询用户[%d]角色失败：%v", userId, err)
				response.FailInternal(w, r)
				return
			}

			if len(roleIds) == 0 {
				logx.WithContext(r.Context()).Infof("用户[%d]没有任何角色，拒绝访问：%s %s", userId, reqMethod, reqPath)
				response.FailForbidden(w, r)
				return
			}

			// 逐个角色检查是否有权限访问当前接口
			// 只要有一个角色有权限就允许通过（OR逻辑）
			hasPermission := false
			for _, roleId := range roleIds {
				role, err := roleModel.FindOneByRoleId(r.Context(), roleId)
				if err != nil || role == nil {
					continue
				}

				// 使用角色编码（code）作为casbin中的subject进行权限检查
				allowed, err := enforcer.Enforce(role.Code, reqPath, reqMethod)
				if err != nil {
					logx.WithContext(r.Context()).Errorf("Casbin Enforce失败：%v", err)
					continue
				}

				if allowed {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				logx.WithContext(r.Context()).Infof(
					"用户[%d]权限不足，拒绝访问：%s %s",
					userId, reqMethod, reqPath,
				)
				response.FailForbidden(w, r)
				return
			}
			logx.Infof("\u001B[31m%s\u001B[0m", "********************************** casbinMiddleware end *********************************")
			next.ServeHTTP(w, r)
		})
	}
}
