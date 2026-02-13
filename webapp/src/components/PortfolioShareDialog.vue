<script setup>
/**
 * 投资组合导出报告对话框
 * 预览导出图（只显示百分比配置）、多模板选择、生成图片
 * 包含健康评分、不包含绝对金额
 */
import { ref, computed, watch } from 'vue'
import { PhX, PhDownloadSimple, PhCopy } from '@phosphor-icons/vue'
import { useAssetStore } from '../stores/assetStore'
import { useI18n } from '../composables/useI18n'
import { useToast } from '../composables/useToast'

const props = defineProps({
  visible: { type: Boolean, default: false },
})

const emit = defineEmits(['close'])

const assetStore = useAssetStore()
const { t } = useI18n()
const { showToast } = useToast()

const canvasRef = ref(null)
const selectedTemplate = ref('dark')
const generating = ref(false)

// 模板配置
const templates = {
  dark: { bg: '#0D1117', bgSec: '#161B22', text: '#E6EDF3', textSec: '#8B949E', accent: '#4B83F0', border: '#30363D' },
  blue: { bg: '#0F172A', bgSec: '#1E293B', text: '#F1F5F9', textSec: '#94A3B8', accent: '#3B82F6', border: '#334155' },
  green: { bg: '#052E16', bgSec: '#14532D', text: '#F0FDF4', textSec: '#86EFAC', accent: '#22C55E', border: '#166534' },
}

// 资产分类百分比
const categories = computed(() => {
  const cex = assetStore.cexAccounts.reduce((sum, a) => sum + a.balance, 0)
  const chain = assetStore.walletAddresses.reduce((sum, w) => sum + w.balance, 0)
  const manual = assetStore.manualAssets.reduce((sum, a) => sum + (a.balance * (a.currency === 'CNY' ? 0.14 : 1)), 0)
  const total = cex + chain + manual
  if (total === 0) return [{ label: 'CEX', pct: 33 }, { label: 'Chain', pct: 34 }, { label: 'Manual', pct: 33 }]
  return [
    { label: 'CEX', pct: Math.round((cex / total) * 100), color: '#4B83F0' },
    { label: t('dashboard.blockchainAssets'), pct: Math.round((chain / total) * 100), color: '#3EA87A' },
    { label: t('dashboard.manualAssets'), pct: Math.round((manual / total) * 100), color: '#F0A44B' },
  ].filter(c => c.pct > 0)
})

// Top 5 持仓百分比
const topHoldings = computed(() => {
  const holdings = []
  for (const acc of assetStore.cexAccounts) {
    if (acc.holdings) acc.holdings.forEach(h => holdings.push({ symbol: h.symbol, value: h.value }))
  }
  for (const w of assetStore.walletAddresses) {
    if (w.holdings) w.holdings.forEach(h => holdings.push({ symbol: h.symbol, value: h.value }))
  }
  const merged = {}
  holdings.forEach(h => {
    if (!merged[h.symbol]) merged[h.symbol] = { symbol: h.symbol, value: 0 }
    merged[h.symbol].value += h.value
  })
  const sorted = Object.values(merged).sort((a, b) => b.value - a.value)
  const total = sorted.reduce((sum, h) => sum + h.value, 0)
  return sorted.slice(0, 5).map(h => ({
    symbol: h.symbol,
    pct: total > 0 ? ((h.value / total) * 100).toFixed(1) : '0',
  }))
})

// 健康评分（简化版）
const healthScore = computed(() => {
  const concentration = assetStore.assetConcentration
  if (concentration.hhi < 1500) return { score: 85, label: t('social.healthGood') }
  if (concentration.hhi < 2500) return { score: 65, label: t('social.healthFair') }
  return { score: 40, label: t('social.healthPoor') }
})

