<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import {
  PhChartLineUp,
  PhWallet,
  PhBank,
  PhCaretDown,
  PhCheck,
  PhCalendar,
  PhCurrencyDollar,
  PhPlus,
  PhMinus
} from '@phosphor-icons/vue'
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
import CalendarHeatmap from '../components/CalendarHeatmap.vue'
import TransactionTimeline from '../components/TransactionTimeline.vue'
import { useAssetStore } from '../stores/assetStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import { useThemeStore } from '../stores/themeStore' // 用于获取主题颜色

// 注册 Chart.js 组件
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const assetStore = useAssetStore()
const { formatCurrency, formatNumber, formatRelativeTime, pricingCurrency, currencies } = useFormatters()
const { t } = useI18n()
const themeStore = useThemeStore()

// 计价货币切换函数
const setPricingCurrency = (currency) => {
  pricingCurrency.value = currency
}

// 页面 Tab：'trend' 资产趋势 / 'transactions' 交易记录
const activeHistoryTab = ref('trend')

// 日历视图模式（用于布局切换）
const calendarViewMode = ref('month')

// 日历组件引用
const calendarRef = ref(null)
// 图表组件引用（用于点击选中后手动刷新）
const chartRef = ref(null)

// 筛选状态
const selectedTimeRange = ref('30D')
const timeRanges = ['7D', '30D', '90D', 'LAST_YEAR', 'ALL']
const showTimeRangeDropdown = ref(false)

// 标记：是否由日历翻月触发的时间范围变化（避免循环重置）
const isCalendarDriven = ref(false)

// 日历选中日期（YYYY-MM-DD），用于图表高亮
const selectedDate = ref(null)

// 图表 hover 中的日期（YYYY-MM-DD），用于日历临时高亮
const hoveredDate = ref(null)

// 日历 hover 中的日期（YYYY-MM-DD），用于图表临时高亮
const calendarHoveredDate = ref(null)

// 汇率相关计算属性
const showExchangeRate = computed(() => pricingCurrency.value !== 'USDC')
const exchangeRates = computed(() => assetStore.exchangeRates)

// 计价货币符号和精度
const pricingCurrencySymbol = computed(() => {
  const symbols = { USDC: '$', BTC: '₿', ETH: 'Ξ', CNY: '¥' }
  return symbols[pricingCurrency.value] || '$'
})

const pricingCurrencyDecimals = computed(() => {
  switch (pricingCurrency.value) {
    case 'BTC': return 5
    case 'ETH': return 4
    default: return 2
  }
})

const selectedAssetType = ref('all') // 'all', 'cex', 'blockchain', 'manual'
const assetTypes = computed(() => [
  { id: 'all', labelKey: 'history.allAssets', icon: PhChartLineUp },
  { id: 'cex', labelKey: 'accounts.cexTab', icon: PhCurrencyDollar },
  { id: 'blockchain', labelKey: 'accounts.blockchainTab', icon: PhWallet },
  { id: 'manual', labelKey: 'accounts.manualTab', icon: PhBank },
])

// 历史数据
const historyData = computed(() => assetStore.historyData)

