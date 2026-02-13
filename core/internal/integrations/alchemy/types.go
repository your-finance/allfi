// Package alchemy Alchemy NFT API 数据类型
package alchemy

// NFT 单个 NFT 资产
type NFT struct {
	ContractAddress string  `json:"contract_address"` // 合约地址
	TokenID         string  `json:"token_id"`         // Token ID
	Name            string  `json:"name"`             // NFT 名称
	Description     string  `json:"description"`      // 描述
	ImageURL        string  `json:"image_url"`        // 图片 URL
	Collection      string  `json:"collection"`       // 收藏集名称
	CollectionSlug  string  `json:"collection_slug"`  // 收藏集标识
	Chain           string  `json:"chain"`            // 所在链
	FloorPrice      float64 `json:"floor_price"`      // Floor Price（ETH 计）
	FloorPriceCurrency string `json:"floor_price_currency"` // Floor Price 计价币种
	FloorPriceUSD   float64 `json:"floor_price_usd"`  // Floor Price（USD 计）
}

// NFTCollection 收藏集汇总
type NFTCollection struct {
	Name           string  `json:"name"`            // 收藏集名称
	Slug           string  `json:"slug"`            // 收藏集标识
	Count          int     `json:"count"`           // 持有数量
	FloorPrice     float64 `json:"floor_price"`     // Floor Price
	FloorCurrency  string  `json:"floor_currency"`  // Floor Price 计价币种
	TotalValueUSD  float64 `json:"total_value_usd"` // 总估值（USD）
	ImageURL       string  `json:"image_url"`       // 收藏集封面
	Chain          string  `json:"chain"`           // 所在链
}

// alchemyNFTResponse Alchemy API 返回的 NFT 列表响应
type alchemyNFTResponse struct {
	OwnedNfts []alchemyNFT `json:"ownedNfts"`
	PageKey   string       `json:"pageKey"`
	TotalCount int         `json:"totalCount"`
}

// alchemyNFT Alchemy API 返回的单个 NFT
type alchemyNFT struct {
	Contract struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"contract"`
	TokenID string `json:"tokenId"`
	Name    string `json:"name"`
	Description string `json:"description"`
	Image struct {
		CachedURL   string `json:"cachedUrl"`
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"image"`
	Collection struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"collection"`
}

// alchemyFloorPriceResponse Alchemy Floor Price 响应
type alchemyFloorPriceResponse struct {
	OpenSea struct {
		FloorPrice    float64 `json:"floorPrice"`
		PriceCurrency string  `json:"priceCurrency"`
	} `json:"openSea"`
}
