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
)

func DeleteApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.PathValue("id")
		apiId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || apiId <= 0 {
			response.FailWithMsg(w, r, "ApiID格式错误")
			return
		}
		l := api.NewDeleteApiLogic(r.Context(), svcCtx)
		err = l.DeleteApi(apiId)
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailWithMsg(w, r, err.Error())
			return
		} else {
			response.OK(w, r)
		}
	}
}
