<script setup>
/**
 * 数据分析页面
 * 4 个分析面板：资产配置趋势、收益率曲线、资产集中度、平台分布
 */
import { ref, computed, onMounted, watch } from 'vue'
import { Line, Doughnut, Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import FeeAnalytics from '../components/FeeAnalytics.vue'
import StrategyPanel from '../components/StrategyPanel.vue'
import BenchmarkPanel from '../components/BenchmarkPanel.vue'
import PortfolioShareDialog from '../components/PortfolioShareDialog.vue'
import { useAssetStore } from '../stores/assetStore'
import { useThemeStore } from '../stores/themeStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

// 注册 Chart.js 组件
ChartJS.register(
  CategoryScale, LinearScale, PointElement, LineElement,
  BarElement, ArcElement, Title, Tooltip, Legend, Filler
)

const assetStore = useAssetStore()
const themeStore = useThemeStore()
const { formatNumber, currencySymbol } = useFormatters()
const { t } = useI18n()

// 时间范围
const selectedTimeRange = ref('30D')
const timeRanges = ['7D', '30D', '90D', '1Y']

// 初始化
onMounted(async () => {
  if (!assetStore.historyData) {
    await assetStore.loadHistory(selectedTimeRange.value)
  }
  if (!assetStore.summary) {
    await assetStore.loadSummary()
  }
})

watch(selectedTimeRange, async (val) => {
  await assetStore.loadHistory(val)
})

const hasData = computed(() => assetStore.hasData)

// 配置分享对话框
const showShareDialog = ref(false)
const colors = computed(() => themeStore.currentTheme.colors)

// ===== 面板 1：资产配置趋势（堆叠面积图） =====
const allocationChartData = computed(() => {
  const hd = assetStore.historyData
  if (!hd || !hd.labels) return { labels: [], datasets: [] }

  const labels = hd.labels
  const totalValues = hd.values || []

  // 模拟各类资产占比随时间变化（基于当前分类比例 + 随机波动）
  const dist = assetStore.platformDistribution
  const cexRatio = (dist[0]?.percentage || 60) / 100
  const chainRatio = (dist[1]?.percentage || 30) / 100
  const manualRatio = (dist[2]?.percentage || 10) / 100

  return {
    labels,
    datasets: [
      {
        label: t('dashboard.cexAssets'),
        data: totalValues.map((v, i) => v * (cexRatio + (Math.sin(i * 0.3) * 0.05))),
        backgroundColor: `${colors.value.accentPrimary}66`,
        borderColor: colors.value.accentPrimary,
        borderWidth: 1,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: t('dashboard.blockchainAssets'),
        data: totalValues.map((v, i) => v * (chainRatio + (Math.cos(i * 0.3) * 0.03))),
        backgroundColor: `${colors.value.accentSecondary}66`,
        borderColor: colors.value.accentSecondary,
        borderWidth: 1,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: t('dashboard.manualAssets'),
        data: totalValues.map(() => totalValues[0] * manualRatio),
        backgroundColor: `${colors.value.accentTertiary}66`,
        borderColor: colors.value.accentTertiary,
        borderWidth: 1,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
      },
    ]
  }
})

const allocationChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      align: 'end',
      labels: { color: colors.value.textSecondary, usePointStyle: true, pointStyle: 'circle', padding: 12, font: { size: 11 } }
    },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
    }
  },
  scales: {
    x: { display: true, grid: { display: false }, ticks: { color: colors.value.textMuted, maxTicksLimit: 6, font: { size: 10 } } },
    y: {
      stacked: true,
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: { color: colors.value.textSecondary, font: { size: 10 }, callback: (v) => `${currencySymbol.value}${formatNumber(v, 0)}` }
    }
  }
}))

