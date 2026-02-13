// Package nft NFT 资产 API 定义
// 提供 NFT 资产列表查询接口
package nft

import "github.com/gogf/gf/v2/frame/g"

// GetAssetsReq 获取 NFT 资产列表请求
type GetAssetsReq struct {
	g.Meta `path:"/nft/assets" method:"get" summary:"获取 NFT 资产列表" tags:"NFT"`
}

// NFTItem NFT 资产条目
type NFTItem struct {
	ID              uint    `json:"id" dc:"NFT ID"`
	ContractAddress string  `json:"contract_address" dc:"合约地址"`
	TokenID         string  `json:"token_id" dc:"Token ID"`
	Name            string  `json:"name" dc:"名称"`
	Description     string  `json:"description" dc:"描述"`
	ImageURL        string  `json:"image_url" dc:"图片 URL"`
	Collection      string  `json:"collection" dc:"收藏集名称"`
	Chain           string  `json:"chain" dc:"所在链"`
	FloorPrice      float64 `json:"floor_price" dc:"地板价（ETH）"`
	EstimatedValue  float64 `json:"estimated_value" dc:"估值（USD）"`
	WalletAddr      string  `json:"wallet_addr" dc:"持有钱包地址"`
}

// GetAssetsRes 获取 NFT 资产列表响应
type GetAssetsRes struct {
	Assets     []NFTItem `json:"assets" dc:"NFT 资产列表"`
	TotalValue float64   `json:"total_value" dc:"总估值（USD）"`
}

// GetCollectionsReq NFT 收藏集统计请求
type GetCollectionsReq struct {
	g.Meta `path:"/nfts/collections" method:"get" summary:"获取 NFT 收藏集统计" tags:"NFT"`
}

// CollectionItem 收藏集统计项
type CollectionItem struct {
	Collection      string  `json:"collection" dc:"收藏集名称"`
	Count           int     `json:"count" dc:"NFT 数量"`
	TotalFloorPrice float64 `json:"total_floor_price" dc:"地板价总值（USD）"`
	Chain           string  `json:"chain" dc:"所在链"`
}

// GetCollectionsRes NFT 收藏集统计响应
type GetCollectionsRes struct {
	Collections []CollectionItem `json:"collections" dc:"收藏集列表"`
	TotalCount  int              `json:"total_count" dc:"NFT 总数"`
	TotalValue  float64          `json:"total_value" dc:"总价值（USD）"`
}
