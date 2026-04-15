// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/api"
	"plating/internal/svc"
	"plating/internal/types"
)

func CreateApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateApiReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := api.NewCreateApiLogic(r.Context(), svcCtx)
		err := l.CreateApi(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
