// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package file

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"go-zero-admin/internal/logic/system/file"
	"go-zero-admin/internal/svc"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewDeleteFileLogic(r.Context(), svcCtx)
		err := l.DeleteFile()
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			response.OK(w, r)
		}
	}
}