// ===== 面板 2：收益率曲线 =====
const returnChartData = computed(() => {
  const hd = assetStore.historyData
  if (!hd || !hd.values || hd.values.length < 2) return { labels: [], datasets: [] }

  const labels = hd.labels
  const base = hd.values[0] || 1
  // 我的收益率
  const myReturns = hd.values.map(v => ((v - base) / base) * 100)
  // BTC/ETH 基准（模拟）
  const btcChange = assetStore.benchmarkComparison?.btc?.change || 0
  const ethChange = assetStore.benchmarkComparison?.eth?.change || 0

  const btcReturns = labels.map((_, i) => {
    const progress = i / Math.max(labels.length - 1, 1)
    return btcChange * progress
  })
  const ethReturns = labels.map((_, i) => {
    const progress = i / Math.max(labels.length - 1, 1)
    return ethChange * progress
  })

  return {
    labels,
    datasets: [
      {
        label: t('analytics.myReturn'),
        data: myReturns,
        borderColor: colors.value.accentPrimary,
        backgroundColor: `${colors.value.accentPrimary}15`,
        borderWidth: 2,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
        pointHoverRadius: 3,
      },
      {
        label: 'BTC',
        data: btcReturns,
        borderColor: '#F7931A',
        backgroundColor: 'transparent',
        borderWidth: 1.5,
        borderDash: [6, 3],
        fill: false,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'ETH',
        data: ethReturns,
        borderColor: '#627EEA',
        backgroundColor: 'transparent',
        borderWidth: 1.5,
        borderDash: [6, 3],
        fill: false,
        tension: 0.4,
        pointRadius: 0,
      }
    ]
  }
})

const returnChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      align: 'end',
      labels: { color: colors.value.textSecondary, usePointStyle: true, pointStyle: 'line', padding: 12, font: { size: 11 } }
    },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      callbacks: {
        label: (ctx) => `${ctx.dataset.label}: ${ctx.raw >= 0 ? '+' : ''}${ctx.raw.toFixed(2)}%`
      }
    }
  },
  scales: {
    x: { display: true, grid: { display: false }, ticks: { color: colors.value.textMuted, maxTicksLimit: 6, font: { size: 10 } } },
    y: {
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: { color: colors.value.textSecondary, font: { size: 10 }, callback: (v) => `${v >= 0 ? '+' : ''}${v.toFixed(1)}%` }
    }
  }
}))

// ===== 面板 3：资产集中度（饼图 + HHI） =====
const concentrationData = computed(() => assetStore.assetConcentration)

const concentrationChartData = computed(() => {
  const top5 = concentrationData.value.top5
  const chartColors = [
    colors.value.accentPrimary,
    colors.value.accentSecondary,
    colors.value.accentTertiary || '#f97316',
    '#8b5cf6',
    '#06b6d4',
    colors.value.textMuted,
  ]
  return {
    labels: top5.map(a => a.symbol),
    datasets: [{
      data: top5.map(a => a.value),
      backgroundColor: top5.map((_, i) => chartColors[i % chartColors.length]),
      borderWidth: 0,
      cutout: '65%',
    }]
  }
})

const concentrationChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: 'rgba(0,0,0,0.8)',
      callbacks: {
        label: (ctx) => {
          const total = ctx.dataset.data.reduce((a, b) => a + b, 0)
          const pct = ((ctx.raw / total) * 100).toFixed(1)
          return `${ctx.label}: ${pct}%`
        }
      }
    }
  }
}

// HHI 等级
const hhiLevel = computed(() => {
  const hhi = concentrationData.value.hhi
  if (hhi < 1500) return { level: 'low', color: 'var(--color-success)' }
  if (hhi < 2500) return { level: 'medium', color: 'var(--color-warning)' }
  return { level: 'high', color: 'var(--color-error)' }
})

// ===== 面板 4：平台分布（水平柱状图） =====
const platformData = computed(() => assetStore.platformDistribution)

