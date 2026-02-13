/**
 * Mock 数据管理
 * 在后端 API 未实现前提供模拟数据
 * 数据持久化到 localStorage，支持完整的 CRUD 操作
 */

// ========== 存储键名 ==========
const STORAGE_KEYS = {
  CEX_ACCOUNTS: 'allfi-cex-accounts',
  WALLET_ADDRESSES: 'allfi-wallet-addresses',
  MANUAL_ASSETS: 'allfi-manual-assets',
  ASSET_SNAPSHOTS: 'allfi-asset-snapshots',
  EXCHANGE_RATES: 'allfi-exchange-rates',
  USER_SETTINGS: 'allfi-user-settings'
}

// ========== 初始 Mock 数据 ==========

// CEX 交易所账户
const initialCexAccounts = [
  {
    id: 1,
    exchange: 'Binance',
    name: 'Binance 主账户',
    balance: 35000,
    lastSync: new Date(Date.now() - 3600000).toISOString(),
    status: 'connected',
    apiKeyMasked: '****B3F2',
    holdings: [
      { symbol: 'BTC', name: 'Bitcoin', balance: 0.5, price: 42000, value: 21000, change24h: 2.1 },
      { symbol: 'ETH', name: 'Ethereum', balance: 3.2, price: 2200, value: 7040, change24h: -0.5 },
      { symbol: 'BNB', name: 'Binance Coin', balance: 20, price: 310, value: 6200, change24h: 3.2 },
      { symbol: 'SOL', name: 'Solana', balance: 8, price: 95, value: 760, change24h: 5.8 }
    ]
  },
  {
    id: 2,
    exchange: 'OKX',
    name: 'OKX 合约账户',
    balance: 12000,
    lastSync: new Date(Date.now() - 7200000).toISOString(),
    status: 'connected',
    apiKeyMasked: '****9A1C',
    holdings: [
      { symbol: 'BTC', name: 'Bitcoin', balance: 0.15, price: 42000, value: 6300, change24h: 2.1 },
      { symbol: 'ETH', name: 'Ethereum', balance: 2.0, price: 2200, value: 4400, change24h: -0.5 },
      { symbol: 'USDC', name: 'USD Coin', balance: 1300, price: 1, value: 1300, change24h: 0 }
    ]
  },
  {
    id: 3,
    exchange: 'Coinbase',
    name: 'Coinbase',
    balance: 3000,
    lastSync: new Date(Date.now() - 86400000).toISOString(),
    status: 'error',
    apiKeyMasked: '****X7K2',
    errorMessage: 'API 密钥已过期',
    holdings: [
      { symbol: 'BTC', name: 'Bitcoin', balance: 0.05, price: 42000, value: 2100, change24h: 2.1 },
      { symbol: 'USDC', name: 'USD Coin', balance: 900, price: 1, value: 900, change24h: 0 }
    ]
  }
]

// 区块链钱包
const initialWalletAddresses = [
  {
    id: 1,
    name: '主钱包',
    address: '0x742d35Cc6634C0532925a3b844Bc9e7595f8fC2c',
    blockchain: 'ETH',
    balance: 45000,
    lastSync: new Date(Date.now() - 1800000).toISOString(),
    status: 'active',
    holdings: [
      { symbol: 'ETH', name: 'Ethereum', balance: 15.5, price: 2200, value: 34100, change24h: -0.5 },
      { symbol: 'UNI', name: 'Uniswap', balance: 500, price: 6.5, value: 3250, change24h: 4.2 },
      { symbol: 'LINK', name: 'Chainlink', balance: 300, price: 14, value: 4200, change24h: 1.8 },
      { symbol: 'AAVE', name: 'Aave', balance: 30, price: 115, value: 3450, change24h: -1.2 }
    ]
  },
  {
    id: 2,
    name: 'DeFi 钱包',
    address: '0x1234567890abcdef1234567890abcdef12345678',
    blockchain: 'ETH',
    balance: 8500,
    lastSync: new Date(Date.now() - 3600000).toISOString(),
    status: 'active',
    holdings: [
      { symbol: 'ETH', name: 'Ethereum', balance: 3.0, price: 2200, value: 6600, change24h: -0.5 },
      { symbol: 'USDC', name: 'USD Coin', balance: 1900, price: 1, value: 1900, change24h: 0 }
    ]
  },
  {
    id: 3,
    name: 'BSC 钱包',
    address: '0xabcdef1234567890abcdef1234567890abcdef12',
    blockchain: 'BSC',
    balance: 3200,
    lastSync: new Date(Date.now() - 5400000).toISOString(),
    status: 'active',
    holdings: [
      { symbol: 'BNB', name: 'Binance Coin', balance: 8, price: 310, value: 2480, change24h: 3.2 },
      { symbol: 'CAKE', name: 'PancakeSwap', balance: 288, price: 2.5, value: 720, change24h: 5.5 }
    ]
  },
  {
    id: 4,
    name: 'SOL 主钱包',
    address: 'DYw8jCTfwHNRJhhmFcbXvVDTqWMEVFBX6ZKUmG5CNSKK',
    blockchain: 'SOL',
    balance: 6500,
    lastSync: new Date(Date.now() - 7200000).toISOString(),
    status: 'active',
    holdings: [
      { symbol: 'SOL', name: 'Solana', balance: 55, price: 95, value: 5225, change24h: 5.8 },
      { symbol: 'RAY', name: 'Raydium', balance: 400, price: 3.2, value: 1280, change24h: 8.2 }
    ]
  },
  {
    id: 5,
    name: '冷钱包',
    address: '0xfedcba0987654321fedcba0987654321fedcba09',
    blockchain: 'ETH',
    balance: 500,
    lastSync: new Date(Date.now() - 172800000).toISOString(),
    status: 'inactive',
    holdings: [
      { symbol: 'ETH', name: 'Ethereum', balance: 0.2, price: 2200, value: 440, change24h: -0.5 },
      { symbol: 'USDC', name: 'USD Coin', balance: 60, price: 1, value: 60, change24h: 0 }
    ]
  }
]

