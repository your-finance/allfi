/**
 * 成就 API 服务
 * 管理成就数据的获取，支持 Mock/Real 自动切换
 */

import { get, post } from './client.js'
import * as mockData from '../data/mockAchievementData.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

export const achievementService = {
  /**
   * 获取所有成就
   * @returns {Promise<Array>}
   */
  async getAchievements() {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.getAchievements()
    }
    const result = await get('/achievements')
    return Array.isArray(result) ? result : (result.achievements || [])
  },

  /**
   * 检查成就进度（触发后端重新评估）
   * @returns {Promise<Object>}
   */
  async checkProgress() {
    if (USE_MOCK) {
      await simulateDelay(500)
      return { checked: true, newly_unlocked: [] }
    }
    return post('/achievements/check')
  }
}