// 图表数据
const lineChartData = computed(() => {
  if (!historyData.value || !historyData.value.labels || !historyData.value.values) {
    return { labels: [], datasets: [] }
  }

  const labels = historyData.value.labels

  // 根据当前计价货币转换数据
  const convertedValues = historyData.value.values.map(val =>
    assetStore.convertValue(val, pricingCurrency.value)
  )

  // 选中 / hover 日期对应的索引（选中优先：相同日期时 hover 让位）
  const selectedIdx = selectedDate.value ? labels.indexOf(selectedDate.value) : -1
  const rawHoverIdx = calendarHoveredDate.value ? labels.indexOf(calendarHoveredDate.value) : -1
  const hoverIdx = (rawHoverIdx >= 0 && rawHoverIdx !== selectedIdx) ? rawHoverIdx : -1
  const hasHighlight = selectedIdx >= 0 || hoverIdx >= 0

  const accentColor = themeStore.currentTheme.colors.accentPrimary

  // 按数据点生成样式数组：选中点强高亮，hover 点轻高亮，其余隐藏
  const pointRadiusArr = labels.map((_, i) =>
    i === selectedIdx ? 6 : i === hoverIdx ? 4 : 0
  )
  const pointHoverRadiusArr = labels.map((_, i) =>
    i === selectedIdx ? 8 : i === hoverIdx ? 6 : 4
  )
  const pointBgArr = labels.map((_, i) =>
    i === selectedIdx ? '#fff' : i === hoverIdx ? accentColor + '40' : accentColor
  )
  const pointBorderArr = labels.map((_, i) =>
    i === selectedIdx ? accentColor : i === hoverIdx ? accentColor : accentColor
  )
  const pointBorderWArr = labels.map((_, i) =>
    i === selectedIdx ? 3 : i === hoverIdx ? 2 : 1
  )

  // 主数据集：资产总价值
  const datasets = [{
    label: t('history.totalValue'),
    data: convertedValues,
    borderColor: accentColor,
    backgroundColor: `${accentColor}08`,
    borderWidth: 1.5,
    fill: true,
    tension: 0.4,
    pointRadius: hasHighlight ? pointRadiusArr : 0,
    pointHoverRadius: hasHighlight ? pointHoverRadiusArr : 4,
    pointBackgroundColor: hasHighlight ? pointBgArr : accentColor,
    pointBorderColor: hasHighlight ? pointBorderArr : accentColor,
    pointBorderWidth: hasHighlight ? pointBorderWArr : 1,
    yAxisID: 'y'
  }]

  // 第二条数据集：汇率变化线（仅在非 USDC 计价时显示）
  if (showExchangeRate.value) {
    const baseRate = exchangeRates.value[pricingCurrency.value] || 1
    // 生成汇率波动数据（模拟）
    const rateData = labels.map(() => baseRate * (0.92 + Math.random() * 0.16))

    datasets.push({
      label: `${pricingCurrency.value}/USDC`,
      data: rateData,
      borderColor: themeStore.currentTheme.colors.accentSecondary,
      backgroundColor: 'transparent',
      borderWidth: 1,
      borderDash: [4, 4],
      fill: false,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 3,
      yAxisID: 'y1'
    })
  }

  return { labels, datasets }
})

// 自定义 Chart.js 插件：选中/hover 点垂直参考线
const verticalLinePlugin = {
  id: 'verticalLine',
  afterDraw(chart) {
    const accentColor = themeStore.currentTheme.colors.accentPrimary
    const labels = chart.data.labels || []
    const meta = chart.getDatasetMeta(0)
    if (!meta) return
    const { ctx, chartArea } = chart

    // 辅助：绘制一条垂直虚线
    const drawLine = (dateStr, opacity, dashPattern) => {
      if (!dateStr) return
      const idx = labels.indexOf(dateStr)
      if (idx < 0 || !meta.data[idx]) return
      ctx.save()
      ctx.beginPath()
      ctx.setLineDash(dashPattern)
      ctx.lineWidth = 1
      ctx.strokeStyle = accentColor + opacity
      ctx.moveTo(meta.data[idx].x, chartArea.top)
      ctx.lineTo(meta.data[idx].x, chartArea.bottom)
      ctx.stroke()
      ctx.restore()
    }

    // 日历 hover 参考线（更淡；与选中日期相同时跳过，避免重叠）
    if (calendarHoveredDate.value !== selectedDate.value) {
      drawLine(calendarHoveredDate.value, '35', [2, 3])
    }
    // 选中日期参考线（更实，优先级最高）
    drawLine(selectedDate.value, '60', [3, 3])
  }
}

