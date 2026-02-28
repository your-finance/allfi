package logic

import (
	"your-finance/allfi/internal/app/cross_chain/service"
)

func init() {
	service.RegisterCrossChain(New())
}
