/**
 * API 服务层入口
 * 统一导出所有 API 服务，自动切换真实 API 和 Mock 数据
 */

import { get, post, put, del, healthCheck, ApiError } from './client.js'
import * as mockData from './mockData.js'

// 是否使用 Mock 数据（后端未就绪时使用）
// 可通过环境变量控制：VITE_USE_MOCK_API=true
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'

// 模拟网络延迟
const simulateDelay = (ms = 500) => new Promise(resolve => setTimeout(resolve, ms))

// ========== 认证服务（PIN 模式） ==========

export const authService = {
  /**
   * 获取认证状态（是否已设置 PIN）
   * @returns {Promise<{pin_set: boolean}>}
   */
  async getStatus() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return { pin_set: false }
    }
    return get('/auth/status')
  },

  /**
   * 首次设置 PIN
   * @param {string} pin - 4-8 位数字 PIN
   * @returns {Promise<{token: string}>}
   */
  async setupPIN(pin) {
    if (USE_MOCK) {
      await simulateDelay(600)
      return { token: 'mock-jwt-token-' + Date.now() }
    }
    return post('/auth/setup', { pin })
  },

  /**
   * PIN 登录
   * @param {string} pin - PIN 码
   * @returns {Promise<{token: string}>}
   */
  async login(pin) {
    if (USE_MOCK) {
      await simulateDelay(600)
      return { token: 'mock-jwt-token-' + Date.now() }
    }
    return post('/auth/login', { pin })
  },

  /**
   * 修改 PIN
   * @param {string} oldPin - 旧 PIN
   * @param {string} newPin - 新 PIN
   * @returns {Promise<void>}
   */
  async changePIN(oldPin, newPin) {
    if (USE_MOCK) {
      await simulateDelay(600)
      return {}
    }
    return post('/auth/change', { current_pin: oldPin, new_pin: newPin })
  }
}

// ========== 资产服务 ==========

export const assetService = {
  /**
   * 获取资产总览
   * 真实 API 模式下并发请求多个接口，聚合为 assetStore 期望的格式
   * @returns {Promise<object>}
   */
  async getSummary() {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.getAssetSummary()
    }
    // 并发请求多个接口，聚合为 Store 期望的格式
    const [summaryData, accountsData, walletsData, manualsData, ratesData] = await Promise.all([
      get('/assets/summary'),
      get('/exchanges/accounts').catch(() => ({ accounts: [] })),
      get('/wallets/addresses').catch(() => ({ addresses: [] })),
      get('/manual/assets').catch(() => ({ assets: [] })),
      get('/rates/current').catch(() => ({ rates: {} })),
    ])

    // 适配 CEX 账户列表（后端可能返回数组或 {accounts: []}）
    const rawAccounts = Array.isArray(accountsData) ? accountsData : (accountsData.accounts || [])
    const cexAccounts = rawAccounts.map(a => ({
      id: a.id,
      exchange: a.exchange_name,
      name: a.label || a.exchange_name,
      balance: 0,
      lastSync: a.updated_at,
      status: a.status || 'connected',
      apiKeyMasked: '',
      holdings: []
    }))

    // 适配钱包地址列表（后端可能返回数组或 {addresses: []}）
    const rawWallets = Array.isArray(walletsData) ? walletsData : (walletsData.addresses || [])
    const walletAddresses = rawWallets.map(w => ({
      id: w.id,
      name: w.label || w.address?.slice(0, 10),
      address: w.address,
      blockchain: w.blockchain,
      balance: 0,
      lastSync: w.updated_at,
      status: 'active',
      holdings: []
    }))

    // 适配手动资产列表（后端可能返回数组或 {assets: []}）
    const rawManuals = Array.isArray(manualsData) ? manualsData : (manualsData.assets || [])
    const manualAssets = rawManuals.map(a => ({
      id: a.id,
      name: a.asset_name,
      type: a.asset_type,
      institution: a.institution || '',
      currency: a.currency,
      balance: a.amount,
      note: a.notes || '',
      lastUpdated: a.updated_at
    }))

    // 适配汇率数据（后端可能返回扁平对象或 {rates: {}}）
    const rates = ratesData.rates || ratesData || {}

    // 从后端 summary 数据中提取各分类总值
    const bySource = summaryData.by_source || {}
    const cexTotal = bySource.cex || 0
    const blockchainTotal = bySource.blockchain || 0
    const manualTotal = bySource.manual || 0
    const totalValue = summaryData.total_value || (cexTotal + blockchainTotal + manualTotal)

    return {
      totalValue,
      change24h: 0,
      changeValue: 0,
      categories: {
        cex: { value: cexTotal, count: cexAccounts.length, percentage: totalValue ? (cexTotal / totalValue) * 100 : 0 },
        blockchain: { value: blockchainTotal, count: walletAddresses.length, percentage: totalValue ? (blockchainTotal / totalValue) * 100 : 0 },
        manual: { value: manualTotal, count: manualAssets.length, percentage: totalValue ? (manualTotal / totalValue) * 100 : 0 }
      },
      cexAccounts,
      walletAddresses,
      manualAssets,
      exchangeRates: { ...rates, lastUpdated: ratesData.last_updated }
    }
  },

  /**
   * 获取资产历史
   * @param {string} timeRange - 时间范围 (7D, 30D, 90D, 1Y)
   * @returns {Promise<{labels: string[], values: number[]}>}
   */
  async getHistory(timeRange = '30D') {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockData.generateHistoryData(timeRange)
    }
    // 后端接受 days 整数参数
    const daysMap = { '7D': 7, '30D': 30, '90D': 90, '1Y': 365 }
    const days = daysMap[timeRange] || 30
    const result = await get('/assets/history', { days })
    // 后端返回 { snapshots: [{date, total_value, currency}] }
    // 转换为前端使用的 { labels, values } 格式
    if (result.snapshots) {
      return {
        labels: result.snapshots.map(s => s.date),
        values: result.snapshots.map(s => s.total_value)
      }
    }
    return result
  },

  /**
   * 刷新所有资产数据
   * 真实 API 模式下先触发后端刷新，再重新聚合获取最新数据
   * @returns {Promise<object>}
   */
  async refresh() {
    if (USE_MOCK) {
      await simulateDelay(2000)
      // 刷新所有账户
      const cexAccounts = mockData.getCexAccounts()
      cexAccounts.forEach(acc => mockData.refreshCexAccount(acc.id))

      const wallets = mockData.getWalletAddresses()
      wallets.forEach(w => mockData.refreshWalletAddress(w.id))

      return mockData.getAssetSummary()
    }
    // 先触发后端刷新，再重新获取聚合数据
    await post('/assets/refresh')
    return this.getSummary()
  },

  /**
   * 获取资产明细列表
   * @param {string} sourceType - 来源类型筛选（cex/blockchain/manual），为空则返回全部
   * @param {string} currency - 计价货币，默认 USD
   * @returns {Promise<object[]>} 资产明细数组
   */
  async getDetails(sourceType = '', currency = 'USD') {
    if (USE_MOCK) {
      await simulateDelay(300)
      // 从 Mock 数据聚合资产明细
      const details = []
      if (!sourceType || sourceType === 'cex') {
        for (const acc of mockData.getCexAccounts()) {
          if (acc.holdings) {
            for (const h of acc.holdings) {
              details.push({ id: h.id || details.length + 1, symbol: h.symbol, amount: h.balance, value: h.value, price: h.price, source: acc.name, source_type: 'cex', updated_at: acc.lastSync })
            }
          }
        }
      }
      if (!sourceType || sourceType === 'blockchain') {
        for (const w of mockData.getWalletAddresses()) {
          if (w.holdings) {
            for (const h of w.holdings) {
              details.push({ id: h.id || details.length + 1, symbol: h.symbol, amount: h.balance, value: h.value, price: h.price, source: w.name, source_type: 'blockchain', updated_at: w.lastSync })
            }
          }
        }
      }
      if (!sourceType || sourceType === 'manual') {
        for (const a of mockData.getManualAssets()) {
          details.push({ id: a.id, symbol: a.currency, amount: a.balance, value: a.balance, price: 1, source: a.name, source_type: 'manual', updated_at: a.lastUpdated })
        }
      }
      return details
    }
    const params = { currency }
    if (sourceType) params.source_type = sourceType
    const result = await get('/assets/details', params)
    return result.assets || result || []
  }
}

