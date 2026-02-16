<script setup>
/**
 * 根组件 — 侧边栏 + 顶栏布局
 * 侧边栏 200px，顶栏 48px，紧凑专业金融风格
 */
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PhChartPieSlice,
  PhWallet,
  PhGear,
  PhList,
  PhX,
  PhCaretLeft,
  PhCaretRight,
  PhArrowsClockwise,
  PhBell,
  PhPaintBrush,
  PhGlobe,
  PhCaretDown,
  PhCheck,
  PhUser,
  PhSignOut,
  PhClockCounterClockwise,
  PhChartLine,
  PhNewspaper,
  PhEye,
  PhEyeSlash
} from '@phosphor-icons/vue'
import { useThemeStore } from './stores/themeStore'
import { useAuthStore } from './stores/authStore'
import { useAssetStore } from './stores/assetStore'
import { useNotificationStore } from './stores/notificationStore'
import { useCommandStore } from './stores/commandStore'
import { useI18n } from './composables/useI18n'
import { useFormatters } from './composables/useFormatters'
import NotificationPanel from './components/NotificationPanel.vue'
import CommandPalette from './components/CommandPalette.vue'
import ToastContainer from './components/ToastContainer.vue'
import BottomNav from './components/BottomNav.vue'
import VersionBadge from './components/VersionBadge.vue'
import { useSystemStore } from './stores/systemStore'
import { marketService } from './api/marketService.js'

const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()
const authStore = useAuthStore()
const assetStore = useAssetStore()
const notifStore = useNotificationStore()
const cmdStore = useCommandStore()
const systemStore = useSystemStore()
const { t } = useI18n()
const { formatCurrency, pricingCurrency, setPricingCurrency, availablePricingCurrencies } = useFormatters()

// 判断是否为认证页面
const isAuthPage = computed(() => {
  return ['/login', '/register', '/2fa'].includes(route.path)
})

// Gas 价格状态（多链）
const gasData = ref(null)
const gasChainIndex = ref(0)
let gasRefreshTimer = null

// 当前选中的链 Gas 数据
const currentGas = computed(() => {
  if (!gasData.value?.prices?.length) return null
  return gasData.value.prices[gasChainIndex.value] || gasData.value.prices[0]
})

// 加载 Gas 价格
const loadGasPrice = async () => {
  try {
    const data = await marketService.getGasPrice()
    if (data.prices) {
      gasData.value = data
    } else {
      // 兼容异常格式
      gasData.value = { prices: [{ chain: 'Ethereum', ...data }] }
    }
  } catch (err) {
    console.error('加载 Gas 价格失败:', err)
  }
}

// 点击 Gas 胶囊切换链
const cycleGasChain = () => {
  if (!gasData.value?.prices?.length) return
  gasChainIndex.value = (gasChainIndex.value + 1) % gasData.value.prices.length
}

// Gas 价格颜色（根据 level 字段或 standard 值）
const gasLevel = computed(() => {
  if (!currentGas.value) return 'unknown'
  // 优先使用后端返回的 level 字段
  if (currentGas.value.level) return currentGas.value.level
  const standard = currentGas.value.standard
  if (standard <= 10) return 'low'
  if (standard <= 30) return 'medium'
  return 'high'
})

// 下拉菜单状态
const isUserMenuOpen = ref(false)
const isCurrencyDropdownOpen = ref(false)
const isLanguageDropdownOpen = ref(false)

// 响应式检测
const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value < 768)
const isTablet = computed(() => windowWidth.value < 1024)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

// 侧边栏状态
const isSidebarCollapsed = ref(false)
const isMobileSidebarOpen = ref(false)

// 导航菜单
const navItems = computed(() => [
  { path: '/dashboard', labelKey: 'nav.dashboard', icon: PhChartPieSlice },
  { path: '/accounts', labelKey: 'nav.accounts', icon: PhWallet },
  { path: '/history', labelKey: 'nav.history', icon: PhClockCounterClockwise },
  { path: '/analytics', labelKey: 'nav.analytics', icon: PhChartLine },
  { path: '/reports', labelKey: 'nav.reports', icon: PhNewspaper },
  { path: '/settings', labelKey: 'nav.settings', icon: PhGear }
])

