// Package model DeFi 仓位模块 - 数据传输对象
package model

// PositionItem DeFi 仓位条目
type PositionItem struct {
	Protocol   string  `json:"protocol"`    // 协议名称（Aave/Compound/Uniswap 等）
	Type       string  `json:"type"`        // 仓位类型（lending/staking/liquidity 等）
	Token      string  `json:"token"`       // 代币符号
	Amount     float64 `json:"amount"`      // 数量
	Value      float64 `json:"value"`       // 价值（计价货币）
	APY        float64 `json:"apy"`         // 年化收益率
	Chain      string  `json:"chain"`       // 所在链
	WalletAddr string  `json:"wallet_addr"` // 钱包地址
}

// ProtocolItem DeFi 协议条目
type ProtocolItem struct {
	Name     string   `json:"name"`      // 协议名称
	Chains   []string `json:"chains"`    // 支持的链列表
	Types    []string `json:"types"`     // 支持的仓位类型
	IsActive bool     `json:"is_active"` // 是否可用
}
