// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tank

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"plating/internal/logic/plate/tank"
	"plating/internal/svc"
)

func GetTankConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := tank.NewGetTankConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetTankConfig()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
