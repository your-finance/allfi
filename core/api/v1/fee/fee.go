// Package fee 费用分析 API 定义
// 提供交易费用（手续费 + Gas）分析接口
package fee

import "github.com/gogf/gf/v2/frame/g"

// GetAnalyticsReq 获取费用分析请求
type GetAnalyticsReq struct {
	g.Meta   `path:"/analytics/fees" method:"get" summary:"获取费用分析" tags:"费用分析"`
	Range    string `json:"range" in:"query" d:"30D" dc:"时间范围（7D/30D/90D）"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// FeeBreakdown 费用构成（对象形式，用于前端环形图）
type FeeBreakdownObj struct {
	CexTradeFee float64 `json:"cexTradeFee" dc:"CEX 交易手续费"`
	GasFee      float64 `json:"gasFee" dc:"链上 Gas 费"`
	WithdrawFee float64 `json:"withdrawFee" dc:"提现手续费"`
}

// MonthlyTrend 月度费用趋势
type MonthlyTrend struct {
	Month       string  `json:"month" dc:"月份（YYYY-MM）"`
	CexTradeFee float64 `json:"cexTradeFee" dc:"CEX 交易手续费"`
	GasFee      float64 `json:"gasFee" dc:"链上 Gas 费"`
	WithdrawFee float64 `json:"withdrawFee" dc:"提现手续费"`
	Total       float64 `json:"total" dc:"当月总费用"`
}

// Suggestion 费用优化建议
type Suggestion struct {
	ID             string  `json:"id" dc:"建议 ID"`
	Type           string  `json:"type" dc:"建议类型（gas/timing/consolidate/exchange）"`
	TitleKey       string  `json:"titleKey" dc:"标题 i18n key"`
	DescKey        string  `json:"descKey" dc:"描述 i18n key"`
	SavingEstimate float64 `json:"savingEstimate" dc:"预估节省金额"`
	Priority       string  `json:"priority" dc:"优先级（high/medium/low）"`
	ImpactScore    int     `json:"impactScore" dc:"影响评分（0-100）"`
}

// GetAnalyticsRes 获取费用分析响应
type GetAnalyticsRes struct {
	Total         float64        `json:"total" dc:"总费用"`
	ChangePercent float64        `json:"changePercent" dc:"较上月变化百分比"`
	Breakdown     FeeBreakdownObj `json:"breakdown" dc:"费用构成"`
	MonthlyTrend  []MonthlyTrend `json:"monthlyTrend" dc:"月度费用趋势"`
	Suggestions   []Suggestion   `json:"suggestions" dc:"费用优化建议"`
}
