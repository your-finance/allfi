// Package logic NFT 资产模块 - 业务逻辑实现
package logic

import (
	"context"
	"math"
	"os"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	nftApi "your-finance/allfi/api/v1/nft"
	"your-finance/allfi/internal/app/nft/dao"
	"your-finance/allfi/internal/app/nft/model"
	"your-finance/allfi/internal/app/nft/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/alchemy"
	"your-finance/allfi/internal/model/entity"

	walletDao "your-finance/allfi/internal/app/wallet/dao"
)

// sNft NFT 资产服务实现
type sNft struct{}

// New 创建 NFT 资产服务实例
func New() service.INft {
	return &sNft{}
}

// cacheTTL NFT 缓存有效期（1 小时）
const cacheTTL = 1 * time.Hour

// GetNFTs 获取用户 NFT 资产列表
//
// 业务逻辑:
// 1. 获取用户所有钱包地址
// 2. 对每个钱包:
//   - 检查缓存: 从 nft_caches 表查询该钱包地址的缓存
//   - 缓存有效则使用缓存数据
//   - 否则从 Alchemy API 获取（当前使用缓存兜底）
//
// 3. 按收藏集过滤（如指定）
// 4. 计算 Floor Price 总值
func (s *sNft) GetNFTs(ctx context.Context, chain string, collection string) (nfts []*model.NFTItem, totalValue float64, err error) {
	userID := consts.GetUserID(ctx)

	// 获取用户所有钱包地址
	var wallets []*entity.WalletAddresses
	err = walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, userID).
		Scan(&wallets)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询钱包地址失败")
	}

	if len(wallets) == 0 {
		g.Log().Info(ctx, "用户无钱包地址，跳过 NFT 查询")
		return []*model.NFTItem{}, 0, nil
	}

	nfts = make([]*model.NFTItem, 0)

	for _, wallet := range wallets {
		// 从缓存获取 NFT 数据
		var caches []*entity.NftCaches
		err = dao.NftCaches.Ctx(ctx).
			Where(dao.NftCaches.Columns().WalletAddress, wallet.Address).
			Scan(&caches)
		if err != nil {
			g.Log().Warning(ctx, "查询 NFT 缓存失败",
				"wallet", wallet.Address,
				"error", err,
			)
			continue
		}

		// 检查缓存是否有效
		if len(caches) > 0 && isCacheValid(caches[0].CachedAt) {
			// 使用缓存数据
			for _, cache := range caches {
				item := cacheToNFTItem(cache, wallet.Address)

				// 按收藏集过滤
				if collection != "" && cache.Collection != collection && cache.CollectionSlug != collection {
					continue
				}

				nfts = append(nfts, item)
				totalValue += float64(cache.FloorPriceUsd)
			}

			g.Log().Debug(ctx, "使用 NFT 缓存数据",
				"wallet", wallet.Address,
				"count", len(caches),
			)
		} else {
			// 缓存过期或不存在，尝试通过 Alchemy API 刷新
			g.Log().Info(ctx, "NFT 缓存过期，尝试从 Alchemy API 刷新",
				"wallet", wallet.Address,
				"cacheCount", len(caches),
			)

			refreshed := false
			apiKey := os.Getenv("ALCHEMY_API_KEY")
			if apiKey != "" {
				// 创建 Alchemy 客户端并获取最新 NFT 数据
				client := alchemy.NewClient(apiKey, nil)
				// 依次查询支持的链
				for _, chainName := range alchemy.SupportedChains() {
					alchemyNFTs, err := client.GetNFTs(ctx, wallet.Address, chainName)
					if err != nil {
						g.Log().Warning(ctx, "Alchemy API 获取 NFT 失败",
							"wallet", wallet.Address,
							"chain", chainName,
							"error", err,
						)
						continue
					}

					if len(alchemyNFTs) == 0 {
						continue
					}

					// 补充 Floor Price
					enrichFloorPrices(ctx, client, alchemyNFTs, chainName)

					// 将 Alchemy 数据写入 nft_caches 表（更新缓存）
					updateNFTCache(ctx, wallet.Address, alchemyNFTs)
					refreshed = true

					// 转换为 NFTItem 并追加到结果
					for _, n := range alchemyNFTs {
						if collection != "" && n.Collection != collection && n.CollectionSlug != collection {
							continue
						}
						item := &model.NFTItem{
							ContractAddress: n.ContractAddress,
							TokenID:         n.TokenID,
							Name:            n.Name,
							Description:     n.Description,
							ImageURL:        n.ImageURL,
							Collection:      n.Collection,
							Chain:           n.Chain,
							FloorPrice:      n.FloorPrice,
							EstimatedValue:  n.FloorPriceUSD,
							WalletAddr:      wallet.Address,
						}
						nfts = append(nfts, item)
						totalValue += n.FloorPriceUSD
					}
				}
			} else {
				g.Log().Warning(ctx, "ALCHEMY_API_KEY 未配置，无法刷新 NFT 缓存")
			}

			// Alchemy 调用失败或未配置时，使用过期缓存兜底
			if !refreshed && len(caches) > 0 {
				g.Log().Info(ctx, "Alchemy 刷新失败，使用过期缓存兜底",
					"wallet", wallet.Address,
					"cacheCount", len(caches),
				)
				for _, cache := range caches {
					item := cacheToNFTItem(cache, wallet.Address)

					if collection != "" && cache.Collection != collection && cache.CollectionSlug != collection {
						continue
					}

					nfts = append(nfts, item)
					totalValue += float64(cache.FloorPriceUsd)
				}
			}
		}
	}

	totalValue = math.Round(totalValue*100) / 100

	g.Log().Info(ctx, "获取 NFT 资产成功",
		"walletCount", len(wallets),
		"nftCount", len(nfts),
		"totalValue", totalValue,
	)

	return nfts, totalValue, nil
}

