<script setup>
/**
 * DeFi 迷你卡片 — Dashboard 紧凑预览
 * 显示 DeFi 总价值、仓位数量、待领奖励，点击跳转详情页
 */
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { PhVault, PhCoins, PhArrowRight, PhTrendUp } from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import { useAssetStore } from '../stores/assetStore'

const router = useRouter()
const { formatNumber } = useFormatters()
const { t } = useI18n()
const assetStore = useAssetStore()

// DeFi 数据
const positions = computed(() => assetStore.defiPositions)
const totalValue = computed(() => assetStore.defiTotalValue)
const totalRewards = computed(() => {
  return positions.value.reduce((sum, p) => sum + (p.rewards?.valueUSD || 0), 0)
})

// 按类型统计
const typeStats = computed(() => {
  const stats = { lp: 0, staking: 0, lending: 0 }
  positions.value.forEach(p => {
    if (stats[p.type] !== undefined) stats[p.type]++
  })
  return stats
})

// 跳转详情页
const goToDetail = () => {
  router.push('/defi')
}
</script>

<template>
  <div class="defi-mini-card" @click="goToDetail">
    <!-- 头部 -->
    <div class="card-header">
      <div class="header-left">
        <div class="icon-box">
          <PhVault :size="16" weight="duotone" />
        </div>
        <span class="card-title">{{ t('defi.title') }}</span>
      </div>
      <button class="view-more">
        <span>{{ t('common.viewAll') }}</span>
        <PhArrowRight :size="12" weight="bold" />
      </button>
    </div>

    <!-- 主要数值 -->
    <div class="main-value">
      <span class="value font-mono">${{ formatNumber(totalValue, 0) }}</span>
      <span class="label">{{ t('defi.totalValue') }}</span>
    </div>

    <!-- 统计行 -->
    <div class="stats-row">
      <div class="stat-item">
        <span class="stat-num font-mono">{{ positions.length }}</span>
        <span class="stat-label">{{ t('defi.positions') }}</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item" v-if="totalRewards > 0">
        <span class="stat-num font-mono rewards">${{ formatNumber(totalRewards, 0) }}</span>
        <span class="stat-label">{{ t('defi.pendingRewards') }}</span>
      </div>
      <div class="stat-item" v-else>
        <span class="stat-num font-mono">{{ typeStats.lp + typeStats.staking }}</span>
        <span class="stat-label">{{ t('defi.activeProtocols') }}</span>
      </div>
    </div>

    <!-- 类型标签 -->
    <div class="type-tags">
      <span v-if="typeStats.lp > 0" class="type-tag lp">
        <PhCoins :size="10" />
        LP {{ typeStats.lp }}
      </span>
      <span v-if="typeStats.staking > 0" class="type-tag staking">
        <PhTrendUp :size="10" />
        Staking {{ typeStats.staking }}
      </span>
      <span v-if="typeStats.lending > 0" class="type-tag lending">
        Lending {{ typeStats.lending }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.defi-mini-card {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  height: 100%;
}

.defi-mini-card:hover {
  border-color: color-mix(in srgb, var(--color-border) 50%, #8B5CF6 50%);
  background: color-mix(in srgb, var(--color-bg-secondary) 95%, #8B5CF6 5%);
}

/* 头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.icon-box {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%);
  border-radius: var(--radius-xs);
  color: white;
}

.card-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.view-more {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  font-size: 0.625rem;
  cursor: pointer;
  transition: color 0.15s ease;
}

.view-more:hover {
  color: #8B5CF6;
}

/* 主要数值 */
.main-value {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.main-value .value {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: -0.02em;
}

.main-value .label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

/* 统计行 */
.stats-row {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-xs) 0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.stat-num {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.stat-num.rewards {
  color: var(--color-success);
}

.stat-label {
  font-size: 0.5625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.stat-divider {
  width: 1px;
  height: 24px;
  background: var(--color-border);
}

/* 类型标签 */
.type-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: auto;
}

.type-tag {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 6px;
  font-size: 0.5625rem;
  font-weight: 500;
  border-radius: var(--radius-xs);
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.type-tag.lp {
  background: color-mix(in srgb, #3B82F6 15%, transparent);
  color: #3B82F6;
}

.type-tag.staking {
  background: color-mix(in srgb, #8B5CF6 15%, transparent);
  color: #8B5CF6;
}

.type-tag.lending {
  background: color-mix(in srgb, #10B981 15%, transparent);
  color: #10B981;
}
</style>
