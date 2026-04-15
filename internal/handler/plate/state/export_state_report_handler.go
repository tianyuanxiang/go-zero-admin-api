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

func ExportStateReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportStateReportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := state.NewExportStateReportLogic(r.Context(), svcCtx)
		err := l.ExportStateReport(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
