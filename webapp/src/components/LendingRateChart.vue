<script setup>
/**
 * 借贷利率历史图表组件
 * 展示存款和借款利率的历史趋势
 */
import { ref, computed, watch } from 'vue'
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler } from 'chart.js'
import { useThemeStore } from '../stores/themeStore'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const props = defineProps({
  history: {
    type: Array,
    required: true
  },
  protocol: {
    type: String,
    required: true
  },
  token: {
    type: String,
    required: true
  }
})

const themeStore = useThemeStore()

// 图表数据
const chartData = computed(() => {
  const labels = props.history.map(h => h.date)
  const supplyData = props.history.map(h => h.supply_apy)
  const borrowStableData = props.history.map(h => h.borrow_apy_stable)
  const borrowVariableData = props.history.map(h => h.borrow_apy_variable)

  return {
    labels,
    datasets: [
      {
        label: '存款 APY',
        data: supplyData,
        borderColor: '#10B981',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 4,
        tension: 0.3,
        fill: true
      },
      {
        label: '借款 APY (稳定)',
        data: borrowStableData,
        borderColor: '#F59E0B',
        backgroundColor: 'rgba(245, 158, 11, 0.1)',
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 4,
        tension: 0.3,
        fill: true
      },
      {
        label: '借款 APY (浮动)',
        data: borrowVariableData,
        borderColor: '#EF4444',
        backgroundColor: 'rgba(239, 68, 68, 0.1)',
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 4,
        tension: 0.3,
        fill: true
      }
    ]
  }
})

// 图表配置
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      labels: {
        color: themeStore.isDark ? '#9CA3AF' : '#6B7280',
        font: { size: 11, family: 'IBM Plex Sans' },
        padding: 12,
        usePointStyle: true,
        pointStyle: 'circle'
      }
    },
    tooltip: {
      backgroundColor: themeStore.isDark ? '#1F2937' : '#FFFFFF',
      titleColor: themeStore.isDark ? '#F9FAFB' : '#111827',
      bodyColor: themeStore.isDark ? '#D1D5DB' : '#4B5563',
      borderColor: themeStore.isDark ? '#374151' : '#E5E7EB',
      borderWidth: 1,
      padding: 12,
      displayColors: true,
      callbacks: {
        label: (context) => {
          return `${context.dataset.label}: ${context.parsed.y.toFixed(2)}%`
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        display: false
      },
      ticks: {
        color: themeStore.isDark ? '#6B7280' : '#9CA3AF',
        font: { size: 10, family: 'IBM Plex Mono' },
        maxRotation: 0
      }
    },
    y: {
      beginAtZero: true,
      grid: {
        color: themeStore.isDark ? '#374151' : '#F3F4F6',
        drawBorder: false
      },
      ticks: {
        color: themeStore.isDark ? '#6B7280' : '#9CA3AF',
        font: { size: 10, family: 'IBM Plex Mono' },
        callback: (value) => `${value}%`
      }
    }
  }
}))
</script>

<template>
  <div class="lending-rate-chart">
    <div class="chart-header">
      <div class="chart-title">
        <span class="protocol-name">{{ protocol }}</span>
        <span class="token-name">{{ token }}</span>
      </div>
      <div class="chart-period">{{ history.length }} 天历史</div>
    </div>
    <div class="chart-container">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<style scoped>
.lending-rate-chart {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-title {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.protocol-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.token-name {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
}

.chart-period {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.chart-container {
  height: 240px;
}
</style>
