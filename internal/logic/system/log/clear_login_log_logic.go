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

type ClearLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearLoginLogLogic {
	return &ClearLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearLoginLogLogic) ClearLoginLog(logId int64) error {

	_, err := l.svcCtx.SysLoginLogModel.FindOne(l.ctx, logId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Errorf("删除登录日志时查询logId[%d]是否存在失败:%v\n", logId, err)
		return err
	}

	if err := l.svcCtx.SysLoginLogModel.Delete(l.ctx, logId); err != nil {
		l.Errorf("删除登录日志失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}

	return err
}