// GetCollections 获取 NFT 收藏集统计（按 collection+chain 分组聚合）
//
// 业务逻辑:
// 1. 查询当前用户所有钱包地址
// 2. 从 nft_caches 表按 collection+chain 分组统计
// 3. 计算总数和总价值
func (s *sNft) GetCollections(ctx context.Context) (*nftApi.GetCollectionsRes, error) {
	userID := consts.GetUserID(ctx)

	// 获取用户所有钱包地址
	var wallets []*entity.WalletAddresses
	err := walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, userID).
		Scan(&wallets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询钱包地址失败")
	}

	if len(wallets) == 0 {
		g.Log().Info(ctx, "用户无钱包地址，跳过 NFT 收藏集统计")
		return &nftApi.GetCollectionsRes{
			Collections: []nftApi.CollectionItem{},
			TotalCount:  0,
			TotalValue:  0,
		}, nil
	}

	// 收集钱包地址列表
	addrs := make([]string, 0, len(wallets))
	for _, w := range wallets {
		addrs = append(addrs, w.Address)
	}

	// 按 collection + chain 分组查询
	type collectionRow struct {
		Collection     string  `json:"collection"`
		Chain          string  `json:"chain"`
		Count          int     `json:"count"`
		TotalFloorPrice float64 `json:"total_floor_price"`
	}

	var rows []collectionRow
	err = dao.NftCaches.Ctx(ctx).
		Fields(
			dao.NftCaches.Columns().Collection,
			dao.NftCaches.Columns().Chain,
			"COUNT(*) as count",
			"SUM(floor_price_usd) as total_floor_price",
		).
		Where(dao.NftCaches.Columns().WalletAddress, addrs).
		Group(dao.NftCaches.Columns().Collection, dao.NftCaches.Columns().Chain).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询 NFT 收藏集统计失败")
	}

	// 构建响应
	collections := make([]nftApi.CollectionItem, 0, len(rows))
	totalCount := 0
	totalValue := 0.0

	for _, row := range rows {
		collections = append(collections, nftApi.CollectionItem{
			Collection:      row.Collection,
			Count:           row.Count,
			TotalFloorPrice: math.Round(row.TotalFloorPrice*100) / 100,
			Chain:           row.Chain,
		})
		totalCount += row.Count
		totalValue += row.TotalFloorPrice
	}

	totalValue = math.Round(totalValue*100) / 100

	g.Log().Info(ctx, "获取 NFT 收藏集统计成功",
		"collectionCount", len(collections),
		"totalCount", totalCount,
		"totalValue", totalValue,
	)

	return &nftApi.GetCollectionsRes{
		Collections: collections,
		TotalCount:  totalCount,
		TotalValue:  totalValue,
	}, nil
}

