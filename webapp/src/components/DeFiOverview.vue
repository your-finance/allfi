<script setup>
/**
 * DeFi 仓位概览面板
 * Dashboard Widget：展示 DeFi 总价值、按协议分组的仓位列表、类型筛选器
 */
import { ref, computed } from 'vue'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import { useAssetStore } from '../stores/assetStore'
import DeFiPositionCard from './DeFiPositionCard.vue'

const { formatNumber } = useFormatters()
const { t } = useI18n()
const assetStore = useAssetStore()

// 当前筛选类型
const filterType = ref('all')
const typeFilters = ['all', 'lp', 'staking', 'lending']

// DeFi 仓位数据（从 assetStore 获取）
const positions = computed(() => assetStore.defiPositions)

// 筛选后的仓位
const filteredPositions = computed(() => {
  if (filterType.value === 'all') return positions.value
  return positions.value.filter(p => p.type === filterType.value)
})

// 按协议分组
const groupedByProtocol = computed(() => {
  const groups = {}
  for (const p of filteredPositions.value) {
    if (!groups[p.protocol]) {
      groups[p.protocol] = {
        protocol: p.protocol,
        protocolIcon: p.protocolIcon,
        positions: [],
        totalValue: 0
      }
    }
    groups[p.protocol].positions.push(p)
    groups[p.protocol].totalValue += p.valueUSD
  }
  // 按总价值降序排列
  return Object.values(groups).sort((a, b) => b.totalValue - a.totalValue)
})

// DeFi 总价值
const totalValue = computed(() => assetStore.defiTotalValue)

// 总奖励
const totalRewards = computed(() => {
  return positions.value.reduce((sum, p) => sum + (p.rewards?.valueUSD || 0), 0)
})

// 类型筛选标签
const getTypeLabel = (type) => {
  const labels = {
    all: t('defi.filterAll'),
    lp: t('defi.typeLp'),
    staking: t('defi.typeStaking'),
    lending: t('defi.typeLending')
  }
  return labels[type] || type
}
</script>

<template>
  <div class="defi-overview">
    <!-- 头部：标题 + 统计 -->
    <div class="overview-header">
      <h3>{{ t('defi.title') }}</h3>
      <div class="overview-stats">
        <div class="stat-chip">
          <span class="stat-chip-label">{{ t('defi.totalValue') }}</span>
          <span class="stat-chip-value font-mono">${{ formatNumber(totalValue) }}</span>
        </div>
        <div v-if="totalRewards > 0" class="stat-chip rewards-chip">
          <span class="stat-chip-label">{{ t('defi.pendingRewards') }}</span>
          <span class="stat-chip-value font-mono">${{ formatNumber(totalRewards) }}</span>
        </div>
      </div>
    </div>

    <!-- 类型筛选器 -->
    <div class="filter-bar">
      <button
        v-for="type in typeFilters"
        :key="type"
        class="filter-btn"
        :class="{ active: filterType === type }"
        @click="filterType = type"
      >
        {{ getTypeLabel(type) }}
        <span class="filter-count">
          {{ type === 'all' ? positions.length : positions.filter(p => p.type === type).length }}
        </span>
      </button>
    </div>

    <!-- 按协议分组的仓位列表 -->
    <div v-if="groupedByProtocol.length > 0" class="protocol-groups">
      <div
        v-for="group in groupedByProtocol"
        :key="group.protocol"
        class="protocol-group"
      >
        <div class="protocol-header">
          <span class="protocol-icon">{{ group.protocolIcon }}</span>
          <span class="protocol-name">{{ group.protocol }}</span>
          <span class="protocol-total font-mono">${{ formatNumber(group.totalValue) }}</span>
        </div>
        <div class="protocol-positions">
          <DeFiPositionCard
            v-for="pos in group.positions"
            :key="pos.id"
            :position="pos"
          />
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      {{ t('defi.noPositions') }}
    </div>
  </div>
</template>

<style scoped>
.defi-overview {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 头部 */
.overview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

.overview-header h3 {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.overview-stats {
  display: flex;
  gap: var(--gap-sm);
}

.stat-chip {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-xs) var(--gap-sm);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
}

.stat-chip-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.stat-chip-value {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.rewards-chip .stat-chip-value {
  color: var(--color-success);
}

/* 筛选器 */
.filter-bar {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
  align-self: flex-start;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 4px;
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

.filter-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.filter-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

.filter-count {
  font-size: 0.625rem;
  opacity: 0.7;
}

.filter-btn.active .filter-count {
  opacity: 1;
}

/* 协议分组 */
.protocol-groups {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.protocol-group {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.protocol-header {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding-bottom: var(--gap-xs);
  border-bottom: 1px solid var(--color-border);
}

.protocol-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
  font-size: 0.5625rem;
  font-weight: 700;
  color: var(--color-text-secondary);
}

.protocol-name {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
  flex: 1;
}

.protocol-total {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.protocol-positions {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--gap-sm);
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: var(--gap-xl);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* 响应式 */
@media (max-width: 768px) {
  .overview-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .protocol-positions {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-wrap: wrap;
  }

  .filter-btn {
    min-height: 36px;
    padding: 6px 12px;
  }
}
</style>
