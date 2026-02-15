<script setup>
/**
 * Dashboard 页面 - 资产总览 (Bento Grid 布局重构)
 * 参考设计：现代金融仪表盘，Grid 布局，模块化卡片
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
  PhSquaresFour,
  PhWallet,
  PhTrendUp,
  PhChartPieSlice,
  PhTarget
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
    borderWidth: 2,
    fill: true,
    tension: 0.4,
    pointRadius: 0,
    pointHoverRadius: 4,
    yAxisID: 'y'
  }]

  if (showExchangeRate.value) {
    datasets.push({
      label: `${selectedPricingCurrency.value}/USDC`,
      data: rateData,
      borderColor: themeStore.currentTheme.colors.accentSecondary,
      backgroundColor: 'transparent',
      borderWidth: 1.5,
      borderDash: [4, 4],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 4,
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
      borderWidth: 2,
      borderDash: [6, 3],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 4,
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
        padding: 10,
        titleFont: { size: 12 },
        bodyFont: { size: 12 },
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
          lineWidth: 0.5,
          borderDash: [4, 4]
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
  <div ref="dashboardRef" class="dashboard-container">
    <!-- 下拉刷新指示器 -->
    <PullToRefresh :pull-distance="pullDistance" :is-refreshing="pullRefreshing" />

    <!-- 顶部导航栏 -->
    <header class="dashboard-navbar">
      <div class="navbar-left">
        <h2 class="page-title">{{ t('dashboard.title') }}</h2>
      </div>
      <div class="navbar-right">
        <button class="icon-btn" @click="showCustomizer = true" :title="t('widgets.customize')">
          <PhGear :size="18" weight="bold" />
        </button>
      </div>
    </header>

    <!-- 空状态引导 -->
    <div v-if="dashboardStore.enabledCount === 0" class="empty-dashboard-state">
      <PhWallet :size="48" weight="duotone" class="empty-icon" />
      <p>{{ t('widgets.emptyHint') }}</p>
      <button class="btn btn-primary" @click="showCustomizer = true">{{ t('widgets.customize') }}</button>
    </div>

    <!-- 核心 Grid 布局 -->
    <div v-else class="bento-grid">
      
      <!-- 1. 资产总览卡片 (Hero Card) -->
      <div v-if="dashboardStore.widgetConfig.assetSummary" class="grid-item hero-card">
        <div class="hero-content">
          <div class="hero-main">
            <span class="label-text">{{ t('dashboard.totalAssets') }}</span>
            <div class="balance-row">
              <h1 class="total-balance font-mono">
                {{ currencySymbol }}{{ formatNumber(totalAssets) }}
              </h1>
              <div class="pnl-badge" :class="totalChange24h >= 0 ? 'positive' : 'negative'">
                <PhCaretUp v-if="totalChange24h >= 0" :size="14" weight="bold" />
                <PhCaretDown v-else :size="14" weight="bold" />
                <span class="font-mono">{{ formatPercent(totalChange24h) }}</span>
                <span class="pnl-value font-mono">({{ currencySymbol }}{{ formatNumber(totalChangeValue) }})</span>
              </div>
            </div>
            <div class="hero-actions">
              <button class="action-btn" @click="showShareCard = true" :title="t('share.title')">
                <PhShareNetwork :size="16" weight="bold" />
                <span>{{ t('share.title') }}</span>
              </button>
              <button class="action-btn" @click="showPriceAlert = true" :title="t('priceAlert.title')">
                <PhBell :size="16" weight="bold" />
                <span>{{ t('priceAlert.title') }}</span>
              </button>
            </div>
          </div>
          
          <div class="hero-stats">
            <div class="stat-item" v-for="cat in assetCategories" :key="cat.id">
              <div class="stat-icon" :style="{ backgroundColor: cat.fixedColor || themeStore.currentTheme.colors[cat.colorKey] }"></div>
              <div class="stat-info">
                <span class="stat-label">{{ t(cat.labelKey) }}</span>
                <span class="stat-num font-mono">{{ currencySymbol }}{{ formatNumber(cat.totalValue, 0) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 2. 趋势图表 (占大版面) -->
      <div v-if="dashboardStore.widgetConfig.trend" class="grid-item trend-card">
        <div class="card-header">
          <div class="header-title">
            <PhTrendUp :size="18" weight="duotone" class="header-icon" />
            <h3>{{ t('dashboard.assetTrend') }}</h3>
          </div>
          <div class="header-controls">
            <!-- 计价货币切换 -->
            <div class="toggle-group">
               <button 
                 v-for="curr in currencies" 
                 :key="curr.code"
                 class="toggle-item"
                 :class="{ active: selectedPricingCurrency === curr.code }"
                 @click="setPricingCurrency(curr.code)"
               >
                 {{ curr.code }}
               </button>
            </div>
            <!-- 时间范围切换 -->
            <div class="toggle-group">
               <button 
                 v-for="range in timeRanges" 
                 :key="range"
                 class="toggle-item"
                 :class="{ active: selectedTimeRange === range }"
                 @click="selectedTimeRange = range"
               >
                 {{ range }}
               </button>
            </div>
          </div>
        </div>
        <div class="chart-wrapper">
          <Line
            :data="lineChartData"
            :options="lineChartOptions"
            :key="selectedPricingCurrency + selectedTimeRange + selectedBenchmark"
          />
        </div>
      </div>

      <!-- 3. 资产分布 (Doughnut) -->
      <div v-if="dashboardStore.widgetConfig.distribution" class="grid-item dist-card">
        <div class="card-header">
          <div class="header-title">
            <PhChartPieSlice :size="18" weight="duotone" class="header-icon" />
            <h3>{{ t('dashboard.assetDistribution') }}</h3>
          </div>
        </div>
        <div class="dist-content">
          <div class="doughnut-box">
            <Doughnut :data="doughnutChartData" :options="doughnutChartOptions" />
          </div>
          <div class="dist-list">
            <div v-for="cat in assetCategories" :key="cat.id" class="dist-item">
              <span class="dot" :style="{ background: cat.fixedColor || themeStore.currentTheme.colors[cat.colorKey] }"></span>
              <span class="name">{{ t(cat.labelKey) }}</span>
              <span class="pct font-mono">{{ totalAssets ? ((cat.totalValue / totalAssets) * 100).toFixed(1) : 0 }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. 健康评分 & 目标 (Insight Group) -->
      <div v-if="dashboardStore.widgetConfig.healthScore" class="grid-item health-card">
        <HealthScoreCard />
      </div>

      <div v-if="dashboardStore.widgetConfig.goals" class="grid-item goals-card">
        <div class="card-header">
          <div class="header-title">
            <PhTarget :size="18" weight="duotone" class="header-icon" />
            <h3>{{ t('goals.title') }}</h3>
          </div>
          <button class="icon-btn-sm" @click="showAddGoal = true">
            <span class="plus-icon">+</span>
          </button>
        </div>
        <div class="goals-content">
          <div v-if="displayedGoals.length > 0" class="goals-stack">
            <GoalCard
              v-for="goal in displayedGoals"
              :key="goal.id"
              :goal="goal"
              @delete="goalStore.removeGoal($event)"
            />
          </div>
          <div v-else class="empty-state-sm">
            {{ t('goals.empty') }}
          </div>
        </div>
      </div>

      <!-- 5. 其他组件概览 -->
      <div v-if="dashboardStore.widgetConfig.defiOverview" class="grid-item defi-card">
        <DeFiOverview />
      </div>
      <div v-if="dashboardStore.widgetConfig.nftOverview" class="grid-item nft-card">
        <NFTOverview />
      </div>

      <!-- 6. 持仓明细 (Full Width) -->
      <div v-if="dashboardStore.widgetConfig.holdings" class="grid-item holdings-card">
        <div class="card-header">
          <div class="header-title">
            <PhListBullets :size="18" weight="duotone" class="header-icon" />
            <h3>{{ t('dashboard.holdingsDetail') || '持仓明细' }}</h3>
          </div>
          <div class="header-actions">
            <!-- 视图切换 -->
            <div class="view-toggles">
              <button 
                class="view-btn" 
                :class="{ active: holdingsViewMode === 'flat' }"
                @click="holdingsViewMode = 'flat'"
              >
                <PhListBullets :size="14" />
              </button>
              <button 
                class="view-btn" 
                :class="{ active: holdingsViewMode === 'grouped' }"
                @click="holdingsViewMode = 'grouped'"
              >
                <PhSquaresFour :size="14" />
              </button>
            </div>
            <!-- 搜索 -->
            <div class="search-input">
              <PhMagnifyingGlass :size="14" class="search-icon" />
              <input 
                type="text" 
                v-model="searchQuery" 
                :placeholder="t('common.search')" 
              />
            </div>
          </div>
        </div>

        <div class="table-container">
          <table class="modern-table">
            <thead>
              <tr>
                <th class="th-expand"></th>
                <th class="th-asset">{{ t('dashboard.asset') }}</th>
                <th class="th-source">{{ t('dashboard.sources') }}</th>
                <th class="th-sortable" @click="toggleSort('price')">
                  {{ t('dashboard.price') }}
                  <span class="sort-indicator" v-if="sortField === 'price'">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th class="th-sortable" @click="toggleSort('change24h')">
                  {{ t('dashboard.change') }}
                  <span class="sort-indicator" v-if="sortField === 'change24h'">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th class="th-right">{{ t('dashboard.balance') }}</th>
                <th class="th-right th-sortable" @click="toggleSort('value')">
                  {{ t('dashboard.value') }}
                  <span class="sort-indicator" v-if="sortField === 'value'">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th class="th-right">{{ t('dashboard.percentage') }}</th>
              </tr>
            </thead>
            
            <tbody v-if="holdingsViewMode === 'flat'">
              <template v-for="holding in paginatedHoldings" :key="holding.symbol">
                <tr class="tr-hover" @click="holding.sources.length > 1 && toggleAssetExpand(holding.symbol)">
                  <td class="td-expand">
                    <button v-if="holding.sources.length > 1" class="expand-toggle">
                      <PhCaretRight :size="12" :class="{ 'rotate-90': expandedAssets.has(holding.symbol) }" />
                    </button>
                  </td>
                  <td>
                    <div class="asset-info">
                      <div class="asset-logo">
                        <img v-if="holding.icon" :src="holding.icon" />
                        <span v-else>{{ holding.symbol.slice(0, 1) }}</span>
                      </div>
                      <div class="asset-text">
                        <div class="symbol">{{ holding.symbol }}</div>
                        <div class="name">{{ holding.name }}</div>
                      </div>
                    </div>
                  </td>
                  <td>
                    <span v-if="holding.sources.length === 1" class="tag">{{ holding.sources[0].source }}</span>
                    <span v-else class="tag multi">{{ holding.sources.length }} sources</span>
                  </td>
                  <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.price) }}</td>
                  <td class="font-mono" :class="holding.change24h >= 0 ? 'text-green' : 'text-red'">
                    {{ holding.change24h >= 0 ? '+' : '' }}{{ holding.change24h }}%
                  </td>
                  <td class="font-mono text-right">{{ formatNumber(holding.balance, 4) }}</td>
                  <td class="font-mono text-right bold">{{ currencySymbol }}{{ formatNumber(holding.value) }}</td>
                  <td class="font-mono text-right text-muted">{{ holdingPercent(holding.value) }}%</td>
                </tr>
                <!-- 展开详情 -->
                <template v-if="expandedAssets.has(holding.symbol) && holding.sources.length > 1">
                  <tr v-for="(source, idx) in holding.sources" :key="idx" class="tr-detail">
                    <td></td>
                    <td class="td-indent">└ {{ source.source }}</td>
                    <td><span class="tag-xs">{{ source.sourceType }}</span></td>
                    <td></td>
                    <td></td>
                    <td class="font-mono text-right">{{ formatNumber(source.balance, 4) }}</td>
                    <td class="font-mono text-right">{{ currencySymbol }}{{ formatNumber(source.value) }}</td>
                    <td></td>
                  </tr>
                </template>
              </template>
            </tbody>
            
            <!-- 分组视图 (简化版逻辑复用) -->
            <tbody v-else>
              <template v-for="group in groupedHoldings" :key="group.id">
                <tr class="tr-group" @click="toggleGroup(group.id)">
                  <td class="td-expand">
                    <PhCaretRight :size="12" :class="{ 'rotate-90': expandedGroups.has(group.id) }" />
                  </td>
                  <td colspan="5" class="group-title">
                    {{ t(group.labelKey) }} <span class="badge-count">{{ group.items.length }}</span>
                  </td>
                  <td class="font-mono text-right bold">{{ currencySymbol }}{{ formatNumber(group.totalValue) }}</td>
                  <td class="font-mono text-right">{{ group.percentage }}%</td>
                </tr>
                <template v-if="expandedGroups.has(group.id)">
                  <tr v-for="holding in group.items" :key="holding.symbol" class="tr-hover">
                     <!-- 与平铺视图相同的行结构，略 -->
                     <td></td>
                     <td>
                        <div class="asset-info">
                          <div class="asset-logo sm">
                            <img v-if="holding.icon" :src="holding.icon" />
                            <span v-else>{{ holding.symbol.slice(0, 1) }}</span>
                          </div>
                          <span class="symbol">{{ holding.symbol }}</span>
                        </div>
                     </td>
                     <td>
                       <span v-if="holding.sources.length === 1" class="tag">{{ holding.sources[0].source }}</span>
                       <span v-else class="tag multi">{{ holding.sources.length }}</span>
                     </td>
                     <td class="font-mono">{{ currencySymbol }}{{ formatNumber(holding.price) }}</td>
                     <td class="font-mono" :class="holding.change24h >= 0 ? 'text-green' : 'text-red'">
                       {{ holding.change24h }}%
                     </td>
                     <td class="font-mono text-right">{{ formatNumber(holding.balance, 4) }}</td>
                     <td class="font-mono text-right">{{ currencySymbol }}{{ formatNumber(holding.value) }}</td>
                     <td class="font-mono text-right text-muted">{{ holdingPercent(holding.value) }}%</td>
                  </tr>
                </template>
              </template>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div class="pagination-bar">
          <span class="page-info">{{ t('common.showingItems', paginationInfo) }}</span>
          <div class="page-ctrl">
            <button :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">{{ t('common.prevPage') }}</button>
            <button :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">{{ t('common.nextPage') }}</button>
          </div>
        </div>
      </div>

    </div>

    <!-- 弹窗组件 -->
    <AssetDetailDrawer :visible="showAssetDrawer" :asset="selectedAsset" @close="closeAssetDrawer" />
    <PriceAlertDialog :visible="showPriceAlert" @close="showPriceAlert = false" />
    <ShareCard :visible="showShareCard" @close="showShareCard = false" />
    <DashboardCustomizer :visible="showCustomizer" @close="showCustomizer = false" />
    <AddGoalDialog :visible="showAddGoal" @close="showAddGoal = false" />
    <OnboardingWizard v-if="showOnboarding" @complete="() => {}" />
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 1600px;
  margin: 0 auto;
  padding-bottom: var(--gap-2xl);
}

/* 顶部导航 */
.dashboard-navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-lg);
}
.page-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text-primary);
}
.icon-btn {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}
.icon-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-accent-primary);
}

