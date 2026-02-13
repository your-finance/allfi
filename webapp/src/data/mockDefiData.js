/**
 * DeFi Mock 数据
 * 模拟 Uniswap LP、Aave 借贷、Lido 质押等 DeFi 仓位数据
 */

// DeFi 仓位 Mock 数据
export const mockDefiPositions = [
  {
    id: 1,
    protocol: 'Uniswap V3',
    protocolIcon: 'UNI',
    type: 'lp',
    chain: 'ETH',
    chainName: 'Ethereum',
    tokens: [
      { symbol: 'ETH', name: 'Ethereum', amount: 2.5, icon: '' },
      { symbol: 'USDC', name: 'USD Coin', amount: 5750, icon: '' }
    ],
    valueUSD: 11500,
    apy: 18.5,
    rewards: {
      token: 'UNI',
      amount: 12.3,
      valueUSD: 86.1
    },
    priceRange: { lower: 1800, upper: 3200 },
    inRange: true,
    createdAt: '2025-12-15T08:30:00Z'
  },
  {
    id: 2,
    protocol: 'Aave V3',
    protocolIcon: 'AAVE',
    type: 'lending',
    chain: 'ETH',
    chainName: 'Ethereum',
    tokens: [
      { symbol: 'USDC', name: 'USD Coin', amount: 20000, icon: '' }
    ],
    valueUSD: 20000,
    apy: 4.2,
    rewards: {
      token: 'AAVE',
      amount: 0.15,
      valueUSD: 45.0
    },
    healthFactor: 2.8,
    createdAt: '2025-11-20T10:00:00Z'
  },
  {
    id: 3,
    protocol: 'Lido',
    protocolIcon: 'LDO',
    type: 'staking',
    chain: 'ETH',
    chainName: 'Ethereum',
    tokens: [
      { symbol: 'stETH', name: 'Lido Staked ETH', amount: 5.0, icon: '' }
    ],
    valueUSD: 11500,
    apy: 3.8,
    rewards: {
      token: 'stETH',
      amount: 0.019,
      valueUSD: 43.7
    },
    createdAt: '2025-10-01T12:00:00Z'
  },
  {
    id: 4,
    protocol: 'PancakeSwap',
    protocolIcon: 'CAKE',
    type: 'lp',
    chain: 'BSC',
    chainName: 'BNB Chain',
    tokens: [
      { symbol: 'BNB', name: 'BNB', amount: 10, icon: '' },
      { symbol: 'USDT', name: 'Tether', amount: 3100, icon: '' }
    ],
    valueUSD: 6200,
    apy: 22.4,
    rewards: {
      token: 'CAKE',
      amount: 8.5,
      valueUSD: 25.5
    },
    createdAt: '2026-01-05T14:00:00Z'
  },
  {
    id: 5,
    protocol: 'Compound V3',
    protocolIcon: 'COMP',
    type: 'lending',
    chain: 'ETH',
    chainName: 'Ethereum',
    tokens: [
      { symbol: 'ETH', name: 'Ethereum', amount: 3.0, icon: '' }
    ],
    valueUSD: 6900,
    apy: 2.1,
    rewards: {
      token: 'COMP',
      amount: 0.8,
      valueUSD: 48.0
    },
    healthFactor: 3.5,
    createdAt: '2026-01-10T09:00:00Z'
  },
  {
    id: 6,
    protocol: 'Rocket Pool',
    protocolIcon: 'RPL',
    type: 'staking',
    chain: 'ETH',
    chainName: 'Ethereum',
    tokens: [
      { symbol: 'rETH', name: 'Rocket Pool ETH', amount: 2.0, icon: '' }
    ],
    valueUSD: 4600,
    apy: 3.5,
    rewards: {
      token: 'RPL',
      amount: 0.5,
      valueUSD: 15.0
    },
    createdAt: '2025-12-20T16:00:00Z'
  }
]

/**
 * 获取所有 DeFi 仓位
 * @returns {Array} DeFi 仓位列表
 */
export function getDefiPositions() {
  return [...mockDefiPositions]
}

/**
 * 按协议筛选 DeFi 仓位
 * @param {string} protocol - 协议名称
 * @returns {Array} 筛选后的仓位列表
 */
export function getDefiPositionsByProtocol(protocol) {
  return mockDefiPositions.filter(p => p.protocol === protocol)
}

/**
 * 按类型筛选 DeFi 仓位
 * @param {string} type - 仓位类型 (lp/staking/lending)
 * @returns {Array} 筛选后的仓位列表
 */
export function getDefiPositionsByType(type) {
  if (!type || type === 'all') return getDefiPositions()
  return mockDefiPositions.filter(p => p.type === type)
}

/**
 * 获取 DeFi 仓位统计
 * @returns {object} 统计数据
 */
export function getDefiStats() {
  const positions = mockDefiPositions
  const totalValue = positions.reduce((sum, p) => sum + p.valueUSD, 0)
  const totalRewards = positions.reduce((sum, p) => sum + (p.rewards?.valueUSD || 0), 0)

  // 按类型统计
  const byType = {}
  for (const p of positions) {
    if (!byType[p.type]) {
      byType[p.type] = { count: 0, value: 0 }
    }
    byType[p.type].count++
    byType[p.type].value += p.valueUSD
  }

  // 按协议统计
  const byProtocol = {}
  for (const p of positions) {
    if (!byProtocol[p.protocol]) {
      byProtocol[p.protocol] = { count: 0, value: 0 }
    }
    byProtocol[p.protocol].count++
    byProtocol[p.protocol].value += p.valueUSD
  }

  return {
    totalValue,
    totalRewards,
    positionCount: positions.length,
    byType,
    byProtocol
  }
}
