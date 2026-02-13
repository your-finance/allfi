<script setup>
/**
 * 资产配置分享卡片
 * 匿名化展示：只显示百分比，不显示绝对金额
 * 包含：资产分布饼图 + 配置比例 + AllFi 标识
 */
import { ref, computed, nextTick } from 'vue'
import { PhX, PhDownloadSimple, PhShareNetwork } from '@phosphor-icons/vue'
import { useAssetStore } from '../stores/assetStore'
import { useI18n } from '../composables/useI18n'
import { useToast } from '../composables/useToast'
import { generateShareImage } from '../composables/useShareImage'

const { t } = useI18n()
const { showToast } = useToast()
const assetStore = useAssetStore()

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close'])

const cardRef = ref(null)
const generating = ref(false)

// 资产分类数据（匿名化，只保留百分比）
const categories = computed(() => {
  const cex = assetStore.cexAccounts.reduce((sum, a) => sum + a.balance, 0)
  const chain = assetStore.walletAddresses.reduce((sum, w) => sum + w.balance, 0)
  const manual = assetStore.manualAssets.reduce((sum, a) => sum + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0)
  const total = cex + chain + manual

  if (total === 0) {
    return [
      { label: t('dashboard.cexAssets'), percent: 33.3, color: '#4B83F0' },
      { label: t('dashboard.blockchainAssets'), percent: 33.3, color: '#3EA87A' },
      { label: t('dashboard.manualAssets'), percent: 33.4, color: '#F0A44B' }
    ]
  }

  return [
    { label: t('dashboard.cexAssets'), percent: ((cex / total) * 100), color: '#4B83F0' },
    { label: t('dashboard.blockchainAssets'), percent: ((chain / total) * 100), color: '#3EA87A' },
    { label: t('dashboard.manualAssets'), percent: ((manual / total) * 100), color: '#F0A44B' }
  ].filter(c => c.percent > 0)
})

// 持仓 Top 5（只显示百分比）
const topHoldings = computed(() => {
  const holdings = []
  for (const acc of assetStore.cexAccounts) {
    if (acc.holdings) {
      for (const h of acc.holdings) {
        holdings.push({ symbol: h.symbol, value: h.value })
      }
    }
  }
  for (const w of assetStore.walletAddresses) {
    if (w.holdings) {
      for (const h of w.holdings) {
        holdings.push({ symbol: h.symbol, value: h.value })
      }
    }
  }

  // 合并相同 symbol
  const merged = {}
  for (const h of holdings) {
    if (!merged[h.symbol]) {
      merged[h.symbol] = { symbol: h.symbol, value: 0 }
    }
    merged[h.symbol].value += h.value
  }

  const sorted = Object.values(merged).sort((a, b) => b.value - a.value)
  const total = sorted.reduce((sum, h) => sum + h.value, 0)

  return sorted.slice(0, 5).map(h => ({
    symbol: h.symbol,
    percent: total > 0 ? ((h.value / total) * 100).toFixed(1) : '0.0'
  }))
})

// 饼图 CSS 渐变（纯 CSS 实现，避免 canvas 嵌套问题）
const pieGradient = computed(() => {
  let current = 0
  const stops = []
  for (const cat of categories.value) {
    stops.push(`${cat.color} ${current}%`)
    current += cat.percent
    stops.push(`${cat.color} ${current}%`)
  }
  return `conic-gradient(${stops.join(', ')})`
})

// 当前日期
const currentDate = computed(() => {
  const d = new Date()
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
})

// 资产总数
const totalAccounts = computed(() => {
  return assetStore.cexAccounts.length + assetStore.walletAddresses.length + assetStore.manualAssets.length
})

// 下载分享图
const downloadImage = async () => {
  if (!cardRef.value) return
  generating.value = true
  await nextTick()
  try {
    await generateShareImage(cardRef.value, `allfi-portfolio-${Date.now()}`)
    showToast(t('share.download') + ' ✓', 'success')
  } catch (e) {
    showToast(t('common.operationFailed'), 'error')
    console.error('生成分享图失败:', e)
  } finally {
    generating.value = false
  }
}
</script>

