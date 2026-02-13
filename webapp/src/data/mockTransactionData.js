/**
 * 交易记录 Mock 数据
 * 模拟 CEX、链上、手动资产的统一交易记录
 * 类型：buy（买入）、sell（卖出）、transfer（转账）、swap（兑换）、deposit（充值）、withdraw（提现）
 */

// 生成模拟交易记录
const now = Date.now()
const DAY = 86400000

export const mockTransactions = [
  // === CEX 交易 ===
  {
    id: 'tx-001',
    type: 'buy',
    timestamp: now - 0.5 * DAY,
    from: { symbol: 'USDC', amount: 5000 },
    to: { symbol: 'BTC', amount: 0.052 },
    fee: { amount: 5.0, currency: 'USDC' },
    chain: null,
    source: 'Binance',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-002',
    type: 'sell',
    timestamp: now - 1.2 * DAY,
    from: { symbol: 'ETH', amount: 2.5 },
    to: { symbol: 'USDC', amount: 5750 },
    fee: { amount: 5.75, currency: 'USDC' },
    chain: null,
    source: 'Binance',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-003',
    type: 'buy',
    timestamp: now - 2 * DAY,
    from: { symbol: 'USDC', amount: 1200 },
    to: { symbol: 'SOL', amount: 8.0 },
    fee: { amount: 1.2, currency: 'USDC' },
    chain: null,
    source: 'OKX',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-004',
    type: 'deposit',
    timestamp: now - 3 * DAY,
    from: { symbol: 'USDC', amount: 10000 },
    to: { symbol: 'USDC', amount: 10000 },
    fee: { amount: 0, currency: 'USDC' },
    chain: 'ETH',
    source: 'Binance',
    sourceType: 'cex',
    txHash: '0xabc123...def456',
    note: '从链上充值到 Binance',
  },
  {
    id: 'tx-005',
    type: 'withdraw',
    timestamp: now - 4.5 * DAY,
    from: { symbol: 'ETH', amount: 5.0 },
    to: { symbol: 'ETH', amount: 5.0 },
    fee: { amount: 0.003, currency: 'ETH' },
    chain: 'ETH',
    source: 'OKX',
    sourceType: 'cex',
    txHash: '0x789abc...123def',
    note: '提现到个人钱包',
  },
  // === 链上交易 ===
  {
    id: 'tx-006',
    type: 'swap',
    timestamp: now - 1 * DAY,
    from: { symbol: 'ETH', amount: 1.5 },
    to: { symbol: 'USDC', amount: 3450 },
    fee: { amount: 0.0045, currency: 'ETH' },
    chain: 'ETH',
    source: 'Uniswap',
    sourceType: 'blockchain',
    txHash: '0xdef789...abc123',
    note: '',
  },
  {
    id: 'tx-007',
    type: 'transfer',
    timestamp: now - 2.5 * DAY,
    from: { symbol: 'USDC', amount: 2000 },
    to: { symbol: 'USDC', amount: 2000 },
    fee: { amount: 0.002, currency: 'ETH' },
    chain: 'ETH',
    source: 'Ethereum',
    sourceType: 'blockchain',
    txHash: '0x456def...789abc',
    note: '转账给朋友',
  },
  {
    id: 'tx-008',
    type: 'swap',
    timestamp: now - 5 * DAY,
    from: { symbol: 'BNB', amount: 10 },
    to: { symbol: 'CAKE', amount: 320 },
    fee: { amount: 0.005, currency: 'BNB' },
    chain: 'BSC',
    source: 'PancakeSwap',
    sourceType: 'blockchain',
    txHash: '0xbsc123...456abc',
    note: '',
  },
  {
    id: 'tx-009',
    type: 'transfer',
    timestamp: now - 6 * DAY,
    from: { symbol: 'SOL', amount: 50 },
    to: { symbol: 'SOL', amount: 50 },
    fee: { amount: 0.00025, currency: 'SOL' },
    chain: 'SOL',
    source: 'Solana',
    sourceType: 'blockchain',
    txHash: '5vSol...abc123',
    note: '转移到冷钱包',
  },
  // === 更多 CEX 交易 ===
  {
    id: 'tx-010',
    type: 'buy',
    timestamp: now - 7 * DAY,
    from: { symbol: 'USDC', amount: 3000 },
    to: { symbol: 'ETH', amount: 1.3 },
    fee: { amount: 3.0, currency: 'USDC' },
    chain: null,
    source: 'Coinbase',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-011',
    type: 'sell',
    timestamp: now - 8 * DAY,
    from: { symbol: 'BTC', amount: 0.1 },
    to: { symbol: 'USDC', amount: 9600 },
    fee: { amount: 9.6, currency: 'USDC' },
    chain: null,
    source: 'Binance',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-012',
    type: 'swap',
    timestamp: now - 10 * DAY,
    from: { symbol: 'USDC', amount: 5000 },
    to: { symbol: 'USDT', amount: 4998 },
    fee: { amount: 0.001, currency: 'ETH' },
    chain: 'ETH',
    source: 'Curve',
    sourceType: 'blockchain',
    txHash: '0xcurve...abc123',
    note: '稳定币兑换',
  },
  {
    id: 'tx-013',
    type: 'buy',
    timestamp: now - 12 * DAY,
    from: { symbol: 'USDC', amount: 800 },
    to: { symbol: 'MATIC', amount: 1000 },
    fee: { amount: 0.8, currency: 'USDC' },
    chain: null,
    source: 'Binance',
    sourceType: 'cex',
    txHash: null,
    note: '',
  },
  {
    id: 'tx-014',
    type: 'withdraw',
    timestamp: now - 15 * DAY,
    from: { symbol: 'USDC', amount: 3000 },
    to: { symbol: 'USDC', amount: 3000 },
    fee: { amount: 1.0, currency: 'USDC' },
    chain: 'ETH',
    source: 'Binance',
    sourceType: 'cex',
    txHash: '0xwithdraw...abc',
    note: '提现到 DeFi',
  },
  {
    id: 'tx-015',
    type: 'transfer',
    timestamp: now - 20 * DAY,
    from: { symbol: 'BTC', amount: 0.05 },
    to: { symbol: 'BTC', amount: 0.05 },
    fee: { amount: 0.00005, currency: 'BTC' },
    chain: 'BTC',
    source: 'Bitcoin',
    sourceType: 'blockchain',
    txHash: 'bc1qxy2...abc',
    note: '转入冷钱包',
  },
  // === 手动记录 ===
  {
    id: 'tx-016',
    type: 'buy',
    timestamp: now - 25 * DAY,
    from: { symbol: 'CNY', amount: 50000 },
    to: { symbol: 'USDC', amount: 7000 },
    fee: { amount: 0, currency: 'CNY' },
    chain: null,
    source: '手动记录',
    sourceType: 'manual',
    txHash: null,
    note: '场外购入 USDC',
  },
]

