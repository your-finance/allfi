// Package defi DeFi 仓位数据结构
package defi

// Position DeFi 仓位
// 统一描述用户在某个 DeFi 协议中的一个仓位
type Position struct {
	Protocol      string  `json:"protocol"`        // 协议标识（lido, aave, uniswap_v2）
	ProtocolName  string  `json:"protocol_name"`   // 协议显示名（Lido Finance）
	Type          string  `json:"type"`            // 仓位类型：staking/lending/lp/vault
	Chain         string  `json:"chain"`           // 所在链：ethereum/polygon/arbitrum 等
	DepositTokens []Token `json:"deposit_tokens"`  // 质押/存入的代币
	ReceiveTokens []Token `json:"receive_tokens"`  // 获得的凭证代币（stETH/rETH/aUSDC）
	ValueUSD      float64 `json:"value_usd"`       // 当前仓位价值（USD）
	APY           float64 `json:"apy"`             // 年化收益率（百分比，如 3.5 表示 3.5%）
	Rewards       []Token `json:"rewards"`         // 未领取的奖励代币

	// 借贷相关字段（仅 Type="lending" 时有效）
	BorrowTokens         []Token `json:"borrow_tokens,omitempty"`          // 借出的代币
	HealthFactor         float64 `json:"health_factor,omitempty"`          // 健康因子（>1 安全，<1 可能被清算）
	LiquidationThreshold float64 `json:"liquidation_threshold,omitempty"`  // 清算阈值（百分比）
	LTV                  float64 `json:"ltv,omitempty"`                    // Loan-to-Value 比率（百分比）
	NetAPY               float64 `json:"net_apy,omitempty"`                // 净收益率（存款APY - 借款APY）
}

// Token 代币信息
type Token struct {
	Symbol   string  `json:"symbol"`    // 代币符号（ETH, stETH, USDC）
	Amount   float64 `json:"amount"`    // 数量
	ValueUSD float64 `json:"value_usd"` // USD 价值
}
