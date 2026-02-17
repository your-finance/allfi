<script setup>
/**
 * NFT 迷你卡片 — Dashboard 紧凑预览
 * 显示 NFT 总估值、数量、收藏集数，点击跳转详情页
 */
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { PhImages, PhArrowRight, PhSparkle } from '@phosphor-icons/vue'
import { useNFTStore } from '../stores/nftStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const nftStore = useNFTStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

// 热门收藏集 Top2
const topCollections = computed(() => nftStore.collections.slice(0, 2))

// 跳转详情页
const goToDetail = () => {
  router.push('/nft')
}

onMounted(() => {
  if (nftStore.nfts.length === 0) {
    nftStore.fetchNFTs()
  }
  nftStore.initialize()
})
</script>

<template>
  <div class="nft-mini-card" @click="goToDetail">
    <!-- 头部 -->
    <div class="card-header">
      <div class="header-left">
        <div class="icon-box">
          <PhImages :size="16" weight="duotone" />
        </div>
        <span class="card-title">{{ t('nft.title') }}</span>
      </div>
      <button class="view-more">
        <span>{{ t('common.viewAll') }}</span>
        <PhArrowRight :size="12" weight="bold" />
      </button>
    </div>

    <!-- 主要数值 -->
    <div class="main-value">
      <span class="value font-mono">${{ formatNumber(nftStore.totalFloorValue, 0) }}</span>
      <span class="label">{{ t('nft.totalValue') }}</span>
    </div>

    <!-- 统计行 -->
    <div class="stats-row">
      <div class="stat-item">
        <span class="stat-num font-mono">{{ nftStore.totalCount }}</span>
        <span class="stat-label">{{ t('nft.items') }}</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-num font-mono">{{ nftStore.collectionCount }}</span>
        <span class="stat-label">{{ t('nft.collections') }}</span>
      </div>
    </div>

    <!-- 热门收藏集预览 -->
    <div v-if="topCollections.length > 0" class="top-preview">
      <div
        v-for="col in topCollections"
        :key="col.name"
        class="collection-chip"
      >
        <PhSparkle :size="10" weight="fill" class="chip-icon" />
        <span class="chip-name">{{ col.name }}</span>
        <span class="chip-count font-mono">{{ col.count }}</span>
      </div>
    </div>

    <!-- 计入总资产开关 -->
    <label class="include-toggle" @click.stop>
      <input
        type="checkbox"
        :checked="nftStore.includeInTotal"
        @change="nftStore.toggleIncludeInTotal()"
      />
      <span class="toggle-text">{{ t('nft.includeInTotal') }}</span>
    </label>
  </div>
</template>

<style scoped>
.nft-mini-card {
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

.nft-mini-card:hover {
  border-color: color-mix(in srgb, var(--color-border) 50%, #EC4899 50%);
  background: color-mix(in srgb, var(--color-bg-secondary) 95%, #EC4899 5%);
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
  background: linear-gradient(135deg, #EC4899 0%, #F472B6 100%);
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
  color: #EC4899;
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

/* 热门收藏集预览 */
.top-preview {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.collection-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xs);
  font-size: 0.625rem;
}

.chip-icon {
  color: #EC4899;
}

.chip-name {
  color: var(--color-text-secondary);
  max-width: 60px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chip-count {
  color: var(--color-text-muted);
  font-weight: 500;
}

/* 计入总资产开关 */
.include-toggle {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  cursor: pointer;
  margin-top: auto;
  padding-top: var(--gap-xs);
  border-top: 1px solid var(--color-border);
}

.include-toggle input[type="checkbox"] {
  width: 12px;
  height: 12px;
  accent-color: #EC4899;
}

.toggle-text {
  font-size: 0.5625rem;
  color: var(--color-text-muted);
}
</style>
