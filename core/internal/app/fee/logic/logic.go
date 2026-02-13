// Package logic 费用分析模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/fee/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册费用分析服务
	service.RegisterFee(New())
}
