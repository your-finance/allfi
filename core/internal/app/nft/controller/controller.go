// Package controller NFT 资产模块 - 控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	nftApi "your-finance/allfi/api/v1/nft"
	"your-finance/allfi/internal/app/nft/service"
)

// NftController NFT 资产控制器
type NftController struct{}

// GetAssets 获取 NFT 资产列表
func (c *NftController) GetAssets(ctx context.Context, req *nftApi.GetAssetsReq) (res *nftApi.GetAssetsRes, err error) {
	// 调用服务层获取 NFT 资产
	nfts, totalValue, err := service.Nft().GetNFTs(ctx, "", "")
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &nftApi.GetAssetsRes{
		TotalValue: totalValue,
		Assets:     make([]nftApi.NFTItem, 0, len(nfts)),
	}

	for _, n := range nfts {
		res.Assets = append(res.Assets, nftApi.NFTItem{
			ID:              n.ID,
			ContractAddress: n.ContractAddress,
			TokenID:         n.TokenID,
			Name:            n.Name,
			Description:     n.Description,
			ImageURL:        n.ImageURL,
			Collection:      n.Collection,
			Chain:           n.Chain,
			FloorPrice:      n.FloorPrice,
			EstimatedValue:  n.EstimatedValue,
			WalletAddr:      n.WalletAddr,
		})
	}

	return res, nil
}

// GetCollections 获取 NFT 收藏集统计
func (c *NftController) GetCollections(ctx context.Context, req *nftApi.GetCollectionsReq) (res *nftApi.GetCollectionsRes, err error) {
	return service.Nft().GetCollections(ctx)
}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&NftController{})
}
