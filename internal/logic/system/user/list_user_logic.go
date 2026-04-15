// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"plating/pkg/xerr"

	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLogic) ListUser(req *types.ListUserReq) (resp *types.ListUserResp, err error) {
	users, total, err := l.svcCtx.SysUserModel.List(l.ctx, req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 收集所有用户id
	userIds := make([]int64, 0, len(users))
	for _, user := range users {
		userIds = append(userIds, user.Id)
	}
	userRoleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserIds(l.ctx, userIds)
	if err != nil {
		l.Logger.Errorf("批量查询用户关联角色失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 再去收集所有角色ID
	// 去重
	roleIdSet := make(map[int64]struct{})
	for _, ur := range userRoleIds {
		roleIdSet[ur.RoleId] = struct{}{}
	}

	roleIds := make([]int64, 0, len(roleIdSet))
	for id := range roleIdSet {
		roleIds = append(roleIds, id)
	}

	// 4. 批量查询所有角色信息
	roleMap := make(map[int64]string) // roleId -> code
	if len(roleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
		if err != nil {
			l.Logger.Errorf("批量查询角色信息失败: %v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		for _, role := range roles {
			roleMap[role.Id] = role.Code
		}
	}

	// 5. 构建 userId -> []code 的映射
	userRoleMap := make(map[int64][]string)
	for _, ur := range userRoleIds {
		if code, ok := roleMap[ur.RoleId]; ok {
			userRoleMap[ur.UserId] = append(userRoleMap[ur.UserId], code)
		}
	}

	// 6. 组装返回数据
	list := make([]types.UserItem, 0, len(users))
	for _, user := range users {
		list = append(list, types.UserItem{
			Id:        user.Id,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			Status:    int(user.Status),
			Roles:     userRoleMap[user.Id],
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &types.ListUserResp{
		Total: total,
		List:  list,
	}, nil
}
