// Package defi DeFi API 定义
// 提供 DeFi 仓位查询和协议列表接口
package defi

import "github.com/gogf/gf/v2/frame/g"

// GetPositionsReq 获取 DeFi 仓位请求
type GetPositionsReq struct {
	g.Meta   `path:"/defi/positions" method:"get" summary:"获取 DeFi 仓位列表" tags:"DeFi"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// PositionItem DeFi 仓位条目
type PositionItem struct {
	Protocol    string  `json:"protocol" dc:"协议名称（Aave/Compound/Uniswap 等）"`
	Type        string  `json:"type" dc:"仓位类型（lending/staking/liquidity 等）"`
	Token       string  `json:"token" dc:"代币符号"`
	Amount      float64 `json:"amount" dc:"数量"`
	Value       float64 `json:"value" dc:"价值（计价货币）"`
	APY         float64 `json:"apy" dc:"年化收益率"`
	Chain       string  `json:"chain" dc:"所在链"`
	WalletAddr  string  `json:"wallet_addr" dc:"钱包地址"`
}

// GetPositionsRes 获取 DeFi 仓位响应
type GetPositionsRes struct {
	Positions  []PositionItem `json:"positions" dc:"仓位列表"`
	TotalValue float64        `json:"total_value" dc:"总价值"`
	Currency   string         `json:"currency" dc:"计价货币"`
}

// GetProtocolsReq 获取 DeFi 协议列表请求
type GetProtocolsReq struct {
	g.Meta `path:"/defi/protocols" method:"get" summary:"获取 DeFi 支持的协议列表" tags:"DeFi"`
}

// ProtocolItem DeFi 协议条目
type ProtocolItem struct {
	Name       string   `json:"name" dc:"协议名称"`
	Chains     []string `json:"chains" dc:"支持的链列表"`
	Types      []string `json:"types" dc:"支持的仓位类型"`
	IsActive   bool     `json:"is_active" dc:"是否可用"`
}

// GetProtocolsRes 获取 DeFi 协议列表响应
type GetProtocolsRes struct {
	Protocols []ProtocolItem `json:"protocols" dc:"协议列表"`
}

// GetStatsReq DeFi 统计请求
type GetStatsReq struct {
	g.Meta `path:"/defi/stats" method:"get" summary:"获取 DeFi 统计" tags:"DeFi"`
}

// GetStatsRes DeFi 统计响应
type GetStatsRes struct {
	TotalValueLocked float64            `json:"total_value_locked" dc:"总锁仓价值"`
	PositionCount    int                `json:"position_count" dc:"仓位数量"`
	ByProtocol       map[string]float64 `json:"by_protocol" dc:"按协议分组价值"`
	ByChain          map[string]float64 `json:"by_chain" dc:"按链分组价值"`
	ByType           map[string]float64 `json:"by_type" dc:"按类型分组价值"`
}
