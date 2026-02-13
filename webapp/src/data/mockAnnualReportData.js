/**
 * 年度投资报告 Mock 数据
 * 包含全年收益、月度曲线、资产配置、里程碑、投资风格等
 */

/**
 * 获取年度报告数据
 * @param {number} year - 年份
 * @returns {Object} 年度报告
 */
export function getAnnualReport(year = 2025) {
  return {
    year,
    // 第 1 页：年度总览
    summary: {
      totalReturn: 45.2,
      totalReturnValue: 28680,
      startValue: 63440,
      endValue: 92120,
      bestMonth: { month: 3, return: 18.5 },
      worstMonth: { month: 6, return: -8.2 },
      tradingDays: 365,
      totalTransactions: 127,
      totalFeesPaid: 892.50,
    },

    // 第 2 页：月度收益曲线
    monthlyReturns: [
      { month: 1, return: 5.2, value: 66742 },
      { month: 2, return: 3.8, value: 69279 },
      { month: 3, return: 18.5, value: 82095 },
      { month: 4, return: -2.1, value: 80371 },
      { month: 5, return: 4.6, value: 84068 },
      { month: 6, return: -8.2, value: 77175 },
      { month: 7, return: 6.3, value: 82037 },
      { month: 8, return: 1.9, value: 83596 },
      { month: 9, return: 7.8, value: 90116 },
      { month: 10, return: -1.5, value: 88764 },
      { month: 11, return: 2.4, value: 90894 },
      { month: 12, return: 1.3, value: 92120 },
    ],

    // 第 3 页：资产配置变化
    allocationChanges: {
      start: [
        { symbol: 'BTC', pct: 45 },
        { symbol: 'ETH', pct: 25 },
        { symbol: 'USDC', pct: 20 },
        { symbol: 'SOL', pct: 5 },
        { symbol: '其他', pct: 5 },
      ],
      end: [
        { symbol: 'BTC', pct: 52 },
        { symbol: 'ETH', pct: 18 },
        { symbol: 'USDC', pct: 12 },
        { symbol: 'SOL', pct: 10 },
        { symbol: '其他', pct: 8 },
      ],
    },

    // 第 4 页：关键里程碑
    milestones: [
      {
        date: `${year}-01-15`,
        title: '开启投资之旅',
        description: '首次在 AllFi 记录资产，总值 $63,440',
        icon: 'start',
      },
      {
        date: `${year}-03-12`,
        title: '单月收益 +18.5%',
        description: '3 月 BTC 大涨，创下年内最佳月收益',
        icon: 'rocket',
      },
      {
        date: `${year}-05-20`,
        title: '资产突破 $80,000',
        description: '持仓首次突破 8 万美元大关',
        icon: 'milestone',
      },
      {
        date: `${year}-09-08`,
        title: '资产突破 $90,000',
        description: '9 月市场回暖，资产再创新高',
        icon: 'milestone',
      },
      {
        date: `${year}-12-31`,
        title: '年终总结',
        description: '全年收益 +45.2%，跑赢 BTC 基准',
        icon: 'trophy',
      },
    ],

    // 第 5 页：投资风格 + 总结
    styleTag: 'steady',
    styleScore: {
      diversification: 78,
      rebalanceFreq: 65,
      holdingPeriod: 82,
      riskControl: 71,
    },
    annualStar: {
      symbol: 'SOL',
      return: 120.5,
      reason: '年初配置 5%，年底涨到 10%，单币种收益最高',
    },
    annualRegret: {
      symbol: 'DOGE',
      return: -35.2,
      reason: '追高买入，未设止损，年底清仓',
    },

    // 基准对比
    benchmarks: {
      btc: 38.6,
      eth: 22.4,
      sp500: 18.2,
    },
  }
}
