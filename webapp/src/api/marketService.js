/**
 * 市场数据 API 服务
 * 提供多链 Gas 价格等市场数据查询
 */
import { get } from './client.js'

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'

// 模拟网络延迟
const simulateDelay = (ms = 300) => new Promise(r => setTimeout(r, ms))

/**
 * 获取拥堵等级
 * @param {number} baseFee - 基础费用 (Gwei)
 * @returns {string} low / medium / high
 */
function getLevel(baseFee) {
  if (baseFee < 15) return 'low'
  if (baseFee < 30) return 'medium'
  return 'high'
}

/**
 * 获取多链 Gas 价格
 * @returns {Promise<{prices: Array<{chain, low, standard, fast, instant, base_fee, level}>, updated_at: number}>}
 */
export async function getGasPrice() {
  if (USE_MOCK) {
    await simulateDelay(200)
    const ethBase = 8 + Math.random() * 20
    return {
      prices: [
        {
          chain: 'Ethereum',
          low: Math.round(ethBase * 10) / 10,
          standard: Math.round((ethBase + 2) * 10) / 10,
          fast: Math.round((ethBase + 5) * 10) / 10,
          instant: Math.round((ethBase + 10) * 10) / 10,
          base_fee: Math.round(ethBase * 10) / 10,
          level: getLevel(ethBase)
        },
        {
          chain: 'BSC',
          low: 1.0,
          standard: 3.0,
          fast: 5.0,
          instant: 8.0,
          base_fee: 1.0,
          level: 'low'
        },
        {
          chain: 'Polygon',
          low: 25.0,
          standard: 30.0,
          fast: 50.0,
          instant: 80.0,
          base_fee: 25.0,
          level: 'medium'
        }
      ],
      updated_at: Math.floor(Date.now() / 1000)
    }
  }

  return get('/market/gas')
}

export const marketService = {
  getGasPrice
}
