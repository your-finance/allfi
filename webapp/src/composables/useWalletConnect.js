/**
 * Web3 钱包连接 Composable
 * 支持 MetaMask / WalletConnect 连接，自动获取地址和链信息
 * Mock 模式下模拟钱包连接，不依赖真实 Web3 库
 */
import { ref, computed } from 'vue'

// 是否使用 Mock 模式
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== 'false'

// 模拟延迟
const delay = (ms = 1000) => new Promise(resolve => setTimeout(resolve, ms))

// 全局连接状态（单例）
const _isConnected = ref(false)
const _currentAddress = ref(null)
const _currentChainId = ref(null)
const _isConnecting = ref(false)
const _error = ref(null)

// 链 ID 映射
const CHAIN_MAP = {
  1: { id: 'ETH', name: 'Ethereum' },
  56: { id: 'BSC', name: 'BNB Chain' },
  137: { id: 'MATIC', name: 'Polygon' },
  42161: { id: 'ARB', name: 'Arbitrum' },
  10: { id: 'OP', name: 'Optimism' },
}

/**
 * Mock 钱包地址生成
 * @returns {string} 模拟的以太坊地址
 */
function generateMockAddress() {
  const chars = '0123456789abcdef'
  let addr = '0x'
  for (let i = 0; i < 40; i++) {
    addr += chars[Math.floor(Math.random() * chars.length)]
  }
  return addr
}

export function useWalletConnect() {
  const isConnected = computed(() => _isConnected.value)
  const currentAddress = computed(() => _currentAddress.value)
  const currentChainId = computed(() => _currentChainId.value)
  const isConnecting = computed(() => _isConnecting.value)
  const error = computed(() => _error.value)

  // 当前链信息
  const currentChain = computed(() => {
    if (!_currentChainId.value) return null
    return CHAIN_MAP[_currentChainId.value] || { id: 'UNKNOWN', name: `Chain ${_currentChainId.value}` }
  })

  /**
   * 连接钱包
   * @param {'metamask'|'walletconnect'} method - 连接方式
   * @returns {Promise<{ address: string, chainId: number }>}
   */
  async function connectWallet(method = 'metamask') {
    _isConnecting.value = true
    _error.value = null

    try {
      if (USE_MOCK) {
        // Mock 模式：模拟连接过程
        await delay(1500)
        const mockAddress = generateMockAddress()
        const mockChainId = 1 // 默认以太坊主网
        _currentAddress.value = mockAddress
        _currentChainId.value = mockChainId
        _isConnected.value = true
        return { address: mockAddress, chainId: mockChainId }
      }

      // 真实模式：调用 MetaMask 或 WalletConnect
      if (method === 'metamask') {
        if (!window.ethereum) {
          throw new Error('未检测到 MetaMask，请安装 MetaMask 浏览器扩展')
        }
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' })
        const chainId = await window.ethereum.request({ method: 'eth_chainId' })
        const address = accounts[0]
        const parsedChainId = parseInt(chainId, 16)

        _currentAddress.value = address
        _currentChainId.value = parsedChainId
        _isConnected.value = true
        return { address, chainId: parsedChainId }
      }

      // WalletConnect 需要 wagmi 库（未安装时回退提示）
      throw new Error('WalletConnect 需要安装 wagmi 依赖，请使用 MetaMask 或手动输入地址')
    } catch (err) {
      _error.value = err.message || '钱包连接失败'
      throw err
    } finally {
      _isConnecting.value = false
    }
  }

  /**
   * 断开钱包连接
   */
  function disconnectWallet() {
    _currentAddress.value = null
    _currentChainId.value = null
    _isConnected.value = false
    _error.value = null
  }

  /**
   * 清除错误
   */
  function clearError() {
    _error.value = null
  }

  return {
    isConnected,
    currentAddress,
    currentChainId,
    currentChain,
    isConnecting,
    error,
    connectWallet,
    disconnectWallet,
    clearError,
    CHAIN_MAP,
  }
}
