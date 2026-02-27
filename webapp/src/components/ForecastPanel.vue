<script setup>
/**
 * 趋势预测面板
 * 展示线性回归预测曲线（中值 + 置信区间阴影）
 */
import { ref, computed, onMounted, watch } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { PhTrendUp, PhTrendDown, PhEquals, PhChartLineUp } from '@phosphor-icons/vue'
import { analyticsService } from '../api/index.js'
import { useFormatters } from '../composables/useFormatters'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

ChartJS.register(
  CategoryScale, LinearScale, PointElement, LineElement,
  Title, Tooltip, Legend, Filler
)

const { formatNumber, currencySymbol } = useFormatters()
const themeStore = useThemeStore()
const { t } = useI18n()

// 预测天数
const selectedDays = ref(14)
const daysOptions = [7, 14, 30]

// 数据
const data = ref(null)
const isLoading = ref(false)
const hasError = ref(false)

const loadData = async () => {
  isLoading.value = true
  hasError.value = false
  try {
    data.value = await analyticsService.getForecast(selectedDays.value)
  } catch (err) {
    console.error('加载预测数据失败:', err)
    hasError.value = true
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)
watch(selectedDays, loadData)

const colors = computed(() => themeStore.currentTheme.colors)

// 趋势相关
const trend = computed(() => data.value?.trend || 'flat')
const confidence = computed(() => data.value?.confidence || 0)
const slope = computed(() => data.value?.slope || 0)
const forecastPoints = computed(() => data.value?.forecast_points || [])

const trendIcon = computed(() => {
  const map = { up: PhTrendUp, down: PhTrendDown, flat: PhEquals }
  return map[trend.value] || PhEquals
})

const trendColor = computed(() => {
  const map = { up: 'var(--color-success)', down: 'var(--color-error)', flat: 'var(--color-text-muted)' }
  return map[trend.value] || 'var(--color-text-muted)'
})

const trendLabel = computed(() => {
  const map = { up: 'trendUp', down: 'trendDown', flat: 'trendFlat' }
  return t(`analytics.forecast.${map[trend.value] || 'trendFlat'}`)
})

// Chart.js 数据
const chartData = computed(() => {
  if (!forecastPoints.value.length) return { labels: [], datasets: [] }

  const labels = forecastPoints.value.map(p => {
    const d = new Date(p.date)
    return `${d.getMonth() + 1}/${d.getDate()}`
  })

  return {
    labels,
    datasets: [
      // 上界（用于填充区间阴影）
      {
        label: t('analytics.forecast.upperBound'),
        data: forecastPoints.value.map(p => p.upper),
        borderColor: 'transparent',
        backgroundColor: 'transparent',
        borderWidth: 0,
        pointRadius: 0,
        fill: false,
      },
      // 下界（填充到上界构成置信区间）
      {
        label: t('analytics.forecast.lowerBound'),
        data: forecastPoints.value.map(p => p.lower),
        borderColor: 'transparent',
        backgroundColor: `${colors.value.accentPrimary}15`,
        borderWidth: 0,
        pointRadius: 0,
        fill: '-1', // 填充到上一个 dataset（upper）
      },
      // 中值预测线
      {
        label: t('analytics.forecast.predictedValue'),
        data: forecastPoints.value.map(p => p.value),
        borderColor: colors.value.accentPrimary,
        backgroundColor: 'transparent',
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 4,
        fill: false,
        tension: 0.3,
      },
    ]
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      filter: (item) => item.datasetIndex === 2, // 仅显示中值
      callbacks: {
        label: (ctx) => `${t('analytics.forecast.predictedValue')}: ${currencySymbol.value}${formatNumber(ctx.raw, 0)}`
      }
    }
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: colors.value.textMuted, maxTicksLimit: 7, font: { size: 10 } }
    },
    y: {
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => `${currencySymbol.value}${formatNumber(v, 0)}`
      }
    }
  }
}))
</script>

<template>
  <section class="forecast-panel">
    <!-- 标题栏 -->
    <div class="panel-header">
      <h3 class="panel-title">
        <PhChartLineUp :size="16" weight="bold" />
        {{ t('analytics.forecast.title') }}
      </h3>
      <div class="range-selector">
        <button
          v-for="d in daysOptions"
          :key="d"
          class="range-btn"
          :class="{ active: selectedDays === d }"
          @click="selectedDays = d"
        >
          {{ d }}{{ t('settings.day') }}
        </button>
      </div>
    </div>

    <!-- 加载中 -->
    <div v-if="isLoading" class="loading-state">
      {{ t('common.loading') }}
    </div>

    <!-- 错误 -->
    <div v-else-if="hasError" class="empty-state">
      <p>{{ t('analytics.forecast.error') }}</p>
    </div>

    <!-- 无数据 -->
    <div v-else-if="!data || forecastPoints.length === 0" class="empty-state">
      <p>{{ t('analytics.forecast.noData') }}</p>
    </div>

    <!-- 主体 -->
    <template v-else>
      <!-- 指标卡片行 -->
      <div class="metrics-row">
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.forecast.trendUp').split('').shift() && '趋势' || 'Trend' }}</span>
          <div class="metric-trend" :style="{ color: trendColor }">
            <component :is="trendIcon" :size="20" weight="bold" />
            <span class="metric-value font-mono">{{ trendLabel }}</span>
          </div>
        </div>
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.forecast.confidence') }}</span>
          <span class="metric-value font-mono">{{ (confidence * 100).toFixed(0) }}%</span>
          <div class="confidence-bar">
            <div class="confidence-fill" :style="{ width: `${confidence * 100}%` }" />
          </div>
        </div>
        <div class="metric-card">
          <span class="metric-label">{{ t('analytics.forecast.dailyGrowth') }}</span>
          <span class="metric-value font-mono" :class="slope >= 0 ? 'positive' : 'negative'">
            {{ slope >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(Math.abs(slope), 0) }}/{{ t('settings.day') }}
          </span>
        </div>
      </div>

      <!-- 图表 -->
      <div class="chart-container">
        <Line :data="chartData" :options="chartOptions" :key="selectedDays" />
      </div>

      <!-- 底部提示 -->
      <div class="forecast-hint">
        <PhEquals :size="12" />
        <span>{{ t('analytics.forecast.hint') }}</span>
      </div>
    </template>
  </section>
</template>

<style scoped>
.forecast-panel {
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
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
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

/* 指标卡片 */
.metrics-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-md);
}

.metric-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.metric-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.metric-value {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.metric-trend {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

/* 置信度条 */
.confidence-bar {
  height: 4px;
  background: var(--color-bg-primary);
  border-radius: 2px;
  margin-top: 4px;
}

.confidence-fill {
  height: 100%;
  background: var(--color-accent-primary);
  border-radius: 2px;
  transition: width 0.4s ease;
}

/* 图表 */
.chart-container {
  height: 240px;
}

/* 状态 */
.positive {
  color: var(--color-success);
}

.negative {
  color: var(--color-error);
}

/* 底部提示 */
.forecast-hint {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

/* 空/加载 */
.loading-state,
.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .metrics-row {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
