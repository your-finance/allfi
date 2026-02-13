<script setup>
/**
 * Login 页面 - PIN 认证
 * 支持首次设置 PIN 和 PIN 登录两种模式
 */
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  PhLock,
  PhSpinner,
  PhWarning,
  PhShieldCheck
} from '@phosphor-icons/vue'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { t } = useI18n()

// 页面模式：setup（首次设置 PIN）或 login（登录）
const mode = ref('login')
const pin = ref('')
const confirmPin = ref('')
const pinError = ref('')
const isChecking = ref(true)

// 初始化：检查认证状态
onMounted(async () => {
  themeStore.initTheme()
  authStore.clearError()
  const isPinSet = await authStore.checkAuthStatus()
  mode.value = isPinSet ? 'login' : 'setup'
  isChecking.value = false
})

// 验证 PIN 格式
const validatePin = (value) => {
  if (!value) return '请输入 PIN'
  if (!/^\d+$/.test(value)) return 'PIN 只能包含数字'
  if (value.length < 4 || value.length > 8) return 'PIN 长度为 4-8 位'
  return ''
}

// 提交表单
const handleSubmit = async () => {
  pinError.value = ''
  authStore.clearError()

  if (mode.value === 'setup') {
    // 首次设置模式
    const err = validatePin(pin.value)
    if (err) { pinError.value = err; return }
    if (pin.value !== confirmPin.value) {
      pinError.value = '两次输入的 PIN 不一致'
      return
    }
    const success = await authStore.setupPIN(pin.value)
    if (success) {
      router.push('/dashboard')
    }
  } else {
    // 登录模式
    if (!pin.value) { pinError.value = '请输入 PIN'; return }
    const success = await authStore.login(pin.value)
    if (success) {
      router.push('/dashboard')
    }
  }
}
</script>

<template>
  <div class="auth-page">
    <!-- 背景装饰 -->
    <div class="auth-bg">
      <div class="bg-gradient"></div>
      <div class="bg-pattern"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="auth-container">
      <div class="auth-card">
        <!-- Logo -->
        <div class="auth-header">
          <div class="logo">
            <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="16" cy="16" r="14" stroke="var(--color-accent-primary)" stroke-width="2"/>
              <path d="M10 16 L16 10 L22 16 L16 22 Z" fill="var(--color-accent-primary)"/>
            </svg>
          </div>
          <h1 class="auth-title">AllFi</h1>
          <p class="auth-subtitle">
            {{ mode === 'setup' ? '首次使用，请设置 PIN 码' : '请输入 PIN 码解锁' }}
          </p>
        </div>

        <!-- 加载中 -->
        <div v-if="isChecking" class="loading-state">
          <PhSpinner :size="24" class="spin" />
          <span>检查认证状态...</span>
        </div>

        <template v-else>
          <!-- 错误提示 -->
          <div v-if="authStore.error || pinError" class="error-alert">
            <PhWarning :size="20" weight="fill" />
            <span>{{ authStore.error || pinError }}</span>
          </div>

          <!-- PIN 表单 -->
          <form @submit.prevent="handleSubmit" class="auth-form">
            <!-- PIN 输入 -->
            <div class="form-group">
              <label class="form-label">
                {{ mode === 'setup' ? '设置 PIN 码（4-8 位数字）' : 'PIN 码' }}
              </label>
              <div class="input-wrapper" :class="{ 'has-error': pinError }">
                <PhLock :size="20" class="input-icon" />
                <input
                  type="password"
                  v-model="pin"
                  class="form-input"
                  placeholder="输入 PIN 码"
                  inputmode="numeric"
                  maxlength="8"
                  autocomplete="off"
                  @input="pinError = ''; authStore.clearError()"
                />
              </div>
            </div>

            <!-- 确认 PIN（仅设置模式） -->
            <div v-if="mode === 'setup'" class="form-group">
              <label class="form-label">确认 PIN 码</label>
              <div class="input-wrapper">
                <PhShieldCheck :size="20" class="input-icon" />
                <input
                  type="password"
                  v-model="confirmPin"
                  class="form-input"
                  placeholder="再次输入 PIN 码"
                  inputmode="numeric"
                  maxlength="8"
                  autocomplete="off"
                  @input="pinError = ''"
                />
              </div>
            </div>

            <!-- 提交按钮 -->
            <button
              type="submit"
              class="btn btn-primary btn-lg btn-block"
              :disabled="authStore.isLoading"
            >
              <PhSpinner v-if="authStore.isLoading" :size="20" class="spin" />
              <span v-else>{{ mode === 'setup' ? '设置 PIN 并进入' : '解锁' }}</span>
            </button>
          </form>

          <!-- 安全提示 -->
          <div class="security-note">
            <PhShieldCheck :size="16" />
            <span>PIN 使用 bcrypt 加密存储，所有数据仅保存在本地</span>
          </div>
        </template>
      </div>

      <!-- 版权 -->
      <p class="auth-copyright">© 2026 AllFi. Self-hosted asset manager.</p>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: var(--color-bg-primary);
}

/* 背景 */
.auth-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.bg-gradient {
  position: absolute;
  inset: 0;
  background: var(--color-bg-primary);
}

.bg-pattern {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(var(--color-border) 1px, transparent 1px),
    linear-gradient(90deg, var(--color-border) 1px, transparent 1px);
  background-size: 60px 60px;
  opacity: 0.15;
}

/* 容器 */
.auth-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 380px;
  padding: var(--gap-lg);
}

.auth-card {
  padding: var(--gap-xl);
  border-radius: var(--radius-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
}

/* 头部 */
.auth-header {
  text-align: center;
  margin-bottom: var(--gap-xl);
}

.logo {
  width: 40px;
  height: 40px;
  margin: 0 auto var(--gap-sm);
}

.logo svg {
  width: 100%;
  height: 100%;
}

.auth-title {
  font-family: var(--font-heading);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-xs);
}

.auth-subtitle {
  color: var(--color-text-secondary);
  font-size: 13px;
}

/* 加载中 */
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  padding: var(--gap-xl) 0;
  color: var(--color-text-secondary);
  font-size: 13px;
}

/* 错误提示 */
.error-alert {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: color-mix(in srgb, var(--color-error) 10%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-error) 30%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-error);
  font-size: 13px;
  margin-bottom: var(--gap-md);
}

/* 表单 */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.form-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 12px;
  color: var(--color-text-muted);
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 16px;
  font-family: var(--font-mono);
  letter-spacing: 4px;
  transition: border-color var(--transition-fast);
}

.form-input::placeholder {
  color: var(--color-text-muted);
  font-size: 13px;
  letter-spacing: normal;
  font-family: var(--font-body);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.input-wrapper.has-error .form-input {
  border-color: var(--color-error);
}

/* 按钮 */
.btn-block {
  width: 100%;
  justify-content: center;
}

.btn-lg {
  padding: 10px 20px;
  font-size: 13px;
}

/* 安全提示 */
.security-note {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  margin-top: var(--gap-lg);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 版权 */
.auth-copyright {
  text-align: center;
  margin-top: var(--gap-md);
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 动画 */
.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 响应式 */
@media (max-width: 480px) {
  .auth-container {
    padding: var(--gap-md);
  }

  .auth-card {
    padding: var(--gap-lg);
  }
}
</style>