// ========== CEX 账户服务 ==========

export const cexService = {
  /**
   * 获取支持的交易所列表
   * @returns {Promise<Array<{id: string, name: string, category: string}>>}
   */
  async getSupportedExchanges() {
    if (USE_MOCK) {
      await simulateDelay(200)
      // Mock 数据，返回常用交易所
      return [
        { id: 'binance', name: 'Binance', category: 'spot' },
        { id: 'okx', name: 'OKX', category: 'spot' },
        { id: 'coinbase', name: 'Coinbase', category: 'spot' },
        { id: 'kraken', name: 'Kraken', category: 'spot' },
        { id: 'kucoin', name: 'KuCoin', category: 'spot' },
        { id: 'gate', name: 'Gate.io', category: 'spot' },
        { id: 'bitget', name: 'Bitget', category: 'spot' },
        { id: 'bybit', name: 'Bybit', category: 'spot' },
        { id: 'mexc', name: 'MEXC', category: 'spot' },
        { id: 'huobi', name: 'Huobi', category: 'spot' },
        { id: 'bingx', name: 'BingX', category: 'spot' },
        { id: 'bitfinex', name: 'Bitfinex', category: 'spot' },
        { id: 'coinex', name: 'CoinEx', category: 'spot' },
        { id: 'binanceusdm', name: 'Binance USD-M Futures', category: 'futures' },
        { id: 'binancecoinm', name: 'Binance COIN-M Futures', category: 'futures' },
        { id: 'bitmex', name: 'BitMEX', category: 'futures' },
        { id: 'dydx', name: 'dYdX', category: 'derivatives' },
        { id: 'hyperliquid', name: 'Hyperliquid', category: 'derivatives' },
      ]
    }
    const result = await get('/exchanges/supported')
    return result.exchanges || []
  },

  /**
   * 获取所有 CEX 账户
   * @returns {Promise<object[]>}
   */
  async getAccounts() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockData.getCexAccounts()
    }
    const result = await get('/exchanges/accounts')
    // 灵活处理数组或包装对象
    const accounts = Array.isArray(result) ? result : (result.accounts || [])
    return accounts.map(a => ({
      id: a.id,
      exchange: a.exchange_name,
      name: a.label || a.exchange_name,
      balance: 0,
      lastSync: a.updated_at,
      status: a.status || 'connected',
      apiKeyMasked: '',
      note: a.note || '',
      holdings: []
    }))
  },

  /**
   * 获取单个账户
   * @param {number} id - 账户ID
   * @returns {Promise<object>}
   */
  async getAccount(id) {
    if (USE_MOCK) {
      await simulateDelay(100)
      const account = mockData.getCexAccountById(id)
      if (!account) throw new ApiError(40401, '账户不存在')
      return account
    }
    const result = await get(`/exchanges/accounts/${id}`)
    const a = result.account || result
    return {
      id: a.id,
      exchange: a.exchange_name,
      name: a.label || a.exchange_name,
      balance: 0,
      lastSync: a.updated_at,
      status: a.status || 'connected',
      apiKeyMasked: '',
      note: a.note || '',
      holdings: []
    }
  },

  /**
   * 添加 CEX 账户
   * @param {object} data - 账户数据
   * @returns {Promise<object>}
   */
  async addAccount(data) {
    if (USE_MOCK) {
      await simulateDelay(500)
      return mockData.addCexAccount(data)
    }
    const result = await post('/exchanges/accounts', data)
    const a = result.account || result
    return {
      id: a.id,
      exchange: a.exchange_name,
      name: a.label || a.exchange_name,
      balance: 0,
      lastSync: a.updated_at,
      status: a.status || 'connected',
      apiKeyMasked: '',
      note: a.note || '',
      holdings: []
    }
  },

  /**
   * 更新 CEX 账户
   * @param {number} id - 账户ID
   * @param {object} data - 更新数据
   * @returns {Promise<object>}
   */
  async updateAccount(id, data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      const account = mockData.updateCexAccount(id, data)
      if (!account) throw new ApiError(40401, '账户不存在')
      return account
    }
    const result = await put(`/exchanges/accounts/${id}`, data)
    const a = result.account || result
    return {
      id: a.id,
      exchange: a.exchange_name,
      name: a.label || a.exchange_name,
      balance: 0,
      lastSync: a.updated_at,
      status: a.status || 'connected',
      apiKeyMasked: '',
      note: a.note || '',
      holdings: []
    }
  },

  /**
   * 删除 CEX 账户
   * @param {number} id - 账户ID
   * @returns {Promise<boolean>}
   */
  async deleteAccount(id) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.deleteCexAccount(id)
    }
    return del(`/exchanges/accounts/${id}`)
  },

  /**
   * 测试交易所连接
   * @param {number} id - 账户ID
   * @returns {Promise<{success: boolean, message: string}>}
   */
  async testConnection(id) {
    if (USE_MOCK) {
      await simulateDelay(1000)
      return { success: true, message: '连接成功' }
    }
    return post(`/exchanges/accounts/${id}/test`)
  },

  /**
   * 刷新 CEX 账户余额
   * @param {number} id - 账户ID
   * @returns {Promise<object>}
   */
  async refreshAccount(id) {
    if (USE_MOCK) {
      await simulateDelay(1500)
      const account = mockData.refreshCexAccount(id)
      if (!account) throw new ApiError(40401, '账户不存在')
      return account
    }
    await post(`/exchanges/accounts/${id}/sync`)
    // sync 返回 {message}，需要再获取余额和账户信息
    const [account, balances] = await Promise.all([
      this.getAccount(id),
      get(`/exchanges/accounts/${id}/balances`).catch(() => null)
    ])
    if (balances) {
      account.balance = balances.total_value || 0
      account.holdings = (balances.balances || []).map(b => ({
        symbol: b.symbol,
        name: b.symbol,
        balance: b.total,
        price: b.total > 0 ? (b.value_usd / b.total) : 0,
        value: b.value_usd,
        change24h: 0
      }))
    }
    return account
  }
}

