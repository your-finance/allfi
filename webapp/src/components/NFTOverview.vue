<script setup>
/**
 * NFT 概览组件 — Dashboard 小组件
 * 显示 NFT 总估值、收藏集数量、热门收藏集 Top3、计入总资产开关
 */
import { onMounted } from 'vue'
import { useNFTStore } from '../stores/nftStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const nftStore = useNFTStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

// 热门收藏集 Top3
const topCollections = () => nftStore.collections.slice(0, 3)

onMounted(() => {
  if (nftStore.nfts.length === 0) {
    nftStore.fetchNFTs()
  }
  nftStore.initialize()
})
</script>

<template>
  <section class="nft-overview">
    <!-- 标题栏 -->
    <div class="section-header">
      <h3 class="section-title">{{ t('nft.title') }}</h3>
      <label class="toggle-label">
        <input
          type="checkbox"
          :checked="nftStore.includeInTotal"
          @change="nftStore.toggleIncludeInTotal()"
        />
        <span class="toggle-text">{{ t('nft.includeInTotal') }}</span>
      </label>
    </div>

    <!-- 统计摘要 -->
    <div class="stats-row">
      <div class="stat-item">
        <span class="stat-label">{{ t('nft.totalValue') }}</span>
        <span class="stat-value font-mono">${{ formatNumber(nftStore.totalFloorValue, 0) }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">{{ t('nft.totalCount') }}</span>
        <span class="stat-value font-mono">{{ nftStore.totalCount }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">{{ t('nft.collectionCount') }}</span>
        <span class="stat-value font-mono">{{ nftStore.collectionCount }}</span>
      </div>
    </div>

    <!-- 热门收藏集 Top3 -->
    <div v-if="topCollections().length > 0" class="top-collections">
      <h4 class="sub-title">{{ t('nft.topCollections') }}</h4>
      <div class="collection-list">
        <div
          v-for="(col, idx) in topCollections()"
          :key="col.name"
          class="collection-row"
        >
          <span class="collection-rank font-mono">{{ idx + 1 }}</span>
          <div class="collection-info">
            <span class="collection-name">{{ col.name }}</span>
            <span class="collection-detail font-mono">
              {{ col.count }} {{ t('nft.items') }} · {{ col.floorPrice }} {{ col.floorCurrency }}
            </span>
          </div>
          <span class="collection-value font-mono">${{ formatNumber(col.totalFloorUSD, 0) }}</span>
        </div>
      </div>
    </div>

    <!-- 无数据 -->
    <div v-else-if="!nftStore.isLoading" class="empty-hint">
      {{ t('nft.noNFTs') }}
    </div>
  </section>
</template>

<style scoped>
.nft-overview {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 标题栏 */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  cursor: pointer;
}

.toggle-label input[type="checkbox"] {
  accent-color: var(--color-accent-primary);
}

.toggle-text {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 统计摘要 */
.stats-row {
  display: flex;
  gap: var(--gap-lg);
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
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 热门收藏集 */
.sub-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-sm);
}

.collection-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.collection-row {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.collection-rank {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-muted);
  width: 20px;
  text-align: center;
}

.collection-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.collection-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.collection-detail {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.collection-value {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  flex-shrink: 0;
}

/* 无数据 */
.empty-hint {
  text-align: center;
  padding: var(--gap-md);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}
</style>
