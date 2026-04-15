package system

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysUserModel = (*customSysUserModel)(nil)

type (
	SysUserModel interface {
		sysUserModel
		List(ctx context.Context, page, pageSize int, keyword string, status int) ([]*SysUser, int64, error)
		UpdatePassword(ctx context.Context, id int64, password string) error
		InsertUser(ctx context.Context, tx *gorm.DB, user *SysUser) (int64, error)
	}

	customSysUserModel struct {
		*defaultSysUserModel
		db *gorm.DB
	}
)

func NewSysUserModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysUserModel {
	return &customSysUserModel{
		defaultSysUserModel: newSysUserModel(conn, c, opts...),
		db:                  db,
	}
}

// UpdatePassword 仅更新用户密码字段。
func (m *customSysUserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	result := m.db.WithContext(ctx).Table("sys_user").
		Where("id = ? AND deleted_at IS NULL", id).
		Update("password", password)
	return result.Error
}

// FindOneByUsername 按用户名查询未软删除的用户。
//
// 覆盖生成代码的同名方法，原因：
//  1. 生成代码使用双层索引缓存，若数据直接写入数据库而未经 Insert 方法，
//     Redis 中可能存有负向缓存（值为"*"），导致返回 ErrNotFound。
//  2. 生成代码的 SQL 缺少 deleted_at IS NULL 过滤。
func (m *customSysUserModel) FindOneByUsername(ctx context.Context, username string) (*SysUser, error) {
	var user SysUser
	result := m.db.WithContext(ctx).Table("sys_user").
		Where("username = ? AND deleted_at IS NULL", username).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// 插入用户返回id
func (m *customSysUserModel) InsertUser(ctx context.Context, tx *gorm.DB, user *SysUser) (int64, error) {
	result := tx.WithContext(ctx).Table("sys_user").Create(&user)
	return result.RowsAffected, result.Error
}

func (m *customSysUserModel) List(ctx context.Context, page, pageSize int, keyword string, status int) ([]*SysUser, int64, error) {
	// 构建查询条件
	db := m.db.WithContext(ctx).Table("sys_user").Where("deleted_at IS NULL")
	// 关键词模糊检索
	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 状态筛选
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	var users []*SysUser
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
