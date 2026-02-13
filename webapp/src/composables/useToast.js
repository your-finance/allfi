/**
 * 全局 Toast 通知系统
 * 提供统一的操作反馈：成功 / 错误 / 信息 / 警告
 */
import { ref } from 'vue'

// 全局 Toast 列表（响应式，跨组件共享）
const toasts = ref([])

// 最多同时显示 3 条
const MAX_TOASTS = 3

let idCounter = 0

/**
 * 显示 Toast 通知
 * @param {string} message - 消息内容
 * @param {'success'|'error'|'info'|'warning'} type - 类型
 * @param {number} duration - 自动消失时间（毫秒），默认 3000
 */
function showToast(message, type = 'info', duration = 3000) {
  const id = ++idCounter
  const toast = { id, message, type, duration }

  toasts.value.push(toast)

  // 超过最大数量时移除最早的
  if (toasts.value.length > MAX_TOASTS) {
    toasts.value.shift()
  }

  // 自动消失
  if (duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, duration)
  }

  return id
}

/**
 * 移除指定 Toast
 * @param {number} id - Toast ID
 */
function removeToast(id) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

/**
 * Toast composable
 */
export function useToast() {
  return {
    toasts,
    showToast,
    removeToast,
  }
}
