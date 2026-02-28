package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 注册跨链交易路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(
		CrossChain.GetTransactions,
		CrossChain.GetAssetFlow,
		CrossChain.GetFeeStats,
		CrossChain.GetBridges,
	)
}
