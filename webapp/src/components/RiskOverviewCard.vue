<script setup>
/**
 * 风险总览卡片
 * 展示核心风险指标：波动率、夏普比率、最大回撤、VaR、CVaR
 */
import { computed } from 'vue'
import { PhTrendUp, PhTrendDown, PhWarning } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
})

const { t } = useI18n()

// 风险等级判断
const riskLevel = computed(() => {
  const vol = props.data.volatility
  if (vol < 15) return { level: 'low', color: 'var(--color-success)' }
  if (vol < 25) return { level: 'medium', color: 'var(--color-warning)' }
  return { level: 'high', color: 'var(--color-error)' }
})

// 夏普比率评级
const sharpeRating = computed(() => {
  const sr = props.data.sharpe_ratio
  if (sr >= 2) return 'excellent'
  if (sr >= 1) return 'good'
  if (sr >= 0.5) return 'fair'
  return 'poor'
})
</script>

<template>
  <div class="risk-overview-card">
    <h3 class="card-title">{{ t('risk.overview') }}</h3>

    <div class="metrics-grid">
      <!-- 波动率 -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.volatility') }}</span>
        <span class="metric-value font-mono" :style="{ color: riskLevel.color }">
          {{ data.volatility.toFixed(1) }}%
        </span>
        <span class="metric-hint" :style="{ color: riskLevel.color }">
          {{ t(`risk.riskLevel.${riskLevel.level}`) }}
        </span>
      </div>

      <!-- 夏普比率 -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.sharpeRatio') }}</span>
        <span class="metric-value font-mono">
          {{ data.sharpe_ratio.toFixed(2) }}
        </span>
        <span class="metric-hint">
          {{ t(`risk.sharpeRating.${sharpeRating}`) }}
        </span>
      </div>

      <!-- 最大回撤 -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.maxDrawdown') }}</span>
        <span class="metric-value font-mono negative">
          {{ data.max_drawdown.toFixed(1) }}%
        </span>
      </div>

      <!-- VaR 95% -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.var95') }}</span>
        <span class="metric-value font-mono negative">
          {{ data.var_95.toFixed(1) }}%
        </span>
      </div>

      <!-- CVaR 95% -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.cvar95') }}</span>
        <span class="metric-value font-mono negative">
          {{ data.cvar_95.toFixed(1) }}%
        </span>
      </div>

      <!-- Beta BTC -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.betaBTC') }}</span>
        <span class="metric-value font-mono">
          {{ data.beta_btc.toFixed(2) }}
        </span>
      </div>

      <!-- Beta ETH -->
      <div class="metric-item">
        <span class="metric-label">{{ t('risk.betaETH') }}</span>
        <span class="metric-value font-mono">
          {{ data.beta_eth.toFixed(2) }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.risk-overview-card {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.card-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-lg);
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.metric-value {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.metric-value.negative {
  color: var(--color-error);
}

.metric-hint {
  font-size: 0.625rem;
  font-weight: 500;
}

/* 响应式 */
@media (max-width: 768px) {
  .metrics-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--gap-md);
  }
}
</style>

