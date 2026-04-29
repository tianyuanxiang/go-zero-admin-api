// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-admin/internal/logic/system/role"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
)

func UpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := role.NewUpdateRoleLogic(r.Context(), svcCtx)
		err := l.UpdateRole(&req)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OK(w, r)
		}
	}
}
