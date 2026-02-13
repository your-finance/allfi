<script setup>
/**
 * 钱包连接按钮组件
 * 未连接：显示"连接钱包"按钮
 * 已连接：显示地址缩写和链标识
 */
import { PhWallet, PhLinkBreak } from '@phosphor-icons/vue'
import { useWalletConnect } from '../composables/useWalletConnect'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const emit = defineEmits(['connect-click'])

const { isConnected, currentAddress, currentChain, disconnectWallet } = useWalletConnect()
const { shortenAddress } = useFormatters()
const { t } = useI18n()
</script>

<template>
  <div class="wallet-connect-btn-wrap">
    <!-- 未连接状态 -->
    <button
      v-if="!isConnected"
      class="btn btn-primary btn-sm wallet-btn"
      @click="emit('connect-click')"
    >
      <PhWallet :size="14" />
      {{ t('walletConnect.connectWallet') }}
    </button>

    <!-- 已连接状态 -->
    <div v-else class="connected-info">
      <span class="chain-tag" v-if="currentChain">{{ currentChain.id }}</span>
      <span class="address-text font-mono">{{ shortenAddress(currentAddress, 4) }}</span>
      <button
        class="disconnect-btn"
        :title="t('walletConnect.disconnect')"
        @click="disconnectWallet"
      >
        <PhLinkBreak :size="12" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.wallet-connect-btn-wrap {
  display: inline-flex;
}

.wallet-btn {
  gap: var(--gap-xs);
}

.connected-info {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 4px 8px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.chain-tag {
  font-size: 0.625rem;
  font-weight: 600;
  padding: 1px 4px;
  border-radius: var(--radius-xs);
  background: rgba(16, 185, 129, 0.15);
  color: var(--color-success);
}

.address-text {
  font-size: 0.75rem;
  color: var(--color-text-primary);
}

.disconnect-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: color var(--transition-fast);
}

.disconnect-btn:hover {
  color: var(--color-error);
}
</style>
