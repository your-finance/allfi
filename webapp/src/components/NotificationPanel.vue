<script setup>
/**
 * 通知面板组件
 * 从右侧滑出，展示通知列表和操作
 */
import { onMounted, onUnmounted } from 'vue'
import {
  PhX,
  PhBellRinging,
  PhCheckCircle,
  PhChartLine,
  PhWarning,
  PhShieldCheck,
  PhChecks
} from '@phosphor-icons/vue'
import { useNotificationStore } from '../stores/notificationStore'
import { useI18n } from '../composables/useI18n'

const notifStore = useNotificationStore()
const { t } = useI18n()

// 通知类型图标映射
const typeIcons = {
  daily_digest: PhChartLine,
  price_alert: PhWarning,
  asset_change: PhBellRinging,
  security_alert: PhShieldCheck
}

// 通知类型样式映射
const typeColors = {
  daily_digest: 'var(--color-accent-primary)',
  price_alert: 'var(--color-warning)',
  asset_change: 'var(--color-info)',
  security_alert: 'var(--color-error)'
}

// 点击通知：标记已读
const handleClickNotification = (notification) => {
  if (!notification.is_read) {
    notifStore.markAsRead(notification.id)
  }
}

// 标记全部已读
const handleMarkAllRead = () => {
  notifStore.markAllAsRead()
}

// 关闭面板
const handleClose = () => {
  notifStore.closePanel()
}

// ESC 键关闭
const handleKeydown = (e) => {
  if (e.key === 'Escape') {
    handleClose()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Transition name="panel">
    <div v-if="notifStore.isPanelOpen" class="panel-overlay" @click.self="handleClose">
      <div class="panel">
        <!-- 头部 -->
        <div class="panel-header">
          <h3 class="panel-title">{{ t('notification.title') || '通知' }}</h3>
          <div class="panel-actions">
            <button
              v-if="notifStore.hasUnread"
              class="mark-all-btn"
              @click="handleMarkAllRead"
              :title="t('notification.markAllRead') || '全部标记已读'"
            >
              <PhChecks :size="16" />
            </button>
            <button class="close-btn" @click="handleClose">
              <PhX :size="16" />
            </button>
          </div>
        </div>

        <!-- 内容 -->
        <div class="panel-body">
          <!-- 加载中 -->
          <div v-if="notifStore.notificationsLoading" class="panel-loading">
            <span>{{ t('common.loading') || '加载中...' }}</span>
          </div>

          <!-- 空状态 -->
          <div v-else-if="notifStore.notifications.length === 0" class="panel-empty">
            <PhBellRinging :size="32" />
            <span>{{ t('notification.empty') || '暂无通知' }}</span>
          </div>

          <!-- 通知列表 -->
          <div v-else class="notification-list">
            <div
              v-for="notification in notifStore.notifications"
              :key="notification.id"
              class="notification-item"
              :class="{ 'notification-unread': !notification.is_read }"
              @click="handleClickNotification(notification)"
            >
              <!-- 类型图标 -->
              <div class="notification-icon" :style="{ color: typeColors[notification.type] }">
                <component :is="typeIcons[notification.type] || PhBellRinging" :size="18" />
              </div>

              <!-- 内容 -->
              <div class="notification-content">
                <div class="notification-title">{{ notification.title }}</div>
                <div class="notification-text">{{ notification.message }}</div>
                <div class="notification-time">{{ notification.created_at }}</div>
              </div>

              <!-- 已读标记 -->
              <div v-if="notification.is_read" class="notification-read">
                <PhCheckCircle :size="14" />
              </div>
              <div v-else class="notification-unread-dot" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* 遮罩 */
.panel-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 300;
  display: flex;
  justify-content: flex-end;
}

/* 面板 */
.panel {
  width: 380px;
  max-width: 100vw;
  height: 100vh;
  background: var(--color-bg-secondary);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-lg);
}

/* 头部 */
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-md) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
  height: 48px;
}

.panel-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-actions {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.mark-all-btn,
.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.mark-all-btn:hover,
.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

/* 内容区 */
.panel-body {
  flex: 1;
  overflow-y: auto;
}

.panel-loading,
.panel-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  padding: var(--gap-2xl);
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 通知列表 */
.notification-list {
  display: flex;
  flex-direction: column;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-sm);
  padding: var(--gap-md) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.notification-item:hover {
  background: var(--color-bg-tertiary);
}

.notification-unread {
  background: rgba(75, 131, 240, 0.04);
}

.notification-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.notification-text {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-time {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  margin-top: 4px;
  font-family: var(--font-mono);
}

.notification-read {
  flex-shrink: 0;
  color: var(--color-text-muted);
}

.notification-unread-dot {
  flex-shrink: 0;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-accent-primary);
  margin-top: 6px;
}

/* 过渡动画 */
.panel-enter-active,
.panel-leave-active {
  transition: opacity var(--transition-base);
}

.panel-enter-active .panel,
.panel-leave-active .panel {
  transition: transform var(--transition-base);
}

.panel-enter-from,
.panel-leave-to {
  opacity: 0;
}

.panel-enter-from .panel,
.panel-leave-to .panel {
  transform: translateX(100%);
}

/* 移动端 */
@media (max-width: 640px) {
  .panel {
    width: 100vw;
  }
}
</style>
