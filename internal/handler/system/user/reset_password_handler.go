// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"go-zero-admin/internal/logic/system/user"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 管理员重置别人的密码

func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.ResetPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}
		if req.Id <= 0 {
			response.FailWithMsg(w, r, "用户ID格式错误")
			return
		}
		if req.NewPassword == "" {
			response.FailWithMsg(w, r, "密码不能为空")
			return
		}

		l := user.NewResetPasswordLogic(r.Context(), svcCtx)
		err := l.ResetPassword(&req)
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
