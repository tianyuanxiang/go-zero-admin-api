// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"go-zero-admin/internal/logic/system/dict"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListDictTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListDictTypeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := dict.NewListDictTypeLogic(r.Context(), svcCtx)
		resp, err := l.ListDictType(&req)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
