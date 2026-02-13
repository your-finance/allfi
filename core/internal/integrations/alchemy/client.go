// Package alchemy Alchemy NFT API 客户端
// 通过 Alchemy API 获取用户 NFT 资产和 Floor Price 估值
package alchemy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"your-finance/allfi/internal/integrations"
)

// Alchemy API 基础 URL 模板（按链区分）
var chainBaseURLs = map[string]string{
	"ethereum": "https://eth-mainnet.g.alchemy.com",
	"polygon":  "https://polygon-mainnet.g.alchemy.com",
}

// Client Alchemy NFT 客户端
type Client struct {
	apiKey      string
	httpClient  *http.Client
	priceClient integrations.PriceClient // 用于 ETH → USD 转换
}

// NewClient 创建 Alchemy 客户端
func NewClient(apiKey string, priceClient integrations.PriceClient) *Client {
	return &Client{
		apiKey:      apiKey,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		priceClient: priceClient,
	}
}

// GetNFTs 获取地址持有的所有 NFT
func (c *Client) GetNFTs(ctx context.Context, ownerAddress string, chain string) ([]NFT, error) {
	baseURL, ok := chainBaseURLs[chain]
	if !ok {
		return nil, fmt.Errorf("不支持的链: %s", chain)
	}

	url := fmt.Sprintf("%s/nft/v3/%s/getNFTsForOwner?owner=%s&withMetadata=true&pageSize=100",
		baseURL, c.apiKey, ownerAddress)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Alchemy API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alchemy API 返回错误: %d", resp.StatusCode)
	}

	var result alchemyNFTResponse
	if err := integrations.ReadJSON(resp, &result); err != nil {
		return nil, err
	}

	// 转换为内部 NFT 类型
	nfts := make([]NFT, 0, len(result.OwnedNfts))
	for _, n := range result.OwnedNfts {
		nfts = append(nfts, NFT{
			ContractAddress: n.Contract.Address,
			TokenID:         n.TokenID,
			Name:            n.Name,
			Description:     n.Description,
			ImageURL:        n.Image.CachedURL,
			Collection:      n.Collection.Name,
			CollectionSlug:  n.Collection.Slug,
			Chain:           chain,
		})
	}

	return nfts, nil
}

// GetFloorPrice 获取收藏集 Floor Price
func (c *Client) GetFloorPrice(ctx context.Context, contractAddress string, chain string) (float64, string, error) {
	baseURL, ok := chainBaseURLs[chain]
	if !ok {
		return 0, "", fmt.Errorf("不支持的链: %s", chain)
	}

	url := fmt.Sprintf("%s/nft/v3/%s/getFloorPrice?contractAddress=%s",
		baseURL, c.apiKey, contractAddress)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("请求 Floor Price 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("Alchemy API 返回错误: %d", resp.StatusCode)
	}

	var result alchemyFloorPriceResponse
	if err := integrations.ReadJSON(resp, &result); err != nil {
		return 0, "", err
	}

	return result.OpenSea.FloorPrice, result.OpenSea.PriceCurrency, nil
}

// GetFloorPriceUSD 获取 Floor Price 并转换为 USD
func (c *Client) GetFloorPriceUSD(ctx context.Context, contractAddress string, chain string) (float64, float64, error) {
	floorPrice, currency, err := c.GetFloorPrice(ctx, contractAddress, chain)
	if err != nil {
		return 0, 0, err
	}

	if floorPrice == 0 {
		return 0, 0, nil
	}

	// 转换为 USD
	priceSymbol := "ETH"
	if currency == "MATIC" {
		priceSymbol = "MATIC"
	}

	var floorPriceUSD float64
	if c.priceClient != nil {
		tokenPrice, err := c.priceClient.GetPrice(ctx, priceSymbol)
		if err == nil {
			floorPriceUSD = floorPrice * tokenPrice
		}
	}

	return floorPrice, floorPriceUSD, nil
}

// SupportedChains 返回支持的链
func SupportedChains() []string {
	chains := make([]string, 0, len(chainBaseURLs))
	for k := range chainBaseURLs {
		chains = append(chains, k)
	}
	return chains
}
