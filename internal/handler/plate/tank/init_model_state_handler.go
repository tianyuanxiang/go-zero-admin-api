// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/tank"
	"plating/internal/svc"
	"plating/internal/types"
)

func InitModelStateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InitModelStateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tank.NewInitModelStateLogic(r.Context(), svcCtx)
		err := l.InitModelState(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
