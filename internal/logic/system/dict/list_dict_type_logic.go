// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictTypeLogic) ListDictType(req *types.ListDictTypeReq) (resp *types.ListDictTypeResp, err error) {
	dictTypes, count, err := l.svcCtx.SysDictTypeModel.List(l.ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		l.Errorf("查询字典类型列表失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]types.DictTypeItem, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		list = append(list, types.DictTypeItem{
			Id:        int(dictType.Id),
			DictName:  dictType.Name,
			DictType:  dictType.Code,
			Remark:    dictType.Remark,
			Status:    int(dictType.Status),
			CreatedAt: dictType.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &types.ListDictTypeResp{
		Total: int(count),
		List:  list,
	}, err
}
