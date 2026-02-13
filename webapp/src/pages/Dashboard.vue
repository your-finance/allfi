<script setup>
/**
 * Dashboard 页面 - 资产总览
 * 一屏平铺式布局：摘要栏 → 图表区 → 持仓表格
 */
import { ref, computed, onMounted, watch, useTemplateRef } from 'vue'
import {
  PhCaretUp,
  PhCaretDown,
  PhMagnifyingGlass,
  PhSortAscending,
  PhSortDescending,
  PhCaretRight,
  PhBell,
  PhShareNetwork,
  PhGear,
  PhListBullets,
  PhSquaresFour
} from '@phosphor-icons/vue'
import { Line, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import AssetDetailDrawer from '../components/AssetDetailDrawer.vue'
import PriceAlertDialog from '../components/PriceAlertDialog.vue'
import ShareCard from '../components/ShareCard.vue'
import HealthScoreCard from '../components/HealthScoreCard.vue'
import GoalCard from '../components/GoalCard.vue'
import AddGoalDialog from '../components/AddGoalDialog.vue'
import DashboardCustomizer from '../components/DashboardCustomizer.vue'
import OnboardingWizard from '../components/OnboardingWizard.vue'
import DeFiOverview from '../components/DeFiOverview.vue'
import NFTOverview from '../components/NFTOverview.vue'
import FeeAnalytics from '../components/FeeAnalytics.vue'
import StrategyPanel from '../components/StrategyPanel.vue'
import PullToRefresh from '../components/PullToRefresh.vue'
import { usePullToRefresh } from '../composables/usePullToRefresh'
import { useFormatters } from '../composables/useFormatters'
import { useThemeStore } from '../stores/themeStore'
import { useAssetStore } from '../stores/assetStore'
import { useI18n } from '../composables/useI18n'
import { useGoalStore } from '../stores/goalStore'
import { useDashboardStore } from '../stores/dashboardStore'
import { useNFTStore } from '../stores/nftStore'

// 注册 Chart.js 组件
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const {
  currentCurrency,
  currencySymbol,
  formatNumber,
  formatPercent,
  currencies
} = useFormatters()

const themeStore = useThemeStore()
const assetStore = useAssetStore()
const goalStore = useGoalStore()
const dashboardStore = useDashboardStore()
const nftStore = useNFTStore()
const { t } = useI18n()

// Widget 配置面板
const showCustomizer = ref(false)

// 目标追踪
const showAddGoal = ref(false)
const displayedGoals = computed(() => goalStore.goalsWithProgress.slice(0, 3))

// ========== 状态 ==========
const selectedTimeRange = ref('30D')
const timeRanges = ['7D', '30D', '90D', '1Y', 'ALL']

// 计价货币
const selectedPricingCurrency = ref('USDC')
const showExchangeRate = computed(() => selectedPricingCurrency.value !== 'USDC')
const exchangeRates = computed(() => assetStore.exchangeRates)

const setPricingCurrency = (currency) => {
  selectedPricingCurrency.value = currency
  currentCurrency.value = currency
}

const pricingCurrencySymbol = computed(() => {
  const currency = currencies.find(c => c.code === selectedPricingCurrency.value)
  return currency?.symbol || '$'
})

const pricingCurrencyDecimals = computed(() => {
  switch (selectedPricingCurrency.value) {
    case 'BTC': return 5
    case 'ETH': return 4
    default: return 2
  }
})

// 资产详情抽屉
const selectedAsset = ref(null)
const showAssetDrawer = ref(false)

// 价格预警对话框
const showPriceAlert = ref(false)

// 分享卡片
const showShareCard = ref(false)

// 基准对比选择
const selectedBenchmark = ref('none') // 'none', 'BTC', 'ETH'

// 智能资产分组
const holdingsViewMode = ref('flat') // 'flat' | 'grouped'

// 资产分类规则
const STABLECOINS = new Set(['USDC', 'USDT', 'DAI', 'BUSD', 'TUSD', 'FRAX', 'LUSD', 'GUSD', 'USDP', 'PYUSD', 'FDUSD', 'cUSD'])
const BLUE_CHIPS = new Set(['BTC', 'ETH', 'SOL', 'BNB', 'XRP', 'ADA', 'DOT', 'AVAX', 'MATIC', 'LINK', 'ATOM', 'LTC', 'NEAR', 'APT', 'SUI', 'TON', 'TRX'])

// 判断资产类别
const classifyAsset = (symbol) => {
  const s = symbol.toUpperCase()
  if (STABLECOINS.has(s)) return 'stablecoin'
  if (BLUE_CHIPS.has(s)) return 'bluechip'
  return 'risk'
}

// 展开/折叠的分组
const expandedGroups = ref(new Set(['stablecoin', 'bluechip', 'risk']))

const toggleGroup = (groupId) => {
  if (expandedGroups.value.has(groupId)) {
    expandedGroups.value.delete(groupId)
  } else {
    expandedGroups.value.add(groupId)
  }
}

// 分组后的持仓数据
const groupedHoldings = computed(() => {
  const groups = {
    stablecoin: { id: 'stablecoin', labelKey: 'dashboard.groupStablecoin', items: [], totalValue: 0 },
    bluechip: { id: 'bluechip', labelKey: 'dashboard.groupBluechip', items: [], totalValue: 0 },
    risk: { id: 'risk', labelKey: 'dashboard.groupRisk', items: [], totalValue: 0 }
  }

  for (const holding of filteredHoldings.value) {
    const category = classifyAsset(holding.symbol)
    groups[category].items.push(holding)
    groups[category].totalValue += holding.value
  }

  // 按总值降序排列每组内的资产
  for (const group of Object.values(groups)) {
    group.items.sort((a, b) => b.value - a.value)
    group.percentage = totalAssets.value ? ((group.totalValue / totalAssets.value) * 100).toFixed(1) : '0.0'
  }

  return [groups.stablecoin, groups.bluechip, groups.risk].filter(g => g.items.length > 0)
})

// 搜索和排序
const searchQuery = ref('')
const sortField = ref('value')
const sortOrder = ref('desc')

// ========== 从 Store 获取数据 ==========
const totalAssets = computed(() => assetStore.totalValue)
const totalChange24h = computed(() => assetStore.change24h)
const totalChangeValue = computed(() => assetStore.changeValue)

// CEX / 链上 / 手动资产
const cexAccounts = computed(() => assetStore.cexAccounts)
const blockchainWallets = computed(() => assetStore.walletAddresses)
const manualAssets = computed(() => assetStore.manualAssets)

// 首次使用引导：未完成引导且无任何账户时显示
const showOnboarding = computed(() =>
  !themeStore.onboardingCompleted &&
  cexAccounts.value.length === 0 &&
  blockchainWallets.value.length === 0 &&
  manualAssets.value.length === 0
)

// DeFi 仓位
const defiPositions = computed(() => assetStore.defiPositions)

// 资产分类汇总
const assetCategories = computed(() => {
  const cats = [
    {
      id: 'cex',
      labelKey: 'dashboard.cexAssets',
      colorKey: 'accentPrimary',
      accounts: cexAccounts.value,
      totalValue: cexAccounts.value.reduce((sum, acc) => sum + acc.balance, 0),
      count: cexAccounts.value.length
    },
    {
      id: 'blockchain',
      labelKey: 'dashboard.blockchainAssets',
      colorKey: 'accentSecondary',
      accounts: blockchainWallets.value,
      totalValue: blockchainWallets.value.reduce((sum, w) => sum + w.balance, 0),
      count: blockchainWallets.value.length
    },
    {
      id: 'manual',
      labelKey: 'dashboard.manualAssets',
      colorKey: 'accentTertiary',
      accounts: manualAssets.value,
      totalValue: manualAssets.value.reduce((sum, a) => sum + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0),
      count: manualAssets.value.length
    }
  ]
  // DeFi 仓位有数据时显示
  if (defiPositions.value.length > 0) {
    cats.push({
      id: 'defi',
      labelKey: 'dashboard.defiAssets',
      fixedColor: '#8B5CF6',
      accounts: [],
      totalValue: assetStore.defiTotalValue,
      count: defiPositions.value.length
    })
  }
  // NFT 资产（用户启用计入总资产且有数据时）
  if (nftStore.includeInTotal && nftStore.totalCount > 0) {
    cats.push({
      id: 'nft',
      labelKey: 'dashboard.nftAssets',
      fixedColor: '#EC4899',
      accounts: [],
      totalValue: nftStore.totalFloorValue,
      count: nftStore.totalCount
    })
  }
  return cats
})

// 汇总所有账户的全部持仓（平铺展示）
const allHoldings = computed(() => {
  const holdings = []
  for (const category of assetCategories.value) {
    for (const account of category.accounts) {
      if (account.holdings) {
        for (const h of account.holdings) {
          holdings.push({
            ...h,
            source: account.name,
            sourceType: category.id
          })
        }
      }
    }
  }
  return holdings
})

// 合并相同资产的持仓（按 symbol 分组）
const mergedHoldings = computed(() => {
  const grouped = {}

  for (const holding of allHoldings.value) {
    const key = holding.symbol

    if (!grouped[key]) {
      grouped[key] = {
        symbol: holding.symbol,
        name: holding.name,
        icon: holding.icon,
        price: holding.price,
        change24h: holding.change24h,
        balance: 0,
        value: 0,
        sources: []
      }
    }

    // 累加余额和价值
    grouped[key].balance += holding.balance
    grouped[key].value += holding.value

    // 记录来源
    grouped[key].sources.push({
      source: holding.source,
      sourceType: holding.sourceType,
      balance: holding.balance,
      value: holding.value
    })
  }

  return Object.values(grouped)
})

// 展开的资产 symbol 集合
const expandedAssets = ref(new Set())

// 切换资产来源展开/收起
const toggleAssetExpand = (symbol) => {
  if (expandedAssets.value.has(symbol)) {
    expandedAssets.value.delete(symbol)
  } else {
    expandedAssets.value.add(symbol)
  }
}

// 搜索 + 排序后的持仓
const filteredHoldings = computed(() => {
  let list = mergedHoldings.value

  // 搜索过滤
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(h =>
      h.symbol.toLowerCase().includes(q) ||
      h.name.toLowerCase().includes(q) ||
      h.sources.some(s => s.source.toLowerCase().includes(q))
    )
  }

  // 排序
  const field = sortField.value
  const order = sortOrder.value === 'asc' ? 1 : -1
  list = [...list].sort((a, b) => {
    const va = a[field] ?? 0
    const vb = b[field] ?? 0
    if (typeof va === 'string') return va.localeCompare(vb) * order
    return (va - vb) * order
  })

  return list
})