/* Bento Grid */
.bento-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: var(--gap-lg);
  grid-auto-flow: dense;
}

.grid-item {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 1. Hero Card */
.hero-card {
  grid-column: span 8;
  padding: var(--gap-xl);
  background: linear-gradient(135deg, var(--color-bg-secondary) 0%, var(--color-bg-elevated) 100%);
}

@media (min-width: 1600px) {
  .hero-card { grid-column: span 12; display: flex; flex-direction: row; justify-content: space-between; align-items: center; }
  .hero-content { display: flex; width: 100%; justify-content: space-between; align-items: center; }
  .hero-main { min-width: 400px; }
  .hero-stats { margin-top: 0 !important; }
}

.label-text {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.balance-row {
  display: flex;
  align-items: baseline;
  gap: var(--gap-md);
  margin: var(--gap-sm) 0 var(--gap-lg);
}
.total-balance {
  font-size: 2.5rem;
  line-height: 1;
  color: var(--color-text-primary);
}
.pnl-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  background: rgba(0,0,0,0.2);
  font-size: 0.875rem;
}
.pnl-badge.positive { color: var(--color-success); background: rgba(46, 189, 133, 0.1); }
.pnl-badge.negative { color: var(--color-error); background: rgba(226, 92, 92, 0.1); }
.pnl-value { opacity: 0.8; font-size: 0.8em; margin-left: 4px; }

.hero-actions {
  display: flex;
  gap: var(--gap-md);
}
.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}
.action-btn:hover { border-color: var(--color-accent-primary); }

