// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package log

import (
	"context"
	"go-zero-admin/internal/model/system"
	"go-zero-admin/pkg/dict"
	"go-zero-admin/pkg/xerr"
	"time"

	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperLogLogic {
	return &ListOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOperLogLogic) ListOperLog(req *types.ListOperLogReq) (resp *types.ListOperLogResp, err error) {
	listReq := &system.OperLogListReq{
		Page:         req.Page,
		PageSize:     req.PageSize,
		Keyword:      req.Keyword,
		BusinessType: req.BusinessType,
		Status:       req.Status,
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

	logs, total, err := l.svcCtx.SysOperLogModel.List(l.ctx, listReq)
	if err != nil {
		l.Errorf("查询操作日志失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]types.OperLogItem, 0, len(logs))
	for _, log := range logs {
		list = append(list, types.OperLogItem{
			Id:         int(log.Id),
			Title:      log.Title,
			OperType:   dict.OperNumToType[log.BusinessType],
			Method:     log.Method,
			ReqMethod:  log.RequestMethod,
			OperName:   log.OperatorName,
			DeptName:   log.DeptName,
			ReqUrl:     log.OperUrl,
			ReqParam:   log.OperParam.String,
			RespResult: log.JsonResult.String,
			Status:     int(log.Status),
			Ip:         log.OperIp,
			OperTime:   log.OperTime.Format(timeLayout),
		})
	}

	return &types.ListOperLogResp{
		Total: total,
		List:  list,
	}, nil
}
