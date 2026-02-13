<script setup>
/**
 * 年度投资报告组件
 * 多页翻页式展示（类似 PPT），支持键盘和触摸翻页
 * 第 1 页：年度收益总览
 * 第 2 页：月度收益曲线
 * 第 3 页：资产配置变化
 * 第 4 页：关键里程碑
 * 第 5 页：投资风格标签 + 总结
 */
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler
} from 'chart.js'
import {
  PhCaretLeft,
  PhCaretRight,
  PhTrophy,
  PhRocketLaunch,
  PhFlag,
  PhStar,
  PhX,
  PhShareNetwork
} from '@phosphor-icons/vue'
import { annualReportService } from '../api/annualReportService.js'
import { useThemeStore } from '../stores/themeStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Filler)

const props = defineProps({
  visible: { type: Boolean, default: false },
  year: { type: Number, default: 2025 },
})

const emit = defineEmits(['close', 'share'])

const themeStore = useThemeStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

const report = ref(null)
const currentPage = ref(0)
const totalPages = 5
const isLoading = ref(false)
const colors = computed(() => themeStore.currentTheme.colors)

// 投资风格标签映射
const styleLabels = {
  steady: 'annualReport.styleSteady',
  aggressive: 'annualReport.styleAggressive',
  conservative: 'annualReport.styleConservative',
  balanced: 'annualReport.styleBalanced',
}

// 里程碑图标映射
const milestoneIcons = {
  start: PhFlag,
  rocket: PhRocketLaunch,
  milestone: PhStar,
  trophy: PhTrophy,
}

// 月份名称
const monthNames = computed(() => [
  t('annualReport.jan'), t('annualReport.feb'), t('annualReport.mar'),
  t('annualReport.apr'), t('annualReport.may'), t('annualReport.jun'),
  t('annualReport.jul'), t('annualReport.aug'), t('annualReport.sep'),
  t('annualReport.oct'), t('annualReport.nov'), t('annualReport.dec'),
])

// 翻页
const goNext = () => {
  if (currentPage.value < totalPages - 1) currentPage.value++
}
const goPrev = () => {
  if (currentPage.value > 0) currentPage.value--
}

// 键盘翻页
const handleKeydown = (e) => {
  if (!props.visible) return
  if (e.key === 'ArrowRight' || e.key === 'ArrowDown') goNext()
  if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') goPrev()
  if (e.key === 'Escape') emit('close')
}

// 触摸翻页
let touchStartX = 0
const handleTouchStart = (e) => { touchStartX = e.touches[0].clientX }
const handleTouchEnd = (e) => {
  const diff = touchStartX - e.changedTouches[0].clientX
  if (Math.abs(diff) > 50) {
    diff > 0 ? goNext() : goPrev()
  }
}

// 月度收益图表
const monthlyChartData = computed(() => {
  if (!report.value) return { labels: [], datasets: [] }
  const data = report.value.monthlyReturns
  return {
    labels: data.map(m => monthNames.value[m.month - 1]),
    datasets: [{
      label: t('annualReport.monthlyReturn'),
      data: data.map(m => m.return),
      borderColor: colors.value.accentPrimary,
      backgroundColor: (ctx) => {
        if (!ctx.chart.chartArea) return 'transparent'
        const gradient = ctx.chart.ctx.createLinearGradient(0, ctx.chart.chartArea.top, 0, ctx.chart.chartArea.bottom)
        gradient.addColorStop(0, `${colors.value.accentPrimary}40`)
        gradient.addColorStop(1, `${colors.value.accentPrimary}05`)
        return gradient
      },
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointRadius: 4,
      pointBackgroundColor: (ctx) => {
        const val = ctx.dataset.data[ctx.dataIndex]
        return val >= 0 ? colors.value.accentPrimary : '#EF4444'
      },
    }]
  }
})

const monthlyChartOptions = computed(() => ({
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
        label: (ctx) => `${ctx.raw >= 0 ? '+' : ''}${ctx.raw}%`
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
        callback: (v) => `${v >= 0 ? '+' : ''}${v}%`
      }
    }
  }
}))

