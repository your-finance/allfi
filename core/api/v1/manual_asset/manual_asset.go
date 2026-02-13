// Package manual_asset 手动资产 API 定义
// 提供手动资产（银行账户/现金/股票/基金）的增删改查接口
package manual_asset

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取手动资产列表请求
type ListReq struct {
	g.Meta `path:"/manual/assets" method:"get" summary:"获取手动资产列表" tags:"手动资产"`
}

// ManualAssetItem 手动资产条目
type ManualAssetItem struct {
	ID          uint    `json:"id" dc:"资产 ID"`
	AssetType   string  `json:"asset_type" dc:"资产类型（cash/bank/stock/fund）"`
	AssetName   string  `json:"asset_name" dc:"资产名称"`
	Amount      float64 `json:"amount" dc:"数量"`
	Currency    string  `json:"currency" dc:"货币"`
	Notes       string  `json:"notes" dc:"备注"`
	Institution string  `json:"institution" dc:"机构名称"`
	CreatedAt   string  `json:"created_at" dc:"创建时间"`
	UpdatedAt   string  `json:"updated_at" dc:"更新时间"`
}

// ListRes 获取手动资产列表响应
type ListRes struct {
	Assets []ManualAssetItem `json:"assets" dc:"手动资产列表"`
}

// CreateReq 添加手动资产请求
type CreateReq struct {
	g.Meta      `path:"/manual/assets" method:"post" summary:"添加手动资产" tags:"手动资产"`
	AssetType   string  `json:"asset_type" v:"required|in:cash,bank,stock,fund" dc:"资产类型"`
	AssetName   string  `json:"asset_name" v:"required|max-length:100" dc:"资产名称"`
	Amount      float64 `json:"amount" v:"required|min:0" dc:"数量"`
	Currency    string  `json:"currency" v:"required" dc:"货币（CNY/USD/EUR/HKD 等）"`
	Notes       string  `json:"notes" dc:"备注说明"`
	Institution string  `json:"institution" dc:"机构名称（银行/券商/基金平台）"`
}

// CreateRes 添加手动资产响应
type CreateRes struct {
	Asset *ManualAssetItem `json:"asset" dc:"新建的手动资产信息"`
}

// UpdateReq 更新手动资产请求
type UpdateReq struct {
	g.Meta      `path:"/manual/assets/{id}" method:"put" summary:"更新手动资产" tags:"手动资产"`
	Id          uint    `json:"id" in:"path" v:"required|min:1" dc:"资产 ID"`
	AssetType   string  `json:"asset_type" dc:"资产类型"`
	AssetName   string  `json:"asset_name" dc:"资产名称"`
	Amount      float64 `json:"amount" dc:"数量"`
	Currency    string  `json:"currency" dc:"货币"`
	Notes       string  `json:"notes" dc:"备注说明"`
	Institution string  `json:"institution" dc:"机构名称"`
}

// UpdateRes 更新手动资产响应
type UpdateRes struct {
	Asset *ManualAssetItem `json:"asset" dc:"更新后的手动资产信息"`
}

// DeleteReq 删除手动资产请求
type DeleteReq struct {
	g.Meta `path:"/manual/assets/{id}" method:"delete" summary:"删除手动资产" tags:"手动资产"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"资产 ID"`
}

// DeleteRes 删除手动资产响应
type DeleteRes struct{}
