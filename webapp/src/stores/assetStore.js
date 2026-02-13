/**
 * 资产数据 Store
 * 管理资产汇总、历史数据、汇率等
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { assetService, rateService } from '../api/index.js'
import { defiService } from '../api/defiService.js'
import { useNFTStore } from './nftStore'

export const useAssetStore = defineStore('asset', () => {
  // ========== 状态 ==========

  // 资产汇总
  const summary = ref(null)
  const summaryLoading = ref(false)
  const summaryError = ref(null)

  // 历史数据
  const historyData = ref(null)
  const historyLoading = ref(false)
  const historyError = ref(null)
  const currentTimeRange = ref('30D')

  // 汇率数据
  const exchangeRates = ref({
    BTC: 42500,
    ETH: 2280,
    CNY: 0.14
  })
  const ratesLoading = ref(false)
  const ratesError = ref(null)

  // DeFi 仓位数据
  const defiPositions = ref([])
  const defiLoading = ref(false)
  const defiError = ref(null)

  // 刷新状态
  const isRefreshing = ref(false)
  const lastRefreshTime = ref(null)

  // 离线状态
  const isOffline = ref(!navigator.onLine)

  // ========== 计算属性 ==========

  // 总资产价值
  const totalValue = computed(() => summary.value?.totalValue || 0)

  // 24小时变化
  const change24h = computed(() => summary.value?.change24h || 0)
  const changeValue = computed(() => summary.value?.changeValue || 0)

  // 分类数据
  const categories = computed(() => summary.value?.categories || {
    cex: { value: 0, count: 0, percentage: 0 },
    blockchain: { value: 0, count: 0, percentage: 0 },
    manual: { value: 0, count: 0, percentage: 0 }
  })

  // CEX 账户列表
  const cexAccounts = computed(() => summary.value?.cexAccounts || [])

  // 钱包地址列表
  const walletAddresses = computed(() => summary.value?.walletAddresses || [])

  // 手动资产列表
  const manualAssets = computed(() => summary.value?.manualAssets || [])

  // DeFi 总价值
  const defiTotalValue = computed(() => {
    return defiPositions.value.reduce((sum, p) => sum + (p.valueUSD || 0), 0)
  })

  // 是否有数据
  const hasData = computed(() => summary.value !== null)

  // 今日盈亏（基于 24h 变化）
  const todayPnL = computed(() => {
    const value = changeValue.value
    const percent = change24h.value
    return {
      value: value || 0,
      percent: percent || 0,
      isPositive: (value || 0) >= 0
    }
  })

  // 基准对比（对比 BTC/ETH 表现）
  const benchmarkComparison = computed(() => {
    const myChange = change24h.value || 0
    // 从持仓中获取 BTC 和 ETH 的 24h 变化
    const holdings = summary.value?.cexAccounts?.flatMap(a => a.holdings || []) || []
    const walletHoldings = summary.value?.walletAddresses?.flatMap(w => w.holdings || []) || []
    const allHoldings = [...holdings, ...walletHoldings]

    const btcHolding = allHoldings.find(h => h.symbol === 'BTC')
    const ethHolding = allHoldings.find(h => h.symbol === 'ETH')

    const btcChange = btcHolding?.change24h || 0
    const ethChange = ethHolding?.change24h || 0

    return {
      myChange,
      btc: { change: btcChange, diff: myChange - btcChange },
      eth: { change: ethChange, diff: myChange - ethChange }
    }
  })

  // 资产集中度分析（Top 5 资产 + HHI 指数）
  const assetConcentration = computed(() => {
    const holdings = []
    for (const acc of cexAccounts.value) {
      if (acc.holdings) {
        for (const h of acc.holdings) {
          holdings.push({ symbol: h.symbol, value: h.value })
        }
      }
    }
    for (const w of walletAddresses.value) {
      if (w.holdings) {
        for (const h of w.holdings) {
          holdings.push({ symbol: h.symbol, value: h.value })
        }
      }
    }

    // 按 symbol 合并
    const grouped = {}
    for (const h of holdings) {
      grouped[h.symbol] = (grouped[h.symbol] || 0) + h.value
    }

    const total = Object.values(grouped).reduce((s, v) => s + v, 0) || 1
    const sorted = Object.entries(grouped)
      .map(([symbol, value]) => ({ symbol, value, percentage: (value / total) * 100 }))
      .sort((a, b) => b.value - a.value)

    const top5 = sorted.slice(0, 5)
    const othersValue = sorted.slice(5).reduce((s, item) => s + item.value, 0)
    if (othersValue > 0) {
      top5.push({ symbol: 'Others', value: othersValue, percentage: (othersValue / total) * 100 })
    }

    // HHI 指数（赫芬达尔指数）：各资产占比平方和 * 10000
    const hhi = sorted.reduce((sum, item) => sum + Math.pow(item.percentage / 100, 2), 0) * 10000

    return { top5, hhi: Math.round(hhi), total }
  })

  // 平台分布分析（含 NFT）
  const platformDistribution = computed(() => {
    const cexTotal = cexAccounts.value.reduce((s, a) => s + a.balance, 0)
    const chainTotal = walletAddresses.value.reduce((s, w) => s + w.balance, 0)
    const manualTotal = manualAssets.value.reduce((s, a) => s + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0)
    const defiTotal = defiTotalValue.value

    // 如果用户选择将 NFT 计入总资产，纳入分布
    let nftTotal = 0
    try {
      const nftStore = useNFTStore()
      if (nftStore.includeInTotal) {
        nftTotal = nftStore.totalFloorValue
      }
    } catch {}

    const total = cexTotal + chainTotal + manualTotal + defiTotal + nftTotal || 1

    const result = [
      { id: 'cex', value: cexTotal, percentage: (cexTotal / total) * 100 },
      { id: 'blockchain', value: chainTotal, percentage: (chainTotal / total) * 100 },
      { id: 'manual', value: manualTotal, percentage: (manualTotal / total) * 100 },
      { id: 'defi', value: defiTotal, percentage: (defiTotal / total) * 100 },
    ]
    if (nftTotal > 0) {
      result.push({ id: 'nft', value: nftTotal, percentage: (nftTotal / total) * 100 })
    }
    return result
  })

  // 收益率（基于历史快照的简单计算）
  const returnRate = computed(() => {
    if (!historyData.value || !historyData.value.values || historyData.value.values.length < 2) {
      return { totalReturn: 0, annualizedReturn: 0 }
    }
    const values = historyData.value.values
    const first = values[0] || 1
    const last = values[values.length - 1]
    const totalReturn = ((last - first) / first) * 100
    // 简单年化：按数据点数近似天数
    const days = values.length
    const annualizedReturn = days > 0 ? (totalReturn / days) * 365 : 0
    return { totalReturn: parseFloat(totalReturn.toFixed(2)), annualizedReturn: parseFloat(annualizedReturn.toFixed(2)) }
  })

  // 是否正在加载
  const isLoading = computed(() => summaryLoading.value || historyLoading.value || ratesLoading.value || defiLoading.value)

  // ========== 资产汇总操作 ==========

  // 离线缓存键
  const CACHE_SUMMARY_KEY = 'allfi-cache-summary'
  const CACHE_HISTORY_KEY = 'allfi-cache-history'

  // 缓存数据到 localStorage
  function cacheData(key, data) {
    try {
      localStorage.setItem(key, JSON.stringify({ data, timestamp: Date.now() }))
    } catch (e) {
      // localStorage 满时忽略
    }
  }

  // 从 localStorage 读取缓存
  function getCachedData(key) {
    try {
      const raw = localStorage.getItem(key)
      if (raw) return JSON.parse(raw).data
    } catch (e) {
      // 解析失败忽略
    }
    return null
  }

  /**
   * 加载资产汇总
   */
  async function loadSummary() {
    summaryLoading.value = true
    summaryError.value = null

    try {
      summary.value = await assetService.getSummary()
      lastRefreshTime.value = new Date()
      // 缓存成功数据
      cacheData(CACHE_SUMMARY_KEY, summary.value)
    } catch (err) {
      // 离线时尝试读取缓存
      const cached = getCachedData(CACHE_SUMMARY_KEY)
      if (cached) {
        summary.value = cached
        isOffline.value = true
      } else {
        summaryError.value = err.message || '加载资产汇总失败'
      }
      console.error('加载资产汇总失败:', err)
    } finally {
      summaryLoading.value = false
    }
  }

  /**
   * 刷新所有资产数据
   */
  async function refreshAll() {
    isRefreshing.value = true
    summaryError.value = null

    try {
      summary.value = await assetService.refresh()
      lastRefreshTime.value = new Date()

      // 同时刷新历史数据
      await loadHistory(currentTimeRange.value)
    } catch (err) {
      summaryError.value = err.message || '刷新资产数据失败'
      throw err
    } finally {
      isRefreshing.value = false
    }
  }

  // ========== 历史数据操作 ==========

  /**
   * 加载历史数据
   * @param {string} timeRange - 时间范围 (7D, 30D, 90D, 1Y, ALL)
   */
  async function loadHistory(timeRange = '30D') {
    historyLoading.value = true
    historyError.value = null
    currentTimeRange.value = timeRange

    try {
      historyData.value = await assetService.getHistory(timeRange)
      // 缓存成功数据
      cacheData(CACHE_HISTORY_KEY, historyData.value)
    } catch (err) {
      // 离线时尝试读取缓存
      const cached = getCachedData(CACHE_HISTORY_KEY)
      if (cached) {
        historyData.value = cached
        isOffline.value = true
      } else {
        historyError.value = err.message || '加载历史数据失败'
      }
      console.error('加载历史数据失败:', err)
    } finally {
      historyLoading.value = false
    }
  }

  /**
   * 切换时间范围
   * @param {string} timeRange - 时间范围
   */
  async function setTimeRange(timeRange) {
    if (timeRange !== currentTimeRange.value) {
      await loadHistory(timeRange)
    }
  }

  // ========== DeFi 仓位操作 ==========

  /**
   * 加载 DeFi 仓位数据
   */
  async function loadDefiPositions() {
    defiLoading.value = true
    defiError.value = null

    try {
      defiPositions.value = await defiService.getPositions()
    } catch (err) {
      defiError.value = err.message || '加载 DeFi 仓位失败'
      console.error('加载 DeFi 仓位失败:', err)
    } finally {
      defiLoading.value = false
    }
  }

  // ========== 汇率操作 ==========

  /**
   * 加载当前汇率
   */
  async function loadRates() {
    ratesLoading.value = true
    ratesError.value = null

    try {
      const rates = await rateService.getCurrentRates()
      exchangeRates.value = rates
    } catch (err) {
      ratesError.value = err.message || '加载汇率失败'
      console.error('加载汇率失败:', err)
    } finally {
      ratesLoading.value = false
    }
  }

  /**
   * 获取特定货币的汇率
   * @param {string} currency - 货币代码
   * @returns {number}
   */
  function getRate(currency) {
    return exchangeRates.value[currency] || 1
  }

  /**
   * 转换货币价值
   * @param {number} value - 原始价值（USDC）
   * @param {string} toCurrency - 目标货币
   * @returns {number}
   */
  function convertValue(value, toCurrency) {
    if (toCurrency === 'USDC') {
      return value
    }
    const rate = getRate(toCurrency)
    return value / rate
  }

  // ========== 初始化 ==========

  /**
   * 初始化所有数据
   */
  async function initialize() {
    await Promise.all([
      loadSummary(),
      loadHistory(currentTimeRange.value),
      loadRates(),
      loadDefiPositions()
    ])
  }

  /**
   * 清除所有错误
   */
  function clearErrors() {
    summaryError.value = null
    historyError.value = null
    ratesError.value = null
  }

  /**
   * 重置状态
   */
  function reset() {
    summary.value = null
    historyData.value = null
    defiPositions.value = []
    lastRefreshTime.value = null
    clearErrors()
  }

  return {
    // 状态
    summary,
    summaryLoading,
    summaryError,
    historyData,
    historyLoading,
    historyError,
    currentTimeRange,
    exchangeRates,
    ratesLoading,
    ratesError,
    isRefreshing,
    lastRefreshTime,
    isOffline,
    defiPositions,
    defiLoading,
    defiError,

    // 计算属性
    totalValue,
    change24h,
    changeValue,
    categories,
    cexAccounts,
    walletAddresses,
    manualAssets,
    defiTotalValue,
    hasData,
    isLoading,
    todayPnL,
    benchmarkComparison,
    assetConcentration,
    platformDistribution,
    returnRate,

    // 资产汇总操作
    loadSummary,
    refreshAll,

    // DeFi 操作
    loadDefiPositions,

    // 历史数据操作
    loadHistory,
    setTimeRange,

    // 汇率操作
    loadRates,
    getRate,
    convertValue,

    // 通用操作
    initialize,
    clearErrors,
    reset
  }
})
