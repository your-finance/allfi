<script setup>
/**
 * 命令面板组件
 * Cmd+K / Ctrl+K 触发，支持模糊搜索、键盘导航、命令执行
 * 增强：资产搜索、快捷筛选前缀 (@来源 #类型 >价值)
 */
import { ref, watch, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  PhMagnifyingGlass,
  PhArrowsClockwise,
  PhPlus,
  PhWallet,
  PhPaintBrush,
  PhChartPieSlice,
  PhGear,
  PhClockCounterClockwise,
  PhChartLine,
  PhArrowRight,
  PhCommand,
  PhEyeSlash,
  PhCurrencyBtc,
  PhCubeTransparent
} from '@phosphor-icons/vue'
import { useCommandStore } from '../stores/commandStore'
import { useAssetStore } from '../stores/assetStore'
import { useThemeStore } from '../stores/themeStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const cmdStore = useCommandStore()
const assetStore = useAssetStore()
const themeStore = useThemeStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

const inputRef = ref(null)

// 图标映射
const iconMap = {
  ArrowsClockwise: PhArrowsClockwise,
  Plus: PhPlus,
  Wallet: PhWallet,
  PaintBrush: PhPaintBrush,
  ChartPieSlice: PhChartPieSlice,
  Gear: PhGear,
  ClockCounterClockwise: PhClockCounterClockwise,
  ChartLine: PhChartLine,
  EyeSlash: PhEyeSlash,
  CurrencyBtc: PhCurrencyBtc,
  CubeTransparent: PhCubeTransparent,
}

// 面板打开时自动聚焦输入框
watch(() => cmdStore.isOpen, (val) => {
  if (val) {
    nextTick(() => inputRef.value?.focus())
  }
})

// 搜索内容变化时重置选中索引
watch(() => cmdStore.searchQuery, () => {
  cmdStore.selectedIndex = 0
})

// 是否有搜索前缀（显示筛选提示）
const hasPrefix = computed(() => {
  const q = cmdStore.searchQuery.trim()
  return q.startsWith('@') || q.startsWith('#') || q.startsWith('>')
})

// 获取类型标签
const getTypeLabel = (type) => {
  const labels = {
    page: t('commandPalette.page'),
    command: t('commandPalette.action'),
    asset: t('commandPalette.asset')
  }
  return labels[type] || type
}

// 键盘事件处理
const handleKeydown = (e) => {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    cmdStore.moveDown()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    cmdStore.moveUp()
  } else if (e.key === 'Enter') {
    e.preventDefault()
    executeSelected()
  } else if (e.key === 'Escape') {
    cmdStore.close()
  }
}

// 执行选中的命令
const executeSelected = () => {
  const item = cmdStore.getSelectedItem()
  if (!item) return

  if (item.type === 'page') {
    router.push(item.path)
  } else if (item.type === 'command') {
    executeCommand(item.id)
  } else if (item.type === 'asset') {
    // 资产搜索结果：跳转到对应页面
    router.push('/accounts')
  }

  cmdStore.close()
}

// 执行命令动作
const executeCommand = (id) => {
  switch (id) {
    case 'refresh':
      assetStore.loadSummary()
      break
    case 'addExchange':
      router.push('/accounts')
      break
    case 'addWallet':
      router.push('/accounts')
      break
    case 'toggleTheme':
      document.body.classList.add('theme-transition')
      themeStore.nextTheme()
      setTimeout(() => document.body.classList.remove('theme-transition'), 280)
      break
    case 'togglePrivacy':
      themeStore.togglePrivacyMode()
      break
  }
}

// 点击条目执行
const handleClickItem = (item, index) => {
  cmdStore.selectedIndex = index
  executeSelected()
}
</script>