// Cmd+K / Ctrl+K 全局快捷键，Ctrl+H 隐私模式切换
const handleGlobalKeydown = (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    cmdStore.toggle()
  }
  if ((e.metaKey || e.ctrlKey) && e.key === 'h') {
    e.preventDefault()
    themeStore.togglePrivacyMode()
  }
}

// 初始化
onMounted(async () => {
  themeStore.initTheme()
  authStore.restoreSession()
  await assetStore.loadSummary()
  notifStore.initialize()
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleGlobalKeydown)
  window.addEventListener('resize', handleResize)

  // 加载 Gas 价格并定时刷新（每 30 秒）
  loadGasPrice()
  gasRefreshTimer = setInterval(loadGasPrice, 30000)

  // 加载版本信息并自动检查更新
  systemStore.loadVersion()
  systemStore.checkForUpdate()
  systemStore.startAutoCheck()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleGlobalKeydown)
  window.removeEventListener('resize', handleResize)
  if (gasRefreshTimer) clearInterval(gasRefreshTimer)
  systemStore.stopAutoCheck()
})

// 点击外部关闭下拉菜单
const handleClickOutside = (e) => {
  if (!e.target.closest('.language-dropdown')) isLanguageDropdownOpen.value = false
  if (!e.target.closest('.user-menu')) isUserMenuOpen.value = false
  if (!e.target.closest('.currency-dropdown')) isCurrencyDropdownOpen.value = false
}

// 操作函数
const toggleUserMenu = () => { isUserMenuOpen.value = !isUserMenuOpen.value }
const toggleCurrencyDropdown = () => { isCurrencyDropdownOpen.value = !isCurrencyDropdownOpen.value }
const toggleLanguageDropdown = () => { isLanguageDropdownOpen.value = !isLanguageDropdownOpen.value }
const toggleSidebar = () => { isSidebarCollapsed.value = !isSidebarCollapsed.value }
const toggleMobileSidebar = () => { isMobileSidebarOpen.value = !isMobileSidebarOpen.value }
const isActiveRoute = (path) => route.path === path

const navigateTo = (path) => {
  router.push(path)
  isMobileSidebarOpen.value = false
}

const selectPricingCurrency = (code) => {
  setPricingCurrency(code)
  isCurrencyDropdownOpen.value = false
}

const selectLanguage = (langCode) => {
  themeStore.setLanguage(langCode)
  isLanguageDropdownOpen.value = false
}

const quickSwitchTheme = () => {
  document.body.classList.add('theme-transition')
  themeStore.nextTheme()
  setTimeout(() => { document.body.classList.remove('theme-transition') }, 280)
}

const handleLogout = () => {
  authStore.logout()
  isUserMenuOpen.value = false
  router.push('/login')
}

// 动态标题
const updateDocumentTitle = () => {
  const titleKey = route.meta.titleKey
  document.title = titleKey ? `${t(titleKey)} - AllFi` : 'AllFi'
}
watch(() => themeStore.currentLanguageCode, updateDocumentTitle)
watch(() => route.path, updateDocumentTitle, { immediate: true })
</script>