const lineChartOptions = computed(() => {
  const colors = themeStore.currentTheme.colors

  const baseOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: { intersect: false, mode: 'index' },
    onClick: (event, elements, chart) => {
      if (elements.length > 0) {
        const idx = elements[0].index
        const dateLabel = chart.data.labels[idx]
        // 点击同一日期取消选中，否则选中新日期
        const newDate = selectedDate.value === dateLabel ? null : dateLabel
        selectedDate.value = newDate
        // 导航日历到对应月份
        if (newDate && calendarRef.value) {
          calendarRef.value.navigateToDate(newDate)
        }
      }
    },
    onHover: (event, elements, chart) => {
      // 图表 hover 时更新日历临时高亮日期
      if (elements.length > 0) {
        hoveredDate.value = chart.data.labels[elements[0].index]
      } else {
        hoveredDate.value = null
      }
    },
    plugins: {
      legend: {
        display: showExchangeRate.value,
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
        callbacks: {
          label: (context) => {
            const label = context.dataset.label || ''
            const value = context.raw
            if (context.datasetIndex === 0) {
              // 资产总值
              return `${label}: ${pricingCurrencySymbol.value}${formatCurrency(value, { showSymbol: false })}`
            } else {
              // 汇率
              return `${label}: $${value?.toFixed(2)}`
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
          maxTicksLimit: 7,
          callback: function(value, index) {
            // 将 YYYY-MM-DD 转为短日期显示（M/D 或 M月）
            const label = this.getLabelForValue(value)
            if (!label || !label.includes('-')) return label
            const parts = label.split('-')
            const m = parseInt(parts[1])
            const d = parseInt(parts[2])
            // 长时间范围只显示月份
            if (selectedTimeRange.value === 'LAST_YEAR' || selectedTimeRange.value === 'ALL') {
              return d === 1 || index === 0 ? `${m}月` : ''
            }
            return `${m}/${d}`
          }
        }
      },
      y: {
        type: 'linear',
        display: true,
        position: 'left',
        grid: {
          color: colors.border,
          lineWidth: 0.5
        },
        ticks: {
          color: colors.textSecondary,
          font: { size: 11 },
          callback: (value) => `${pricingCurrencySymbol.value}${formatCurrency(value, { showSymbol: false })}`
        }
      }
    }
  }

  // 非 USDC 计价时添加右侧 Y 轴（汇率轴）
  if (showExchangeRate.value) {
    baseOptions.scales.y1 = {
      type: 'linear',
      display: true,
      position: 'right',
      grid: { drawOnChartArea: false },
      ticks: {
        color: colors.accentSecondary,
        font: { size: 10 },
        callback: (value) => `$${value?.toFixed(2)}`
      }
    }
  }

  return baseOptions
})

// 日历热力图数据：从历史快照中计算每日变化百分比
const heatmapData = computed(() => {
  if (!historyData.value || !historyData.value.labels || !historyData.value.values) {
    return {}
  }
  const { labels, values } = historyData.value
  const result = {}
  for (let i = 1; i < labels.length; i++) {
    // labels 统一为 YYYY-MM-DD 格式
    const dateStr = labels[i]
    if (values[i - 1] > 0) {
      const change = ((values[i] - values[i - 1]) / values[i - 1]) * 100
      result[dateStr] = parseFloat(change.toFixed(2))
    }
  }
  return result
})

// 历史资金变化记录（模拟数据）
// 不按事件类型分类，只记录资金快照的变化
const historyRecords = ref([
  {
    id: 1,
    source: 'CEX',
    sourceName: 'Binance',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'increase', asset: 'BTC', amount: 3.12 },
      { action: 'decrease', asset: 'USDT', amount: 32876 }
    ]
  },
  {
    id: 2,
    source: 'Blockchain',
    sourceName: 'MetaMask',
    timestamp: new Date(Date.now() - 5 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'decrease', asset: 'ETH', amount: 2.5 }
    ]
  },
  {
    id: 3,
    source: 'CEX',
    sourceName: 'OKX',
    timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'increase', asset: 'SOL', amount: 150 },
      { action: 'decrease', asset: 'USDC', amount: 3200 }
    ]
  },
  {
    id: 4,
    source: 'Manual',
    sourceName: '工资卡',
    timestamp: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'increase', asset: 'CNY', amount: 15000 }
    ]
  },
  {
    id: 5,
    source: 'CEX',
    sourceName: 'Binance',
    timestamp: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'increase', asset: 'ETH', amount: 8.5 },
      { action: 'decrease', asset: 'USDT', amount: 18500 }
    ]
  },
  {
    id: 6,
    source: 'Blockchain',
    sourceName: '冷钱包',
    timestamp: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
    changes: [
      { action: 'increase', asset: 'BTC', amount: 0.25 }
    ]
  }
])

