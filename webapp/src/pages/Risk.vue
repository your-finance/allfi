<script setup>
/**
 * 风险管理页面
 * 展示资产组合的风险指标、回撤分析、Beta 系数等
 */
import { ref, onMounted, watch } from 'vue'
import RiskOverviewCard from '../components/RiskOverviewCard.vue'
import RiskMetricsChart from '../components/RiskMetricsChart.vue'
import DrawdownChart from '../components/DrawdownChart.vue'
import BetaComparisonCard from '../components/BetaComparisonCard.vue'
import RiskAlertPanel from '../components/RiskAlertPanel.vue'
import { riskService } from '../api/riskService.js'
import { useI18n } from '../composables/useI18n'

const { t } = useI18n()

// 时间范围
const selectedPeriod = ref('30d')
const periods = ['7d', '30d', '90d', '1y']

// 数据状态
const isLoading = ref(false)
const riskOverview = ref(null)
const riskMetrics = ref(null)
const drawdownData = ref(null)
const betaData = ref(null)
const riskAlerts = ref([])

// 加载所有数据
const loadAllData = async () => {
  isLoading.value = true
  try {
    const [overview, metrics, drawdown, beta, alerts] = await Promise.all([
      riskService.getRiskOverview(selectedPeriod.value),
      riskService.getRiskMetrics(selectedPeriod.value),
      riskService.getDrawdown(selectedPeriod.value),
      riskService.getBeta(selectedPeriod.value),
      riskService.getRiskAlerts()
    ])

    riskOverview.value = overview
    riskMetrics.value = metrics
    drawdownData.value = drawdown
    betaData.value = beta
    riskAlerts.value = alerts
  } catch (error) {
    console.error('加载风险数据失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 切换时间周期
watch(selectedPeriod, () => {
  loadAllData()
})

onMounted(() => {
  loadAllData()
})
</script>

<template>
  <div class="risk-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">{{ t('nav.risk') }}</h2>
      <div class="period-selector">
        <button
          v-for="period in periods"
          :key="period"
          class="period-btn"
          :class="{ active: selectedPeriod === period }"
          @click="selectedPeriod = period"
        >
          {{ t(`risk.period_${period}`) }}
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-state">
      {{ t('common.loading') }}
    </div>

    <!-- 主内容 -->
    <template v-else-if="riskOverview">
      <!-- 风险总览 -->
      <RiskOverviewCard :data="riskOverview" />

      <!-- 第一行：风险指标趋势 + 回撤曲线 -->
      <div class="panels-row">
        <RiskMetricsChart v-if="riskMetrics" :data="riskMetrics" />
        <DrawdownChart v-if="drawdownData" :data="drawdownData" />
      </div>

      <!-- 第二行：Beta 对比 + 风险预警 -->
      <div class="panels-row">
        <BetaComparisonCard v-if="betaData" :data="betaData" />
        <RiskAlertPanel :alerts="riskAlerts" />
      </div>
    </template>

    <!-- 无数据状态 -->
    <div v-else class="empty-state">
      <p>{{ t('risk.noData') }}</p>
    </div>
  </div>
</template>

<style scoped>
.risk-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  max-width: 1400px;
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
}

.page-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.period-selector {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.period-btn {
  padding: 4px 10px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.period-btn:hover {
  color: var(--color-text-primary);
}

.period-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* 面板行 */
.panels-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-lg);
}

/* 加载和空状态 */
.loading-state,
.empty-state {
  padding: var(--gap-2xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .panels-row {
    grid-template-columns: 1fr;
  }
}
</style>