onMounted(async () => {
  document.addEventListener('keydown', handleKeydown)
  isLoading.value = true
  try {
    report.value = await annualReportService.getAnnualReport(props.year)
  } catch (e) {
    console.error('加载年度报告失败:', e)
  } finally {
    isLoading.value = false
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="report-overlay" @click.self="emit('close')">
      <div
        class="report-container"
        @touchstart="handleTouchStart"
        @touchend="handleTouchEnd"
      >
        <!-- 头部 -->
        <div class="report-header">
          <h2 class="report-title">{{ year }} {{ t('annualReport.title') }}</h2>
          <div class="header-actions">
            <button class="action-btn" @click="emit('share')">
              <PhShareNetwork :size="16" />
            </button>
            <button class="action-btn" @click="emit('close')">
              <PhX :size="16" />
            </button>
          </div>
        </div>

        <!-- 加载中 -->
        <div v-if="isLoading" class="loading-state">
          {{ t('common.loading') }}
        </div>

        <!-- 报告内容 -->
        <div v-else-if="report" class="report-body">
          <!-- 第 1 页：年度总览 -->
          <div v-show="currentPage === 0" class="page page-summary">
            <div class="page-label">{{ t('annualReport.pageSummary') }}</div>
            <div class="big-return font-mono" :class="report.summary.totalReturn >= 0 ? 'positive' : 'negative'">
              {{ report.summary.totalReturn >= 0 ? '+' : '' }}{{ report.summary.totalReturn }}%
            </div>
            <div class="summary-subtitle">{{ t('annualReport.yearlyReturn') }}</div>

            <div class="summary-grid">
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.startValue') }}</span>
                <span class="summary-value font-mono">${{ formatNumber(report.summary.startValue, 0) }}</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.endValue') }}</span>
                <span class="summary-value font-mono">${{ formatNumber(report.summary.endValue, 0) }}</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.bestMonth') }}</span>
                <span class="summary-value font-mono positive">
                  {{ monthNames[report.summary.bestMonth.month - 1] }} +{{ report.summary.bestMonth.return }}%
                </span>
              </div>
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.worstMonth') }}</span>
                <span class="summary-value font-mono negative">
                  {{ monthNames[report.summary.worstMonth.month - 1] }} {{ report.summary.worstMonth.return }}%
                </span>
              </div>
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.totalTx') }}</span>
                <span class="summary-value font-mono">{{ report.summary.totalTransactions }}</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">{{ t('annualReport.totalFees') }}</span>
                <span class="summary-value font-mono">${{ formatNumber(report.summary.totalFeesPaid, 2) }}</span>
              </div>
            </div>

            <!-- 基准对比 -->
            <div class="benchmark-bar">
              <span class="benchmark-label">{{ t('annualReport.vsBenchmarks') }}</span>
              <span class="benchmark-item">BTC <span class="font-mono">+{{ report.benchmarks.btc }}%</span></span>
              <span class="benchmark-item">ETH <span class="font-mono">+{{ report.benchmarks.eth }}%</span></span>
              <span class="benchmark-item">S&amp;P 500 <span class="font-mono">+{{ report.benchmarks.sp500 }}%</span></span>
            </div>
          </div>

          <!-- 第 2 页：月度收益曲线 -->
          <div v-show="currentPage === 1" class="page page-monthly">
            <div class="page-label">{{ t('annualReport.pageMonthly') }}</div>
            <div class="chart-wrapper">
              <Line :data="monthlyChartData" :options="monthlyChartOptions" />
            </div>
            <div class="monthly-list">
              <div
                v-for="m in report.monthlyReturns"
                :key="m.month"
                class="monthly-item"
              >
                <span class="ml-month">{{ monthNames[m.month - 1] }}</span>
                <span class="ml-return font-mono" :class="m.return >= 0 ? 'positive' : 'negative'">
                  {{ m.return >= 0 ? '+' : '' }}{{ m.return }}%
                </span>
                <span class="ml-value font-mono">${{ formatNumber(m.value, 0) }}</span>
              </div>
            </div>
          </div>

          <!-- 第 3 页：资产配置变化 -->
          <div v-show="currentPage === 2" class="page page-allocation">
            <div class="page-label">{{ t('annualReport.pageAllocation') }}</div>
            <div class="alloc-compare">
              <div class="alloc-col">
                <h4 class="alloc-title">{{ t('annualReport.yearStart') }}</h4>
                <div
                  v-for="item in report.allocationChanges.start"
                  :key="'s-' + item.symbol"
                  class="alloc-row"
                >
                  <span class="alloc-symbol">{{ item.symbol }}</span>
                  <div class="alloc-bar-track">
                    <div class="alloc-bar" :style="{ width: item.pct + '%' }" />
                  </div>
                  <span class="alloc-pct font-mono">{{ item.pct }}%</span>
                </div>
              </div>
              <div class="alloc-arrow">&rarr;</div>
              <div class="alloc-col">
                <h4 class="alloc-title">{{ t('annualReport.yearEnd') }}</h4>
                <div
                  v-for="item in report.allocationChanges.end"
                  :key="'e-' + item.symbol"
                  class="alloc-row"
                >
                  <span class="alloc-symbol">{{ item.symbol }}</span>
                  <div class="alloc-bar-track">
                    <div class="alloc-bar accent" :style="{ width: item.pct + '%' }" />
                  </div>
                  <span class="alloc-pct font-mono">{{ item.pct }}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 第 4 页：关键里程碑 -->
          <div v-show="currentPage === 3" class="page page-milestones">
            <div class="page-label">{{ t('annualReport.pageMilestones') }}</div>
            <div class="milestone-timeline">
              <div
                v-for="(ms, idx) in report.milestones"
                :key="idx"
                class="milestone-item"
              >
                <div class="ms-icon">
                  <component :is="milestoneIcons[ms.icon] || PhStar" :size="16" weight="bold" />
                </div>
                <div class="ms-line" v-if="idx < report.milestones.length - 1" />
                <div class="ms-content">
                  <div class="ms-date">{{ ms.date }}</div>
                  <div class="ms-title">{{ ms.title }}</div>
                  <div class="ms-desc">{{ ms.description }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 第 5 页：投资风格 + 总结 -->
          <div v-show="currentPage === 4" class="page page-style">
            <div class="page-label">{{ t('annualReport.pageStyle') }}</div>
            <div class="style-tag">{{ t(styleLabels[report.styleTag]) }}</div>
            <div class="style-scores">
              <div v-for="(score, key) in report.styleScore" :key="key" class="score-item">
                <span class="score-label">{{ t(`annualReport.score_${key}`) }}</span>
                <div class="score-bar-track">
                  <div class="score-bar" :style="{ width: score + '%' }" />
                </div>
                <span class="score-value font-mono">{{ score }}</span>
              </div>
            </div>

            <!-- 年度之星 & 年度遗憾 -->
            <div class="star-regret">
              <div class="sr-card star">
                <div class="sr-header">
                  <PhTrophy :size="16" />
                  <span>{{ t('annualReport.annualStar') }}</span>
                </div>
                <div class="sr-symbol">{{ report.annualStar.symbol }}</div>
                <div class="sr-return font-mono positive">+{{ report.annualStar.return }}%</div>
                <div class="sr-reason">{{ report.annualStar.reason }}</div>
              </div>
              <div class="sr-card regret">
                <div class="sr-header">
                  <PhFlag :size="16" />
                  <span>{{ t('annualReport.annualRegret') }}</span>
                </div>
                <div class="sr-symbol">{{ report.annualRegret.symbol }}</div>
                <div class="sr-return font-mono negative">{{ report.annualRegret.return }}%</div>
                <div class="sr-reason">{{ report.annualRegret.reason }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 底部导航 -->
        <div class="report-nav">
          <button class="nav-btn" :disabled="currentPage === 0" @click="goPrev">
            <PhCaretLeft :size="16" />
          </button>
          <div class="page-dots">
            <span
              v-for="i in totalPages"
              :key="i"
              class="dot"
              :class="{ active: currentPage === i - 1 }"
              @click="currentPage = i - 1"
            />
          </div>
          <button class="nav-btn" :disabled="currentPage === totalPages - 1" @click="goNext">
            <PhCaretRight :size="16" />
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.report-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.report-container {
  width: 100%;
  max-width: 640px;
  max-height: 90vh;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-lg) var(--gap-xl);
  border-bottom: 1px solid var(--color-border);
}

.report-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.header-actions {
  display: flex;
  gap: var(--gap-xs);
}

.action-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--gap-xs);
  border-radius: var(--radius-xs);
  transition: color var(--transition-fast);
}

.action-btn:hover {
  color: var(--color-text-primary);
}

.loading-state {
  padding: var(--gap-2xl);
  text-align: center;
  color: var(--color-text-muted);
}

/* 报告主体 */
.report-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--gap-xl);
}

