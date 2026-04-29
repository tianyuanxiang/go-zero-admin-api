// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"go-zero-admin/internal/logic/auth"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := auth.NewChangePasswordLogic(r.Context(), svcCtx)
		err := l.ChangePassword(&req)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OK(w, r)
		}
	}
}
