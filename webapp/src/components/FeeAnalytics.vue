<script setup>
/**
 * 费用分析面板
 * 展示费用摘要、构成环形图、月度趋势折线图、智能建议
 */
import { ref, computed, onMounted, watch } from 'vue'
import { Doughnut, Line } from 'vue-chartjs'
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
import { PhCaretUp, PhCaretDown, PhLightbulb } from '@phosphor-icons/vue'
import { feeService } from '../api/feeService.js'
import { useFormatters } from '../composables/useFormatters'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

ChartJS.register(
  CategoryScale, LinearScale, PointElement, LineElement,
  ArcElement, Title, Tooltip, Legend, Filler
)

const { formatNumber } = useFormatters()
const themeStore = useThemeStore()
const { t } = useI18n()

// 时间范围
const selectedRange = ref('30D')
const rangeOptions = ['7D', '30D', '90D']

// 数据
const feeData = ref(null)
const isLoading = ref(false)

// 加载数据
const loadData = async () => {
  isLoading.value = true
  try {
    feeData.value = await feeService.getFeeSummary(selectedRange.value)
  } catch (err) {
    console.error('加载费用数据失败:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)
watch(selectedRange, loadData)

const colors = computed(() => themeStore.currentTheme.colors)

// 费用构成环形图
const doughnutData = computed(() => {
  if (!feeData.value) return { labels: [], datasets: [] }
  const { breakdown } = feeData.value
  return {
    labels: [t('fee.cexTradeFee'), t('fee.gasFee'), t('fee.withdrawFee')],
    datasets: [{
      data: [breakdown.cexTradeFee, breakdown.gasFee, breakdown.withdrawFee],
      backgroundColor: [
        colors.value.accentPrimary,
        colors.value.accentSecondary,
        colors.value.accentTertiary || '#F59E0B',
      ],
      borderWidth: 0,
      hoverOffset: 4,
    }]
  }
})

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      callbacks: {
        label: (ctx) => `${ctx.label}: $${ctx.raw.toFixed(2)}`
      }
    }
  }
}

// 月度趋势折线图
const trendData = computed(() => {
  if (!feeData.value || !feeData.value.monthlyTrend) return { labels: [], datasets: [] }
  const trend = feeData.value.monthlyTrend
  return {
    labels: trend.map(m => m.month.slice(5)), // 只取 MM
    datasets: [{
      label: t('fee.totalFee'),
      data: trend.map(m => m.total),
      borderColor: colors.value.accentPrimary,
      backgroundColor: `${colors.value.accentPrimary}15`,
      borderWidth: 1.5,
      fill: true,
      tension: 0.4,
      pointRadius: 3,
      pointHoverRadius: 5,
    }]
  }
})

const trendOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      callbacks: {
        label: (ctx) => `$${ctx.raw.toFixed(2)}`
      }
    }
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: colors.value.textMuted, font: { size: 10 } }
    },
    y: {
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => `$${v}`
      }
    }
  }
}))

// 费用构成百分比
const breakdownPercentages = computed(() => {
  if (!feeData.value) return []
  const { breakdown, total } = feeData.value
  return [
    { label: t('fee.cexTradeFee'), value: breakdown.cexTradeFee, pct: (breakdown.cexTradeFee / total * 100).toFixed(1), color: colors.value.accentPrimary },
    { label: t('fee.gasFee'), value: breakdown.gasFee, pct: (breakdown.gasFee / total * 100).toFixed(1), color: colors.value.accentSecondary },
    { label: t('fee.withdrawFee'), value: breakdown.withdrawFee, pct: (breakdown.withdrawFee / total * 100).toFixed(1), color: colors.value.accentTertiary || '#F59E0B' },
  ]
})
</script>

