<script setup>
/**
 * 回撤曲线图
 * 展示资产组合的历史回撤情况
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
  labels: props.data.labels,
  datasets: [
    {
      label: t('risk.drawdown'),
      data: props.data.drawdown,
      borderColor: colors.value.error,
      backgroundColor: `${colors.value.error}15`,
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3
    }
  ]
}))

// 图表配置
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
      callbacks: {
        label: (ctx) => `${t('risk.drawdown')}: ${ctx.raw.toFixed(2)}%`
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
      display: true,
      grid: {
        color: colors.value.border,
        lineWidth: 0.5
      },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => `${v.toFixed(0)}%`
      }
    }
  }
}))
</script>

<template>
  <div class="drawdown-chart">
    <h3 class="chart-title">{{ t('risk.drawdownCurve') }}</h3>
    <div class="chart-container">
      <Line v-if="hasEnoughData" :data="chartData" :options="chartOptions" />
      <div v-else class="empty-chart">
        {{ t('risk.insufficientData') }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.drawdown-chart {
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
  height: 240px;
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
