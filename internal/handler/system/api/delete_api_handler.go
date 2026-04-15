// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/api"
	"plating/internal/svc"
)

func DeleteApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := api.NewDeleteApiLogic(r.Context(), svcCtx)
		err := l.DeleteApi()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
