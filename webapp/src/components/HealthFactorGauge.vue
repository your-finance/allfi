<script setup>
/**
 * 健康因子仪表盘组件
 * 可视化展示健康因子状态
 */
import { computed } from 'vue'
import { PhWarning, PhCheckCircle, PhXCircle } from '@phosphor-icons/vue'

const props = defineProps({
  healthFactor: {
    type: Number,
    required: true
  },
  size: {
    type: String,
    default: 'medium' // small, medium, large
  }
})

// 尺寸配置
const sizeConfig = computed(() => {
  const configs = {
    small: { radius: 40, strokeWidth: 6, fontSize: '1rem' },
    medium: { radius: 60, strokeWidth: 8, fontSize: '1.5rem' },
    large: { radius: 80, strokeWidth: 10, fontSize: '2rem' }
  }
  return configs[props.size] || configs.medium
})

// 健康因子状态
const status = computed(() => {
  const hf = props.healthFactor
  if (hf === 0) return { level: 'none', color: '#6B7280', label: '无借款', icon: PhCheckCircle }
  if (hf >= 2.0) return { level: 'safe', color: '#10B981', label: '安全', icon: PhCheckCircle }
  if (hf >= 1.5) return { level: 'medium', color: '#F59E0B', label: '中等', icon: PhWarning }
  if (hf >= 1.2) return { level: 'high', color: '#EF4444', label: '高风险', icon: PhWarning }
  return { level: 'critical', color: '#DC2626', label: '危险', icon: PhXCircle }
})

// 圆环进度
const progress = computed(() => {
  if (props.healthFactor === 0) return 0
  // 健康因子 1.0 = 0%, 3.0 = 100%
  const normalized = Math.min(Math.max((props.healthFactor - 1.0) / 2.0, 0), 1)
  return normalized * 100
})

// SVG 参数
const circumference = computed(() => 2 * Math.PI * sizeConfig.value.radius)
const strokeDashoffset = computed(() => {
  return circumference.value * (1 - progress.value / 100)
})
</script>

<template>
  <div class="health-factor-gauge" :class="size">
    <svg
      :width="sizeConfig.radius * 2 + 20"
      :height="sizeConfig.radius * 2 + 20"
      class="gauge-svg"
    >
      <!-- 背景圆环 -->
      <circle
        :cx="sizeConfig.radius + 10"
        :cy="sizeConfig.radius + 10"
        :r="sizeConfig.radius"
        class="gauge-bg"
        :stroke-width="sizeConfig.strokeWidth"
      />
      <!-- 进度圆环 -->
      <circle
        :cx="sizeConfig.radius + 10"
        :cy="sizeConfig.radius + 10"
        :r="sizeConfig.radius"
        class="gauge-progress"
        :stroke="status.color"
        :stroke-width="sizeConfig.strokeWidth"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="strokeDashoffset"
      />
    </svg>

    <!-- 中心内容 -->
    <div class="gauge-content">
      <component :is="status.icon" :size="size === 'small' ? 16 : size === 'large' ? 28 : 20" :style="{ color: status.color }" weight="fill" />
      <div class="gauge-value font-mono" :style="{ fontSize: sizeConfig.fontSize, color: status.color }">
        {{ healthFactor > 0 ? healthFactor.toFixed(2) : '-' }}
      </div>
      <div class="gauge-label">{{ status.label }}</div>
    </div>
  </div>
</template>

<style scoped>
.health-factor-gauge {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.gauge-svg {
  transform: rotate(-90deg);
}

.gauge-bg {
  fill: none;
  stroke: var(--color-bg-tertiary);
}

.gauge-progress {
  fill: none;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.6s ease, stroke 0.3s ease;
}

.gauge-content {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.gauge-value {
  font-weight: 700;
  line-height: 1;
}

.gauge-label {
  font-size: 0.625rem;
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

/* 尺寸变体 */
.health-factor-gauge.small .gauge-label {
  font-size: 0.5625rem;
}

.health-factor-gauge.large .gauge-label {
  font-size: 0.75rem;
}
</style>
