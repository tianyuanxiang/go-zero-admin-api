// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"go-zero-admin/pkg/encrypt"
	"go-zero-admin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"

	systemmodel "go-zero-admin/internal/model/system"
	"go-zero-admin/internal/svc"
	"go-zero-admin/internal/types"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateUser 创建新用户。
//
// 业务流程：
//  1. 检查用户名是否已存在
//  2. 对密码进行bcrypt加密
//  3. 插入用户记录
//  4. 分配初始角色（如果提供了roleIds）
//
// 参数：
//   - req : 创建用户请求体
//
// 返回：
//   - int64 : 新创建的用户ID
//   - error : 业务错误
func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) error {
	// 1. 检查用户名唯一性
	existUser, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound {
		l.Logger.Errorf("查询用户名[%s]是否存在失败：%v", req.Username, err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	if existUser != nil {
		return xerr.NewCodeError(xerr.ErrUsernameDuplicate)
	}
	// 2.对密码进行加密
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		l.Errorf("密码加密失败：%v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	// 默认状态为启用
	status := req.Status
	if status == 0 {
		status = 1
	}
	// 3.插入用户
	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.插入用户
		userId, err := l.svcCtx.SysUserModel.InsertUserTrans(l.ctx, tx, &systemmodel.SysUser{
			Username: req.Username,
			Password: hashedPassword,
			Nickname: req.Nickname,
			Email:    req.Email,
			Phone:    req.Phone,
			Avatar:   req.Avatar,
			Status:   int64(status),
			Remark:   req.Remark,
		})
		if err != nil {
			l.Errorf("插入用户记录失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 2. 分配初始角色
		if len(req.RoleIds) > 0 {
			if err = l.svcCtx.SysUserRoleModel.AssignRolesTrans(l.ctx, tx, userId, req.RoleIds); err != nil {
				l.Errorf("为新用户[%d]分配角色失败：%v", userId, err)
				// 角色分配失败不影响用户创建成功，仅记录日志
			}
		}
		return err
	})
	if err != nil {
		l.Logger.Errorf("创建用户事务执行失败: %v", err)
		return xerr.NewCodeError(xerr.ErrInternal)
	}
	return nil
}
