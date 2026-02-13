/**
 * NFT API 服务
 * 管理 NFT 资产数据的获取，支持 Mock/Real 自动切换
 */

import { get } from './client.js'
import * as mockNFTData from '../data/mockNFTData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const nftService = {
  /**
   * 获取所有 NFT
   * @returns {Promise<{assets: Array, total_value: number}>}
   */
  async getNFTs() {
    if (USE_MOCK) {
      await simulateDelay(400)
      return mockNFTData.getNFTs()
    }
    return get('/nft/assets')
  },

  /**
   * 按收藏集筛选 NFT（本地过滤）
   * @param {string} collection - 收藏集名称
   * @returns {Promise<Array>}
   */
  async getNFTsByCollection(collection) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockNFTData.getNFTsByCollection(collection)
    }
    // 后端 /nft/assets 无 collection 参数，获取全部后本地过滤
    const result = await get('/nft/assets')
    if (collection && result.assets) {
      return result.assets.filter(nft => nft.collection === collection)
    }
    return result.assets || []
  },

  /**
   * 获取收藏集统计
   * @returns {Promise<Array>}
   */
  async getCollectionStats() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockNFTData.getCollectionStats()
    }
    const result = await get('/nfts/collections')
    // 后端返回 {collections: [...], total_count, total_value} 或直接数组
    return Array.isArray(result) ? result : (result.collections || [])
  }
}
