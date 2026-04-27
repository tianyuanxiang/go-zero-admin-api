// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"
	"plating/pkg/response"
	"plating/pkg/xerr"

	"plating/internal/logic/system/role"
	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRolePermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := role.NewUpdateRolePermissionsLogic(r.Context(), svcCtx)
		err := l.UpdateRolePermissions(&req)
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
			return
		} else {
			response.OK(w, r)
		}
	}
}
