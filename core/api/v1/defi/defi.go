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

// GetLendingPositionsReq 获取借贷仓位请求
type GetLendingPositionsReq struct {
	g.Meta `path:"/defi/lending/positions" method:"get" summary:"获取借贷仓位列表" tags:"DeFi"`
}

// LendingPositionItem 借贷仓位条目
type LendingPositionItem struct {
	Protocol             string  `json:"protocol" dc:"协议名称"`
	Chain                string  `json:"chain" dc:"所在链"`
	WalletAddr           string  `json:"wallet_addr" dc:"钱包地址"`
	SupplyToken          string  `json:"supply_token" dc:"存款代币"`
	SupplyAmount         float64 `json:"supply_amount" dc:"存款数量"`
	SupplyValueUSD       float64 `json:"supply_value_usd" dc:"存款价值（USD）"`
	SupplyAPY            float64 `json:"supply_apy" dc:"存款年化收益率"`
	BorrowToken          string  `json:"borrow_token,omitempty" dc:"借款代币"`
	BorrowAmount         float64 `json:"borrow_amount" dc:"借款数量"`
	BorrowValueUSD       float64 `json:"borrow_value_usd" dc:"借款价值（USD）"`
	BorrowAPY            float64 `json:"borrow_apy" dc:"借款年化利率"`
	HealthFactor         float64 `json:"health_factor" dc:"健康因子"`
	LiquidationThreshold float64 `json:"liquidation_threshold" dc:"清算阈值"`
	LTV                  float64 `json:"ltv" dc:"Loan-to-Value 比率"`
	NetAPY               float64 `json:"net_apy" dc:"净收益率"`
}

// GetLendingPositionsRes 获取借贷仓位响应
type GetLendingPositionsRes struct {
	Positions []LendingPositionItem `json:"positions" dc:"借贷仓位列表"`
}

// GetLendingRatesReq 获取借贷利率请求
type GetLendingRatesReq struct {
	g.Meta   `path:"/defi/lending/rates" method:"get" summary:"获取借贷利率" tags:"DeFi"`
	Protocol string `json:"protocol" in:"query" dc:"协议名称（可选）"`
	Chain    string `json:"chain" in:"query" dc:"链名称（可选）"`
}

// LendingRateItem 借贷利率条目
type LendingRateItem struct {
	Protocol          string  `json:"protocol" dc:"协议名称"`
	Chain             string  `json:"chain" dc:"所在链"`
	Token             string  `json:"token" dc:"代币符号"`
	SupplyAPY         float64 `json:"supply_apy" dc:"存款年化收益率"`
	BorrowAPYStable   float64 `json:"borrow_apy_stable" dc:"稳定借款利率"`
	BorrowAPYVariable float64 `json:"borrow_apy_variable" dc:"浮动借款利率"`
	UtilizationRate   float64 `json:"utilization_rate" dc:"利用率"`
}

// GetLendingRatesRes 获取借贷利率响应
type GetLendingRatesRes struct {
	Rates []LendingRateItem `json:"rates" dc:"借贷利率列表"`
}

// GetLendingRateHistoryReq 获取借贷利率历史请求
type GetLendingRateHistoryReq struct {
	g.Meta   `path:"/defi/lending/rate-history" method:"get" summary:"获取借贷利率历史" tags:"DeFi"`
	Protocol string `json:"protocol" in:"query" v:"required" dc:"协议名称"`
	Token    string `json:"token" in:"query" v:"required" dc:"代币符号"`
	Days     int    `json:"days" in:"query" d:"30" dc:"历史天数"`
}

// LendingRateHistoryItem 借贷利率历史条目
type LendingRateHistoryItem struct {
	Date              string  `json:"date" dc:"日期"`
	SupplyAPY         float64 `json:"supply_apy" dc:"存款年化收益率"`
	BorrowAPYStable   float64 `json:"borrow_apy_stable" dc:"稳定借款利率"`
	BorrowAPYVariable float64 `json:"borrow_apy_variable" dc:"浮动借款利率"`
	UtilizationRate   float64 `json:"utilization_rate" dc:"利用率"`
}

// GetLendingRateHistoryRes 获取借贷利率历史响应
type GetLendingRateHistoryRes struct {
	History []LendingRateHistoryItem `json:"history" dc:"利率历史"`
}

// GetLendingOptimizationReq 获取借贷策略优化请求
type GetLendingOptimizationReq struct {
	g.Meta `path:"/defi/lending/optimization" method:"get" summary:"获取借贷策略优化建议" tags:"DeFi"`
}

// LendingRecommendation 借贷优化建议
type LendingRecommendation struct {
	Action       string  `json:"action" dc:"操作类型"`
	FromProtocol string  `json:"from_protocol,omitempty" dc:"源协议"`
	ToProtocol   string  `json:"to_protocol" dc:"目标协议"`
	Token        string  `json:"token" dc:"代币"`
	Amount       float64 `json:"amount" dc:"数量"`
	Reason       string  `json:"reason" dc:"原因说明"`
	ExpectedGain float64 `json:"expected_gain" dc:"预期收益提升"`
}

// GetLendingOptimizationRes 获取借贷策略优化响应
type GetLendingOptimizationRes struct {
	CurrentPositions []LendingPositionItem   `json:"current_positions" dc:"当前仓位"`
	Recommendations  []LendingRecommendation `json:"recommendations" dc:"优化建议"`
	PotentialGain    float64                 `json:"potential_gain" dc:"潜在收益提升"`
	RiskLevel        string                  `json:"risk_level" dc:"风险等级"`
	Summary          string                  `json:"summary" dc:"总结说明"`
}