<template>
  <!-- 认证页面 -->
  <div v-if="isAuthPage" class="auth-layout">
    <RouterView />
  </div>

  <!-- 应用页面 -->
  <div v-else class="app-container">
    <!-- 移动端遮罩 -->
    <Transition name="fade">
      <div v-if="isMobileSidebarOpen" class="mobile-overlay" @click="toggleMobileSidebar" />
    </Transition>

    <!-- 侧边栏 -->
    <aside
      class="sidebar"
      :class="{
        'sidebar-collapsed': isSidebarCollapsed,
        'sidebar-mobile-open': isMobileSidebarOpen
      }"
    >
      <!-- Logo -->
      <div class="sidebar-header">
        <div class="sidebar-header-top">
          <div class="logo" @click="navigateTo('/dashboard')">
            <div class="logo-icon">
              <svg viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="1" y="1" width="26" height="26" rx="5" stroke="var(--color-accent-primary)" stroke-width="1.5" fill="none"/>
                <path d="M8 14 L14 8 L20 14 L14 20 Z" fill="var(--color-accent-primary)" opacity="0.8"/>
              </svg>
            </div>
            <span v-if="!isSidebarCollapsed" class="logo-text">AllFi</span>
          </div>

          <button class="collapse-btn desktop-only" @click="toggleSidebar">
            <PhCaretLeft v-if="!isSidebarCollapsed" :size="14" />
            <PhCaretRight v-else :size="14" />
          </button>

          <button class="close-btn mobile-only" @click="toggleMobileSidebar">
            <PhX :size="16" />
          </button>
        </div>

        <VersionBadge :collapsed="isSidebarCollapsed" />
      </div>

      <!-- 导航 -->
      <nav class="sidebar-nav">
        <ul>
          <li v-for="item in navItems" :key="item.path">
            <button
              class="nav-item"
              :class="{ 'nav-item-active': isActiveRoute(item.path) }"
              @click="navigateTo(item.path)"
              :title="isSidebarCollapsed ? t(item.labelKey) : ''"
            >
              <component :is="item.icon" :size="18" />
              <span v-if="!isSidebarCollapsed" class="nav-label">{{ t(item.labelKey) }}</span>
              <div v-if="isActiveRoute(item.path)" class="nav-indicator" />
            </button>
          </li>
        </ul>
      </nav>

      <!-- 底部 -->
      <div class="sidebar-footer">
        <!-- 主题切换 -->
        <button
          v-if="!isSidebarCollapsed"
          class="theme-switch-btn"
          @click="quickSwitchTheme"
          :title="`当前: ${themeStore.currentTheme.name}`"
        >
          <PhPaintBrush :size="14" />
          <span>{{ themeStore.currentTheme.name }}</span>
        </button>
        <button
          v-else
          class="theme-switch-btn-icon"
          @click="quickSwitchTheme"
          :title="`当前: ${themeStore.currentTheme.name}`"
        >
          <PhPaintBrush :size="16" />
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 顶栏 -->
      <header class="top-bar">
        <button class="menu-btn mobile-only" @click="toggleMobileSidebar">
          <PhList :size="20" />
        </button>

        <h1 class="page-title desktop-only">{{ route.meta.titleKey ? t(route.meta.titleKey) : 'AllFi' }}</h1>

        <!-- 总资产 -->
        <div class="total-assets">
          <span class="assets-label">{{ t('topBar.totalAssets') }}</span>
          <span class="assets-value font-mono">{{ formatCurrency(assetStore.totalValue) }}</span>

          <div class="currency-dropdown">
            <button class="currency-btn" @click="toggleCurrencyDropdown">
              <span>{{ pricingCurrency }}</span>
              <PhCaretDown :size="12" :class="{ 'rotated': isCurrencyDropdownOpen }" />
            </button>
            <Transition name="dropdown">
              <div v-if="isCurrencyDropdownOpen" class="dropdown-menu">
                <button
                  v-for="c in availablePricingCurrencies" :key="c"
                  class="dropdown-item"
                  :class="{ 'dropdown-item-active': pricingCurrency === c }"
                  @click="selectPricingCurrency(c)"
                >
                  <span>{{ c }}</span>
                  <PhCheck v-if="pricingCurrency === c" :size="14" />
                </button>
              </div>
            </Transition>
          </div>
        </div>

        <!-- 加油站 Gas 价格胶囊（点击切换链） -->
        <div
          v-if="currentGas"
          class="gas-capsule"
          :class="`gas-${gasLevel}`"
          :title="t('topBar.gasTooltip')"
          @click="cycleGasChain"
        >
          <span class="gas-chain-tag">{{ currentGas.chain }}</span>
          <span class="gas-value font-mono">{{ Math.round(currentGas.standard) }}</span>
          <span class="gas-unit">Gwei</span>
        </div>

        <!-- 右侧操作区 -->
        <div class="top-bar-actions">
          <!-- 隐私模式（主播模式） -->
          <button
            class="action-btn privacy-btn"
            :class="{ 'privacy-active': themeStore.privacyMode }"
            :title="themeStore.privacyMode ? t('topBar.privacyOn') : t('topBar.privacyOff')"
            @click="themeStore.togglePrivacyMode()"
          >
            <PhEyeSlash v-if="themeStore.privacyMode" :size="16" weight="bold" />
            <PhEye v-else :size="16" />
          </button>

          <!-- 语言 -->
          <div class="language-dropdown">
            <button class="action-btn" @click="toggleLanguageDropdown" :title="themeStore.currentLanguage.name">
              <PhGlobe :size="16" />
              <span class="lang-code">{{ themeStore.currentLanguage.short }}</span>
            </button>
            <Transition name="dropdown">
              <div v-if="isLanguageDropdownOpen" class="dropdown-menu">
                <button
                  v-for="lang in themeStore.availableLanguages" :key="lang.code"
                  class="dropdown-item"
                  :class="{ 'dropdown-item-active': themeStore.currentLanguageCode === lang.code }"
                  @click="selectLanguage(lang.code)"
                >
                  <span class="lang-flag">{{ lang.flag }}</span>
                  <span class="lang-name">{{ lang.name }}</span>
                  <PhCheck v-if="themeStore.currentLanguageCode === lang.code" :size="14" />
                </button>
              </div>
            </Transition>
          </div>

          <!-- 刷新 -->
          <button
            class="action-btn"
            :title="t('topBar.refreshData')"
            @click="assetStore.refreshAll()"
            :disabled="assetStore.isRefreshing"
          >
            <PhArrowsClockwise :size="16" :class="{ 'animate-spin': assetStore.isRefreshing }" />
          </button>

          <!-- 通知 -->
          <button class="action-btn" :title="t('topBar.notifications')" @click="notifStore.togglePanel()">
            <PhBell :size="16" />
            <span v-if="notifStore.hasUnread" class="notification-badge">{{ notifStore.unreadBadge }}</span>
            <span v-if="notifStore.hasUnread" class="notification-dot" />
          </button>

          <!-- 用户 -->
          <div class="user-menu" v-if="authStore.isLoggedIn">
            <button class="user-btn" @click="toggleUserMenu">
              <div class="user-avatar"><PhUser :size="14" /></div>
              <span class="user-name desktop-only">{{ authStore.userName }}</span>
              <PhCaretDown :size="12" class="desktop-only" :class="{ 'rotated': isUserMenuOpen }" />
            </button>
            <Transition name="dropdown">
              <div v-if="isUserMenuOpen" class="user-dropdown">
                <div class="user-info">
                  <span class="user-email">{{ authStore.userEmail }}</span>
                </div>
                <div class="dropdown-divider" />
                <button class="dropdown-item logout" @click="handleLogout">
                  <PhSignOut :size="16" />
                  <span>{{ t('auth.logout') }}</span>
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </header>

      <!-- 离线提示条 -->
      <div v-if="assetStore.isOffline" class="offline-bar">
        {{ t('common.offlineMode') }}
      </div>

      <!-- 路由内容 -->
      <div class="page-content">
        <RouterView v-slot="{ Component }">
          <component :is="Component" />
        </RouterView>
      </div>
    </main>

    <!-- 移动端底部导航 -->
    <BottomNav v-if="isMobile" />

    <!-- 通知面板 -->
    <NotificationPanel />

    <!-- 命令面板 -->
    <CommandPalette />

    <!-- 全局 Toast 通知 -->
    <ToastContainer />
  </div>
