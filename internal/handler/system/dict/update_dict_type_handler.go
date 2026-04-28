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
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateDictTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		dictTypeId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || dictTypeId <= 0 {
			response.FailWithMsg(w, r, "字典类型ID格式错误")
			return
		}
		var req types.UpdateDictTypeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}
		req.Id = dictTypeId
		l := dict.NewUpdateDictTypeLogic(r.Context(), svcCtx)
		err = l.UpdateDictType(&req)
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
