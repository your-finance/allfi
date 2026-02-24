/**
 * 交易记录 Mock 数据
 * 字段格式与后端 API 保持一致：
 * tx_type, symbol, amount, price, total, fee, fee_coin, tx_hash, timestamp(ISO)
 */

const now = new Date()
const iso = (offsetMs) => new Date(Date.now() - offsetMs).toISOString()
const DAY = 86400000

export const mockTransactions = [
  // === CEX 交易 ===
  {
    id: 1,
    source: 'Binance',
    tx_type: 'buy',
    symbol: 'BTC',
    amount: 0.052,
    price: 96153.85,
    total: 5000,
    fee: 5.0,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(0.5 * DAY),
  },
  {
    id: 2,
    source: 'Binance',
    tx_type: 'sell',
    symbol: 'ETH',
    amount: 2.5,
    price: 2300,
    total: 5750,
    fee: 5.75,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(1.2 * DAY),
  },
  {
    id: 3,
    source: 'OKX',
    tx_type: 'buy',
    symbol: 'SOL',
    amount: 8.0,
    price: 150,
    total: 1200,
    fee: 1.2,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(2 * DAY),
  },
  {
    id: 4,
    source: 'Binance',
    tx_type: 'deposit',
    symbol: 'USDC',
    amount: 10000,
    price: 1,
    total: 10000,
    fee: 0,
    fee_coin: 'USDC',
    tx_hash: '0xabc123def456',
    timestamp: iso(3 * DAY),
  },
  {
    id: 5,
    source: 'OKX',
    tx_type: 'withdraw',
    symbol: 'ETH',
    amount: 5.0,
    price: 2300,
    total: 11500,
    fee: 0.003,
    fee_coin: 'ETH',
    tx_hash: '0x789abc123def',
    timestamp: iso(4.5 * DAY),
  },
  // === 链上交易 ===
  {
    id: 6,
    source: 'Uniswap',
    tx_type: 'swap',
    symbol: 'ETH',
    amount: 1.5,
    price: 2300,
    total: 3450,
    fee: 0.0045,
    fee_coin: 'ETH',
    tx_hash: '0xdef789abc123',
    timestamp: iso(1 * DAY),
  },
  {
    id: 7,
    source: 'Ethereum',
    tx_type: 'transfer',
    symbol: 'USDC',
    amount: 2000,
    price: 1,
    total: 2000,
    fee: 0.002,
    fee_coin: 'ETH',
    tx_hash: '0x456def789abc',
    timestamp: iso(2.5 * DAY),
  },
  {
    id: 8,
    source: 'PancakeSwap',
    tx_type: 'swap',
    symbol: 'BNB',
    amount: 10,
    price: 320,
    total: 3200,
    fee: 0.005,
    fee_coin: 'BNB',
    tx_hash: '0xbsc123456abc',
    timestamp: iso(5 * DAY),
  },
  {
    id: 9,
    source: 'Solana',
    tx_type: 'transfer',
    symbol: 'SOL',
    amount: 50,
    price: 150,
    total: 7500,
    fee: 0.00025,
    fee_coin: 'SOL',
    tx_hash: '5vSolabc123',
    timestamp: iso(6 * DAY),
  },
  // === 更多 CEX 交易 ===
  {
    id: 10,
    source: 'Coinbase',
    tx_type: 'buy',
    symbol: 'ETH',
    amount: 1.3,
    price: 2307.69,
    total: 3000,
    fee: 3.0,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(7 * DAY),
  },
  {
    id: 11,
    source: 'Binance',
    tx_type: 'sell',
    symbol: 'BTC',
    amount: 0.1,
    price: 96000,
    total: 9600,
    fee: 9.6,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(8 * DAY),
  },
  {
    id: 12,
    source: 'Curve',
    tx_type: 'swap',
    symbol: 'USDC',
    amount: 5000,
    price: 1,
    total: 5000,
    fee: 0.001,
    fee_coin: 'ETH',
    tx_hash: '0xcurveabc123',
    timestamp: iso(10 * DAY),
  },
  {
    id: 13,
    source: 'Binance',
    tx_type: 'buy',
    symbol: 'MATIC',
    amount: 1000,
    price: 0.8,
    total: 800,
    fee: 0.8,
    fee_coin: 'USDC',
    tx_hash: '',
    timestamp: iso(12 * DAY),
  },
  {
    id: 14,
    source: 'Binance',
    tx_type: 'withdraw',
    symbol: 'USDC',
    amount: 3000,
    price: 1,
    total: 3000,
    fee: 1.0,
    fee_coin: 'USDC',
    tx_hash: '0xwithdrawabc',
    timestamp: iso(15 * DAY),
  },
  {
    id: 15,
    source: 'Bitcoin',
    tx_type: 'transfer',
    symbol: 'BTC',
    amount: 0.05,
    price: 96000,
    total: 4800,
    fee: 0.00005,
    fee_coin: 'BTC',
    tx_hash: 'bc1qxy2abc',
    timestamp: iso(20 * DAY),
  },
  {
    id: 16,
    source: '手动记录',
    tx_type: 'buy',
    symbol: 'USDC',
    amount: 7000,
    price: 7.14,
    total: 50000,
    fee: 0,
    fee_coin: 'CNY',
    tx_hash: '',
    timestamp: iso(25 * DAY),
  },
]

/**
 * 获取交易记录（支持筛选和分页）
 * 返回格式与后端保持一致：{ transactions, total, page, page_size }
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
    filtered = filtered.filter(tx => tx.tx_type === type)
  }

  // 按来源类型筛选（mock 中通过 source 名称判断）
  if (sourceType !== 'all') {
    const cexSources = ['Binance', 'OKX', 'Coinbase', 'Bybit', 'Kraken']
    const blockchainSources = ['Ethereum', 'Uniswap', 'PancakeSwap', 'Curve', 'Solana', 'Bitcoin']
    if (sourceType === 'cex') {
      filtered = filtered.filter(tx => cexSources.includes(tx.source))
    } else if (sourceType === 'blockchain') {
      filtered = filtered.filter(tx => blockchainSources.includes(tx.source))
    } else if (sourceType === 'manual') {
      filtered = filtered.filter(tx => !cexSources.includes(tx.source) && !blockchainSources.includes(tx.source))
    }
  }

  // 按时间范围筛选
  if (timeRange !== 'all') {
    const days = parseInt(timeRange)
    if (!isNaN(days)) {
      const cutoff = new Date(Date.now() - days * DAY).toISOString()
      filtered = filtered.filter(tx => tx.timestamp >= cutoff)
    }
  }

  // 搜索（按代币名称或来源）
  if (search) {
    const q = search.toLowerCase()
    filtered = filtered.filter(tx =>
      tx.symbol.toLowerCase().includes(q) ||
      tx.source.toLowerCase().includes(q)
    )
  }

  // 按时间降序排序
  filtered.sort((a, b) => b.timestamp.localeCompare(a.timestamp))

  // 分页
  const total = filtered.length
  const start = (page - 1) * pageSize
  const transactions = filtered.slice(start, start + pageSize)

  return { transactions, total, page, page_size: pageSize }
}
