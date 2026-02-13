/**
 * 投资组合健康评分 Composable
 * 基于持仓数据计算 4 个维度的评分，返回总分 + 等级 + 各维度详情
 */
import { computed } from 'vue'
import { useAssetStore } from '../stores/assetStore'

// 稳定币列表
const STABLECOINS = ['USDC', 'USDT', 'DAI', 'BUSD', 'TUSD', 'USDP', 'FRAX', 'LUSD', 'GUSD']

// 高市值币种列表（BTC、ETH + 稳定币）
const HIGH_CAP_COINS = ['BTC', 'ETH', ...STABLECOINS]

// 评分等级映射
const GRADES = [
  { min: 80, grade: 'excellent', color: 'var(--color-success)' },
  { min: 60, grade: 'good', color: 'var(--color-warning)' },
  { min: 40, grade: 'warning', color: '#f97316' },
  { min: 0, grade: 'danger', color: 'var(--color-error)' },
]

/**
 * 计算投资组合健康评分
 */
export function useHealthScore() {
  const assetStore = useAssetStore()

  // 收集所有持仓
  const allHoldings = computed(() => {
    const holdings = []
    const cex = assetStore.cexAccounts || []
    const wallets = assetStore.walletAddresses || []
    const manual = assetStore.manualAssets || []

    for (const acc of cex) {
      if (acc.holdings) {
        for (const h of acc.holdings) {
          holdings.push({ symbol: h.symbol, value: h.value, source: acc.name, sourceType: 'cex' })
        }
      }
    }
    for (const w of wallets) {
      if (w.holdings) {
        for (const h of w.holdings) {
          holdings.push({ symbol: h.symbol, value: h.value, source: w.name, sourceType: 'blockchain' })
        }
      }
    }
    for (const m of manual) {
      holdings.push({ symbol: m.currency || 'OTHER', value: m.balance * (m.currency === 'CNY' ? 0.14 : 1), source: m.name, sourceType: 'manual' })
    }

    return holdings
  })

  // 总资产价值
  const totalValue = computed(() => allHoldings.value.reduce((sum, h) => sum + (h.value || 0), 0))

  // 是否有数据
  const hasData = computed(() => totalValue.value > 0 && allHoldings.value.length > 0)

  // ===== 维度 1：现金缓冲（25 分）=====
  // 稳定币占比 >= 20% 得满分
  const cashBufferScore = computed(() => {
    if (!hasData.value) return { score: 0, maxScore: 25, ratio: 0 }
    const stableValue = allHoldings.value
      .filter(h => STABLECOINS.includes(h.symbol.toUpperCase()))
      .reduce((sum, h) => sum + h.value, 0)
    const ratio = stableValue / totalValue.value
    // 线性映射：0% -> 0 分, 20% -> 25 分
    const score = Math.min(25, Math.round((ratio / 0.2) * 25))
    return { score, maxScore: 25, ratio }
  })

  // ===== 维度 2：集中度风险（30 分）=====
  // 最大单一资产占比 <= 25% 得满分
  const concentrationScore = computed(() => {
    if (!hasData.value) return { score: 0, maxScore: 30, maxRatio: 0 }
    // 按 symbol 合并
    const grouped = {}
    for (const h of allHoldings.value) {
      const key = h.symbol.toUpperCase()
      grouped[key] = (grouped[key] || 0) + h.value
    }
    const maxValue = Math.max(...Object.values(grouped))
    const maxRatio = maxValue / totalValue.value
    // 最大占比 <= 25% 满分，100% -> 0 分，线性
    const score = maxRatio <= 0.25
      ? 30
      : Math.max(0, Math.round(30 * (1 - (maxRatio - 0.25) / 0.75)))
    return { score, maxScore: 30, maxRatio }
  })

  // ===== 维度 3：平台分散度（20 分）=====
  // 资产分布在 3+ 个平台得满分
  const platformScore = computed(() => {
    if (!hasData.value) return { score: 0, maxScore: 20, platformCount: 0 }
    const sources = new Set(allHoldings.value.map(h => h.source))
    const platformCount = sources.size
    // 1 平台 -> 7 分, 2 平台 -> 14 分, 3+ -> 20 分
    const score = Math.min(20, Math.round((platformCount / 3) * 20))
    return { score, maxScore: 20, platformCount }
  })

  // ===== 维度 4：波动性评估（25 分）=====
  // 高市值币种占比 >= 60% 得满分
  const volatilityScore = computed(() => {
    if (!hasData.value) return { score: 0, maxScore: 25, highCapRatio: 0 }
    const highCapValue = allHoldings.value
      .filter(h => HIGH_CAP_COINS.includes(h.symbol.toUpperCase()))
      .reduce((sum, h) => sum + h.value, 0)
    const highCapRatio = highCapValue / totalValue.value
    // 线性映射：0% -> 0 分, 60% -> 25 分
    const score = Math.min(25, Math.round((highCapRatio / 0.6) * 25))
    return { score, maxScore: 25, highCapRatio }
  })

  // 总分与等级
  const healthScore = computed(() => {
    const total = cashBufferScore.value.score
      + concentrationScore.value.score
      + platformScore.value.score
      + volatilityScore.value.score

    const gradeInfo = GRADES.find(g => total >= g.min) || GRADES[GRADES.length - 1]

    const dimensions = [
      {
        id: 'cashBuffer',
        ...cashBufferScore.value,
      },
      {
        id: 'concentration',
        ...concentrationScore.value,
      },
      {
        id: 'platform',
        ...platformScore.value,
      },
      {
        id: 'volatility',
        ...volatilityScore.value,
      },
    ]

    // 找出最薄弱维度（得分比例最低）用于建议
    let weakest = null
    if (hasData.value) {
      let minRatio = 1
      for (const d of dimensions) {
        const ratio = d.score / d.maxScore
        if (ratio < minRatio) {
          minRatio = ratio
          weakest = d.id
        }
      }
    }

    return {
      total,
      grade: gradeInfo.grade,
      color: gradeInfo.color,
      hasData: hasData.value,
      dimensions,
      weakest,
    }
  })

  return {
    healthScore,
  }
}
