<script setup>
/**
 * Register 页面 - 用户注册
 */
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  PhEnvelope, 
  PhLock, 
  PhUser,
  PhEye, 
  PhEyeSlash,
  PhSpinner,
  PhWarning,
  PhCheck
} from '@phosphor-icons/vue'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { t } = useI18n()

// 表单数据
const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const agreeTerms = ref(false)

// 表单验证错误
const errors = ref({})

// 初始化主题
onMounted(() => {
  themeStore.initTheme()
})

// 密码强度检查
const passwordStrength = ref(0)
const checkPasswordStrength = () => {
  let strength = 0
  if (password.value.length >= 6) strength++
  if (password.value.length >= 10) strength++
  if (/[A-Z]/.test(password.value)) strength++
  if (/[0-9]/.test(password.value)) strength++
  if (/[^A-Za-z0-9]/.test(password.value)) strength++
  passwordStrength.value = strength
}

// 验证表单
const validateForm = () => {
  errors.value = {}
  
  if (!name.value.trim()) {
    errors.value.name = t('auth.nameRequired')
  }
  
  if (!email.value) {
    errors.value.email = t('auth.emailRequired')
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errors.value.email = t('auth.emailInvalid')
  }
  
  if (!password.value) {
    errors.value.password = t('auth.passwordRequired')
  } else if (password.value.length < 6) {
    errors.value.password = t('auth.passwordTooShort')
  }
  
  if (password.value !== confirmPassword.value) {
    errors.value.confirmPassword = t('auth.passwordMismatch')
  }
  
  if (!agreeTerms.value) {
    errors.value.terms = t('auth.termsRequired')
  }
  
  return Object.keys(errors.value).length === 0
}

// 提交注册
const handleSubmit = async () => {
  if (!validateForm()) return
  
  const success = await authStore.register({
    name: name.value,
    email: email.value,
    password: password.value
  })
  
  if (success) {
    router.push('/dashboard')
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
    
    <!-- 注册卡片 -->
    <div class="auth-container">
      <div class="auth-card glass-card">
        <!-- Logo -->
        <div class="auth-header">
          <div class="logo">
            <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="16" cy="16" r="14" stroke="url(#logo-gradient)" stroke-width="2"/>
              <path d="M10 16 L16 10 L22 16 L16 22 Z" fill="url(#logo-gradient)"/>
              <defs>
                <linearGradient id="logo-gradient" x1="0" y1="0" x2="32" y2="32">
                  <stop offset="0%" :stop-color="themeStore.currentTheme.colors.accentPrimary"/>
                  <stop offset="100%" :stop-color="themeStore.currentTheme.colors.accentSecondary"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h1 class="auth-title">{{ t('auth.createAccount') }}</h1>
          <p class="auth-subtitle">{{ t('auth.registerSubtitle') }}</p>
        </div>
        
        <!-- 错误提示 -->
        <div v-if="authStore.error" class="error-alert">
          <PhWarning :size="20" weight="fill" />
          <span>{{ authStore.error }}</span>
        </div>
        
        <!-- 注册表单 -->
        <form @submit.prevent="handleSubmit" class="auth-form">
          <!-- 用户名 -->
          <div class="form-group">
            <label class="form-label">{{ t('auth.name') }}</label>
            <div class="input-wrapper" :class="{ 'has-error': errors.name }">
              <PhUser :size="20" class="input-icon" />
              <input
                type="text"
                v-model="name"
                class="form-input"
                :placeholder="t('auth.namePlaceholder')"
              />
            </div>
            <span v-if="errors.name" class="error-text">{{ errors.name }}</span>
          </div>
          
          <!-- 邮箱 -->
          <div class="form-group">
            <label class="form-label">{{ t('auth.email') }}</label>
            <div class="input-wrapper" :class="{ 'has-error': errors.email }">
              <PhEnvelope :size="20" class="input-icon" />
              <input
                type="email"
                v-model="email"
                class="form-input"
                :placeholder="t('auth.emailPlaceholder')"
              />
            </div>
            <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
          </div>
          
          <!-- 密码 -->
          <div class="form-group">
            <label class="form-label">{{ t('auth.password') }}</label>
            <div class="input-wrapper" :class="{ 'has-error': errors.password }">
              <PhLock :size="20" class="input-icon" />
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="password"
                class="form-input"
                :placeholder="t('auth.passwordPlaceholder')"
                @input="checkPasswordStrength"
              />
              <button 
                type="button" 
                class="password-toggle"
                @click="showPassword = !showPassword"
              >
                <PhEyeSlash v-if="showPassword" :size="20" />
                <PhEye v-else :size="20" />
              </button>
            </div>
            <!-- 密码强度指示 -->
            <div v-if="password" class="password-strength">
              <div class="strength-bars">
                <div 
                  v-for="i in 5" 
                  :key="i"
                  class="strength-bar"
                  :class="{ 
                    'active': i <= passwordStrength,
                    'weak': passwordStrength <= 2,
                    'medium': passwordStrength === 3,
                    'strong': passwordStrength >= 4
                  }"
                />
              </div>
              <span class="strength-text">
                {{ passwordStrength <= 2 ? t('auth.passwordWeak') : 
                   passwordStrength === 3 ? t('auth.passwordMedium') : t('auth.passwordStrong') }}
              </span>
            </div>
            <span v-if="errors.password" class="error-text">{{ errors.password }}</span>
          </div>
          
          <!-- 确认密码 -->
          <div class="form-group">
            <label class="form-label">{{ t('auth.confirmPassword') }}</label>
            <div class="input-wrapper" :class="{ 'has-error': errors.confirmPassword }">
              <PhLock :size="20" class="input-icon" />
              <input
                type="password"
                v-model="confirmPassword"
                class="form-input"
                :placeholder="t('auth.confirmPasswordPlaceholder')"
              />
              <PhCheck 
                v-if="confirmPassword && password === confirmPassword" 
                :size="20" 
                class="input-icon-right success"
                weight="bold"
              />
            </div>
            <span v-if="errors.confirmPassword" class="error-text">{{ errors.confirmPassword }}</span>
          </div>
          
          <!-- 服务条款 -->
          <div class="form-group">
            <label class="checkbox-label" :class="{ 'has-error': errors.terms }">
              <input type="checkbox" v-model="agreeTerms" class="checkbox" />
              <span class="checkbox-text">
                {{ t('auth.agreeTerms') }}
                <a href="#" class="auth-link">{{ t('auth.termsOfService') }}</a>
              </span>
            </label>
            <span v-if="errors.terms" class="error-text">{{ errors.terms }}</span>
          </div>
          
          <!-- 注册按钮 -->
          <button 
            type="submit" 
            class="btn btn-primary btn-lg btn-block"
            :disabled="authStore.isLoading"
          >
            <PhSpinner v-if="authStore.isLoading" :size="20" class="spin" />
            <span v-else>{{ t('auth.register') }}</span>
          </button>
        </form>
        
        <!-- 登录链接 -->
        <div class="auth-footer">
          <span>{{ t('auth.hasAccount') }}</span>
          <router-link to="/login" class="auth-link">{{ t('auth.login') }}</router-link>
        </div>
      </div>
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
  padding: var(--gap-xl) 0;
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
}

