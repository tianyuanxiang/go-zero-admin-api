// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"
	"plating/pkg/response"

	"plating/internal/logic/system/menu"
	"plating/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前用户有权限的菜单树
func GetCurrentUserMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetCurrentUserMenusLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentUserMenus()
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
