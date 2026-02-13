// Package v1 手动资产 API 定义
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/model/entity"
)

// AddManualAssetReq 添加手动资产请求
type AddManualAssetReq struct {
	g.Meta      `path:"/manual-assets" method:"post" summary:"添加手动资产" tags:"手动资产"`
	AssetType   string  `json:"asset_type" v:"required|in:cash,bank,stock,fund" dc:"资产类型"`
	AssetName   string  `json:"asset_name" v:"required|max-length:100" dc:"资产名称"`
	Institution string  `json:"institution" v:"max-length:100" dc:"机构名称（银行/券商/基金平台）"`
	Amount      float64 `json:"amount" v:"required|min:0" dc:"数量"`
	UnitPrice   float64 `json:"unit_price" v:"min:0" dc:"单价（可选）"`
	Currency    string  `json:"currency" v:"required|in:CNY,USD,EUR,HKD,JPY,GBP,SGD" dc:"货币"`
	Description string  `json:"description" v:"max-length:500" dc:"备注说明"`
}

// AddManualAssetRes 添加手动资产响应
type AddManualAssetRes struct {
	Asset *entity.ManualAsset `json:"asset" dc:"手动资产信息"`
}

// UpdateManualAssetReq 更新手动资产请求
type UpdateManualAssetReq struct {
	g.Meta      `path:"/manual-assets/:id" method:"put" summary:"更新手动资产" tags:"手动资产"`
	AssetId     uint    `json:"id" v:"required|min:1" dc:"资产ID"`
	Amount      float64 `json:"amount" v:"required|min:0" dc:"数量"`
	UnitPrice   float64 `json:"unit_price" v:"min:0" dc:"单价"`
	Description string  `json:"description" v:"max-length:500" dc:"备注说明"`
}

// UpdateManualAssetRes 更新手动资产响应
type UpdateManualAssetRes struct {
	Asset *entity.ManualAsset `json:"asset" dc:"手动资产信息"`
}

// ListManualAssetsReq 获取手动资产列表请求
type ListManualAssetsReq struct {
	g.Meta    `path:"/manual-assets" method:"get" summary:"获取手动资产列表" tags:"手动资产"`
	AssetType string `json:"asset_type" v:"in:cash,bank,stock,fund" dc:"资产类型（可选）"`
}

// ListManualAssetsRes 获取手动资产列表响应
type ListManualAssetsRes struct {
	Assets []*entity.ManualAsset `json:"assets" dc:"手动资产列表"`
}

// DeleteManualAssetReq 删除手动资产请求
type DeleteManualAssetReq struct {
	g.Meta  `path:"/manual-assets/:id" method:"delete" summary:"删除手动资产" tags:"手动资产"`
	AssetId uint `json:"id" v:"required|min:1" dc:"资产ID"`
}

// DeleteManualAssetRes 删除手动资产响应
type DeleteManualAssetRes struct{}

// GetManualAssetStatsReq 获取手动资产统计请求
type GetManualAssetStatsReq struct {
	g.Meta `path:"/manual-assets/stats" method:"get" summary:"获取手动资产统计" tags:"手动资产"`
}

// GetManualAssetStatsRes 获取手动资产统计响应
type GetManualAssetStatsRes struct {
	TotalValue float64            `json:"total_value" dc:"总价值（CNY）"`
	ByType     map[string]float64 `json:"by_type" dc:"按类型分类"`
	ByCurrency map[string]float64 `json:"by_currency" dc:"按货币分类"`
	AssetCount int                `json:"asset_count" dc:"资产数量"`
}
