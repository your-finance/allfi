<script setup>
/**
 * 风险预警面板
 * 展示系统生成的风险提示和建议
 */
import { computed } from 'vue'
import { PhWarning, PhInfo, PhCheckCircle } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  alerts: {
    type: Array,
    default: () => []
  }
})

const { t } = useI18n()

// 严重程度图标和颜色
const severityConfig = {
  critical: { icon: PhWarning, color: 'var(--color-error)' },
  warning: { icon: PhWarning, color: 'var(--color-warning)' },
  info: { icon: PhInfo, color: 'var(--color-info)' }
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (hours < 1) return t('common.justNow')
  if (hours < 24) return t('common.hoursAgo', { hours })
  return t('common.daysAgo', { days })
}
</script>

<template>
  <div class="risk-alert-panel">
    <h3 class="panel-title">{{ t('risk.alerts') }}</h3>

    <div v-if="alerts.length > 0" class="alert-list">
      <div
        v-for="alert in alerts"
        :key="alert.id"
        class="alert-item"
        :class="`severity-${alert.severity}`"
      >
        <component
          :is="severityConfig[alert.severity]?.icon || PhInfo"
          :size="16"
          class="alert-icon"
          :style="{ color: severityConfig[alert.severity]?.color }"
        />
        <div class="alert-content">
          <p class="alert-message">{{ alert.message }}</p>
          <span class="alert-time">{{ formatTime(alert.timestamp) }}</span>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <PhCheckCircle :size="32" style="color: var(--color-success)" />
      <p>{{ t('risk.noAlerts') }}</p>
    </div>
  </div>
</template>

<style scoped>
.risk-alert-panel {
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.panel-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.alert-item {
  display: flex;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  border-left: 3px solid transparent;
}

.alert-item.severity-critical {
  border-left-color: var(--color-error);
}

.alert-item.severity-warning {
  border-left-color: var(--color-warning);
}

.alert-item.severity-info {
  border-left-color: var(--color-info);
}

.alert-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.alert-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.alert-message {
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  line-height: 1.4;
}

.alert-time {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-xl);
  text-align: center;
}

.empty-state p {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
}
</style>