<template>
  <Transition name="cmd">
    <div v-if="cmdStore.isOpen" class="cmd-overlay" @click.self="cmdStore.close()">
      <div class="cmd-panel" @keydown="handleKeydown">
        <!-- 搜索框 -->
        <div class="cmd-input-wrap">
          <PhMagnifyingGlass :size="16" class="cmd-search-icon" />
          <input
            ref="inputRef"
            v-model="cmdStore.searchQuery"
            type="text"
            class="cmd-input"
            :placeholder="t('commandPalette.placeholder') || '搜索命令、页面或资产...'"
            autocomplete="off"
          >
          <kbd class="cmd-kbd">ESC</kbd>
        </div>

        <!-- 快捷筛选提示 -->
        <div v-if="!cmdStore.searchQuery" class="cmd-hints-bar">
          <span class="prefix-hint">
            <kbd>@</kbd> {{ t('commandPalette.filterBySource') }}
          </span>
          <span class="prefix-hint">
            <kbd>#</kbd> {{ t('commandPalette.filterByType') }}
          </span>
          <span class="prefix-hint">
            <kbd>&gt;</kbd> {{ t('commandPalette.filterByValue') }}
          </span>
        </div>

        <!-- 结果列表 -->
        <div class="cmd-list" v-if="cmdStore.filteredItems.length > 0">
          <div
            v-for="(item, index) in cmdStore.filteredItems"
            :key="item.id"
            class="cmd-item"
            :class="{ 'cmd-item-active': index === cmdStore.selectedIndex }"
            @click="handleClickItem(item, index)"
            @mouseenter="cmdStore.selectedIndex = index"
          >
            <component
              :is="iconMap[item.icon] || PhArrowRight"
              :size="16"
              class="cmd-item-icon"
            />
            <div class="cmd-item-content">
              <span class="cmd-item-label">{{ item.label }}</span>
              <!-- 资产附加信息：余额和来源 -->
              <span v-if="item.type === 'asset' && item.value" class="cmd-item-meta font-mono">
                ${{ formatNumber(item.value) }}
                <span class="cmd-item-source">{{ item.source }}</span>
              </span>
            </div>
            <span class="cmd-item-type" :class="`type-${item.type}`">
              {{ getTypeLabel(item.type) }}
            </span>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="cmd-empty">
          {{ t('commandPalette.empty') || '没有匹配的结果' }}
        </div>

        <!-- 底部提示 -->
        <div class="cmd-footer">
          <span class="cmd-hint">
            <kbd>↑</kbd><kbd>↓</kbd> {{ t('commandPalette.navigate') || '导航' }}
          </span>
          <span class="cmd-hint">
            <kbd>↵</kbd> {{ t('commandPalette.execute') || '执行' }}
          </span>
          <span class="cmd-hint">
            <kbd>ESC</kbd> {{ t('commandPalette.close') || '关闭' }}
          </span>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* 遮罩 */
.cmd-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 400;
  display: flex;
  justify-content: center;
  padding-top: 20vh;
}

/* 面板 */
.cmd-panel {
  width: 520px;
  max-width: 90vw;
  max-height: 460px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  align-self: flex-start;
}

/* 搜索框 */
.cmd-input-wrap {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-md) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.cmd-search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.cmd-input {
  flex: 1;
  border: none;
  background: transparent;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  font-family: var(--font-body);
}

.cmd-input::placeholder {
  color: var(--color-text-muted);
}

.cmd-kbd {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  background: var(--color-bg-tertiary);
  color: var(--color-text-muted);
  border: 1px solid var(--color-border);
  font-family: var(--font-mono);
}

/* 快捷筛选提示 */
.cmd-hints-bar {
  display: flex;
  gap: var(--gap-md);
  padding: var(--gap-xs) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-tertiary);
}

.prefix-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.prefix-hint kbd {
  font-size: 0.625rem;
  padding: 1px 4px;
  border-radius: 2px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  font-family: var(--font-mono);
  font-weight: 600;
  color: var(--color-accent-primary);
}

/* 结果列表 */
.cmd-list {
  flex: 1;
  overflow-y: auto;
  padding: var(--gap-xs) 0;
}

.cmd-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-lg);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.cmd-item:hover,
.cmd-item-active {
  background: var(--color-bg-tertiary);
}

.cmd-item-icon {
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.cmd-item-active .cmd-item-icon {
  color: var(--color-accent-primary);
}

.cmd-item-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  min-width: 0;
}

.cmd-item-label {
  font-size: 13px;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cmd-item-meta {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
}

.cmd-item-source {
  color: var(--color-text-muted);
  opacity: 0.7;
  margin-left: 4px;
}

.cmd-item-type {
  font-size: 11px;
  color: var(--color-text-muted);
  padding: 1px 6px;
  border-radius: var(--radius-xs);
  background: var(--color-bg-tertiary);
  white-space: nowrap;
}

.cmd-item-type.type-asset {
  background: rgba(139, 92, 246, 0.12);
  color: #8B5CF6;
}

/* 空状态 */
.cmd-empty {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 13px;
}

/* 底部提示 */
.cmd-footer {
  display: flex;
  align-items: center;
  gap: var(--gap-lg);
  padding: var(--gap-sm) var(--gap-lg);
  border-top: 1px solid var(--color-border);
}

.cmd-hint {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 11px;
  color: var(--color-text-muted);
}

.cmd-hint kbd {
  font-size: 10px;
  padding: 1px 4px;
  border-radius: 2px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  font-family: var(--font-mono);
}

/* 过渡动画 */
.cmd-enter-active,
.cmd-leave-active {
  transition: opacity 120ms ease;
}

.cmd-enter-active .cmd-panel,
.cmd-leave-active .cmd-panel {
  transition: transform 120ms ease, opacity 120ms ease;
}

.cmd-enter-from,
.cmd-leave-to {
  opacity: 0;
}

.cmd-enter-from .cmd-panel {
  transform: scale(0.96) translateY(-8px);
  opacity: 0;
}

.cmd-leave-to .cmd-panel {
  transform: scale(0.96);
  opacity: 0;
}
</style>
