// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"go-zero-admin/pkg/response"
	"net/http"
	"strings"

	"go-zero-admin/internal/logic/auth"
	"go-zero-admin/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Authorization头获取token（用于加入黑名单）
		tokenStr := r.Header.Get("Authorization")
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		}
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout(tokenStr)
		if err != nil {
			response.FailInternal(w, r)
		} else {
			response.OK(w, r)
		}
	}
}
