<script setup>
/**
 * 资产详情抽屉组件
 * 展示单个资产的详细信息和历史走势
 */
import { ref, computed, watch, onMounted } from 'vue'
import { 
  PhX, 
  PhCaretUp, 
  PhCaretDown,
  PhArrowsClockwise,
  PhCopy,
  PhCheck,
  PhArrowUpRight,
  PhArrowDownLeft
} from '@phosphor-icons/vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip
} from 'chart.js'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'
import { useFormatters } from '../composables/useFormatters'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip
)

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  asset: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close'])

const themeStore = useThemeStore()
const { t } = useI18n()
const { formatNumber, currencySymbol } = useFormatters()

// 时间范围
const selectedRange = ref('30D')
const timeRanges = ['24H', '7D', '30D', '90D', '1Y']

// 生成模拟走势数据
const generateChartData = (range) => {
  const points = {
    '24H': 24,
    '7D': 7,
    '30D': 30,
    '90D': 90,
    '1Y': 12
  }[range]
  
  const labels = []
  const data = []
  const baseValue = props.asset?.value || 1000
  
  for (let i = 0; i < points; i++) {
    if (range === '24H') {
      labels.push(`${i}:00`)
    } else if (range === '1Y') {
      const date = new Date()
      date.setMonth(date.getMonth() - (points - 1 - i))
      labels.push(date.toLocaleDateString('zh-CN', { month: 'short' }))
    } else {
      const date = new Date()
      date.setDate(date.getDate() - (points - 1 - i))
      labels.push(date.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' }))
    }
    
    // 随机波动
    const variance = baseValue * 0.15
    data.push(baseValue + (Math.random() - 0.5) * variance)
  }
  
  return { labels, data }
}

// 图表数据
const chartData = computed(() => {
  const { labels, data } = generateChartData(selectedRange.value)
  const gradient = themeStore.currentTheme.colors.accentPrimary
  
  return {
    labels,
    datasets: [{
      data,
      borderColor: gradient,
      backgroundColor: `${gradient}20`,
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 4,
      pointHoverBackgroundColor: gradient,
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
    }]
  }
})

// 图表选项
const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index'
  },
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: themeStore.currentTheme.colors.bgElevated,
      titleColor: themeStore.currentTheme.colors.textPrimary,
      bodyColor: themeStore.currentTheme.colors.textSecondary,
      borderColor: themeStore.currentTheme.colors.border,
      borderWidth: 1,
      padding: 12,
      displayColors: false,
      callbacks: {
        label: (context) => `${currencySymbol.value}${formatNumber(context.raw)}`
      }
    }
  },
  scales: {
    x: {
      display: true,
      grid: { display: false },
      ticks: {
        color: themeStore.currentTheme.colors.textMuted,
        maxTicksLimit: 6
      }
    },
    y: {
      display: true,
      position: 'right',
      grid: {
        color: `${themeStore.currentTheme.colors.border}50`
      },
      ticks: {
        color: themeStore.currentTheme.colors.textMuted,
        callback: (value) => `${currencySymbol.value}${formatNumber(value, 0)}`
      }
    }
  }
}))

// 模拟交易历史
const transactions = computed(() => {
  if (!props.asset) return []
  return [
    { type: 'in', amount: 0.5, value: 15000, time: '2024-01-10 14:32', from: 'Binance' },
    { type: 'out', amount: 0.1, value: 3200, time: '2024-01-08 09:15', to: 'DeFi Wallet' },
    { type: 'in', amount: 0.3, value: 9500, time: '2024-01-05 18:45', from: 'OKX' },
    { type: 'in', amount: 0.2, value: 6800, time: '2024-01-02 11:20', from: 'Coinbase' }
  ]
})

// 关闭抽屉
const close = () => {
  emit('close')
}
</script>

