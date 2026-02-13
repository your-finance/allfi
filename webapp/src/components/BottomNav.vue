<script setup>
/**
 * 移动端底部导航栏
 * 仅在移动端（<768px）显示，提供 5 个核心页面的快速切换
 * 支持 safe-area-inset-bottom（iPhone 刘海/底部安全区）
 */
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PhChartPieSlice,
  PhWallet,
  PhClockCounterClockwise,
  PhChartLine,
  PhGear
} from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

// 底部导航项（5 个核心页面）
const tabs = computed(() => [
  { path: '/dashboard', labelKey: 'nav.dashboard', icon: PhChartPieSlice },
  { path: '/accounts', labelKey: 'nav.accounts', icon: PhWallet },
  { path: '/history', labelKey: 'nav.history', icon: PhClockCounterClockwise },
  { path: '/analytics', labelKey: 'nav.analytics', icon: PhChartLine },
  { path: '/settings', labelKey: 'nav.settings', icon: PhGear }
])

const isActive = (path) => route.path === path

const navigateTo = (path) => {
  if (route.path !== path) {
    router.push(path)
  }
}
</script>

<template>
  <nav class="bottom-nav" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.path"
      class="bottom-nav-item"
      :class="{ 'bottom-nav-item-active': isActive(tab.path) }"
      role="tab"
      :aria-selected="isActive(tab.path)"
      :aria-label="t(tab.labelKey)"
      @click="navigateTo(tab.path)"
    >
      <component :is="tab.icon" :size="20" :weight="isActive(tab.path) ? 'fill' : 'regular'" />
      <span class="bottom-nav-label">{{ t(tab.labelKey) }}</span>
    </button>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: stretch;
  background: var(--color-bg-secondary);
  border-top: 1px solid var(--color-border);
  z-index: 100;
  /* iPhone 安全区适配 */
  padding-bottom: env(safe-area-inset-bottom, 0px);
}

.bottom-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  /* 最小触摸区域 44px */
  min-height: 52px;
  padding: 6px 0;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
  user-select: none;
}

.bottom-nav-item:active {
  background: var(--color-bg-tertiary);
}

.bottom-nav-item-active {
  color: var(--color-accent-primary);
}

.bottom-nav-label {
  font-size: 0.625rem;
  font-weight: 500;
  line-height: 1;
  white-space: nowrap;
}
</style>
