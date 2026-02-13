/**
 * NFT Mock 数据
 * 模拟 CryptoPunks、BAYC、Art Blocks 等 NFT 持仓数据
 */

export const mockNFTs = [
  {
    id: 1,
    name: 'CryptoPunk #7804',
    collection: 'CryptoPunks',
    collectionSlug: 'cryptopunks',
    tokenId: '7804',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 49.5,
    floorCurrency: 'ETH',
    floorPriceUSD: 113850,
    lastSalePrice: 42.0,
    lastSaleCurrency: 'ETH',
    rarity: 'Legendary',
    attributes: [
      { trait: 'Type', value: 'Alien' },
      { trait: 'Accessory', value: 'Pipe' }
    ]
  },
  {
    id: 2,
    name: 'Bored Ape #3749',
    collection: 'Bored Ape Yacht Club',
    collectionSlug: 'bayc',
    tokenId: '3749',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 12.8,
    floorCurrency: 'ETH',
    floorPriceUSD: 29440,
    lastSalePrice: 10.5,
    lastSaleCurrency: 'ETH',
    rarity: 'Rare',
    attributes: [
      { trait: 'Background', value: 'Blue' },
      { trait: 'Fur', value: 'Gold' },
      { trait: 'Eyes', value: 'Laser Eyes' }
    ]
  },
  {
    id: 3,
    name: 'Fidenza #313',
    collection: 'Art Blocks',
    collectionSlug: 'artblocks',
    tokenId: '313',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 8.2,
    floorCurrency: 'ETH',
    floorPriceUSD: 18860,
    lastSalePrice: 7.0,
    lastSaleCurrency: 'ETH',
    rarity: 'Epic',
    attributes: [
      { trait: 'Collection', value: 'Fidenza' },
      { trait: 'Artist', value: 'Tyler Hobbs' }
    ]
  },
  {
    id: 4,
    name: 'Azuki #5289',
    collection: 'Azuki',
    collectionSlug: 'azuki',
    tokenId: '5289',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 5.1,
    floorCurrency: 'ETH',
    floorPriceUSD: 11730,
    lastSalePrice: 4.8,
    lastSaleCurrency: 'ETH',
    rarity: 'Uncommon',
    attributes: [
      { trait: 'Type', value: 'Human' },
      { trait: 'Hair', value: 'Pink Hairband' }
    ]
  },
  {
    id: 5,
    name: 'Bored Ape #8821',
    collection: 'Bored Ape Yacht Club',
    collectionSlug: 'bayc',
    tokenId: '8821',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 12.8,
    floorCurrency: 'ETH',
    floorPriceUSD: 29440,
    lastSalePrice: 11.2,
    lastSaleCurrency: 'ETH',
    rarity: 'Common',
    attributes: [
      { trait: 'Background', value: 'Orange' },
      { trait: 'Fur', value: 'Brown' }
    ]
  },
  {
    id: 6,
    name: 'Pudgy Penguin #1234',
    collection: 'Pudgy Penguins',
    collectionSlug: 'pudgypenguins',
    tokenId: '1234',
    imageUrl: '',
    chain: 'ETH',
    chainName: 'Ethereum',
    floorPrice: 8.5,
    floorCurrency: 'ETH',
    floorPriceUSD: 19550,
    lastSalePrice: 7.8,
    lastSaleCurrency: 'ETH',
    rarity: 'Rare',
    attributes: [
      { trait: 'Background', value: 'Green' },
      { trait: 'Body', value: 'Tuxedo' }
    ]
  }
]

/**
 * 获取所有 NFT
 * @returns {Array}
 */
export function getNFTs() {
  return [...mockNFTs]
}

/**
 * 按收藏集筛选 NFT
 * @param {string} collection - 收藏集名称
 * @returns {Array}
 */
export function getNFTsByCollection(collection) {
  if (!collection || collection === 'all') return getNFTs()
  return mockNFTs.filter(n => n.collectionSlug === collection)
}

/**
 * 获取收藏集统计
 * @returns {Array}
 */
export function getCollectionStats() {
  const groups = {}
  for (const nft of mockNFTs) {
    if (!groups[nft.collection]) {
      groups[nft.collection] = {
        name: nft.collection,
        slug: nft.collectionSlug,
        count: 0,
        totalFloorUSD: 0,
        chain: nft.chain
      }
    }
    groups[nft.collection].count++
    groups[nft.collection].totalFloorUSD += nft.floorPriceUSD
  }
  return Object.values(groups).sort((a, b) => b.totalFloorUSD - a.totalFloorUSD)
}
