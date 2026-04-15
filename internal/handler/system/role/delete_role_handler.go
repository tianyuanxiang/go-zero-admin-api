// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/role"
	"plating/internal/svc"
)

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewDeleteRoleLogic(r.Context(), svcCtx)
		err := l.DeleteRole()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
