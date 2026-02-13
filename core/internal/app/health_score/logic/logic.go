// Package logic 资产健康评分模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/health_score/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册健康评分服务
	service.RegisterHealthScore(New())
}
