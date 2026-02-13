<script setup>
/**
 * 年报分享图生成组件
 * 基于 Canvas 生成精美分享图片
 * 包含年度收益率、资产配置、风格标签
 * 不包含绝对金额（隐私保护）
 */
import { ref, onMounted, computed } from 'vue'
import { PhDownloadSimple, PhX } from '@phosphor-icons/vue'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  report: { type: Object, required: true },
})

const emit = defineEmits(['close'])

const themeStore = useThemeStore()
const { t } = useI18n()

const canvasRef = ref(null)
const colors = computed(() => themeStore.currentTheme.colors)

// 投资风格标签
const styleLabels = {
  steady: 'annualReport.styleSteady',
  aggressive: 'annualReport.styleAggressive',
  conservative: 'annualReport.styleConservative',
  balanced: 'annualReport.styleBalanced',
}

// 绘制分享图
const drawShareImage = () => {
  const canvas = canvasRef.value
  if (!canvas || !props.report) return

  const ctx = canvas.getContext('2d')
  const w = 1200
  const h = 630
  canvas.width = w * 2
  canvas.height = h * 2
  ctx.scale(2, 2)

  // 背景
  ctx.fillStyle = colors.value.bgPrimary || '#0D1117'
  ctx.fillRect(0, 0, w, h)

  // 顶部渐变条
  const grad = ctx.createLinearGradient(0, 0, w, 0)
  grad.addColorStop(0, colors.value.accentPrimary || '#4B83F0')
  grad.addColorStop(1, colors.value.accentSecondary || '#3EA87A')
  ctx.fillStyle = grad
  ctx.fillRect(0, 0, w, 4)

  // AllFi 标题
  ctx.fillStyle = colors.value.textMuted || '#8B949E'
  ctx.font = '500 14px "DM Sans", sans-serif'
  ctx.fillText('AllFi', 48, 48)

  // 年份标题
  ctx.fillStyle = colors.value.textPrimary || '#E6EDF3'
  ctx.font = '700 28px "DM Sans", sans-serif'
  ctx.fillText(`${props.report.year} ${t('annualReport.title')}`, 48, 90)

  // 年度收益率（大字）
  const returnText = `${props.report.summary.totalReturn >= 0 ? '+' : ''}${props.report.summary.totalReturn}%`
  const returnColor = props.report.summary.totalReturn >= 0
    ? (colors.value.accentPrimary || '#10B981')
    : '#EF4444'
  ctx.fillStyle = returnColor
  ctx.font = '700 72px "IBM Plex Mono", monospace'
  ctx.fillText(returnText, 48, 190)

  // 收益标签
  ctx.fillStyle = colors.value.textMuted || '#8B949E'
  ctx.font = '400 16px "IBM Plex Sans", sans-serif'
  ctx.fillText(t('annualReport.yearlyReturn'), 48, 220)

  // 基准对比
  const benchmarks = props.report.benchmarks
  const benchmarkY = 270
  ctx.font = '500 14px "IBM Plex Sans", sans-serif'
  ctx.fillStyle = colors.value.textSecondary || '#8B949E'
  ctx.fillText(`vs BTC +${benchmarks.btc}%  |  vs ETH +${benchmarks.eth}%  |  vs S&P 500 +${benchmarks.sp500}%`, 48, benchmarkY)

  // 分割线
  ctx.strokeStyle = colors.value.border || '#30363D'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(48, 300)
  ctx.lineTo(w - 48, 300)
  ctx.stroke()

  // 月度收益条形图（简化版）
  const months = props.report.monthlyReturns
  const barWidth = (w - 136) / 12
  const barBaseY = 440
  const maxReturn = Math.max(...months.map(m => Math.abs(m.return)))
  const barScale = 100 / maxReturn

  ctx.font = '400 10px "IBM Plex Sans", sans-serif'
  const monthLabels = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12']

  months.forEach((m, i) => {
    const x = 48 + i * barWidth + barWidth * 0.15
    const bw = barWidth * 0.7
    const bh = Math.abs(m.return) * barScale
    const y = m.return >= 0 ? barBaseY - bh : barBaseY

    ctx.fillStyle = m.return >= 0
      ? (colors.value.accentPrimary || '#10B981')
      : '#EF4444'
    ctx.globalAlpha = 0.8
    ctx.beginPath()
    ctx.roundRect(x, y, bw, bh, 2)
    ctx.fill()
    ctx.globalAlpha = 1

    // 月份标签
    ctx.fillStyle = colors.value.textMuted || '#8B949E'
    ctx.textAlign = 'center'
    ctx.fillText(monthLabels[i], x + bw / 2, barBaseY + 16)
  })
  ctx.textAlign = 'left'

  // 投资风格标签
  ctx.fillStyle = colors.value.accentPrimary || '#4B83F0'
  ctx.font = '700 20px "DM Sans", sans-serif'
  ctx.fillText(t(styleLabels[props.report.styleTag]), 48, 510)

  // 年度之星
  ctx.fillStyle = colors.value.textSecondary || '#8B949E'
  ctx.font = '400 13px "IBM Plex Sans", sans-serif'
  ctx.fillText(
    `${t('annualReport.annualStar')}: ${props.report.annualStar.symbol} +${props.report.annualStar.return}%`,
    48, 540
  )

  // 底部水印
  ctx.fillStyle = colors.value.textMuted || '#484F58'
  ctx.font = '400 12px "IBM Plex Sans", sans-serif'
  ctx.fillText('Powered by AllFi — #MyAllFiYear', 48, h - 24)

  // 右下角年份
  ctx.textAlign = 'right'
  ctx.fillStyle = colors.value.textMuted || '#484F58'
  ctx.font = '600 48px "DM Sans", sans-serif'
  ctx.globalAlpha = 0.1
  ctx.fillText(String(props.report.year), w - 48, h - 20)
  ctx.globalAlpha = 1
  ctx.textAlign = 'left'
}

// 下载图片
const downloadImage = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  canvas.toBlob((blob) => {
    if (!blob) return
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `allfi-${props.report.year}-report.png`
    link.click()
    URL.revokeObjectURL(url)
  }, 'image/png')
}

onMounted(() => {
  if (props.report) drawShareImage()
})
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="share-overlay" @click.self="emit('close')">
      <div class="share-panel">
        <div class="share-header">
          <h3>{{ t('annualReport.shareTitle') }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="18" />
          </button>
        </div>

        <div class="canvas-wrapper">
          <canvas ref="canvasRef" class="share-canvas" />
        </div>

        <div class="share-hint">{{ t('annualReport.shareHint') }}</div>

        <div class="share-actions">
          <button class="btn btn-primary" @click="downloadImage">
            <PhDownloadSimple :size="14" />
            {{ t('annualReport.downloadImage') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.share-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  padding: var(--gap-lg);
}

.share-panel {
  width: 100%;
  max-width: 640px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--gap-xl);
}

.share-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-lg);
}

.share-header h3 {
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

.canvas-wrapper {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  margin-bottom: var(--gap-md);
}

.share-canvas {
  width: 100%;
  height: auto;
  display: block;
}

.share-hint {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  text-align: center;
  margin-bottom: var(--gap-md);
}

.share-actions {
  display: flex;
  justify-content: center;
}

.share-actions .btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

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
