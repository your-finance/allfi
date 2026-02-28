<script setup>
/**
 * 借贷仓位卡片组件
 * 展示单个借贷仓位的详细信息
 */
import { computed } from 'vue'
import { PhBank, PhTrendUp, PhTrendDown, PhWarning } from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'

const props = defineProps({
  position: {
    type: Object,
    required: true
  }
})

const { formatNumber } = useFormatters()

// 健康因子风险等级
const healthStatus = computed(() => {
  const hf = props.position.health_factor
  if (hf === 0 || !props.position.borrow_amount) return { level: 'safe', color: 'var(--color-text-muted)', label: '无借款' }
  if (hf >= 2.0) return { level: 'safe', color: 'var(--color-success)', label: '安全' }
  if (hf >= 1.5) return { level: 'medium', color: 'var(--color-warning)', label: '中等' }
  if (hf >= 1.2) return { level: 'high', color: 'var(--color-error)', label: '高风险' }
  return { level: 'critical', color: '#DC2626', label: '危险' }
})

// 净收益率颜色
const netApyColor = computed(() => {
  return props.position.net_apy >= 0 ? 'var(--color-success)' : 'var(--color-error)'
})
</script>

<template>
  <div class="lending-position-card">
    <!-- 头部：协议和链 -->
    <div class="card-header">
      <div class="protocol-info">
        <PhBank :size="16" weight="duotone" />
        <span class="protocol-name">{{ position.protocol }}</span>
        <span class="chain-badge">{{ position.chain }}</span>
      </div>
      <div class="health-badge" :style="{ color: healthStatus.color }">
        <PhWarning v-if="healthStatus.level !== 'safe'" :size="14" weight="fill" />
        <span>{{ healthStatus.label }}</span>
      </div>
    </div>

    <!-- 存款信息 -->
    <div class="position-section supply">
      <div class="section-header">
        <PhTrendUp :size="14" weight="duotone" class="section-icon supply-icon" />
        <span class="section-label">存款</span>
      </div>
      <div class="section-content">
        <div class="token-info">
          <span class="token-symbol">{{ position.supply_token }}</span>
          <span class="token-amount font-mono">{{ formatNumber(position.supply_amount, 4) }}</span>
        </div>
        <div class="value-info">
          <span class="value font-mono">${{ formatNumber(position.supply_value_usd, 2) }}</span>
          <span class="apy positive font-mono">{{ formatNumber(position.supply_apy, 2) }}% APY</span>
        </div>
      </div>
    </div>

    <!-- 借款信息 -->
    <div v-if="position.borrow_amount > 0" class="position-section borrow">
      <div class="section-header">
        <PhTrendDown :size="14" weight="duotone" class="section-icon borrow-icon" />
        <span class="section-label">借款</span>
      </div>
      <div class="section-content">
        <div class="token-info">
          <span class="token-symbol">{{ position.borrow_token }}</span>
          <span class="token-amount font-mono">{{ formatNumber(position.borrow_amount, 4) }}</span>
        </div>
        <div class="value-info">
          <span class="value font-mono">${{ formatNumber(position.borrow_value_usd, 2) }}</span>
          <span class="apy negative font-mono">{{ formatNumber(position.borrow_apy, 2) }}% APY</span>
        </div>
      </div>
    </div>

    <!-- 底部：关键指标 -->
    <div class="card-footer">
      <div class="metric">
        <span class="metric-label">健康因子</span>
        <span class="metric-value font-mono" :style="{ color: healthStatus.color }">
          {{ position.health_factor > 0 ? formatNumber(position.health_factor, 2) : '-' }}
        </span>
      </div>
      <div class="metric">
        <span class="metric-label">LTV</span>
        <span class="metric-value font-mono">{{ formatNumber(position.ltv, 1) }}%</span>
      </div>
      <div class="metric">
        <span class="metric-label">净收益</span>
        <span class="metric-value font-mono" :style="{ color: netApyColor }">
          {{ formatNumber(position.net_apy, 2) }}%
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.lending-position-card {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: border-color 0.15s ease;
}

.lending-position-card:hover {
  border-color: var(--color-accent-primary);
}

/* 头部 */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
}

.protocol-info {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.protocol-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.chain-badge {
  font-size: 0.625rem;
  font-weight: 500;
  color: var(--color-text-muted);
  background: var(--color-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  text-transform: uppercase;
}

.health-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.6875rem;
  font-weight: 600;
}

/* 仓位区块 */
.position-section {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 4px;
}

.section-label {
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.section-icon {
  opacity: 0.6;
}

.supply-icon {
  color: var(--color-success);
}

.borrow-icon {
  color: var(--color-error);
}

.section-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.token-info {
  display: flex;
  align-items: baseline;
  gap: var(--gap-xs);
}

.token-symbol {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.token-amount {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.value-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.value {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.apy {
  font-size: 0.75rem;
  font-weight: 500;
}

/* 底部指标 */
.card-footer {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-sm);
  padding-top: var(--gap-sm);
  border-top: 1px solid var(--color-border);
}

.metric {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.metric-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.metric-value {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}
</style>
