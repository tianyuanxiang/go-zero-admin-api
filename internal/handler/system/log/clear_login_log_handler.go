// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/log"
	"go-zero-admin/internal/svc"
)

func ClearLoginLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		logId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || logId <= 0 {
			response.FailWithMsg(w, r, "登录日志ID格式错误")
			return
		}

		l := log.NewClearLoginLogLogic(r.Context(), svcCtx)
		err = l.ClearLoginLog(logId)
		if err != nil {
			if codeErr, ok := err.(*xerr.CodeError); ok {
				response.Fail(w, r, codeErr.Code, codeErr.Msg)
				return
			}
			response.FailInternal(w, r)
			return
		} else {
			response.OK(w, r)
		}
	}
}
