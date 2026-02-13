<script setup>
/**
 * Reports 页面 - 资产报告
 * 三个标签：日报/周报 | 月报 | 报告对比
 * 月报含交互式图表，对比功能支持两个月数据比较
 */
import { ref, computed, onMounted } from 'vue'
import {
  PhCalendarBlank,
  PhCaretUp,
  PhCaretDown,
  PhArrowClockwise,
  PhChartLine,
  PhBookOpenText,
  PhArrowsLeftRight
} from '@phosphor-icons/vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Filler
} from 'chart.js'
import AnnualReport from '../components/AnnualReport.vue'
import AnnualReportShare from '../components/AnnualReportShare.vue'
import { reportService } from '../api/index.js'
import { annualReportService } from '../api/annualReportService.js'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler)

const { currencySymbol, formatNumber, formatPercent } = useFormatters()
const { t } = useI18n()

// 当前标签页
const activeTab = ref('timeline') // timeline | monthly | compare

// ========== 日报/周报 ==========
const reports = ref([])
const loading = ref(false)
const generating = ref(false)
const filterType = ref('')

// 年度报告
const showAnnualReport = ref(false)
const showAnnualShare = ref(false)
const annualReportData = ref(null)
const selectedYear = ref(new Date().getFullYear() - 1)

const openAnnualReport = async () => {
  showAnnualReport.value = true
  try {
    annualReportData.value = await annualReportService.getAnnualReport(selectedYear.value)
  } catch (e) {
    console.error('加载年度报告失败:', e)
  }
}

const openAnnualShare = () => {
  if (annualReportData.value) {
    showAnnualReport.value = false
    showAnnualShare.value = true
  }
}

const filteredReports = computed(() => {
  if (!filterType.value) return reports.value
  return reports.value.filter(r => r.type === filterType.value)
})

const loadReports = async () => {
  loading.value = true
  try {
    const result = await reportService.getReports('', 30)
    // 后端返回 { reports: [...] }，Mock 直接返回数组
    reports.value = Array.isArray(result) ? result : (result.reports || [])
  } catch (e) {
    console.error('加载报告失败:', e)
  } finally {
    loading.value = false
  }
}

const generateReport = async (type) => {
  generating.value = true
  try {
    const result = await reportService.generateReport(type)
    // 后端返回 { report: {...} }
    const report = result.report || result
    reports.value.unshift(report)
  } catch (e) {
    console.error('生成报告失败:', e)
  } finally {
    generating.value = false
  }
}

// ========== 月报 ==========
const monthlyReport = ref(null)
const monthlyLoading = ref(false)
const selectedMonth = ref((() => {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth()).padStart(2, '0') // 上月
  return `${y}-${m}`
})())

// 月份选项（最近 12 个月）
const monthOptions = computed(() => {
  const opts = []
  const now = new Date()
  for (let i = 1; i <= 12; i++) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const val = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
    const label = d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' })
    opts.push({ value: val, label })
  }
  return opts
})

const loadMonthlyReport = async () => {
  monthlyLoading.value = true
  try {
    monthlyReport.value = await reportService.getMonthlyReport(selectedMonth.value)
  } catch (e) {
    console.error('加载月报失败:', e)
  } finally {
    monthlyLoading.value = false
  }
}

// 月度累计收益图表配置
const monthlyChartData = computed(() => {
  if (!monthlyReport.value?.dailyReturns) return null
  const dr = monthlyReport.value.dailyReturns
  return {
    labels: dr.map(d => d.day),
    datasets: [{
      label: t('reports.cumReturn'),
      data: dr.map(d => d.cumReturnPct),
      borderColor: 'var(--color-accent-primary)',
      backgroundColor: 'rgba(75, 131, 240, 0.08)',
      borderWidth: 2,
      pointRadius: 0,
      pointHoverRadius: 4,
      pointHoverBorderWidth: 2,
      fill: true,
      tension: 0.3,
    }],
  }
})

const monthlyChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    mode: 'index',
    intersect: false,
  },
  plugins: {
    tooltip: {
      callbacks: {
        label: (ctx) => `${ctx.dataset.label}: ${ctx.parsed.y >= 0 ? '+' : ''}${ctx.parsed.y.toFixed(2)}%`,
      },
    },
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { font: { size: 10 }, maxTicksLimit: 10 },
    },
    y: {
      grid: { color: 'rgba(128,128,128,0.1)' },
      ticks: {
        font: { size: 10 },
        callback: (val) => `${val >= 0 ? '+' : ''}${val}%`,
      },
    },
  },
}

// ========== 报告对比 ==========
const compareData = ref(null)
const compareLoading = ref(false)
const compareReportId1 = ref('')
const compareReportId2 = ref('')

// 报告选项（从已加载报告列表生成）
const reportOptions = computed(() => {
  return reports.value.map(r => ({
    value: r.id,
    label: `${r.type === 'daily' ? '日报' : r.type === 'weekly' ? '周报' : r.type === 'monthly' ? '月报' : '年报'} - ${r.created_at?.slice(0, 10) || r.period || r.id}`
  }))
})

// 初始化对比报告
const initCompareMonths = () => {
  if (reportOptions.value.length >= 2) {
    compareReportId1.value = reportOptions.value[0].value
    compareReportId2.value = reportOptions.value[1].value
  }
}

const loadCompare = async () => {
  if (!compareReportId1.value || !compareReportId2.value) return
  compareLoading.value = true
  try {
    compareData.value = await reportService.compareReports(compareReportId1.value, compareReportId2.value)
  } catch (e) {
    console.error('加载对比数据失败:', e)
  } finally {
    compareLoading.value = false
  }
}

// 格式化时间
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const formatPeriod = (report) => {
  if (report.type === 'weekly') {
    return report.period
  }
  const d = new Date(report.period)
  return d.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })
}

// 月份标签
const monthLabel = (val) => {
  const opt = monthOptions.value.find(o => o.value === val)
  return opt ? opt.label : val
}

onMounted(() => {
  loadReports()
  initCompareMonths()
})
</script>

