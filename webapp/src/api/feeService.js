/**
 * 费用分析 API 服务
 * 管理费用数据的获取，支持 Mock/Real 自动切换
 */

import { get } from './client.js'
import * as mockFeeData from '../data/mockFeeData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const feeService = {
  /**
   * 获取费用分析
   * @param {string} timeRange - 时间范围 (7d/30d/90d/1y)
   * @param {string} currency - 计价货币
   * @returns {Promise<Object>}
   */
  async getFeeSummary(timeRange = '30d', currency = 'USD') {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockFeeData.getFeeSummary(timeRange)
    }
    return get('/analytics/fees', { range: timeRange, currency })
  },
}