/* 头部 */
.auth-header {
  text-align: center;
  margin-bottom: var(--gap-lg);
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
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-xs);
}

.auth-subtitle {
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

.input-icon-right {
  position: absolute;
  right: 12px;
  pointer-events: none;
}

.input-icon-right.success {
  color: var(--color-success);
}

.form-input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 13px;
  transition: border-color var(--transition-fast);
}

.form-input::placeholder {
  color: var(--color-text-muted);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.input-wrapper.has-error .form-input {
  border-color: var(--color-error);
}

.password-toggle {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: 4px;
  display: flex;
}

.password-toggle:hover {
  color: var(--color-text-secondary);
}

.error-text {
  font-size: 12px;
  color: var(--color-error);
}

/* 密码强度 */
.password-strength {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-top: var(--gap-xs);
}

.strength-bars {
  display: flex;
  gap: 3px;
}

.strength-bar {
  width: 28px;
  height: 3px;
  background: var(--color-bg-elevated);
  border-radius: 1px;
  transition: background var(--transition-fast);
}

.strength-bar.active.weak { background: var(--color-error); }
.strength-bar.active.medium { background: var(--color-warning); }
.strength-bar.active.strong { background: var(--color-success); }

.strength-text {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* 复选框 */
.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-xs);
  cursor: pointer;
}

.checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent-primary);
  margin-top: 2px;
  flex-shrink: 0;
}

.checkbox-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  line-height: 1.4;
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

/* 底部 */
.auth-footer {
  text-align: center;
  margin-top: var(--gap-lg);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
  font-size: 13px;
  color: var(--color-text-secondary);
}

.auth-link {
  color: var(--color-accent-primary);
  text-decoration: none;
  font-weight: 500;
  margin-left: var(--gap-xs);
}

.auth-link:hover {
  text-decoration: underline;
}

/* 动画 */
.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