.hero-stats {
  margin-top: var(--gap-xl);
  display: flex;
  gap: var(--gap-xl);
}
.stat-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}
.stat-icon {
  width: 4px;
  height: 24px;
  border-radius: 2px;
}
.stat-info {
  display: flex;
  flex-direction: column;
}
.stat-label { font-size: 0.75rem; color: var(--color-text-muted); }
.stat-num { font-size: 0.9em; font-weight: 500; }

/* 2. Trend Card */
.trend-card {
  grid-column: span 8;
  min-height: 320px;
  padding: var(--gap-lg);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-lg);
}
.header-title {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}
.header-icon { color: var(--color-accent-primary); }
.header-title h3 { font-size: 1rem; font-weight: 600; }

.toggle-group {
  display: flex;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-md);
}
.toggle-item {
  padding: 4px 12px;
  font-size: 0.75rem;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
}
.toggle-item.active {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
  box-shadow: var(--shadow-sm);
}
.chart-wrapper {
  flex: 1;
  position: relative;
  width: 100%;
}

/* 3. Distribution Card */
.dist-card {
  grid-column: span 4;
  padding: var(--gap-lg);
}
.dist-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-lg);
}
.doughnut-box {
  width: 160px;
  height: 160px;
}
.dist-list {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}
.dist-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 0.8125rem;
}
.dot { width: 8px; height: 8px; border-radius: 50%; margin-right: 8px; }
.name { flex: 1; color: var(--color-text-secondary); }
.pct { font-weight: 500; }