<template>
  <section class="fee-analytics">
    <!-- 标题栏 + 时间范围选择 -->
    <div class="panel-header">
      <h3 class="panel-title">{{ t('fee.title') }}</h3>
      <div class="range-selector">
        <button
          v-for="r in rangeOptions"
          :key="r"
          class="range-btn"
          :class="{ active: selectedRange === r }"
          @click="selectedRange = r"
        >
          {{ r }}
        </button>
      </div>
    </div>

    <template v-if="feeData">
      <!-- 费用摘要 -->
      <div class="fee-summary">
        <div class="summary-main">
          <span class="summary-label">{{ t('fee.totalSpent') }}</span>
          <span class="summary-value font-mono">${{ formatNumber(feeData.total, 2) }}</span>
        </div>
        <div class="summary-change" :class="feeData.changePercent > 0 ? 'negative' : 'positive'">
          <PhCaretUp v-if="feeData.changePercent > 0" :size="12" />
          <PhCaretDown v-else :size="12" />
          <span class="font-mono">{{ Math.abs(feeData.changePercent) }}%</span>
          <span class="change-label">{{ t('fee.vsLastMonth') }}</span>
        </div>
      </div>

      <!-- 图表区域 -->
      <div class="charts-row">
        <!-- 构成环形图 -->
        <div class="chart-block">
          <h4 class="chart-label">{{ t('fee.composition') }}</h4>
          <div class="doughnut-wrapper">
            <Doughnut :data="doughnutData" :options="doughnutOptions" />
          </div>
          <div class="breakdown-list">
            <div v-for="item in breakdownPercentages" :key="item.label" class="breakdown-item">
              <span class="dot" :style="{ background: item.color }"></span>
              <span class="breakdown-label">{{ item.label }}</span>
              <span class="breakdown-value font-mono">${{ formatNumber(item.value, 2) }}</span>
              <span class="breakdown-pct font-mono">{{ item.pct }}%</span>
            </div>
          </div>
        </div>

        <!-- 月度趋势折线图 -->
        <div class="chart-block">
          <h4 class="chart-label">{{ t('fee.monthlyTrend') }}</h4>
          <div class="line-wrapper">
            <Line :data="trendData" :options="trendOptions" :key="selectedRange" />
          </div>
        </div>
      </div>

      <!-- 智能建议 -->
      <div v-if="feeData.suggestions && feeData.suggestions.length > 0" class="suggestions">
        <h4 class="suggestions-title">
          <PhLightbulb :size="14" weight="fill" />
          {{ t('fee.suggestions') }}
        </h4>
        <div class="suggestion-list">
          <div
            v-for="s in feeData.suggestions"
            :key="s.id"
            class="suggestion-item"
          >
            <div class="suggestion-content">
              <span class="suggestion-title">{{ t(s.titleKey) }}</span>
              <span class="suggestion-desc">{{ t(s.descKey) }}</span>
            </div>
            <span v-if="s.savingEstimate" class="saving-badge font-mono positive">
              ~${{ s.savingEstimate }}/{{ t('fee.month') }}
            </span>
          </div>
        </div>
      </div>
    </template>

    <!-- 加载中 -->
    <div v-else-if="isLoading" class="loading-state">
      {{ t('common.loading') }}
    </div>
  </section>
</template>

<style scoped>
.fee-analytics {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 标题栏 */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.range-selector {
  display: flex;
  gap: 4px;
}

.range-btn {
  padding: 3px 10px;
  font-size: 0.6875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.range-btn:hover {
  border-color: var(--color-accent-primary);
}

.range-btn.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

/* 费用摘要 */
.fee-summary {
  display: flex;
  align-items: flex-end;
  gap: var(--gap-lg);
}

.summary-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.summary-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.summary-value {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.summary-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 500;
}

.summary-change.positive {
  color: var(--color-success);
}

.summary-change.negative {
  color: var(--color-error);
}

.change-label {
  color: var(--color-text-muted);
  font-weight: 400;
}

/* 图表区域 */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-lg);
}

.chart-block {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.chart-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.doughnut-wrapper {
  height: 160px;
  display: flex;
  justify-content: center;
}

.line-wrapper {
  height: 160px;
}

/* 构成明细 */
.breakdown-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.breakdown-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  font-size: 0.6875rem;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.breakdown-label {
  flex: 1;
  color: var(--color-text-secondary);
}

.breakdown-value {
  color: var(--color-text-primary);
  font-weight: 500;
}

.breakdown-pct {
  color: var(--color-text-muted);
  width: 40px;
  text-align: right;
}

/* 智能建议 */
.suggestions {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

.suggestions-title {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-warning);
}

.suggestion-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.suggestion-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.suggestion-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.suggestion-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.saving-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  padding: 2px 8px;
  background: rgba(16, 185, 129, 0.1);
  border-radius: var(--radius-xs);
  white-space: nowrap;
}

/* 加载 */
.loading-state {
  text-align: center;
  padding: var(--gap-xl);
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

@media (max-width: 768px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}
</style>
