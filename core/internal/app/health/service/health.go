// =================================================================================
// 健康检查服务接口定义
// 提供服务健康状态检查能力，包括数据库连接、版本信息等
// =================================================================================

package service

import (
	"context"
)

// IHealth 健康检查服务接口
type IHealth interface {
	// GetHealthStatus 获取服务健康状态
	// 检查数据库连接并返回版本信息
	GetHealthStatus(ctx context.Context) (map[string]interface{}, error)
}

// localHealth 健康检查服务实例（延迟注入）
var localHealth IHealth

// Health 获取健康检查服务实例
// 如果服务未注册，会触发 panic
func Health() IHealth {
	if localHealth == nil {
		panic("IHealth 服务未注册，请检查 logic/health 包的 init 函数")
	}
	return localHealth
}

// RegisterHealth 注册健康检查服务实现
// 由 logic 层在 init 函数中调用
func RegisterHealth(i IHealth) {
	localHealth = i
}
