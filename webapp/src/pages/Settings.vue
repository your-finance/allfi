<script setup>
/**
 * Settings 页面 - 系统设置
 * 管理用户偏好、主题切换、显示设置、安全配置等
 */
import { ref, computed, onMounted, watch } from 'vue'
import {
  PhPalette,
  PhShieldCheck,
  PhKey,
  PhBell,
  PhDatabase,
  PhTrash,
  PhDownloadSimple,
  PhArrowsClockwise,
  PhCheck,
  PhMoon,
  PhSun,
  PhCaretRight,
  PhPaintBrush,
  PhQrCode,
  PhCopy,
  PhWarning,
  PhLock,
  PhTrophy,
  PhClockClockwise,
  PhInfo,
  PhArrowDown,
  PhArrowCounterClockwise
} from '@phosphor-icons/vue'
import AchievementPanel from '../components/AchievementPanel.vue'
import { useAchievementStore } from '../stores/achievementStore'
import { useSystemStore } from '../stores/systemStore'
import { useThemeStore } from '../stores/themeStore'
import { useAuthStore } from '../stores/authStore'
import { useNotificationStore } from '../stores/notificationStore'
import { useI18n } from '../composables/useI18n'
import { useExportData } from '../composables/useExportData'
import { useToast } from '../composables/useToast'
import { settingsService } from '../api/index.js'
import { transactionService } from '../api/transactionService.js'

// 使用主题 Store 和多语言
const themeStore = useThemeStore()
const authStore = useAuthStore()
const notifStore = useNotificationStore()
const { t } = useI18n()
const { exportAsCSV, exportAsJSON, getDefaultFilename } = useExportData()
const { showToast } = useToast()

// 通知偏好保存状态
const achStore = useAchievementStore()
const systemStore = useSystemStore()
const notifSaving = ref(false)

// 成就面板
const showAchievements = ref(false)

