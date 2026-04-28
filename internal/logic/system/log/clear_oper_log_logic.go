// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ClearOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearOperLogLogic {
	return &ClearOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearOperLogLogic) ClearOperLog(logId int64) error {
	_, err := l.svcCtx.SysOperLogModel.FindOne(l.ctx, logId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Errorf("删除操作日志时查询logId[%d]是否存在失败:%v\n", logId, err)
		return err
	}

	if err := l.svcCtx.SysOperLogModel.Delete(l.ctx, logId); err != nil {
		l.Errorf("删除操作日志失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return err
}