// 手动资产（含机构信息）
const initialManualAssets = [
  {
    id: 1,
    name: '港币储蓄',
    type: 'bank',
    institution: 'HSBC 汇丰银行',
    currency: 'HKD',
    balance: 50000,
    note: '',
    lastUpdated: new Date(Date.now() - 604800000).toISOString()
  },
  {
    id: 2,
    name: '美元储蓄',
    type: 'bank',
    institution: 'HSBC 汇丰银行',
    currency: 'USD',
    balance: 10000,
    note: '',
    lastUpdated: new Date(Date.now() - 604800000).toISOString()
  },
  {
    id: 3,
    name: '工资卡',
    type: 'bank',
    institution: '招商银行',
    currency: 'CNY',
    balance: 80000,
    note: '',
    lastUpdated: new Date(Date.now() - 1209600000).toISOString()
  },
  {
    id: 4,
    name: '美股账户',
    type: 'stock',
    institution: '富途证券',
    currency: 'USD',
    balance: 25000,
    note: '',
    lastUpdated: new Date(Date.now() - 2592000000).toISOString()
  },
  {
    id: 5,
    name: '基金组合',
    type: 'fund',
    institution: '天天基金',
    currency: 'CNY',
    balance: 120000,
    note: '',
    lastUpdated: new Date(Date.now() - 864000000).toISOString()
  },
  {
    id: 6,
    name: '家中保险柜',
    type: 'cash',
    institution: '',
    currency: 'CNY',
    balance: 10000,
    note: '应急资金',
    lastUpdated: new Date(Date.now() - 2592000000).toISOString()
  }
]

// 汇率数据
const initialExchangeRates = {
  BTC: 42500,
  ETH: 2280,
  CNY: 0.14,
  lastUpdated: new Date().toISOString()
}

// ========== 数据访问函数 ==========

/**
 * 获取存储数据
 * @param {string} key - 存储键名
 * @param {any} defaultValue - 默认值
 * @returns {any}
 */
function getStorageData(key, defaultValue) {
  try {
    const data = localStorage.getItem(key)
    return data ? JSON.parse(data) : defaultValue
  } catch {
    return defaultValue
  }
}

/**
 * 设置存储数据
 * @param {string} key - 存储键名
 * @param {any} data - 数据
 */
function setStorageData(key, data) {
  localStorage.setItem(key, JSON.stringify(data))
}

/**
 * 初始化 Mock 数据（如果不存在）
 */
export function initMockData() {
  if (!localStorage.getItem(STORAGE_KEYS.CEX_ACCOUNTS)) {
    setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, initialCexAccounts)
  }
  if (!localStorage.getItem(STORAGE_KEYS.WALLET_ADDRESSES)) {
    setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, initialWalletAddresses)
  }
  if (!localStorage.getItem(STORAGE_KEYS.MANUAL_ASSETS)) {
    setStorageData(STORAGE_KEYS.MANUAL_ASSETS, initialManualAssets)
  }
  if (!localStorage.getItem(STORAGE_KEYS.EXCHANGE_RATES)) {
    setStorageData(STORAGE_KEYS.EXCHANGE_RATES, initialExchangeRates)
  }
}

// ========== CEX 账户操作 ==========

