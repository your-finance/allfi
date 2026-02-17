/**
 * 仪表盘布局配置 Store
 * 管理 Widget 的显示/隐藏偏好，持久化到 localStorage
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const STORAGE_KEY = 'allfi-dashboard-layout'

// 默认布局：全部显示
const DEFAULT_LAYOUT = {
  assetSummary: true,
  distribution: true,
  trend: true,
  healthScore: true,
  goals: true,
  holdings: true,
  defiOverview: true,
  nftOverview: true,
  feeAnalytics: true,
}

export const useDashboardStore = defineStore('dashboard', () => {
  // 从 localStorage 恢复配置
  const loadFromStorage = () => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        // 合并默认值（防止新增 Widget 时缺少字段）
        return { ...DEFAULT_LAYOUT, ...parsed }
      }
    } catch (e) {
      console.warn('加载仪表盘布局配置失败:', e)
    }
    return { ...DEFAULT_LAYOUT }
  }

  const widgetConfig = ref(loadFromStorage())

  // 已启用的 Widget 数量
  const enabledCount = computed(() =>
    Object.values(widgetConfig.value).filter(Boolean).length
  )

  // 至少保留一个 Widget
  const canDisable = computed(() => enabledCount.value > 1)

  // 切换 Widget 显示状态
  const toggleWidget = (name) => {
    if (!widgetConfig.value[name] || canDisable.value) {
      widgetConfig.value[name] = !widgetConfig.value[name]
      saveToStorage()
    }
  }

  // 恢复默认布局
  const resetLayout = () => {
    widgetConfig.value = { ...DEFAULT_LAYOUT }
    saveToStorage()
  }

  // 保存到 localStorage
  const saveToStorage = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(widgetConfig.value))
    } catch (e) {
      console.warn('保存仪表盘布局配置失败:', e)
    }
  }

  return {
    widgetConfig,
    enabledCount,
    canDisable,
    toggleWidget,
    resetLayout,
  }
})
