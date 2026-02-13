// Package service NFT 资产模块 - Service 层接口定义
package service

import (
	"context"

	nftApi "your-finance/allfi/api/v1/nft"
	"your-finance/allfi/internal/app/nft/model"
)

// INft NFT 资产服务接口
type INft interface {
	// GetNFTs 获取用户 NFT 资产列表
	// chain: 指定链名（ethereum 等），空字符串默认为 ethereum
	// collection: 按收藏集过滤，空字符串不过滤
	// 返回 NFT 列表和总估值
	GetNFTs(ctx context.Context, chain string, collection string) (nfts []*model.NFTItem, totalValue float64, err error)

	// GetCollections 获取 NFT 收藏集统计（按 collection+chain 分组聚合）
	GetCollections(ctx context.Context) (*nftApi.GetCollectionsRes, error)
}

var localNft INft

// Nft 获取 NFT 资产服务实例
func Nft() INft {
	if localNft == nil {
		panic("INft 服务未注册，请检查 logic/nft 包的 init 函数")
	}
	return localNft
}

// RegisterNft 注册 NFT 资产服务实现
// 由 logic 层在 init 函数中调用
func RegisterNft(i INft) {
	localNft = i
}
