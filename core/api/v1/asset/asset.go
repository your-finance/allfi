// Package asset 资产 API 定义
// 提供资产概览、明细、历史趋势、强制刷新接口
package asset

import "github.com/gogf/gf/v2/frame/g"

// GetSummaryReq 获取资产概览请求
type GetSummaryReq struct {
	g.Meta   `path:"/assets/summary" method:"get" summary:"获取资产概览" tags:"资产"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币（USD/BTC/ETH/CNY）"`
}

// GetSummaryRes 获取资产概览响应
type GetSummaryRes struct {
	TotalValue float64            `json:"total_value" dc:"总资产价值"`
	Currency   string             `json:"currency" dc:"计价货币"`
	BySource   map[string]float64 `json:"by_source" dc:"按来源分类（cex/blockchain/manual）"`
	UpdatedAt  string             `json:"updated_at" dc:"最后更新时间"`
}

// GetDetailsReq 获取资产明细请求
type GetDetailsReq struct {
	g.Meta     `path:"/assets/details" method:"get" summary:"获取资产明细列表" tags:"资产"`
	SourceType string `json:"source_type" in:"query" dc:"来源类型筛选（cex/blockchain/manual）"`
	Currency   string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// AssetDetailItem 资产明细条目
type AssetDetailItem struct {
	ID         uint    `json:"id" dc:"资产 ID"`
	Symbol     string  `json:"symbol" dc:"币种符号"`
	Amount     float64 `json:"amount" dc:"持有数量"`
	Value      float64 `json:"value" dc:"价值（计价货币）"`
	Price      float64 `json:"price" dc:"单价"`
	Source     string  `json:"source" dc:"来源（交易所名/链名/手动）"`
	SourceType string  `json:"source_type" dc:"来源类型"`
	UpdatedAt  string  `json:"updated_at" dc:"更新时间"`
}

// GetDetailsRes 获取资产明细响应
type GetDetailsRes struct {
	Assets []AssetDetailItem `json:"assets" dc:"资产明细列表"`
}

// GetHistoryReq 获取资产历史趋势请求
type GetHistoryReq struct {
	g.Meta   `path:"/assets/history" method:"get" summary:"获取资产历史趋势" tags:"资产"`
	Days     int    `json:"days" in:"query" d:"30" dc:"查询天数"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// SnapshotItem 资产快照条目
type SnapshotItem struct {
	Date       string  `json:"date" dc:"日期"`
	TotalValue float64 `json:"total_value" dc:"总资产价值"`
	Currency   string  `json:"currency" dc:"计价货币"`
}

// GetHistoryRes 获取资产历史响应
type GetHistoryRes struct {
	Snapshots []SnapshotItem `json:"snapshots" dc:"资产快照列表"`
}

// RefreshReq 强制刷新所有资产请求
type RefreshReq struct {
	g.Meta `path:"/assets/refresh" method:"post" summary:"强制刷新所有资产" tags:"资产"`
}

// RefreshRes 强制刷新所有资产响应
type RefreshRes struct {
	Message        string `json:"message" dc:"刷新结果消息"`
	RefreshedCount int    `json:"refreshed_count" dc:"刷新的资产数量"`
}