// ========== 钱包地址服务 ==========

export const walletService = {
  /**
   * 获取所有钱包地址
   * @returns {Promise<object[]>}
   */
  async getAddresses() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockData.getWalletAddresses()
    }
    const result = await get('/wallets/addresses')
    const addresses = Array.isArray(result) ? result : (result.addresses || [])
    return addresses.map(w => ({
      id: w.id,
      name: w.label || `${w.address?.slice(0, 6)}...${w.address?.slice(-4)}`,
      address: w.address,
      blockchain: w.blockchain,
      balance: 0,
      lastSync: w.updated_at,
      status: 'active',
      holdings: []
    }))
  },

  /**
   * 获取单个钱包
   * @param {number} id - 钱包ID
   * @returns {Promise<object>}
   */
  async getAddress(id) {
    if (USE_MOCK) {
      await simulateDelay(100)
      const wallet = mockData.getWalletAddressById(id)
      if (!wallet) throw new ApiError(40401, '钱包不存在')
      return wallet
    }
    const result = await get(`/wallets/addresses/${id}`)
    const w = result.address || result
    return {
      id: w.id,
      name: w.label || `${w.address?.slice(0, 6)}...${w.address?.slice(-4)}`,
      address: w.address,
      blockchain: w.blockchain,
      balance: 0,
      lastSync: w.updated_at,
      status: 'active',
      holdings: []
    }
  },

  /**
   * 添加钱包地址
   * @param {object} data - 钱包数据
   * @returns {Promise<object>}
   */
  async addAddress(data) {
    if (USE_MOCK) {
      await simulateDelay(500)
      return mockData.addWalletAddress(data)
    }
    const result = await post('/wallets/addresses', data)
    const w = result.address || result
    return {
      id: w.id,
      name: w.label || `${w.address?.slice(0, 6)}...${w.address?.slice(-4)}`,
      address: w.address,
      blockchain: w.blockchain,
      balance: 0,
      lastSync: w.updated_at,
      status: 'active',
      holdings: []
    }
  },

  /**
   * 更新钱包地址
   * @param {number} id - 钱包ID
   * @param {object} data - 更新数据
   * @returns {Promise<object>}
   */
  async updateAddress(id, data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      const wallet = mockData.updateWalletAddress(id, data)
      if (!wallet) throw new ApiError(40401, '钱包不存在')
      return wallet
    }
    const result = await put(`/wallets/addresses/${id}`, data)
    const w = result.address || result
    return {
      id: w.id,
      name: w.label || `${w.address?.slice(0, 6)}...${w.address?.slice(-4)}`,
      address: w.address,
      blockchain: w.blockchain,
      balance: 0,
      lastSync: w.updated_at,
      status: 'active',
      holdings: []
    }
  },

  /**
   * 删除钱包地址
   * @param {number} id - 钱包ID
   * @returns {Promise<boolean>}
   */
  async deleteAddress(id) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.deleteWalletAddress(id)
    }
    return del(`/wallets/addresses/${id}`)
  },

  /**
   * 刷新钱包余额
   * @param {number} id - 钱包ID
   * @returns {Promise<object>}
   */
  async refreshAddress(id) {
    if (USE_MOCK) {
      await simulateDelay(1500)
      const wallet = mockData.refreshWalletAddress(id)
      if (!wallet) throw new ApiError(40401, '钱包不存在')
      return wallet
    }
    await post(`/wallets/addresses/${id}/sync`)
    const [wallet, balances] = await Promise.all([
      this.getAddress(id),
      get(`/wallets/addresses/${id}/balances`).catch(() => null)
    ])
    if (balances) {
      wallet.balance = balances.total_value_usd || 0
      wallet.holdings = []
      if (balances.native_balance > 0) {
        wallet.holdings.push({
          symbol: wallet.blockchain === 'ethereum' ? 'ETH' : wallet.blockchain.toUpperCase(),
          name: wallet.blockchain,
          balance: balances.native_balance,
          value: 0,
          change24h: 0
        })
      }
      if (balances.token_balances) {
        for (const [symbol, amount] of Object.entries(balances.token_balances)) {
          wallet.holdings.push({ symbol, name: symbol, balance: amount, value: 0, change24h: 0 })
        }
      }
    }
    return wallet
  },

  /**
   * 批量导入钱包地址
   * @param {string[]} addresses - 地址列表
   * @param {string} blockchain - 区块链类型
   * @returns {Promise<{imported: number, failed: number}>}
   */
  async batchImport(addresses, blockchain) {
    if (USE_MOCK) {
      await simulateDelay(1500)
      // 模拟批量导入结果
      let imported = 0
      let failedCount = 0
      const seen = new Set()
      for (const addr of addresses) {
        const lower = addr.toLowerCase().trim()
        if (!lower) continue
        if (seen.has(lower) || !/^0x[0-9a-fA-F]{40}$/.test(lower)) {
          failedCount++
          continue
        }
        seen.add(lower)
        imported++
      }
      return { imported, failed: failedCount }
    }
    // 后端期望 { addresses: [{blockchain, address, label}] }
    const batchAddresses = addresses
      .map(addr => addr.trim())
      .filter(Boolean)
      .map(addr => ({
        blockchain,
        address: addr,
        label: `${addr.slice(0, 6)}...${addr.slice(-4)}`
      }))
    return post('/wallets/batch', { addresses: batchAddresses })
  }
}

