<script setup>
/**
 * NFT 画廊组件
 * 网格展示 NFT 列表，支持收藏集分组和搜索
 */
import { ref, computed } from 'vue'
import { PhMagnifyingGlass } from '@phosphor-icons/vue'
import { useNFTStore } from '../stores/nftStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import NFTCard from './NFTCard.vue'

const nftStore = useNFTStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

// 搜索关键词
const searchQuery = ref('')

// 搜索过滤后的 NFT 列表
const filteredNFTs = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return nftStore.nfts
  return nftStore.nfts.filter(nft =>
    nft.name.toLowerCase().includes(q) ||
    nft.collection.toLowerCase().includes(q)
  )
})

// 按收藏集分组
const groupedByCollection = computed(() => {
  const groups = {}
  for (const nft of filteredNFTs.value) {
    if (!groups[nft.collection]) {
      groups[nft.collection] = {
        name: nft.collection,
        floorPrice: nft.floorPrice,
        floorCurrency: nft.floorCurrency,
        totalFloorUSD: 0,
        nfts: []
      }
    }
    groups[nft.collection].totalFloorUSD += nft.floorPriceUSD
    groups[nft.collection].nfts.push(nft)
  }
  return Object.values(groups).sort((a, b) => b.totalFloorUSD - a.totalFloorUSD)
})
</script>

<template>
  <div class="nft-gallery">
    <!-- 搜索栏 -->
    <div class="gallery-toolbar">
      <div class="search-box">
        <PhMagnifyingGlass :size="14" class="search-icon" />
        <input
          type="text"
          v-model="searchQuery"
          :placeholder="t('nft.searchPlaceholder')"
          class="search-input"
        />
      </div>
      <span class="gallery-count font-mono">
        {{ filteredNFTs.length }} {{ t('nft.items') }}
      </span>
    </div>

    <!-- 无数据提示 -->
    <div v-if="filteredNFTs.length === 0" class="empty-state">
      <p>{{ t('nft.noNFTs') }}</p>
    </div>

    <!-- 按收藏集分组展示 -->
    <div v-else class="collection-groups">
      <div
        v-for="group in groupedByCollection"
        :key="group.name"
        class="collection-group"
      >
        <!-- 收藏集头部 -->
        <div class="collection-header">
          <div class="collection-title">
            <span class="collection-name">{{ group.name }}</span>
            <span class="collection-count font-mono">{{ group.nfts.length }} {{ t('nft.items') }}</span>
          </div>
          <span class="collection-value font-mono">
            ${{ formatNumber(group.totalFloorUSD, 0) }}
          </span>
        </div>

        <!-- NFT 网格 -->
        <div class="nft-grid">
          <NFTCard
            v-for="nft in group.nfts"
            :key="nft.id"
            :nft="nft"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.nft-gallery {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* 工具栏 */
.gallery-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
}

.search-box {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: 6px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  flex: 1;
  max-width: 320px;
}

.search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.search-input {
  background: none;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-size: 0.8125rem;
  width: 100%;
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.gallery-count {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* 无数据 */
.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 收藏集分组 */
.collection-groups {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xl);
}

.collection-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
  margin-bottom: var(--gap-md);
}

.collection-title {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.collection-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.collection-count {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.collection-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* NFT 网格 */
.nft-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: var(--gap-md);
}
</style>
