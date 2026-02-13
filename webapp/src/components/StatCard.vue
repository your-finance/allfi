<script setup>
/**
 * 统计数据卡片组件
 * 用于展示关键 KPI 指标
 */
import { computed } from 'vue'
import { PhTrendUp, PhTrendDown } from '@phosphor-icons/vue'

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  change: {
    type: Number,
    default: null
  },
  changeLabel: {
    type: String,
    default: '24h'
  },
  icon: {
    type: Object,
    default: null
  },
  variant: {
    type: String,
    default: 'default',
    validator: (v) => ['default', 'cyan', 'purple', 'amber', 'emerald'].includes(v)
  }
})

// 变化方向
const isPositive = computed(() => props.change >= 0)
const changeFormatted = computed(() => {
  if (props.change === null) return null
  const sign = isPositive.value ? '+' : ''
  return `${sign}${props.change.toFixed(2)}%`
})

// 图标颜色映射（低饱和度金融风格）
const iconColors = {
  default: '#8A919E',
  cyan: '#4B83F0',
  purple: '#8B7CC8',
  amber: '#C4952B',
  emerald: '#2EBD85'
}
</script>

<template>
  <div class="stat-card">
    <div class="stat-card-content">
      <!-- 图标 -->
      <div 
        v-if="icon" 
        class="stat-icon"
        :style="{ color: iconColors[variant] }"
      >
        <component :is="icon" :size="24" weight="duotone" />
      </div>
      
      <!-- 数据区 -->
      <div class="stat-data">
        <span class="stat-label">{{ label }}</span>
        <span class="stat-value font-mono">{{ value }}</span>
        
        <!-- 变化指示器 -->
        <div 
          v-if="changeFormatted !== null" 
          class="stat-change"
          :class="isPositive ? 'change-positive' : 'change-negative'"
        >
          <PhTrendUp v-if="isPositive" :size="14" weight="bold" />
          <PhTrendDown v-else :size="14" weight="bold" />
          <span>{{ changeFormatted }}</span>
          <span class="change-label">{{ changeLabel }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stat-card {
  position: relative;
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.stat-card-content {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-sm);
}

.stat-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

.stat-data {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.stat-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-change {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 12px;
  font-weight: 500;
  margin-top: 2px;
}

.change-positive {
  color: var(--color-success);
}

.change-negative {
  color: var(--color-error);
}

.change-label {
  color: var(--color-text-muted);
  font-weight: 400;
  margin-left: 2px;
}

@media (max-width: 640px) {
  .stat-value {
    font-size: 20px;
  }
}
</style>