// 更新通知偏好
const updateNotifPreference = async (field, value) => {
  notifSaving.value = true
  try {
    await notifStore.updatePreferences({ [field]: value })
    showToast(t('settings.saveSuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    notifSaving.value = false
  }
}

// 页面加载时获取通知偏好和系统信息
onMounted(() => {
  notifStore.loadPreferences()
  achStore.fetchAchievements()
  systemStore.loadVersion()
  systemStore.checkForUpdate()
  systemStore.loadUpdateHistory()
  systemStore.startAutoCheck()
})

// 操作状态
const isExporting = ref(false)
const isClearing = ref(false)
const isResetting = ref(false)

// 导出格式弹窗
const showExportDialog = ref(false)

// ========== 2FA 状态管理 ==========
const is2FAEnabled = computed(() => authStore.user?.has2FA || false)
const show2FASetupModal = ref(false)
const show2FADisableModal = ref(false)
const setupStep = ref(1) // 1: 显示二维码, 2: 验证
const verifyCode = ref('')
const verifyError = ref('')
const isVerifying = ref(false)
const disableCode = ref('')
const disableError = ref('')
const isDisabling = ref(false)
const codeCopied = ref(false)

// 模拟密钥
const twoFASecret = ref('JBSWY3DPEHPK3PXP')

// 开始设置2FA
const startSetup2FA = () => {
  show2FASetupModal.value = true
  setupStep.value = 1
  verifyCode.value = ''
  verifyError.value = ''
}

// 复制密钥
const copySecret = () => {
  navigator.clipboard.writeText(twoFASecret.value)
  codeCopied.value = true
  setTimeout(() => codeCopied.value = false, 2000)
}

// 验证并启用2FA
const verifyAndEnable2FA = async () => {
  if (verifyCode.value.length !== 6) {
    verifyError.value = t('settings.twoFA.invalidCode')
    return
  }

  isVerifying.value = true
  verifyError.value = ''

  try {
    // 调用 authStore 启用 2FA
    await authStore.enable2FA(verifyCode.value)
    show2FASetupModal.value = false
  } catch (err) {
    verifyError.value = err.message || t('settings.twoFA.verifyFailed')
  } finally {
    isVerifying.value = false
  }
}

// 打开禁用2FA对话框
const openDisable2FA = () => {
  show2FADisableModal.value = true
  disableCode.value = ''
  disableError.value = ''
}

// 禁用2FA
const disable2FA = async () => {
  if (disableCode.value.length !== 6) {
    disableError.value = t('settings.twoFA.invalidCode')
    return
  }

  isDisabling.value = true
  disableError.value = ''

  try {
    // 调用 authStore 禁用 2FA
    await authStore.disable2FA(disableCode.value)
    show2FADisableModal.value = false
  } catch (err) {
    disableError.value = err.message || t('settings.twoFA.disableFailed')
  } finally {
    isDisabling.value = false
  }
}

// 设置状态
const settings = ref({
  // 显示设置
  language: 'zh-CN',
  currency: 'USDC',
  
  // 数据设置
  autoRefresh: true,
  refreshInterval: 300,
  historyRetention: 90,

  // 安全设置
  showBalances: true,
  confirmOperations: true
})

// 保存状态
const isSaving = ref(false)
const saveSuccess = ref(false)

// 语言选项（从store获取）
const languageOptions = themeStore.availableLanguages.map(lang => ({
  value: lang.code,
  label: lang.name
}))

// 货币选项
const currencyOptions = [
  { value: 'USDC', label: 'USDC', symbol: '$' },
  { value: 'BTC', label: 'Bitcoin', symbol: '₿' },
  { value: 'ETH', label: 'Ethereum', symbol: 'Ξ' },
  { value: 'CNY', label: '人民币', symbol: '¥' }
]

// 刷新间隔选项（使用翻译）
const refreshIntervalOptions = computed(() => [
  { value: 60, label: `1 ${t('settings.minute')}` },
  { value: 300, label: `5 ${t('settings.minute')}` },
  { value: 600, label: `10 ${t('settings.minute')}` },
  { value: 1800, label: `30 ${t('settings.minute')}` },
  { value: 3600, label: `1 ${t('settings.hour')}` }
])

// 保留天数选项（使用翻译）
const retentionOptions = computed(() => [
  { value: 30, label: `30 ${t('settings.day')}` },
  { value: 60, label: `60 ${t('settings.day')}` },
  { value: 90, label: `90 ${t('settings.day')}` },
  { value: 180, label: `180 ${t('settings.day')}` },
  { value: 365, label: `1 ${t('settings.year')}` }
])

// 获取字体名称（从CSS字体栈中提取第一个）
const getFontName = (fontStack) => {
  if (!fontStack) return ''
  // 提取第一个字体名称，去除引号
  const match = fontStack.match(/['"]?([^'"]+)['"]?/)
  return match ? match[1] : fontStack
}

// 切换主题
const selectTheme = (themeId) => {
  // 添加过渡类
  document.body.classList.add('theme-transition')
  themeStore.setTheme(themeId)
  
  // 移除过渡类
  setTimeout(() => {
    document.body.classList.remove('theme-transition')
  }, 300)
}

// 保存设置
const saveSettings = async () => {
  isSaving.value = true
  await new Promise(resolve => setTimeout(resolve, 800))
  isSaving.value = false
  saveSuccess.value = true
  setTimeout(() => {
    saveSuccess.value = false
  }, 2000)
}

// 打开导出格式选择弹窗
const openExportDialog = () => {
  showExportDialog.value = true
}

// 执行导出
const doExport = (format) => {
  isExporting.value = true
  try {
    const filename = getDefaultFilename()
    if (format === 'csv') {
      exportAsCSV(filename)
    } else {
      exportAsJSON(filename)
    }
    showToast(t('export.success'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    isExporting.value = false
    showExportDialog.value = false
  }
}

// 清除缓存
const clearCache = async () => {
  isClearing.value = true
  try {
    await settingsService.clearCache()
    showToast(t('settings.cacheCleared'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    isClearing.value = false
  }
}

// 重置设置
const resetSettings = async () => {
  if (!confirm(t('settings.resetConfirm') || '确定要重置所有设置吗？此操作不可撤销。')) {
    return
  }

  isResetting.value = true
  try {
    await settingsService.resetSettings()
    themeStore.resetTheme()
    // 重置本地设置
    settings.value = {
      language: 'zh-CN',
      currency: 'USDC',
      autoRefresh: true,
      refreshInterval: 300,
      historyRetention: 90,
      showBalances: true,
      confirmOperations: true
    }
    showToast(t('settings.resetSuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    isResetting.value = false
  }
}

// 处理版本回滚
const handleRollback = (version) => {
  const msg = t('system.rollbackConfirm').replace('{version}', version)
  if (!confirm(msg)) return
  systemStore.rollback(version)
}

// ========== 交易同步设置 ==========
const txSyncSettings = ref({
  enabled: false,
  interval_minutes: 360,
  lookback_days: 90
})
const txSyncLoading = ref(false)
const txSyncSaving = ref(false)
const txSyncing = ref(false)

// 同步间隔选项
const txSyncIntervalOptions = computed(() => [
  { value: 60, label: t('settings.txSyncInterval1h') },
  { value: 360, label: t('settings.txSyncInterval6h') },
  { value: 1440, label: t('settings.txSyncInterval24h') }
])

// 回溯天数选项
const txSyncLookbackOptions = computed(() => [
  { value: 30, label: `30 ${t('settings.day')}` },
  { value: 60, label: `60 ${t('settings.day')}` },
  { value: 90, label: `90 ${t('settings.day')}` },
  { value: 180, label: `180 ${t('settings.day')}` },
  { value: 365, label: `1 ${t('settings.year')}` }
])

// 加载交易同步设置
const loadTxSyncSettings = async () => {
  txSyncLoading.value = true
  try {
    const data = await transactionService.getSyncSettings()
    txSyncSettings.value = {
      enabled: data.enabled || false,
      interval_minutes: data.interval_minutes || 360,
      lookback_days: data.lookback_days || 90
    }
  } catch {
    // 静默失败，使用默认值
  } finally {
    txSyncLoading.value = false
  }
}

// 更新交易同步设置
const updateTxSyncSetting = async (field, value) => {
  txSyncSaving.value = true
  try {
    txSyncSettings.value[field] = value
    await transactionService.updateSyncSettings({ [field]: value })
    showToast(t('settings.saveSuccess'), 'success')
  } catch {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    txSyncSaving.value = false
  }
}

// 手动触发同步
const triggerSync = async () => {
  txSyncing.value = true
  try {
    const result = await transactionService.syncTransactions()
    const count = result?.synced_count || 0
    showToast(t('settings.txSyncSuccess').replace('{count}', count), 'success')
  } catch {
    showToast(t('settings.txSyncFailed'), 'error')
  } finally {
    txSyncing.value = false
  }
}

// 页面加载时也获取同步设置
onMounted(() => {
  loadTxSyncSettings()
})
</script>

<template>
  <div class="settings-page">
    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-info">
        <p class="header-subtitle">{{ t('settings.subtitle') }}</p>
      </div>
      <button 
        class="btn btn-primary"
        :disabled="isSaving"
        @click="saveSettings"
      >
        <PhArrowsClockwise v-if="isSaving" :size="18" class="spin" />
        <PhCheck v-else-if="saveSuccess" :size="18" />
        <span>{{ saveSuccess ? t('common.save') : t('settings.saveSettings') }}</span>
      </button>
    </header>

    <div class="settings-grid">
      <!-- 🎨 主题设置 - 新增! -->
      <section class="settings-section full-width">
        <div class="section-header">
          <PhPaintBrush :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.themeSection') }}</h2>
        </div>
        
        <div class="themes-grid">
          <button
            v-for="theme in themeStore.availableThemes"
            :key="theme.id"
            class="theme-card"
            :class="{ active: themeStore.currentThemeId === theme.id }"
            @click="selectTheme(theme.id)"
          >
            <!-- 预览色块 -->
            <div class="theme-preview" :style="{ background: theme.preview }">
              <span class="theme-icon">{{ theme.icon }}</span>
            </div>
            
            <!-- 主题信息 -->
            <div class="theme-info">
              <span class="theme-name" :style="{ fontFamily: theme.fonts?.heading }">{{ theme.name }}</span>
              <span class="theme-desc">{{ theme.description }}</span>
              <span class="theme-font">{{ getFontName(theme.fonts?.heading) }}</span>
            </div>
            
            <!-- 选中标识 -->
            <div v-if="themeStore.currentThemeId === theme.id" class="theme-check">
              <PhCheck :size="16" weight="bold" />
            </div>
            
            <!-- 模式标签 -->
            <span class="theme-mode" :class="theme.mode">
              <PhMoon v-if="theme.mode === 'dark'" :size="12" />
              <PhSun v-else :size="12" />
              {{ theme.mode === 'dark' ? t('settings.darkMode') : t('settings.lightMode') }}
            </span>
          </button>
        </div>
      </section>

      <!-- 显示设置 -->
      <section class="settings-section">
        <div class="section-header">
          <PhPalette :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.displaySection') }}</h2>
        </div>
        
        <div class="glass-card settings-card">
          <!-- 语言 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.language') }}</span>
              <span class="setting-desc">{{ t('settings.languageDesc') }}</span>
            </div>
            <select 
              :value="themeStore.currentLanguageCode"
              @change="themeStore.setLanguage($event.target.value)"
              class="input select setting-select"
            >
              <option 
                v-for="lang in languageOptions" 
                :key="lang.value" 
                :value="lang.value"
              >
                {{ lang.label }}
              </option>
            </select>
          </div>

          <div class="setting-divider" />

          <!-- 计价货币 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.currency') }}</span>
              <span class="setting-desc">{{ t('settings.currencyDesc') }}</span>
            </div>
            <div class="currency-selector">
              <button
                v-for="curr in currencyOptions"
                :key="curr.value"
                class="currency-btn"
                :class="{ active: settings.currency === curr.value }"
                @click="settings.currency = curr.value"
              >
                <span class="currency-symbol">{{ curr.symbol }}</span>
                <span class="currency-code">{{ curr.value }}</span>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 数据设置 -->
      <section class="settings-section">
        <div class="section-header">
          <PhDatabase :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.dataSection') }}</h2>
        </div>
        
        <div class="glass-card settings-card">
          <!-- 自动刷新 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.autoRefresh') }}</span>
              <span class="setting-desc">{{ t('settings.autoRefreshDesc') }}</span>
            </div>
            <label class="toggle">
              <input 
                type="checkbox" 
                v-model="settings.autoRefresh"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 刷新间隔 -->
          <div class="setting-item" :class="{ disabled: !settings.autoRefresh }">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.refreshInterval') }}</span>
              <span class="setting-desc">{{ t('settings.refreshIntervalDesc') }}</span>
            </div>
            <select 
              v-model="settings.refreshInterval"
              class="input select setting-select"
              :disabled="!settings.autoRefresh"
            >
              <option 
                v-for="interval in refreshIntervalOptions" 
                :key="interval.value" 
                :value="interval.value"
              >
                {{ interval.label }}
              </option>
            </select>
          </div>

          <div class="setting-divider" />

          <!-- 历史数据保留 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.historyRetention') }}</span>
              <span class="setting-desc">{{ t('settings.historyRetentionDesc') }}</span>
            </div>
            <select 
              v-model="settings.historyRetention"
              class="input select setting-select"
            >
              <option 
                v-for="ret in retentionOptions" 
                :key="ret.value" 
                :value="ret.value"
              >
                {{ ret.label }}
              </option>
            </select>
          </div>
        </div>
      </section>

      <!-- 交易同步设置 -->
      <section class="settings-section">
        <div class="section-header">
          <PhClockClockwise :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.txSyncSection') }}</h2>
        </div>

        <div class="glass-card settings-card">
          <!-- 同步开关 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.txSyncEnabled') }}</span>
              <span class="setting-desc">{{ t('settings.txSyncEnabledDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="txSyncSettings.enabled"
                @change="updateTxSyncSetting('enabled', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <!-- 开启同步后的警告提示 -->
          <div v-if="txSyncSettings.enabled" class="sync-warning">
            <PhWarning :size="16" />
            <span>{{ t('settings.txSyncWarning') }}</span>
          </div>

          <div class="setting-divider" />

          <!-- 同步间隔 -->
          <div class="setting-item" :class="{ disabled: !txSyncSettings.enabled }">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.txSyncInterval') }}</span>
              <span class="setting-desc">{{ t('settings.txSyncIntervalDesc') }}</span>
            </div>
            <select
              class="input select setting-select"
              :value="txSyncSettings.interval_minutes"
              :disabled="!txSyncSettings.enabled"
              @change="updateTxSyncSetting('interval_minutes', Number($event.target.value))"
            >
              <option
                v-for="opt in txSyncIntervalOptions"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
          </div>

          <div class="setting-divider" />

          <!-- 回溯天数 -->
          <div class="setting-item" :class="{ disabled: !txSyncSettings.enabled }">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.txSyncLookback') }}</span>
              <span class="setting-desc">{{ t('settings.txSyncLookbackDesc') }}</span>
            </div>
            <select
              class="input select setting-select"
              :value="txSyncSettings.lookback_days"
              :disabled="!txSyncSettings.enabled"
              @change="updateTxSyncSetting('lookback_days', Number($event.target.value))"
            >
              <option
                v-for="opt in txSyncLookbackOptions"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
          </div>

          <div class="setting-divider" />

          <!-- 手动同步按钮 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.txSyncNow') }}</span>
              <span class="setting-desc">
                {{ txSyncSettings.enabled ? t('settings.txSyncEnabledDesc') : t('settings.txSyncDisabledHint') }}
              </span>
            </div>
            <button
              class="btn btn-primary btn-sm"
              :disabled="!txSyncSettings.enabled || txSyncing"
              @click="triggerSync"
            >
              <PhArrowsClockwise v-if="txSyncing" :size="16" class="spin" />
              <span>{{ txSyncing ? t('settings.txSyncSyncing') : t('settings.txSyncNow') }}</span>
            </button>
          </div>
        </div>
      </section>

      <!-- 通知设置 -->
      <section class="settings-section">
        <div class="section-header">
          <PhBell :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.notificationSection') }}</h2>
        </div>

        <div class="glass-card settings-card">
          <!-- 推送通知 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.pushEnabled') }}</span>
              <span class="setting-desc">{{ t('settings.pushEnabledDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="notifStore.preferences?.push_enabled"
                @change="updateNotifPreference('push_enabled', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 邮件通知 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.emailEnabled') }}</span>
              <span class="setting-desc">{{ t('settings.emailEnabledDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="notifStore.preferences?.email_enabled"
                @change="updateNotifPreference('email_enabled', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 价格预警 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.priceAlert') }}</span>
              <span class="setting-desc">{{ t('settings.priceAlertDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="notifStore.preferences?.price_alert"
                @change="updateNotifPreference('price_alert', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 资产变动通知 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.assetAlert') }}</span>
              <span class="setting-desc">{{ t('settings.assetAlertDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="notifStore.preferences?.portfolio_alert"
                @change="updateNotifPreference('portfolio_alert', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 系统通知 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.systemNotice') }}</span>
              <span class="setting-desc">{{ t('settings.systemNoticeDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="notifStore.preferences?.system_notice"
                @change="updateNotifPreference('system_notice', $event.target.checked)"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
      </section>

      <!-- 安全设置 -->
      <section class="settings-section">
        <div class="section-header">
          <PhShieldCheck :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.securitySection') }}</h2>
        </div>
        
        <div class="glass-card settings-card">
          <!-- 两步验证 (2FA) -->
          <div class="setting-item two-fa-item">
            <div class="setting-info">
              <span class="setting-label">
                <PhLock :size="18" class="inline-icon" />
                {{ t('settings.twoFA.title') }}
              </span>
              <span class="setting-desc">{{ t('settings.twoFA.description') }}</span>
            </div>
            <div class="two-fa-status">
              <span 
                class="status-badge"
                :class="is2FAEnabled ? 'enabled' : 'disabled'"
              >
                {{ is2FAEnabled ? t('settings.twoFA.enabled') : t('settings.twoFA.disabled') }}
              </span>
              <button 
                v-if="is2FAEnabled"
                class="btn btn-danger btn-sm"
                @click="openDisable2FA"
              >
                {{ t('settings.twoFA.disable') }}
              </button>
              <button 
                v-else
                class="btn btn-primary btn-sm"
                @click="startSetup2FA"
              >
                {{ t('settings.twoFA.enable') }}
              </button>
            </div>
          </div>

          <div class="setting-divider" />

          <!-- 显示余额（隐私模式） -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.showBalance') }}</span>
              <span class="setting-desc">{{ t('settings.showBalanceDesc') }}</span>
            </div>
            <label class="toggle">
              <input
                type="checkbox"
                :checked="!themeStore.privacyMode"
                @change="themeStore.togglePrivacyMode()"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-divider" />

          <!-- 操作确认 -->
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">{{ t('settings.confirmSensitive') }}</span>
              <span class="setting-desc">{{ t('settings.confirmSensitiveDesc') }}</span>
            </div>
            <label class="toggle">
              <input 
                type="checkbox" 
                v-model="settings.confirmOperations"
              >
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
      </section>

      <!-- 2FA 设置模态框 -->
      <Teleport to="body">
        <Transition name="modal">
          <div v-if="show2FASetupModal" class="modal-overlay" @click.self="show2FASetupModal = false">
            <div class="modal-content glass-card two-fa-modal">
              <!-- 步骤1: 显示二维码和密钥 -->
              <template v-if="setupStep === 1">
                <div class="modal-header">
                  <PhQrCode :size="32" weight="duotone" class="modal-icon" />
                  <h3>{{ t('settings.twoFA.setupTitle') }}</h3>
                  <p class="modal-subtitle">{{ t('settings.twoFA.setupStep1') }}</p>
                </div>
                
                <div class="modal-body">
                  <!-- QR Code 占位 -->
                  <div class="qr-code-container">
                    <div class="qr-placeholder">
                      <PhQrCode :size="80" />
                      <span>{{ t('settings.twoFA.scanQR') }}</span>
                    </div>
                  </div>
                  
                  <!-- 密钥显示 -->
                  <div class="secret-key">
                    <span class="secret-label">{{ t('settings.twoFA.secretKey') }}</span>
                    <div class="secret-value">
                      <code>{{ twoFASecret }}</code>
                      <button class="copy-btn" @click="copySecret">
                        <PhCheck v-if="codeCopied" :size="16" />
                        <PhCopy v-else :size="16" />
                      </button>
                    </div>
                    <span class="secret-hint">{{ t('settings.twoFA.manualEntry') }}</span>
                  </div>
                </div>
                
                <div class="modal-actions">
                  <button class="btn btn-secondary" @click="show2FASetupModal = false">
                    {{ t('common.cancel') }}
                  </button>
                  <button class="btn btn-primary" @click="setupStep = 2">
                    {{ t('settings.twoFA.next') }}
                  </button>
                </div>
              </template>
              
              <!-- 步骤2: 验证代码 -->
              <template v-else>
                <div class="modal-header">
                  <PhKey :size="32" weight="duotone" class="modal-icon" />
                  <h3>{{ t('settings.twoFA.verifyTitle') }}</h3>
                  <p class="modal-subtitle">{{ t('settings.twoFA.setupStep2') }}</p>
                </div>
                
                <div class="modal-body">
                  <div v-if="verifyError" class="error-alert">
                    <PhWarning :size="18" />
                    <span>{{ verifyError }}</span>
                  </div>
                  
                  <div class="verify-input-group">
                    <label>{{ t('settings.twoFA.enterCode') }}</label>
                    <input
                      type="text"
                      v-model="verifyCode"
                      class="verify-input"
                      maxlength="6"
                      placeholder="000000"
                      inputmode="numeric"
                    />
                  </div>
                </div>
                
                <div class="modal-actions">
                  <button class="btn btn-secondary" @click="setupStep = 1">
                    {{ t('settings.twoFA.back') }}
                  </button>
                  <button 
                    class="btn btn-primary" 
                    @click="verifyAndEnable2FA"
                    :disabled="isVerifying || verifyCode.length !== 6"
                  >
                    <PhArrowsClockwise v-if="isVerifying" :size="18" class="spin" />
                    <span v-else>{{ t('settings.twoFA.verify') }}</span>
                  </button>
                </div>
              </template>
            </div>
          </div>
        </Transition>
      </Teleport>

      <!-- 2FA 禁用确认模态框 -->
      <Teleport to="body">
        <Transition name="modal">
          <div v-if="show2FADisableModal" class="modal-overlay" @click.self="show2FADisableModal = false">
            <div class="modal-content glass-card two-fa-modal">
              <div class="modal-header warning">
                <PhWarning :size="32" weight="fill" class="modal-icon warning" />
                <h3>{{ t('settings.twoFA.disableTitle') }}</h3>
                <p class="modal-subtitle">{{ t('settings.twoFA.disableWarning') }}</p>
              </div>
              
              <div class="modal-body">
                <div v-if="disableError" class="error-alert">
                  <PhWarning :size="18" />
                  <span>{{ disableError }}</span>
                </div>
                
                <div class="verify-input-group">
                  <label>{{ t('settings.twoFA.enterCodeToDisable') }}</label>
                  <input
                    type="text"
                    v-model="disableCode"
                    class="verify-input"
                    maxlength="6"
                    placeholder="000000"
                    inputmode="numeric"
                  />
                </div>
              </div>
              
              <div class="modal-actions">
                <button class="btn btn-secondary" @click="show2FADisableModal = false">
                  {{ t('common.cancel') }}
                </button>
                <button 
                  class="btn btn-danger" 
                  @click="disable2FA"
                  :disabled="isDisabling || disableCode.length !== 6"
                >
                  <PhArrowsClockwise v-if="isDisabling" :size="18" class="spin" />
                  <span v-else>{{ t('settings.twoFA.confirmDisable') }}</span>
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>

      <!-- 导出格式选择弹窗 -->
      <Teleport to="body">
        <Transition name="modal">
          <div v-if="showExportDialog" class="modal-overlay" @click.self="showExportDialog = false">
            <div class="modal-content glass-card export-modal">
              <div class="modal-header">
                <PhDownloadSimple :size="28" weight="duotone" class="modal-icon" />
                <h3>{{ t('export.title') }}</h3>
                <p class="modal-subtitle">{{ t('export.subtitle') }}</p>
              </div>

              <div class="modal-body">
                <div class="export-options">
                  <button class="export-option-btn" @click="doExport('csv')">
                    <div class="export-option-icon csv">CSV</div>
                    <div class="export-option-info">
                      <span class="export-option-title">{{ t('export.csvTitle') }}</span>
                      <span class="export-option-desc">{{ t('export.csvDesc') }}</span>
                    </div>
                  </button>
                  <button class="export-option-btn" @click="doExport('json')">
                    <div class="export-option-icon json">JSON</div>
                    <div class="export-option-info">
                      <span class="export-option-title">{{ t('export.jsonTitle') }}</span>
                      <span class="export-option-desc">{{ t('export.jsonDesc') }}</span>
                    </div>
                  </button>
                </div>
              </div>

              <div class="modal-actions">
                <button class="btn btn-secondary" @click="showExportDialog = false">
                  {{ t('common.cancel') }}
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </Teleport>

      <!-- 成就系统 -->
      <section class="settings-section full-width">
        <div class="section-header">
          <PhTrophy :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('achievement.title') }}</h2>
        </div>
        <div class="glass-card settings-card">
          <div class="achievement-entry" @click="showAchievements = true">
            <div class="ach-progress">
              <span class="ach-count font-mono">{{ achStore.unlockedCount }}</span>
              <span class="ach-total"> / {{ achStore.totalCount }}</span>
              <span class="ach-label">{{ t('achievement.unlocked') }}</span>
            </div>
            <PhCaretRight :size="16" class="action-arrow" />
          </div>
        </div>
      </section>

      <!-- 成就面板弹窗 -->
      <AchievementPanel
        :visible="showAchievements"
        @close="showAchievements = false"
      />

      <!-- 关于 AllFi -->
      <section class="settings-section full-width">
        <div class="section-header">
          <PhInfo :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('system.aboutTitle') }}</h2>
        </div>

        <div class="glass-card settings-card">
          <!-- 版本信息区 -->
          <div class="about-version-block">
            <div class="version-current">
              <span class="version-label">{{ t('system.currentVersion') }}</span>
              <span class="version-number font-mono">v{{ systemStore.versionInfo?.version || __APP_VERSION__ }}</span>
            </div>

            <!-- 新版本提示 -->
            <div v-if="systemStore.hasUpdate" class="update-banner">
              <div class="update-badge">
                <PhArrowDown :size="18" />
                <div>
                  <span class="update-text">{{ t('system.newVersionAvailable') }}</span>
                  <span class="update-version font-mono">v{{ systemStore.updateInfo?.latest_version }}</span>
                </div>
              </div>
              <button
                class="btn btn-primary"
                :disabled="systemStore.isUpdating"
                @click="systemStore.applyUpdate(systemStore.updateInfo.latest_version)"
              >
                <PhArrowsClockwise v-if="systemStore.isUpdating" :size="18" class="spin" />
                <PhArrowDown v-else :size="18" />
                <span>{{ systemStore.isUpdating ? t('system.updating') : t('system.updateNow') }}</span>
              </button>
            </div>

            <!-- 已是最新版本 -->
            <div v-else-if="systemStore.updateInfo && !systemStore.hasUpdate" class="up-to-date-badge">
              <PhCheck :size="16" />
              <span>{{ t('system.upToDate') }}</span>
            </div>

            <!-- 操作按钮 -->
            <div class="about-actions">
              <button
                class="btn btn-secondary btn-sm"
                :disabled="systemStore.isChecking"
                @click="systemStore.checkForUpdate()"
              >
                <PhArrowsClockwise v-if="systemStore.isChecking" :size="16" class="spin" />
                <span>{{ systemStore.isChecking ? t('system.checking') : t('system.checkUpdate') }}</span>
              </button>
              <a
                v-if="systemStore.updateInfo?.release_url"
                :href="systemStore.updateInfo.release_url"
                target="_blank"
                rel="noopener"
                class="changelog-link"
              >
                {{ t('system.viewChangelog') }} ↗
              </a>
            </div>
          </div>

          <!-- 更新进度（更新中时显示） -->
          <template v-if="systemStore.isUpdating">
            <div class="setting-divider" />
            <div class="update-progress-block">
              <span class="setting-label">{{ t('system.updateProgress') }}</span>
              <div class="progress-bar-container">
                <div
                  class="progress-bar-fill"
                  :style="{ width: (systemStore.updateStatus?.total > 0 ? (systemStore.updateStatus.step / systemStore.updateStatus.total * 100) : 0) + '%' }"
                ></div>
              </div>
              <div class="progress-info">
                <span>{{ systemStore.updateStatus?.step_name }}</span>
                <span class="font-mono">{{ systemStore.updateStatus?.step }}/{{ systemStore.updateStatus?.total }}</span>
              </div>
            </div>
          </template>

          <!-- 版本回滚 -->
          <template v-if="systemStore.rollbackTargets.length > 0">
            <div class="setting-divider" />
            <div class="rollback-block">
              <span class="setting-label">{{ t('system.rollbackTitle') }}</span>
              <div
                v-for="record in systemStore.rollbackTargets"
                :key="record.version"
                class="rollback-item"
              >
                <div class="rollback-info">
                  <span class="font-mono">v{{ record.version }}</span>
                  <span class="rollback-date">{{ record.timestamp }}</span>
                </div>
                <button
                  class="btn btn-secondary btn-sm"
                  :disabled="systemStore.isUpdating"
                  @click="handleRollback(record.version)"
                >
                  <PhArrowCounterClockwise :size="14" />
                  <span>{{ t('system.rollbackTo') }}</span>
                </button>
              </div>
            </div>
          </template>

          <div class="setting-divider" />

          <!-- 系统信息 -->
          <div class="system-info-block">
            <div class="info-row">
              <span class="info-label">{{ t('system.runMode') }}</span>
              <span class="info-value">
                {{ systemStore.versionInfo?.run_mode === 'docker' ? t('system.runModeDocker') : t('system.runModeHost') }}
              </span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('system.buildTime') }}</span>
              <span class="info-value font-mono">{{ systemStore.versionInfo?.build_time || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('system.gitCommit') }}</span>
              <span class="info-value font-mono">{{ systemStore.versionInfo?.git_commit || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ t('system.goVersion') }}</span>
              <span class="info-value font-mono">{{ systemStore.versionInfo?.go_version || '-' }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 数据管理 -->
      <section class="settings-section full-width">
        <div class="section-header">
          <PhDatabase :size="20" weight="duotone" />
          <h2 class="section-title">{{ t('settings.dataManageSection') }}</h2>
        </div>
        
        <div class="glass-card settings-card">
          <div class="data-actions">
            <button class="data-action-btn" @click="openExportDialog">
              <div class="action-icon export">
                <PhDownloadSimple :size="24" />
              </div>
              <div class="action-info">
                <span class="action-title">{{ t('settings.exportData') }}</span>
                <span class="action-desc">{{ t('settings.exportDataDesc') }}</span>
              </div>
              <PhCaretRight :size="20" class="action-arrow" />
            </button>
            
            <button class="data-action-btn" @click="clearCache">
              <div class="action-icon cache">
                <PhArrowsClockwise :size="24" />
              </div>
              <div class="action-info">
                <span class="action-title">{{ t('settings.clearCache') }}</span>
                <span class="action-desc">{{ t('settings.clearCacheDesc') }}</span>
              </div>
              <PhCaretRight :size="20" class="action-arrow" />
            </button>
            
            <button class="data-action-btn danger" @click="resetSettings">
              <div class="action-icon reset">
                <PhTrash :size="24" />
              </div>
              <div class="action-info">
                <span class="action-title">{{ t('settings.resetSettings') }}</span>
                <span class="action-desc">{{ t('settings.resetSettingsDesc') }}</span>
              </div>
              <PhCaretRight :size="20" class="action-arrow" />
            </button>
          </div>
        </div>
      </section>

    </div>

  </div>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xl);
  max-width: 1400px;
  margin: 0 auto;
}

/* ================== 页面头部 ================== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--gap-md);
}

.header-subtitle {
  color: var(--color-text-muted);
  font-size: 13px;
}

/* ================== 设置网格 ================== */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--gap-xl);
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.settings-section.full-width {
  grid-column: 1 / -1;
}

.section-header {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  color: var(--color-accent-primary);
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-heading);
}

/* ================== 主题选择器 ================== */
.themes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--gap-md);
}

.theme-card {
  position: relative;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 0;
  cursor: pointer;
  transition: border-color var(--transition-fast);
  overflow: hidden;
  text-align: left;
}

.theme-card:hover {
  border-color: var(--color-border-hover);
}

.theme-card.active {
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 1px var(--color-accent-primary);
}

.theme-preview {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.theme-icon {
  font-size: 1.5rem;
  color: white;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.theme-info {
  padding: var(--gap-sm) var(--gap-md);
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.theme-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.theme-desc {
  font-size: 11px;
  color: var(--color-text-muted);
}

.theme-font {
  font-size: 10px;
  color: var(--color-accent-primary);
  font-family: var(--font-mono);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  opacity: 0.7;
}

.theme-check {
  position: absolute;
  top: var(--gap-xs);
  right: var(--gap-xs);
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-accent-primary);
  color: var(--color-text-inverse);
  border-radius: 50%;
}

.theme-mode {
  position: absolute;
  top: var(--gap-xs);
  left: var(--gap-xs);
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
}

.theme-mode.dark {
  background: rgba(0, 0, 0, 0.5);
  color: white;
}

.theme-mode.light {
  background: rgba(255, 255, 255, 0.9);
  color: #0F172A;
}

/* ================== 设置卡片 ================== */
.settings-card {
  padding: 0;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-lg);
  padding: var(--gap-md) var(--gap-lg);
  transition: opacity var(--transition-fast);
}

.setting-item.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.setting-info {
  flex: 1;
  min-width: 0;
}

.setting-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 1px;
}

.setting-desc {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
}

.setting-divider {
  height: 1px;
  background: var(--color-border);
  margin: 0 var(--gap-lg);
}

.setting-select {
  width: auto;
  min-width: 130px;
}

/* 同步警告提示 */
.sync-warning {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-lg);
  margin: 0 var(--gap-lg);
  background: color-mix(in srgb, var(--color-warning) 8%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-warning) 25%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-warning);
  font-size: 12px;
  line-height: 1.4;
}

/* 子设置项（缩进） */
.setting-sub-item {
  padding-left: calc(var(--gap-lg) + var(--gap-md));
  background: color-mix(in srgb, var(--color-bg-tertiary) 30%, transparent);
}

/* 垂直布局设置项 */
.setting-item-vertical {
  flex-direction: column;
  align-items: stretch;
  gap: var(--gap-sm);
}

/* 阈值输入 */
.threshold-input {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.threshold-input .input {
  width: 72px;
  text-align: center;
  font-family: var(--font-mono);
}

.threshold-unit {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

/* Webhook 输入 */
.webhook-input {
  width: 100%;
  font-family: var(--font-mono);
  font-size: 12px;
}

/* 货币选择器 */
.currency-selector {
  display: flex;
  gap: var(--gap-xs);
}

.currency-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  background: color-mix(in srgb, var(--color-bg-tertiary) 50%, transparent);
  border: 1px solid transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.currency-btn:hover {
  border-color: var(--color-border);
  color: var(--color-text-secondary);
}

.currency-btn.active {
  background: color-mix(in srgb, var(--color-accent-primary) 15%, transparent);
  border-color: var(--color-accent-primary);
  color: var(--color-accent-primary);
}

.currency-symbol {
  font-size: 14px;
  font-weight: 600;
  font-family: var(--font-mono);
}

.currency-code {
  font-size: 11px;
  font-weight: 500;
}

/* Toggle 开关 */
.toggle {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
  cursor: pointer;
  flex-shrink: 0;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  inset: 0;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: 9999px;
  transition: all var(--transition-fast);
}

.toggle-slider::before {
  content: '';
  position: absolute;
  width: 16px;
  height: 16px;
  left: 2px;
  bottom: 2px;
  background: var(--color-text-muted);
  border-radius: 50%;
  transition: all var(--transition-fast);
}

.toggle input:checked + .toggle-slider {
  background: color-mix(in srgb, var(--color-accent-primary) 20%, transparent);
  border-color: var(--color-accent-primary);
}

.toggle input:checked + .toggle-slider::before {
  transform: translateX(18px);
  background: var(--color-accent-primary);
}

/* ================== 成就入口 ================== */
.achievement-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-md) var(--gap-lg);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.achievement-entry:hover {
  background: color-mix(in srgb, var(--color-text-muted) 8%, transparent);
}

.ach-progress {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.ach-count {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-warning);
}

.ach-total {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.ach-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-left: var(--gap-xs);
}

/* ================== 数据操作 ================== */
.data-actions {
  display: flex;
  flex-direction: column;
}

.data-action-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md) var(--gap-lg);
  background: transparent;
  border: none;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: background var(--transition-fast);
  text-align: left;
}

.data-action-btn:last-child {
  border-bottom: none;
}

.data-action-btn:hover {
  background: color-mix(in srgb, var(--color-text-muted) 8%, transparent);
}

.data-action-btn.danger:hover {
  background: color-mix(in srgb, var(--color-error) 8%, transparent);
}

.action-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

.action-icon.export {
  background: color-mix(in srgb, var(--color-accent-primary) 12%, transparent);
  color: var(--color-accent-primary);
}

.action-icon.cache {
  background: color-mix(in srgb, var(--color-accent-secondary) 12%, transparent);
  color: var(--color-accent-secondary);
}

.action-icon.reset {
  background: color-mix(in srgb, var(--color-error) 12%, transparent);
  color: var(--color-error);
}

.action-info {
  flex: 1;
}

.action-title {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 1px;
}

.action-desc {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
}

.action-arrow {
  color: var(--color-text-muted);
}


/* ================== 动画 ================== */
.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ================== 响应式 ================== */
@media (max-width: 1024px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }

  .settings-section.full-width {
    grid-column: 1;
  }

  .themes-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--gap-sm);
  }

  .currency-selector,
  .setting-select {
    width: 100%;
  }

  .currency-selector {
    flex-wrap: wrap;
  }

  .currency-btn {
    flex: 1;
    min-width: 60px;
  }

  .settings-footer {
    flex-wrap: wrap;
    gap: var(--gap-xs);
  }

  .footer-divider {
    display: none;
  }

  .themes-grid {
    grid-template-columns: 1fr;
  }
}

/* ========== 2FA 设置样式 ========== */
.two-fa-item {
  flex-wrap: wrap;
  gap: var(--gap-sm);
}

.two-fa-status {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.status-badge {
  padding: 3px 10px;
  border-radius: var(--radius-xs);
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.status-badge.enabled {
  background: color-mix(in srgb, var(--color-success) 15%, transparent);
  color: var(--color-success);
}

.status-badge.disabled {
  background: color-mix(in srgb, var(--color-text-muted) 15%, transparent);
  color: var(--color-text-muted);
}

.inline-icon {
  vertical-align: middle;
  margin-right: 4px;
}

.btn-sm {
  padding: 5px 12px;
  font-size: 12px;
}

.btn-danger {
  background: var(--color-error);
  border-color: var(--color-error);
  color: white;
}

.btn-danger:hover {
  background: color-mix(in srgb, var(--color-error) 85%, black);
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.modal-content {
  width: 100%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
  padding: var(--gap-xl);
  border-radius: var(--radius-lg);
}

.two-fa-modal .modal-header {
  text-align: center;
  margin-bottom: var(--gap-lg);
}

.two-fa-modal .modal-icon {
  color: var(--color-accent-primary);
  margin-bottom: var(--gap-sm);
}

.two-fa-modal .modal-icon.warning {
  color: var(--color-warning);
}

.two-fa-modal h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-heading);
  margin-bottom: var(--gap-xs);
}

.modal-subtitle {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.modal-body {
  margin-bottom: var(--gap-lg);
}

.qr-code-container {
  display: flex;
  justify-content: center;
  margin-bottom: var(--gap-md);
}

.qr-placeholder {
  width: 160px;
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  background: var(--color-bg-tertiary);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 12px;
}

.secret-key {
  text-align: center;
}

.secret-label {
  display: block;
  font-size: 11px;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-xs);
}

.secret-value {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  background: var(--color-bg-tertiary);
  padding: var(--gap-sm) var(--gap-md);
  border-radius: var(--radius-sm);
  margin-bottom: var(--gap-xs);
}

.secret-value code {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: 2px;
}

.copy-btn {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: 4px;
  display: flex;
  transition: color var(--transition-fast);
}

.copy-btn:hover {
  color: var(--color-accent-primary);
}

.secret-hint {
  font-size: 11px;
  color: var(--color-text-muted);
}

.error-alert {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: color-mix(in srgb, var(--color-error) 10%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-error) 30%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-error);
  font-size: 13px;
  margin-bottom: var(--gap-md);
}

.verify-input-group {
  text-align: center;
}

.verify-input-group label {
  display: block;
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-sm);
}

.verify-input {
  width: 100%;
  max-width: 180px;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 6px;
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  transition: border-color var(--transition-fast);
}

.verify-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.modal-actions {
  display: flex;
  gap: var(--gap-sm);
  justify-content: flex-end;
}

.modal-actions .btn {
  min-width: 90px;
}

/* 模态框动画 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.15s ease;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
  transition: transform 0.15s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.97);
}

/* ========== 导出弹窗 ========== */
.export-modal .modal-header {
  text-align: center;
  margin-bottom: var(--gap-lg);
}

.export-modal .modal-icon {
  color: var(--color-accent-primary);
  margin-bottom: var(--gap-sm);
}

.export-modal h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-heading);
  margin-bottom: var(--gap-xs);
}

.export-options {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.export-option-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md) var(--gap-lg);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  text-align: left;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.export-option-btn:hover {
  border-color: var(--color-accent-primary);
  background: color-mix(in srgb, var(--color-accent-primary) 5%, var(--color-bg-tertiary));
}

.export-option-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  font-size: 0.6875rem;
  font-weight: 600;
  font-family: var(--font-mono);
  flex-shrink: 0;
}

.export-option-icon.csv {
  background: color-mix(in srgb, var(--color-success) 15%, transparent);
  color: var(--color-success);
}

.export-option-icon.json {
  background: color-mix(in srgb, var(--color-accent-primary) 15%, transparent);
  color: var(--color-accent-primary);
}

.export-option-info {
  flex: 1;
}

.export-option-title {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: 1px;
}

.export-option-desc {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ================== 关于 AllFi ================== */
.about-version-block {
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-md);
}

.version-current {
  text-align: center;
}

.version-label {
  display: block;
  font-size: 12px;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-xs);
}

.version-number {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: -0.02em;
}

.update-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  width: 100%;
  padding: var(--gap-md);
  background: color-mix(in srgb, var(--color-success) 8%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-success) 25%, transparent);
  border-radius: var(--radius-md);
}

.update-badge {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  color: var(--color-success);
}

.update-text {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-success);
}

.update-version {
  font-size: 12px;
  color: var(--color-success);
  opacity: 0.8;
}

.up-to-date-badge {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-xs) var(--gap-sm);
  border-radius: var(--radius-sm);
  background: color-mix(in srgb, var(--color-success) 10%, transparent);
  color: var(--color-success);
  font-size: 12px;
  font-weight: 500;
}