</template>

<style scoped>
/* ================== 认证布局 ================== */
.auth-layout {
  min-height: 100vh;
}

/* ================== 容器 ================== */
.app-container {
  display: flex;
  min-height: 100vh;
}

/* ================== 侧边栏 ================== */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 200px;
  background: var(--color-bg-secondary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: width var(--transition-base);
}

.sidebar-collapsed {
  width: 56px;
}

/* 侧边栏头部 */
.sidebar-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: var(--gap-sm) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.sidebar-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  cursor: pointer;
}

.logo-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-family: var(--font-heading);
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.collapse-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.collapse-btn:hover {
  color: var(--color-text-primary);
  border-color: var(--color-border-hover);
}

/* 折叠状态 */
.sidebar-collapsed .sidebar-header {
  padding: var(--gap-sm);
  align-items: center;
}

.sidebar-collapsed .sidebar-header-top {
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.sidebar-collapsed .logo {
  justify-content: center;
}

.sidebar-collapsed .collapse-btn {
  width: 100%;
  justify-content: center;
}

.sidebar-collapsed .sidebar-nav {
  padding: var(--gap-sm);
}

.sidebar-collapsed .nav-item {
  justify-content: center;
  padding: 10px;
}

.sidebar-collapsed .sidebar-footer {
  padding: var(--gap-sm);
}

/* 导航 */
.sidebar-nav {
  flex: 1;
  padding: var(--gap-sm) var(--gap-sm);
  overflow-y: auto;
}

.sidebar-nav ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
  position: relative;
  text-align: left;
  font-size: 0.8125rem;
}

