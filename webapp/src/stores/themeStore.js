/**
 * 主题和语言状态管理 Store
 * 使用 Pinia 管理全局主题、语言状态
 */
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { themes, defaultTheme, themeList, applyTheme } from '../themes'

// 本地存储key
const THEME_STORAGE_KEY = 'allfi-theme'
const LANGUAGE_STORAGE_KEY = 'allfi-language'
const PRIVACY_STORAGE_KEY = 'allfi-privacy-mode'
const ONBOARDING_STORAGE_KEY = 'allfi-onboarding-completed'

// 支持的语言列表（3国语言）
export const languages = [
  { code: 'zh-CN', name: '简体中文', flag: '🇨🇳', short: '中' },
  { code: 'zh-TW', name: '繁體中文', flag: '🇹🇼', short: '繁' },
  { code: 'en-US', name: 'English', flag: '🇺🇸', short: 'EN' }
]

// 默认语言
const defaultLanguage = 'zh-CN'

export const useThemeStore = defineStore('theme', () => {
  // ==================== 主题状态 ====================
  // 当前主题ID
  const currentThemeId = ref(defaultTheme)
  
  // 当前主题配置
  const currentTheme = computed(() => {
    return themes[currentThemeId.value] || themes[defaultTheme]
  })
  
  // 是否为深色模式
  const isDarkMode = computed(() => {
    return currentTheme.value.mode === 'dark'
  })
  
  // 所有主题列表
  const availableThemes = computed(() => themeList)
  
  // ==================== 隐私模式 ====================
  const privacyMode = ref(false)

  // ==================== 引导状态 ====================
  const onboardingCompleted = ref(false)

  // ==================== 语言状态 ====================
  // 当前语言
  const currentLanguageCode = ref(defaultLanguage)
  
  // 当前语言配置
  const currentLanguage = computed(() => {
    return languages.find(l => l.code === currentLanguageCode.value) || languages[0]
  })
  
  // 所有语言列表
  const availableLanguages = computed(() => languages)

  /**
   * 初始化主题
   * 从本地存储读取或使用默认主题
   */
  function initTheme() {
    // 初始化主题
    const savedTheme = localStorage.getItem(THEME_STORAGE_KEY)
    if (savedTheme && themes[savedTheme]) {
      currentThemeId.value = savedTheme
    }
    applyTheme(currentThemeId.value)
    
    // 初始化语言
    const savedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY)
    if (savedLanguage && languages.find(l => l.code === savedLanguage)) {
      currentLanguageCode.value = savedLanguage
    }
    // 设置 HTML lang 属性
    document.documentElement.lang = currentLanguageCode.value

    // 初始化隐私模式
    const savedPrivacy = localStorage.getItem(PRIVACY_STORAGE_KEY)
    if (savedPrivacy === 'true') {
      privacyMode.value = true
    }

    // 初始化引导状态
    const savedOnboarding = localStorage.getItem(ONBOARDING_STORAGE_KEY)
    if (savedOnboarding === 'true') {
      onboardingCompleted.value = true
    }
  }
  
  /**
   * 切换主题
   * @param {string} themeId - 主题ID
   */
  function setTheme(themeId) {
    if (!themes[themeId]) {
      console.warn(`主题 "${themeId}" 不存在`)
      return false
    }
    
    currentThemeId.value = themeId
    localStorage.setItem(THEME_STORAGE_KEY, themeId)
    applyTheme(themeId)
    
    return true
  }
  
  /**
   * 切换语言
   * @param {string} langCode - 语言代码
   */
  function setLanguage(langCode) {
    const lang = languages.find(l => l.code === langCode)
    if (!lang) {
      console.warn(`语言 "${langCode}" 不支持`)
      return false
    }
    
    currentLanguageCode.value = langCode
    localStorage.setItem(LANGUAGE_STORAGE_KEY, langCode)
    document.documentElement.lang = langCode
    
    return true
  }
  
  /**
   * 切换到下一个主题（循环）
   */
  function nextTheme() {
    const themeIds = Object.keys(themes)
    const currentIndex = themeIds.indexOf(currentThemeId.value)
    const nextIndex = (currentIndex + 1) % themeIds.length
    setTheme(themeIds[nextIndex])
  }
  
  /**
   * 切换深色/浅色模式
   * 在当前模式下选择对应的主题
   */
  function toggleMode() {
    // 找到另一个模式的第一个主题
    const targetMode = isDarkMode.value ? 'light' : 'dark'
    const targetTheme = themeList.find(t => t.mode === targetMode)
    if (targetTheme) {
      setTheme(targetTheme.id)
    }
  }
  
  /**
   * 切换隐私模式
   */
  function togglePrivacyMode() {
    privacyMode.value = !privacyMode.value
    localStorage.setItem(PRIVACY_STORAGE_KEY, String(privacyMode.value))
  }

  /**
   * 完成引导
   */
  function completeOnboarding() {
    onboardingCompleted.value = true
    localStorage.setItem(ONBOARDING_STORAGE_KEY, 'true')
  }

  /**
   * 重置为默认主题
   */
  function resetTheme() {
    setTheme(defaultTheme)
    localStorage.removeItem(THEME_STORAGE_KEY)
  }
  
  return {
    // Theme State
    currentThemeId,
    currentTheme,
    isDarkMode,
    availableThemes,
    
    // Privacy Mode
    privacyMode,
    togglePrivacyMode,

    // Onboarding
    onboardingCompleted,
    completeOnboarding,

    // Language State
    currentLanguageCode,
    currentLanguage,
    availableLanguages,
    
    // Theme Actions
    initTheme,
    setTheme,
    nextTheme,
    toggleMode,
    resetTheme,
    
    // Language Actions
    setLanguage
  }
})
