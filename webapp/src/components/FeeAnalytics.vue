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
import { PhCaretUp, PhCaretDown, PhLightbulb, PhCheckCircle, PhArrowRight, PhInfo, PhArrowUUpRight, PhGlobe, PhTimer, PhWallet } from '@phosphor-icons/vue'
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

// 建议状态管理
const dismissedSuggestions = ref(new Set())

const dismissSuggestion = (id) => {
  dismissedSuggestions.value.add(id)
}

// 建议优先级颜色
const getPriorityColor = (priority) => {
  switch (priority) {
    case 'high': return 'rgba(239, 68, 68, 0.15)'
    case 'medium': return 'rgba(245, 158, 11, 0.15)'
    case 'low': return 'rgba(59, 130, 246, 0.15)'
    default: return 'rgba(107, 114, 128, 0.15)'
  }
}

const getPriorityBorderColor = (priority) => {
  switch (priority) {
    case 'high': return 'rgba(239, 68, 68, 0.5)'
    case 'medium': return 'rgba(245, 158, 11, 0.5)'
    case 'low': return 'rgba(59, 130, 246, 0.5)'
    default: return 'rgba(107, 114, 128, 0.3)'
  }
}

const getPriorityLabel = (priority) => {
  switch (priority) {
    case 'high': return '高优先级'
    case 'medium': return '中优先级'
    case 'low': return '低优先级'
    default: return ''
  }
}

// 根据建议类型获取对应图标
const getIconForSuggestion = (type) => {
  const iconMap = {
    'gas': PhGlobe,
    'timing': PhTimer,
    'consolidate': PhWallet,
    'exchange': PhArrowUUpRight
  }
  return iconMap[type] || PhLightbulb
}
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
      <div v-if="feeData.suggestions && feeData.suggestions.filter(s => !dismissedSuggestions.has(s.id)).length > 0" class="suggestions-section">
        <div class="suggestions-header">
          <div class="suggestions-title-row">
            <div class="suggestions-icon">
              <PhLightbulb :size="18" weight="fill" />
            </div>
            <div class="suggestions-header-text">
              <h4 class="suggestions-title">{{ t('fee.suggestions') }}</h4>
              <span class="suggestions-subtitle">基于您的交易模式分析</span>
            </div>
          </div>
          <div class="suggestions-stats">
            <div class="stat-box">
              <span class="stat-value">${{ formatNumber(feeData.suggestions.reduce((sum, s) => sum + (s.savingEstimate || 0), 0), 2) }}</span>
              <span class="stat-label">预估月节省</span>
            </div>
          </div>
        </div>

        <div class="suggestions-grid">
          <div
            v-for="s in feeData.suggestions.filter(s => !dismissedSuggestions.has(s.id))"
            :key="s.id"
            class="suggestion-card"
            :style="{
              background: getPriorityColor(s.priority),
              borderColor: getPriorityBorderColor(s.priority)
            }"
          >
            <!-- 卡片头部：优先级和关闭按钮 -->
            <div class="suggestion-card-header">
              <span class="priority-badge" :class="`priority-${s.priority}`">
                {{ getPriorityLabel(s.priority) }}
              </span>
              <button class="dismiss-btn" @click="dismissSuggestion(s.id)" title="忽略此建议">
                <PhInfo :size="14" />
              </button>
            </div>

            <!-- 卡片内容：图标 + 文字 -->
            <div class="suggestion-card-body">
              <div class="suggestion-icon-box">
                <component :is="getIconForSuggestion(s.type)" :size="24" weight="duotone" />
              </div>
              <div class="suggestion-text">
                <h5 class="suggestion-card-title">{{ t(s.titleKey) }}</h5>
                <p class="suggestion-card-desc">{{ t(s.descKey) }}</p>
              </div>
            </div>

            <!-- 卡片底部：节省金额 + 操作按钮 -->
            <div class="suggestion-card-footer">
              <div class="saving-highlight">
                <span class="saving-label">可节省</span>
                <span class="saving-value font-mono">~${{ s.savingEstimate }}</span>
                <span class="saving-period">/月</span>
              </div>
              <button class="action-btn">
                <span>查看详情</span>
                <PhArrowRight :size="14" weight="bold" />
              </button>
            </div>

            <!-- 进度条装饰 -->
            <div class="suggestion-progress-bar">
              <div class="suggestion-progress-fill" :style="{ width: `${s.impactScore}%` }"></div>
            </div>
          </div>
        </div>

        <!-- 已忽略建议提示 -->
        <div v-if="dismissedSuggestions.size > 0" class="dismissed-hint">
          <button class="restore-hint" @click="dismissedSuggestions.clear()">
            <PhCheckCircle :size="14" />
            <span>已忽略 {{ dismissedSuggestions.size }} 条建议，点击恢复</span>
          </button>
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

