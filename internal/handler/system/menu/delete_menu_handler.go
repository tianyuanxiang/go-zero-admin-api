// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"
	"plating/pkg/response"
	"plating/pkg/xerr"
	"strconv"

	"plating/internal/logic/system/menu"
	"plating/internal/svc"
)

func DeleteMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		menuId, err := strconv.ParseInt(idStr, 10, 64)

		if err != nil || menuId <= 0 {
			response.FailWithMsg(w, r, "菜单ID格式错误")
			return
		}
		l := menu.NewDeleteMenuLogic(r.Context(), svcCtx)
		err = l.DeleteMenu(menuId)
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
