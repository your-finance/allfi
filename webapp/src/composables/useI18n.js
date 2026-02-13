/**
 * 多语言 Composable
 * 提供响应式的翻译函数
 */
import { computed } from 'vue'
import { useThemeStore } from '../stores/themeStore'
import { translations, createTranslator } from '../i18n'

/**
 * 使用多语言翻译
 * @returns {{ t: function, locale: string }}
 */
export function useI18n() {
  const themeStore = useThemeStore()
  
  // 当前语言代码
  const locale = computed(() => themeStore.currentLanguageCode)
  
  // 翻译函数 - 响应式
  const t = computed(() => {
    return createTranslator(locale.value)
  })
  
  // 获取翻译
  const translate = (key, params) => {
    return t.value(key, params)
  }
  
  // 获取整个翻译对象（用于下拉选项等）
  const getTranslations = (section) => {
    const dict = translations[locale.value] || translations['zh-CN']
    return dict[section] || {}
  }
  
  return {
    t: translate,
    locale,
    getTranslations
  }
}