// ========== 手动资产服务 ==========

export const manualAssetService = {
  /**
   * 获取所有手动资产
   * @returns {Promise<object[]>}
   */
  async getAssets() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return mockData.getManualAssets()
    }
    const result = await get('/manual/assets')
    const assets = Array.isArray(result) ? result : (result.assets || [])
    return assets.map(a => ({
      id: a.id,
      name: a.asset_name || a.name,
      type: a.asset_type || a.type,
      institution: a.institution || '',
      currency: a.currency,
      balance: a.amount ?? a.balance,
      note: a.notes || a.note || '',
      lastUpdated: a.updated_at
    }))
  },

  /**
   * 获取单个资产
   * @param {number} id - 资产ID
   * @returns {Promise<object>}
   */
  async getAsset(id) {
    if (USE_MOCK) {
      await simulateDelay(100)
      const asset = mockData.getManualAssetById(id)
      if (!asset) throw new ApiError(40401, '资产不存在')
      return asset
    }
    const result = await get(`/manual/assets/${id}`)
    const a = result.asset || result
    return {
      id: a.id,
      name: a.asset_name || a.name,
      type: a.asset_type || a.type,
      institution: a.institution || '',
      currency: a.currency,
      balance: a.amount ?? a.balance,
      note: a.notes || a.note || '',
      lastUpdated: a.updated_at
    }
  },

  /**
   * 添加手动资产
   * @param {object} data - 资产数据
   * @returns {Promise<object>}
   */
  async addAsset(data) {
    if (USE_MOCK) {
      await simulateDelay(500)
      return mockData.addManualAsset(data)
    }
    // 将前端字段名映射为后端字段名
    const payload = {
      asset_type: data.type || data.asset_type,
      asset_name: data.name || data.asset_name,
      amount: data.balance ?? data.amount,
      currency: data.currency,
      notes: data.note || data.notes || '',
      institution: data.institution || ''
    }
    const result = await post('/manual/assets', payload)
    const a = result.asset || result
    return {
      id: a.id,
      name: a.asset_name || a.name,
      type: a.asset_type || a.type,
      institution: a.institution || '',
      currency: a.currency,
      balance: a.amount ?? a.balance,
      note: a.notes || a.note || '',
      lastUpdated: a.updated_at
    }
  },

  /**
   * 更新手动资产
   * @param {number} id - 资产ID
   * @param {object} data - 更新数据
   * @returns {Promise<object>}
   */
  async updateAsset(id, data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      const asset = mockData.updateManualAsset(id, data)
      if (!asset) throw new ApiError(40401, '资产不存在')
      return asset
    }
    const payload = {
      asset_type: data.type || data.asset_type,
      asset_name: data.name || data.asset_name,
      amount: data.balance ?? data.amount,
      currency: data.currency,
      notes: data.note || data.notes || '',
      institution: data.institution || ''
    }
    const result = await put(`/manual/assets/${id}`, payload)
    const a = result.asset || result
    return {
      id: a.id,
      name: a.asset_name || a.name,
      type: a.asset_type || a.type,
      institution: a.institution || '',
      currency: a.currency,
      balance: a.amount ?? a.balance,
      note: a.notes || a.note || '',
      lastUpdated: a.updated_at
    }
  },

  /**
   * 删除手动资产
   * @param {number} id - 资产ID
   * @returns {Promise<boolean>}
   */
  async deleteAsset(id) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.deleteManualAsset(id)
    }
    return del(`/manual/assets/${id}`)
  }
}

// ========== 汇率服务 ==========