export function getCexAccounts() {
  return getStorageData(STORAGE_KEYS.CEX_ACCOUNTS, initialCexAccounts)
}

export function getCexAccountById(id) {
  const accounts = getCexAccounts()
  return accounts.find(a => a.id === id)
}

export function addCexAccount(account) {
  const accounts = getCexAccounts()
  const newAccount = {
    ...account,
    id: Date.now(),
    lastSync: new Date().toISOString(),
    status: 'connected',
    holdings: []
  }
  accounts.push(newAccount)
  setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, accounts)
  return newAccount
}

export function updateCexAccount(id, updates) {
  const accounts = getCexAccounts()
  const index = accounts.findIndex(a => a.id === id)
  if (index !== -1) {
    accounts[index] = { ...accounts[index], ...updates }
    setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, accounts)
    return accounts[index]
  }
  return null
}

export function deleteCexAccount(id) {
  const accounts = getCexAccounts()
  const filtered = accounts.filter(a => a.id !== id)
  setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, filtered)
  return true
}

export function refreshCexAccount(id) {
  const accounts = getCexAccounts()
  const index = accounts.findIndex(a => a.id === id)
  if (index !== -1) {
    // 模拟刷新：更新时间并随机调整余额
    const randomChange = (Math.random() - 0.5) * 0.1
    accounts[index].balance = Math.round(accounts[index].balance * (1 + randomChange))
    accounts[index].lastSync = new Date().toISOString()
    accounts[index].status = 'connected'
    delete accounts[index].errorMessage
    setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, accounts)
    return accounts[index]
  }
  return null
}

// ========== 钱包地址操作 ==========

export function getWalletAddresses() {
  return getStorageData(STORAGE_KEYS.WALLET_ADDRESSES, initialWalletAddresses)
}

export function getWalletAddressById(id) {
  const wallets = getWalletAddresses()
  return wallets.find(w => w.id === id)
}

export function addWalletAddress(wallet) {
  const wallets = getWalletAddresses()
  const newWallet = {
    ...wallet,
    id: Date.now(),
    lastSync: new Date().toISOString(),
    status: 'active',
    balance: 0,
    holdings: []
  }
  wallets.push(newWallet)
  setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, wallets)
  return newWallet
}

export function updateWalletAddress(id, updates) {
  const wallets = getWalletAddresses()
  const index = wallets.findIndex(w => w.id === id)
  if (index !== -1) {
    wallets[index] = { ...wallets[index], ...updates }
    setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, wallets)
    return wallets[index]
  }
  return null
}

export function deleteWalletAddress(id) {
  const wallets = getWalletAddresses()
  const filtered = wallets.filter(w => w.id !== id)
  setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, filtered)
  return true
}

export function refreshWalletAddress(id) {
  const wallets = getWalletAddresses()
  const index = wallets.findIndex(w => w.id === id)
  if (index !== -1) {
    const randomChange = (Math.random() - 0.5) * 0.1
    wallets[index].balance = Math.round(wallets[index].balance * (1 + randomChange))
    wallets[index].lastSync = new Date().toISOString()
    setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, wallets)
    return wallets[index]
  }
  return null
}

// ========== 手动资产操作 ==========

export function getManualAssets() {
  return getStorageData(STORAGE_KEYS.MANUAL_ASSETS, initialManualAssets)
}

export function getManualAssetById(id) {
  const assets = getManualAssets()
  return assets.find(a => a.id === id)
}

export function addManualAsset(asset) {
  const assets = getManualAssets()
  const newAsset = {
    ...asset,
    id: Date.now(),
    lastUpdated: new Date().toISOString()
  }
  assets.push(newAsset)
  setStorageData(STORAGE_KEYS.MANUAL_ASSETS, assets)
  return newAsset
}

export function updateManualAsset(id, updates) {
  const assets = getManualAssets()
  const index = assets.findIndex(a => a.id === id)
  if (index !== -1) {
    assets[index] = { ...assets[index], ...updates, lastUpdated: new Date().toISOString() }
    setStorageData(STORAGE_KEYS.MANUAL_ASSETS, assets)
    return assets[index]
  }
  return null
}

export function deleteManualAsset(id) {
  const assets = getManualAssets()
  const filtered = assets.filter(a => a.id !== id)
  setStorageData(STORAGE_KEYS.MANUAL_ASSETS, filtered)
  return true
}

// ========== 汇率操作 ==========

export function getExchangeRates() {
  return getStorageData(STORAGE_KEYS.EXCHANGE_RATES, initialExchangeRates)
}

export function updateExchangeRates(rates) {
  const current = getExchangeRates()
  const updated = { ...current, ...rates, lastUpdated: new Date().toISOString() }
  setStorageData(STORAGE_KEYS.EXCHANGE_RATES, updated)
  return updated
}

