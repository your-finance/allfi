/**
 * useWebPush - 组合式函数
 * 封装 WebPush 浏览器推送订阅 / 取消订阅的完整流程
 */
import { ref, onMounted } from 'vue'
import { notificationService } from '../api/index.js'

/**
 * Base64 URL 字符串转 Uint8Array
 * 用于将 VAPID 公钥转换为 applicationServerKey
 */
function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
    const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')
    const rawData = window.atob(base64)
    const outputArray = new Uint8Array(rawData.length)
    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i)
    }
    return outputArray
}

export function useWebPush() {
    const isSupported = ref(false)
    const isSubscribed = ref(false)
    const isLoading = ref(false)
    const error = ref(null)
    const permissionState = ref('default') // 'granted' | 'denied' | 'default'

    // 检测浏览器支持和当前订阅状态
    onMounted(async () => {
        isSupported.value = 'serviceWorker' in navigator && 'PushManager' in window

        if (!isSupported.value) return

        try {
            permissionState.value = Notification.permission

            // 检查现有订阅
            const registration = await navigator.serviceWorker.ready
            const subscription = await registration.pushManager.getSubscription()
            isSubscribed.value = !!subscription
        } catch (err) {
            console.warn('检查 WebPush 状态失败:', err)
        }
    })

    /**
     * 注册 Service Worker（如果尚未注册）
     */
    async function ensureServiceWorker() {
        const registrations = await navigator.serviceWorker.getRegistrations()
        if (registrations.length > 0) {
            return registrations[0]
        }
        // 注册一个最简的 SW 用于接收推送
        return navigator.serviceWorker.register('/sw.js')
    }

    /**
     * 订阅 WebPush 推送
     * 完整流程：
     * 1. 请求通知权限
     * 2. 注册/获取 ServiceWorker
     * 3. 获取 VAPID 公钥
     * 4. 创建 PushSubscription
     * 5. 发送订阅信息到后端
     */
    async function subscribePush() {
        if (!isSupported.value) {
            error.value = '当前浏览器不支持推送通知'
            return false
        }

        isLoading.value = true
        error.value = null

        try {
            // 1. 请求权限
            const permission = await Notification.requestPermission()
            permissionState.value = permission

            if (permission !== 'granted') {
                error.value = '用户拒绝了通知权限'
                return false
            }

            // 2. 获取 ServiceWorker
            const registration = await ensureServiceWorker()
            await navigator.serviceWorker.ready

            // 3. 获取 VAPID 公钥
            const { vapid_public_key } = await notificationService.getVapidKey()
            if (!vapid_public_key) {
                error.value = 'VAPID 公钥未配置'
                return false
            }

            // 4. 创建推送订阅
            const subscription = await registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: urlBase64ToUint8Array(vapid_public_key),
            })

            // 5. 发送到后端
            await notificationService.subscribePush(subscription)

            isSubscribed.value = true
            return true
        } catch (err) {
            console.error('WebPush 订阅失败:', err)
            error.value = err.message || '订阅失败'
            return false
        } finally {
            isLoading.value = false
        }
    }

    /**
     * 取消订阅 WebPush
     */
    async function unsubscribePush() {
        if (!isSupported.value) return false

        isLoading.value = true
        error.value = null

        try {
            const registration = await navigator.serviceWorker.ready
            const subscription = await registration.pushManager.getSubscription()

            if (subscription) {
                const endpoint = subscription.endpoint
                // 先在浏览器端取消
                await subscription.unsubscribe()
                // 再通知后端
                await notificationService.unsubscribePush(endpoint)
            }

            isSubscribed.value = false
            return true
        } catch (err) {
            console.error('WebPush 取消订阅失败:', err)
            error.value = err.message || '取消订阅失败'
            return false
        } finally {
            isLoading.value = false
        }
    }

    return {
        isSupported,
        isSubscribed,
        isLoading,
        error,
        permissionState,
        subscribePush,
        unsubscribePush,
    }
}
