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

type ListApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApiLogic) ListApi(req *types.ListApiReq) (resp *types.ListApiResp, err error) {
	apis, count, err := l.svcCtx.SysApiModel.List(l.ctx, req.Page, req.PageSize, req.Group, req.Keyword)
	if err != nil {
		l.Errorf("查询接口列表失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]types.ApiItem, 0, len(apis))
	for _, api := range apis {
		list = append(list, types.ApiItem{
			Id:        api.Id,
			ApiName:   api.ApiName,
			ApiPath:   api.ApiPath,
			Method:    api.Method,
			Group:     api.ApiGroup,
			Remark:    api.Description,
			CreatedAt: api.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &types.ListApiResp{
		Total: count,
		List:  list,
	}, err
}
