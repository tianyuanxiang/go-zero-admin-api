// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"plating/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	systemmodel "plating/internal/model/system"
	"plating/internal/svc"
	"plating/internal/types"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) error {
	// 检查用户是否存在
	_, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户失败，用户ID[%d]：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	// 更新用户信息
	if err = l.svcCtx.SysUserModel.Update(l.ctx, &systemmodel.SysUser{
		Id:       req.Id,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   int64(req.Status),
		Remark:   req.Remark,
	}); err != nil {
		l.Errorf("更新用户[%d]信息失败：%v", req.Id, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
