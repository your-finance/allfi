<script setup>
/**
 * 投资组合健康评分卡片
 * 环形进度条 + 4 个维度进度条 + 一句话建议
 */
import { computed } from 'vue'
import { useHealthScore } from '../composables/useHealthScore'
import { useI18n } from '../composables/useI18n'

const { healthScore } = useHealthScore()
const { t } = useI18n()

// 环形进度条参数
const RADIUS = 52
const CIRCUMFERENCE = 2 * Math.PI * RADIUS
const strokeDashoffset = computed(() => {
  if (!healthScore.value.hasData) return CIRCUMFERENCE
  const progress = healthScore.value.total / 100
  return CIRCUMFERENCE * (1 - progress)
})

// 等级文字
const gradeLabel = computed(() => {
  return t(`healthScore.grade.${healthScore.value.grade}`)
})

// 维度颜色（根据得分比例）
const dimensionColor = (score, maxScore) => {
  const ratio = score / maxScore
  if (ratio >= 0.8) return 'var(--color-success)'
  if (ratio >= 0.5) return 'var(--color-warning)'
  if (ratio >= 0.3) return '#f97316'
  return 'var(--color-error)'
}

// 建议文字
const advice = computed(() => {
  if (!healthScore.value.hasData) return ''
  const w = healthScore.value.weakest
  if (!w) return ''
  return t(`healthScore.advice.${w}`)
})
</script>

<template>
  <div class="health-card">
    <div class="health-header">
      <h3>{{ t('healthScore.title') }}</h3>
      <span class="health-grade" :style="{ color: healthScore.color }">
        {{ gradeLabel }}
      </span>
    </div>

    <div class="health-body">
      <!-- 环形进度条 -->
      <div class="ring-wrapper">
        <svg class="ring-svg" viewBox="0 0 120 120">
          <!-- 背景圆 -->
          <circle
            cx="60" cy="60" :r="RADIUS"
            fill="none"
            stroke="var(--color-bg-tertiary)"
            stroke-width="8"
          />
          <!-- 进度圆 -->
          <circle
            v-if="healthScore.hasData"
            cx="60" cy="60" :r="RADIUS"
            fill="none"
            :stroke="healthScore.color"
            stroke-width="8"
            stroke-linecap="round"
            :stroke-dasharray="CIRCUMFERENCE"
            :stroke-dashoffset="strokeDashoffset"
            class="ring-progress"
          />
        </svg>
        <div class="ring-center">
          <span class="ring-score font-mono" :style="{ color: healthScore.color }">
            {{ healthScore.hasData ? healthScore.total : '--' }}
          </span>
          <span class="ring-label">/ 100</span>
        </div>
      </div>

      <!-- 维度详情 -->
      <div class="dimensions">
        <div
          v-for="dim in healthScore.dimensions"
          :key="dim.id"
          class="dim-row"
        >
          <div class="dim-header">
            <span class="dim-name">{{ t(`healthScore.dim.${dim.id}`) }}</span>
            <span class="dim-score font-mono">
              {{ healthScore.hasData ? dim.score : '--' }}/{{ dim.maxScore }}
            </span>
          </div>
          <div class="dim-bar-bg">
            <div
              class="dim-bar-fill"
              :style="{
                width: healthScore.hasData ? `${(dim.score / dim.maxScore) * 100}%` : '0%',
                background: healthScore.hasData ? dimensionColor(dim.score, dim.maxScore) : 'var(--color-bg-tertiary)',
              }"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 建议 -->
    <div v-if="advice" class="health-advice">
      {{ advice }}
    </div>
    <div v-else-if="!healthScore.hasData" class="health-advice health-no-data">
      {{ t('healthScore.noData') }}
    </div>
  </div>
</template>

<style scoped>
.health-card {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.health-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-md);
}

.health-header h3 {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.health-grade {
  font-size: 0.75rem;
  font-weight: 600;
}

.health-body {
  display: flex;
  gap: var(--gap-xl);
  align-items: center;
}

/* 环形进度条 */
.ring-wrapper {
  position: relative;
  width: 120px;
  height: 120px;
  flex-shrink: 0;
}

.ring-svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.ring-progress {
  transition: stroke-dashoffset 0.6s ease;
}

.ring-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.ring-score {
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1;
}

.ring-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-top: 2px;
}

/* 维度 */
.dimensions {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.dim-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.dim-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dim-name {
  font-size: 0.6875rem;
  color: var(--color-text-secondary);
}

.dim-score {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.dim-bar-bg {
  height: 4px;
  background: var(--color-bg-tertiary);
  border-radius: 2px;
  overflow: hidden;
}

.dim-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.6s ease, background 0.3s ease;
}

/* 建议 */
.health-advice {
  margin-top: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.6875rem;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  line-height: 1.5;
}

.health-no-data {
  color: var(--color-text-muted);
  text-align: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .health-body {
    flex-direction: column;
  }

  .ring-wrapper {
    width: 100px;
    height: 100px;
  }

  .ring-score {
    font-size: 1.25rem;
  }
}
</style>