<template>
  <div class="reports-page">
    <!-- 页头 -->
    <section class="page-header">
      <div class="header-left">
        <PhChartLine :size="20" weight="bold" />
        <h2>{{ t('reports.title') }}</h2>
      </div>

      <!-- 标签切换 -->
      <div class="tab-group">
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'timeline' }"
          @click="activeTab = 'timeline'"
        >
          {{ t('reports.tabTimeline') }}
        </button>
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'monthly' }"
          @click="activeTab = 'monthly'; if (!monthlyReport) loadMonthlyReport()"
        >
          {{ t('reports.tabMonthly') }}
        </button>
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'compare' }"
          @click="activeTab = 'compare'"
        >
          {{ t('reports.tabCompare') }}
        </button>
      </div>
    </section>

    <!-- ========== 日报/周报 标签页 ========== -->
    <template v-if="activeTab === 'timeline'">
      <!-- 过滤和生成 -->
      <div class="toolbar">
        <div class="filter-group">
          <button
            class="filter-btn"
            :class="{ active: filterType === '' }"
            @click="filterType = ''"
          >
            {{ t('reports.all') }}
          </button>
          <button
            class="filter-btn"
            :class="{ active: filterType === 'daily' }"
            @click="filterType = 'daily'"
          >
            {{ t('reports.daily') }}
          </button>
          <button
            class="filter-btn"
            :class="{ active: filterType === 'weekly' }"
            @click="filterType = 'weekly'"
          >
            {{ t('reports.weekly') }}
          </button>
        </div>

        <button
          class="generate-btn"
          :disabled="generating"
          @click="generateReport('daily')"
        >
          <PhArrowClockwise :size="14" :class="{ spinning: generating }" />
          {{ t('reports.generate') }}
        </button>
      </div>

      <!-- 年度报告入口 -->
      <section class="annual-entry">
        <div class="annual-card" @click="openAnnualReport">
          <PhBookOpenText :size="24" class="annual-icon" />
          <div class="annual-info">
            <div class="annual-title">{{ selectedYear }} {{ t('annualReport.title') }}</div>
            <div class="annual-desc">{{ t('annualReport.entryDesc') }}</div>
          </div>
          <PhCaretDown :size="14" class="annual-arrow" />
        </div>
      </section>

      <!-- 报告列表 -->
      <section class="report-list">
        <div v-if="loading" class="loading-state">
          {{ t('common.loading') }}
        </div>

        <div v-else-if="filteredReports.length === 0" class="empty-state">
          {{ t('reports.empty') }}
        </div>

        <div
          v-else
          v-for="report in filteredReports"
          :key="report.id"
          class="report-card"
        >
          <div class="timeline-marker">
            <div class="marker-dot" :class="report.type === 'weekly' ? 'weekly' : 'daily'" />
            <div class="marker-line" />
          </div>

          <div class="report-content">
            <div class="report-header">
              <div class="report-meta">
                <span class="report-type-badge" :class="report.type">
                  {{ report.type === 'weekly' ? t('reports.weekly') : t('reports.daily') }}
                </span>
                <span class="report-period">{{ formatPeriod(report) }}</span>
              </div>
              <span class="report-date">
                <PhCalendarBlank :size="12" />
                {{ formatDate(report.generated_at) }}
              </span>
            </div>

            <div class="report-body">
              <div class="stat-row">
                <span class="stat-label">{{ t('reports.totalValue') }}</span>
                <span class="stat-value font-mono">
                  {{ currencySymbol }}{{ formatNumber(report.total_value) }}
                </span>
              </div>

              <div class="stat-row">
                <span class="stat-label">{{ t('reports.change') }}</span>
                <span
                  class="stat-value font-mono"
                  :class="report.change >= 0 ? 'change-positive' : 'change-negative'"
                >
                  <PhCaretUp v-if="report.change >= 0" :size="12" weight="bold" />
                  <PhCaretDown v-else :size="12" weight="bold" />
                  {{ report.change >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(Math.abs(report.change)) }}
                  ({{ report.change_percent >= 0 ? '+' : '' }}{{ formatPercent(report.change_percent) }})
                </span>
              </div>

              <div v-if="report.btc_benchmark || report.eth_benchmark" class="benchmark-row">
                <span class="benchmark-item" v-if="report.btc_benchmark">
                  BTC:
                  <span :class="report.btc_benchmark >= 0 ? 'change-positive' : 'change-negative'" class="font-mono">
                    {{ report.btc_benchmark >= 0 ? '+' : '' }}{{ formatPercent(report.btc_benchmark) }}
                  </span>
                </span>
                <span class="benchmark-item" v-if="report.eth_benchmark">
                  ETH:
                  <span :class="report.eth_benchmark >= 0 ? 'change-positive' : 'change-negative'" class="font-mono">
                    {{ report.eth_benchmark >= 0 ? '+' : '' }}{{ formatPercent(report.eth_benchmark) }}
                  </span>
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </template>

    <!-- ========== 月报标签页 ========== -->
    <template v-else-if="activeTab === 'monthly'">
      <div class="toolbar">
        <select
          class="month-select"
          v-model="selectedMonth"
          @change="loadMonthlyReport"
        >
          <option v-for="opt in monthOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div v-if="monthlyLoading" class="loading-state">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="monthlyReport" class="monthly-content">
        <!-- 月度收益概览 -->
        <div class="monthly-overview">
          <div class="overview-main">
            <div class="overview-label">{{ t('reports.monthlyReturn') }}</div>
            <div
              class="overview-value font-mono"
              :class="monthlyReport.totalReturn >= 0 ? 'change-positive' : 'change-negative'"
            >
              {{ monthlyReport.totalReturn >= 0 ? '+' : '' }}{{ monthlyReport.totalReturn }}%
            </div>
          </div>
          <div class="overview-grid">
            <div class="overview-item">
              <span class="ov-label">{{ t('reports.startValue') }}</span>
              <span class="ov-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.startValue) }}</span>
            </div>
            <div class="overview-item">
              <span class="ov-label">{{ t('reports.endValue') }}</span>
              <span class="ov-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.endValue) }}</span>
            </div>
            <div class="overview-item" v-if="monthlyReport.btcBenchmark != null">
              <span class="ov-label">BTC</span>
              <span
                class="ov-value font-mono"
                :class="monthlyReport.btcBenchmark >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ monthlyReport.btcBenchmark >= 0 ? '+' : '' }}{{ monthlyReport.btcBenchmark }}%
              </span>
            </div>
            <div class="overview-item" v-if="monthlyReport.ethBenchmark != null">
              <span class="ov-label">ETH</span>
              <span
                class="ov-value font-mono"
                :class="monthlyReport.ethBenchmark >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ monthlyReport.ethBenchmark >= 0 ? '+' : '' }}{{ monthlyReport.ethBenchmark }}%
              </span>
            </div>
          </div>
        </div>

        <!-- 累计收益曲线 -->
        <div class="chart-section">
          <h3 class="section-label">{{ t('reports.cumReturn') }}</h3>
          <div class="chart-container" v-if="monthlyChartData">
            <Line :data="monthlyChartData" :options="monthlyChartOptions" />
          </div>
        </div>

        <!-- 配置变化 -->
        <div class="allocation-section">
          <h3 class="section-label">{{ t('reports.allocationChange') }}</h3>
          <div class="alloc-compare">
            <div class="alloc-col">
              <div class="alloc-title">{{ t('reports.monthStart') }}</div>
              <div
                v-for="item in monthlyReport.allocation.start"
                :key="'s-' + item.symbol"
                class="alloc-bar-row"
              >
                <span class="alloc-symbol">{{ item.symbol }}</span>
                <div class="alloc-bar-track">
                  <div class="alloc-bar" :style="{ width: item.pct + '%' }" />
                </div>
                <span class="alloc-pct font-mono">{{ item.pct }}%</span>
              </div>
            </div>
            <div class="alloc-col">
              <div class="alloc-title">{{ t('reports.monthEnd') }}</div>
              <div
                v-for="item in monthlyReport.allocation.end"
                :key="'e-' + item.symbol"
                class="alloc-bar-row"
              >
                <span class="alloc-symbol">{{ item.symbol }}</span>
                <div class="alloc-bar-track">
                  <div class="alloc-bar" :style="{ width: item.pct + '%' }" />
                </div>
                <span class="alloc-pct font-mono">{{ item.pct }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 费用汇总 -->
        <div class="fee-section">
          <h3 class="section-label">{{ t('reports.feeSummary') }}</h3>
          <div class="fee-grid">
            <div class="fee-item">
              <span class="fee-label">{{ t('reports.totalFee') }}</span>
              <span class="fee-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.feeSummary.totalFee) }}</span>
            </div>
            <div class="fee-item">
              <span class="fee-label">{{ t('reports.tradingFee') }}</span>
              <span class="fee-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.feeSummary.tradingFee) }}</span>
            </div>
            <div class="fee-item">
              <span class="fee-label">{{ t('reports.gasFee') }}</span>
              <span class="fee-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.feeSummary.gasFee) }}</span>
            </div>
            <div class="fee-item">
              <span class="fee-label">{{ t('reports.withdrawFee') }}</span>
              <span class="fee-value font-mono">{{ currencySymbol }}{{ formatNumber(monthlyReport.feeSummary.withdrawFee) }}</span>
            </div>
          </div>
        </div>

        <!-- 建议 -->
        <div class="suggestions-section" v-if="monthlyReport.suggestions?.length">
          <h3 class="section-label">{{ t('reports.suggestions') }}</h3>
          <ul class="suggestion-list">
            <li v-for="(s, i) in monthlyReport.suggestions" :key="i">{{ s }}</li>
          </ul>
        </div>
      </div>
    </template>

    <!-- ========== 报告对比标签页 ========== -->
    <template v-else-if="activeTab === 'compare'">
      <div class="toolbar compare-toolbar">
        <select v-model="compareReportId1" class="month-select">
          <option v-for="opt in reportOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <PhArrowsLeftRight :size="18" class="compare-icon" />
        <select v-model="compareReportId2" class="month-select">
          <option v-for="opt in reportOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <button class="generate-btn" :disabled="compareLoading" @click="loadCompare">
          {{ t('reports.runCompare') }}
        </button>
      </div>

      <div v-if="compareLoading" class="loading-state">
        {{ t('common.loading') }}
      </div>

      <div v-else-if="compareData" class="compare-content">
        <div class="compare-table">
          <div class="compare-row compare-header-row">
            <div class="compare-cell label-cell">{{ t('reports.metric') }}</div>
            <div class="compare-cell">{{ compareData.report_1?.period || `#${compareData.report_1?.id}` }}</div>
            <div class="compare-cell">{{ compareData.report_2?.period || `#${compareData.report_2?.id}` }}</div>
          </div>

          <div class="compare-row">
            <div class="compare-cell label-cell">{{ t('reports.totalValue') }}</div>
            <div class="compare-cell font-mono">{{ currencySymbol }}{{ formatNumber(compareData.report_1?.total_value) }}</div>
            <div class="compare-cell font-mono">{{ currencySymbol }}{{ formatNumber(compareData.report_2?.total_value) }}</div>
          </div>

          <div class="compare-row">
            <div class="compare-cell label-cell">{{ t('reports.change') }}</div>
            <div class="compare-cell">
              <span
                class="font-mono"
                :class="(compareData.report_1?.change || 0) >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ (compareData.report_1?.change || 0) >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(compareData.report_1?.change) }}
              </span>
            </div>
            <div class="compare-cell">
              <span
                class="font-mono"
                :class="(compareData.report_2?.change || 0) >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ (compareData.report_2?.change || 0) >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(compareData.report_2?.change) }}
              </span>
            </div>
          </div>

          <div class="compare-row">
            <div class="compare-cell label-cell">{{ t('reports.changePercent') }}</div>
            <div class="compare-cell">
              <span
                class="font-mono"
                :class="(compareData.report_1?.change_percent || 0) >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ (compareData.report_1?.change_percent || 0) >= 0 ? '+' : '' }}{{ compareData.report_1?.change_percent?.toFixed(2) }}%
              </span>
            </div>
            <div class="compare-cell">
              <span
                class="font-mono"
                :class="(compareData.report_2?.change_percent || 0) >= 0 ? 'change-positive' : 'change-negative'"
              >
                {{ (compareData.report_2?.change_percent || 0) >= 0 ? '+' : '' }}{{ compareData.report_2?.change_percent?.toFixed(2) }}%
              </span>
            </div>
          </div>

          <div class="compare-row">
            <div class="compare-cell label-cell">{{ t('reports.valueDiff') }}</div>
            <div class="compare-cell font-mono" colspan="2">
              <span :class="(compareData.value_diff || 0) >= 0 ? 'change-positive' : 'change-negative'">
                {{ (compareData.value_diff || 0) >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(compareData.value_diff) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        {{ t('reports.compareHint') }}
      </div>
    </template>

    <!-- 年度报告弹窗 -->
    <AnnualReport
      :visible="showAnnualReport"
      :year="selectedYear"
      @close="showAnnualReport = false"
      @share="openAnnualShare"
    />

    <AnnualReportShare
      v-if="annualReportData"
      :visible="showAnnualShare"
      :report="annualReportData"
      @close="showAnnualShare = false"
    />
  </div>
</template>

<style scoped>
.reports-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  max-width: 800px;
}

/* 页头 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  color: var(--color-text-primary);
}

.header-left h2 {
  font-size: 1rem;
  font-weight: 600;
}

/* 标签切换 */
.tab-group {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.tab-btn {
  padding: 5px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.tab-btn:hover {
  color: var(--color-text-primary);
}

.tab-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* 工具栏 */
.toolbar {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

.compare-toolbar {
  align-items: center;
}

.compare-icon {
  color: var(--color-text-muted);
}

.month-select {
  padding: 5px 10px;
  font-size: 0.8125rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  cursor: pointer;
}

.filter-group {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.filter-btn {
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.filter-btn:hover {
  color: var(--color-text-primary);
}

.filter-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

.generate-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.generate-btn:hover:not(:disabled) {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

.generate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 年度报告入口 */
.annual-entry {
  margin-bottom: var(--gap-sm);
}

.annual-card {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md) var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.annual-card:hover {
  border-color: var(--color-accent-primary);
}

.annual-icon {
  color: var(--color-accent-primary);
  flex-shrink: 0;
}

.annual-info { flex: 1; }

.annual-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.annual-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-top: 2px;
}

.annual-arrow {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

/* 报告列表 */
.report-list {
  display: flex;
  flex-direction: column;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--gap-2xl);
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.report-card {
  display: flex;
  gap: var(--gap-md);
}

.timeline-marker {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 20px;
  flex-shrink: 0;
}

.marker-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--color-accent-primary);
  flex-shrink: 0;
  margin-top: var(--gap-md);
}

.marker-dot.weekly {
  width: 12px;
  height: 12px;
  background: var(--color-accent-secondary, var(--color-accent-primary));
}

.marker-line {
  width: 2px;
  flex: 1;
  background: var(--color-border);
  min-height: 16px;
}

.report-content {
  flex: 1;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-md) var(--gap-lg);
  margin-bottom: var(--gap-sm);
  transition: border-color var(--transition-fast);
}

.report-content:hover {
  border-color: var(--color-border-hover, var(--color-border));
}

.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-sm);
}

.report-meta {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.report-type-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: var(--radius-xs);
  text-transform: uppercase;
}

.report-type-badge.daily {
  background: rgba(75, 131, 240, 0.1);
  color: var(--color-accent-primary);
}

.report-type-badge.weekly {
  background: rgba(62, 168, 122, 0.1);
  color: var(--color-success);
}

.report-period {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.report-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.report-body {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-label {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.stat-value {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 2px;
}

.benchmark-row {
  display: flex;
  gap: var(--gap-md);
  padding-top: var(--gap-xs);
  border-top: 1px solid var(--color-border);
  margin-top: var(--gap-xs);
}

.benchmark-item {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* ========== 月报 ========== */
.monthly-content {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

.monthly-overview {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-lg);
}

.overview-main {
  text-align: center;
  margin-bottom: var(--gap-md);
}

.overview-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-xs);
}

.overview-value {
  font-size: 2rem;
  font-weight: 700;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--gap-md);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

.overview-item {
  text-align: center;
}

.ov-label {
  display: block;
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-bottom: 2px;
}

.ov-value {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 图表区域 */
.chart-section {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-lg);
}

.section-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.chart-container {
  height: 200px;
}

/* 配置变化 */
.allocation-section {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-lg);
}

.alloc-compare {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-lg);
}

.alloc-title {
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: var(--gap-sm);
}

.alloc-bar-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-xs);
}

.alloc-symbol {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  width: 40px;
  flex-shrink: 0;
}

.alloc-bar-track {
  flex: 1;
  height: 6px;
  background: var(--color-bg-tertiary);
  border-radius: 3px;
  overflow: hidden;
}

.alloc-bar {
  height: 100%;
  background: var(--color-accent-primary);
  border-radius: 3px;
  transition: width 0.4s ease;
}

.alloc-pct {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  width: 32px;
  text-align: right;
  flex-shrink: 0;
}

/* 费用汇总 */
.fee-section {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-lg);
}

.fee-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--gap-sm);
}

.fee-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap-xs) 0;
}

.fee-label {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.fee-value {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 建议 */
.suggestions-section {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--gap-lg);
}

.suggestion-list {
  margin: 0;
  padding-left: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.suggestion-list li {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

/* ========== 报告对比 ========== */
.compare-content {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.compare-table {
  display: flex;
  flex-direction: column;
}

.compare-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  border-bottom: 1px solid var(--color-border);
}

.compare-row:last-child {
  border-bottom: none;
}

.compare-header-row {
  background: var(--color-bg-tertiary);
}

.compare-header-row .compare-cell {
  font-weight: 600;
  font-size: 0.75rem;
  color: var(--color-text-primary);
}

.compare-cell {
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.compare-cell.label-cell {
  font-weight: 500;
  color: var(--color-text-primary);
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .overview-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .alloc-compare {
    grid-template-columns: 1fr;
  }

  .fee-grid {
    grid-template-columns: 1fr;
  }

  .compare-toolbar {
    flex-wrap: wrap;
  }
}
</style>
