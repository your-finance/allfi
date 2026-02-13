/**
 * DeFi API 服务
 * 管理 DeFi 协议仓位数据的获取，支持 Mock/Real 自动切换
 */

import { get } from './client.js'
import * as mockDefiData from '../data/mockDefiData.js'

// 是否使用 Mock 数据
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'

// 模拟网络延迟
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const defiService = {
  /**
   * 获取所有 DeFi 仓位
   * @returns {Promise<Array>} DeFi 仓位列表
   */
  async getPositions() {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockDefiData.getDefiPositions()
    }
    const result = await get('/defi/positions')
    // 后端返回 {positions: [...], total_value} 或直接数组
    return Array.isArray(result) ? result : (result.positions || [])
  },

  /**
   * 按协议筛选 DeFi 仓位
   * @param {string} protocol - 协议名称
   * @returns {Promise<Array>} 筛选后的仓位列表
   */
  async getPositionsByProtocol(protocol) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockDefiData.getDefiPositionsByProtocol(protocol)
    }
    const result = await get('/defi/positions', { protocol })
    return Array.isArray(result) ? result : (result.positions || [])
  },

  /**
   * 按类型筛选 DeFi 仓位
   * @param {string} type - 仓位类型 (lp/staking/lending/all)
   * @returns {Promise<Array>} 筛选后的仓位列表
   */
  async getPositionsByType(type) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockDefiData.getDefiPositionsByType(type)
    }
    const result = await get('/defi/positions', { type })
    return Array.isArray(result) ? result : (result.positions || [])
  },

  /**
   * 获取 DeFi 仓位统计
   * @returns {Promise<object>} 统计数据
   */
  async getStats() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockDefiData.getDefiStats()
    }
    return get('/defi/stats')
  },

  /**
   * 获取支持的 DeFi 协议列表
   * @returns {Promise<object[]>} 协议列表 [{name, chains, types, is_active}]
   */
  async getProtocols() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return [
        { name: 'Lido', chains: ['Ethereum'], types: ['staking'], is_active: true },
        { name: 'RocketPool', chains: ['Ethereum'], types: ['staking'], is_active: true },
        { name: 'Aave', chains: ['Ethereum', 'Polygon', 'Arbitrum'], types: ['lending', 'borrowing'], is_active: true },
        { name: 'Compound', chains: ['Ethereum'], types: ['lending', 'borrowing'], is_active: true },
        { name: 'Uniswap V3', chains: ['Ethereum', 'Polygon', 'Arbitrum', 'Optimism'], types: ['liquidity'], is_active: true },
        { name: 'Uniswap V2', chains: ['Ethereum'], types: ['liquidity'], is_active: true },
        { name: 'Curve', chains: ['Ethereum'], types: ['liquidity'], is_active: true },
      ]
    }
    const result = await get('/defi/protocols')
    return result.protocols || result || []
  }
}