<template>
  <Transition name="drawer">
    <div v-if="visible" class="drawer-overlay" @click.self="close">
      <div class="drawer-panel">
        <!-- 头部 -->
        <div class="drawer-header">
          <div class="asset-info" v-if="asset">
            <div class="asset-icon">
              <img 
                v-if="asset.icon"
                :src="asset.icon" 
                :alt="asset.symbol"
              />
              <span v-else class="asset-symbol">{{ asset.symbol?.slice(0, 2) }}</span>
            </div>
            <div class="asset-meta">
              <h2 class="asset-name">{{ asset.name }}</h2>
              <span class="asset-symbol-badge">{{ asset.symbol }}</span>
            </div>
          </div>
          <button class="close-btn" @click="close">
            <PhX :size="24" />
          </button>
        </div>
        
        <!-- 价值概览 -->
        <div class="value-section" v-if="asset">
          <div class="current-value">
            <span class="value-label">{{ t('dashboard.holdings') }}</span>
            <span class="value-amount font-mono">
              {{ currencySymbol }}{{ formatNumber(asset.value) }}
            </span>
          </div>
          <div class="value-stats">
            <div class="stat-item">
              <span class="stat-label">{{ t('dashboard.balance') }}</span>
              <span class="stat-value font-mono">{{ asset.balance }} {{ asset.symbol }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">{{ t('dashboard.price') }}</span>
              <span class="stat-value font-mono">{{ currencySymbol }}{{ formatNumber(asset.price) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">{{ t('dashboard.change') }}</span>
              <span 
                class="stat-value font-mono"
                :class="asset.change24h >= 0 ? 'positive' : 'negative'"
              >
                <PhCaretUp v-if="asset.change24h >= 0" :size="14" weight="bold" />
                <PhCaretDown v-else :size="14" weight="bold" />
                {{ Math.abs(asset.change24h) }}%
              </span>
            </div>
          </div>
        </div>
        
        <!-- 走势图 -->
        <div class="chart-section">
          <div class="chart-header">
            <h3>{{ t('dashboard.assetTrend') }}</h3>
            <div class="time-range-selector">
              <button
                v-for="range in timeRanges"
                :key="range"
                class="range-btn"
                :class="{ 'active': selectedRange === range }"
                @click="selectedRange = range"
              >
                {{ range }}
              </button>
            </div>
          </div>
          <div class="chart-container">
            <Line :data="chartData" :options="chartOptions" />
          </div>
        </div>
        
        <!-- 交易历史 -->
        <div class="transactions-section">
          <h3>{{ t('dashboard.recentTransactions') || '近期交易' }}</h3>
          <div class="transactions-list">
            <div 
              v-for="(tx, index) in transactions" 
              :key="index"
              class="transaction-item"
            >
              <div class="tx-icon" :class="tx.type">
                <PhArrowDownLeft v-if="tx.type === 'in'" :size="16" />
                <PhArrowUpRight v-else :size="16" />
              </div>
              <div class="tx-info">
                <span class="tx-type">
                  {{ tx.type === 'in' ? (tx.from ? `从 ${tx.from} 转入` : '转入') : (tx.to ? `转出至 ${tx.to}` : '转出') }}
                </span>
                <span class="tx-time">{{ tx.time }}</span>
              </div>
              <div class="tx-amount" :class="tx.type">
                <span class="amount font-mono">
                  {{ tx.type === 'in' ? '+' : '-' }}{{ tx.amount }} {{ asset?.symbol }}
                </span>
                <span class="value font-mono">≈ {{ currencySymbol }}{{ formatNumber(tx.value) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* 抽屉动画 */
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity var(--transition-base);
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .drawer-panel,
.drawer-leave-to .drawer-panel {
  transform: translateX(100%);
}

/* 遮罩层 */
.drawer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
  justify-content: flex-end;
}

/* 抽屉面板 */
.drawer-panel {
  width: 100%;
  max-width: 440px;
  height: 100%;
  overflow-y: auto;
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* 头部 */
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.asset-info {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.asset-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--color-bg-elevated);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.asset-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.asset-symbol {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text-primary);
}

.asset-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.asset-name {
  font-family: var(--font-heading);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.asset-symbol-badge {
  font-size: 11px;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.close-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.close-btn:hover {
  color: var(--color-text-primary);
}

/* 价值概览 */
.value-section {
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-md);
  padding: var(--gap-md);
}

.current-value {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: var(--gap-md);
  padding-bottom: var(--gap-md);
  border-bottom: 1px solid var(--color-border);
}

.value-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.value-amount {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.value-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-sm);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 11px;
  color: var(--color-text-muted);
}

.stat-value {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.stat-value.positive {
  color: var(--color-success);
}

.stat-value.negative {
  color: var(--color-error);
}

/* 图表区域 */
.chart-section h3,
.transactions-section h3 {
  font-family: var(--font-heading);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-sm);
}

.chart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-sm);
}

.time-range-selector {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 3px;
  border-radius: var(--radius-sm);
}

.range-btn {
  padding: 4px 10px;
  font-size: 11px;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.range-btn:hover {
  color: var(--color-text-primary);
}

.range-btn.active {
  background: var(--color-accent-primary);
  color: white;
}

.chart-container {
  height: 180px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-md);
  padding: var(--gap-sm);
}

/* 交易历史 */
.transactions-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.transaction-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.tx-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.tx-icon.in {
  background: color-mix(in srgb, var(--color-success) 15%, transparent);
  color: var(--color-success);
}

.tx-icon.out {
  background: color-mix(in srgb, var(--color-error) 15%, transparent);
  color: var(--color-error);
}

.tx-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.tx-type {
  font-size: 13px;
  color: var(--color-text-primary);
}

.tx-time {
  font-size: 11px;
  color: var(--color-text-muted);
}

.tx-amount {
  text-align: right;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.tx-amount .amount {
  font-size: 13px;
  font-weight: 500;
}

.tx-amount.in .amount {
  color: var(--color-success);
}

.tx-amount.out .amount {
  color: var(--color-error);
}

.tx-amount .value {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 响应式 */
@media (max-width: 520px) {
  .drawer-panel {
    max-width: 100%;
  }

  .value-stats {
    grid-template-columns: 1fr;
  }
}
</style>