export const rateService = {
  /**
   * 获取当前汇率
   * @returns {Promise<object>}
   */
  async getCurrentRates() {
    if (USE_MOCK) {
      await simulateDelay(100)
      return mockData.getExchangeRates()
    }
    const result = await get('/rates/current')
    // 后端返回 {rates: {BTC: 64000, ...}, base, last_updated, source, is_cached}
    // Store 期望扁平格式 {BTC: 64000, ETH: 2280, lastUpdated: "..."}
    // 灵活处理：后端可能直接返回 rates 对象或包装对象
    const rates = result.rates || result
    return {
      ...(typeof rates === 'object' ? rates : {}),
      lastUpdated: result.last_updated ? new Date(result.last_updated).toISOString() : new Date().toISOString()
    }
  },

  /**
   * 获取历史汇率
   * @param {string} from - 源货币
   * @param {string} to - 目标货币
   * @param {string} range - 时间范围
   * @returns {Promise<object>}
   */
  async getHistoryRates(from, to, range = '30D') {
    if (USE_MOCK) {
      await simulateDelay(200)
      // 生成模拟历史汇率数据
      const baseRate = mockData.getExchangeRates()[from] || 1
      const points = range === '7D' ? 7 : range === '30D' ? 30 : 90
      const labels = []
      const values = []

      for (let i = points; i >= 0; i--) {
        const date = new Date()
        date.setDate(date.getDate() - i)
        labels.push(date.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' }))
        values.push(baseRate * (0.9 + Math.random() * 0.2))
      }

      return { from, to, labels, values }
    }
    // 后端接受 base, quote, days 参数
    const daysMap = { '7D': 7, '30D': 30, '90D': 90 }
    const result = await get('/rates/history', { base: from, quote: to, days: daysMap[range] || 30 })
    // 后端返回 {history: [{date, rate}], base, quote, days}
    // 前端期望 {from, to, labels, values}
    const history = result.history || []
    return {
      from,
      to,
      labels: history.map(h => h.date),
      values: history.map(h => h.rate)
    }
  },

  /**
   * 获取币种实时价格
   * @param {string} symbols - 币种列表（逗号分隔），如 'BTC,ETH,USDC'
   * @returns {Promise<object[]>} 价格列表 [{symbol, price_usd, change_24h, last_updated}]
   */
  async getPrices(symbols = 'BTC,ETH') {
    if (USE_MOCK) {
      await simulateDelay(200)
      const mockPrices = {
        BTC: { price_usd: 96000 + (Math.random() - 0.5) * 2000, change_24h: (Math.random() - 0.4) * 5 },
        ETH: { price_usd: 2300 + (Math.random() - 0.5) * 200, change_24h: (Math.random() - 0.4) * 6 },
        USDC: { price_usd: 1.0, change_24h: 0 },
        USDT: { price_usd: 1.0, change_24h: 0 },
        BNB: { price_usd: 420 + (Math.random() - 0.5) * 40, change_24h: (Math.random() - 0.4) * 4 },
        SOL: { price_usd: 180 + (Math.random() - 0.5) * 30, change_24h: (Math.random() - 0.4) * 7 },
      }
      const symbolList = symbols.split(',').map(s => s.trim().toUpperCase())
      return symbolList.map(s => ({
        symbol: s,
        price_usd: mockPrices[s]?.price_usd || 0,
        change_24h: mockPrices[s]?.change_24h || 0,
        last_updated: Date.now()
      }))
    }
    const result = await get('/rates/prices', { symbols })
    return result.prices || result || []
  },

  /**
   * 强制刷新汇率缓存
   * @returns {Promise<{message: string}>}
   */
  async refreshRates() {
    if (USE_MOCK) {
      await simulateDelay(500)
      return { message: '汇率缓存刷新成功' }
    }
    return post('/rates/refresh')
  }
}

// ========== 设置服务 ==========

export const settingsService = {
  /**
   * 获取用户设置
   * @returns {Promise<object>}
   */
  async getSettings() {
    if (USE_MOCK) {
      await simulateDelay(100)
      return mockData.getUserSettings()
    }
    const result = await get('/users/settings')
    // 后端返回 {settings: {key: value}} 或直接返回设置对象
    return result.settings || result
  },

  /**
   * 更新用户设置
   * @param {object} settings - 设置键值对
   * @returns {Promise<{message: string}>}
   */
  async updateSettings(settings) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return mockData.updateUserSettings(settings)
    }
    // 后端期望 { settings: map[string]string } 包装
    return put('/users/settings', { settings })
  },

  /**
   * 导出用户数据
   * @returns {Promise<object>}
   */
  async exportData() {
    if (USE_MOCK) {
      await simulateDelay(500)
      return mockData.exportAllData()
    }
    return get('/users/export')
  },

  /**
   * 清除缓存
   * @returns {Promise<boolean>}
   */
  async clearCache() {
    if (USE_MOCK) {
      await simulateDelay(300)
      // 只清除临时数据，保留用户设置
      return true
    }
    return post('/users/clear-cache')
  },

  /**
   * 重置所有设置
   * @returns {Promise<boolean>}
   */
  async resetSettings() {
    if (USE_MOCK) {
      await simulateDelay(500)
      mockData.resetToDefaults()
      return true
    }
    return post('/users/reset-settings')
  },

  /**
   * 获取 API Key 配置列表
   * @returns {Promise<{keys: Array}>}
   */
  async getAPIKeys() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return {
        keys: [
          { provider: 'etherscan', display_name: 'Etherscan', configured: false, masked_key: '', description: '以太坊区块链浏览器 API' },
          { provider: 'bscscan', display_name: 'BscScan', configured: false, masked_key: '', description: 'BSC 区块链浏览器 API' },
          { provider: 'coingecko', display_name: 'CoinGecko', configured: false, masked_key: '', description: '加密货币价格 API（可选）' }
        ]
      }
    }
    return get('/system/apikeys')
  },

  /**
   * 更新 API Key
   * @param {string} provider - 服务商标识
   * @param {string} apiKey - API Key
   * @returns {Promise<{success: boolean}>}
   */
  async updateAPIKey(provider, apiKey) {
    if (USE_MOCK) {
      await simulateDelay(500)
      return { success: true }
    }
    return put('/system/apikeys', { provider, api_key: apiKey })
  },

  /**
   * 删除 API Key
   * @param {string} provider - 服务商标识
   * @returns {Promise<{success: boolean}>}
   */
  async deleteAPIKey(provider) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return { success: true }
    }
    return del(`/system/apikeys?provider=${provider}`)
  }
}

// ========== 通知服务 ==========