// ========== 分页 ==========
const currentPage = ref(1)
const itemsPerPage = ref(20)
const availablePageSizes = [10, 20, 50, 100]

const totalPages = computed(() => Math.ceil(filteredHoldings.value.length / itemsPerPage.value))

const paginatedHoldings = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredHoldings.value.slice(start, end)
})

const paginationInfo = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value + 1
  const end = Math.min(start + itemsPerPage.value - 1, filteredHoldings.value.length)
  const total = filteredHoldings.value.length
  return { start, end, total }
})

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const changePageSize = (size) => {
  itemsPerPage.value = size
  currentPage.value = 1
}

// 监听搜索和排序,重置到第一页
watch([searchQuery, sortField, sortOrder], () => {
  currentPage.value = 1
})

// 切换排序
const toggleSort = (field) => {
  if (sortField.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortOrder.value = 'desc'
  }
}

// 计算占比
const holdingPercent = (value) => {
  if (!totalAssets.value) return 0
  return ((value / totalAssets.value) * 100).toFixed(1)
}

// ========== 图表数据 ==========

// 资产分布饼图
const doughnutChartData = computed(() => ({
  labels: assetCategories.value.map(a => t(a.labelKey)),
  datasets: [{
    data: assetCategories.value.map(a => a.totalValue),
    backgroundColor: assetCategories.value.map(a => a.fixedColor || themeStore.currentTheme.colors[a.colorKey]),
    borderWidth: 0,
    cutout: '72%'
  }]
}))

const doughnutChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: 'rgba(0,0,0,0.8)',
      titleFont: { size: 12 },
      bodyFont: { size: 11 },
      padding: 8,
      callbacks: {
        label: (context) => {
          const value = context.raw
          const total = context.dataset.data.reduce((a, b) => a + b, 0)
          const percentage = ((value / total) * 100).toFixed(1)
          return `${currencySymbol.value}${formatNumber(value)} (${percentage}%)`
        }
      }
    }
  }
}

// 趋势图数据
const generateTrendData = () => {
  const historyData = assetStore.historyData
  if (!historyData || !historyData.labels || !historyData.values) {
    return { labels: [], assetData: [], rateData: [] }
  }

  const labels = historyData.labels
  const baseRate = exchangeRates.value[selectedPricingCurrency.value] || 1

  let assetData
  if (selectedPricingCurrency.value === 'USDC') {
    assetData = historyData.values
  } else {
    assetData = historyData.values.map(v => v / baseRate)
  }

  const rateData = []
  if (showExchangeRate.value) {
    for (let i = 0; i < labels.length; i++) {
      const rateVariation = baseRate * (0.92 + Math.random() * 0.16)
      rateData.push(rateVariation)
    }
  }

  return { labels, assetData, rateData }
}

const lineChartData = computed(() => {
  const { labels, assetData, rateData } = generateTrendData()

  const datasets = [{
    label: t('dashboard.totalAssets'),
    data: assetData,
    borderColor: themeStore.currentTheme.colors.accentPrimary,
    backgroundColor: themeStore.currentTheme.colors.chartGradientStart,
    borderWidth: 1.5,
    fill: true,
    tension: 0.4,
    pointRadius: 0,
    pointHoverRadius: 3,
    yAxisID: 'y'
  }]

  if (showExchangeRate.value) {
    datasets.push({
      label: `${selectedPricingCurrency.value}/USDC`,
      data: rateData,
      borderColor: themeStore.currentTheme.colors.accentSecondary,
      backgroundColor: 'transparent',
      borderWidth: 1,
      borderDash: [4, 4],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3,
      yAxisID: 'y1'
    })
  }

  // 基准对比线（BTC 或 ETH）
  if (selectedBenchmark.value !== 'none' && assetData.length > 0) {
    const benchmarkData = []
    const baseValue = assetData[0] || 1
    // 基于基准币种的 24h 变化率模拟历史走势（归一化到资产初始值）
    const benchmarkChange = selectedBenchmark.value === 'BTC'
      ? (assetStore.benchmarkComparison.btc.change || 0)
      : (assetStore.benchmarkComparison.eth.change || 0)

    for (let i = 0; i < assetData.length; i++) {
      // 模拟基准走势：基于每日变化率线性分布
      const progress = i / Math.max(assetData.length - 1, 1)
      const dailyFactor = 1 + (benchmarkChange / 100) * progress
      benchmarkData.push(baseValue * dailyFactor)
    }

    datasets.push({
      label: `${t('dashboard.benchmarkVs')} ${selectedBenchmark.value}`,
      data: benchmarkData,
      borderColor: selectedBenchmark.value === 'BTC' ? '#F7931A' : '#627EEA',
      backgroundColor: 'transparent',
      borderWidth: 1.5,
      borderDash: [6, 3],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3,
      yAxisID: 'y'
    })
  }

  return { labels, datasets }
})

const lineChartOptions = computed(() => {
  const colors = themeStore.currentTheme.colors
  const baseOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: { intersect: false, mode: 'index' },
    plugins: {
      legend: {
        display: showExchangeRate.value || selectedBenchmark.value !== 'none',
        position: 'top',
        align: 'end',
        labels: {
          color: colors.textSecondary,
          usePointStyle: true,
          pointStyle: 'line',
          padding: 12,
          font: { size: 11 }
        }
      },
      tooltip: {
        backgroundColor: colors.bgElevated,
        titleColor: colors.textPrimary,
        bodyColor: colors.textSecondary,
        borderColor: colors.border,
        borderWidth: 1,
        padding: 8,
        titleFont: { size: 11 },
        bodyFont: { size: 11 },
        callbacks: {
          label: (context) => {
            const label = context.dataset.label || ''
            const value = context.raw
            if (context.datasetIndex === 0) {
              return `${label}: ${pricingCurrencySymbol.value}${formatNumber(value, pricingCurrencyDecimals.value)}`
            } else {
              return `${label}: $${formatNumber(value, 2)}`
            }
          }
        }
      }
    },
    scales: {
      x: {
        display: true,
        grid: { display: false },
        ticks: {
          color: colors.textMuted,
          maxTicksLimit: 6,
          font: { size: 10 }
        }
      },
      y: {
        type: 'linear',
        display: true,
        position: 'left',
        grid: {
          color: colors.border,
          lineWidth: 0.5
        },
        ticks: {
          color: colors.textSecondary,
          font: { size: 10 },
          callback: (value) => `${pricingCurrencySymbol.value}${formatNumber(value, pricingCurrencyDecimals.value)}`
        }
      }
    }
  }

  if (showExchangeRate.value) {
    baseOptions.scales.y1 = {
      type: 'linear',
      display: true,
      position: 'right',
      grid: { drawOnChartArea: false },
      ticks: {
        color: colors.accentSecondary,
        font: { size: 10 },
        callback: (value) => `$${formatNumber(value, 2)}`
      }
    }
  }

  return baseOptions
})

