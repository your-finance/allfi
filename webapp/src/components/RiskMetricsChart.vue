<script setup>
/**
 * 风险指标历史趋势图
 * 展示波动率和夏普比率的历史变化
 */
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
})

const themeStore = useThemeStore()
const { t } = useI18n()
const colors = computed(() => themeStore.currentTheme.colors)

// 检查是否有足够的数据点（至少 2 个）
const hasEnoughData = computed(() => {
  const labels = props.data?.labels
  return Array.isArray(labels) && labels.length >= 2
})

// 图表数据
const chartData = computed(() => ({
  labels: props.data.labels || [],
  datasets: [
    {
      label: t('risk.volatility'),
      data: props.data.volatility || [],
      borderColor: colors.value.accentPrimary,
      backgroundColor: `${colors.value.accentPrimary}15`,
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3,
      yAxisID: 'y'
    },
    {
      label: t('risk.sharpeRatio'),
      data: props.data.sharpe || [],
      borderColor: colors.value.accentSecondary,
      backgroundColor: 'transparent',
      borderWidth: 2,
      borderDash: [6, 3],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3,
      yAxisID: 'y1'
    }
  ]
}))

// 图表配置
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      align: 'end',
      labels: {
        color: colors.value.textSecondary,
        usePointStyle: true,
        pointStyle: 'line',
        padding: 12,
        font: { size: 11 }
      }
    },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      callbacks: {
        label: (ctx) => {
          const label = ctx.dataset.label
          const value = ctx.raw.toFixed(2)
          return ctx.datasetIndex === 0 ? `${label}: ${value}%` : `${label}: ${value}`
        }
      }
    }
  },
  scales: {
    x: {
      display: true,
      grid: { display: false },
      ticks: {
        color: colors.value.textMuted,
        maxTicksLimit: 8,
        font: { size: 10 }
      }
    },
    y: {
      type: 'linear',
      display: true,
      position: 'left',
      grid: {
        color: colors.value.border,
        lineWidth: 0.5
      },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => `${v.toFixed(0)}%`
      },
      title: {
        display: true,
        text: t('risk.volatility'),
        color: colors.value.textMuted,
        font: { size: 10 }
      }
    },
    y1: {
      type: 'linear',
      display: true,
      position: 'right',
      grid: { display: false },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => v.toFixed(1)
      },
      title: {
        display: true,
        text: t('risk.sharpeRatio'),
        color: colors.value.textMuted,
        font: { size: 10 }
      }
    }
  }
}))
</script>

<template>
  <div class="risk-metrics-chart">
    <h3 class="chart-title">{{ t('risk.metricsHistory') }}</h3>
    <div class="chart-container">
      <Line v-if="hasEnoughData" :data="chartData" :options="chartOptions" />
      <div v-else class="empty-chart">
        {{ t('risk.insufficientData') }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.risk-metrics-chart {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.chart-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.chart-container {
  height: 280px;
}

.empty-chart {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}
</style>
