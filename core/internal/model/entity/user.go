package entity

import (
	"time"
)

// User 用户实体
// 存储系统用户信息，支持多用户
type User struct {
	Id           uint           `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Username     string         `gorm:"size:50;not null;uniqueIndex:idx_username;comment:用户名" json:"username"`
	Email        string         `gorm:"size:100;not null;uniqueIndex:idx_email;comment:邮箱" json:"email"`
	PasswordHash string         `gorm:"size:255;not null;comment:密码哈希" json:"-"`
	Nickname     string         `gorm:"size:50;comment:昵称" json:"nickname"`
	Avatar       string         `gorm:"size:255;comment:头像URL" json:"avatar"`
	Status       string         `gorm:"size:20;not null;default:active;comment:状态(active/disabled/locked)" json:"status"`
	LastLoginAt  *time.Time     `gorm:"comment:最后登录时间" json:"last_login_at"`
	LastLoginIp  string         `gorm:"size:45;comment:最后登录IP" json:"last_login_ip"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
