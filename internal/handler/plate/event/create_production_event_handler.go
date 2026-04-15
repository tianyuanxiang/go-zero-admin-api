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

func CreateProductionEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateProductionEventReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := event.NewCreateProductionEventLogic(r.Context(), svcCtx)
		err := l.CreateProductionEvent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
