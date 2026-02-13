// Package model 认证模块 - 数据传输对象 (DTO)
package model

// 认证相关常量
const (
	// ConfigKeyPINHash PIN 哈希存储键
	ConfigKeyPINHash = "auth.pin_hash"
	// ConfigKeyJWTSecret JWT 密钥存储键
	ConfigKeyJWTSecret = "auth.jwt_secret"
	// ConfigKeyFailCount 连续失败次数存储键
	ConfigKeyFailCount = "auth.fail_count"
	// ConfigKeyLockUntil 锁定截止时间存储键
	ConfigKeyLockUntil = "auth.lock_until"

	// PINMinLength PIN 最小长度
	PINMinLength = 4
	// PINMaxLength PIN 最大长度
	PINMaxLength = 20
	// MaxFailCount 最大连续失败次数
	MaxFailCount = 5
	// LockDurationMinutes 锁定时间（分钟）
	LockDurationMinutes = 15
	// TokenExpireHours JWT Token 有效期（小时）
	TokenExpireHours = 24
)
