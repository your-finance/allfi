<script setup>
/**
 * DeFi 仓位详情页
 * 完整展示所有 DeFi 仓位，支持筛选、排序、详细信息
 */
import { ref, computed } from 'vue'
import { PhVault, PhFunnel, PhArrowLeft, PhCoins, PhTrendUp, PhBank } from '@phosphor-icons/vue'
import { useRouter } from 'vue-router'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import { useAssetStore } from '../stores/assetStore'
import DeFiPositionCard from '../components/DeFiPositionCard.vue'

const router = useRouter()
const { formatNumber } = useFormatters()
const { t } = useI18n()
const assetStore = useAssetStore()

// 筛选类型
const filterType = ref('all')
const typeFilters = ['all', 'lp', 'staking', 'lending']

// DeFi 仓位数据
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
  return Object.values(groups).sort((a, b) => b.totalValue - a.totalValue)
})

// 统计数据
const totalValue = computed(() => assetStore.defiTotalValue)
const totalRewards = computed(() => {
  return positions.value.reduce((sum, p) => sum + (p.rewards?.valueUSD || 0), 0)
})
const typeStats = computed(() => {
  const stats = { lp: 0, staking: 0, lending: 0 }
  positions.value.forEach(p => {
    if (stats[p.type] !== undefined) stats[p.type]++
  })
  return stats
})

// 类型标签
const getTypeLabel = (type) => {
  const labels = {
    all: t('defi.filterAll'),
    lp: t('defi.typeLp'),
    staking: t('defi.typeStaking'),
    lending: t('defi.typeLending')
  }
  return labels[type] || type
}

const getTypeIcon = (type) => {
  const icons = { lp: PhCoins, staking: PhTrendUp, lending: PhBank }
  return icons[type] || PhVault
}

// 返回
const goBack = () => {
  router.push('/dashboard')
}
</script>

<template>
  <div class="defi-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="back-btn" @click="goBack">
        <PhArrowLeft :size="18" weight="bold" />
      </button>
      <div class="header-content">
        <div class="header-title">
          <div class="title-icon">
            <PhVault :size="20" weight="duotone" />
          </div>
          <h1>{{ t('defi.title') }}</h1>
        </div>
        <p class="header-desc">{{ t('defi.pageDesc') || '管理您的 DeFi 仓位和收益' }}</p>
      </div>
    </header>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card primary">
        <span class="stat-label">{{ t('defi.totalValue') }}</span>
        <span class="stat-value font-mono">${{ formatNumber(totalValue, 0) }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">{{ t('defi.positions') }}</span>
        <span class="stat-value font-mono">{{ positions.length }}</span>
      </div>
      <div class="stat-card rewards" v-if="totalRewards > 0">
        <span class="stat-label">{{ t('defi.pendingRewards') }}</span>
        <span class="stat-value font-mono">${{ formatNumber(totalRewards, 0) }}</span>
      </div>
      <div class="stat-card" v-else>
        <span class="stat-label">{{ t('defi.protocols') }}</span>
        <span class="stat-value font-mono">{{ groupedByProtocol.length }}</span>
      </div>
    </div>

    <!-- 类型统计 -->
    <div class="type-summary">
      <div
        v-for="(count, type) in typeStats"
        :key="type"
        class="type-item"
        :class="{ active: filterType === type }"
        @click="filterType = type"
      >
        <component :is="getTypeIcon(type)" :size="16" weight="duotone" />
        <span class="type-name">{{ getTypeLabel(type) }}</span>
        <span class="type-count font-mono">{{ count }}</span>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <div class="filter-bar">
        <PhFunnel :size="14" class="filter-icon" />
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
    </div>

    <!-- 仓位列表 -->
    <div v-if="groupedByProtocol.length > 0" class="protocol-list">
      <div
        v-for="group in groupedByProtocol"
        :key="group.protocol"
        class="protocol-section"
      >
        <div class="protocol-header">
          <div class="protocol-info">
            <span class="protocol-icon">{{ group.protocolIcon }}</span>
            <span class="protocol-name">{{ group.protocol }}</span>
          </div>
          <span class="protocol-total font-mono">${{ formatNumber(group.totalValue, 0) }}</span>
        </div>
        <div class="positions-grid">
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
      <PhVault :size="48" weight="duotone" class="empty-icon" />
      <p>{{ t('defi.noPositions') }}</p>
    </div>
  </div>
</template>

<style scoped>
.defi-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* 页面头部 */
.page-header {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-md);
}

.back-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.back-btn:hover {
  border-color: var(--color-accent-primary);
  color: var(--color-accent-primary);
}

.header-content {
  flex: 1;
}

.header-title {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.title-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%);
  border-radius: var(--radius-sm);
  color: white;
}

.header-title h1 {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.header-desc {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  margin-top: var(--gap-xs);
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--gap-md);
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.stat-card.primary {
  background: linear-gradient(135deg, color-mix(in srgb, #8B5CF6 15%, var(--color-bg-secondary)) 0%, var(--color-bg-secondary) 100%);
  border-color: color-mix(in srgb, #8B5CF6 30%, var(--color-border));
}

.stat-card.rewards .stat-value {
  color: var(--color-success);
}

.stat-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

/* 类型统计 */
.type-summary {
  display: flex;
  gap: var(--gap-sm);
}

.type-item {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s ease;
}

.type-item:hover {
  border-color: #8B5CF6;
}

.type-item.active {
  background: color-mix(in srgb, #8B5CF6 10%, var(--color-bg-secondary));
  border-color: #8B5CF6;
}

.type-item.active .type-name {
  color: #8B5CF6;
}

.type-name {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.type-count {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
  background: var(--color-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
}

/* 筛选器 */
.filter-section {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 3px;
  border-radius: var(--radius-sm);
}

.filter-icon {
  color: var(--color-text-muted);
  margin: 0 var(--gap-xs);
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: all 0.15s ease;
}

.filter-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.filter-btn.active {
  background: #8B5CF6;
  color: white;
}

.filter-count {
  font-size: 0.625rem;
  opacity: 0.7;
}

/* 协议列表 */
.protocol-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

.protocol-section {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.protocol-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
}

.protocol-info {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.protocol-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 700;
}

.protocol-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.protocol-total {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.positions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--gap-md);
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--gap-2xl);
  text-align: center;
  color: var(--color-text-muted);
  gap: var(--gap-md);
}

.empty-icon {
  opacity: 0.3;
  color: #8B5CF6;
}

/* 响应式 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .type-summary {
    flex-wrap: wrap;
  }

  .positions-grid {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-wrap: wrap;
  }
}
</style>
