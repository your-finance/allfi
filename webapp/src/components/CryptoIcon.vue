<script setup>
/**
 * 加密货币图标组件
 * 根据币种符号显示对应颜色和缩写
 */
import { computed } from 'vue'

const props = defineProps({
  symbol: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  }
})

// 币种颜色映射
const colorMap = {
  BTC: { bg: '#F7931A', text: '#FFFFFF' },
  ETH: { bg: '#627EEA', text: '#FFFFFF' },
  USDC: { bg: '#26A17B', text: '#FFFFFF' },
  USDC: { bg: '#2775CA', text: '#FFFFFF' },
  BNB: { bg: '#F3BA2F', text: '#1E2026' },
  SOL: { bg: '#9945FF', text: '#FFFFFF' },
  XRP: { bg: '#23292F', text: '#FFFFFF' },
  ADA: { bg: '#0033AD', text: '#FFFFFF' },
  DOGE: { bg: '#C3A634', text: '#FFFFFF' },
  DOT: { bg: '#E6007A', text: '#FFFFFF' },
  MATIC: { bg: '#8247E5', text: '#FFFFFF' },
  AVAX: { bg: '#E84142', text: '#FFFFFF' },
  LINK: { bg: '#2A5ADA', text: '#FFFFFF' },
  UNI: { bg: '#FF007A', text: '#FFFFFF' },
  ATOM: { bg: '#2E3148', text: '#FFFFFF' },
  LTC: { bg: '#345D9D', text: '#FFFFFF' },
  TRX: { bg: '#FF0013', text: '#FFFFFF' },
  NEAR: { bg: '#000000', text: '#FFFFFF' },
  APT: { bg: '#000000', text: '#FFFFFF' },
  ARB: { bg: '#28A0F0', text: '#FFFFFF' },
  OP: { bg: '#FF0420', text: '#FFFFFF' },
  CNY: { bg: '#DE2910', text: '#FFFFFF' },
  USD: { bg: '#22C55E', text: '#FFFFFF' },
}

// 获取币种颜色
const colors = computed(() => {
  return colorMap[props.symbol.toUpperCase()] || { bg: '#64748B', text: '#FFFFFF' }
})

// 尺寸映射
const sizeClasses = {
  sm: 'w-6 h-6 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-12 h-12 text-base'
}

const sizeClass = computed(() => sizeClasses[props.size])

// 获取显示缩写（最多3个字符）
const displaySymbol = computed(() => {
  return props.symbol.slice(0, 3).toUpperCase()
})
</script>

<template>
  <div 
    class="crypto-icon"
    :class="sizeClass"
    :style="{ 
      backgroundColor: colors.bg,
      color: colors.text
    }"
  >
    {{ displaySymbol }}
  </div>
</template>

<style scoped>
.crypto-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-weight: 700;
  font-family: var(--font-mono);
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}
</style>