.page {
  min-height: 360px;
}

.page-label {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-accent-primary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: var(--gap-lg);
}

/* 第 1 页：总览 */
.big-return {
  font-size: 2.5rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: var(--gap-xs);
}

.big-return.positive { color: var(--color-success); }
.big-return.negative { color: var(--color-error); }

.summary-subtitle {
  text-align: center;
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-xl);
}

.summary-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-md);
  margin-bottom: var(--gap-lg);
}

.summary-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.summary-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.summary-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.summary-value.positive { color: var(--color-success); }
.summary-value.negative { color: var(--color-error); }

.benchmark-bar {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.benchmark-label {
  font-weight: 600;
  color: var(--color-text-muted);
  margin-right: auto;
}

.benchmark-item .font-mono {
  color: var(--color-text-primary);
}

/* 第 2 页：月度曲线 */
.chart-wrapper {
  height: 200px;
  margin-bottom: var(--gap-lg);
}

.monthly-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-xs);
}

.monthly-item {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-xs) var(--gap-sm);
  font-size: 0.6875rem;
}

.ml-month {
  color: var(--color-text-muted);
  min-width: 28px;
}

.ml-return {
  font-weight: 600;
  min-width: 42px;
  text-align: right;
}

.ml-return.positive { color: var(--color-success); }
.ml-return.negative { color: var(--color-error); }

