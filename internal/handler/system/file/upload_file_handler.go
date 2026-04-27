// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/system/file"
	"plating/internal/svc"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile()
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OKJsonCtx(r.Context(), w, resp)
		}
	}
}