export const notificationService = {
  /**
   * 获取通知列表
   * @param {number} page - 页码
   * @param {number} pageSize - 每页数量
   * @returns {Promise<{list: object[], pagination: object}>}
   */
  async getNotifications(page = 1, pageSize = 20) {
    if (USE_MOCK) {
      await simulateDelay(200)
      // 生成模拟通知数据
      const types = ['daily_digest', 'price_alert', 'asset_change']
      const titles = ['每日资产摘要', 'BTC 价格突破 $50,000', '资产总值变动超过 5%', 'ETH 24h 涨幅 8.2%']
      const contents = [
        '总资产: $125,430.00 | CEX: $80,200 | 链上: $35,230 | 传统: $10,000',
        'BTC 当前价格 $50,123，较昨日上涨 3.2%',
        '您的资产总值在过去 24 小时内变动了 5.3%',
        'ETH 当前价格 $2,850，24h 涨幅 8.2%'
      ]
      const list = Array.from({ length: Math.min(pageSize, 10) }, (_, i) => ({
        id: (page - 1) * pageSize + i + 1,
        type: types[i % types.length],
        title: titles[i % titles.length],
        message: contents[i % contents.length],
        is_read: i >= 3,
        created_at: new Date(Date.now() - i * 3600000).toISOString().replace('T', ' ').slice(0, 19)
      }))
      return {
        list,
        pagination: { page, page_size: pageSize, total: 25, total_pages: 2, has_next: page < 2, has_prev: page > 1 }
      }
    }
    // 使用 params 对象传递 GET 参数，而非拼接查询字符串
    return get('/notifications', { page, page_size: pageSize })
  },

  /**
   * 获取未读通知数量
   * @returns {Promise<{count: number}>}
   */
  async getUnreadCount() {
    if (USE_MOCK) {
      await simulateDelay(100)
      return { count: 3 }
    }
    return get('/notifications/unread-count')
  },

  /**
   * 标记通知为已读
   * @param {number} id - 通知 ID
   * @returns {Promise<void>}
   */
  async markRead(id) {
    if (USE_MOCK) {
      await simulateDelay(100)
      return true
    }
    return post(`/notifications/${id}/read`)
  },

  /**
   * 标记所有通知为已读
   * @returns {Promise<void>}
   */
  async markAllRead() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return true
    }
    return post('/notifications/read-all')
  },

  /**
   * 获取通知偏好
   * @returns {Promise<{preferences: {email_enabled, push_enabled, price_alert, portfolio_alert, system_notice}}>}
   */
  async getPreferences() {
    if (USE_MOCK) {
      await simulateDelay(100)
      return {
        preferences: {
          email_enabled: false,
          push_enabled: true,
          price_alert: true,
          portfolio_alert: true,
          system_notice: true
        }
      }
    }
    return get('/notifications/preferences')
  },

  /**
   * 更新通知偏好
   * @param {object} data - 偏好数据 { email_enabled, push_enabled, price_alert, portfolio_alert, system_notice }
   * @returns {Promise<{preferences: object}>}
   */
  async updatePreferences(data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return { preferences: { ...data } }
    }
    return put('/notifications/preferences', data)
  },

  // ===== WebPush 推送 =====

  /**
   * 获取 VAPID 公钥
   * @returns {Promise<{vapid_public_key: string}>}
   */
  async getVapidKey() {
    if (USE_MOCK) {
      await simulateDelay(100)
      return { vapid_public_key: '' }
    }
    return get('/notifications/webpush/vapid')
  },

  /**
   * 订阅 WebPush 推送
   * @param {PushSubscription} subscription - 浏览器 PushSubscription 对象
   * @returns {Promise<void>}
   */
  async subscribePush(subscription) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return {}
    }
    const sub = subscription.toJSON()
    return post('/notifications/webpush/subscribe', {
      endpoint: sub.endpoint,
      keys: {
        p256dh: sub.keys.p256dh,
        auth: sub.keys.auth
      }
    })
  },

  /**
   * 取消订阅 WebPush 推送
   * @param {string} endpoint - 推送端点 URL
   * @returns {Promise<void>}
   */
  async unsubscribePush(endpoint) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return {}
    }
    return post('/notifications/webpush/unsubscribe', { endpoint })
  }
}

// ========== 报告服务 ==========

export const reportService = {
  /**
   * 获取报告列表
   * @param {string} type - 报告类型 (daily/weekly)
   * @param {number} limit - 限制数量
   * @returns {Promise<object[]>}
   */
  async getReports(type = '', limit = 20) {
    if (USE_MOCK) {
      await simulateDelay(200)
      const now = new Date()
      const reports = []
      for (let i = 0; i < Math.min(limit, 10); i++) {
        const date = new Date(now - i * 86400000)
        const isWeekly = i % 7 === 0 && i > 0
        reports.push({
          id: i + 1,
          type: isWeekly ? 'weekly' : 'daily',
          period: isWeekly
            ? `${date.getFullYear()}-W${String(Math.ceil((date.getDate()) / 7)).padStart(2, '0')}`
            : date.toISOString().slice(0, 10),
          total_value: 125000 + (Math.random() - 0.5) * 10000,
          change: (Math.random() - 0.4) * 3000,
          change_percent: (Math.random() - 0.4) * 5,
          btc_benchmark: (Math.random() - 0.5) * 3,
          eth_benchmark: (Math.random() - 0.5) * 5,
          generated_at: date.toISOString(),
          created_at: date.toISOString()
        })
      }
      if (type) {
        return reports.filter(r => r.type === type)
      }
      return reports
    }
    const params = {}
    if (type) params.type = type
    if (limit) params.limit = limit
    const result = await get('/reports', params)
    // 后端返回 {reports: [...]} 或直接数组
    return Array.isArray(result) ? result : (result.reports || [])
  },

  /**
   * 获取报告详情
   * @param {number} id - 报告 ID
   * @returns {Promise<object>}
   */
  async getReport(id) {
    if (USE_MOCK) {
      await simulateDelay(100)
      return { id, type: 'daily', period: '2026-02-09', total_value: 125430, change: 1230, change_percent: 0.99 }
    }
    const result = await get(`/reports/${id}`)
    return result.report || result
  },

  /**
   * 手动生成报告
   * @param {string} type - 报告类型 (daily/weekly/monthly/annual)
   * @returns {Promise<object>}
   */
  async generateReport(type = 'daily') {
    if (USE_MOCK) {
      await simulateDelay(1000)
      return { report: { id: Date.now(), type, title: `${type} 报告`, total_value: 125000, pnl: 500, pnl_percent: 0.4, created_at: new Date().toISOString() } }
    }
    // 后端接受 POST body { type }
    return post('/reports/generate', { type })
  },

  /**
   * 获取月度报告
   * @param {string} month - 月份 (YYYY-MM)
   * @returns {Promise<object>}
   */
  async getMonthlyReport(month) {
    if (USE_MOCK) {
      await simulateDelay(300)
      const baseValue = 120000
      const dailyReturns = []
      const daysInMonth = 30
      let cumReturn = 0
      for (let i = 1; i <= daysInMonth; i++) {
        const dailyPct = (Math.random() - 0.45) * 2.5
        cumReturn += dailyPct
        dailyReturns.push({
          day: i,
          date: `${month}-${String(i).padStart(2, '0')}`,
          returnPct: +dailyPct.toFixed(2),
          cumReturnPct: +cumReturn.toFixed(2),
        })
      }
      return {
        month,
        totalReturn: +cumReturn.toFixed(2),
        totalReturnAmount: +(baseValue * cumReturn / 100).toFixed(0),
        startValue: baseValue,
        endValue: +(baseValue * (1 + cumReturn / 100)).toFixed(0),
        dailyReturns,
        allocation: {
          start: [
            { symbol: 'BTC', pct: 42 },
            { symbol: 'ETH', pct: 25 },
            { symbol: 'USDC', pct: 18 },
            { symbol: '其他', pct: 15 },
          ],
          end: [
            { symbol: 'BTC', pct: 45 },
            { symbol: 'ETH', pct: 23 },
            { symbol: 'USDC', pct: 17 },
            { symbol: '其他', pct: 15 },
          ],
        },
        feeSummary: {
          totalFee: 186.5,
          tradingFee: 112.3,
          gasFee: 54.2,
          withdrawFee: 20.0,
        },
        suggestions: [
          '本月 BTC 配比上升 3%，注意再平衡',
          'Gas 费支出较上月减少 15%，继续保持',
          '建议关注 ETH 持仓变化趋势',
        ],
        btcBenchmark: +(Math.random() * 10 - 3).toFixed(2),
        ethBenchmark: +(Math.random() * 15 - 5).toFixed(2),
      }
    }
    return get(`/reports/monthly/${month}`)
  },

  /**
   * 报告对比
   * @param {number} reportId1 - 报告 1 ID
   * @param {number} reportId2 - 报告 2 ID
   * @returns {Promise<object>}
   */
  async compareReports(reportId1, reportId2) {
    if (USE_MOCK) {
      await simulateDelay(300)
      const genReport = (id) => ({
        id,
        type: 'daily',
        period: '2026-02-' + String(id).padStart(2, '0'),
        total_value: 120000 + Math.floor(Math.random() * 25000),
        change: +((Math.random() - 0.4) * 3000).toFixed(2),
        change_percent: +((Math.random() - 0.4) * 5).toFixed(2),
        generated_at: new Date().toISOString()
      })
      return {
        report_1: genReport(reportId1),
        report_2: genReport(reportId2),
        value_diff: +((Math.random() - 0.5) * 5000).toFixed(2),
        change_diff: +((Math.random() - 0.5) * 3).toFixed(2)
      }
    }
    return get('/reports/compare', { report_id_1: reportId1, report_id_2: reportId2 })
  },
}

