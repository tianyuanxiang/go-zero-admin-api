// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/tank"
	"plating/internal/svc"
)

func DeleteTankConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := tank.NewDeleteTankConfigLogic(r.Context(), svcCtx)
		err := l.DeleteTankConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
