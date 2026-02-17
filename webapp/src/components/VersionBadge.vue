<script setup>
/**
 * 版本徽章组件
 * 显示在侧边栏 logo 下方，点击弹出 Popover 显示版本信息和检查更新
 */
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { PhArrowsClockwise, PhArrowSquareOut, PhInfo, PhDownloadSimple } from '@phosphor-icons/vue'
import { useSystemStore } from '../stores/systemStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  collapsed: Boolean
})

const { t } = useI18n()
const systemStore = useSystemStore()

// Popover 显示状态
const showPopover = ref(false)
// 组件根元素引用
const badgeRef = ref(null)

// 切换 Popover
const togglePopover = () => {
  showPopover.value = !showPopover.value
}

// 手动检查更新
const handleCheck = async () => {
  await systemStore.checkForUpdate()
}

// 执行更新
const handleUpdate = async () => {
  if (systemStore.updateInfo?.latest_version) {
    await systemStore.applyUpdate(systemStore.updateInfo.latest_version)
  }
}

// 点击外部关闭 Popover
const handleClickOutside = (e) => {
  if (badgeRef.value && !badgeRef.value.contains(e.target)) {
    showPopover.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 当前版本号：优先使用后端 API 返回的版本，降级使用构建时注入的版本
const appVersion = computed(() => {
  return systemStore.versionInfo?.version || __APP_VERSION__
})

// GitHub Releases 地址
const releasesUrl = 'https://github.com/your-finance/allfi/releases'
</script>

<template>
  <div ref="badgeRef" class="version-badge-wrapper">
    <!-- 折叠状态：图标 + tooltip -->
    <button
      v-if="collapsed"
      class="version-badge-icon"
      :title="`v${appVersion}`"
      @click="togglePopover"
    >
      <PhInfo :size="14" />
      <span v-if="systemStore.hasUpdate" class="update-dot" />
    </button>

    <!-- 展开状态：版本号文字 -->
    <button
      v-else
      class="version-badge"
      @click="togglePopover"
    >
      <span class="version-badge-text">v{{ appVersion }}</span>
      <span v-if="systemStore.hasUpdate" class="update-dot" />
    </button>

    <!-- Popover 弹出层 -->
    <Transition name="popover-fade">
      <div v-if="showPopover" class="version-popover">
        <!-- 当前版本 -->
        <div class="popover-row">
          <span class="popover-label">{{ t('system.currentVersion') }}</span>
          <span class="popover-value">v{{ appVersion }}</span>
        </div>

        <!-- 最新版本（检查后显示） -->
        <div v-if="systemStore.updateInfo" class="popover-row">
          <span class="popover-label">{{ t('system.latestVersion') }}</span>
          <span class="popover-value">v{{ systemStore.updateInfo.latest_version }}</span>
        </div>

        <!-- 状态文字 -->
        <div class="popover-status">
          <span v-if="systemStore.isChecking" class="status-checking">
            {{ t('system.checkUpdate') }}...
          </span>
          <span v-else-if="systemStore.hasUpdate" class="status-has-update">
            {{ t('system.newVersionAvailable') }}
          </span>
          <span v-else-if="systemStore.updateInfo" class="status-up-to-date">
            {{ t('system.upToDate') }}
          </span>
        </div>

        <!-- 操作按钮 -->
        <div class="popover-actions">
          <!-- 立即更新按钮（有新版本时显示） -->
          <button
            v-if="systemStore.hasUpdate"
            class="update-now-btn"
            :disabled="systemStore.isUpdating"
            @click="handleUpdate"
          >
            <PhDownloadSimple :size="13" />
            <span>{{ systemStore.isUpdating ? t('system.updating') : t('system.updateNow') }}</span>
          </button>

          <button
            class="check-btn"
            :disabled="systemStore.isChecking"
            @click="handleCheck"
          >
            <PhArrowsClockwise
              :size="13"
              :class="{ 'spin': systemStore.isChecking }"
            />
            <span>{{ t('system.checkUpdate') }}</span>
          </button>

          <a
            class="releases-link"
            :href="releasesUrl"
            target="_blank"
            rel="noopener noreferrer"
          >
            <PhArrowSquareOut :size="13" />
            <span>{{ t('system.viewReleases') }}</span>
          </a>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.version-badge-wrapper {
  position: relative;
}

/* 展开状态的版本徽章 */
.version-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 8px;
  border-radius: var(--radius-sm);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.version-badge:hover {
  border-color: var(--color-border-hover);
  background: var(--color-bg-secondary);
}

.version-badge-text {
  font-size: 10px;
  font-family: var(--font-mono);
  color: var(--color-text-muted);
  line-height: 1.6;
}

/* 折叠状态的图标按钮 */
.version-badge-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast), border-color var(--transition-fast);
}

.version-badge-icon:hover {
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

/* 更新指示圆点 */
.update-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-accent-primary);
  flex-shrink: 0;
}

.version-badge-icon .update-dot {
  position: absolute;
  top: 4px;
  right: 4px;
}

/* Popover 弹出层 */
.version-popover {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 1000;
  min-width: 220px;
  padding: var(--gap-sm);
  background: var(--color-bg-primary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.popover-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.popover-label {
  font-size: 12px;
  color: var(--color-text-muted);
}

.popover-value {
  font-size: 12px;
  font-family: var(--font-mono);
  color: var(--color-text-primary);
  font-weight: 500;
}

/* 状态文字 */
.popover-status {
  padding: 6px 0;
  font-size: 11px;
  text-align: center;
  border-top: 1px solid var(--color-border);
  margin-top: 4px;
}

.status-checking {
  color: var(--color-text-muted);
}

.status-has-update {
  color: var(--color-accent-primary);
  font-weight: 500;
}

.status-up-to-date {
  color: var(--color-success, #22c55e);
}

/* 操作按钮区 */
.popover-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 6px;
  border-top: 1px solid var(--color-border);
  margin-top: 4px;
}

/* 立即更新按钮（醒目样式） */
.update-now-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), opacity var(--transition-fast);
  background: var(--color-accent-primary);
  color: #fff;
  border: none;
}

.update-now-btn:hover {
  opacity: 0.9;
}

.update-now-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.check-btn,
.releases-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  cursor: pointer;
  transition: background var(--transition-fast);
  text-decoration: none;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
}

.check-btn:hover,
.releases-link:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.check-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 旋转动画 */
.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Popover 过渡动画 */
.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
