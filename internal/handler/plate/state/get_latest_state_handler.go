// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package state

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/state"
	"plating/internal/svc"
)

func GetLatestStateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := state.NewGetLatestStateLogic(r.Context(), svcCtx)
		resp, err := l.GetLatestState()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