/**
 * 获取交易记录（支持筛选和分页）
 * @param {Object} options - 筛选选项
 * @returns {Object} 分页结果
 */
export function getTransactions(options = {}) {
  const {
    type = 'all',
    sourceType = 'all',
    timeRange = 'all',
    search = '',
    page = 1,
    pageSize = 20,
  } = options

  let filtered = [...mockTransactions]

  // 按类型筛选
  if (type !== 'all') {
    filtered = filtered.filter(tx => tx.type === type)
  }

  // 按来源类型筛选
  if (sourceType !== 'all') {
    filtered = filtered.filter(tx => tx.sourceType === sourceType)
  }

  // 按时间范围筛选
  if (timeRange !== 'all') {
    const days = parseInt(timeRange)
    if (!isNaN(days)) {
      const cutoff = now - days * DAY
      filtered = filtered.filter(tx => tx.timestamp >= cutoff)
    }
  }

  // 搜索（按代币名称或来源）
  if (search) {
    const q = search.toLowerCase()
    filtered = filtered.filter(tx =>
      tx.from.symbol.toLowerCase().includes(q) ||
      tx.to.symbol.toLowerCase().includes(q) ||
      tx.source.toLowerCase().includes(q) ||
      (tx.note && tx.note.toLowerCase().includes(q))
    )
  }

  // 按时间降序排序
  filtered.sort((a, b) => b.timestamp - a.timestamp)

  // 分页
  const total = filtered.length
  const totalPages = Math.ceil(total / pageSize)
  const start = (page - 1) * pageSize
  const items = filtered.slice(start, start + pageSize)

  return { items, total, totalPages, page }
}
