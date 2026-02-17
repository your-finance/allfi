<script setup>
/**
 * 仪表盘 Widget 配置面板
 * 齿轮图标点击弹出，控制各 Widget 的显示/隐藏
 */
import { useDashboardStore } from '../stores/dashboardStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: Boolean,
})

const emit = defineEmits(['close'])

const dashboardStore = useDashboardStore()
const { t } = useI18n()

// Widget 定义列表
const widgets = [
  { key: 'assetSummary', labelKey: 'widgets.assetSummary', icon: '📊' },
  { key: 'trend', labelKey: 'widgets.trend', icon: '📈' },
  { key: 'distribution', labelKey: 'widgets.distribution', icon: '🍩' },
  { key: 'healthScore', labelKey: 'widgets.healthScore', icon: '💚' },
  { key: 'goals', labelKey: 'widgets.goals', icon: '🎯' },
  { key: 'holdings', labelKey: 'widgets.holdings', icon: '📋' },
  { key: 'defiOverview', labelKey: 'widgets.defiOverview', icon: '🔗' },
  { key: 'nftOverview', labelKey: 'widgets.nftOverview', icon: '🖼' },
  { key: 'feeAnalytics', labelKey: 'widgets.feeAnalytics', icon: '💰' },
]

const handleToggle = (key) => {
  dashboardStore.toggleWidget(key)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="customizer-overlay" @click.self="emit('close')">
        <div class="customizer-panel">
          <div class="customizer-header">
            <h3 class="customizer-title">{{ t('widgets.customize') }}</h3>
            <button class="close-btn" @click="emit('close')">&times;</button>
          </div>

          <div class="widget-list">
            <label
              v-for="w in widgets"
              :key="w.key"
              class="widget-item"
              :class="{ disabled: dashboardStore.widgetConfig[w.key] && !dashboardStore.canDisable }"
            >
              <span class="widget-icon">{{ w.icon }}</span>
              <span class="widget-name">{{ t(w.labelKey) }}</span>
              <input
                type="checkbox"
                class="widget-toggle"
                :checked="dashboardStore.widgetConfig[w.key]"
                :disabled="dashboardStore.widgetConfig[w.key] && !dashboardStore.canDisable"
                @change="handleToggle(w.key)"
              />
            </label>
          </div>

          <div class="customizer-footer">
            <button class="btn btn-ghost btn-sm" @click="dashboardStore.resetLayout()">
              {{ t('widgets.resetDefault') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.customizer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 400;
  display: flex;
  justify-content: flex-end;
}

.customizer-panel {
  width: 320px;
  max-width: 90vw;
  height: 100%;
  background: var(--color-bg-secondary);
  border-left: 1px solid var(--color-border);
  padding: var(--gap-xl);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.customizer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--gap-lg);
}

.customizer-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 1.125rem;
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: background var(--transition-fast), color var(--transition-fast);
}

.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

/* Widget 列表 */
.widget-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
  flex: 1;
}

.widget-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.widget-item:hover {
  border-color: var(--color-accent-primary);
}

.widget-item.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.widget-icon {
  font-size: 1rem;
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.widget-name {
  flex: 1;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.widget-toggle {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent-primary);
  cursor: pointer;
}

.widget-toggle:disabled {
  cursor: not-allowed;
}

/* 底部 */
.customizer-footer {
  margin-top: var(--gap-lg);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: flex-end;
}

/* 过渡 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 120ms ease;
}

.modal-enter-active .customizer-panel,
.modal-leave-active .customizer-panel {
  transition: transform 180ms ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .customizer-panel {
  transform: translateX(100%);
}

.modal-leave-to .customizer-panel {
  transform: translateX(100%);
}
</style>
