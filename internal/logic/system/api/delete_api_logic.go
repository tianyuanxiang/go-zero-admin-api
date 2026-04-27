// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package api

import (
	"context"
	"go-zero-admin/internal/logic/system/helper"
	"go-zero-admin/pkg/xerr"

	"go-zero-admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApiLogic) DeleteApi(apiId int64) error {
	_, err := l.svcCtx.SysApiModel.FindOne(l.ctx, apiId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Logger.Errorf("删除api时查询apiId[%d]是否存在失败:%v\n", apiId, err)
		return err
	}

	roleIds, err := l.svcCtx.SysRoleApiModel.ListRoleIdsByApiId(l.ctx, apiId)
	if err != nil {
		l.Logger.Errorf("查询API[%d]关联角色失败：%v", apiId, err)
	}
	// 1.开启删除事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.删除api
		if err := l.svcCtx.SysApiModel.SoftDeleteApiTrans(l.ctx, tx, apiId); err != nil {
			l.Logger.Errorf("删除api失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 2.删除role_api关联关系
		if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiByApiIdTrans(l.ctx, tx, apiId); err != nil {
			l.Logger.Errorf("删除role_api关联关系失败%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return err
	})
	if err != nil {
		l.Logger.Errorf("创建删除api事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	// 清除Casbin关联
	if err == nil {
		helper.RebuildCasbinByRoleIds(l.ctx, l.svcCtx, roleIds)
	}
	return err
}