/* 4. Insight Cards */
.health-card { grid-column: span 4; padding: 0; border: none; background: transparent; }
.goals-card { grid-column: span 4; padding: var(--gap-lg); }
.defi-card, .nft-card { grid-column: span 6; }

.icon-btn-sm {
  width: 24px; height: 24px;
  border-radius: 50%;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.icon-btn-sm:hover { border-color: var(--color-accent-primary); color: var(--color-accent-primary); }

.goals-stack { display: flex; flex-direction: column; gap: var(--gap-sm); }
.empty-state-sm { text-align: center; padding: var(--gap-lg); color: var(--color-text-muted); font-size: 0.875rem; }

/* 5. Holdings Card */
.holdings-card {
  grid-column: span 12;
  padding: var(--gap-lg);
}
.header-actions {
  display: flex;
  gap: var(--gap-md);
  align-items: center;
}
.view-toggles {
  display: flex;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-md);
  padding: 2px;
}
.view-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
  cursor: pointer;
}
.view-btn.active {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}
.search-input {
  display: flex;
  align-items: center;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 0 8px;
  width: 200px;
}
.search-input input {
  background: transparent;
  border: none;
  padding: 6px;
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  width: 100%;
}
.search-input:focus-within { border-color: var(--color-accent-primary); }
.search-icon { color: var(--color-text-muted); }

