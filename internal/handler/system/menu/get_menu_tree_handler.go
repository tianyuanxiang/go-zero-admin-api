// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"go-zero-admin/pkg/response"
	"net/http"

	"go-zero-admin/internal/logic/system/menu"
	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetMenuTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetMenuTree()
		if err != nil {
			response.FailWithMsg(w, r, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
