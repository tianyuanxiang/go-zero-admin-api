// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package event

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/event"
	"plating/internal/svc"
)

func DeleteProductionEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := event.NewDeleteProductionEventLogic(r.Context(), svcCtx)
		err := l.DeleteProductionEvent()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
