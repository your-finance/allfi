<script setup>
/**
 * 目标卡片组件
 * 显示单个目标的进度条、当前值 vs 目标值、预估达成日期
 */
import { PhTrash } from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  goal: {
    type: Object,
    required: true,
  }
})

const emit = defineEmits(['delete'])

const { formatNumber, currencySymbol } = useFormatters()
const { t } = useI18n()

// 进度条颜色
const progressColor = (pct) => {
  if (pct >= 80) return 'var(--color-success)'
  if (pct >= 50) return 'var(--color-accent-primary)'
  if (pct >= 25) return 'var(--color-warning)'
  return 'var(--color-error)'
}

// 目标类型标签
const typeLabel = (type) => {
  return t(`goals.type.${type}`)
}

// 格式化目标值
const formatGoalValue = (goal) => {
  if (goal.type === 'asset_value') {
    return `${currencySymbol.value}${formatNumber(goal.targetValue)}`
  }
  if (goal.type === 'holding_amount') {
    return `${formatNumber(goal.targetValue, 4)} ${goal.currency}`
  }
  if (goal.type === 'return_rate') {
    return `${goal.targetValue}%`
  }
  return String(goal.targetValue)
}

// 格式化当前值
const formatCurrentValue = (goal) => {
  if (goal.type === 'asset_value') {
    return `${currencySymbol.value}${formatNumber(goal.currentValue)}`
  }
  if (goal.type === 'holding_amount') {
    return `${formatNumber(goal.currentValue, 4)} ${goal.currency}`
  }
  if (goal.type === 'return_rate') {
    return `${goal.currentValue.toFixed(1)}%`
  }
  return String(goal.currentValue)
}
</script>

<template>
  <div class="goal-card">
    <div class="goal-header">
      <div class="goal-title-row">
        <span class="goal-title">{{ goal.title }}</span>
        <span class="goal-type">{{ typeLabel(goal.type) }}</span>
      </div>
      <button class="goal-delete-btn" @click="emit('delete', goal.id)" :title="t('common.delete')">
        <PhTrash :size="14" />
      </button>
    </div>

    <!-- 进度条 -->
    <div class="goal-progress-wrap">
      <div class="goal-progress-bg">
        <div
          class="goal-progress-fill"
          :style="{
            width: `${goal.progress}%`,
            background: progressColor(goal.progress),
          }"
        />
      </div>
      <span class="goal-progress-pct font-mono" :style="{ color: progressColor(goal.progress) }">
        {{ goal.progress }}%
      </span>
    </div>

    <!-- 当前值 vs 目标值 -->
    <div class="goal-values">
      <div class="goal-current">
        <span class="goal-val-label">{{ t('goals.current') }}</span>
        <span class="goal-val font-mono">{{ formatCurrentValue(goal) }}</span>
      </div>
      <div class="goal-target">
        <span class="goal-val-label">{{ t('goals.target') }}</span>
        <span class="goal-val font-mono">{{ formatGoalValue(goal) }}</span>
      </div>
    </div>

    <!-- 预估达成日期 -->
    <div v-if="goal.estimatedDate" class="goal-estimate">
      {{ t('goals.estimatedDate') }}: {{ goal.estimatedDate }}
    </div>
    <div v-else-if="goal.progress >= 100" class="goal-estimate goal-achieved">
      {{ t('goals.achieved') }}
    </div>
  </div>
</template>

<style scoped>
.goal-card {
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.goal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-sm);
}

.goal-title-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.goal-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.goal-type {
  font-size: 0.625rem;
  color: var(--color-text-muted);
  padding: 1px 4px;
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xs);
}

.goal-delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: background var(--transition-fast), color var(--transition-fast);
}

.goal-delete-btn:hover {
  background: var(--color-bg-elevated);
  color: var(--color-error);
}

/* 进度条 */
.goal-progress-wrap {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-sm);
}

.goal-progress-bg {
  flex: 1;
  height: 6px;
  background: var(--color-bg-elevated);
  border-radius: 3px;
  overflow: hidden;
}

.goal-progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s ease;
}

.goal-progress-pct {
  font-size: 0.75rem;
  font-weight: 600;
  min-width: 40px;
  text-align: right;
}

/* 当前值 vs 目标值 */
.goal-values {
  display: flex;
  justify-content: space-between;
  gap: var(--gap-md);
}

.goal-current,
.goal-target {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.goal-val-label {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.goal-val {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

/* 预估达成 */
.goal-estimate {
  margin-top: var(--gap-sm);
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.goal-achieved {
  color: var(--color-success);
  font-weight: 600;
}
</style>
