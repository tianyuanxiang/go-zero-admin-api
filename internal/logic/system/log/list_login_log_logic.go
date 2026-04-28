// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"
	"go-zero-admin/internal/model/system"
	"go-zero-admin/pkg/xerr"
	"time"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginLogLogic {
	return &ListLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoginLogLogic) ListLoginLog(req *types.ListLoginLogReq) (resp *types.ListLoginLogResp, err error) {
	listReq := &system.LoginLogListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	}

	// 解析时间范围（可选参数）
	const timeLayout = "2006-01-02 15:04:05"
	if req.StartTime != "" {
		t, err := time.ParseInLocation(timeLayout, req.StartTime, time.Local)
		if err == nil {
			listReq.StartTime = t
		}
	}
	if req.EndTime != "" {
		t, err := time.ParseInLocation(timeLayout, req.EndTime, time.Local)
		if err == nil {
			listReq.EndTime = t
		}
	}

	logs, total, err := l.svcCtx.SysLoginLogModel.List(l.ctx, listReq)
	if err != nil {
		l.Errorf("查询登录日志失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]types.LoginLogItem, 0, len(logs))
	for _, log := range logs {
		list = append(list, types.LoginLogItem{
			Id:        int(log.Id),
			UserId:    int(log.UserId),
			Username:  log.Username,
			Ip:        log.Ip,
			Location:  log.Location,
			UserAgent: log.Browser,
			OS:        log.Os,
			Status:    int(log.Status),
			Msg:       log.Msg,
			LoginTime: log.LoginTime.Format(timeLayout),
		})
	}

	return &types.ListLoginLogResp{
		Total: int(total),
		List:  list,
	}, nil
}
