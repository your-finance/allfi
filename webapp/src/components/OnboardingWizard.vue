<script setup>
/**
 * 首次使用引导向导
 * 3 步：欢迎介绍 → 添加账户 → 设置偏好
 */
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'
import { useFormatters } from '../composables/useFormatters'

const emit = defineEmits(['complete'])

const router = useRouter()
const themeStore = useThemeStore()
const { t } = useI18n()
const { currencies, currentCurrency } = useFormatters()

const currentStep = ref(0)

// 计价货币选择
const selectedCurrency = ref(currentCurrency.value || 'USDC')

// 主题选择
const selectedTheme = ref(themeStore.currentThemeId)

// 跳过引导
const skip = () => {
  themeStore.completeOnboarding()
  emit('complete')
}

// 下一步
const next = () => {
  if (currentStep.value < 2) {
    currentStep.value++
  } else {
    finish()
  }
}

// 上一步
const prev = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 跳转到添加账户
const goToAddAccount = (type) => {
  themeStore.completeOnboarding()
  emit('complete')
  router.push('/accounts')
}

// 完成引导
const finish = () => {
  // 应用选择的偏好
  if (selectedCurrency.value !== currentCurrency.value) {
    currentCurrency.value = selectedCurrency.value
  }
  if (selectedTheme.value !== themeStore.currentThemeId) {
    themeStore.setTheme(selectedTheme.value)
  }
  themeStore.completeOnboarding()
  emit('complete')
}
</script>

<template>
  <div class="onboarding-overlay">
    <div class="onboarding-card">
      <!-- Step 1: 欢迎 -->
      <div v-if="currentStep === 0" class="step-content">
        <div class="welcome-logo">
          <svg viewBox="0 0 56 56" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="2" y="2" width="52" height="52" rx="10" stroke="var(--color-accent-primary)" stroke-width="2" fill="none"/>
            <path d="M16 28 L28 16 L40 28 L28 40 Z" fill="var(--color-accent-primary)" opacity="0.8"/>
          </svg>
        </div>
        <h2 class="step-title">{{ t('onboarding.welcomeTitle') }}</h2>
        <p class="step-desc">{{ t('onboarding.welcomeDesc') }}</p>
        <div class="step-features">
          <div class="feature-item">
            <span class="feature-icon">📊</span>
            <span>{{ t('onboarding.feature1') }}</span>
          </div>
          <div class="feature-item">
            <span class="feature-icon">🔒</span>
            <span>{{ t('onboarding.feature2') }}</span>
          </div>
          <div class="feature-item">
            <span class="feature-icon">🌐</span>
            <span>{{ t('onboarding.feature3') }}</span>
          </div>
        </div>
      </div>

      <!-- Step 2: 添加账户 -->
      <div v-if="currentStep === 1" class="step-content">
        <h2 class="step-title">{{ t('onboarding.addAccountTitle') }}</h2>
        <p class="step-desc">{{ t('onboarding.addAccountDesc') }}</p>
        <div class="account-options">
          <button class="account-option" @click="goToAddAccount('cex')">
            <span class="option-icon">🏦</span>
            <div class="option-text">
              <span class="option-name">{{ t('dashboard.cexAssets') }}</span>
              <span class="option-desc">{{ t('onboarding.cexHint') }}</span>
            </div>
          </button>
          <button class="account-option" @click="goToAddAccount('wallet')">
            <span class="option-icon">💎</span>
            <div class="option-text">
              <span class="option-name">{{ t('dashboard.blockchainAssets') }}</span>
              <span class="option-desc">{{ t('onboarding.walletHint') }}</span>
            </div>
          </button>
          <button class="account-option" @click="goToAddAccount('manual')">
            <span class="option-icon">🏠</span>
            <div class="option-text">
              <span class="option-name">{{ t('dashboard.manualAssets') }}</span>
              <span class="option-desc">{{ t('onboarding.manualHint') }}</span>
            </div>
          </button>
        </div>
      </div>

      <!-- Step 3: 设置偏好 -->
      <div v-if="currentStep === 2" class="step-content">
        <h2 class="step-title">{{ t('onboarding.preferencesTitle') }}</h2>
        <p class="step-desc">{{ t('onboarding.preferencesDesc') }}</p>

        <div class="pref-group">
          <label class="pref-label">{{ t('settings.currency') }}</label>
          <div class="pref-options">
            <button
              v-for="c in currencies"
              :key="c.code"
              class="pref-btn"
              :class="{ active: selectedCurrency === c.code }"
              @click="selectedCurrency = c.code"
            >
              {{ c.symbol }} {{ c.code }}
            </button>
          </div>
        </div>

        <div class="pref-group">
          <label class="pref-label">{{ t('settings.themeSection') }}</label>
          <div class="pref-options">
            <button
              v-for="theme in themeStore.availableThemes"
              :key="theme.id"
              class="pref-btn"
              :class="{ active: selectedTheme === theme.id }"
              @click="selectedTheme = theme.id"
            >
              {{ theme.name }}
            </button>
          </div>
        </div>
      </div>

      <!-- 底部：进度 + 按钮 -->
      <div class="onboarding-footer">
        <button class="btn-skip" @click="skip">{{ t('onboarding.skip') }}</button>

        <div class="step-dots">
          <span
            v-for="i in 3"
            :key="i"
            class="dot"
            :class="{ active: currentStep === i - 1 }"
          />
        </div>

        <div class="footer-actions">
          <button v-if="currentStep > 0" class="btn btn-ghost" @click="prev">
            {{ t('onboarding.prev') }}
          </button>
          <button class="btn btn-primary" @click="next">
            {{ currentStep === 2 ? t('onboarding.finish') : t('onboarding.next') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.onboarding-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--gap-lg);
}

.onboarding-card {
  width: 480px;
  max-width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--gap-2xl);
  display: flex;
  flex-direction: column;
}

/* 步骤内容 */
.step-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--gap-md);
}

.welcome-logo svg {
  width: 56px;
  height: 56px;
}

.step-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.step-desc {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
  max-width: 380px;
}

/* 特性列表 */
.step-features {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  margin-top: var(--gap-sm);
  width: 100%;
  max-width: 320px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  text-align: left;
}

.feature-icon {
  font-size: 1rem;
  flex-shrink: 0;
}

/* 账户选项 */
.account-options {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  width: 100%;
  margin-top: var(--gap-sm);
}

.account-option {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  text-align: left;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.account-option:hover {
  border-color: var(--color-accent-primary);
  background: var(--color-bg-elevated);
}

.option-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.option-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.option-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.option-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 偏好设置 */
.pref-group {
  width: 100%;
  text-align: left;
  margin-top: var(--gap-sm);
}

.pref-label {
  display: block;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-xs);
}

.pref-options {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap-xs);
}

.pref-btn {
  padding: 6px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast), border-color var(--transition-fast);
}

.pref-btn:hover {
  color: var(--color-text-primary);
  border-color: var(--color-accent-primary);
}

.pref-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
  border-color: var(--color-accent-primary);
}

/* 底部 */
.onboarding-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: var(--gap-xl);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

.btn-skip {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  transition: color var(--transition-fast);
}

.btn-skip:hover {
  color: var(--color-text-primary);
}

.step-dots {
  display: flex;
  gap: 6px;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-bg-tertiary);
  transition: background var(--transition-fast);
}

.dot.active {
  background: var(--color-accent-primary);
}

.footer-actions {
  display: flex;
  gap: var(--gap-xs);
}
</style>
