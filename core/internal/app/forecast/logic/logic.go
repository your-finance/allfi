// Package logic 趋势预测模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/forecast/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册趋势预测服务
	service.RegisterForecast(New())
}
