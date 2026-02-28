<script setup>
/**
 * 借贷策略优化器组件
 * 展示最优借贷策略推荐
 */
import { ref, computed } from 'vue'
import { PhLightbulb, PhArrowRight, PhTrendUp, PhWarning } from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'

const props = defineProps({
  optimization: {
    type: Object,
    required: true
  }
})

const { formatNumber } = useFormatters()

// 风险等级配置
const riskConfig = computed(() => {
  const configs = {
    low: { color: 'var(--color-success)', label: '低风险', icon: '✓' },
    medium: { color: 'var(--color-warning)', label: '中等风险', icon: '⚠' },
    high: { color: 'var(--color-error)', label: '高风险', icon: '⚠' }
  }
  return configs[props.optimization.risk_level] || configs.low
})

// 操作类型标签
const getActionLabel = (action) => {
  const labels = {
    migrate: '迁移',
    rebalance: '再平衡',
    reduce_borrow: '减少借款'
  }
  return labels[action] || action
}

// 操作类型颜色
const getActionColor = (action) => {
  const colors = {
    migrate: 'var(--color-accent-primary)',
    rebalance: 'var(--color-info)',
    reduce_borrow: 'var(--color-warning)'
  }
  return colors[action] || 'var(--color-text-secondary)'
}
</script>

<template>
  <div class="lending-optimizer">
    <!-- 头部：总结 -->
    <div class="optimizer-header">
      <div class="header-icon">
        <PhLightbulb :size="20" weight="duotone" />
      </div>
      <div class="header-content">
        <h3 class="header-title">策略优化建议</h3>
        <p class="header-summary">{{ optimization.summary }}</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-label">当前仓位</span>
        <span class="stat-value font-mono">{{ optimization.current_positions.length }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">优化建议</span>
        <span class="stat-value font-mono">{{ optimization.recommendations.length }}</span>
      </div>
      <div class="stat-card highlight">
        <span class="stat-label">潜在收益提升</span>
        <span class="stat-value font-mono positive">${{ formatNumber(optimization.potential_gain, 2) }}</span>
        <span class="stat-hint">年化</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">风险等级</span>
        <span class="stat-value" :style="{ color: riskConfig.color }">
          {{ riskConfig.icon }} {{ riskConfig.label }}
        </span>
      </div>
    </div>

    <!-- 优化建议列表 -->
    <div v-if="optimization.recommendations.length > 0" class="recommendations">
      <h4 class="recommendations-title">具体建议</h4>
      <div class="recommendation-list">
        <div
          v-for="(rec, index) in optimization.recommendations"
          :key="index"
          class="recommendation-item"
        >
          <div class="recommendation-header">
            <span class="action-badge" :style="{ color: getActionColor(rec.action) }">
              {{ getActionLabel(rec.action) }}
            </span>
            <span v-if="rec.expected_gain > 0" class="gain-badge positive font-mono">
              +${{ formatNumber(rec.expected_gain, 2) }}/年
            </span>
          </div>
          <div class="recommendation-body">
            <div class="protocol-flow">
              <span class="protocol-name">{{ rec.from_protocol || rec.to_protocol }}</span>
              <PhArrowRight v-if="rec.from_protocol" :size="14" class="arrow-icon" />
              <span v-if="rec.from_protocol" class="protocol-name">{{ rec.to_protocol }}</span>
            </div>
            <div class="token-info">
              <span class="token-symbol">{{ rec.token }}</span>
              <span class="token-amount font-mono">{{ formatNumber(rec.amount, 4) }}</span>
            </div>
          </div>
          <p class="recommendation-reason">{{ rec.reason }}</p>
        </div>
      </div>
    </div>

    <!-- 无建议提示 -->
    <div v-else class="no-recommendations">
      <PhTrendUp :size="32" weight="duotone" class="empty-icon" />
      <p>当前策略已优化，暂无改进建议</p>
    </div>
  </div>
</template>

<style scoped>
.lending-optimizer {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 头部 */
.optimizer-header {
  display: flex;
  gap: var(--gap-md);
}

.header-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  border-radius: var(--radius-sm);
  color: white;
  flex-shrink: 0;
}

.header-content {
  flex: 1;
}

.header-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 4px;
}

.header-summary {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--gap-md);
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.stat-card.highlight {
  background: linear-gradient(135deg, color-mix(in srgb, #10B981 15%, var(--color-bg-tertiary)) 0%, var(--color-bg-tertiary) 100%);
  border: 1px solid color-mix(in srgb, #10B981 30%, var(--color-border));
}

.stat-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.stat-hint {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

/* 建议列表 */
.recommendations {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.recommendations-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.recommendation-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.recommendation-item {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
  padding: var(--gap-md);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.recommendation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.action-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.gain-badge {
  font-size: 0.75rem;
  font-weight: 600;
}

.recommendation-body {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.protocol-flow {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.protocol-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.arrow-icon {
  color: var(--color-text-muted);
}

.token-info {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.token-symbol {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.token-amount {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.recommendation-reason {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
  margin-top: 4px;
}

/* 空状态 */
.no-recommendations {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  gap: var(--gap-sm);
}

.empty-icon {
  opacity: 0.3;
  color: var(--color-success);
}

/* 响应式 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
