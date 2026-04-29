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

func DeleteDictTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		dictTypeId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || dictTypeId <= 0 {
			response.FailWithMsg(w, r, "字典类型ID格式错误")
			return
		}

		l := dict.NewDeleteDictTypeLogic(r.Context(), svcCtx)

		err = l.DeleteDictType(dictTypeId)
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
