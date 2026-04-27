// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"

	"go-zero-admin/internal/logic/system/role"
	"go-zero-admin/internal/svc"
)

func ListAllRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewListAllRoleLogic(r.Context(), svcCtx)
		resp, err := l.ListAllRole()
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
			return
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
