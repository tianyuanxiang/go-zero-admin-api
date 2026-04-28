// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"go-zero-admin/pkg/response"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/log"
	"go-zero-admin/internal/svc"
)

func ClearOperLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		logId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || logId <= 0 {
			response.FailWithMsg(w, r, "操作日志ID格式错误")
			return
		}

		l := log.NewClearOperLogLogic(r.Context(), svcCtx)
		err = l.ClearOperLog(logId)
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OK(w, r)
		}
	}
}
