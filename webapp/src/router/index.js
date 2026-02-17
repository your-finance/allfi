/**
 * 路由配置
 * 定义应用的路由结构和导航规则
 */
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

// 页面组件（懒加载）
const Login = () => import('../pages/Login.vue')
const Register = () => import('../pages/Register.vue')
const TwoFactorAuth = () => import('../pages/TwoFactorAuth.vue')
const Dashboard = () => import('../pages/Dashboard.vue')
const Accounts = () => import('../pages/Accounts.vue')
const Settings = () => import('../pages/Settings.vue')
const History = () => import('../pages/History.vue')
const Analytics = () => import('../pages/Analytics.vue')
const Reports = () => import('../pages/Reports.vue')
const DeFi = () => import('../pages/DeFi.vue')
const NFT = () => import('../pages/NFT.vue')

// 路由配置
const routes = [
  // 认证页面（无需登录）
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: {
      requiresAuth: false,
      titleKey: 'auth.login'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: {
      requiresAuth: false,
      titleKey: 'auth.register'
    }
  },
  {
    path: '/2fa',
    name: 'TwoFactorAuth',
    component: TwoFactorAuth,
    meta: {
      requiresAuth: false,
      requires2FA: true,
      titleKey: 'auth.twoFactorAuth'
    }
  },
  
  // 应用页面（需要登录）
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.dashboard'
    }
  },
  {
    path: '/history',
    name: 'History',
    component: History,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.history'
    }
  },
  {
    path: '/analytics',
    name: 'Analytics',
    component: Analytics,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.analytics'
    }
  },
  {
    path: '/accounts',
    name: 'Accounts',
    component: Accounts,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.accounts'
    }
  },
  {
    path: '/reports',
    name: 'Reports',
    component: Reports,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.reports'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.settings'
    }
  },
  {
    path: '/defi',
    name: 'DeFi',
    component: DeFi,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.defi'
    }
  },
  {
    path: '/nft',
    name: 'NFT',
    component: NFT,
    meta: {
      requiresAuth: true,
      titleKey: 'nav.nft'
    }
  },

  // 404 重定向
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - PIN 认证检查
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // 恢复会话
  if (!authStore.isAuthenticated) {
    authStore.restoreSession()
  }

  const requiresAuth = to.meta.requiresAuth !== false

  // 需要认证的页面
  if (requiresAuth) {
    if (!authStore.isLoggedIn) {
      return next('/login')
    }
  }

  // 已登录用户访问登录/注册页面，重定向到 dashboard
  if ((to.path === '/login' || to.path === '/register') && authStore.isLoggedIn) {
    return next('/dashboard')
  }

  next()
})

export default router
