// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package state

import (
	"context"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportStateReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportStateReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportStateReportLogic {
	return &ExportStateReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportStateReportLogic) ExportStateReport(req *types.ExportStateReportReq) error {
	// todo: add your logic here and delete this line

	return nil
}
