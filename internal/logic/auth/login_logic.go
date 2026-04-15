// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"net/http"
	"plating/internal/common"
	"plating/internal/middleware"
	"plating/pkg/encrypt"
	"plating/pkg/jwtx"
	"plating/pkg/xerr"
	"time"

	systemmodel "plating/internal/model/system"
	"plating/internal/svc"
	"plating/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 处理用户登录请求。
//
/* 业务流程：
1. 校验用户名是否存在
2. 校验账号状态是否正常
3. 校验密码是否正确
4. 生成JWT双Token（Access + Refresh）
5. 查询用户角色列表
6. 查询用户可访问的菜单树
7. 记录登录日志
8. 返回登录结果
*/
func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {
	// 1. 按用户名查询用户
	user, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if err == systemmodel.ErrNotFound {
			l.recordLoginLog(req.Username, 0, r, 0, "用户名或密码错误")
			return nil, xerr.NewCodeError(xerr.ErrPasswordWrong)
		}
		l.Logger.Errorf("查询用户[%s]失败：%v", req.Username, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 2. 检查账号是否被禁用
	if user.Status != 1 {
		l.recordLoginLog(req.Username, user.Id, r, 0, "账号已被禁用")
		return nil, xerr.NewCodeError(xerr.ErrAccountDisabled)
	}

	// 3. 校验密码（bcrypt比对）
	if !encrypt.CheckPassword(req.Password, user.Password) {
		l.recordLoginLog(req.Username, user.Id, r, 0, "密码错误")
		return nil, xerr.NewCodeError(xerr.ErrPasswordWrong)
	}

	// 4. 生成访问令牌和刷新令牌
	accessToken, err := jwtx.GenerateToken(user.Id, user.Username, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Logger.Errorf("生成AccessToken失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	refreshToken, err := jwtx.GenerateRefreshToken(user.Id, user.Username, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.RefreshExpire)
	if err != nil {
		l.Logger.Errorf("生成RefreshToken失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 5. 查询用户拥有的角色列表
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, user.Id)
	if err != nil {
		l.Logger.Errorf("查询用户[%s]的角色信息失败：%v", user.Username, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	roleCodes := make([]string, len(roleIds)+1)
	for _, roleId := range roleIds {
		role, roleErr := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if roleErr != nil || role == nil || role.Status != 1 {
			continue
		}
		roleCodes = append(roleCodes, role.Code)
	}

	// 6. 查询用户的菜单权限（合并所有角色的菜单）
	menuIds, err := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, roleIds)
	if err != nil {
		l.Logger.Errorf("查询用户菜单失败：%v", err)
		menuIds = []int64{}
	}

	menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, menuIds)
	if err != nil {
		l.Logger.Errorf("查询菜单详情失败：%v", err)
		menus = []*systemmodel.SysMenu{}
	}

	// 将菜单列表构建为树形结构
	menuTree := buildMenuTree(menus, 0)

	// 7. 记录成功登录日志
	l.recordLoginLog(req.Username, user.Id, r, 1, "登录成功")

	return &types.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
		UserInfo: types.UserInfo{
			UserId:   user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Roles:    roleCodes,
			Menus:    menuTree,
		},
	}, nil
}

// recordLoginLog 记录登录日志（异步写入，不影响主流程）。
func (l *LoginLogic) recordLoginLog(username string, userId int64, r *http.Request, status int64, msg string) {
	go func() {
		clientIP := middleware.GetClientIP(r)
		userAgent := r.Header.Get("User-Agent")
		browser := common.ExtractBrowser(userAgent)
		osName := common.ExtractOS(userAgent)

		loginLog := &systemmodel.SysLoginLog{
			UserId:    userId,
			Username:  username,
			Ip:        clientIP,
			Location:  "",
			Browser:   browser,
			Os:        osName,
			Status:    status,
			Msg:       msg,
			LoginTime: time.Now(),
		}

		if _, insertErr := l.svcCtx.SysLoginLogModel.Insert(l.ctx, loginLog); insertErr != nil {
			l.Logger.Errorf("记录登录日志失败：%v", insertErr)
		}
	}()
}

// buildMenuTree 将扁平菜单列表递归构建为树形结构。
func buildMenuTree(menus []*systemmodel.SysMenu, parentId int64) []types.MenuItem {
	result := make([]types.MenuItem, 0)
	for _, m := range menus {
		if m.ParentId != parentId {
			continue
		}
		menuNode := types.MenuItem{
			Id:        m.Id,
			ParentId:  m.ParentId,
			MenuName:  m.Name,
			Path:      m.Path,
			Component: m.Component,
			Icon:      m.Icon,
			MenuType:  m.Type,
			Perms:     m.Permission,
			Sort:      m.Sort,
		}
		children := buildMenuTree(menus, m.Id)
		if len(children) > 0 {
			menuNode.Children = children
		}
		result = append(result, menuNode)
	}
	return result
}
