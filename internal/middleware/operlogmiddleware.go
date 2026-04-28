package middleware

import (
	"bytes"
	"database/sql"
	"go-zero-admin/internal/model/system"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

// responseWriterWrapper 包装 http.ResponseWriter 以捕获响应内容。
type responseWriterWrapper struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// WriteHeader 覆盖原始WriteHeader，记录状态码。
func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write 覆盖原始Write，同时记录响应内容。
func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// OperLogMiddleware 操作日志记录中间件。
//
// 记录所有需要审计的请求信息，包括：
// - 请求路径、方法、参数
// - 操作人员信息（从Context获取）
// - 响应结果（截取前500字符，防止日志过大）
// - 操作时间
//
// 注意：此中间件应在 AuthMiddleware 之后注册，确保能从 Context 获取用户信息。
//
// 参数：
//   - db : 数据库连接（用于写入操作日志表）
//
// 返回：
//   - func(http.Handler) http.Handler : 标准中间件函数
func OperLogMiddleware(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB) func(http.Handler) http.Handler {
	operLogModel := system.NewSysOperLogModel(conn, c, db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 记录操作开始时间
			operTime := time.Now()

			// 读取请求体（需要先缓存，因为Body只能读一次）
			var requestParam string
			if r.Body != nil && r.Method != http.MethodGet {
				bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 2048))
				if err == nil {
					requestParam = string(bodyBytes)
					// 将读取的内容重新写回Body，供后续Handler读取
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			// 包装ResponseWriter以捕获响应内容
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			// 执行后续处理链
			next.ServeHTTP(wrapper, r)

			// 异步写入操作日志，避免影响接口响应时间
			go func() {
				userId := GetUserIdFromCtx(r.Context())
				username := GetUsernameFromCtx(r.Context())

				// 截取响应内容（最多500字符）
				respBody := wrapper.body.String()
				if len(respBody) > 500 {
					respBody = respBody[:500] + "...(已截断)"
				}

				// 获取客户端真实IP（考虑代理转发）
				clientIp := getClientIP(r)

				// 判断操作是否成功（HTTP状态码2xx为成功）
				operStatus := 1
				if wrapper.statusCode >= 400 {
					operStatus = 0
				}

				// 确定操作类型（根据HTTP方法）
				businessType := getBusinessType(r.Method)

				operLog := &system.SysOperLog{
					Title:         getOperTitle(r.URL.Path),
					BusinessType:  int64(businessType),
					Method:        r.URL.Path,
					RequestMethod: r.Method,
					OperatorType:  1, // 1=后台用户
					OperatorName:  username,
					OperatorId:    userId,
					OperUrl:       r.URL.Path,
					OperIp:        clientIp,
					OperParam:     StringToNullString(requestParam),
					JsonResult:    StringToNullString(respBody),
					Status:        int64(operStatus),
					OperTime:      operTime,
				}

				ctx := r.Context()
				if _, err := operLogModel.Insert(ctx, operLog); err != nil {
					logx.WithContext(ctx).Errorf("写入操作日志失败：%v", err)
				}
			}()
		})
	}
}

// getClientIP 获取客户端真实IP地址。
//
// 优先从 X-Forwarded-For 头获取（代理转发场景），
// 其次从 X-Real-IP 头获取，最后从 RemoteAddr 截取。
func getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For（可能包含多个IP，取第一个）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// 检查X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 从RemoteAddr截取IP（格式：ip:port）
	remoteAddr := r.RemoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx > 0 {
		return remoteAddr[:idx]
	}

	return remoteAddr
}

// getBusinessType 根据HTTP方法判断业务操作类型。
//
// 返回值含义：0=其他，1=新增，2=修改，3=删除，4=查询
func getBusinessType(method string) int {
	switch strings.ToUpper(method) {
	case http.MethodPost:
		return 1 // 新增
	case http.MethodPut, http.MethodPatch:
		return 2 // 修改
	case http.MethodDelete:
		return 3 // 删除
	case http.MethodGet:
		return 4 // 查询
	default:
		return 0 // 其他
	}
}

// getOperTitle 根据请求路径生成操作模块标题。
//
// 简单提取路径中的模块名称作为日志标题。
func getOperTitle(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		// 格式：/api/module/resource -> 取 module/resource
		return parts[1] + "/" + parts[2]
	}
	if len(parts) >= 2 {
		return parts[1]
	}
	return path
}

func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			Valid: false,
		}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
