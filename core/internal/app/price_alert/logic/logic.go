// Package logic 价格预警模块 - Logic 层导入文件
// 在包加载时自动注册价格预警服务实现
package logic

import (
	"your-finance/allfi/internal/app/price_alert/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册价格预警服务
	service.RegisterPriceAlert(New())
}