// ========== 资产汇总 ==========

export function getAssetSummary() {
  const cexAccounts = getCexAccounts()
  const walletAddresses = getWalletAddresses()
  const manualAssets = getManualAssets()
  const rates = getExchangeRates()

  // 计算 CEX 总值
  const cexTotal = cexAccounts.reduce((sum, acc) => sum + acc.balance, 0)

  // 计算链上总值
  const blockchainTotal = walletAddresses.reduce((sum, w) => sum + w.balance, 0)

  // 计算手动资产总值（转换为 USD）
  const manualTotal = manualAssets.reduce((sum, a) => {
    if (a.currency === 'CNY') {
      return sum + a.balance * rates.CNY
    }
    return sum + a.balance
  }, 0)

  const totalValue = cexTotal + blockchainTotal + manualTotal

  // 计算24小时变化（模拟）
  const change24h = (Math.random() - 0.4) * 5 // -2% 到 +3%
  const changeValue = totalValue * (change24h / 100)

  return {
    totalValue,
    change24h,
    changeValue,
    categories: {
      cex: {
        value: cexTotal,
        count: cexAccounts.length,
        percentage: (cexTotal / totalValue) * 100
      },
      blockchain: {
        value: blockchainTotal,
        count: walletAddresses.length,
        percentage: (blockchainTotal / totalValue) * 100
      },
      manual: {
        value: manualTotal,
        count: manualAssets.length,
        percentage: (manualTotal / totalValue) * 100
      }
    },
    cexAccounts,
    walletAddresses,
    manualAssets,
    exchangeRates: rates
  }
}

// ========== 历史数据生成 ==========

export function generateHistoryData(timeRange = '30D') {
  // LAST_YEAR 和 1Y 统一按 365 天处理（日粒度）
  const points = timeRange === '7D' ? 7 :
                 timeRange === '30D' ? 30 :
                 timeRange === '90D' ? 90 :
                 (timeRange === 'LAST_YEAR' || timeRange === '1Y') ? 365 : 730

  const summary = getAssetSummary()
  const baseValue = summary.totalValue
  const labels = []
  const values = []

  // 使用固定种子让同一天数据保持一致
  let seed = baseValue
  const seededRandom = () => {
    seed = (seed * 9301 + 49297) % 233280
    return seed / 233280
  }

  for (let i = points; i >= 0; i--) {
    const date = new Date()
    date.setDate(date.getDate() - i)
    // 统一使用 YYYY-MM-DD 格式，便于日历热力图映射
    const y = date.getFullYear()
    const m = String(date.getMonth() + 1).padStart(2, '0')
    const d = String(date.getDate()).padStart(2, '0')
    labels.push(`${y}-${m}-${d}`)
    // 模拟历史数据：基于当前值的随机波动
    values.push(baseValue * (0.85 + seededRandom() * 0.3))
  }

  return { labels, values }
}

// ========== 用户设置 ==========

export function getUserSettings() {
  return getStorageData(STORAGE_KEYS.USER_SETTINGS, {
    language: 'zh-CN',
    currency: 'USDC',
    autoRefresh: true,
    refreshInterval: 300,
    historyRetention: 90,
    priceAlerts: true,
    syncNotifications: true,
    emailDigest: false,
    showBalances: true,
    confirmOperations: true
  })
}

export function updateUserSettings(settings) {
  const current = getUserSettings()
  const updated = { ...current, ...settings }
  setStorageData(STORAGE_KEYS.USER_SETTINGS, updated)
  return updated
}

// ========== 数据导出 ==========

export function exportAllData() {
  return {
    cexAccounts: getCexAccounts(),
    walletAddresses: getWalletAddresses(),
    manualAssets: getManualAssets(),
    exchangeRates: getExchangeRates(),
    userSettings: getUserSettings(),
    exportedAt: new Date().toISOString()
  }
}

// ========== 清除缓存 ==========

export function clearAllData() {
  Object.values(STORAGE_KEYS).forEach(key => {
    localStorage.removeItem(key)
  })
}

// ========== 重置数据 ==========

export function resetToDefaults() {
  setStorageData(STORAGE_KEYS.CEX_ACCOUNTS, initialCexAccounts)
  setStorageData(STORAGE_KEYS.WALLET_ADDRESSES, initialWalletAddresses)
  setStorageData(STORAGE_KEYS.MANUAL_ASSETS, initialManualAssets)
  setStorageData(STORAGE_KEYS.EXCHANGE_RATES, initialExchangeRates)
  localStorage.removeItem(STORAGE_KEYS.USER_SETTINGS)
}

// 初始化
initMockData()
