/**
 * 策略状态管理 Store
 * 管理再平衡、定投、止盈止损策略
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { strategyService } from '../api/strategyService.js'

export const useStrategyStore = defineStore('strategy', () => {
  // 状态
  const strategies = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // 活跃策略
  const activeStrategies = computed(() =>
    strategies.value.filter(s => s.status === 'active')
  )

  // 按类型分组
  const byType = computed(() => ({
    rebalance: strategies.value.filter(s => s.type === 'rebalance'),
    dca: strategies.value.filter(s => s.type === 'dca'),
    alert: strategies.value.filter(s => s.type === 'alert'),
  }))

  /**
   * 获取策略列表
   */
  async function fetchStrategies() {
    isLoading.value = true
    error.value = null
    try {
      strategies.value = await strategyService.getStrategies()
    } catch (err) {
      error.value = err.message || '加载策略失败'
      console.error('加载策略失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 添加策略（Mock 模式在前端添加）
   */
  function addStrategy(strategy) {
    const newStrategy = {
      ...strategy,
      id: `stg-${Date.now()}`,
      status: 'active',
      createdAt: new Date().toISOString().split('T')[0],
      lastTriggeredAt: null,
    }
    strategies.value.unshift(newStrategy)
  }

  /**
   * 切换策略状态
   */
  function toggleStrategy(id) {
    const s = strategies.value.find(s => s.id === id)
    if (s) {
      s.status = s.status === 'active' ? 'paused' : 'active'
    }
  }

  /**
   * 删除策略
   */
  function deleteStrategy(id) {
    strategies.value = strategies.value.filter(s => s.id !== id)
  }

  /**
   * 获取再平衡建议
   */
  async function getRebalanceSuggestion(strategyId) {
    return strategyService.getRebalanceSuggestion(strategyId)
  }

  /**
   * 重置
   */
  function reset() {
    strategies.value = []
    error.value = null
  }

  return {
    strategies,
    isLoading,
    error,
    activeStrategies,
    byType,
    fetchStrategies,
    addStrategy,
    toggleStrategy,
    deleteStrategy,
    getRebalanceSuggestion,
    reset,
  }
})
