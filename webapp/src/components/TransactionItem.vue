<script setup>
/**
 * 交易记录单条组件
 * 显示交易类型图标、时间、资产金额、来源、手续费
 */
import {
  PhArrowDown,
  PhArrowUp,
  PhArrowsLeftRight,
  PhArrowRight,
  PhDownloadSimple,
  PhUploadSimple
} from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  tx: { type: Object, required: true }
})

const { formatNumber } = useFormatters()
const { t } = useI18n()

// 交易类型配置
const typeConfig = {
  buy: { icon: PhArrowDown, color: '#10B981', labelKey: 'transaction.typeBuy' },
  sell: { icon: PhArrowUp, color: '#EF4444', labelKey: 'transaction.typeSell' },
  swap: { icon: PhArrowsLeftRight, color: '#3B82F6', labelKey: 'transaction.typeSwap' },
  transfer: { icon: PhArrowRight, color: '#F59E0B', labelKey: 'transaction.typeTransfer' },
  deposit: { icon: PhDownloadSimple, color: '#10B981', labelKey: 'transaction.typeDeposit' },
  withdraw: { icon: PhUploadSimple, color: '#EF4444', labelKey: 'transaction.typeWithdraw' },
}

const config = typeConfig[props.tx.type] || typeConfig.transfer

// 格式化时间（HH:MM）
const formatTime = (ts) => {
  const d = new Date(ts)
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// 区块浏览器链接
const explorerUrl = (tx) => {
  if (!tx.txHash || !tx.chain) return null
  const explorers = {
    ETH: 'https://etherscan.io/tx/',
    BSC: 'https://bscscan.com/tx/',
    SOL: 'https://solscan.io/tx/',
    MATIC: 'https://polygonscan.com/tx/',
    BTC: 'https://mempool.space/tx/',
  }
  const base = explorers[tx.chain]
  return base ? `${base}${tx.txHash}` : null
}
</script>

<template>
  <div class="tx-item">
    <!-- 类型图标 -->
    <div class="tx-icon" :style="{ background: config.color + '18', color: config.color }">
      <component :is="config.icon" :size="16" weight="bold" />
    </div>

    <!-- 主内容 -->
    <div class="tx-main">
      <div class="tx-row-top">
        <!-- 类型标签 + 来源 -->
        <div class="tx-type-source">
          <span class="tx-type" :style="{ color: config.color }">{{ t(config.labelKey) }}</span>
          <span class="tx-source">{{ tx.source }}</span>
          <span v-if="tx.chain" class="tx-chain-tag">{{ tx.chain }}</span>
        </div>
        <!-- 时间 -->
        <span class="tx-time font-mono">{{ formatTime(tx.timestamp) }}</span>
      </div>

      <div class="tx-row-bottom">
        <!-- 资产变动 -->
        <div class="tx-amounts font-mono">
          <span v-if="tx.type === 'transfer'" class="tx-amount">
            {{ formatNumber(tx.from.amount, 4) }} {{ tx.from.symbol }}
          </span>
          <template v-else>
            <span class="tx-amount-from">-{{ formatNumber(tx.from.amount, 4) }} {{ tx.from.symbol }}</span>
            <span class="tx-arrow">&rarr;</span>
            <span class="tx-amount-to">+{{ formatNumber(tx.to.amount, 4) }} {{ tx.to.symbol }}</span>
          </template>
        </div>
        <!-- 手续费 -->
        <div class="tx-fee font-mono" v-if="tx.fee && tx.fee.amount > 0">
          {{ t('transaction.fee') }}: {{ formatNumber(tx.fee.amount, 6) }} {{ tx.fee.currency }}
        </div>
      </div>

      <!-- 备注 / 区块浏览器链接 -->
      <div v-if="tx.note || explorerUrl(tx)" class="tx-row-extra">
        <span v-if="tx.note" class="tx-note">{{ tx.note }}</span>
        <a
          v-if="explorerUrl(tx)"
          :href="explorerUrl(tx)"
          target="_blank"
          rel="noopener noreferrer"
          class="tx-explorer-link"
        >
          {{ t('transaction.viewOnExplorer') }} &rarr;
        </a>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tx-item {
  display: flex;
  gap: var(--gap-md);
  padding: var(--gap-md);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.tx-item:hover {
  background: var(--color-bg-tertiary);
}

/* 图标 */
.tx-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* 主内容 */
.tx-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.tx-row-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tx-type-source {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.tx-type {
  font-size: 0.8125rem;
  font-weight: 600;
}

.tx-source {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.tx-chain-tag {
  font-size: 0.5625rem;
  padding: 1px 4px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
  color: var(--color-accent-primary);
  font-weight: 500;
}

.tx-time {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 金额行 */
.tx-row-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tx-amounts {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  font-size: 0.8125rem;
}

.tx-amount {
  color: var(--color-text-primary);
  font-weight: 500;
}

.tx-amount-from {
  color: var(--color-text-secondary);
}

.tx-arrow {
  color: var(--color-text-muted);
  font-size: 0.75rem;
}

.tx-amount-to {
  color: var(--color-text-primary);
  font-weight: 500;
}

.tx-fee {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 额外信息行 */
.tx-row-extra {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  margin-top: 2px;
}

.tx-note {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  font-style: italic;
}

.tx-explorer-link {
  font-size: 0.6875rem;
  color: var(--color-accent-primary);
  text-decoration: none;
}

.tx-explorer-link:hover {
  text-decoration: underline;
}
</style>
