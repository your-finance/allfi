/**
 * NFT 状态管理 Store
 * 管理 NFT 资产数据、收藏集分组、Floor Price 估值
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { nftService } from '../api/nftService.js'

export const useNFTStore = defineStore('nft', () => {
  // 状态
  const nfts = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // 是否将 NFT 价值计入总资产（用户可配置）
  const includeInTotal = ref(false)

  // NFT 总数
  const totalCount = computed(() => nfts.value.length)

  // 总 Floor Price 估值（USD）
  const totalFloorValue = computed(() => {
    return nfts.value.reduce((sum, n) => sum + (n.estimated_value || n.floorPriceUSD || 0), 0)
  })

  // 收藏集分组
  const collections = computed(() => {
    const groups = {}
    for (const nft of nfts.value) {
      if (!groups[nft.collection]) {
        groups[nft.collection] = {
          name: nft.collection,
          slug: nft.collectionSlug,
          count: 0,
          totalFloorUSD: 0,
          floorPrice: nft.floorPrice,
          floorCurrency: nft.floorCurrency,
          chain: nft.chain,
          nfts: []
        }
      }
      groups[nft.collection].count++
      groups[nft.collection].totalFloorUSD += (nft.estimated_value || nft.floorPriceUSD || 0)
      groups[nft.collection].nfts.push(nft)
    }
    // 按总价值降序
    return Object.values(groups).sort((a, b) => b.totalFloorUSD - a.totalFloorUSD)
  })

  // 收藏集数量
  const collectionCount = computed(() => collections.value.length)

  /**
   * 获取 NFT 数据
   */
  async function fetchNFTs() {
    isLoading.value = true
    error.value = null

    try {
      const result = await nftService.getNFTs()
      // 后端返回 { assets: [...], total_value }，Mock 直接返回数组
      nfts.value = Array.isArray(result) ? result : (result.assets || [])
    } catch (err) {
      error.value = err.message || '加载 NFT 数据失败'
      console.error('加载 NFT 数据失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 切换是否计入总资产
   */
  function toggleIncludeInTotal() {
    includeInTotal.value = !includeInTotal.value
    // 持久化到 localStorage
    try {
      localStorage.setItem('allfi_nft_include_total', String(includeInTotal.value))
    } catch {}
  }

  /**
   * 初始化（读取本地配置）
   */
  function initialize() {
    try {
      const saved = localStorage.getItem('allfi_nft_include_total')
      if (saved === 'true') {
        includeInTotal.value = true
      }
    } catch {}
  }

  /**
   * 重置
   */
  function reset() {
    nfts.value = []
    error.value = null
  }

  return {
    nfts,
    isLoading,
    error,
    includeInTotal,
    totalCount,
    totalFloorValue,
    collections,
    collectionCount,
    fetchNFTs,
    toggleIncludeInTotal,
    initialize,
    reset,
  }
})
