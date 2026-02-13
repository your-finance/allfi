/**
 * 自动化策略 Mock 数据
 * 策略类型：rebalance（再平衡）、dca（定投）、alert（止盈止损）
 */

export const mockStrategies = [
  // 再平衡策略
  {
    id: 'stg-001',
    type: 'rebalance',
    name: '核心资产再平衡',
    status: 'active',        // active / paused / triggered
    createdAt: '2025-12-01',
    lastTriggeredAt: '2026-02-05',
    config: {
      targets: [
        { symbol: 'BTC', targetPct: 50, currentPct: 58.2 },
        { symbol: 'ETH', targetPct: 30, currentPct: 24.5 },
        { symbol: 'USDC', targetPct: 20, currentPct: 17.3 },
      ],
      deviationThreshold: 5, // 偏离超过 5% 触发提醒
    },
  },
  // 定投策略
  {
    id: 'stg-002',
    type: 'dca',
    name: '每月定投 BTC',
    status: 'active',
    createdAt: '2025-10-15',
    lastTriggeredAt: '2026-02-01',
    config: {
      symbol: 'BTC',
      amount: 1000,          // 每次 $1000
      currency: 'USDC',
      frequency: 'monthly',  // monthly / weekly / biweekly
      nextDate: '2026-03-01',
      totalInvested: 5000,
      totalPurchased: 0.052,
      avgPrice: 96153.85,
    },
  },
  // 止盈止损策略
  {
    id: 'stg-003',
    type: 'alert',
    name: 'BTC 止盈提醒',
    status: 'active',
    createdAt: '2026-01-10',
    lastTriggeredAt: null,
    config: {
      symbol: 'BTC',
      direction: 'above',    // above / below
      targetPrice: 120000,
      currentPrice: 96000,
      note: 'BTC 达到 12 万时考虑减仓',
    },
  },
  {
    id: 'stg-004',
    type: 'alert',
    name: 'ETH 止损提醒',
    status: 'paused',
    createdAt: '2026-01-20',
    lastTriggeredAt: null,
    config: {
      symbol: 'ETH',
      direction: 'below',
      targetPrice: 1800,
      currentPrice: 2300,
      note: 'ETH 跌破 1800 时考虑止损',
    },
  },
  {
    id: 'stg-005',
    type: 'dca',
    name: '每周定投 ETH',
    status: 'paused',
    createdAt: '2025-11-01',
    lastTriggeredAt: '2026-01-15',
    config: {
      symbol: 'ETH',
      amount: 200,
      currency: 'USDC',
      frequency: 'weekly',
      nextDate: '2026-02-17',
      totalInvested: 2800,
      totalPurchased: 1.22,
      avgPrice: 2295.08,
    },
  },
]

/**
 * 获取所有策略
 * @returns {Array}
 */
export function getStrategies() {
  return [...mockStrategies]
}

/**
 * 获取再平衡建议
 * @param {string} strategyId
 * @returns {Object}
 */
export function getRebalanceSuggestion(strategyId) {
  const strategy = mockStrategies.find(s => s.id === strategyId)
  if (!strategy || strategy.type !== 'rebalance') return null

  const suggestions = []
  for (const target of strategy.config.targets) {
    const diff = target.currentPct - target.targetPct
    if (Math.abs(diff) > strategy.config.deviationThreshold) {
      suggestions.push({
        symbol: target.symbol,
        action: diff > 0 ? 'sell' : 'buy',
        diffPct: parseFloat(Math.abs(diff).toFixed(1)),
        description: diff > 0
          ? `减持 ${target.symbol} ${Math.abs(diff).toFixed(1)}%`
          : `增持 ${target.symbol} ${Math.abs(diff).toFixed(1)}%`
      })
    }
  }
  return { strategyId, suggestions }
}
