// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/log"
	"plating/internal/svc"
	"plating/internal/types"
)

func ListOperLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListOperLogReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := log.NewListOperLogLogic(r.Context(), svcCtx)
		resp, err := l.ListOperLog(&req)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OKJsonCtx(r.Context(), w, resp)
		}
	}
}
