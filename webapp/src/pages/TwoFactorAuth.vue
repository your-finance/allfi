<script setup>
/**
 * 2FA 验证页面
 * 输入 6 位验证码完成二次验证
 */
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { 
  PhShieldCheck,
  PhSpinner,
  PhWarning,
  PhArrowLeft
} from '@phosphor-icons/vue'
import { useAuthStore } from '../stores/authStore'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { t } = useI18n()

// 6位验证码
const code = ref(['', '', '', '', '', ''])
const inputRefs = ref([])

// 合并验证码
const fullCode = computed(() => code.value.join(''))

// 初始化
onMounted(() => {
  themeStore.initTheme()
  
  // 检查是否需要 2FA
  if (!authStore.isAuthenticated || !authStore.requires2FA) {
    router.replace('/login')
    return
  }
  
  // 聚焦第一个输入框
  nextTick(() => {
    if (inputRefs.value[0]) {
      inputRefs.value[0].focus()
    }
  })
})

// 处理输入
const handleInput = (index, event) => {
  const value = event.target.value
  
  // 只允许数字
  if (!/^\d*$/.test(value)) {
    code.value[index] = ''
    return
  }
  
  // 只保留最后一个数字
  code.value[index] = value.slice(-1)
  
  // 自动跳到下一个
  if (value && index < 5) {
    inputRefs.value[index + 1]?.focus()
  }
  
  // 如果填满了自动提交
  if (fullCode.value.length === 6) {
    handleSubmit()
  }
}

// 处理删除键
const handleKeydown = (index, event) => {
  if (event.key === 'Backspace' && !code.value[index] && index > 0) {
    inputRefs.value[index - 1]?.focus()
  }
}

// 处理粘贴
const handlePaste = (event) => {
  event.preventDefault()
  const paste = event.clipboardData.getData('text')
  const digits = paste.replace(/\D/g, '').slice(0, 6)
  
  digits.split('').forEach((digit, i) => {
    if (i < 6) code.value[i] = digit
  })
  
  // 聚焦到最后一个填充的位置或提交
  if (digits.length === 6) {
    handleSubmit()
  } else {
    inputRefs.value[digits.length]?.focus()
  }
}

// 提交验证
const handleSubmit = async () => {
  if (fullCode.value.length !== 6) return
  
  const success = await authStore.verify2FA(fullCode.value)
  
  if (success) {
    router.push('/dashboard')
  } else {
    // 清空验证码
    code.value = ['', '', '', '', '', '']
    inputRefs.value[0]?.focus()
  }
}

// 重发验证码状态
const isResending = ref(false)
const resendSuccess = ref(false)
const resendCooldown = ref(0)
let cooldownTimer = null

// 重发验证码
const resendCode = async () => {
  if (resendCooldown.value > 0 || isResending.value) return

  isResending.value = true
  try {
    await authStore.resend2FACode()
    resendSuccess.value = true

    // 设置冷却时间（60秒）
    resendCooldown.value = 60
    cooldownTimer = setInterval(() => {
      resendCooldown.value--
      if (resendCooldown.value <= 0) {
        clearInterval(cooldownTimer)
        resendSuccess.value = false
      }
    }, 1000)
  } catch (err) {
    console.error('重发验证码失败:', err)
  } finally {
    isResending.value = false
  }
}

// 返回登录
const goBack = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="auth-page">
    <!-- 背景装饰 -->
    <div class="auth-bg">
      <div class="bg-gradient"></div>
      <div class="bg-pattern"></div>
    </div>
    
    <!-- 验证卡片 -->
    <div class="auth-container">
      <div class="auth-card glass-card">
        <!-- 返回按钮 -->
        <button class="back-btn" @click="goBack">
          <PhArrowLeft :size="20" />
          <span>{{ t('auth.backToLogin') }}</span>
        </button>
        
        <!-- 图标 -->
        <div class="auth-header">
          <div class="icon-wrapper">
            <PhShieldCheck :size="48" weight="duotone" />
          </div>
          <h1 class="auth-title">{{ t('auth.twoFactorAuth') }}</h1>
          <p class="auth-subtitle">{{ t('auth.twoFactorSubtitle') }}</p>
          <p class="auth-email">{{ authStore.userEmail }}</p>
        </div>
        
        <!-- 错误提示 -->
        <div v-if="authStore.error" class="error-alert">
          <PhWarning :size="20" weight="fill" />
          <span>{{ authStore.error }}</span>
        </div>
        
        <!-- 验证码输入 -->
        <form @submit.prevent="handleSubmit" class="code-form">
          <div class="code-inputs" @paste="handlePaste">
            <input
              v-for="(_, index) in 6"
              :key="index"
              :ref="el => inputRefs[index] = el"
              type="text"
              inputmode="numeric"
              maxlength="1"
              class="code-input"
              :value="code[index]"
              @input="handleInput(index, $event)"
              @keydown="handleKeydown(index, $event)"
              :disabled="authStore.isLoading"
            />
          </div>
          
          <!-- 提交按钮 -->
          <button 
            type="submit" 
            class="btn btn-primary btn-lg btn-block"
            :disabled="fullCode.length !== 6 || authStore.isLoading"
          >
            <PhSpinner v-if="authStore.isLoading" :size="20" class="spin" />
            <span v-else>{{ t('auth.verify') }}</span>
          </button>
        </form>
        
        <!-- 重发链接 -->
        <div class="auth-footer">
          <template v-if="resendSuccess">
            <span class="resend-success">{{ t('auth.codeSent') || '验证码已发送' }}</span>
          </template>
          <template v-else>
            <span>{{ t('auth.noCode') }}</span>
            <button
              type="button"
              class="resend-btn"
              :disabled="resendCooldown > 0 || isResending"
              @click="resendCode"
            >
              <template v-if="isResending">{{ t('auth.sending') || '发送中...' }}</template>
              <template v-else-if="resendCooldown > 0">{{ t('auth.resendIn') || '重新发送' }} ({{ resendCooldown }}s)</template>
              <template v-else>{{ t('auth.resendCode') }}</template>
            </button>
          </template>
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

/* 返回按钮 */
.back-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  background: none;
  border: none;
  color: var(--color-text-secondary);
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  margin-bottom: var(--gap-md);
  transition: color var(--transition-fast);
}

.back-btn:hover {
  color: var(--color-text-primary);
}

/* 头部 */
.auth-header {
  text-align: center;
  margin-bottom: var(--gap-xl);
}

.icon-wrapper {
  width: 56px;
  height: 56px;
  margin: 0 auto var(--gap-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--color-accent-primary) 12%, transparent);
  border-radius: 50%;
  color: var(--color-accent-primary);
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
  margin-bottom: var(--gap-xs);
}

.auth-email {
  color: var(--color-text-primary);
  font-weight: 500;
  font-size: 13px;
  font-family: var(--font-mono);
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

/* 验证码输入 */
.code-form {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

.code-inputs {
  display: flex;
  justify-content: center;
  gap: var(--gap-sm);
}

.code-input {
  width: 42px;
  height: 48px;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 600;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  transition: border-color var(--transition-fast);
}

.code-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.code-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.resend-btn {
  background: none;
  border: none;
  color: var(--color-accent-primary);
  font-weight: 500;
  cursor: pointer;
  margin-left: var(--gap-xs);
  font-size: inherit;
}

.resend-btn:hover:not(:disabled) {
  text-decoration: underline;
}

.resend-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.resend-success {
  color: var(--color-success);
  font-weight: 500;
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
@media (max-width: 400px) {
  .code-input {
    width: 36px;
    height: 42px;
    font-size: 18px;
  }
}
</style>
