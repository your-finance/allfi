/**
 * 策略 API 服务
 * 管理自动化策略数据的获取，支持 Mock/Real 自动切换
 */

import { get, post, put, del } from './client.js'
import * as mockStrategyData from '../data/mockStrategyData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const strategyService = {
  /**
   * 获取所有策略
   * @returns {Promise<Array>}
   */
  async getStrategies() {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockStrategyData.getStrategies()
    }
    const result = await get('/strategies')
    return Array.isArray(result) ? result : (result.strategies || [])
  },

  /**
   * 获取再平衡建议
   * @param {string} strategyId
   * @returns {Promise<Object>}
   */
  async getRebalanceSuggestion(strategyId) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockStrategyData.getRebalanceSuggestion(strategyId)
    }
    return get(`/strategies/${strategyId}/rebalance`)
  },

  /**
   * 创建策略
   * @param {Object} data - 策略数据
   * @returns {Promise<Object>}
   */
  async createStrategy(data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return { id: 'stg-' + Date.now(), ...data, status: 'active', createdAt: new Date().toISOString() }
    }
    const result = await post('/strategies', data)
    return result.strategy || result
  },

  /**
   * 更新策略
   * @param {number|string} id - 策略 ID
   * @param {Object} data - 更新数据
   * @returns {Promise<Object>}
   */
  async updateStrategy(id, data) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return { id, ...data }
    }
    const result = await put(`/strategies/${id}`, data)
    return result.strategy || result
  },

  /**
   * 删除策略
   * @param {number|string} id - 策略 ID
   * @returns {Promise<void>}
   */
  async deleteStrategy(id) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return true
    }
    return del(`/strategies/${id}`)
  },

  /**
   * 获取策略分析
   * @param {number|string} id - 策略 ID
   * @returns {Promise<Object>}
   */
  async getAnalysis(id) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return {
        strategyId: id,
        performance: { total_return: (Math.random() - 0.3) * 20, sharpe_ratio: 0.5 + Math.random() * 1.5 },
        risk: { max_drawdown: -(5 + Math.random() * 15), volatility: 10 + Math.random() * 20 }
      }
    }
    const result = await get(`/strategies/${id}/analysis`)
    return result.analysis || result
  }
}