.nav-item:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.nav-item-active {
  background: rgba(75, 131, 240, 0.08);
  color: var(--color-accent-primary);
}

.nav-item-active:hover {
  background: rgba(75, 131, 240, 0.12);
}

.nav-label {
  font-weight: 500;
}

.nav-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 16px;
  background: var(--color-accent-primary);
  border-radius: 0 2px 2px 0;
}

/* 侧边栏底部 */
.sidebar-footer {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-md) var(--gap-lg);
  border-top: 1px solid var(--color-border);
}

.theme-switch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  width: 100%;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.theme-switch-btn:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

.theme-switch-btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  margin: 0 auto;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.theme-switch-btn-icon:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

/* ================== 主内容区 ================== */
.main-content {
  flex: 1;
  margin-left: 200px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left var(--transition-base);
}

.sidebar-collapsed ~ .main-content {
  margin-left: 56px;
}

/* ================== 顶栏 ================== */
.top-bar {
  position: sticky;
  top: 0;
  display: flex;
  align-items: center;
  padding: 0 var(--gap-xl);
  background: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border);
  z-index: 50;
  height: 48px;
  gap: var(--gap-lg);
}

.page-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
}

.top-bar-actions {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

/* ================== 总资产 ================== */
.total-assets {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
  font-size: 0.8125rem;
}

.assets-label {
  color: var(--color-text-muted);
  font-size: 0.75rem;
}

.assets-value {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* ================== Gas 价格胶囊 ================== */
.gas-capsule {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 12px;
  font-size: 0.6875rem;
  border: 1px solid var(--color-border);
  background: var(--color-bg-tertiary);
  white-space: nowrap;
  cursor: pointer;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.gas-chain-tag {
  font-size: 0.5625rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  line-height: 1;
}

.gas-value {
  font-weight: 600;
  color: var(--color-text-primary);
}

.gas-unit {
  color: var(--color-text-muted);
  font-size: 0.625rem;
}

/* Gas 等级颜色 */
.gas-low {
  border-color: color-mix(in srgb, var(--color-success) 40%, transparent);
}

.gas-low .gas-value {
  color: var(--color-success);
}

.gas-medium {
  border-color: color-mix(in srgb, var(--color-warning) 40%, transparent);
}

.gas-medium .gas-value {
  color: var(--color-warning);
}

.gas-high {
  border-color: color-mix(in srgb, var(--color-error) 40%, transparent);
}

.gas-high .gas-value {
  color: var(--color-error);
}

/* ================== 操作按钮 ================== */
.action-btn {
  position: relative;
  display: flex;
  align-items: center;
  gap: 4px;
  height: 32px;
  padding: 0 8px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
  font-size: 0.75rem;
}

.action-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

/* 隐私模式按钮 */
.privacy-btn.privacy-active {
  color: var(--color-accent-primary);
  background: rgba(75, 131, 240, 0.1);
}

.lang-code {
  font-weight: 500;
}

.notification-dot {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-error);
}

.notification-badge {
  position: absolute;
  top: 2px;
  right: 0;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: 8px;
  background: var(--color-error);
  color: #fff;
  font-size: 0.625rem;
  font-weight: 600;
  font-family: var(--font-mono);
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* ================== 下拉菜单 ================== */
.currency-dropdown,
.language-dropdown {
  position: relative;
}

.currency-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: border-color var(--transition-fast);
  font-size: 0.75rem;
  font-weight: 500;
}

.currency-btn:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  right: 0;
  min-width: 140px;
  padding: 4px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  z-index: 200;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: 6px 10px;
  border-radius: var(--radius-xs);
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: background var(--transition-fast);
  text-align: left;
}

.dropdown-item:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dropdown-item-active {
  color: var(--color-accent-primary);
}

.lang-flag { font-size: 0.875rem; }
.lang-name { flex: 1; }

/* ================== 用户菜单 ================== */
.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.user-btn:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

.user-avatar {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: 50%;
  color: var(--color-text-muted);
}

.user-name {
  font-size: 0.8125rem;
  font-weight: 500;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rotated {
  transform: rotate(180deg);
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  right: 0;
  min-width: 180px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  z-index: 100;
  overflow: hidden;
}

.user-info {
  padding: 8px 12px;
}

.user-email {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.dropdown-divider {
  height: 1px;
  background: var(--color-border);
}

.user-dropdown .dropdown-item.logout:hover {
  color: var(--color-error);
}

/* ================== 离线提示条 ================== */
.offline-bar {
  padding: var(--gap-xs) var(--gap-xl);
  background: var(--color-warning);
  color: #000;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
}

/* ================== 页面内容 ================== */
.page-content {
  flex: 1;
  padding: var(--gap-xl);
}

/* ================== 移动端 ================== */
.mobile-only { display: none; }
.mobile-overlay { display: none; }

@media (max-width: 1024px) {
  .desktop-only { display: none; }
  .mobile-only { display: flex; }

  .mobile-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 90;
  }

  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar-mobile-open {
    transform: translateX(0);
  }

  .main-content {
    margin-left: 0;
  }

  .sidebar-collapsed ~ .main-content {
    margin-left: 0;
  }

  .menu-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    background: transparent;
    border: 1px solid var(--color-border);
    color: var(--color-text-secondary);
    cursor: pointer;
  }

  .close-btn {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    background: var(--color-bg-tertiary);
    border: none;
    color: var(--color-text-secondary);
    cursor: pointer;
  }

  .page-content {
    padding: var(--gap-lg);
  }

  .total-assets {
    margin-left: 0;
  }
}

@media (max-width: 768px) {
  /* 移动端触摸优化：所有交互元素最小 44px */
  .nav-item {
    min-height: 44px;
    padding: 10px 12px;
  }

  .action-btn {
    min-width: 44px;
    min-height: 44px;
  }

  .currency-btn {
    min-height: 44px;
  }

  .user-btn {
    min-height: 44px;
  }

  .dropdown-item {
    min-height: 44px;
    padding: 10px;
  }

  .theme-switch-btn,
  .theme-switch-btn-icon {
    min-height: 44px;
  }

  /* 移动端顶栏紧凑 */
  .assets-label {
    display: none;
  }

  /* 底部导航栏占位（52px 导航 + safe-area） */
  .page-content {
    padding-bottom: calc(52px + env(safe-area-inset-bottom, 0px) + var(--gap-md));
  }
}

@media (max-width: 640px) {
  .top-bar {
    padding: 0 var(--gap-lg);
    gap: var(--gap-sm);
  }

  .total-assets {
    font-size: 0.75rem;
  }

  .assets-value {
    font-size: 0.8125rem;
  }

  .page-content {
    padding: var(--gap-md);
  }

  /* 小屏幕隐藏语言代码文字和 Gas 胶囊 */
  .lang-code {
    display: none;
  }

  .gas-capsule {
    display: none;
  }
}

/* ================== 过渡动画 ================== */
.fade-enter-active, .fade-leave-active {
  transition: opacity var(--transition-fast);
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.dropdown-enter-active, .dropdown-leave-active {
  transition: opacity var(--transition-fast);
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
}
</style>
