// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/role"
	"go-zero-admin/internal/svc"
)

func GetRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		roleId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || roleId <= 0 {
			response.FailWithMsg(w, r, "角色ID格式错误")
			return
		}
		l := role.NewGetRoleLogic(r.Context(), svcCtx)
		resp, err := l.GetRole(roleId)
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
