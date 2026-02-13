/**
 * 目标追踪 Store
 * 管理用户设定的投资目标，通过后端 API 持久化
 * 支持资产目标、持仓目标和收益目标
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAssetStore } from './assetStore'
import { goalService } from '../api/index.js'

export const useGoalStore = defineStore('goal', () => {
  // 目标列表
  const goals = ref([])
  const loading = ref(false)
  const error = ref(null)

  // 从后端加载目标
  async function loadGoals() {
    loading.value = true
    error.value = null
    try {
      goals.value = await goalService.getGoals()
    } catch (err) {
      error.value = err.message || '加载目标失败'
      console.error('加载目标失败:', err)
      // 降级：尝试从 localStorage 读取旧数据
      try {
        const saved = localStorage.getItem('allfi-goals')
        if (saved) goals.value = JSON.parse(saved)
      } catch { /* 忽略 */ }
    } finally {
      loading.value = false
    }
  }

  // 添加目标
  async function addGoal(goal) {
    try {
      const created = await goalService.createGoal(goal)
      goals.value.push(created)
    } catch (err) {
      console.error('创建目标失败:', err)
      // 降级到本地
      goals.value.push({
        id: String(Date.now()),
        ...goal,
        createdAt: new Date().toISOString(),
      })
    }
  }

  // 删除目标
  async function removeGoal(id) {
    try {
      await goalService.deleteGoal(id)
    } catch (err) {
      console.error('删除目标失败:', err)
    }
    goals.value = goals.value.filter(g => g.id !== id)
  }

  // 更新目标
  async function updateGoal(id, updates) {
    try {
      const updated = await goalService.updateGoal(id, updates)
      const index = goals.value.findIndex(g => g.id === id)
      if (index !== -1) {
        goals.value[index] = { ...goals.value[index], ...updated }
      }
    } catch (err) {
      console.error('更新目标失败:', err)
      // 降级：仅本地更新
      const index = goals.value.findIndex(g => g.id === id)
      if (index !== -1) {
        goals.value[index] = { ...goals.value[index], ...updates }
      }
    }
  }

  // 带进度的目标列表（计算逻辑保持不变）
  const goalsWithProgress = computed(() => {
    const assetStore = useAssetStore()
    const totalValue = assetStore.totalValue || 0
    const change24h = assetStore.change24h || 0

    return goals.value.map(goal => {
      let currentValue = 0
      let progress = 0

      switch (goal.type) {
        case 'asset_value': {
          // 总资产达到目标值
          currentValue = totalValue
          progress = goal.targetValue > 0 ? (currentValue / goal.targetValue) * 100 : 0
          break
        }
        case 'holding_amount': {
          // 持有指定数量的某币种
          const holdings = [
            ...(assetStore.cexAccounts || []).flatMap(a => a.holdings || []),
            ...(assetStore.walletAddresses || []).flatMap(w => w.holdings || []),
          ]
          const totalBalance = holdings
            .filter(h => h.symbol === goal.currency)
            .reduce((sum, h) => sum + h.balance, 0)
          currentValue = totalBalance
          progress = goal.targetValue > 0 ? (totalBalance / goal.targetValue) * 100 : 0
          break
        }
        case 'return_rate': {
          // 收益率目标（使用 24h 变化率近似）
          currentValue = change24h * 365 / 100
          progress = goal.targetValue > 0 ? (Math.abs(currentValue) / Math.abs(goal.targetValue)) * 100 : 0
          break
        }
      }

      progress = Math.min(progress, 100)

      // 预估达成时间（基于 30 日平均增速线性外推）
      let estimatedDate = null
      if (progress > 0 && progress < 100 && goal.type !== 'return_rate') {
        const remaining = goal.targetValue - currentValue
        const dailyRate = (change24h / 100) || 0.001
        if (dailyRate > 0 && remaining > 0) {
          const daysNeeded = Math.ceil(remaining / (currentValue * dailyRate))
          if (daysNeeded > 0 && daysNeeded < 3650) {
            const est = new Date()
            est.setDate(est.getDate() + daysNeeded)
            estimatedDate = est.toISOString().slice(0, 10)
          }
        }
      }

      return {
        ...goal,
        currentValue,
        progress: parseFloat(progress.toFixed(1)),
        estimatedDate,
      }
    })
  })

  return {
    goals,
    loading,
    error,
    goalsWithProgress,
    loadGoals,
    addGoal,
    removeGoal,
    updateGoal,
  }
})
