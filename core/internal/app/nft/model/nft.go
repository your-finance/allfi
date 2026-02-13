// Package model NFT 资产模块 - 数据传输对象
package model

// NFTItem NFT 资产条目
type NFTItem struct {
	ID              uint    `json:"id"`               // 记录 ID
	ContractAddress string  `json:"contract_address"` // 合约地址
	TokenID         string  `json:"token_id"`         // Token ID
	Name            string  `json:"name"`             // 名称
	Description     string  `json:"description"`      // 描述
	ImageURL        string  `json:"image_url"`        // 图片 URL
	Collection      string  `json:"collection"`       // 收藏集名称
	Chain           string  `json:"chain"`            // 所在链
	FloorPrice      float64 `json:"floor_price"`      // 地板价（ETH）
	EstimatedValue  float64 `json:"estimated_value"`  // 估值（USD）
	WalletAddr      string  `json:"wallet_addr"`      // 持有钱包地址
}
