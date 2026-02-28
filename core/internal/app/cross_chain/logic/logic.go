package logic

import (
	"your-finance/allfi/internal/integrations/bridge"
)

// New 创建跨链交易业务逻辑实例
func New() *sCrossChain {
	return &sCrossChain{
		bridgeManager: bridge.NewBridgeManager(),
	}
}
