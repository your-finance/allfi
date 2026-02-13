// Package v1 资产 API 定义
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/model/entity"
)

// GetAssetSummaryReq 获取资产概览请求
type GetAssetSummaryReq struct {
	g.Meta   `path:"/assets/summary" method:"get" summary:"获取资产概览" tags:"资产"`
	Currency string `json:"currency" v:"in:USDC,BTC,ETH,CNY" dc:"计价货币（默认USDC）"`
}

// GetAssetSummaryRes 获取资产概览响应
type GetAssetSummaryRes struct {
	Summary *AssetSummary `json:"summary" dc:"资产概览"`
}

// AssetSummary 资产概览
type AssetSummary struct {
	TotalValue     float64             `json:"total_value" dc:"总资产价值"`
	Currency       string              `json:"currency" dc:"计价货币"`
	BySource       map[string]float64  `json:"by_source" dc:"按来源分类（CEX/链上/手动）"`
	UpdatedAt      string              `json:"updated_at" dc:"更新时间"`
}

// GetAssetDetailsReq 获取资产明细请求
type GetAssetDetailsReq struct {
	g.Meta     `path:"/assets/details" method:"get" summary:"获取资产明细列表" tags:"资产"`
	SourceType string `json:"source_type" v:"in:cex,blockchain,manual" dc:"来源类型（可选）"`
	Currency   string `json:"currency" v:"in:USDC,BTC,ETH,CNY" dc:"计价货币（默认USDC）"`
}

// GetAssetDetailsRes 获取资产明细响应
type GetAssetDetailsRes struct {
	Assets []*entity.AssetDetail `json:"assets" dc:"资产明细列表"`
}

// RefreshAssetsReq 刷新资产请求
type RefreshAssetsReq struct {
	g.Meta `path:"/assets/refresh" method:"post" summary:"强制刷新所有资产" tags:"资产"`
}

// RefreshAssetsRes 刷新资产响应
type RefreshAssetsRes struct {
	Message   string `json:"message" dc:"刷新结果消息"`
	RefreshedCount int `json:"refreshed_count" dc:"刷新的资产数量"`
}

// GetAssetHistoryReq 获取资产历史请求
type GetAssetHistoryReq struct {
	g.Meta    `path:"/assets/history" method:"get" summary:"获取资产历史趋势" tags:"资产"`
	TimeRange string `json:"time_range" v:"in:7d,30d,90d,1y" dc:"时间范围（默认30天）"`
	Currency  string `json:"currency" v:"in:USDC,BTC,ETH,CNY" dc:"计价货币（默认USDC）"`
}

// GetAssetHistoryRes 获取资产历史响应
type GetAssetHistoryRes struct {
	Snapshots []*entity.AssetSnapshot `json:"snapshots" dc:"资产快照列表"`
}
