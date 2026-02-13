// Package logic 健康检查业务逻辑
// 检查数据库连接状态并返回服务版本信息
package logic

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/health/service"
)

// sHealth 健康检查服务实现
type sHealth struct{}

// New 创建健康检查服务实例
func New() service.IHealth {
	return &sHealth{}
}

// GetHealthStatus 获取服务健康状态
//
// 功能说明:
// 1. 检查数据库连接状态
// 2. 返回服务版本信息
// 3. 返回运行时配置
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - map[string]interface{}: 健康状态数据
//   - error: 错误信息
func (s *sHealth) GetHealthStatus(ctx context.Context) (map[string]interface{}, error) {
	// 从配置读取版本信息
	appName := g.Cfg().MustGet(ctx, "app.name", "AllFi").String()
	appVersion := g.Cfg().MustGet(ctx, "app.version", "0.1.0").String()
	appEnv := g.Cfg().MustGet(ctx, "app.env", "development").String()

	// 检查数据库连接
	dbStatus := "ok"
	db := g.DB()
	if db == nil {
		dbStatus = "not_initialized"
	} else {
		// 尝试 ping 数据库
		if err := db.PingMaster(); err != nil {
			g.Log().Error(ctx, "数据库连接失败", "error", err)
			dbStatus = "error"
		}
	}

	// 构造健康状态响应
	status := map[string]interface{}{
		"status":  "ok",
		"name":    appName,
		"version": appVersion,
		"env":     appEnv,
		"database": map[string]string{
			"status": dbStatus,
		},
	}

	return status, nil
}
