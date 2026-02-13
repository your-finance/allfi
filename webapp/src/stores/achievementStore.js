/**
 * 成就状态管理 Store
 * 管理成就列表、解锁状态、解锁动画
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { achievementService } from '../api/achievementService.js'

export const useAchievementStore = defineStore('achievement', () => {
  // 状态
  const achievements = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  // 解锁动画状态
  const pendingUnlock = ref(null)

  // 已解锁数量
  const unlockedCount = computed(() =>
    achievements.value.filter(a => a.unlocked).length
  )

  // 总数量
  const totalCount = computed(() => achievements.value.length)

  // 按类别分组
  const byCategory = computed(() => ({
    milestone: achievements.value.filter(a => a.category === 'milestone'),
    persistence: achievements.value.filter(a => a.category === 'persistence'),
    investment: achievements.value.filter(a => a.category === 'investment'),
  }))

  /**
   * 获取成就列表
   */
  async function fetchAchievements() {
    isLoading.value = true
    error.value = null
    try {
      achievements.value = await achievementService.getAchievements()
    } catch (err) {
      error.value = err.message || '加载成就失败'
      console.error('加载成就失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 检查新成就（Mock: 模拟解锁第一个未解锁成就）
   */
  function checkAchievements() {
    const unlockedNew = achievements.value.find(a => !a.unlocked)
    if (unlockedNew) {
      unlockedNew.unlocked = true
      unlockedNew.unlockedAt = new Date().toISOString().split('T')[0]
      showUnlockAnimation(unlockedNew)
    }
  }

  /**
   * 触发解锁动画
   * @param {Object} achievement - 解锁的成就
   */
  function showUnlockAnimation(achievement) {
    pendingUnlock.value = achievement
  }

  /**
   * 关闭解锁动画
   */
  function dismissUnlock() {
    pendingUnlock.value = null
  }

  /**
   * 重置
   */
  function reset() {
    achievements.value = []
    error.value = null
    pendingUnlock.value = null
  }

  return {
    achievements,
    isLoading,
    error,
    unlockedCount,
    totalCount,
    byCategory,
    pendingUnlock,
    fetchAchievements,
    checkAchievements,
    showUnlockAnimation,
    dismissUnlock,
    reset,
  }
})