// ========== 价格预警服务 ==========

export const priceAlertService = {
  /**
   * 获取预警列表
   * @returns {Promise<object[]>}
   */
  async getAlerts() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return [
        { id: 1, symbol: 'BTC', condition: 'above', target_price: 100000, is_active: true, triggered: false, note: '', created_at: new Date().toISOString() },
        { id: 2, symbol: 'ETH', condition: 'below', target_price: 2000, is_active: true, triggered: false, note: '抄底提醒', created_at: new Date().toISOString() },
        { id: 3, symbol: 'BTC', condition: 'below', target_price: 50000, is_active: false, triggered: true, triggered_at: new Date(Date.now() - 86400000).toISOString(), note: '', created_at: new Date(Date.now() - 172800000).toISOString() }
      ]
    }
    const result = await get('/alerts')
    // 后端返回 {alerts: [...]}，前端期望数组
    return Array.isArray(result) ? result : (result.alerts || [])
  },

  /**
   * 创建预警
   * @param {object} data - { symbol, condition, target_price, note }
   * @returns {Promise<object>}
   */
  async createAlert(data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return { id: Date.now(), ...data, is_active: true, triggered: false, created_at: new Date().toISOString() }
    }
    const result = await post('/alerts', data)
    return result.alert || result
  },

  /**
   * 更新预警
   * @param {number} id - 预警 ID
   * @param {object} data - 更新数据
   * @returns {Promise<object>}
   */
  async updateAlert(id, data) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return { id, ...data }
    }
    const result = await put(`/alerts/${id}`, data)
    return result.alert || result
  },

  /**
   * 删除预警
   * @param {number} id - 预警 ID
   * @returns {Promise<void>}
   */
  async deleteAlert(id) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return true
    }
    return del(`/alerts/${id}`)
  }
}

// ========== 分析服务 ==========

export const analyticsService = {
  /**
   * 获取每日盈亏
   * @param {string} range - 时间范围（7d/30d/90d/1y）
   */
  async getDailyPnL(range = '30d') {
    if (USE_MOCK) {
      await simulateDelay(300)
      const daysCount = { '7d': 7, '30d': 30, '90d': 90, '1y': 365 }[range] || 30
      const daily = []
      let totalValue = 50000 + Math.random() * 10000
      for (let i = 0; i < daysCount; i++) {
        const date = new Date(Date.now() - (daysCount - i) * 86400000)
        const pnl = (Math.random() - 0.45) * 800
        totalValue += pnl
        daily.push({
          date: date.toISOString().split('T')[0],
          pnl: +pnl.toFixed(2),
          pnl_percent: +((pnl / totalValue) * 100).toFixed(2),
          total_value: +totalValue.toFixed(2)
        })
      }
      return { daily }
    }
    // 后端接受 days 整数参数
    const daysMap = { '7d': 7, '30d': 30, '90d': 90, '1y': 365 }
    const days = daysMap[range] || 30
    return get('/analytics/pnl/daily', { days })
  },

  /**
   * 获取盈亏摘要（后端无参数）
   */
  async getPnLSummary() {
    if (USE_MOCK) {
      await simulateDelay(200)
      return {
        total_pnl: (Math.random() - 0.3) * 5000,
        total_pnl_percent: (Math.random() - 0.3) * 10,
        pnl_7d: (Math.random() - 0.3) * 2000,
        pnl_30d: (Math.random() - 0.3) * 5000,
        pnl_90d: (Math.random() - 0.3) * 10000,
        best_day: '2026-02-05',
        worst_day: '2026-01-20',
        best_day_pnl: 1200,
        worst_day_pnl: -800
      }
    }
    return get('/analytics/pnl/summary')
  },

  /**
   * 获取趋势预测
   * @param {number} days - 预测天数
   */
  async getForecast(days = 30) {
    if (USE_MOCK) {
      await simulateDelay(400)
      const forecastPoints = []
      let value = 55000
      for (let i = 1; i <= days; i++) {
        const date = new Date(Date.now() + i * 86400000)
        value += (Math.random() - 0.4) * 300
        forecastPoints.push({
          date: date.toISOString().split('T')[0],
          value: +value.toFixed(2),
          lower: +(value * 0.95).toFixed(2),
          upper: +(value * 1.05).toFixed(2)
        })
      }
      return {
        forecast_points: forecastPoints,
        trend: 'up',
        confidence: +(0.65 + Math.random() * 0.3).toFixed(2),
        slope: +(50 + Math.random() * 100).toFixed(2),
        days
      }
    }
    return get('/analytics/forecast', { days })
  },

  /**
   * 获取资产归因分析
   * @param {string} range - 时间范围（7d/30d/90d/1y）
   * @param {string} currency - 计价货币
   */
  async getAttribution(range = '7d', currency = 'USD') {
    if (USE_MOCK) {
      await simulateDelay(300)
      return {
        total_return: (Math.random() - 0.3) * 3000,
        total_percent: (Math.random() - 0.3) * 5,
        attributions: [],
        days: { '7d': 7, '30d': 30, '90d': 90, '1y': 365 }[range] || 7,
        currency
      }
    }
    // 后端接受 days 和 currency 参数
    const daysMap = { '7d': 7, '30d': 30, '90d': 90, '1y': 365 }
    const days = daysMap[range] || 7
    return get('/analytics/attribution', { days, currency })
  }
}

