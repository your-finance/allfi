// Package alchemy Alchemy NFT 客户端测试
package alchemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSupportedChains 测试支持的链列表
func TestSupportedChains(t *testing.T) {
	chains := SupportedChains()
	assert.Contains(t, chains, "ethereum")
	assert.Contains(t, chains, "polygon")
	assert.Len(t, chains, 2)
}

// TestNewClient 测试客户端创建
func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", nil)
	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.apiKey)
}

// TestNFTCollectionGrouping 测试 NFT 收藏集分组逻辑
func TestNFTCollectionGrouping(t *testing.T) {
	nfts := []NFT{
		{Collection: "Bored Ape Yacht Club", CollectionSlug: "bayc", ContractAddress: "0x1", FloorPriceUSD: 50000},
		{Collection: "Bored Ape Yacht Club", CollectionSlug: "bayc", ContractAddress: "0x1", FloorPriceUSD: 50000},
		{Collection: "CryptoPunks", CollectionSlug: "cryptopunks", ContractAddress: "0x2", FloorPriceUSD: 80000},
	}

	// 按收藏集分组
	groups := make(map[string]*NFTCollection)
	for _, nft := range nfts {
		if _, ok := groups[nft.Collection]; !ok {
			groups[nft.Collection] = &NFTCollection{
				Name: nft.Collection,
				Slug: nft.CollectionSlug,
			}
		}
		groups[nft.Collection].Count++
		groups[nft.Collection].TotalValueUSD += nft.FloorPriceUSD
	}

	assert.Len(t, groups, 2)
	assert.Equal(t, 2, groups["Bored Ape Yacht Club"].Count)
	assert.Equal(t, float64(100000), groups["Bored Ape Yacht Club"].TotalValueUSD)
	assert.Equal(t, 1, groups["CryptoPunks"].Count)
}