<template>
  <Transition name="fade">
    <div v-if="visible" class="dialog-overlay" @click.self="emit('close')">
      <div class="share-dialog">
        <!-- 标题栏 -->
        <div class="dialog-header">
          <h3>{{ t('share.title') }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="16" />
          </button>
        </div>

        <div class="dialog-body">
          <!-- 分享卡片预览 -->
          <div ref="cardRef" class="share-card">
            <!-- 卡片头部 -->
            <div class="card-header">
              <div class="brand">
                <svg viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg" class="brand-icon">
                  <rect x="1" y="1" width="26" height="26" rx="5" stroke="#4B83F0" stroke-width="1.5" fill="none"/>
                  <path d="M8 14 L14 8 L20 14 L14 20 Z" fill="#4B83F0" opacity="0.8"/>
                </svg>
                <span class="brand-name">AllFi</span>
              </div>
              <span class="card-date">{{ currentDate }}</span>
            </div>

            <!-- 资产配置饼图 + 图例 -->
            <div class="card-main">
              <div class="pie-section">
                <div class="pie-chart" :style="{ background: pieGradient }">
                  <div class="pie-center">
                    <span class="pie-count">{{ totalAccounts }}</span>
                    <span class="pie-label">{{ t('share.accounts') }}</span>
                  </div>
                </div>
              </div>

              <div class="legend-section">
                <div class="legend-title">{{ t('share.allocation') }}</div>
                <div v-for="cat in categories" :key="cat.label" class="legend-item">
                  <div class="legend-dot" :style="{ background: cat.color }" />
                  <span class="legend-name">{{ cat.label }}</span>
                  <span class="legend-percent">{{ cat.percent.toFixed(1) }}%</span>
                </div>

                <!-- Top 持仓 -->
                <div v-if="topHoldings.length > 0" class="top-holdings">
                  <div class="holdings-title">{{ t('share.topHoldings') }}</div>
                  <div v-for="h in topHoldings" :key="h.symbol" class="holding-item">
                    <span class="holding-symbol">{{ h.symbol }}</span>
                    <span class="holding-percent">{{ h.percent }}%</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 卡片底部 -->
            <div class="card-footer">
              <span class="footer-tag">#AllFi</span>
              <span class="footer-text">{{ t('share.footerText') }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="dialog-actions">
            <button class="btn btn-primary" :disabled="generating" @click="downloadImage">
              <PhDownloadSimple :size="14" />
              {{ generating ? t('share.generating') : t('share.download') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.share-dialog {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  width: 680px;
  max-width: 95vw;
  max-height: 90vh;
  overflow-y: auto;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.dialog-header h3 {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
}

.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dialog-body {
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

/* ========== 分享卡片 ========== */
.share-card {
  width: 600px;
  height: 315px;
  margin: 0 auto;
  padding: 28px 32px;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  color: #e2e8f0;
  font-family: 'IBM Plex Sans', -apple-system, sans-serif;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-icon {
  width: 24px;
  height: 24px;
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
  color: #f1f5f9;
  letter-spacing: 0.5px;
}

.card-date {
  font-size: 12px;
  color: #94a3b8;
}

/* 主体 */
.card-main {
  flex: 1;
  display: flex;
  gap: 32px;
  align-items: center;
}

.pie-section {
  flex-shrink: 0;
}

.pie-chart {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  position: relative;
}

.pie-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: #1e293b;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.pie-count {
  font-size: 24px;
  font-weight: 700;
  color: #f1f5f9;
  font-family: 'IBM Plex Mono', monospace;
}

.pie-label {
  font-size: 10px;
  color: #94a3b8;
  margin-top: 2px;
}

/* 图例 */
.legend-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.legend-title {
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  flex-shrink: 0;
}

.legend-name {
  flex: 1;
  font-size: 13px;
  color: #cbd5e1;
}

.legend-percent {
  font-size: 14px;
  font-weight: 600;
  color: #f1f5f9;
  font-family: 'IBM Plex Mono', monospace;
}

/* Top 持仓 */
.top-holdings {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(148, 163, 184, 0.15);
}

.holdings-title {
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}

.holding-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1px 0;
}

.holding-symbol {
  font-size: 12px;
  color: #cbd5e1;
  font-weight: 500;
}

.holding-percent {
  font-size: 12px;
  color: #94a3b8;
  font-family: 'IBM Plex Mono', monospace;
}

/* 底部 */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid rgba(148, 163, 184, 0.1);
}

.footer-tag {
  font-size: 13px;
  font-weight: 600;
  color: #4B83F0;
}

.footer-text {
  font-size: 11px;
  color: #64748b;
}

/* 操作 */
.dialog-actions {
  display: flex;
  justify-content: center;
  gap: var(--gap-sm);
}

.btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 20px;
  font-size: 0.8125rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-primary {
  background: var(--color-accent-primary);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 过渡 */
.fade-enter-active, .fade-leave-active {
  transition: opacity var(--transition-fast);
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* 响应式 */
@media (max-width: 640px) {
  .share-card {
    width: 100%;
    height: auto;
    aspect-ratio: 1200 / 630;
    padding: 20px;
  }

  .pie-chart {
    width: 120px;
    height: 120px;
  }

  .pie-center {
    width: 72px;
    height: 72px;
  }

  .pie-count {
    font-size: 18px;
  }
}
</style>
