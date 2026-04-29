// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"go-zero-admin/pkg/response"
	"go-zero-admin/pkg/xerr"
	"net/http"
	"strconv"

	"go-zero-admin/internal/logic/system/user"
	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		userId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || userId <= 0 {
			response.FailWithMsg(w, r, "用户ID格式错误")
			return
		}

		l := user.NewDeleteUserLogic(r.Context(), svcCtx)
		err = l.DeleteUser(userId)
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