// ========== 交互 ==========
const openAssetDetail = (asset) => {
  selectedAsset.value = asset
  showAssetDrawer.value = true
}

const closeAssetDrawer = () => {
  showAssetDrawer.value = false
  selectedAsset.value = null
}

// ========== 下拉刷新 ==========
const dashboardRef = useTemplateRef('dashboardRef')
const { pullDistance, isRefreshing: pullRefreshing } = usePullToRefresh(
  () => assetStore.refreshAll(),
  dashboardRef
)

// ========== 初始化 ==========
onMounted(async () => {
  await assetStore.initialize()
  goalStore.loadGoals()
})

watch(selectedTimeRange, async (newRange) => {
  await assetStore.setTimeRange(newRange)
})
</script>

<template>
  <div ref="dashboardRef" class="dashboard">
    <!-- 下拉刷新指示器 -->
    <PullToRefresh :pull-distance="pullDistance" :is-refreshing="pullRefreshing" />

    <!-- 顶部：标题 + 配置齿轮 -->
    <div class="dashboard-header">
      <h2 class="page-title">{{ t('dashboard.title') }}</h2>
      <button class="gear-btn" @click="showCustomizer = true" :title="t('widgets.customize')">
        <PhGear :size="16" weight="bold" />
      </button>
    </div>

    <!-- 全部隐藏时的引导 -->
    <div v-if="dashboardStore.enabledCount === 0" class="empty-dashboard">
      <p>{{ t('widgets.emptyHint') }}</p>
      <button class="btn btn-primary" @click="showCustomizer = true">{{ t('widgets.customize') }}</button>
    </div>

    <!-- 摘要栏：总资产 + 今日盈亏 + 分类分布 -->
    <section v-if="dashboardStore.widgetConfig.assetSummary" class="summary-bar">
      <div class="summary-left">
        <span class="stat-label">{{ t('dashboard.totalAssets') }}</span>
        <div class="total-row">
          <h2 class="stat-value-hero font-mono">
            {{ currencySymbol }}{{ formatNumber(totalAssets) }}
          </h2>
          <span
            class="change-badge"
            :class="totalChange24h >= 0 ? 'change-positive' : 'change-negative'"
          >
            <PhCaretUp v-if="totalChange24h >= 0" :size="12" weight="bold" />
            <PhCaretDown v-else :size="12" weight="bold" />
            {{ formatPercent(totalChange24h) }}
            <span class="change-abs">({{ currencySymbol }}{{ formatNumber(totalChangeValue) }})</span>
          </span>
          <span class="change-period">24h</span>
        </div>
      </div>
      <!-- 今日盈亏卡片 -->
      <div class="pnl-card" :class="assetStore.todayPnL.isPositive ? 'pnl-positive' : 'pnl-negative'">
        <span class="pnl-label">{{ t('dashboard.todayPnL') }}</span>
        <div class="pnl-value font-mono">
          <span>{{ assetStore.todayPnL.isPositive ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(assetStore.todayPnL.value) }}</span>
        </div>
        <span class="pnl-percent font-mono" :class="assetStore.todayPnL.isPositive ? 'change-positive' : 'change-negative'">
          <PhCaretUp v-if="assetStore.todayPnL.isPositive" :size="10" weight="bold" />
          <PhCaretDown v-else :size="10" weight="bold" />
          {{ formatPercent(assetStore.todayPnL.percent) }}
        </span>
      </div>
      <div class="summary-right">
        <button class="alert-btn" @click="showShareCard = true" :title="t('share.title')">
          <PhShareNetwork :size="16" weight="bold" />
        </button>
        <button class="alert-btn" @click="showPriceAlert = true" :title="t('priceAlert.title')">
          <PhBell :size="16" weight="bold" />
        </button>
        <div
          v-for="cat in assetCategories"
          :key="cat.id"
          class="category-chip"
        >
          <span
            class="chip-dot"
            :style="{ background: cat.fixedColor || themeStore.currentTheme.colors[cat.colorKey] }"
          />
          <span class="chip-label">{{ t(cat.labelKey) }}</span>
          <span class="chip-value font-mono">{{ currencySymbol }}{{ formatNumber(cat.totalValue, 0) }}</span>
        </div>
      </div>
    </section>

    <!-- 图表区：趋势图（60%） + 饼图（40%） -->
    <section v-if="dashboardStore.widgetConfig.trend || dashboardStore.widgetConfig.distribution" class="charts-row">
      <!-- 趋势图 -->
      <div v-if="dashboardStore.widgetConfig.trend" class="glass-card trend-panel">
        <div class="panel-header">
          <h3>{{ t('dashboard.assetTrend') }}</h3>
          <div class="panel-controls">
            <!-- 计价货币选择 -->
            <div class="selector-group">
              <button
                v-for="curr in currencies"
                :key="curr.code"
                class="selector-btn"
                :class="{ active: selectedPricingCurrency === curr.code }"
                @click="setPricingCurrency(curr.code)"
              >
                {{ curr.symbol }} {{ curr.code }}
              </button>
            </div>
            <!-- 基准对比选择 -->
            <div class="selector-group">
              <button
                class="selector-btn"
                :class="{ active: selectedBenchmark === 'none' }"
                @click="selectedBenchmark = 'none'"
              >
                {{ t('dashboard.benchmarkNone') }}
              </button>
              <button
                class="selector-btn"
                :class="{ active: selectedBenchmark === 'BTC' }"
                @click="selectedBenchmark = 'BTC'"
              >
                vs BTC
              </button>
              <button
                class="selector-btn"
                :class="{ active: selectedBenchmark === 'ETH' }"
                @click="selectedBenchmark = 'ETH'"
              >
                vs ETH
              </button>
            </div>
            <!-- 时间范围选择 -->
            <div class="selector-group">
              <button
                v-for="range in timeRanges"
                :key="range"
                class="selector-btn"
                :class="{ active: selectedTimeRange === range }"
                @click="selectedTimeRange = range"
              >
                {{ range }}
              </button>
            </div>
          </div>
        </div>
        <div class="chart-container">
          <Line
            :data="lineChartData"
            :options="lineChartOptions"
            :key="selectedPricingCurrency + selectedTimeRange + selectedBenchmark"
          />
        </div>
      </div>

      <!-- 资产分布饼图 -->
      <div v-if="dashboardStore.widgetConfig.distribution" class="glass-card dist-panel">
        <div class="panel-header">
          <h3>{{ t('dashboard.assetDistribution') }}</h3>
        </div>
        <div class="doughnut-wrapper">
          <div class="doughnut-chart">
            <Doughnut :data="doughnutChartData" :options="doughnutChartOptions" />
          </div>
          <div class="dist-legend">
            <div
              v-for="cat in assetCategories"
              :key="cat.id"
              class="legend-row"
            >
              <span
                class="legend-dot"
                :style="{ background: cat.fixedColor || themeStore.currentTheme.colors[cat.colorKey] }"
              />
              <span class="legend-name">{{ t(cat.labelKey) }}</span>
              <span class="legend-pct font-mono">
                {{ totalAssets ? ((cat.totalValue / totalAssets) * 100).toFixed(1) : 0 }}%
              </span>
              <span class="legend-val font-mono">
                {{ currencySymbol }}{{ formatNumber(cat.totalValue, 0) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 健康评分 + 目标追踪 -->
    <section v-if="dashboardStore.widgetConfig.healthScore || dashboardStore.widgetConfig.goals" class="insight-row">
      <HealthScoreCard v-if="dashboardStore.widgetConfig.healthScore" />
      <div v-if="dashboardStore.widgetConfig.goals" class="goals-panel">
        <div class="goals-header">
          <h3>{{ t('goals.title') }}</h3>
          <button class="btn btn-ghost btn-sm" @click="showAddGoal = true">
            {{ t('goals.addGoal') }}
          </button>
        </div>
        <div v-if="displayedGoals.length > 0" class="goals-list">
          <GoalCard
            v-for="goal in displayedGoals"
            :key="goal.id"
            :goal="goal"
            @delete="goalStore.removeGoal($event)"
          />
        </div>
        <div v-else class="goals-empty">
          {{ t('goals.empty') }}
        </div>
      </div>
    </section>

    <!-- DeFi 仓位概览 -->
    <DeFiOverview v-if="dashboardStore.widgetConfig.defiOverview" />

    <!-- NFT 资产概览 -->
    <NFTOverview v-if="dashboardStore.widgetConfig.nftOverview" />

    <!-- 费用分析 -->
    <FeeAnalytics v-if="dashboardStore.widgetConfig.feeAnalytics" />

    <!-- 策略面板 -->
    <StrategyPanel v-if="dashboardStore.widgetConfig.strategyPanel" />

    <!-- 添加目标对话框 -->
    <AddGoalDialog :visible="showAddGoal" @close="showAddGoal = false" />

    <!-- 持仓明细表格 -->
    <section v-if="dashboardStore.widgetConfig.holdings" class="glass-card holdings-panel">
      <div class="panel-header">
        <h3>{{ t('dashboard.holdingsDetail') || '持仓明细' }}</h3>
        <div class="panel-controls">
          <!-- 视图切换：平铺 / 分组 -->
          <div class="selector-group">
            <button
              class="selector-btn"
              :class="{ active: holdingsViewMode === 'flat' }"
              @click="holdingsViewMode = 'flat'"
              :title="t('dashboard.viewFlat')"
            >
              <PhListBullets :size="14" />
            </button>
            <button
              class="selector-btn"
              :class="{ active: holdingsViewMode === 'grouped' }"
              @click="holdingsViewMode = 'grouped'"
              :title="t('dashboard.viewGrouped')"
            >
              <PhSquaresFour :size="14" />
            </button>
          </div>
          <div class="search-box">
            <PhMagnifyingGlass :size="14" />
            <input
              type="text"
              v-model="searchQuery"
              :placeholder="t('common.search') || '搜索...'"
            />
          </div>
        </div>
      </div>

      <table class="table holdings-table">
        <thead>
          <tr>
            <th class="col-expand"></th>
            <th class="col-asset">{{ t('dashboard.asset') }}</th>
            <th class="col-source">{{ t('dashboard.sources') }}</th>
            <th class="col-sortable" @click="toggleSort('price')">
              {{ t('dashboard.price') }}
              <PhSortAscending v-if="sortField === 'price' && sortOrder === 'asc'" :size="12" />
              <PhSortDescending v-else-if="sortField === 'price' && sortOrder === 'desc'" :size="12" />
            </th>
            <th class="col-sortable" @click="toggleSort('change24h')">
              {{ t('dashboard.change') }}
              <PhSortAscending v-if="sortField === 'change24h' && sortOrder === 'asc'" :size="12" />
              <PhSortDescending v-else-if="sortField === 'change24h' && sortOrder === 'desc'" :size="12" />
            </th>
            <th class="col-sortable" @click="toggleSort('balance')">
              {{ t('dashboard.balance') }}
              <PhSortAscending v-if="sortField === 'balance' && sortOrder === 'asc'" :size="12" />
              <PhSortDescending v-else-if="sortField === 'balance' && sortOrder === 'desc'" :size="12" />
            </th>
            <th class="col-sortable" @click="toggleSort('value')">
              {{ t('dashboard.value') }}
              <PhSortAscending v-if="sortField === 'value' && sortOrder === 'asc'" :size="12" />
              <PhSortDescending v-else-if="sortField === 'value' && sortOrder === 'desc'" :size="12" />
            </th>
            <th class="col-pct">{{ t('dashboard.percentage') }}</th>
          </tr>
        </thead>
        <!-- 平铺视图 -->
        <tbody v-if="holdingsViewMode === 'flat'">
          <template v-for="holding in paginatedHoldings" :key="holding.symbol">
            <tr class="holding-row" @click="holding.sources.length > 1 && toggleAssetExpand(holding.symbol)">
              <td class="col-expand">
                <button
                  v-if="holding.sources.length > 1"
                  class="expand-btn"
                  @click.stop="toggleAssetExpand(holding.symbol)"
                >
                  <PhCaretRight
                    :size="14"
                    :class="{ 'rotated-down': expandedAssets.has(holding.symbol) }"
                  />
                </button>
              </td>
              <td class="col-asset">
                <div class="asset-cell">
                  <div class="asset-icon">
                    <img v-if="holding.icon" :src="holding.icon" :alt="holding.symbol" />
                    <span v-else class="icon-fallback">{{ holding.symbol.slice(0, 2) }}</span>
                  </div>
                  <div class="asset-meta">
                    <span class="asset-name">{{ holding.name }}</span>
                    <span class="asset-symbol">{{ holding.symbol }}</span>
                  </div>
                </div>
              </td>
              <td class="col-source">
                <span v-if="holding.sources.length === 1" class="source-tag">
                  {{ holding.sources[0].source }}
                </span>
                <span v-else class="sources-count">
                  {{ t('dashboard.mergedAssets', { count: holding.sources.length }) }}
                </span>
              </td>
              <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.price) }}</td>
              <td
                class="font-mono"
                :class="holding.change24h >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ holding.change24h >= 0 ? '+' : '' }}{{ holding.change24h }}%
              </td>
              <td class="font-mono">{{ formatNumber(holding.balance, 4) }}</td>
              <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.value) }}</td>
              <td class="font-mono col-pct-val">{{ holdingPercent(holding.value) }}%</td>
            </tr>

            <template v-if="expandedAssets.has(holding.symbol) && holding.sources.length > 1">
              <tr
                v-for="(source, idx) in holding.sources"
                :key="`${holding.symbol}-${idx}`"
                class="source-detail-row"
              >
                <td class="col-expand"></td>
                <td class="col-asset">
                  <div class="source-detail-indent">
                    <span class="source-detail-label">└ {{ source.source }}</span>
                  </div>
                </td>
                <td class="col-source">
                  <span class="source-type-tag">{{ source.sourceType }}</span>
                </td>
                <td class="font-mono">-</td>
                <td class="font-mono">-</td>
                <td class="font-mono">{{ formatNumber(source.balance, 4) }}</td>
                <td class="font-mono">{{ currencySymbol }}{{ formatNumber(source.value) }}</td>
                <td class="font-mono col-pct-val">{{ holdingPercent(source.value) }}%</td>
              </tr>
            </template>
          </template>

          <tr v-if="paginatedHoldings.length === 0">
            <td colspan="8" class="empty-row">
              {{ searchQuery ? t('common.noResults') : t('common.noData') }}
            </td>
          </tr>
        </tbody>

        <!-- 智能分组视图 -->
        <tbody v-else>
          <template v-for="group in groupedHoldings" :key="group.id">
            <!-- 分组头 -->
            <tr class="group-header-row" @click="toggleGroup(group.id)">
              <td class="col-expand">
                <button class="expand-btn" @click.stop="toggleGroup(group.id)">
                  <PhCaretRight
                    :size="14"
                    :class="{ 'rotated-down': expandedGroups.has(group.id) }"
                  />
                </button>
              </td>
              <td colspan="5">
                <div class="group-header-cell">
                  <span class="group-name">{{ t(group.labelKey) }}</span>
                  <span class="group-count">{{ group.items.length }}</span>
                </div>
              </td>
              <td class="font-mono group-total">{{ currencySymbol }}{{ formatNumber(group.totalValue) }}</td>
              <td class="font-mono col-pct-val group-total">{{ group.percentage }}%</td>
            </tr>

            <!-- 分组内的资产 -->
            <template v-if="expandedGroups.has(group.id)">
              <tr
                v-for="holding in group.items"
                :key="`${group.id}-${holding.symbol}`"
                class="holding-row grouped-item-row"
              >
                <td class="col-expand"></td>
                <td class="col-asset">
                  <div class="asset-cell">
                    <div class="asset-icon">
                      <img v-if="holding.icon" :src="holding.icon" :alt="holding.symbol" />
                      <span v-else class="icon-fallback">{{ holding.symbol.slice(0, 2) }}</span>
                    </div>
                    <div class="asset-meta">
                      <span class="asset-name">{{ holding.name }}</span>
                      <span class="asset-symbol">{{ holding.symbol }}</span>
                    </div>
                  </div>
                </td>
                <td class="col-source">
                  <span v-if="holding.sources.length === 1" class="source-tag">
                    {{ holding.sources[0].source }}
                  </span>
                  <span v-else class="sources-count">
                    {{ t('dashboard.mergedAssets', { count: holding.sources.length }) }}
                  </span>
                </td>
                <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.price) }}</td>
                <td
                  class="font-mono"
                  :class="holding.change24h >= 0 ? 'change-positive' : 'change-negative'"
                >
                  {{ holding.change24h >= 0 ? '+' : '' }}{{ holding.change24h }}%
                </td>
                <td class="font-mono">{{ formatNumber(holding.balance, 4) }}</td>
                <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.value) }}</td>
                <td class="font-mono col-pct-val">{{ holdingPercent(holding.value) }}%</td>
              </tr>
            </template>
          </template>

          <tr v-if="groupedHoldings.length === 0">
            <td colspan="8" class="empty-row">
              {{ searchQuery ? t('common.noResults') : t('common.noData') }}
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 分页控件 -->
      <div v-if="filteredHoldings.length > 0" class="pagination">
        <div class="pagination-info">
          {{ t('common.showingItems', paginationInfo) }}
        </div>

        <div class="pagination-controls">
          <button
            class="page-btn"
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
          >
            {{ t('common.prevPage') }}
          </button>

          <div class="page-numbers">
            <button
              v-for="page in totalPages"
              :key="page"
              v-show="Math.abs(page - currentPage) <= 2 || page === 1 || page === totalPages"
              class="page-num"
              :class="{ active: page === currentPage }"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
          </div>

          <button
            class="page-btn"
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
          >
            {{ t('common.nextPage') }}
          </button>
        </div>

        <div class="page-size-selector">
          <span>{{ t('common.itemsPerPage') }}:</span>
          <select v-model="itemsPerPage" @change="changePageSize($event.target.value)">
            <option v-for="size in availablePageSizes" :key="size" :value="size">
              {{ size }}
            </option>
          </select>
        </div>
      </div>
    </section>

    <!-- 资产详情抽屉 -->
    <AssetDetailDrawer
      :visible="showAssetDrawer"
      :asset="selectedAsset"
      @close="closeAssetDrawer"
    />

    <!-- 价格预警对话框 -->
    <PriceAlertDialog
      :visible="showPriceAlert"
      @close="showPriceAlert = false"
    />

    <!-- 分享卡片 -->
    <ShareCard
      :visible="showShareCard"
      @close="showShareCard = false"
    />

    <!-- Widget 配置面板 -->
    <DashboardCustomizer
      :visible="showCustomizer"
      @close="showCustomizer = false"
    />

    <!-- 首次使用引导 -->
    <OnboardingWizard v-if="showOnboarding" @complete="() => {}" />
  </div>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  max-width: 1400px;
}