// 绘制 Canvas 分享图
const drawCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const theme = templates[selectedTemplate.value]
  const ctx = canvas.getContext('2d')
  const w = 1200, h = 630
  canvas.width = w * 2
  canvas.height = h * 2
  ctx.scale(2, 2)

  // 背景
  ctx.fillStyle = theme.bg
  ctx.fillRect(0, 0, w, h)

  // 顶部强调条
  ctx.fillStyle = theme.accent
  ctx.fillRect(0, 0, w, 3)

  // AllFi 品牌
  ctx.fillStyle = theme.textSec
  ctx.font = '500 14px "DM Sans", sans-serif'
  ctx.fillText('AllFi', 48, 44)

  // 标题
  ctx.fillStyle = theme.text
  ctx.font = '700 24px "DM Sans", sans-serif'
  ctx.fillText(t('social.shareTitle'), 48, 80)

  // 日期
  ctx.fillStyle = theme.textSec
  ctx.font = '400 12px "IBM Plex Sans", sans-serif'
  ctx.fillText(new Date().toLocaleDateString('zh-CN'), w - 160, 44)

  // 分割线
  ctx.strokeStyle = theme.border
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(48, 100)
  ctx.lineTo(w - 48, 100)
  ctx.stroke()

  // 左侧：资产配比饼图（CSS conic-gradient 无法用于 Canvas，改用弧形）
  const cx = 200, cy = 260, r = 80
  let startAngle = -Math.PI / 2
  const catColors = ['#4B83F0', '#3EA87A', '#F0A44B']
  categories.value.forEach((cat, i) => {
    const sliceAngle = (cat.pct / 100) * Math.PI * 2
    ctx.beginPath()
    ctx.moveTo(cx, cy)
    ctx.arc(cx, cy, r, startAngle, startAngle + sliceAngle)
    ctx.closePath()
    ctx.fillStyle = catColors[i] || '#888'
    ctx.fill()
    startAngle += sliceAngle
  })
  // 中心圆
  ctx.beginPath()
  ctx.arc(cx, cy, 48, 0, Math.PI * 2)
  ctx.fillStyle = theme.bg
  ctx.fill()

  // 健康评分在中心
  ctx.fillStyle = theme.accent
  ctx.font = '700 28px "IBM Plex Mono", monospace'
  ctx.textAlign = 'center'
  ctx.fillText(String(healthScore.value.score), cx, cy + 4)
  ctx.fillStyle = theme.textSec
  ctx.font = '400 10px "IBM Plex Sans", sans-serif'
  ctx.fillText(healthScore.value.label, cx, cy + 20)
  ctx.textAlign = 'left'

  // 图例
  let legendY = 180
  categories.value.forEach((cat, i) => {
    ctx.fillStyle = catColors[i] || '#888'
    ctx.beginPath()
    ctx.roundRect(48, legendY - 6, 10, 10, 2)
    ctx.fill()
    ctx.fillStyle = theme.text
    ctx.font = '500 13px "IBM Plex Sans", sans-serif'
    ctx.fillText(`${cat.label}  ${cat.pct}%`, 66, legendY + 3)
    legendY += 24
  })

  // 右侧：Top 5 持仓
  const rightX = 360
  ctx.fillStyle = theme.textSec
  ctx.font = '600 12px "IBM Plex Sans", sans-serif'
  ctx.fillText('TOP HOLDINGS', rightX, 140)

  topHoldings.value.forEach((h, i) => {
    const y = 168 + i * 36
    // 进度条背景
    ctx.fillStyle = theme.bgSec
    ctx.beginPath()
    ctx.roundRect(rightX, y, 500, 24, 4)
    ctx.fill()
    // 进度条
    ctx.fillStyle = theme.accent
    ctx.globalAlpha = 0.3
    ctx.beginPath()
    ctx.roundRect(rightX, y, Math.max(500 * (parseFloat(h.pct) / 100), 4), 24, 4)
    ctx.fill()
    ctx.globalAlpha = 1
    // 文字
    ctx.fillStyle = theme.text
    ctx.font = '600 12px "IBM Plex Sans", sans-serif'
    ctx.fillText(h.symbol, rightX + 10, y + 16)
    ctx.fillStyle = theme.textSec
    ctx.font = '500 12px "IBM Plex Mono", monospace'
    ctx.textAlign = 'right'
    ctx.fillText(`${h.pct}%`, rightX + 490, y + 16)
    ctx.textAlign = 'left'
  })

  // 底部分割线
  ctx.strokeStyle = theme.border
  ctx.beginPath()
  ctx.moveTo(48, h - 60)
  ctx.lineTo(w - 48, h - 60)
  ctx.stroke()

  // 底部
  ctx.fillStyle = theme.accent
  ctx.font = '600 13px "DM Sans", sans-serif'
  ctx.fillText('AllFi Portfolio Report', 48, h - 30)
  ctx.fillStyle = theme.textSec
  ctx.font = '400 11px "IBM Plex Sans", sans-serif'
  ctx.textAlign = 'right'
  ctx.fillText('AllFi — Your Portfolio, Your Data', w - 48, h - 30)
  ctx.textAlign = 'left'
}