/* Table Styling */
.table-container { overflow-x: auto; margin-top: var(--gap-md); }
.modern-table { width: 100%; border-collapse: separate; border-spacing: 0; }
.modern-table th {
  text-align: left;
  padding: 12px;
  font-size: 0.75rem;
  color: var(--color-text-muted);
  font-weight: 500;
  border-bottom: 1px solid var(--color-border);
}
.modern-table td {
  padding: 12px;
  font-size: 0.875rem;
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
}
.th-right, .text-right { text-align: right; }
.th-sortable { cursor: pointer; user-select: none; }
.th-sortable:hover { color: var(--color-text-primary); }

.tr-hover:hover { background: var(--color-bg-tertiary); cursor: pointer; }

/* Asset Cell */
.asset-info { display: flex; align-items: center; gap: 10px; }
.asset-logo {
  width: 32px; height: 32px;
  background: var(--color-bg-elevated);
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  overflow: hidden;
}
.asset-logo img { width: 20px; height: 20px; }
.asset-logo span { font-size: 0.75rem; font-weight: 600; color: var(--color-text-secondary); }
.asset-logo.sm { width: 24px; height: 24px; }
.asset-text { display: flex; flex-direction: column; }
.asset-text .symbol { font-weight: 500; }
.asset-text .name { font-size: 0.75rem; color: var(--color-text-muted); }

