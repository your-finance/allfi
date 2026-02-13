// Package logic 基准对比模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/benchmark/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册基准对比服务
	service.RegisterBenchmark(New())
}
