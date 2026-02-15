/**
 * 系统管理 Store
 * 管理版本信息、更新检测、更新状态
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { systemService } from '../api/systemService.js'

export const useSystemStore = defineStore('system', () => {
  // 版本信息
  const versionInfo = ref(null)
  // 更新信息（来自 checkUpdate）
  const updateInfo = ref(null)
  // 更新状态
  const updateStatus = ref({ state: 'idle', step: 0, total: 3, step_name: '', message: '' })
  // 更新历史
  const updateHistory = ref([])
  // 加载状态
  const isChecking = ref(false)
  const isUpdating = ref(false)
  // 轮询定时器
  let statusPollTimer = null
  let autoCheckTimer = null

  // 是否有新版本
  const hasUpdate = computed(() => updateInfo.value?.has_update === true)

  // 可回滚的版本列表（排除当前版本）
  const rollbackTargets = computed(() => {
    if (!updateHistory.value?.length) return []
    const current = versionInfo.value?.version
    return updateHistory.value.filter(r => r.version !== current && r.status === 'success')
  })

  // 加载版本信息
  async function loadVersion() {
    try {
      versionInfo.value = await systemService.getVersion()
    } catch (e) {
      console.error('获取版本信息失败:', e)
    }
  }

  // 检查更新
  async function checkForUpdate() {
    isChecking.value = true
    try {
      updateInfo.value = await systemService.checkUpdate()
    } catch (e) {
      console.error('检查更新失败:', e)
    } finally {
      isChecking.value = false
    }
  }

  // 执行更新
  async function applyUpdate(targetVersion) {
    isUpdating.value = true
    try {
      await systemService.applyUpdate(targetVersion)
      startStatusPolling()
    } catch (e) {
      console.error('执行更新失败:', e)
      isUpdating.value = false
    }
  }

  // 执行回滚
  async function rollback(targetVersion) {
    isUpdating.value = true
    try {
      await systemService.rollback(targetVersion)
      startStatusPolling()
    } catch (e) {
      console.error('执行回滚失败:', e)
      isUpdating.value = false
    }
  }

  // 轮询更新状态
  function startStatusPolling() {
    stopStatusPolling()
    statusPollTimer = setInterval(async () => {
      try {
        const status = await systemService.getUpdateStatus()
        updateStatus.value = status
        if (status.state === 'completed' || status.state === 'failed') {
          stopStatusPolling()
          isUpdating.value = false
          if (status.state === 'completed') {
            await loadUpdateHistory()
          }
        }
      } catch (e) {
        // 后端可能正在重启，继续轮询
      }
    }, 2000)
  }

  function stopStatusPolling() {
    if (statusPollTimer) {
      clearInterval(statusPollTimer)
      statusPollTimer = null
    }
  }

  // 加载更新历史
  async function loadUpdateHistory() {
    try {
      const res = await systemService.getUpdateHistory()
      updateHistory.value = res.records || []
    } catch (e) {
      console.error('获取更新历史失败:', e)
    }
  }

  // 启动定时检查（每小时）
  function startAutoCheck() {
    stopAutoCheck()
    autoCheckTimer = setInterval(checkForUpdate, 3600 * 1000)
  }

  function stopAutoCheck() {
    if (autoCheckTimer) {
      clearInterval(autoCheckTimer)
      autoCheckTimer = null
    }
  }

  return {
    versionInfo,
    updateInfo,
    updateStatus,
    updateHistory,
    isChecking,
    isUpdating,
    hasUpdate,
    rollbackTargets,
    loadVersion,
    checkForUpdate,
    applyUpdate,
    rollback,
    loadUpdateHistory,
    startAutoCheck,
    stopAutoCheck,
  }
})
