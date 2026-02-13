/**
 * 基准对比 Mock 数据
 * 用户收益率 vs BTC/ETH 基准指数
 */

/**
 * 获取基准对比数据
 * @param {string} period - 周期: 7d / 30d / 90d / 1y
 * @returns {Object} 基准对比结果
 */
export function getBenchmark(period = '30d') {
  // 根据不同周期模拟不同的收益率数据
  const data = {
    '7d': {
      user_return: 2.3,
      benchmarks: [
        { name: 'Bitcoin', return_percent: 1.8, user_outperform: 0.5 },
        { name: 'Ethereum', return_percent: 3.1, user_outperform: -0.8 },
      ],
    },
    '30d': {
      user_return: 8.5,
      benchmarks: [
        { name: 'Bitcoin', return_percent: 12.3, user_outperform: -3.8 },
        { name: 'Ethereum', return_percent: 6.7, user_outperform: 1.8 },
      ],
    },
    '90d': {
      user_return: 15.2,
      benchmarks: [
        { name: 'Bitcoin', return_percent: 22.1, user_outperform: -6.9 },
        { name: 'Ethereum', return_percent: 18.4, user_outperform: -3.2 },
      ],
    },
    '1y': {
      user_return: 42.8,
      benchmarks: [
        { name: 'Bitcoin', return_percent: 65.3, user_outperform: -22.5 },
        { name: 'Ethereum', return_percent: 38.9, user_outperform: 3.9 },
      ],
    },
  }

  const d = data[period] || data['30d']
  const now = new Date()
  const days = period === '7d' ? 7 : period === '30d' ? 30 : period === '90d' ? 90 : 365
  const start = new Date(now.getTime() - days * 24 * 60 * 60 * 1000)

  return {
    period,
    ...d,
    start_date: start.toISOString(),
    end_date: now.toISOString(),
  }
}
