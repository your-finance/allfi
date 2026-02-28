package logic

import (
	"your-finance/allfi/internal/app/cross_chain/service"
	"your-finance/allfi/internal/integrations/bridge"
)

func init() {
	service.RegisterCrossChain(New())
}

// New 创建跨链交易业务逻辑实例
func New() *sCrossChain {
	return &sCrossChain{
		bridgeManager: bridge.NewBridgeManager(),
	}
}
