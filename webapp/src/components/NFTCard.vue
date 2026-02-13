<script setup>
/**
 * NFT 卡片组件
 * 显示单个 NFT 信息：图片占位、名称、收藏集、地板价、稀有度
 */
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  nft: { type: Object, required: true }
})

const { formatNumber } = useFormatters()
const { t } = useI18n()

// 稀有度颜色映射
const rarityColor = {
  Legendary: '#FF6B35',
  Epic: '#A855F7',
  Rare: '#3B82F6',
  Uncommon: '#10B981',
  Common: '#6B7280'
}
</script>

<template>
  <div class="nft-card">
    <!-- 图片占位区域 -->
    <div class="nft-image">
      <img v-if="nft.imageUrl" :src="nft.imageUrl" :alt="nft.name" />
      <div v-else class="image-placeholder">
        <span class="placeholder-text">{{ nft.collection?.charAt(0) || '?' }}</span>
      </div>
      <!-- 稀有度标签 -->
      <span
        v-if="nft.rarity"
        class="rarity-badge"
        :style="{ background: rarityColor[nft.rarity] || '#6B7280' }"
      >
        {{ nft.rarity }}
      </span>
    </div>

    <!-- 信息区域 -->
    <div class="nft-info">
      <div class="nft-name">{{ nft.name }}</div>
      <div class="nft-collection">{{ nft.collection }}</div>

      <div class="nft-meta">
        <!-- 地板价 -->
        <div class="meta-item">
          <span class="meta-label">{{ t('nft.floorPrice') }}</span>
          <span class="meta-value font-mono">
            {{ formatNumber(nft.floorPrice, 2) }} {{ nft.floorCurrency }}
          </span>
        </div>

        <!-- 链标签 -->
        <div class="meta-item">
          <span class="meta-label">{{ t('nft.chain') }}</span>
          <span class="chain-tag">{{ nft.chain }}</span>
        </div>
      </div>

      <!-- Floor Price USD -->
      <div class="nft-value font-mono">
        ${{ formatNumber(nft.floorPriceUSD, 0) }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.nft-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: border-color var(--transition-fast);
}

.nft-card:hover {
  border-color: var(--color-accent-primary);
}

/* 图片区域 */
.nft-image {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  background: var(--color-bg-tertiary);
  overflow: hidden;
}

.nft-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
}

.placeholder-text {
  font-size: 2rem;
  font-weight: 600;
  color: var(--color-text-muted);
  opacity: 0.5;
}

.rarity-badge {
  position: absolute;
  top: var(--gap-sm);
  right: var(--gap-sm);
  padding: 2px 8px;
  border-radius: var(--radius-xs);
  font-size: 0.625rem;
  font-weight: 600;
  color: #fff;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 信息区域 */
.nft-info {
  padding: var(--gap-md);
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.nft-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nft-collection {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nft-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: var(--gap-xs);
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.meta-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.meta-value {
  font-size: 0.75rem;
  color: var(--color-text-primary);
}

.chain-tag {
  font-size: 0.6875rem;
  color: var(--color-accent-primary);
  font-weight: 500;
}

.nft-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-top: var(--gap-xs);
  padding-top: var(--gap-xs);
  border-top: 1px solid var(--color-border);
}
</style>
