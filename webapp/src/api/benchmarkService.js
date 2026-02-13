/**
 * 基准对比 API 服务
 * 获取用户收益率与 BTC/ETH 等基准指数的对比数据
 * 替代原有的排行榜功能，适用于自托管单用户场景
 */

import { get } from './client.js'
import * as mockData from '../data/mockBenchmarkData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const benchmarkService = {
  /**
   * 获取基准对比数据
   * @param {string} range - 时间范围: 7d / 30d / 90d / 1y
   * @returns {Promise<Object>}
   */
  async getBenchmark(range = '30d') {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.getBenchmark(range)
    }
    return get('/benchmark', { range })
  },
}
