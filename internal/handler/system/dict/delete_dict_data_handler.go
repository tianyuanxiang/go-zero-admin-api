// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/dict"
	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteDictDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		dictDataId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || dictDataId <= 0 {
			response.FailWithMsg(w, r, "字典数据ID格式错误")
			return
		}
		l := dict.NewDeleteDictDataLogic(r.Context(), svcCtx)
		err = l.DeleteDictData(dictDataId)
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
