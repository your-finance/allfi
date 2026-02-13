<script setup>
/**
 * 日历热力图组件（增强版）
 * 支持月视图（带盈亏数字）和年视图（12 个月概览）
 * iOS/macOS 日历风格，CSS Grid 布局
 */
import { ref, computed, watch } from 'vue'
import {
  PhCaretLeft,
  PhCaretRight,
  PhEye,
  PhEyeSlash,
  PhCalendarBlank,
  PhSquaresFour
} from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'
import { useFormatters } from '../composables/useFormatters'

const props = defineProps({
  // 日期 -> 变化百分比的映射，如 { '2026-02-01': 1.5, '2026-02-02': -0.8 }
  data: {
    type: Object,
    default: () => ({})
  },
  // 当前选中的日期（YYYY-MM-DD），用于高亮显示
  selectedDate: {
    type: String,
    default: null
  },
  // 图表 hover 中的日期（YYYY-MM-DD），用于临时高亮
  hoveredDate: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['update:viewMode', 'update:month', 'select-date', 'hover-date'])

const { t } = useI18n()
const { formatPercent, formatCurrency, currencySymbol } = useFormatters()

// === 视图模式 ===
const viewMode = ref('month') // 'month' | 'year'
const showDetails = ref(false) // 是否在月视图格子中显示盈亏数字

// 切换视图模式时通知父组件
const setViewMode = (mode) => {
  viewMode.value = mode
  emit('update:viewMode', mode)
}

// === 月份导航 ===
const monthOffset = ref(0)

// 当前显示的月份
const currentMonth = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() + monthOffset.value)
  d.setDate(1)
  return d
})

// 当前显示的年份（年视图用）
const currentYear = computed(() => currentMonth.value.getFullYear())

// 月份标题
const monthLabel = computed(() => {
  return currentMonth.value.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' })
})

// 星期标签
const weekdayLabels = computed(() => [
  t('heatmap.sun') || '日',
  t('heatmap.mon') || '一',
  t('heatmap.tue') || '二',
  t('heatmap.wed') || '三',
  t('heatmap.thu') || '四',
  t('heatmap.fri') || '五',
  t('heatmap.sat') || '六',
])

// 迷你星期标签（年视图用，单字）
const miniWeekdayLabels = computed(() =>
  weekdayLabels.value.map(w => w.charAt(0))
)

