/**
 * 认证状态管理 Store
 * 单用户 PIN 认证模式
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '../api/index.js'

export const useAuthStore = defineStore('auth', () => {
  // 认证状态
  const isAuthenticated = ref(false)
  const pinSet = ref(false)
  const token = ref(null)

  // 加载状态
  const isLoading = ref(false)
  const error = ref(null)

  // 兼容旧代码的计算属性
  const isLoggedIn = computed(() => isAuthenticated.value)
  const requires2FA = ref(false)
  const is2FAVerified = ref(true)
  const user = ref(null)
  const userEmail = computed(() => '')
  const userName = computed(() => 'AllFi User')

  /**
   * 检查认证状态（是否已设置 PIN）
   * @returns {Promise<boolean>} 是否已设置 PIN
   */
  async function checkAuthStatus() {
    try {
      const result = await authService.getStatus()
      pinSet.value = result.pin_set
      return result.pin_set
    } catch {
      // 后端不可用时默认不需要认证
      pinSet.value = false
      return false
    }
  }

  /**
   * 设置 PIN（首次使用）
   * @param {string} pin - 4-8 位数字
   * @returns {Promise<boolean>}
   */
  async function setupPIN(pin) {
    isLoading.value = true
    error.value = null

    try {
      const result = await authService.setupPIN(pin)
      token.value = result.token
      isAuthenticated.value = true
      pinSet.value = true
      // 保存到 localStorage
      localStorage.setItem('allfi-auth', JSON.stringify({
        token: result.token,
        isAuthenticated: true
      }))
      return true
    } catch (err) {
      error.value = err.message || '设置 PIN 失败'
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * PIN 登录
   * @param {string} pin - PIN 码
   * @returns {Promise<boolean>}
   */
  async function login(pin) {
    isLoading.value = true
    error.value = null

    try {
      const result = await authService.login(pin)
      token.value = result.token
      isAuthenticated.value = true
      // 保存到 localStorage
      localStorage.setItem('allfi-auth', JSON.stringify({
        token: result.token,
        isAuthenticated: true
      }))
      return true
    } catch (err) {
      error.value = err.message || '登录失败'
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 修改 PIN
   * @param {string} oldPin - 旧 PIN
   * @param {string} newPin - 新 PIN
   * @returns {Promise<boolean>}
   */
  async function changePIN(oldPin, newPin) {
    isLoading.value = true
    error.value = null

    try {
      await authService.changePIN(oldPin, newPin)
      return true
    } catch (err) {
      error.value = err.message || '修改 PIN 失败'
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 登出
   */
  function logout() {
    token.value = null
    isAuthenticated.value = false
    error.value = null
    localStorage.removeItem('allfi-auth')
  }

  /**
   * 从 localStorage 恢复会话
   */
  function restoreSession() {
    const saved = localStorage.getItem('allfi-auth')
    if (saved) {
      try {
        const data = JSON.parse(saved)
        token.value = data.token
        isAuthenticated.value = data.isAuthenticated || false
      } catch {
        localStorage.removeItem('allfi-auth')
      }
    }
  }

  /**
   * 清除错误
   */
  function clearError() {
    error.value = null
  }

  // 兼容旧代码的空操作
  async function register() { return false }
  async function verify2FA() { return true }
  async function resend2FACode() { return true }
  async function enable2FA() { return false }
  async function disable2FA() { return false }

  return {
    // State
    user,
    isAuthenticated,
    requires2FA,
    is2FAVerified,
    isLoading,
    error,
    pinSet,
    token,

    // Computed
    isLoggedIn,
    userEmail,
    userName,

    // Actions
    checkAuthStatus,
    setupPIN,
    login,
    changePIN,
    register,
    verify2FA,
    resend2FACode,
    enable2FA,
    disable2FA,
    logout,
    restoreSession,
    clearError
  }
})