/* ========== 顶部标题 ========== */
.dashboard-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.gear-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.gear-btn:hover {
  background: var(--color-bg-elevated);
  color: var(--color-accent-primary);
}

/* 空仪表盘引导 */
.empty-dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-2xl);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  text-align: center;
}

.empty-dashboard p {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

/* ========== 摘要栏 ========== */
.summary-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-xl);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.summary-left {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.total-row {
  display: flex;
  align-items: baseline;
  gap: var(--gap-sm);
}

.total-row .stat-value-hero {
  font-size: 1.5rem;
}

.change-badge {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 0.8125rem;
  font-weight: 500;
}

.change-abs {
  color: var(--color-text-secondary);
  font-weight: 400;
}

.change-period {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-left: 2px;
}

/* 今日盈亏卡片 */
.pnl-card {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  min-width: 120px;
}

.pnl-card.pnl-positive {
  border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
}

.pnl-card.pnl-negative {
  border-color: color-mix(in srgb, var(--color-error) 30%, transparent);
}

.pnl-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.pnl-value {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.pnl-positive .pnl-value {
  color: var(--color-success);
}

.pnl-negative .pnl-value {
  color: var(--color-error);
}

.pnl-percent {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 0.6875rem;
  font-weight: 500;
}

.summary-right {
  display: flex;
  align-items: center;
  gap: var(--gap-lg);
}

.alert-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast), border-color var(--transition-fast);
}

