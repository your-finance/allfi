// Package logic 策略引擎模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/strategy/service"
)

// init 在包加载时自动注册策略服务
func init() {
	service.RegisterStrategy(New())
}
