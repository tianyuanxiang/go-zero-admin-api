// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/api"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func UpdateApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		apiId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || apiId <= 0 {
			response.FailWithMsg(w, r, "api ID格式错误")
			return
		}
		var req types.UpdateApiReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		req.Id = apiId
		l := api.NewUpdateApiLogic(r.Context(), svcCtx)
		err = l.UpdateApi(&req)
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
		} else {
			response.OK(w, r)
		}
	}
}