.ml-value {
  color: var(--color-text-secondary);
  margin-left: auto;
}

/* 第 3 页：配置变化 */
.alloc-compare {
  display: flex;
  gap: var(--gap-md);
  align-items: flex-start;
}

.alloc-col {
  flex: 1;
}

.alloc-arrow {
  font-size: 1.25rem;
  color: var(--color-text-muted);
  align-self: center;
  margin-top: var(--gap-xl);
}

.alloc-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-md);
}

.alloc-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-sm);
}

.alloc-symbol {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
  min-width: 36px;
}

.alloc-bar-track {
  flex: 1;
  height: 8px;
  background: var(--color-bg-tertiary);
  border-radius: 4px;
  overflow: hidden;
}

.alloc-bar {
  height: 100%;
  background: var(--color-text-muted);
  border-radius: 4px;
  transition: width 0.6s ease;
}

.alloc-bar.accent {
  background: var(--color-accent-primary);
}

.alloc-pct {
  font-size: 0.6875rem;
  color: var(--color-text-secondary);
  min-width: 30px;
  text-align: right;
}

/* 第 4 页：里程碑 */
.milestone-timeline {
  display: flex;
  flex-direction: column;
}

.milestone-item {
  display: flex;
  gap: var(--gap-md);
  position: relative;
}

.ms-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-accent-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  z-index: 1;
}

.ms-line {
  position: absolute;
  left: 15px;
  top: 32px;
  bottom: -8px;
  width: 2px;
  background: var(--color-border);
}

.ms-content {
  flex: 1;
  padding-bottom: var(--gap-lg);
}

.ms-date {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-bottom: 2px;
}

.ms-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.ms-desc {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

/* 第 5 页：投资风格 */
.style-tag {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-accent-primary);
  text-align: center;
  margin-bottom: var(--gap-xl);
}

.style-scores {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  margin-bottom: var(--gap-xl);
}

.score-item {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
}

.score-label {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  min-width: 80px;
}

.score-bar-track {
  flex: 1;
  height: 6px;
  background: var(--color-bg-tertiary);
  border-radius: 3px;
  overflow: hidden;
}

.score-bar {
  height: 100%;
  background: var(--color-accent-primary);
  border-radius: 3px;
  transition: width 0.8s ease;
}

.score-value {
  font-size: 0.75rem;
  color: var(--color-text-primary);
  font-weight: 600;
  min-width: 24px;
  text-align: right;
}

.star-regret {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-md);
}

.sr-card {
  padding: var(--gap-md);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
}

.sr-card.star {
  background: rgba(16, 185, 129, 0.04);
}

.sr-card.regret {
  background: rgba(239, 68, 68, 0.04);
}

.sr-header {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-sm);
}

.sr-symbol {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.sr-return {
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: var(--gap-xs);
}

.sr-return.positive { color: var(--color-success); }
.sr-return.negative { color: var(--color-error); }

.sr-reason {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  line-height: 1.4;
}

/* 底部导航 */
.report-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-lg);
  padding: var(--gap-md) var(--gap-xl);
  border-top: 1px solid var(--color-border);
}

.nav-btn {
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: var(--gap-xs);
  border-radius: var(--radius-xs);
  transition: color var(--transition-fast);
}

.nav-btn:hover:not(:disabled) {
  color: var(--color-text-primary);
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-dots {
  display: flex;
  gap: var(--gap-sm);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-bg-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.dot.active {
  background: var(--color-accent-primary);
}

/* 过渡 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

/* 响应式 */
@media (max-width: 640px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }

  .alloc-compare {
    flex-direction: column;
  }

  .alloc-arrow {
    transform: rotate(90deg);
    align-self: center;
    margin-top: 0;
  }

  .star-regret {
    grid-template-columns: 1fr;
  }

  .monthly-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