// 今天日期
const today = new Date()
today.setHours(0, 0, 0, 0)
const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`

// === 单月日历数据生成 ===
function generateMonthCells(year, month) {
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDayOfWeek = firstDay.getDay()

  const cells = []

  // 前面的空白格
  for (let i = 0; i < startDayOfWeek; i++) {
    cells.push({ type: 'empty' })
  }

  // 日期格子
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const cellDate = new Date(year, month, d)
    const isFuture = cellDate > today
    const isToday = dateStr === todayStr
    const changePercent = props.data[dateStr]
    const hasData = changePercent !== undefined && changePercent !== null && !isFuture

    cells.push({
      type: 'day',
      day: d,
      date: dateStr,
      isFuture,
      isToday,
      hasData,
      changePercent: hasData ? changePercent : null,
      color: hasData ? getColor(changePercent) : null,
    })
  }

  return cells
}

// 当前月的格子数据
const calendarCells = computed(() => {
  const year = currentMonth.value.getFullYear()
  const month = currentMonth.value.getMonth()
  return generateMonthCells(year, month)
})

// === 月度统计 ===
const monthStats = computed(() => {
  const cells = calendarCells.value.filter(c => c.hasData)
  const winDays = cells.filter(c => c.changePercent > 0).length
  const lossDays = cells.filter(c => c.changePercent < 0).length
  const flatDays = cells.filter(c => c.changePercent === 0).length
  const totalDays = cells.length
  const winRate = totalDays > 0 ? (winDays / totalDays * 100) : 0
  const totalChange = cells.reduce((sum, c) => sum + (c.changePercent || 0), 0)

  return { winDays, lossDays, flatDays, totalDays, winRate, totalChange }
})

// === 年视图：12 个月数据 ===
const yearMonths = computed(() => {
  const year = currentYear.value
  const months = []

  for (let m = 0; m < 12; m++) {
    const cells = generateMonthCells(year, m)
    const dataCells = cells.filter(c => c.hasData)
    const winDays = dataCells.filter(c => c.changePercent > 0).length
    const lossDays = dataCells.filter(c => c.changePercent < 0).length
    const totalChange = dataCells.reduce((sum, c) => sum + (c.changePercent || 0), 0)

    // 月份名称
    const monthDate = new Date(year, m, 1)
    const label = monthDate.toLocaleDateString('zh-CN', { month: 'short' })

    months.push({
      index: m,
      label,
      cells,
      winDays,
      lossDays,
      totalDays: dataCells.length,
      totalChange,
    })
  }

  return months
})

// === 颜色映射 ===
function getColor(percent) {
  if (percent > 3) return 'rgba(34, 197, 94, 0.42)'
  if (percent > 0) return 'rgba(34, 197, 94, 0.20)'
  if (percent === 0) return 'var(--color-bg-tertiary)'
  if (percent > -3) return 'rgba(239, 68, 68, 0.20)'
  return 'rgba(239, 68, 68, 0.42)'
}

// === Tooltip ===
const hoveredCell = ref(null)
const tooltipPos = ref({ x: 0, y: 0 })

const showTooltip = (cell, event) => {
  if (!cell.hasData) return
  hoveredCell.value = cell
  const rect = event.target.getBoundingClientRect()
  tooltipPos.value = {
    x: rect.left + rect.width / 2,
    y: rect.top - 8,
  }
  // 通知父组件当前 hover 的日期（用于图表联动高亮）
  emit('hover-date', cell.date)
}

const hideTooltip = () => {
  hoveredCell.value = null
  emit('hover-date', null)
}

// 点击日期格子：选中/取消选中
const selectDate = (cell) => {
  if (!cell.hasData) return
  // 点击同一天取消选中，否则选中新日期
  const newDate = props.selectedDate === cell.date ? null : cell.date
  emit('select-date', newDate)
}

// === 导航 ===

// 通知父组件当前显示的月份已改变
function emitMonthChange() {
  const d = new Date()
  d.setMonth(d.getMonth() + monthOffset.value)
  emit('update:month', {
    year: d.getFullYear(),
    month: d.getMonth(), // 0-indexed
    monthOffset: monthOffset.value,
  })
}

const prevMonth = () => {
  monthOffset.value--
  emitMonthChange()
}
const nextMonth = () => {
  if (monthOffset.value < 0) {
    monthOffset.value++
    emitMonthChange()
  }
}
const canGoNext = computed(() => monthOffset.value < 0)

// 年视图前后年切换
const prevYear = () => {
  monthOffset.value -= 12
  emitMonthChange()
}
const nextYear = () => {
  if (monthOffset.value < -11) monthOffset.value += 12
  else monthOffset.value = 0
  emitMonthChange()
}
const canGoNextYear = computed(() => monthOffset.value < -11)

// 年视图点击月份 → 进入月视图
const goToMonth = (monthIndex) => {
  const targetDate = new Date(currentYear.value, monthIndex, 1)
  const now = new Date()
  now.setDate(1)
  now.setHours(0, 0, 0, 0)
  const diff = (targetDate.getFullYear() - now.getFullYear()) * 12 + (targetDate.getMonth() - now.getMonth())
  monthOffset.value = diff
  setViewMode('month')
  emitMonthChange()
}

// 允许父组件重置月份偏移（时间范围变更时）
const resetToCurrentMonth = () => {
  monthOffset.value = 0
}

// 导航到指定日期所在月份（供图表点击联动使用）
const navigateToDate = (dateStr) => {
  if (!dateStr) return
  const [y, m] = dateStr.split('-').map(Number)
  const now = new Date()
  const offset = (y - now.getFullYear()) * 12 + (m - 1 - now.getMonth())
  monthOffset.value = offset
  // 如果在年视图，切换回月视图
  if (viewMode.value === 'year') {
    setViewMode('month')
  }
}

// 暴露方法给父组件
defineExpose({ resetToCurrentMonth, navigateToDate })

// 格式化百分比（紧凑显示）
function fmtPct(val) {
  if (val === null || val === undefined) return ''
  const abs = Math.abs(val)
  if (abs >= 10) return `${val > 0 ? '+' : ''}${val.toFixed(0)}`
  return `${val > 0 ? '+' : ''}${val.toFixed(1)}`
}
</script>

<template>
  <div class="cal-card">
    <!-- 头部工具栏 -->
    <div class="cal-toolbar">
      <!-- 视图切换（iOS 风格 Segmented Control） -->
      <div class="cal-segmented">
        <button
          class="seg-btn"
          :class="{ active: viewMode === 'month' }"
          @click="setViewMode('month')"
        >
          <PhCalendarBlank :size="13" />
          {{ t('heatmap.monthView') || '月' }}
        </button>
        <button
          class="seg-btn"
          :class="{ active: viewMode === 'year' }"
          @click="setViewMode('year')"
        >
          <PhSquaresFour :size="13" />
          {{ t('heatmap.yearView') || '年' }}
        </button>
      </div>

      <!-- 月视图时显示详情开关 -->
      <button
        v-if="viewMode === 'month'"
        class="cal-toggle-btn"
        :class="{ active: showDetails }"
        @click="showDetails = !showDetails"
        :title="showDetails ? '隐藏盈亏数字' : '显示盈亏数字'"
      >
        <PhEye v-if="showDetails" :size="14" />
        <PhEyeSlash v-else :size="14" />
        <span class="toggle-label">{{ showDetails ? (t('heatmap.hideNumbers') || '隐藏') : (t('heatmap.showNumbers') || '数值') }}</span>
      </button>
    </div>

    <!-- ============ 月视图 ============ -->
    <template v-if="viewMode === 'month'">
      <!-- 月份导航 -->
      <div class="cal-nav">
        <button class="cal-nav-btn" @click="prevMonth">
          <PhCaretLeft :size="14" />
        </button>
        <span class="cal-nav-title">{{ monthLabel }}</span>
        <button class="cal-nav-btn" :disabled="!canGoNext" @click="nextMonth">
          <PhCaretRight :size="14" />
        </button>
      </div>

      <!-- 星期标签 -->
      <div class="cal-weekdays">
        <span v-for="(day, i) in weekdayLabels" :key="i" class="cal-wd">{{ day }}</span>
      </div>

      <!-- 月历网格 -->
      <div class="cal-grid" :class="{ 'cal-grid-detail': showDetails }">
        <div
          v-for="(cell, i) in calendarCells"
          :key="i"
          class="cal-cell"
          :class="{
            'cell-empty': cell.type === 'empty',
            'cell-future': cell.isFuture,
            'cell-no-data': cell.type === 'day' && !cell.hasData && !cell.isFuture,
            'cell-has-data': cell.hasData,
            'cell-today': cell.isToday,
            'cell-positive': cell.hasData && cell.changePercent > 0,
            'cell-negative': cell.hasData && cell.changePercent < 0,
            'cell-selected': cell.date === selectedDate,
            'cell-chart-hover': cell.date === hoveredDate && cell.date !== selectedDate,
          }"
          :style="cell.hasData ? { background: cell.color } : {}"
          @mouseenter="showTooltip(cell, $event)"
          @mouseleave="hideTooltip"
          @click="selectDate(cell)"
        >
          <template v-if="cell.type === 'day'">
            <span class="cell-day" :class="{ 'today-ring': cell.isToday }">{{ cell.day }}</span>
            <span
              v-if="showDetails && cell.hasData"
              class="cell-pct font-mono"
              :class="cell.changePercent >= 0 ? 'pct-up' : 'pct-down'"
            >{{ fmtPct(cell.changePercent) }}%</span>
          </template>
        </div>
      </div>

      <!-- 月度统计条 -->
      <div v-if="monthStats.totalDays > 0" class="cal-stats">
        <div class="stat-item">
          <span class="stat-dot stat-dot-win" />
          <span class="stat-label">{{ t('heatmap.winDays') || '盈' }}</span>
          <span class="stat-val font-mono">{{ monthStats.winDays }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-dot stat-dot-loss" />
          <span class="stat-label">{{ t('heatmap.lossDays') || '亏' }}</span>
          <span class="stat-val font-mono">{{ monthStats.lossDays }}</span>
        </div>
        <div class="stat-divider" />
        <div class="stat-item">
          <span class="stat-label">{{ t('heatmap.winRate') || '胜率' }}</span>
          <span class="stat-val font-mono" :class="monthStats.winRate >= 50 ? 'pct-up' : 'pct-down'">
            {{ monthStats.winRate.toFixed(0) }}%
          </span>
        </div>
        <div class="stat-divider" />
        <div class="stat-item">
          <span class="stat-label">{{ t('heatmap.totalPnl') || '累计' }}</span>
          <span class="stat-val font-mono" :class="monthStats.totalChange >= 0 ? 'pct-up' : 'pct-down'">
            {{ monthStats.totalChange >= 0 ? '+' : '' }}{{ monthStats.totalChange.toFixed(2) }}%
          </span>
        </div>
      </div>

      <!-- 图例 -->
      <div class="cal-legend">
        <span class="legend-label">{{ t('heatmap.loss') || '亏损' }}</span>
        <div class="legend-colors">
          <span class="legend-swatch" style="background: rgba(239, 68, 68, 0.42)" />
          <span class="legend-swatch" style="background: rgba(239, 68, 68, 0.20)" />
          <span class="legend-swatch" style="background: var(--color-bg-tertiary)" />
          <span class="legend-swatch" style="background: rgba(34, 197, 94, 0.20)" />
          <span class="legend-swatch" style="background: rgba(34, 197, 94, 0.42)" />
        </div>
        <span class="legend-label">{{ t('heatmap.profit') || '盈利' }}</span>
      </div>
    </template>

    <!-- ============ 年视图 ============ -->
    <template v-if="viewMode === 'year'">
      <!-- 年份导航 -->
      <div class="cal-nav">
        <button class="cal-nav-btn" @click="prevYear">
          <PhCaretLeft :size="14" />
        </button>
        <span class="cal-nav-title">{{ currentYear }}{{ t('heatmap.yearSuffix') || '年' }}</span>
        <button class="cal-nav-btn" :disabled="!canGoNextYear" @click="nextYear">
          <PhCaretRight :size="14" />
        </button>
      </div>

      <!-- 12 个月网格 -->
      <div class="year-grid">
        <div
          v-for="m in yearMonths"
          :key="m.index"
          class="year-month-card"
          @click="goToMonth(m.index)"
        >
          <!-- 月标题 -->
          <div class="ym-header">
            <span class="ym-label">{{ m.label }}</span>
            <span
              v-if="m.totalDays > 0"
              class="ym-pnl font-mono"
              :class="m.totalChange >= 0 ? 'pct-up' : 'pct-down'"
            >{{ m.totalChange >= 0 ? '+' : '' }}{{ m.totalChange.toFixed(1) }}%</span>
          </div>

          <!-- 迷你星期标签 -->
          <div class="ym-weekdays">
            <span v-for="(wd, i) in miniWeekdayLabels" :key="i" class="ym-wd">{{ wd }}</span>
          </div>

          <!-- 迷你日历网格 -->
          <div class="ym-grid">
            <div
              v-for="(cell, ci) in m.cells"
              :key="ci"
              class="ym-cell"
              :class="{
                'ym-empty': cell.type === 'empty',
                'ym-future': cell.isFuture,
                'ym-today': cell.isToday,
              }"
              :style="cell.hasData ? { background: cell.color } : {}"
            />
          </div>

          <!-- 迷你统计 -->
          <div v-if="m.totalDays > 0" class="ym-stats">
            <span class="ym-stat">
              <span class="stat-dot stat-dot-win" />{{ m.winDays }}
            </span>
            <span class="ym-stat">
              <span class="stat-dot stat-dot-loss" />{{ m.lossDays }}
            </span>
          </div>
        </div>
      </div>
    </template>

    <!-- Tooltip（Teleport 到 body） -->
    <Teleport to="body">
      <Transition name="tooltip-fade">
        <div
          v-if="hoveredCell"
          class="heatmap-tooltip"
          :style="{ left: tooltipPos.x + 'px', top: tooltipPos.y + 'px' }"
        >
          <span class="tooltip-date">{{ hoveredCell.date }}</span>
          <span
            class="tooltip-change font-mono"
            :class="hoveredCell.changePercent >= 0 ? 'positive' : 'negative'"
          >
            {{ hoveredCell.changePercent >= 0 ? '+' : '' }}{{ hoveredCell.changePercent.toFixed(2) }}%
          </span>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
/* === 卡片容器 === */
.cal-card {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* === 工具栏 === */
.cal-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-md);
}

/* iOS 风格 Segmented Control */
.cal-segmented {
  display: flex;
  padding: 2px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  gap: 1px;
}

.seg-btn {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 4px 10px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-muted);
  background: transparent;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.seg-btn.active {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

.seg-btn:hover:not(.active) {
  color: var(--color-text-secondary);
}

/* 详情开关按钮 */
.cal-toggle-btn {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 4px 8px;
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  background: var(--color-bg-tertiary);
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.cal-toggle-btn:hover {
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

.cal-toggle-btn.active {
  color: var(--color-accent-primary);
  border-color: color-mix(in srgb, var(--color-accent-primary) 30%, transparent);
  background: color-mix(in srgb, var(--color-accent-primary) 8%, var(--color-bg-tertiary));
}

.toggle-label {
  white-space: nowrap;
}

/* === 月份/年份导航 === */
.cal-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-md);
  margin-bottom: var(--gap-sm);
}

.cal-nav-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  min-width: 100px;
  text-align: center;
}

.cal-nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.cal-nav-btn:hover:not(:disabled) {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

.cal-nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* === 星期标签 === */
.cal-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  margin-bottom: 3px;
}

.cal-wd {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  text-align: center;
  padding: 2px 0;
  font-weight: 500;
}

/* === 月视图网格 === */
.cal-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 3px;
}

/* 详情模式：格子更高以容纳百分比 */
.cal-grid-detail .cal-cell:not(.cell-empty) {
  aspect-ratio: auto;
  min-height: 44px;
}

/* 普通模式：正方形格子 */
.cal-grid:not(.cal-grid-detail) .cal-cell {
  aspect-ratio: 1;
}

.cal-cell {
  border-radius: var(--radius-xs);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: default;
  position: relative;
  transition: opacity var(--transition-fast);
  gap: 1px;
}

.cell-empty {
  background: transparent;
}

.cell-future {
  background: transparent;
  opacity: 0.25;
}

.cell-no-data {
  background: var(--color-bg-tertiary);
}

.cell-has-data {
  cursor: pointer;
}

.cell-has-data:hover:not(.cell-selected) {
  opacity: 0.75;
}

/* 选中状态：白色实线边框 + 微放大 */
.cell-selected {
  outline: 2px solid var(--color-text-primary);
  outline-offset: -1px;
  z-index: 2;
  transform: scale(1.08);
}

/* 图表 hover 联动：虚线边框 + 轻微放大 */
.cell-chart-hover {
  outline: 1.5px dashed var(--color-accent-primary);
  outline-offset: -1px;
  z-index: 1;
  transform: scale(1.05);
  transition: outline var(--transition-fast), transform var(--transition-fast);
}

/* 今天指示器 */
.cell-today:not(.cell-selected):not(.cell-chart-hover) {
  outline: 1.5px solid var(--color-accent-primary);
  outline-offset: -1px;
  z-index: 1;
}

/* 日期数字 */
.cell-day {
  font-size: 0.625rem;
  font-weight: 500;
  color: var(--color-text-primary);
  line-height: 1;
}

.cell-no-data .cell-day {
  color: var(--color-text-muted);
}

.cell-future .cell-day {
  color: var(--color-text-muted);
}

/* 今天日期加粗 */
.cell-today .cell-day {
  font-weight: 600;
  color: var(--color-accent-primary);
}

/* 盈亏百分比文字 */
.cell-pct {
  font-size: 0.5rem;
  font-weight: 600;
  line-height: 1;
  letter-spacing: -0.02em;
}

.pct-up {
  color: var(--color-success);
}

.pct-down {
  color: var(--color-error);
}

/* === 月度统计条 === */
.cal-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  margin-top: var(--gap-sm);
  padding: var(--gap-xs) var(--gap-sm);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 3px;
}

.stat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.stat-dot-win {
  background: var(--color-success);
}

.stat-dot-loss {
  background: var(--color-error);
}

.stat-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.stat-val {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.stat-divider {
  width: 1px;
  height: 12px;
  background: var(--color-border);
}

/* === 图例 === */
.cal-legend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-xs);
  margin-top: var(--gap-sm);
}

.legend-label {
  font-size: 0.5625rem;
  color: var(--color-text-muted);
}

.legend-colors {
  display: flex;
  gap: 2px;
}

.legend-swatch {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}

/* === 年视图 === */
.year-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--gap-sm);
}

.year-month-card {
  padding: var(--gap-sm);
  background: var(--color-bg-tertiary);
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.year-month-card:hover {
  border-color: var(--color-border-hover, var(--color-border));
  background: var(--color-bg-elevated);
}

/* 年视图月标题 */
.ym-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.ym-label {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.ym-pnl {
  font-size: 0.5625rem;
  font-weight: 600;
}

/* 年视图迷你星期标签 */
.ym-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  margin-bottom: 2px;
}

.ym-wd {
  font-size: 0.5rem;
  color: var(--color-text-muted);
  text-align: center;
  line-height: 1.2;
}

/* 年视图迷你网格 */
.ym-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1.5px;
}

.ym-cell {
  aspect-ratio: 1;
  border-radius: 1px;
  background: var(--color-bg-secondary);
}

.ym-empty {
  background: transparent;
}

.ym-future {
  background: transparent;
  opacity: 0.2;
}

.ym-today {
  outline: 1px solid var(--color-accent-primary);
  z-index: 1;
}

/* 年视图迷你统计 */
.ym-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-xs);
  margin-top: 4px;
  font-size: 0.5625rem;
  font-family: var(--font-mono);
  color: var(--color-text-muted);
}

.ym-stat {
  display: flex;
  align-items: center;
  gap: 2px;
}

/* === Tooltip 过渡动画 === */
.tooltip-fade-enter-active,
.tooltip-fade-leave-active {
  transition: opacity 120ms ease;
}

.tooltip-fade-enter-from,
.tooltip-fade-leave-to {
  opacity: 0;
}

/* === 响应式 === */
@media (max-width: 480px) {
  .year-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>

<style>
/* Tooltip（非 scoped，因为 Teleport 到 body） */
.heatmap-tooltip {
  position: fixed;
  transform: translate(-50%, -100%);
  padding: 6px 10px;
  background: var(--color-bg-elevated, #1e1e2e);
  border: 1px solid var(--color-border, #333);
  border-radius: var(--radius-sm, 4px);
  z-index: 500;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  pointer-events: none;
}

.tooltip-date {
  font-size: 0.625rem;
  color: var(--color-text-muted, #888);
}

.tooltip-change {
  font-size: 0.75rem;
  font-weight: 600;
}

.tooltip-change.positive {
  color: var(--color-success, #22c55e);
}

.tooltip-change.negative {
  color: var(--color-error, #ef4444);
}
</style>
