<script setup>
/**
 * 资产配置分享卡片
 * 支持：数值/比例切换、多种精美样式、完整图片导出
 */
import { ref, computed, nextTick, watch } from 'vue'
import { PhX, PhDownloadSimple, PhPercent, PhCurrencyDollar } from '@phosphor-icons/vue'
import { useAssetStore } from '../stores/assetStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import { useToast } from '../composables/useToast'
import html2canvas from 'html2canvas'

const { t } = useI18n()
const { showToast } = useToast()
const { formatNumber } = useFormatters()
const assetStore = useAssetStore()

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close'])

const cardRef = ref(null)
const generating = ref(false)

// 显示模式：percent（百分比）或 value（数值）
const displayMode = ref('percent')

// 卡片样式
const cardStyles = [
  { id: 'midnight', name: '午夜深蓝', gradient: 'linear-gradient(135deg, #0c1929 0%, #1a365d 50%, #0f172a 100%)', accent: '#60a5fa' },
  { id: 'aurora', name: '极光幻彩', gradient: 'linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%)', accent: '#a78bfa' },
  { id: 'ember', name: '琥珀暖阳', gradient: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f0e17 100%)', accent: '#fbbf24' },
  { id: 'forest', name: '翡翠森林', gradient: 'linear-gradient(135deg, #0d1f22 0%, #1a3a3a 50%, #0f1f1f 100%)', accent: '#34d399' },
  { id: 'rose', name: '玫瑰金辉', gradient: 'linear-gradient(135deg, #1f1a24 0%, #2d2438 50%, #1a1520 100%)', accent: '#f472b6' }
]
const selectedStyle = ref('midnight')

const currentStyle = computed(() => cardStyles.find(s => s.id === selectedStyle.value) || cardStyles[0])

// 资产总值
const totalValue = computed(() => {
  const cex = assetStore.cexAccounts.reduce((sum, a) => sum + a.balance, 0)
  const chain = assetStore.walletAddresses.reduce((sum, w) => sum + w.balance, 0)
  const manual = assetStore.manualAssets.reduce((sum, a) => sum + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0)
  return cex + chain + manual
})

// 资产分类数据
const categories = computed(() => {
  const cex = assetStore.cexAccounts.reduce((sum, a) => sum + a.balance, 0)
  const chain = assetStore.walletAddresses.reduce((sum, w) => sum + w.balance, 0)
  const manual = assetStore.manualAssets.reduce((sum, a) => sum + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0)
  const total = cex + chain + manual

  if (total === 0) {
    return [
      { label: t('dashboard.cexAssets'), value: 0, percent: 33.3, color: '#60a5fa' },
      { label: t('dashboard.blockchainAssets'), value: 0, percent: 33.3, color: '#34d399' },
      { label: t('dashboard.manualAssets'), value: 0, percent: 33.4, color: '#fbbf24' }
    ]
  }

  return [
    { label: t('dashboard.cexAssets'), value: cex, percent: ((cex / total) * 100), color: '#60a5fa' },
    { label: t('dashboard.blockchainAssets'), value: chain, percent: ((chain / total) * 100), color: '#34d399' },
    { label: t('dashboard.manualAssets'), value: manual, percent: ((manual / total) * 100), color: '#fbbf24' }
  ].filter(c => c.percent > 0)
})

// 持仓 Top 5
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
    value: h.value,
    percent: total > 0 ? ((h.value / total) * 100).toFixed(1) : '0.0'
  }))
})

// 饼图 CSS 渐变
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

// 格式化显示值
const formatDisplayValue = (value, percent) => {
  if (displayMode.value === 'percent') {
    return `${percent.toFixed(1)}%`
  }
  return `$${formatNumber(value, 0)}`
}

