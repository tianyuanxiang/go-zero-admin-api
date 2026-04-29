// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/menu"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		menuId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || menuId <= 0 {
			response.FailWithMsg(w, r, "菜单ID格式错误")
			return
		}
		var req types.UpdateMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		req.Id = menuId
		l := menu.NewUpdateMenuLogic(r.Context(), svcCtx)
		err = l.UpdateMenu(&req)
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