// 监听模板变化重绘
watch([selectedTemplate, () => props.visible], () => {
  if (props.visible) {
    requestAnimationFrame(drawCanvas)
  }
})

// 下载图片
const downloadImage = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  generating.value = true
  canvas.toBlob((blob) => {
    if (!blob) { generating.value = false; return }
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `allfi-portfolio-${Date.now()}.png`
    link.click()
    URL.revokeObjectURL(url)
    generating.value = false
    showToast(t('social.downloadSuccess'), 'success')
  }, 'image/png')
}

// 复制到剪贴板
const copyToClipboard = async () => {
  const canvas = canvasRef.value
  if (!canvas) return
  try {
    const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/png'))
    if (blob) {
      await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })])
      showToast(t('common.copySuccess'), 'success')
    }
  } catch (e) {
    showToast(t('common.operationFailed'), 'error')
  }
}
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="share-overlay" @click.self="emit('close')">
      <div class="share-dialog">
        <!-- 头部 -->
        <div class="dialog-header">
          <h3>{{ t('social.sharePortfolio') }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="18" />
          </button>
        </div>

        <div class="dialog-body">
          <!-- 模板选择 -->
          <div class="template-selector">
            <span class="selector-label">{{ t('social.selectTemplate') }}</span>
            <div class="template-options">
              <button
                v-for="(_, key) in templates"
                :key="key"
                class="template-btn"
                :class="{ active: selectedTemplate === key }"
                @click="selectedTemplate = key"
              >
                <span class="template-preview" :style="{ background: templates[key].bg, borderColor: templates[key].accent }" />
                {{ t(`social.template_${key}`) }}
              </button>
            </div>
          </div>

          <!-- Canvas 预览 -->
          <div class="canvas-wrapper">
            <canvas ref="canvasRef" class="share-canvas" />
          </div>

          <div class="privacy-hint">{{ t('social.noAmountsHint') }}</div>

          <!-- 操作按钮 -->
          <div class="dialog-actions">
            <button class="btn btn-secondary" @click="copyToClipboard">
              <PhCopy :size="14" />
              {{ t('social.copyImage') }}
            </button>
            <button class="btn btn-primary" :disabled="generating" @click="downloadImage">
              <PhDownloadSimple :size="14" />
              {{ t('social.downloadImage') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.share-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.share-dialog {
  width: 100%;
  max-width: 680px;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-lg) var(--gap-xl);
  border-bottom: 1px solid var(--color-border);
}

.dialog-header h3 {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--gap-xs);
}

.close-btn:hover {
  color: var(--color-text-primary);
}

.dialog-body {
  padding: var(--gap-xl);
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

/* 模板选择 */
.template-selector {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
}

.selector-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.template-options {
  display: flex;
  gap: var(--gap-sm);
}

.template-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 4px 12px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.template-btn:hover {
  border-color: var(--color-accent-primary);
}

.template-btn.active {
  border-color: var(--color-accent-primary);
  color: var(--color-accent-primary);
}

.template-preview {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  border: 1px solid;
}

/* Canvas */
.canvas-wrapper {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.share-canvas {
  width: 100%;
  height: auto;
  display: block;
}

.privacy-hint {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  text-align: center;
}

/* 操作 */
.dialog-actions {
  display: flex;
  justify-content: center;
  gap: var(--gap-sm);
}

.dialog-actions .btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 8px 16px;
  font-size: 0.8125rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: none;
  transition: opacity var(--transition-fast);
}

.btn-primary {
  background: var(--color-accent-primary);
  color: #fff;
}

.btn-primary:hover:not(:disabled) { opacity: 0.9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

.btn-secondary {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover { opacity: 0.9; }

/* 过渡 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