// 下载分享图 - 修复版本
const downloadImage = async () => {
  if (!cardRef.value) return
  generating.value = true

  await nextTick()

  try {
    // 获取卡片实际尺寸
    const rect = cardRef.value.getBoundingClientRect()

    const canvas = await html2canvas(cardRef.value, {
      width: rect.width,
      height: rect.height,
      scale: 2,
      backgroundColor: null,
      useCORS: true,
      logging: false,
      allowTaint: true,
      foreignObjectRendering: false,
      removeContainer: true,
      imageTimeout: 0,
      onclone: (clonedDoc, element) => {
        // 确保克隆的元素有正确的尺寸
        element.style.width = `${rect.width}px`
        element.style.height = `${rect.height}px`
        element.style.transform = 'none'
        element.style.position = 'relative'
      }
    })

    // 转换为 blob 并下载
    canvas.toBlob((blob) => {
      if (!blob) {
        showToast(t('common.operationFailed'), 'error')
        return
      }
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `allfi-portfolio-${Date.now()}.png`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)
      showToast(t('share.download') + ' ✓', 'success')
    }, 'image/png', 1.0)
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
          <!-- 配置选项 -->
          <div class="config-section">
            <!-- 数值/比例切换 -->
            <div class="config-group">
              <span class="config-label">{{ t('share.displayMode') || '显示方式' }}</span>
              <div class="toggle-group">
                <button
                  class="toggle-btn"
                  :class="{ active: displayMode === 'percent' }"
                  @click="displayMode = 'percent'"
                >
                  <PhPercent :size="14" />
                  <span>{{ t('share.percent') || '比例' }}</span>
                </button>
                <button
                  class="toggle-btn"
                  :class="{ active: displayMode === 'value' }"
                  @click="displayMode = 'value'"
                >
                  <PhCurrencyDollar :size="14" />
                  <span>{{ t('share.value') || '数值' }}</span>
                </button>
              </div>
            </div>

            <!-- 样式选择 -->
            <div class="config-group">
              <span class="config-label">{{ t('share.cardStyle') || '卡片样式' }}</span>
              <div class="style-picker">
                <button
                  v-for="style in cardStyles"
                  :key="style.id"
                  class="style-option"
                  :class="{ active: selectedStyle === style.id }"
                  :style="{ background: style.gradient }"
                  :title="style.name"
                  @click="selectedStyle = style.id"
                >
                  <span class="style-dot" :style="{ background: style.accent }"></span>
                </button>
              </div>
            </div>
          </div>

          <!-- 分享卡片预览 -->
          <div
            ref="cardRef"
            class="share-card"
            :style="{ background: currentStyle.gradient }"
          >
            <!-- 装饰元素 -->
            <div class="card-decoration">
              <div class="deco-circle deco-1" :style="{ background: currentStyle.accent }"></div>
              <div class="deco-circle deco-2" :style="{ background: currentStyle.accent }"></div>
              <div class="deco-line" :style="{ background: currentStyle.accent }"></div>
            </div>

            <!-- 卡片头部 -->
            <div class="card-header">
              <div class="brand">
                <div class="brand-icon" :style="{ borderColor: currentStyle.accent }">
                  <svg viewBox="0 0 24 24" fill="none">
                    <path d="M6 12 L12 6 L18 12 L12 18 Z" :fill="currentStyle.accent" opacity="0.9"/>
                  </svg>
                </div>
                <span class="brand-name">AllFi</span>
              </div>
              <div class="header-right">
                <span class="card-date">{{ currentDate }}</span>
              </div>
            </div>

            <!-- 总资产（数值模式显示） -->
            <div v-if="displayMode === 'value'" class="total-section">
              <span class="total-label">{{ t('dashboard.totalAssets') }}</span>
              <span class="total-value" :style="{ color: currentStyle.accent }">${{ formatNumber(totalValue, 0) }}</span>
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
                  <span class="legend-percent">{{ formatDisplayValue(cat.value, cat.percent) }}</span>
                </div>

                <!-- Top 持仓 -->
                <div v-if="topHoldings.length > 0" class="top-holdings">
                  <div class="holdings-title">{{ t('share.topHoldings') }}</div>
                  <div v-for="h in topHoldings" :key="h.symbol" class="holding-item">
                    <span class="holding-symbol">{{ h.symbol }}</span>
                    <span class="holding-percent">
                      {{ displayMode === 'percent' ? `${h.percent}%` : `$${formatNumber(h.value, 0)}` }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 卡片底部 -->
            <div class="card-footer">
              <span class="footer-tag" :style="{ color: currentStyle.accent }">#AllFi</span>
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
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.share-dialog {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  width: 720px;
  max-width: 95vw;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.dialog-header h3 {
  font-size: 1rem;
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
  transition: all 0.15s ease;
}

.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dialog-body {
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* ========== 配置选项 ========== */
.config-section {
  display: flex;
  gap: var(--gap-xl);
  flex-wrap: wrap;
}

.config-group {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.config-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.toggle-group {
  display: flex;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  padding: 3px;
  border: 1px solid var(--color-border);
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.15s ease;
}

.toggle-btn:hover {
  color: var(--color-text-secondary);
}

.toggle-btn.active {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 样式选择器 */
.style-picker {
  display: flex;
  gap: 8px;
}

.style-option {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.style-option::before {
  content: '';
  position: absolute;
  inset: 0;
  background: inherit;
  filter: brightness(0.8);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.style-option:hover::before {
  opacity: 1;
}

.style-option.active {
  border-color: var(--color-accent-primary);
  transform: scale(1.05);
}

.style-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  position: relative;
  z-index: 1;
}

/* ========== 分享卡片 ========== */
.share-card {
  width: 640px;
  height: 360px;
  margin: 0 auto;
  padding: 32px 36px;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  color: #e2e8f0;
  font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, sans-serif;
  position: relative;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

/* 装饰元素 */
.card-decoration {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.deco-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.06;
  filter: blur(60px);
}

.deco-1 {
  width: 300px;
  height: 300px;
  top: -100px;
  right: -50px;
}

.deco-2 {
  width: 200px;
  height: 200px;
  bottom: -80px;
  left: -40px;
}

.deco-line {
  position: absolute;
  width: 1px;
  height: 100%;
  right: 180px;
  top: 0;
  opacity: 0.08;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  width: 32px;
  height: 32px;
  border: 1.5px solid;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.03);
}

.brand-icon svg {
  width: 18px;
  height: 18px;
}

.brand-name {
  font-size: 18px;
  font-weight: 700;
  color: #f8fafc;
  letter-spacing: 0.5px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-date {
  font-size: 12px;
  color: #94a3b8;
  font-weight: 500;
}

/* 总资产 */
.total-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 12px;
  position: relative;
  z-index: 1;
}

.total-label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 500;
}

.total-value {
  font-size: 28px;
  font-weight: 700;
  font-family: 'SF Mono', 'Fira Code', monospace;
  letter-spacing: -0.02em;
}

/* 主体 */
.card-main {
  flex: 1;
  display: flex;
  gap: 36px;
  align-items: center;
  position: relative;
  z-index: 1;
}

.pie-section {
  flex-shrink: 0;
}

.pie-chart {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  position: relative;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.pie-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 84px;
  height: 84px;
  border-radius: 50%;
  background: rgba(15, 23, 42, 0.95);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);
}

.pie-count {
  font-size: 26px;
  font-weight: 700;
  color: #f1f5f9;
  font-family: 'SF Mono', 'Fira Code', monospace;
  line-height: 1;
}

.pie-label {
  font-size: 9px;
  color: #64748b;
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 500;
}

/* 图例 */
.legend-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.legend-title {
  font-size: 10px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 4px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  flex-shrink: 0;
}

.legend-name {
  flex: 1;
  font-size: 13px;
  color: #cbd5e1;
  font-weight: 500;
}

.legend-percent {
  font-size: 14px;
  font-weight: 600;
  color: #f1f5f9;
  font-family: 'SF Mono', 'Fira Code', monospace;
}

/* Top 持仓 */
.top-holdings {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(148, 163, 184, 0.12);
}

.holdings-title {
  font-size: 10px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 6px;
}

.holding-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 0;
}

.holding-symbol {
  font-size: 12px;
  color: #94a3b8;
  font-weight: 600;
}

.holding-percent {
  font-size: 12px;
  color: #64748b;
  font-family: 'SF Mono', 'Fira Code', monospace;
}

/* 底部 */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid rgba(148, 163, 184, 0.08);
  position: relative;
  z-index: 1;
}

.footer-tag {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.footer-text {
  font-size: 11px;
  color: #475569;
  font-weight: 500;
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
  gap: 6px;
  padding: 10px 24px;
  font-size: 0.8125rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary {
  background: var(--color-accent-primary);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(var(--color-accent-primary-rgb), 0.3);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

/* 过渡 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .share-dialog {
    width: 95vw;
  }

  .config-section {
    flex-direction: column;
    gap: var(--gap-md);
  }

  .share-card {
    width: 100%;
    height: auto;
    aspect-ratio: 16 / 9;
    padding: 24px;
  }

  .card-main {
    gap: 24px;
  }

  .pie-chart {
    width: 100px;
    height: 100px;
  }

  .pie-center {
    width: 60px;
    height: 60px;
  }

  .pie-count {
    font-size: 18px;
  }

  .total-value {
    font-size: 22px;
  }
}
</style>
