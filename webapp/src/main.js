/**
 * AllFi 前端应用入口
 * 使用 Vue 3 Composition API + Vue Router + Pinia
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './styles/main.css'

// 创建 Vue 应用实例
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()

// 注册插件
app.use(pinia)
app.use(router)

// 挂载应用
app.mount('#app')
