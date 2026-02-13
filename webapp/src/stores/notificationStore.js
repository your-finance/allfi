/**
 * 通知状态 Store
 * 管理通知列表、未读数量和通知偏好设置
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { notificationService } from '../api/index.js'

export const useNotificationStore = defineStore('notification', () => {
  // ========== 状态 ==========

  // 通知列表
  const notifications = ref([])
  const notificationsLoading = ref(false)
  const notificationsError = ref(null)

  // 分页
  const pagination = ref({
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 0
  })

  // 未读数量
  const unreadCount = ref(0)

  // 通知偏好
  const preferences = ref(null)
  const preferencesLoading = ref(false)

  // 面板状态
  const isPanelOpen = ref(false)

  // ========== 计算属性 ==========

  // 是否有未读通知
  const hasUnread = computed(() => unreadCount.value > 0)

  // 未读数量显示文本（超过 99 显示 99+）
  const unreadBadge = computed(() => {
    if (unreadCount.value === 0) return ''
    if (unreadCount.value > 99) return '99+'
    return String(unreadCount.value)
  })

  // ========== 操作 ==========

  /**
   * 加载通知列表
   * @param {number} page - 页码
   */
  async function loadNotifications(page = 1) {
    notificationsLoading.value = true
    notificationsError.value = null

    try {
      const result = await notificationService.getNotifications(page, pagination.value.pageSize)
      notifications.value = result.list || []
      if (result.pagination) {
        pagination.value = {
          page: result.pagination.page,
          pageSize: result.pagination.page_size,
          total: result.pagination.total,
          totalPages: result.pagination.total_pages
        }
      }
    } catch (err) {
      notificationsError.value = err.message || '加载通知列表失败'
      console.error('加载通知列表失败:', err)
    } finally {
      notificationsLoading.value = false
    }
  }

  /**
   * 加载未读数量
   */
  async function loadUnreadCount() {
    try {
      const result = await notificationService.getUnreadCount()
      unreadCount.value = result.count || 0
    } catch (err) {
      console.error('获取未读数量失败:', err)
    }
  }

  /**
   * 标记通知为已读
   * @param {number} id - 通知 ID
   */
  async function markAsRead(id) {
    try {
      await notificationService.markRead(id)
      // 更新本地状态
      const notification = notifications.value.find(n => n.id === id)
      if (notification && !notification.is_read) {
        notification.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (err) {
      console.error('标记已读失败:', err)
    }
  }

  /**
   * 标记所有通知为已读
   */
  async function markAllAsRead() {
    try {
      await notificationService.markAllRead()
      notifications.value.forEach(n => { n.is_read = true })
      unreadCount.value = 0
    } catch (err) {
      console.error('标记全部已读失败:', err)
    }
  }

  /**
   * 加载通知偏好
   */
  async function loadPreferences() {
    preferencesLoading.value = true
    try {
      const result = await notificationService.getPreferences()
      // 后端返回 { preferences: {...} }，解包
      preferences.value = result.preferences || result
    } catch (err) {
      console.error('加载通知偏好失败:', err)
    } finally {
      preferencesLoading.value = false
    }
  }

  /**
   * 更新通知偏好
   * @param {object} data - 偏好数据
   */
  async function updatePreferences(data) {
    try {
      const result = await notificationService.updatePreferences(data)
      // 后端返回 { preferences: {...} }，解包
      preferences.value = result.preferences || result
    } catch (err) {
      console.error('更新通知偏好失败:', err)
      throw err
    }
  }

  /**
   * 切换面板显示
   */
  function togglePanel() {
    isPanelOpen.value = !isPanelOpen.value
    if (isPanelOpen.value) {
      loadNotifications()
    }
  }

  /**
   * 关闭面板
   */
  function closePanel() {
    isPanelOpen.value = false
  }

  /**
   * 初始化
   */
  async function initialize() {
    await loadUnreadCount()
  }

  return {
    // 状态
    notifications,
    notificationsLoading,
    notificationsError,
    pagination,
    unreadCount,
    preferences,
    preferencesLoading,
    isPanelOpen,

    // 计算属性
    hasUnread,
    unreadBadge,

    // 操作
    loadNotifications,
    loadUnreadCount,
    markAsRead,
    markAllAsRead,
    loadPreferences,
    updatePreferences,
    togglePanel,
    closePanel,
    initialize
  }
})