const platformChartData = computed(() => {
  const platformLabels = [t('dashboard.cexAssets'), t('dashboard.blockchainAssets'), t('dashboard.manualAssets')]
  const platformColors = [colors.value.accentPrimary, colors.value.accentSecondary, colors.value.accentTertiary || '#f97316']
  return {
    labels: platformLabels,
    datasets: [{
      data: platformData.value.map(p => p.value),
      backgroundColor: platformColors,
      borderWidth: 0,
      borderRadius: 3,
      barThickness: 24,
    }]
  }
})

const platformChartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  indexAxis: 'y',
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      callbacks: {
        label: (ctx) => `${currencySymbol.value}${formatNumber(ctx.raw)}`
      }
    }
  },
  scales: {
    x: {
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: { color: colors.value.textSecondary, font: { size: 10 }, callback: (v) => `${currencySymbol.value}${formatNumber(v, 0)}` }
    },
    y: {
      grid: { display: false },
      ticks: { color: colors.value.textPrimary, font: { size: 11, weight: 500 } }
    }
  }
}))
</script>

<template>
  <div class="analytics-page">
    <!-- 时间范围选择 -->
    <div class="analytics-header">
      <h2 class="page-title">{{ t('nav.analytics') }}</h2>
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

    <template v-if="hasData">
      <!-- 摘要指标 -->
      <div class="metrics-row">
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.totalReturn') }}</span>
          <span class="metric-value font-mono" :class="assetStore.returnRate.totalReturn >= 0 ? 'positive' : 'negative'">
            {{ assetStore.returnRate.totalReturn >= 0 ? '+' : '' }}{{ assetStore.returnRate.totalReturn }}%
          </span>
        </div>
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.annualizedReturn') }}</span>
          <span class="metric-value font-mono" :class="assetStore.returnRate.annualizedReturn >= 0 ? 'positive' : 'negative'">
            {{ assetStore.returnRate.annualizedReturn >= 0 ? '+' : '' }}{{ assetStore.returnRate.annualizedReturn }}%
          </span>
        </div>
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.hhiIndex') }}</span>
          <span class="metric-value font-mono" :style="{ color: hhiLevel.color }">
            {{ concentrationData.hhi }}
          </span>
          <span class="metric-hint" :style="{ color: hhiLevel.color }">
            {{ t(`analytics.hhi.${hhiLevel.level}`) }}
          </span>
        </div>
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.assetCount') }}</span>
          <span class="metric-value font-mono">{{ concentrationData.top5.length }}</span>
        </div>
      </div>

      <!-- 上排：资产配置趋势 + 收益率曲线 -->
      <div class="panels-row">
        <div class="panel">
          <h3>{{ t('analytics.allocationTrend') }}</h3>
          <div class="chart-container">
            <Line :data="allocationChartData" :options="allocationChartOptions" :key="'alloc-' + selectedTimeRange" />
          </div>
        </div>
        <div class="panel">
          <h3>{{ t('analytics.returnCurve') }}</h3>
          <div class="chart-container">
            <Line :data="returnChartData" :options="returnChartOptions" :key="'return-' + selectedTimeRange" />
          </div>
        </div>
      </div>

      <!-- 下排：资产集中度 + 平台分布 -->
      <div class="panels-row">
        <div class="panel">
          <h3>{{ t('analytics.concentration') }}</h3>
          <div class="concentration-layout">
            <div class="concentration-chart">
              <Doughnut :data="concentrationChartData" :options="concentrationChartOptions" />
            </div>
            <div class="concentration-legend">
              <div
                v-for="(asset, i) in concentrationData.top5"
                :key="asset.symbol"
                class="conc-legend-row"
              >
                <span class="conc-dot" :style="{ background: concentrationChartData.datasets[0].backgroundColor[i] }" />
                <span class="conc-name">{{ asset.symbol }}</span>
                <span class="conc-pct font-mono">{{ asset.percentage.toFixed(1) }}%</span>
                <span class="conc-val font-mono">{{ currencySymbol }}{{ formatNumber(asset.value, 0) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="panel">
          <h3>{{ t('analytics.platformDist') }}</h3>
          <div class="chart-container chart-bar">
            <Bar :data="platformChartData" :options="platformChartOptions" />
          </div>
          <!-- 平台占比文字 -->
          <div class="platform-stats">
            <div v-for="p in platformData" :key="p.id" class="platform-stat">
              <span class="platform-label">{{ t(`dashboard.${p.id === 'blockchain' ? 'blockchainAssets' : p.id === 'cex' ? 'cexAssets' : 'manualAssets'}`) }}</span>
              <span class="platform-pct font-mono">{{ p.percentage.toFixed(1) }}%</span>
            </div>
          </div>
        </div>
      </div>
      <!-- 面板 5：费用分析 -->
      <FeeAnalytics />

      <!-- 面板 6：策略面板 -->
      <StrategyPanel />

      <!-- 面板 7：基准对比 + 导出报告 -->
      <div class="panels-row">
        <BenchmarkPanel />
        <div class="panel share-entry-panel">
          <h3>{{ t('social.sharePortfolio') }}</h3>
          <p class="share-desc">{{ t('social.shareDesc') }}</p>
          <button class="btn btn-primary btn-share" @click="showShareDialog = true">
            {{ t('social.generateShareImage') }}
          </button>
        </div>
      </div>

    </template>

    <!-- 无数据状态 -->
    <div v-else class="empty-state">
      <p>{{ t('analytics.noData') }}</p>
    </div>

    <!-- 配置分享对话框 -->
    <PortfolioShareDialog
      :visible="showShareDialog"
      @close="showShareDialog = false"
    />
  </div>
</template>

<style scoped>
.analytics-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  max-width: 1400px;
}

/* 头部 */
.analytics-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
}

