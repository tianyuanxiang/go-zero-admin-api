// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"
	"go-zero-admin/internal/model/system"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApiLogic) CreateApi(req *types.CreateApiReq) error {
	// 检查路径+方法组合唯一性
	exist, err := l.svcCtx.SysApiModel.FindOneByApiPathMethod(l.ctx, req.ApiPath, req.Method)
	if err != nil && err != sqlx.ErrNotFound {
		l.Errorf("查询接口[%s %s]是否存在失败：%v", req.Method, req.ApiPath, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if exist != nil {
		return xerr.NewCodeErrorMsg(xerr.ErrDuplicate, "该路径和方法组合已存在")
	}

	_, err = l.svcCtx.SysApiModel.Insert(l.ctx, &system.SysApi{
		ApiPath:     req.ApiPath,
		ApiName:     req.ApiName,
		Method:      req.Method,
		ApiGroup:    req.Group,
		Description: req.Remark,
	})
	if err != nil {
		l.Errorf("插入接口记录失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return nil
}
