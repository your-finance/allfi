<script setup>
/**
 * 下拉刷新指示器组件
 * 显示下拉进度和刷新状态，配合 usePullToRefresh 使用
 */
import { computed } from 'vue'
import { PhArrowsClockwise } from '@phosphor-icons/vue'

const props = defineProps({
  /** 当前下拉距离（px） */
  pullDistance: { type: Number, default: 0 },
  /** 是否正在刷新 */
  isRefreshing: { type: Boolean, default: false },
  /** 触发阈值（px），用于计算进度 */
  threshold: { type: Number, default: 60 }
})

// 下拉进度（0-1）
const progress = computed(() => Math.min(props.pullDistance / (props.threshold / 2.5), 1))

// 是否可见
const isVisible = computed(() => props.pullDistance > 2 || props.isRefreshing)

// 图标旋转角度
const rotation = computed(() => progress.value * 180)
</script>

<template>
  <Transition name="pull-indicator">
    <div v-if="isVisible" class="pull-indicator" :style="{ height: pullDistance + 'px' }">
      <div class="pull-icon" :class="{ 'pull-refreshing': isRefreshing }">
        <PhArrowsClockwise
          :size="20"
          :style="{ transform: `rotate(${rotation}deg)` }"
          :class="{ 'animate-spin': isRefreshing }"
        />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.pull-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: height 0.2s ease;
}

.pull-icon {
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}

.pull-refreshing {
  color: var(--color-accent-primary);
}

.pull-indicator-leave-active {
  transition: height 0.2s ease, opacity 0.2s ease;
}

.pull-indicator-leave-to {
  height: 0 !important;
  opacity: 0;
}
</style>
