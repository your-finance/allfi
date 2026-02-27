<script setup>
/**
 * Beta 系数对比卡片
 * 展示资产组合相对于不同基准的 Beta 系数和相关性
 */
import { computed } from 'vue'
import { PhTrendUp, PhTrendDown } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
})

const { t } = useI18n()

// Beta 解读
const betaInterpretation = (beta) => {
  if (beta > 1.2) return { text: t('risk.betaHigh'), color: 'var(--color-warning)' }
  if (beta > 0.8) return { text: t('risk.betaModerate'), color: 'var(--color-text-secondary)' }
  return { text: t('risk.betaLow'), color: 'var(--color-success)' }
}
</script>

<template>
  <div class="beta-comparison-card">
    <h3 class="card-title">{{ t('risk.betaComparison') }}</h3>

    <!-- 组合 Beta -->
    <div class="portfolio-beta">
      <span class="beta-label">{{ t('risk.portfolioBeta') }}</span>
      <span class="beta-value font-mono">{{ data.portfolio_beta.toFixed(2) }}</span>
    </div>

    <!-- 基准对比列表 -->
    <div class="benchmark-list">
      <div
        v-for="bench in data.benchmarks"
        :key="bench.name"
        class="benchmark-item"
      >
        <div class="bench-header">
          <span class="bench-name">{{ bench.name }}</span>
          <span class="bench-beta font-mono">β = {{ bench.beta.toFixed(2) }}</span>
        </div>
        <div class="bench-stats">
          <div class="stat-item">
            <span class="stat-label">{{ t('risk.correlation') }}</span>
            <span class="stat-value font-mono">{{ (bench.correlation * 100).toFixed(0) }}%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 说明 -->
    <div class="beta-hint">
      {{ t('risk.betaHint') }}
    </div>
  </div>
</template>

<style scoped>
.beta-comparison-card {
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

.portfolio-beta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-md) var(--gap-lg);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  margin-bottom: var(--gap-md);
}

.beta-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.beta-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.benchmark-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-md);
}

.benchmark-item {
  padding: var(--gap-sm) var(--gap-md);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.benchmark-item:hover {
  background: var(--color-bg-tertiary);
}

.bench-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-xs);
}

.bench-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.bench-beta {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.bench-stats {
  display: flex;
  gap: var(--gap-md);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.stat-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.stat-value {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.beta-hint {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  padding-top: var(--gap-sm);
  border-top: 1px solid var(--color-border);
}
</style>
