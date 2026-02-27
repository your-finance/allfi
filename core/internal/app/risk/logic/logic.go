// Package logic 风险管理模块 - 服务注册
package logic

import (
	"your-finance/allfi/internal/app/risk/service"
)

func init() {
	service.RegisterRisk(New())
}
