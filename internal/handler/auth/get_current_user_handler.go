// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"

	"go-zero-admin/internal/logic/auth"
	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetCurrentUserLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentUser()
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
