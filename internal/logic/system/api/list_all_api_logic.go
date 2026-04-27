// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllApiLogic {
	return &ListAllApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllApiLogic) ListAllApi() (resp *types.ListAllApiResp, err error) {
	apis, err := l.svcCtx.SysApiModel.ListByIds(l.ctx, nil)
	if err != nil {
		l.Errorf("查询全部角色失败: %v\n", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]types.ApiOption, 0, len(apis))
	for _, api := range apis {
		list = append(list, types.ApiOption{
			Id:      api.Id,
			ApiName: api.ApiName,
			ApiPath: api.ApiPath,
			Remark:  api.Description,
		})
	}
	return &types.ListAllApiResp{
		ListAll: list,
	}, err
}
