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

// LendingPositionItem 借贷仓位条目
type LendingPositionItem struct {
	Protocol             string  `json:"protocol"`               // 协议名称（Aave/Compound）
	Chain                string  `json:"chain"`                  // 所在链
	WalletAddr           string  `json:"wallet_addr"`            // 钱包地址
	SupplyToken          string  `json:"supply_token"`           // 存款代币
	SupplyAmount         float64 `json:"supply_amount"`          // 存款数量
	SupplyValueUSD       float64 `json:"supply_value_usd"`       // 存款价值（USD）
	SupplyAPY            float64 `json:"supply_apy"`             // 存款年化收益率
	BorrowToken          string  `json:"borrow_token,omitempty"` // 借款代币
	BorrowAmount         float64 `json:"borrow_amount"`          // 借款数量
	BorrowValueUSD       float64 `json:"borrow_value_usd"`       // 借款价值（USD）
	BorrowAPY            float64 `json:"borrow_apy"`             // 借款年化利率
	HealthFactor         float64 `json:"health_factor"`          // 健康因子
	LiquidationThreshold float64 `json:"liquidation_threshold"`  // 清算阈值
	LTV                  float64 `json:"ltv"`                    // Loan-to-Value 比率
	NetAPY               float64 `json:"net_apy"`                // 净收益率
}

// LendingRateItem 借贷利率条目
type LendingRateItem struct {
	Protocol          string  `json:"protocol"`            // 协议名称
	Chain             string  `json:"chain"`               // 所在链
	Token             string  `json:"token"`               // 代币符号
	SupplyAPY         float64 `json:"supply_apy"`          // 存款年化收益率
	BorrowAPYStable   float64 `json:"borrow_apy_stable"`   // 稳定借款利率
	BorrowAPYVariable float64 `json:"borrow_apy_variable"` // 浮动借款利率
	UtilizationRate   float64 `json:"utilization_rate"`    // 利用率
}

// LendingRateHistoryItem 借贷利率历史条目
type LendingRateHistoryItem struct {
	Date              string  `json:"date"`                // 日期
	SupplyAPY         float64 `json:"supply_apy"`          // 存款年化收益率
	BorrowAPYStable   float64 `json:"borrow_apy_stable"`   // 稳定借款利率
	BorrowAPYVariable float64 `json:"borrow_apy_variable"` // 浮动借款利率
	UtilizationRate   float64 `json:"utilization_rate"`    // 利用率
}

// LendingOptimizationResult 借贷策略优化结果
type LendingOptimizationResult struct {
	CurrentPositions []*LendingPositionItem       `json:"current_positions"` // 当前仓位
	Recommendations  []*LendingRecommendation     `json:"recommendations"`   // 优化建议
	PotentialGain    float64                      `json:"potential_gain"`    // 潜在收益提升（年化USD）
	RiskLevel        string                       `json:"risk_level"`        // 风险等级（low/medium/high）
	Summary          string                       `json:"summary"`           // 总结说明
}

// LendingRecommendation 借贷优化建议
type LendingRecommendation struct {
	Action      string  `json:"action"`       // 操作类型（migrate/rebalance/reduce_borrow）
	FromProtocol string `json:"from_protocol,omitempty"` // 源协议
	ToProtocol   string  `json:"to_protocol"`   // 目标协议
	Token        string  `json:"token"`         // 代币
	Amount       float64 `json:"amount"`        // 数量
	Reason       string  `json:"reason"`        // 原因说明
	ExpectedGain float64 `json:"expected_gain"` // 预期收益提升（年化USD）
}

// LendingHealthResult 健康因子监控结果
type LendingHealthResult struct {
	HealthyCount    int                   `json:"healthy_count"`     // 健康仓位数量
	AtRiskCount     int                   `json:"at_risk_count"`     // 风险仓位数量
	CriticalCount   int                   `json:"critical_count"`    // 危险仓位数量
	AtRiskPositions []*LendingHealthItem  `json:"at_risk_positions"` // 风险仓位列表
}

// LendingHealthItem 健康因子监控条目
type LendingHealthItem struct {
	Protocol             string  `json:"protocol"`               // 协议名称
	Chain                string  `json:"chain"`                  // 所在链
	WalletAddr           string  `json:"wallet_addr"`            // 钱包地址
	HealthFactor         float64 `json:"health_factor"`          // 健康因子
	LiquidationThreshold float64 `json:"liquidation_threshold"`  // 清算阈值
	SupplyValueUSD       float64 `json:"supply_value_usd"`       // 存款价值（USD）
	BorrowValueUSD       float64 `json:"borrow_value_usd"`       // 借款价值（USD）
	RiskLevel            string  `json:"risk_level"`             // 风险等级（low/medium/high/critical）
}