// 过滤后的记录：选中日期时只显示当天记录
const filteredRecords = computed(() => {
  if (!selectedDate.value) return historyRecords.value
  // 将 selectedDate (YYYY-MM-DD) 与记录的 timestamp 比较日期部分
  return historyRecords.value.filter(record => {
    const recordDate = new Date(record.timestamp)
    const y = recordDate.getFullYear()
    const m = String(recordDate.getMonth() + 1).padStart(2, '0')
    const d = String(recordDate.getDate()).padStart(2, '0')
    return `${y}-${m}-${d}` === selectedDate.value
  })
})

// 选中日期的友好显示格式
const selectedDateLabel = computed(() => {
  if (!selectedDate.value) return ''
  const [y, m, d] = selectedDate.value.split('-')
  return `${parseInt(m)}月${parseInt(d)}日`
})

// === 日历 ↔ 时间范围联动 ===

// 日历翻月时：自动扩展时间范围以覆盖该月
const handleCalendarMonthChange = ({ year, month, monthOffset: offset }) => {
  if (offset === 0) return // 当前月无需变更

  // 计算该月 1 日距今的天数
  const monthStart = new Date(year, month, 1)
  const now = new Date()
  now.setHours(0, 0, 0, 0)
  const daysBack = Math.ceil((now - monthStart) / (86400000))

  // 选择能覆盖该月的最小时间范围
  let targetRange
  if (daysBack <= 7) targetRange = '7D'
  else if (daysBack <= 30) targetRange = '30D'
  else if (daysBack <= 90) targetRange = '90D'
  else if (daysBack <= 365) targetRange = 'LAST_YEAR'
  else targetRange = 'ALL'

  // 如果需要扩展范围，标记为日历触发后更新
  if (targetRange !== selectedTimeRange.value) {
    isCalendarDriven.value = true
    selectedTimeRange.value = targetRange
  }
}

// 监听时间范围变化，加载对应的历史数据
watch(selectedTimeRange, async (newRange) => {
  await assetStore.loadHistory(newRange)

  // 清除选中日期（新数据范围可能不包含之前选中的日期）
  selectedDate.value = null

  // 如果是用户通过下拉菜单切换（非日历驱动），重置日历到当前月
  if (!isCalendarDriven.value && calendarRef.value) {
    calendarRef.value.resetToCurrentMonth()
  }
  isCalendarDriven.value = false
}, { immediate: true })

// 监听货币变化，重新加载历史数据（确保图表数据正确转换）
watch(pricingCurrency, async () => {
  await assetStore.loadHistory(selectedTimeRange.value)
})

// 选中日期或日历 hover 日期变化时手动刷新图表（更新高亮点和参考线）
watch([selectedDate, calendarHoveredDate], () => {
  if (chartRef.value?.chart) {
    chartRef.value.chart.update('none') // 'none' 禁用动画，立即刷新
  }
})

