/**
 * Service Worker for WebPush notifications
 * 注意：SW 必须位于根目录才能控制整个 scope
 */

/* eslint-env serviceworker */

// 监听推送事件
self.addEventListener('push', (event) => {
    let data = { title: 'AllFi', body: '您有一条新通知' }

    if (event.data) {
        try {
            data = event.data.json()
        } catch {
            data.body = event.data.text()
        }
    }

    const options = {
        body: data.body,
        icon: '/favicon.ico',
        badge: '/favicon.ico',
        vibrate: [200, 100, 200],
        data: { url: data.url || '/' },
    }

    event.waitUntil(self.registration.showNotification(data.title, options))
})

// 点击通知时打开应用
self.addEventListener('notificationclick', (event) => {
    event.notification.close()

    const url = event.notification.data?.url || '/'
    event.waitUntil(
        self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
            // 如果已有窗口则 focus
            for (const client of clientList) {
                if (client.url.includes(self.location.origin) && 'focus' in client) {
                    client.navigate(url)
                    return client.focus()
                }
            }
            // 否则打开新窗口
            return self.clients.openWindow(url)
        })
    )
})
