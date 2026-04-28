// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-admin/internal/logic/system/log"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
)

func ListOperLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListOperLogReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}
		l := log.NewListOperLogLogic(r.Context(), svcCtx)
		resp, err := l.ListOperLog(&req)
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
			return
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
