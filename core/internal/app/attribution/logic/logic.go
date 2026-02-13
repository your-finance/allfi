// Package logic 资产归因分析模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/attribution/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册资产归因分析服务
	service.RegisterAttribution(New())
}
