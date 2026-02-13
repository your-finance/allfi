<script setup>
/**
 * 全局 Toast 通知容器
 * 固定在右上角，从上往下堆叠，进入/退出动画
 */
import { useToast } from '../composables/useToast'

const { toasts, removeToast } = useToast()

// 类型对应的图标
const typeIcons = {
  success: '✓',
  error: '✕',
  warning: '⚠',
  info: 'ℹ',
}
</script>

<template>
  <Teleport to="body">
    <div class="toast-container" aria-live="polite">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="toast-item"
          :class="`toast-${toast.type}`"
          @click="removeToast(toast.id)"
        >
          <span class="toast-icon">{{ typeIcons[toast.type] }}</span>
          <span class="toast-message">{{ toast.message }}</span>
          <button class="toast-close" @click.stop="removeToast(toast.id)">&times;</button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: var(--gap-lg);
  right: var(--gap-lg);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  max-width: 360px;
  pointer-events: none;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  border-radius: var(--radius-sm);
  border: 1px solid;
  font-size: 0.8125rem;
  cursor: pointer;
  pointer-events: auto;
  min-width: 240px;
}

/* 类型样式 */
.toast-success {
  background: color-mix(in srgb, var(--color-success) 15%, var(--color-bg-secondary));
  border-color: color-mix(in srgb, var(--color-success) 40%, transparent);
  color: var(--color-success);
}

.toast-error {
  background: color-mix(in srgb, var(--color-error) 15%, var(--color-bg-secondary));
  border-color: color-mix(in srgb, var(--color-error) 40%, transparent);
  color: var(--color-error);
}

.toast-warning {
  background: color-mix(in srgb, var(--color-warning) 15%, var(--color-bg-secondary));
  border-color: color-mix(in srgb, var(--color-warning) 40%, transparent);
  color: var(--color-warning);
}

.toast-info {
  background: color-mix(in srgb, var(--color-info) 15%, var(--color-bg-secondary));
  border-color: color-mix(in srgb, var(--color-info) 40%, transparent);
  color: var(--color-info);
}

.toast-icon {
  font-size: 0.875rem;
  font-weight: 700;
  flex-shrink: 0;
  width: 18px;
  text-align: center;
}

.toast-message {
  flex: 1;
  color: var(--color-text-primary);
  line-height: 1.4;
}

.toast-close {
  background: none;
  border: none;
  color: var(--color-text-muted);
  font-size: 1rem;
  cursor: pointer;
  padding: 0 4px;
  flex-shrink: 0;
  transition: color var(--transition-fast);
}

.toast-close:hover {
  color: var(--color-text-primary);
}

/* 进入/退出动画 */
.toast-enter-active {
  transition: all 200ms ease-out;
}

.toast-leave-active {
  transition: all 150ms ease-in;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(40px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(40px);
}

.toast-move {
  transition: transform 200ms ease;
}
</style>