.page-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.selector-group {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.selector-btn {
  padding: 4px 10px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.selector-btn:hover {
  color: var(--color-text-primary);
}

.selector-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* 摘要指标 */
.metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--gap-md);
}

.metric-card {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.metric-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.metric-value {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.metric-value.positive {
  color: var(--color-success);
}

.metric-value.negative {
  color: var(--color-error);
}

.metric-hint {
  font-size: 0.625rem;
  font-weight: 500;
}

/* 面板行 */
.panels-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-lg);
}

.panel {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.panel h3 {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.chart-container {
  height: 220px;
}

.chart-bar {
  height: 140px;
}

/* 集中度布局 */
.concentration-layout {
  display: flex;
  gap: var(--gap-lg);
  align-items: center;
}

.concentration-chart {
  width: 140px;
  height: 140px;
  flex-shrink: 0;
}

.concentration-legend {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.conc-legend-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.conc-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.conc-name {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  flex: 1;
}

.conc-pct {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  min-width: 40px;
  text-align: right;
}

.conc-val {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-primary);
  min-width: 80px;
  text-align: right;
}

/* 平台统计 */
.platform-stats {
  display: flex;
  justify-content: space-between;
  gap: var(--gap-md);
  margin-top: var(--gap-md);
  padding-top: var(--gap-sm);
  border-top: 1px solid var(--color-border);
}

.platform-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.platform-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.platform-pct {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 空状态 */
.empty-state {
  padding: var(--gap-2xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 分享入口 */
.share-entry-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: var(--gap-md);
  text-align: center;
}

.share-desc {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  max-width: 280px;
}

.btn-share {
  padding: 8px 20px;
  font-size: 0.8125rem;
  font-weight: 500;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  background: var(--color-accent-primary);
  color: #fff;
  transition: opacity var(--transition-fast);
}

.btn-share:hover {
  opacity: 0.9;
}

/* 响应式 */
@media (max-width: 768px) {
  .metrics-row {
    grid-template-columns: repeat(2, 1fr);
  }

  .panels-row {
    grid-template-columns: 1fr;
  }

  .concentration-layout {
    flex-direction: column;
  }
}
</style>