// isCacheValid 检查缓存是否仍然有效
func isCacheValid(cachedAt time.Time) bool {
	return time.Since(cachedAt) < cacheTTL
}

// cacheToNFTItem 将缓存记录转换为 NFT 条目
func cacheToNFTItem(cache *entity.NftCaches, walletAddr string) *model.NFTItem {
	return &model.NFTItem{
		ID:              uint(cache.Id),
		ContractAddress: cache.ContractAddress,
		TokenID:         cache.TokenId,
		Name:            cache.Name,
		Description:     cache.Description,
		ImageURL:        cache.ImageUrl,
		Collection:      cache.Collection,
		Chain:           cache.Chain,
		FloorPrice:      float64(cache.FloorPrice),
		EstimatedValue:  float64(cache.FloorPriceUsd),
		WalletAddr:      walletAddr,
	}
}

// enrichFloorPrices 为 Alchemy 返回的 NFT 列表补充 Floor Price（USD 估值）
// 按合约地址去重查询，避免重复 API 调用
func enrichFloorPrices(ctx context.Context, client *alchemy.Client, nfts []alchemy.NFT, chain string) {
	// 收集唯一合约地址
	contracts := make(map[string]bool)
	for _, nft := range nfts {
		contracts[nft.ContractAddress] = true
	}

	// 查询每个合约的 Floor Price
	type floorInfo struct {
		price    float64
		priceUSD float64
		currency string
	}
	floorPrices := make(map[string]floorInfo)

	for addr := range contracts {
		price, priceUSD, err := client.GetFloorPriceUSD(ctx, addr, chain)
		if err != nil {
			g.Log().Debug(ctx, "获取 Floor Price 失败",
				"contract", addr,
				"chain", chain,
				"error", err,
			)
			continue
		}
		floorPrices[addr] = floorInfo{price: price, priceUSD: priceUSD, currency: "ETH"}
	}

	// 回写到 NFT 列表
	for i := range nfts {
		if fp, ok := floorPrices[nfts[i].ContractAddress]; ok {
			nfts[i].FloorPrice = fp.price
			nfts[i].FloorPriceUSD = fp.priceUSD
			nfts[i].FloorPriceCurrency = fp.currency
		}
	}
}

// updateNFTCache 将 Alchemy 获取的 NFT 数据写入 nft_caches 表
// 先删除该钱包地址的旧缓存，再批量插入新数据
func updateNFTCache(ctx context.Context, walletAddress string, nfts []alchemy.NFT) {
	// 删除该钱包的旧缓存
	_, err := dao.NftCaches.Ctx(ctx).
		Where(dao.NftCaches.Columns().WalletAddress, walletAddress).
		Delete()
	if err != nil {
		g.Log().Warning(ctx, "删除旧 NFT 缓存失败",
			"wallet", walletAddress,
			"error", err,
		)
		return
	}

	// 批量插入新缓存
	now := time.Now()
	for _, nft := range nfts {
		cache := entity.NftCaches{
			UserId:          consts.GetUserID(ctx),
			WalletAddress:   walletAddress,
			ContractAddress: nft.ContractAddress,
			TokenId:         nft.TokenID,
			Name:            nft.Name,
			Description:     nft.Description,
			ImageUrl:        nft.ImageURL,
			Collection:      nft.Collection,
			CollectionSlug:  nft.CollectionSlug,
			Chain:           nft.Chain,
			FloorPrice:      float32(nft.FloorPrice),
			FloorCurrency:   nft.FloorPriceCurrency,
			FloorPriceUsd:   float32(nft.FloorPriceUSD),
			CachedAt:        now,
		}
		_, insertErr := dao.NftCaches.Ctx(ctx).Data(cache).Insert()
		if insertErr != nil {
			g.Log().Warning(ctx, "插入 NFT 缓存失败",
				"wallet", walletAddress,
				"nft", nft.Name,
				"error", insertErr,
			)
		}
	}

	g.Log().Info(ctx, "NFT 缓存已更新",
		"wallet", walletAddress,
		"count", len(nfts),
	)
}