// ========== 目标追踪服务 ==========

export const goalService = {
  /**
   * 获取所有投资目标
   * @returns {Promise<object[]>}
   */
  async getGoals() {
    if (USE_MOCK) {
      await simulateDelay(200)
      try {
        const saved = localStorage.getItem('allfi-goals')
        return saved ? JSON.parse(saved) : []
      } catch { return [] }
    }
    const result = await get('/goals')
    const goals = Array.isArray(result) ? result : (result.goals || [])
    return goals.map(g => ({
      id: String(g.id),
      title: g.title,
      type: g.type,
      targetValue: g.target_value,
      currency: g.currency || 'USD',
      deadline: g.deadline || null,
      createdAt: g.created_at
    }))
  },

  /**
   * 创建投资目标
   * @param {object} data - 目标数据
   * @returns {Promise<object>}
   */
  async createGoal(data) {
    if (USE_MOCK) {
      await simulateDelay(300)
      return { id: String(Date.now()), ...data, createdAt: new Date().toISOString() }
    }
    const payload = {
      title: data.title,
      type: data.type,
      target_value: data.targetValue,
      currency: data.currency || 'USD',
      deadline: data.deadline || null
    }
    const result = await post('/goals', payload)
    const g = result.goal || result
    return {
      id: String(g.id),
      title: g.title,
      type: g.type,
      targetValue: g.target_value,
      currency: g.currency || 'USD',
      deadline: g.deadline || null,
      createdAt: g.created_at
    }
  },

  /**
   * 更新投资目标
   * @param {string} id - 目标 ID
   * @param {object} data - 更新数据
   * @returns {Promise<object>}
   */
  async updateGoal(id, data) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return { id, ...data }
    }
    const payload = {}
    if (data.title !== undefined) payload.title = data.title
    if (data.type !== undefined) payload.type = data.type
    if (data.targetValue !== undefined) payload.target_value = data.targetValue
    if (data.currency !== undefined) payload.currency = data.currency
    if (data.deadline !== undefined) payload.deadline = data.deadline
    const result = await put(`/goals/${id}`, payload)
    const g = result.goal || result
    return {
      id: String(g.id),
      title: g.title,
      type: g.type,
      targetValue: g.target_value,
      currency: g.currency || 'USD',
      deadline: g.deadline || null,
      createdAt: g.created_at
    }
  },

  /**
   * 删除投资目标
   * @param {string} id - 目标 ID
   * @returns {Promise<void>}
   */
  async deleteGoal(id) {
    if (USE_MOCK) {
      await simulateDelay(200)
      return true
    }
    return del(`/goals/${id}`)
  }
}

// ========== 投资组合服务 ==========

export const portfolioService = {
  /**
   * 获取资产健康评分（后端计算）
   * @param {string} currency - 计价货币，默认 USD
   * @returns {Promise<{overall_score, level, details, currency, updated_at}>}
   */
  async getHealthScore(currency = 'USD') {
    if (USE_MOCK) {
      await simulateDelay(400)
      return {
        overall_score: 60 + Math.random() * 35,
        level: 'good',
        details: [
          { category: '集中度风险', score: 55 + Math.random() * 30, weight: 0.25, description: '单个资产占比偏高', suggestion: '考虑分散投资，建议将最大单个资产占比控制在 30% 以下' },
          { category: '平台风险', score: 70 + Math.random() * 25, weight: 0.25, description: '资产分散在多个平台', suggestion: '继续保持平台多元化' },
          { category: '波动性风险', score: 60 + Math.random() * 30, weight: 0.25, description: '组合波动率处于中等水平', suggestion: '可适度增加稳定币配置以降低风险' },
          { category: '缓冲能力', score: 65 + Math.random() * 30, weight: 0.25, description: '现金储备充足', suggestion: '维持当前现金配置水平' },
        ],
        currency,
        updated_at: new Date().toISOString()
      }
    }
    return get('/portfolio/health', { currency })
  }
}

// ========== 系统服务 ==========

export const systemService = {
  /**
   * 健康检查
   * @returns {Promise<boolean>}
   */
  async healthCheck() {
    return healthCheck()
  }
}

// 导出所有服务
export default {
  auth: authService,
  asset: assetService,
  cex: cexService,
  wallet: walletService,
  manualAsset: manualAssetService,
  rate: rateService,
  settings: settingsService,
  notification: notificationService,
  priceAlert: priceAlertService,
  report: reportService,
  analytics: analyticsService,
  goal: goalService,
  portfolio: portfolioService,
  system: systemService
}
