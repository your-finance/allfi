<script setup>
/**
 * NFT 资产详情页
 * 完整展示所有 NFT 收藏，支持按收藏集分组、网格/列表视图
 */
import { ref, computed, onMounted } from 'vue'
import { PhImages, PhArrowLeft, PhSquaresFour, PhListBullets, PhSparkle } from '@phosphor-icons/vue'
import { useRouter } from 'vue-router'
import { useNFTStore } from '../stores/nftStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const nftStore = useNFTStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

// 视图模式
const viewMode = ref('grid') // 'grid' | 'list'

// 选中的收藏集筛选
const selectedCollection = ref('all')

// 收藏集列表
const collections = computed(() => nftStore.collections)

// 筛选后的 NFT
const filteredNFTs = computed(() => {
  if (selectedCollection.value === 'all') return nftStore.nfts
  return nftStore.nfts.filter(nft => nft.collection === selectedCollection.value)
})

// 返回
const goBack = () => {
  router.push('/dashboard')
}

onMounted(() => {
  if (nftStore.nfts.length === 0) {
    nftStore.fetchNFTs()
  }
  nftStore.initialize()
})
</script>

<template>
  <div class="nft-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <button class="back-btn" @click="goBack">
        <PhArrowLeft :size="18" weight="bold" />
      </button>
      <div class="header-content">
        <div class="header-title">
          <div class="title-icon">
            <PhImages :size="20" weight="duotone" />
          </div>
          <h1>{{ t('nft.title') }}</h1>
        </div>
        <p class="header-desc">{{ t('nft.pageDesc') || '浏览和管理您的 NFT 收藏' }}</p>
      </div>
      <div class="header-actions">
        <label class="include-toggle">
          <input
            type="checkbox"
            :checked="nftStore.includeInTotal"
            @change="nftStore.toggleIncludeInTotal()"
          />
          <span>{{ t('nft.includeInTotal') }}</span>
        </label>
      </div>
    </header>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card primary">
        <span class="stat-label">{{ t('nft.totalValue') }}</span>
        <span class="stat-value font-mono">${{ formatNumber(nftStore.totalFloorValue, 0) }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">{{ t('nft.totalCount') }}</span>
        <span class="stat-value font-mono">{{ nftStore.totalCount }}</span>
      </div>
      <div class="stat-card">
        <span class="stat-label">{{ t('nft.collectionCount') }}</span>
        <span class="stat-value font-mono">{{ nftStore.collectionCount }}</span>
      </div>
    </div>

    <!-- 收藏集筛选 + 视图切换 -->
    <div class="toolbar">
      <div class="collection-filter">
        <button
          class="collection-btn"
          :class="{ active: selectedCollection === 'all' }"
          @click="selectedCollection = 'all'"
        >
          {{ t('common.all') }}
          <span class="count font-mono">{{ nftStore.totalCount }}</span>
        </button>
        <button
          v-for="col in collections"
          :key="col.name"
          class="collection-btn"
          :class="{ active: selectedCollection === col.name }"
          @click="selectedCollection = col.name"
        >
          {{ col.name }}
          <span class="count font-mono">{{ col.count }}</span>
        </button>
      </div>
      <div class="view-toggles">
        <button
          class="view-btn"
          :class="{ active: viewMode === 'grid' }"
          @click="viewMode = 'grid'"
        >
          <PhSquaresFour :size="16" />
        </button>
        <button
          class="view-btn"
          :class="{ active: viewMode === 'list' }"
          @click="viewMode = 'list'"
        >
          <PhListBullets :size="16" />
        </button>
      </div>
    </div>

    <!-- NFT 网格视图 -->
    <div v-if="viewMode === 'grid'" class="nft-grid">
      <div
        v-for="nft in filteredNFTs"
        :key="nft.id"
        class="nft-card"
      >
        <div class="nft-image">
          <img v-if="nft.image" :src="nft.image" :alt="nft.name" />
          <div v-else class="nft-placeholder">
            <PhSparkle :size="32" weight="duotone" />
          </div>
        </div>
        <div class="nft-info">
          <span class="nft-collection">{{ nft.collection }}</span>
          <span class="nft-name">{{ nft.name }}</span>
          <div class="nft-price">
            <span class="floor-label">Floor</span>
            <span class="floor-value font-mono">{{ nft.floorPrice }} {{ nft.floorCurrency }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- NFT 列表视图 -->
    <div v-else class="nft-list">
      <div class="list-header">
        <span class="col-image"></span>
        <span class="col-name">{{ t('nft.name') }}</span>
        <span class="col-collection">{{ t('nft.collection') }}</span>
        <span class="col-floor">{{ t('nft.floorPrice') }}</span>
        <span class="col-value">{{ t('nft.value') }}</span>
      </div>
      <div
        v-for="nft in filteredNFTs"
        :key="nft.id"
        class="list-row"
      >
        <div class="col-image">
          <div class="list-thumb">
            <img v-if="nft.image" :src="nft.image" :alt="nft.name" />
            <PhSparkle v-else :size="16" weight="duotone" />
          </div>
        </div>
        <span class="col-name">{{ nft.name }}</span>
        <span class="col-collection">{{ nft.collection }}</span>
        <span class="col-floor font-mono">{{ nft.floorPrice }} {{ nft.floorCurrency }}</span>
        <span class="col-value font-mono">${{ formatNumber(nft.floorValueUSD || 0, 0) }}</span>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="filteredNFTs.length === 0" class="empty-state">
      <PhImages :size="48" weight="duotone" class="empty-icon" />
      <p>{{ t('nft.noNFTs') }}</p>
    </div>
  </div>
</template>

<style scoped>
.nft-page {
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
  background: linear-gradient(135deg, #EC4899 0%, #F472B6 100%);
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

.header-actions {
  flex-shrink: 0;
}

.include-toggle {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.include-toggle input {
  accent-color: #EC4899;
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
  background: linear-gradient(135deg, color-mix(in srgb, #EC4899 15%, var(--color-bg-secondary)) 0%, var(--color-bg-secondary) 100%);
  border-color: color-mix(in srgb, #EC4899 30%, var(--color-border));
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

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

.collection-filter {
  display: flex;
  gap: var(--gap-xs);
  flex-wrap: wrap;
}

.collection-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.15s ease;
}

.collection-btn:hover {
  border-color: #EC4899;
}

.collection-btn.active {
  background: color-mix(in srgb, #EC4899 10%, var(--color-bg-secondary));
  border-color: #EC4899;
  color: #EC4899;
}

.collection-btn .count {
  font-size: 0.625rem;
  background: var(--color-bg-tertiary);
  padding: 1px 5px;
  border-radius: var(--radius-xs);
}

.view-toggles {
  display: flex;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  padding: 2px;
}

.view-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.15s ease;
}

.view-btn.active {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

/* NFT 网格 */
.nft-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--gap-md);
}

.nft-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: all 0.2s ease;
  cursor: pointer;
}

.nft-card:hover {
  border-color: #EC4899;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(236, 72, 153, 0.1);
}

.nft-image {
  aspect-ratio: 1;
  background: var(--color-bg-tertiary);
  overflow: hidden;
}

.nft-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.nft-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  opacity: 0.3;
}

.nft-info {
  padding: var(--gap-sm);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nft-collection {
  font-size: 0.625rem;
  color: #EC4899;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.nft-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nft-price {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  margin-top: var(--gap-xs);
}

.floor-label {
  font-size: 0.5625rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
}

.floor-value {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* NFT 列表 */
.nft-list {
  display: flex;
  flex-direction: column;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.list-header {
  display: grid;
  grid-template-columns: 48px 1fr 1fr 100px 100px;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.list-row {
  display: grid;
  grid-template-columns: 48px 1fr 1fr 100px 100px;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  align-items: center;
  border-bottom: 1px solid var(--color-border);
  transition: background 0.15s ease;
}

.list-row:last-child {
  border-bottom: none;
}

.list-row:hover {
  background: var(--color-bg-tertiary);
}

.list-thumb {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-xs);
  overflow: hidden;
  background: var(--color-bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.list-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.col-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.col-collection {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.col-floor, .col-value {
  font-size: 0.75rem;
  color: var(--color-text-primary);
  text-align: right;
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
  color: #EC4899;
}

/* 响应式 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .nft-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .list-header,
  .list-row {
    grid-template-columns: 40px 1fr 80px;
  }

  .col-collection,
  .col-floor {
    display: none;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .collection-filter {
    overflow-x: auto;
    flex-wrap: nowrap;
    padding-bottom: var(--gap-xs);
  }
}
</style>
