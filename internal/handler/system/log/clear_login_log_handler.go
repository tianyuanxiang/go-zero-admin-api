// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/log"
	"plating/internal/svc"
)

func ClearLoginLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := log.NewClearLoginLogLogic(r.Context(), svcCtx)
		err := l.ClearLoginLog()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