/* 智能建议 - 重新设计 */
.suggestions-section {
  margin-top: var(--gap-xl);
  padding-top: var(--gap-lg);
  border-top: 1px solid var(--color-border);
}

/* 建议头部 */
.suggestions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-md);
  flex-wrap: wrap;
  gap: var(--gap-md);
}

.suggestions-title-row {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
}

.suggestions-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, var(--color-accent-primary) 0%, color-mix(in srgb, var(--color-accent-primary) 80%, var(--color-accent-secondary) 100%) 100%);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 12px color-mix(in srgb, var(--color-accent-primary) 30%, transparent);
}

.suggestions-header-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.suggestions-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.suggestions-subtitle {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.suggestions-stats {
  display: flex;
  gap: var(--gap-md);
}

.stat-box {
  text-align: center;
  padding: var(--gap-sm) var(--gap-lg);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
}

.stat-value {
  display: block;
  font-size: 1rem;
  font-weight: 700;
  color: var(--color-success);
  font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
}

.stat-label {
  display: block;
  font-size: 0.625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-top: 2px;
}

/* 建议卡片网格 */
.suggestions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--gap-md);
}

.suggestion-card {
  position: relative;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-md);
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  overflow: hidden;
  transition: all 0.25s ease;
}

.suggestion-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--color-accent-primary), var(--color-accent-secondary));
}

.suggestion-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

/* 卡片头部 */
.suggestion-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.priority-badge {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: var(--radius-xs);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.priority-badge.priority-high {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.priority-badge.priority-medium {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.priority-badge.priority-low {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.dismiss-btn {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius-xs);
  transition: all 0.15s ease;
}

.dismiss-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-tertiary);
}

/* 卡片主体 */
.suggestion-card-body {
  display: flex;
  gap: var(--gap-md);
  align-items: flex-start;
}

.suggestion-icon-box {
  width: 48px;
  height: 48px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--color-accent-primary);
  border: 1px solid var(--color-border);
}

.suggestion-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.suggestion-card-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  line-height: 1.3;
}

.suggestion-card-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin: 0;
  line-height: 1.4;
}

/* 卡片底部 */
.suggestion-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--gap-sm);
  border-top: 1px solid color-mix(in srgb, var(--color-border) 50%, transparent);
}

.saving-highlight {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.saving-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.saving-value {
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--color-success);
}

.saving-period {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-accent-primary);
  border: none;
  border-radius: var(--radius-sm);
  color: white;
  font-size: 0.6875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

.action-btn:hover {
  background: color-mix(in srgb, var(--color-accent-primary) 85%, black);
  transform: translateX(2px);
}

/* 进度条装饰 */
.suggestion-progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--color-bg-tertiary);
}

.suggestion-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-success), var(--color-accent-primary));
  border-radius: 1px;
  transition: width 0.5s ease;
}

/* 已忽略提示 */
.dismissed-hint {
  margin-top: var(--gap-md);
  display: flex;
  justify-content: center;
}

.restore-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-bg-tertiary);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  font-size: 0.6875rem;
  cursor: pointer;
  transition: all 0.15s ease;
}

.restore-hint:hover {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: white;
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

@media (max-width: 768px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
  .suggestions-grid {
    grid-template-columns: 1fr;
  }
  .suggestions-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .suggestions-stats {
    width: 100%;
    justify-content: center;
  }
  .suggestion-card-footer {
    flex-direction: column;
    gap: var(--gap-sm);
    align-items: stretch;
  }
  .action-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