.alert-btn:hover {
  background: var(--color-bg-elevated);
  color: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
}

.category-chip {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-xs) var(--gap-sm);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
}

.chip-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.chip-label {
  font-size: 0.6875rem;
  color: var(--color-text-secondary);
}

.chip-value {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

/* ========== 健康评分 + 目标追踪 ========== */
.insight-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-lg);
}

.goals-panel {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
}

.goals-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-md);
}

.goals-header h3 {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.btn-sm {
  padding: 4px 10px;
  font-size: 0.6875rem;
}

.goals-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  flex: 1;
}

.goals-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* ========== 图表区 ========== */
.charts-row {
  display: grid;
  grid-template-columns: 3fr 2fr;
  gap: var(--gap-lg);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  margin-bottom: var(--gap-md);
  flex-wrap: wrap;
}

.panel-header h3 {
  font-size: 0.8125rem;
  font-weight: 600;
}

.panel-controls {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.trend-panel {
  padding: var(--gap-lg);
}

.chart-container {
  height: 220px;
}

/* 选择器按钮组 */
.selector-group {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.selector-btn {
  padding: 4px 8px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
  white-space: nowrap;
}

.selector-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.selector-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* 饼图面板 */
.dist-panel {
  padding: var(--gap-lg);
}

.doughnut-wrapper {
  display: flex;
  align-items: center;
  gap: var(--gap-lg);
  height: 220px;
}

.doughnut-chart {
  width: 140px;
  height: 140px;
  flex-shrink: 0;
}

.dist-legend {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.legend-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-name {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  flex: 1;
}

.legend-pct {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  min-width: 40px;
  text-align: right;
}

.legend-val {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-primary);
  min-width: 80px;
  text-align: right;
}

/* ========== 持仓表格 ========== */
.holdings-panel {
  padding: var(--gap-lg);
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 4px 8px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
}

.search-box:focus-within {
  border-color: var(--color-accent-primary);
}

.search-box input {
  background: none;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-size: 0.75rem;
  width: 140px;
}

.search-box input::placeholder {
  color: var(--color-text-muted);
}

.holdings-table {
  margin-top: var(--gap-sm);
}

.holdings-table th {
  white-space: nowrap;
}

.holdings-table tbody tr {
  cursor: pointer;
  height: 40px;
}

.col-sortable {
  cursor: pointer;
  user-select: none;
}

.col-sortable:hover {
  color: var(--color-text-primary);
}

/* 资产列 */
.asset-cell {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.asset-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-elevated);
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.asset-icon img {
  width: 20px;
  height: 20px;
  object-fit: contain;
}

.icon-fallback {
  font-size: 0.625rem;
  font-weight: 600;
  color: var(--color-text-muted);
}

.asset-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.asset-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.asset-symbol {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 来源标签 */
.source-tag {
  font-size: 0.6875rem;
  color: var(--color-text-secondary);
  padding: 1px 6px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
}

.col-pct {
  text-align: right;
}

.col-pct-val {
  text-align: right;
  color: var(--color-text-secondary);
}

.empty-row {
  text-align: center;
  color: var(--color-text-muted);
  padding: var(--gap-xl) !important;
}

/* 展开按钮列 */
.col-expand {
  width: 32px;
  text-align: center;
  padding: 0 !important;
}

.expand-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: background var(--transition-fast), color var(--transition-fast);
}

.expand-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.expand-btn svg {
  transition: transform var(--transition-fast);
}

.rotated-down {
  transform: rotate(90deg);
}

.holding-row {
  cursor: pointer;
}

/* 来源数量标签 */
.sources-count {
  font-size: 0.6875rem;
  color: var(--color-accent-primary);
  padding: 2px 6px;
  background: rgba(75, 131, 240, 0.1);
  border-radius: var(--radius-xs);
  font-weight: 500;
}

/* 来源明细行 */
.source-detail-row {
  background: var(--color-bg-tertiary);
  font-size: 0.75rem;
}

.source-detail-row td {
  padding-top: 4px !important;
  padding-bottom: 4px !important;
}

.source-detail-indent {
  padding-left: var(--gap-md);
  color: var(--color-text-secondary);
}

.source-detail-label {
  font-size: 0.75rem;
}

.source-type-tag {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  padding: 1px 4px;
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xs);
  text-transform: uppercase;
}

/* ========== 智能分组视图 ========== */
.group-header-row {
  background: var(--color-bg-tertiary);
  cursor: pointer;
  border-bottom: 1px solid var(--color-border);
}

.group-header-row:hover {
  background: var(--color-bg-elevated);
}

.group-header-row td {
  padding: var(--gap-sm) var(--gap-md) !important;
}

.group-header-cell {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.group-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.group-count {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  padding: 1px 6px;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-xs);
}

.group-total {
  font-weight: 600;
  color: var(--color-text-primary) !important;
}

.grouped-item-row td:first-child {
  border-left: 2px solid var(--color-accent-primary);
}

/* ========== 分页 ========== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-lg);
  margin-top: var(--gap-lg);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
  flex-wrap: wrap;
}

.pagination-info {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.page-btn {
  padding: 4px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.page-btn:hover:not(:disabled) {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  border-color: var(--color-border-hover);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  align-items: center;
  gap: 2px;
}

.page-num {
  min-width: 28px;
  height: 28px;
  padding: 0 6px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.page-num:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.page-num.active {
  background: var(--color-accent-primary);
  color: #fff;
  border-color: var(--color-accent-primary);
}

.page-size-selector {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.page-size-selector select {
  padding: 4px 8px;
  font-size: 0.75rem;
  color: var(--color-text-primary);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.page-size-selector select:hover {
  border-color: var(--color-border-hover);
}

.page-size-selector select:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

/* ========== 响应式 ========== */
@media (max-width: 768px) {
  .summary-bar {
    flex-direction: column;
    align-items: flex-start;
    padding: var(--gap-md);
  }

  .total-row {
    flex-wrap: wrap;
  }

  .total-row .stat-value-hero {
    font-size: 1.25rem;
  }

  .summary-right {
    flex-wrap: wrap;
    gap: var(--gap-sm);
    width: 100%;
  }

  .category-chip {
    flex: 1;
    min-width: 0;
  }

  .pnl-card {
    width: 100%;
  }

  .insight-row {
    grid-template-columns: 1fr;
  }

  .charts-row {
    grid-template-columns: 1fr;
  }

  .doughnut-wrapper {
    flex-direction: column;
    height: auto;
  }

  .panel-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .panel-controls {
    flex-wrap: wrap;
    width: 100%;
  }

  .selector-btn {
    min-height: 36px;
    padding: 6px 10px;
  }

  /* 移动端按钮触摸区域 */
  .gear-btn,
  .alert-btn,
  .btn-sm {
    min-width: 44px;
    min-height: 44px;
  }

  .sort-btn {
    min-width: 44px;
    min-height: 44px;
  }

  /* 移动端隐藏部分列 */
  .col-source,
  .holdings-table th:nth-child(2),
  .holdings-table td:nth-child(3),
  .holdings-table th:nth-child(3) {
    display: none;
  }

  /* 持仓行更紧凑 */
  .holdings-table td,
  .holdings-table th {
    padding: var(--gap-sm);
    font-size: 0.75rem;
  }
}

@media (max-width: 480px) {
  .change-abs {
    display: none;
  }

  .change-period {
    display: none;
  }
}
</style>