// 点击外部关闭时间范围下拉菜单
const handleClickOutside = (e) => {
  if (!e.target.closest('.time-range-dropdown')) {
    showTimeRangeDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="history-page">
    <!-- Tab 切换：资产趋势 / 交易记录 -->
    <div class="history-tabs">
      <button
        class="history-tab"
        :class="{ active: activeHistoryTab === 'trend' }"
        @click="activeHistoryTab = 'trend'"
      >
        <PhChartLineUp :size="16" />
        {{ t('history.trendTab') }}
      </button>
      <button
        class="history-tab"
        :class="{ active: activeHistoryTab === 'transactions' }"
        @click="activeHistoryTab = 'transactions'"
      >
        <PhWallet :size="16" />
        {{ t('history.transactionsTab') }}
      </button>
    </div>

    <!-- 交易记录时间线 -->
    <TransactionTimeline v-if="activeHistoryTab === 'transactions'" />

    <!-- === 资产趋势（原有内容） === -->
    <template v-if="activeHistoryTab === 'trend'">
    <!-- 页面头部和筛选器 -->
    <header class="page-header">
      <div class="header-filters">
        <!-- 资产类型筛选 -->
        <div class="filter-group">
          <label class="filter-label">{{ t('history.assetType') }}:</label>
          <div class="type-filter-buttons">
            <button 
              v-for="type in assetTypes" 
              :key="type.id"
              class="filter-button"
              :class="{ 'active': selectedAssetType === type.id }"
              @click="selectedAssetType = type.id"
            >
              <component :is="type.icon" :size="18" weight="duotone" />
              <span class="type-label">{{ t(type.labelKey) }}</span>
            </button>
          </div>
        </div>
        
        <!-- 时间范围筛选 -->
        <div class="filter-group">
          <label class="filter-label">{{ t('history.timeRange') }}:</label>
          <div class="time-range-dropdown">
            <button 
              class="range-display-btn"
              @click="showTimeRangeDropdown = !showTimeRangeDropdown"
            >
              <PhCalendar :size="18" />
              <span>{{ selectedTimeRange === 'ALL' ? t('history.allTime') : (selectedTimeRange === 'LAST_YEAR' ? t('history.lastYear') : t('history.lastDays', { days: parseInt(selectedTimeRange) })) }}</span>
              <PhCaretDown :size="14" :class="{ 'arrow-rotated': showTimeRangeDropdown }" />
            </button>
            <Transition name="dropdown">
              <div v-if="showTimeRangeDropdown" class="dropdown-menu">
                <button
                  v-for="range in timeRanges"
                  :key="range"
                  class="dropdown-item"
                  :class="{ 'dropdown-item-active': selectedTimeRange === range }"
                  @click="selectedTimeRange = range; showTimeRangeDropdown = false;"
                >
                  {{ range === 'ALL' ? t('history.allTime') : (range === 'LAST_YEAR' ? t('history.lastYear') : t('history.lastDays', { days: parseInt(range) })) }}
                  <PhCheck v-if="selectedTimeRange === range" :size="16" />
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </header>

    <!-- 趋势区域：日历 + 曲线同行（年视图时堆叠） -->
    <div class="trend-layout" :class="{ 'trend-layout-year': calendarViewMode === 'year' }">
      <!-- 左侧：日历热力图 -->
      <div class="trend-calendar">
        <CalendarHeatmap
          ref="calendarRef"
          :data="heatmapData"
          :selected-date="selectedDate"
          :hovered-date="hoveredDate"
          @update:view-mode="calendarViewMode = $event"
          @update:month="handleCalendarMonthChange"
          @select-date="selectedDate = $event"
          @hover-date="calendarHoveredDate = $event"
        />
      </div>

      <!-- 右侧：趋势曲线 -->
      <section class="chart-section">
        <div class="chart-header">
          <h3>{{ t('history.assetValueTrend') }}</h3>
          <div class="chart-header-right">
            <!-- 计价货币选择器 -->
            <div class="currency-selector">
              <label class="selector-label">{{ t('dashboard.pricingCurrency') }}:</label>
              <div class="selector-buttons">
                <button
                  v-for="curr in currencies"
                  :key="curr.code"
                  class="selector-btn"
                  :class="{ active: pricingCurrency === curr.code }"
                  @click="setPricingCurrency(curr.code)"
                >
                  {{ curr.symbol }} {{ curr.code }}
                </button>
              </div>
            </div>

            <span class="current-value">{{ t('history.currentValue') }}: {{ formatCurrency(assetStore.totalValue) }}</span>
            <span v-if="showExchangeRate" class="rate-indicator">
              <span class="rate-dot" />
              {{ pricingCurrency }}/USDC
            </span>
          </div>
        </div>
        <div class="chart-container" @mouseleave="hoveredDate = null">
          <Line ref="chartRef" :data="lineChartData" :options="lineChartOptions" :plugins="[verticalLinePlugin]" :key="selectedTimeRange + pricingCurrency + (showExchangeRate ? 'rate' : '')" />
        </div>
      </section>
    </div>

    <!-- 历史资金变化记录 -->
    <section class="records-section">
      <div class="section-header">
        <h3 class="section-title">{{ t('history.fundChanges') || '资金变化记录' }}</h3>
        <Transition name="filter-fade">
          <button
            v-if="selectedDate"
            class="filter-badge"
            @click="selectedDate = null"
          >
            <PhCalendar :size="12" />
            <span class="filter-date">{{ selectedDateLabel }}</span>
            <span class="filter-count font-mono">{{ filteredRecords.length }}/{{ historyRecords.length }}</span>
            <span class="filter-clear">&times;</span>
          </button>
        </Transition>
      </div>

      <div v-if="filteredRecords.length === 0" class="empty-state">
        <span class="text-muted">{{ selectedDate ? `${selectedDateLabel} 无资金变化记录` : (t('history.noRecordsFound') || '暂无历史记录') }}</span>
        <button v-if="selectedDate" class="btn-ghost-sm" @click="selectedDate = null">查看全部</button>
      </div>

      <div v-else class="records-list">
        <div
          v-for="record in filteredRecords"
          :key="record.id"
          class="record-item"
        >
          <!-- 记录头部 -->
          <div class="record-header">
            <div class="record-source">
              <span class="source-type">{{ record.source }}</span>
              <span class="source-name">{{ record.sourceName }}</span>
            </div>
            <span class="record-time">{{ formatRelativeTime(record.timestamp) }}</span>
          </div>

          <!-- 资金变化列表 -->
          <div class="changes-list">
            <div
              v-for="(change, idx) in record.changes"
              :key="idx"
              class="change-item"
              :class="change.action"
            >
              <PhPlus v-if="change.action === 'increase'" :size="14" weight="bold" />
              <PhMinus v-if="change.action === 'decrease'" :size="14" weight="bold" />
              <span class="change-asset">{{ change.asset }}</span>
              <span class="change-amount font-mono">{{ formatNumber(change.amount) }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>
    </template>
  </div>
</template>
<style scoped>
/* Tab 切换 */
.history-tabs {
  display: flex;
  gap: var(--gap-xs);
  border-bottom: 1px solid var(--color-border);
  padding-bottom: var(--gap-xs);
}

.history-tab {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-muted);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.history-tab:hover {
  color: var(--color-text-primary);
}

.history-tab.active {
  color: var(--color-accent-primary);
  border-bottom-color: var(--color-accent-primary);
}

.history-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xl);
  max-width: 1400px;
  margin: 0 auto;
}

/* ================== 页面头部和筛选器 ================== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--gap-md);
  padding-bottom: var(--gap-md);
  border-bottom: 1px solid var(--color-border);
}

.header-filters {
  display: flex;
  gap: var(--gap-lg);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.filter-label {
  font-size: 12px;
  color: var(--color-text-muted);
  font-weight: 500;
  white-space: nowrap;
}

.type-filter-buttons {
  display: flex;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  padding: 3px;
  gap: 2px;
}

.filter-button {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 5px 10px;
  border-radius: var(--radius-xs);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.filter-button:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.filter-button.active {
  background: var(--color-accent-primary);
  color: var(--color-text-inverse);
}

.type-label {
  white-space: nowrap;
}

/* 时间范围下拉菜单 */
.time-range-dropdown {
  position: relative;
}

.range-display-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 5px 10px;
  border-radius: var(--radius-sm);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.range-display-btn:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

.arrow-rotated {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  right: 0;
  min-width: 150px;
  padding: 4px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  z-index: 200;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 7px 10px;
  border-radius: var(--radius-xs);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--transition-fast);
  text-align: left;
}

