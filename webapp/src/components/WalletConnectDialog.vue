<script setup>
/**
 * 钱包连接对话框
 * 步骤 1：选择连接方式（MetaMask / WalletConnect / 手动输入）
 * 步骤 2：等待钱包授权
 * 步骤 3：显示检测到的地址和链，用户确认添加
 */
import { ref, watch } from 'vue'
import { PhX, PhWallet, PhQrCode, PhPencilSimpleLine, PhSpinnerGap, PhCheck, PhWarning } from '@phosphor-icons/vue'
import { useWalletConnect } from '../composables/useWalletConnect'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'connected'])

const { isConnecting, error, currentChain, connectWallet, clearError, CHAIN_MAP } = useWalletConnect()
const { shortenAddress } = useFormatters()
const { t } = useI18n()

// 步骤：select / connecting / success
const step = ref('select')
const connectedAddress = ref('')
const connectedChainId = ref(null)
const walletName = ref('')

// 可用的连接方式
const connectMethods = [
  { id: 'metamask', labelKey: 'walletConnect.metamask', icon: PhWallet, desc: 'walletConnect.metamaskDesc' },
  { id: 'walletconnect', labelKey: 'walletConnect.walletconnect', icon: PhQrCode, desc: 'walletConnect.walletconnectDesc' },
]

// 重置对话框
const resetDialog = () => {
  step.value = 'select'
  connectedAddress.value = ''
  connectedChainId.value = null
  walletName.value = ''
  clearError()
}

// 监听可见性变化
watch(() => props.visible, (val) => {
  if (!val) resetDialog()
})

// 选择连接方式并发起连接
const handleMethodSelect = async (method) => {
  step.value = 'connecting'
  clearError()

  try {
    const result = await connectWallet(method)
    connectedAddress.value = result.address
    connectedChainId.value = result.chainId
    // 自动生成钱包名称
    const chain = CHAIN_MAP[result.chainId]
    walletName.value = `${chain?.name || 'Unknown'} Wallet`
    step.value = 'success'
  } catch (err) {
    step.value = 'select'
  }
}

// 确认添加钱包
const confirmAdd = () => {
  const chain = CHAIN_MAP[connectedChainId.value]
  emit('connected', {
    address: connectedAddress.value,
    blockchain: chain?.id || 'ETH',
    name: walletName.value || `${chain?.name || 'Web3'} Wallet`,
    chainId: connectedChainId.value
  })
  emit('close')
}
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-content">
        <!-- 头部 -->
        <div class="modal-header">
          <h3>{{ t('walletConnect.title') }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="20" />
          </button>
        </div>

        <div class="modal-body">
          <!-- 步骤 1：选择连接方式 -->
          <div v-if="step === 'select'" class="step-select">
            <p class="step-desc">{{ t('walletConnect.selectMethod') }}</p>
            <div class="method-list">
              <button
                v-for="m in connectMethods"
                :key="m.id"
                class="method-btn"
                @click="handleMethodSelect(m.id)"
              >
                <component :is="m.icon" :size="24" class="method-icon" />
                <div class="method-info">
                  <span class="method-label">{{ t(m.labelKey) }}</span>
                  <span class="method-desc">{{ t(m.desc) }}</span>
                </div>
              </button>
            </div>

            <!-- 错误提示 -->
            <div v-if="error" class="error-box">
              <PhWarning :size="16" weight="fill" />
              <span>{{ error }}</span>
            </div>
          </div>

          <!-- 步骤 2：连接中 -->
          <div v-else-if="step === 'connecting'" class="step-connecting">
            <PhSpinnerGap :size="32" class="spinner" />
            <p class="connecting-text">{{ t('walletConnect.connecting') }}</p>
            <p class="connecting-hint">{{ t('walletConnect.confirmInWallet') }}</p>
          </div>

          <!-- 步骤 3：连接成功 -->
          <div v-else-if="step === 'success'" class="step-success">
            <div class="success-icon">
              <PhCheck :size="24" weight="bold" />
            </div>
            <p class="success-text">{{ t('walletConnect.connected') }}</p>

            <!-- 检测到的信息 -->
            <div class="detected-info">
              <div class="info-row">
                <span class="info-label">{{ t('walletConnect.detectedAddress') }}</span>
                <span class="info-value font-mono">{{ shortenAddress(connectedAddress, 8) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">{{ t('walletConnect.detectedChain') }}</span>
                <span class="info-value">{{ currentChain?.name || 'Unknown' }}</span>
              </div>
            </div>

            <!-- 钱包名称输入 -->
            <div class="form-group">
              <label class="input-label">{{ t('walletConnect.walletNameLabel') }}</label>
              <input
                type="text"
                v-model="walletName"
                class="input-field"
                :placeholder="t('walletConnect.walletNamePlaceholder')"
              />
            </div>
          </div>
        </div>

        <!-- 底部按钮 -->
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="emit('close')">
            {{ t('common.cancel') }}
          </button>
          <button
            v-if="step === 'success'"
            class="btn btn-primary"
            @click="confirmAdd"
          >
            {{ t('walletConnect.addToSystem') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.modal-content {
  width: 100%;
  max-width: 420px;
  padding: var(--gap-xl);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-lg);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
}

.modal-header h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.close-btn:hover {
  color: var(--color-text-primary);
}

.modal-body {
  margin-bottom: var(--gap-lg);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

/* 步骤 1：选择连接方式 */
.step-desc {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-md);
}

.method-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.method-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: border-color var(--transition-fast);
  text-align: left;
}

.method-btn:hover {
  border-color: var(--color-accent-primary);
}

.method-icon {
  color: var(--color-accent-primary);
  flex-shrink: 0;
}

.method-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.method-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.method-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.error-box {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: rgba(226, 92, 92, 0.1);
  border-radius: var(--radius-sm);
  color: var(--color-error);
  font-size: 0.75rem;
  margin-top: var(--gap-md);
}

/* 步骤 2：连接中 */
.step-connecting {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--gap-xl) 0;
  gap: var(--gap-md);
}

.spinner {
  color: var(--color-accent-primary);
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.connecting-text {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.connecting-hint {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* 步骤 3：连接成功 */
.step-success {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-md);
}

.success-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(16, 185, 129, 0.15);
  border-radius: 50%;
  color: var(--color-success);
}

.success-text {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-success);
}

.detected-info {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.info-value {
  font-size: 0.75rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.form-group {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.input-label {
  font-size: 0.75rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.input-field {
  width: 100%;
  padding: 8px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.8125rem;
}

.input-field:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

/* 动画 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
