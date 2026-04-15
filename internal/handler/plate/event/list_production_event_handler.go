// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/event"
	"plating/internal/svc"
	"plating/internal/types"
)

func ListProductionEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := event.NewListProductionEventLogic(r.Context(), svcCtx)
		resp, err := l.ListProductionEvent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
