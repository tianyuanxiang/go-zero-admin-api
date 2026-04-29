// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-admin/internal/logic/system/file"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
)

func ListFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := file.NewListFileLogic(r.Context(), svcCtx)
		resp, err := l.ListFile(&req)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
