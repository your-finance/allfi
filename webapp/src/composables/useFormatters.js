/**
 * 格式化工具 Composable
 * 提供货币、数字、日期等格式化功能
 */
import { ref, computed } from 'vue'
import { useThemeStore } from '../stores/themeStore'

// 支持的货币配置
export const currencies = [
  { code: 'USDC', symbol: '$', name: 'USDC' },
  { code: 'BTC', symbol: '₿', name: 'Bitcoin' },
  { code: 'ETH', symbol: 'Ξ', name: 'Ethereum' },
  { code: 'CNY', symbol: '¥', name: '人民币' }
]

// 隐私模式遮盖字符
const PRIVACY_MASK = '••••'

export function useFormatters() {
  const themeStore = useThemeStore()
  const currentCurrency = ref('USDC')
  
  // 获取当前货币符号
  const currencySymbol = computed(() => {
    const currency = currencies.find(c => c.code === currentCurrency.value)
    return currency?.symbol || '$'
  })
  
  /**
   * 格式化货币数值
   * @param {number} value - 数值
   * @param {object} options - 格式化选项
   */
  const formatCurrency = (value, options = {}) => {
    const {
      decimals = 2,
      showSymbol = true,
      compact = false
    } = options

    // 隐私模式下遮盖金额
    if (themeStore.privacyMode) {
      return showSymbol ? `${currencySymbol.value}${PRIVACY_MASK}` : PRIVACY_MASK
    }

    if (value === null || value === undefined || isNaN(value)) {
      return showSymbol ? `${currencySymbol.value}0.00` : '0.00'
    }
    
    let formatted
    
    if (compact && Math.abs(value) >= 1000000000) {
      formatted = (value / 1000000000).toFixed(2) + 'B'
    } else if (compact && Math.abs(value) >= 1000000) {
      formatted = (value / 1000000).toFixed(2) + 'M'
    } else if (compact && Math.abs(value) >= 1000) {
      formatted = (value / 1000).toFixed(2) + 'K'
    } else {
      formatted = new Intl.NumberFormat('en-US', {
        minimumFractionDigits: decimals,
        maximumFractionDigits: decimals
      }).format(value)
    }
    
    return showSymbol ? `${currencySymbol.value}${formatted}` : formatted
  }
  
  /**
   * 格式化数字（带千分位）
   * @param {number} value - 数值
   * @param {number} decimals - 小数位数
   */
  const formatNumber = (value, decimals = 2) => {
    // 隐私模式下遮盖数字
    if (themeStore.privacyMode) {
      return PRIVACY_MASK
    }

    if (value === null || value === undefined || isNaN(value)) {
      return '0'
    }
    
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: decimals,
      maximumFractionDigits: decimals
    }).format(value)
  }
  
  /**
   * 格式化百分比
   * @param {number} value - 数值（已经是百分比形式）
   * @param {boolean} showSign - 是否显示正负号
   */
  const formatPercent = (value, showSign = true) => {
    if (value === null || value === undefined || isNaN(value)) {
      return '0.00%'
    }
    
    const sign = showSign && value > 0 ? '+' : ''
    return `${sign}${value.toFixed(2)}%`
  }
  
  /**
   * 格式化日期
   * @param {Date|string|number} date - 日期
   * @param {string} format - 格式类型
   */
  const formatDate = (date, format = 'short') => {
    const d = new Date(date)
    
    if (isNaN(d.getTime())) {
      return '-'
    }
    
    const options = {
      short: { month: 'short', day: 'numeric' },
      medium: { month: 'short', day: 'numeric', year: 'numeric' },
      long: { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' },
      time: { hour: '2-digit', minute: '2-digit' },
      datetime: { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }
    }
    
    return new Intl.DateTimeFormat('zh-CN', options[format] || options.short).format(d)
  }
  
  /**
   * 格式化相对时间
   * @param {Date|string|number} date - 日期
   */
  const formatRelativeTime = (date) => {
    const d = new Date(date)
    const now = new Date()
    const diff = now - d
    
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)
    
    if (minutes < 1) return '刚刚'
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    if (days < 7) return `${days}天前`
    
    return formatDate(date, 'short')
  }
  
  /**
   * 缩短地址显示
   * @param {string} address - 钱包地址
   * @param {number} chars - 前后保留字符数
   */
  const shortenAddress = (address, chars = 4) => {
    if (!address) return ''
    return `${address.slice(0, chars + 2)}...${address.slice(-chars)}`
  }

  /**
   * 设置定价货币
   * @param {string} code - 货币代码
   */
  const setPricingCurrency = (code) => {
    if (currencies.find(c => c.code === code)) {
      currentCurrency.value = code
    }
  }

  // 可用的定价货币列表
  const availablePricingCurrencies = computed(() =>
    currencies.map(c => c.code)
  )

  return {
    currentCurrency,
    pricingCurrency: currentCurrency, // 别名，用于计价货币场景
    currencySymbol,
    currencies,
    availablePricingCurrencies,
    setPricingCurrency,
    formatCurrency,
    formatNumber,
    formatPercent,
    formatDate,
    formatRelativeTime,
    shortenAddress
  }
}