.tag {
  font-size: 0.75rem; padding: 2px 6px;
  background: var(--color-bg-tertiary);
  border-radius: 4px; color: var(--color-text-secondary);
}
.tag.multi { background: rgba(75, 131, 240, 0.1); color: var(--color-accent-primary); }
.text-green { color: var(--color-success); }
.text-red { color: var(--color-error); }
.text-muted { color: var(--color-text-muted); }
.bold { font-weight: 600; }

.expand-toggle {
  background: transparent; border: none; color: var(--color-text-muted);
  cursor: pointer; padding: 4px;
}
.rotate-90 { transform: rotate(90deg); transition: transform 0.2s; }

/* Detail Row */
.tr-detail { background: var(--color-bg-tertiary); font-size: 0.8125rem; }
.td-indent { padding-left: 24px !important; color: var(--color-text-secondary); }

/* Group Row */
.tr-group { background: var(--color-bg-tertiary); cursor: pointer; font-weight: 600; }
.badge-count { font-size: 0.75rem; background: var(--color-bg-secondary); padding: 1px 6px; border-radius: 4px; margin-left: 6px; font-weight: normal; color: var(--color-text-muted); }

/* Pagination */
.pagination-bar {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: var(--gap-lg);
}
.page-info { font-size: 0.8125rem; color: var(--color-text-secondary); }
.page-ctrl button {
  padding: 6px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  cursor: pointer;
  margin-left: 8px;
}
.page-ctrl button:disabled { opacity: 0.5; cursor: not-allowed; }

/* Responsive */
@media (max-width: 1200px) {
  .hero-card { grid-column: span 12; }
  .trend-card { grid-column: span 12; }
  .dist-card { grid-column: span 12; display: flex; flex-direction: row; justify-content: space-around; }
  .doughnut-box { margin: 0; }
  .health-card, .goals-card { grid-column: span 6; }
}

@media (max-width: 768px) {
  .bento-grid { display: flex; flex-direction: column; }
  .dashboard-navbar { flex-direction: column; align-items: flex-start; gap: var(--gap-md); }
  .navbar-right { position: absolute; top: var(--gap-lg); right: var(--gap-lg); }
  .hero-card { padding: var(--gap-lg); }
  .total-balance { font-size: 2rem; }
  .hero-stats { flex-direction: column; gap: var(--gap-md); }
  .chart-wrapper { height: 250px; }
  .dist-card { flex-direction: column; }
}
</style>
