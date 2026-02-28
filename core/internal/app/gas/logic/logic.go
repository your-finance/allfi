// Package logic Gas 优化业务逻辑
package logic

import (
	"your-finance/allfi/internal/app/gas/service"
)

func init() {
	service.RegisterGas(New())
}