.dropdown-item:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dropdown-item-active {
  background: color-mix(in srgb, var(--color-accent-primary) 10%, transparent);
  color: var(--color-accent-primary);
}

.dropdown-item-active:hover {
  background: color-mix(in srgb, var(--color-accent-primary) 15%, transparent);
}

/* ================== 趋势区域：日历 + 曲线同行 ================== */
.trend-layout {
  display: grid;
  grid-template-columns: minmax(260px, 320px) 1fr;
  gap: var(--gap-lg);
  align-items: start;
}

.trend-calendar {
  position: sticky;
  top: var(--gap-lg);
}

/* 年视图时日历占全宽，堆叠布局 */
.trend-layout-year {
  grid-template-columns: 1fr;
}

.trend-layout-year .trend-calendar {
  position: static;
}

/* ================== 趋势图区 ================== */
.chart-section {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-md);
  flex-wrap: wrap;
}

.chart-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-heading);
}

.chart-header-right {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

/* 计价货币选择器 */
.currency-selector {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.selector-label {
  font-size: 11px;
  color: var(--color-text-muted);
  font-weight: 500;
  white-space: nowrap;
}

.selector-buttons {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.selector-btn {
  padding: 4px 8px;
  border-radius: var(--radius-xs);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.selector-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.selector-btn.active {
  background: var(--color-accent-primary);
  color: var(--color-text-inverse);
}

.current-value {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-accent-primary);
  font-family: var(--font-mono);
}

.rate-indicator {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 11px;
  color: var(--color-accent-secondary);
  font-family: var(--font-mono);
}

.rate-dot {
  width: 12px;
  height: 0;
  border-top: 2px dashed var(--color-accent-secondary);
}

.chart-container {
  height: 280px;
  cursor: crosshair;
}

/* ================== 历史资金变化记录 ================== */
.records-section {
  padding: var(--gap-lg);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-lg);
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-heading);
}

/* 过滤徽章 */
.filter-badge {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 8px 4px 6px;
  background: color-mix(in srgb, var(--color-accent-primary) 12%, var(--color-bg-tertiary));
  border: 1px solid color-mix(in srgb, var(--color-accent-primary) 25%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-accent-primary);
  font-size: 0.6875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.filter-badge:hover {
  background: color-mix(in srgb, var(--color-accent-primary) 18%, var(--color-bg-tertiary));
  border-color: color-mix(in srgb, var(--color-accent-primary) 40%, transparent);
}

.filter-date {
  font-weight: 600;
}

.filter-count {
  color: var(--color-text-muted);
  font-size: 0.625rem;
}

.filter-clear {
  font-size: 0.875rem;
  line-height: 1;
  opacity: 0.6;
  margin-left: 2px;
}

.filter-badge:hover .filter-clear {
  opacity: 1;
}

/* 过滤徽章动画 */
.filter-fade-enter-active,
.filter-fade-leave-active {
  transition: all 180ms ease;
}

.filter-fade-enter-from,
.filter-fade-leave-to {
  opacity: 0;
  transform: translateX(8px);
}

/* 空状态按钮 */
.btn-ghost-sm {
  margin-top: var(--gap-xs);
  padding: 4px 12px;
  font-size: 0.75rem;
  color: var(--color-accent-primary);
  background: transparent;
  border: 1px solid color-mix(in srgb, var(--color-accent-primary) 30%, transparent);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-ghost-sm:hover {
  background: color-mix(in srgb, var(--color-accent-primary) 10%, transparent);
  border-color: var(--color-accent-primary);
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.record-item {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-md);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  background: var(--color-bg-secondary);
  transition: border-color var(--transition-fast);
}

.record-item:hover {
  border-color: var(--color-border-hover);
}

.record-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  padding-bottom: var(--gap-xs);
  border-bottom: 1px solid var(--color-border);
}

.record-source {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.source-type {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.source-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.record-time {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
  font-family: var(--font-mono);
}

.changes-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap-xs);
}

.change-item {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.8125rem;
  padding: 6px 10px;
  border-radius: var(--radius-sm);
  background: var(--color-bg-tertiary);
  border: 1px solid transparent;
  transition: all var(--transition-fast);
}

.change-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.change-item.increase {
  color: var(--color-success);
  border-color: color-mix(in srgb, var(--color-success) 20%, transparent);
  background: color-mix(in srgb, var(--color-success) 5%, var(--color-bg-tertiary));
}

.change-item.decrease {
  color: var(--color-error);
  border-color: color-mix(in srgb, var(--color-error) 20%, transparent);
  background: color-mix(in srgb, var(--color-error) 5%, var(--color-bg-tertiary));
}

.change-asset {
  font-weight: 600;
  letter-spacing: 0.3px;
}

.change-amount {
  font-weight: 600;
  font-family: var(--font-mono);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-xl);
  color: var(--color-text-muted);
  font-size: 13px;
  text-align: center;
}

.event-timeline::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 1px;
  background: var(--color-border);
}