.about-actions {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.changelog-link {
  font-size: 12px;
  color: var(--color-accent-primary);
  text-decoration: none;
  transition: opacity var(--transition-fast);
}

.changelog-link:hover {
  opacity: 0.8;
}

/* 更新进度 */
.update-progress-block {
  padding: var(--gap-md) var(--gap-lg);
}

.progress-bar-container {
  width: 100%;
  height: 6px;
  background: var(--color-bg-tertiary);
  border-radius: 3px;
  margin: var(--gap-sm) 0;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: var(--color-accent-primary);
  border-radius: 3px;
  transition: width 0.5s ease;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-muted);
}

/* 版本回滚 */
.rollback-block {
  padding: var(--gap-md) var(--gap-lg);
}

.rollback-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-sm) 0;
  border-bottom: 1px solid var(--color-border);
}

.rollback-item:last-child {
  border-bottom: none;
}

.rollback-info {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.rollback-date {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 系统信息 */
.system-info-block {
  padding: var(--gap-md) var(--gap-lg);
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap-xs) 0;
}

.info-row + .info-row {
  border-top: 1px solid color-mix(in srgb, var(--color-border) 50%, transparent);
}

.info-label {
  font-size: 12px;
  color: var(--color-text-muted);
}

.info-value {
  font-size: 12px;
  color: var(--color-text-secondary);
}
</style>
