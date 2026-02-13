// Package service 用户模块服务接口
// 定义用户设置和缓存管理的服务接口
package service

import (
	"context"

	userApi "your-finance/allfi/api/v1/user"
)

// IUser 用户服务接口
type IUser interface {
	// GetSettings 获取用户设置
	GetSettings(ctx context.Context) (*userApi.GetSettingsRes, error)
	// UpdateSettings 更新用户设置
	UpdateSettings(ctx context.Context, settings map[string]string) error
	// ResetSettings 重置所有用户设置（删除 system_config 中 user.* 的记录）
	ResetSettings(ctx context.Context) error
	// ClearCache 清除缓存
	ClearCache(ctx context.Context) error
	// ExportData 导出用户数据（不含敏感信息）
	ExportData(ctx context.Context) (*userApi.ExportDataRes, error)
}

var localUser IUser

// User 获取用户服务实例
func User() IUser {
	if localUser == nil {
		panic("IUser 服务未注册，请检查 logic/user 包的 init 函数")
	}
	return localUser
}

// RegisterUser 注册用户服务实现
// 由 logic 层在 init 函数中调用
func RegisterUser(i IUser) {
	localUser = i
}
