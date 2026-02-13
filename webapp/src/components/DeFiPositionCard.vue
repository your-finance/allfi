<script setup>
/**
 * DeFi 仓位卡片组件
 * 展示单个 DeFi 协议仓位信息：协议、类型、代币、价值、APY、奖励
 */
import { computed } from 'vue'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  position: {
    type: Object,
    required: true
  }
})

const { formatNumber } = useFormatters()
const { t } = useI18n()

// 仓位类型标签
const typeLabel = computed(() => {
  const labels = {
    lp: t('defi.typeLp'),
    staking: t('defi.typeStaking'),
    lending: t('defi.typeLending')
  }
  return labels[props.position.type] || props.position.type
})

// 仓位类型样式类
const typeClass = computed(() => `type-${props.position.type}`)

// 链配置
const chainColors = {
  ETH: '#627EEA',
  BSC: '#F3BA2F',
  SOL: '#9945FF',
  MATIC: '#8247E5',
  ARB: '#28A0F0',
  OP: '#FF0420'
}

const chainColor = computed(() => chainColors[props.position.chain] || '#888')

// 代币列表展示文本
const tokenPairLabel = computed(() => {
  return props.position.tokens.map(t => t.symbol).join(' / ')
})
</script>

<template>
  <div class="defi-card">
    <!-- 头部：协议 + 类型标签 + 链标识 -->
    <div class="card-top">
      <div class="protocol-info">
        <span class="protocol-icon">{{ position.protocolIcon }}</span>
        <span class="protocol-name">{{ position.protocol }}</span>
      </div>
      <div class="card-badges">
        <span class="type-badge" :class="typeClass">{{ typeLabel }}</span>
        <span class="chain-badge" :style="{ borderColor: chainColor, color: chainColor }">
          {{ position.chain }}
        </span>
      </div>
    </div>

    <!-- 代币列表 -->
    <div class="token-pair">
      <span class="token-label">{{ tokenPairLabel }}</span>
    </div>

    <!-- 价值和 APY -->
    <div class="card-stats">
      <div class="stat-item">
        <span class="stat-label">{{ t('defi.value') }}</span>
        <span class="stat-value font-mono">${{ formatNumber(position.valueUSD) }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">APY</span>
        <span class="stat-value font-mono apy-value">{{ position.apy }}%</span>
      </div>
    </div>

    <!-- 奖励信息 -->
    <div v-if="position.rewards && position.rewards.valueUSD > 0" class="rewards-row">
      <span class="rewards-label">{{ t('defi.pendingRewards') }}</span>
      <span class="rewards-value font-mono">
        {{ formatNumber(position.rewards.amount, 4) }} {{ position.rewards.token }}
        <span class="rewards-usd">(≈${{ formatNumber(position.rewards.valueUSD) }})</span>
      </span>
    </div>

    <!-- LP 特有信息：价格区间 -->
    <div v-if="position.type === 'lp' && position.priceRange" class="extra-info">
      <span class="extra-label">{{ t('defi.priceRange') }}</span>
      <span class="extra-value font-mono">
        ${{ formatNumber(position.priceRange.lower) }} — ${{ formatNumber(position.priceRange.upper) }}
      </span>
      <span class="range-status" :class="position.inRange ? 'in-range' : 'out-range'">
        {{ position.inRange ? t('defi.inRange') : t('defi.outOfRange') }}
      </span>
    </div>

    <!-- 借贷特有信息：健康因子 -->
    <div v-if="position.type === 'lending' && position.healthFactor" class="extra-info">
      <span class="extra-label">{{ t('defi.healthFactor') }}</span>
      <span class="extra-value font-mono" :class="position.healthFactor >= 2 ? 'hf-safe' : 'hf-warn'">
        {{ position.healthFactor.toFixed(2) }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.defi-card {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: border-color var(--transition-fast);
}

.defi-card:hover {
  border-color: var(--color-accent-primary);
}

/* 头部 */
.card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-sm);
}

.protocol-info {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.protocol-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.625rem;
  font-weight: 700;
  color: var(--color-text-secondary);
}

.protocol-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.card-badges {
  display: flex;
  gap: var(--gap-xs);
}

.type-badge {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.type-badge.type-lp {
  background: rgba(139, 92, 246, 0.15);
  color: #8B5CF6;
}

.type-badge.type-staking {
  background: rgba(16, 185, 129, 0.15);
  color: #10B981;
}

.type-badge.type-lending {
  background: rgba(59, 130, 246, 0.15);
  color: #3B82F6;
}

.chain-badge {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  border: 1px solid;
}

/* 代币对 */
.token-pair {
  padding: var(--gap-xs) 0;
}

.token-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

/* 统计数据 */
.card-stats {
  display: flex;
  gap: var(--gap-lg);
  padding: var(--gap-sm) 0;
  border-top: 1px solid var(--color-border);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.stat-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.apy-value {
  color: var(--color-success);
}

/* 奖励 */
.rewards-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-sm);
  padding: var(--gap-xs) var(--gap-sm);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.rewards-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.rewards-value {
  font-size: 0.75rem;
  color: var(--color-text-primary);
}

.rewards-usd {
  color: var(--color-text-muted);
  font-size: 0.6875rem;
}

/* 额外信息 */
.extra-info {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  font-size: 0.75rem;
}

.extra-label {
  color: var(--color-text-muted);
  font-size: 0.6875rem;
}

.extra-value {
  color: var(--color-text-secondary);
}

.range-status {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 1px 4px;
  border-radius: var(--radius-xs);
}

.range-status.in-range {
  background: rgba(16, 185, 129, 0.15);
  color: var(--color-success);
}

.range-status.out-range {
  background: rgba(226, 92, 92, 0.15);
  color: var(--color-error);
}

.hf-safe {
  color: var(--color-success);
}

.hf-warn {
  color: var(--color-warning);
}

/* 响应式 */
@media (max-width: 768px) {
  .card-stats {
    gap: var(--gap-md);
  }
}
</style>
