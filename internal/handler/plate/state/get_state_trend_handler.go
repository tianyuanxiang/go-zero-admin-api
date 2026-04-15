// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package state

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/state"
	"plating/internal/svc"
	"plating/internal/types"
)

func GetStateTrendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetStateTrendReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := state.NewGetStateTrendLogic(r.Context(), svcCtx)
		resp, err := l.GetStateTrend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
