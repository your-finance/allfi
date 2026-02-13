/**
 * 交易记录 API 服务
 * 管理交易记录数据的获取、同步设置，支持 Mock/Real 自动切换
 */

import { get, post, put } from './client.js'
import * as mockTransactionData from '../data/mockTransactionData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const transactionService = {
  /**
   * 获取交易记录（分页 + 筛选）
   * @param {Object} options - 筛选选项
   * @returns {Promise<Object>}
   */
  async getTransactions(options = {}) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockTransactionData.getTransactions(options)
    }
    return get('/transactions', options)
  },

  /**
   * 获取交易同步设置
   * @returns {Promise<Object>} 同步设置
   */
  async getSyncSettings() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return {
        enabled: false,
        interval_minutes: 360,
        lookback_days: 90,
        sources: []
      }
    }
    const result = await get('/settings/tx-sync')
    // 后端返回 {settings: {auto_sync, sync_interval, last_sync_at}} 或直接返回设置对象
    const s = result.settings || result
    return {
      enabled: s.auto_sync || false,
      interval_minutes: s.sync_interval || 360,
      lookback_days: 90,
      sources: [],
      last_sync_at: s.last_sync_at
    }
  },

  /**
   * 更新交易同步设置
   * @param {Object} settings - 更新的设置项
   * @returns {Promise<void>}
   */
  async updateSyncSettings(settings) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return
    }
    return put('/settings/tx-sync', settings)
  },

  /**
   * 手动触发交易记录同步
   * @returns {Promise<Object>} 同步结果
   */
  async syncTransactions() {
    if (USE_MOCK) {
      await simulateDelay(1500)
      return { synced_count: 12 }
    }
    return post('/transactions/sync')
  },

  /**
   * 获取交易统计
   * @returns {Promise<Object>} 交易统计数据
   */
  async getTransactionStats() {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockTransactionData.getTransactionStats?.() || {}
    }
    return get('/transactions/stats')
  },
}
