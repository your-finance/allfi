/**
 * 费用分析 Mock 数据
 * 模拟 CEX 手续费、Gas 费、提现费等费用数据
 */

// 费用分类汇总
const feeBreakdown = {
  cexTradeFee: 45.80,    // CEX 交易手续费（USD）
  gasFee: 128.50,        // 链上 Gas 费（USD）
  withdrawFee: 18.20,    // 提现手续费（USD）
}

// 月度费用趋势（最近 6 个月）
const monthlyTrend = [
  { month: '2025-09', cexTradeFee: 35.2, gasFee: 95.4, withdrawFee: 12.0, total: 142.6 },
  { month: '2025-10', cexTradeFee: 42.1, gasFee: 110.8, withdrawFee: 15.5, total: 168.4 },
  { month: '2025-11', cexTradeFee: 38.6, gasFee: 88.3, withdrawFee: 10.2, total: 137.1 },
  { month: '2025-12', cexTradeFee: 52.3, gasFee: 145.6, withdrawFee: 22.8, total: 220.7 },
  { month: '2026-01', cexTradeFee: 48.9, gasFee: 132.1, withdrawFee: 16.3, total: 197.3 },
  { month: '2026-02', cexTradeFee: 45.8, gasFee: 128.5, withdrawFee: 18.2, total: 192.5 },
]

// 智能建议
const suggestions = [
  {
    id: 'gas-l2',
    type: 'saving',
    titleKey: 'fee.suggestionGasL2Title',
    descKey: 'fee.suggestionGasL2Desc',
    savingEstimate: 85,
  },
  {
    id: 'batch-withdraw',
    type: 'saving',
    titleKey: 'fee.suggestionBatchTitle',
    descKey: 'fee.suggestionBatchDesc',
    savingEstimate: 12,
  },
  {
    id: 'maker-fee',
    type: 'info',
    titleKey: 'fee.suggestionMakerTitle',
    descKey: 'fee.suggestionMakerDesc',
    savingEstimate: 20,
  },
]

/**
 * 获取费用统计摘要
 * @param {string} timeRange - 时间范围：'7D' / '30D' / '90D'
 * @returns {Object}
 */
export function getFeeSummary(timeRange = '30D') {
  // 根据时间范围调整数据（模拟）
  const multiplier = timeRange === '7D' ? 0.25 : timeRange === '90D' ? 3.0 : 1.0
  const total = (feeBreakdown.cexTradeFee + feeBreakdown.gasFee + feeBreakdown.withdrawFee) * multiplier

  // 上月对比
  const lastMonthTotal = monthlyTrend[monthlyTrend.length - 2]?.total || 0
  const currentMonthTotal = monthlyTrend[monthlyTrend.length - 1]?.total || 0
  const changePercent = lastMonthTotal > 0
    ? ((currentMonthTotal - lastMonthTotal) / lastMonthTotal * 100)
    : 0

  return {
    total: parseFloat((total).toFixed(2)),
    breakdown: {
      cexTradeFee: parseFloat((feeBreakdown.cexTradeFee * multiplier).toFixed(2)),
      gasFee: parseFloat((feeBreakdown.gasFee * multiplier).toFixed(2)),
      withdrawFee: parseFloat((feeBreakdown.withdrawFee * multiplier).toFixed(2)),
    },
    changePercent: parseFloat(changePercent.toFixed(1)),
    monthlyTrend,
    suggestions,
  }
}