.event-item {
  display: flex;
  margin-bottom: var(--gap-lg);
  position: relative;
}

.event-item:last-child {
  margin-bottom: 0;
}

.event-date {
  position: absolute;
  left: -100px;
  width: 80px;
  text-align: right;
  font-size: 12px;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.event-date .date {
  font-weight: 500;
  color: var(--color-text-secondary);
}

.event-date .time {
  font-size: 11px;
}

.event-content {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-sm);
  position: relative;
  flex: 1;
}

.event-content::before {
  content: '';
  position: absolute;
  left: -24px;
  top: 4px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-accent-primary);
  border: 2px solid var(--color-bg-secondary);
  z-index: 1;
}

.event-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: color-mix(in srgb, var(--color-accent-primary) 12%, transparent);
  color: var(--color-accent-primary);
  flex-shrink: 0;
}

.event-details {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.event-description {
  font-size: 13px;
  color: var(--color-text-primary);
  margin-bottom: 2px;
  line-height: 1.4;
}

.event-value {
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 500;
}

.event-value.positive { color: var(--color-success); }
.event-value.negative { color: var(--color-error); }

/* ================== 响应式 ================== */
@media (max-width: 960px) {
  .trend-layout {
    grid-template-columns: 1fr;
  }

  .trend-calendar {
    position: static;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-filters {
    width: 100%;
    flex-direction: column;
    gap: var(--gap-md);
  }

  .filter-group {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
  }

  .type-filter-buttons {
    width: 100%;
    justify-content: stretch;
    flex-wrap: wrap;
  }

  .filter-button {
    flex: 1;
    justify-content: center;
  }

  .time-range-dropdown {
    width: 100%;
  }

  .range-display-btn {
    width: 100%;
    justify-content: space-between;
  }

  .event-timeline {
    padding-left: 0;
  }

  .event-timeline::before {
    left: 15px;
  }

  .event-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .event-date {
    position: static;
    width: auto;
    text-align: left;
    margin-left: 50px;
    margin-bottom: var(--gap-xs);
    flex-direction: row;
    align-items: center;
    gap: var(--gap-xs);
  }

  .event-date .time {
    display: none;
  }

  .event-content {
    padding-left: 50px;
  }

  .event-content::before {
    left: 10px;
  }
}

@media (max-width: 480px) {
  .type-filter-buttons {
    flex-direction: column;
  }
}
</style>
