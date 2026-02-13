/**
 * 账户管理 Store
 * 管理 CEX 账户、链上钱包、手动资产的状态和操作
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { cexService, walletService, manualAssetService } from '../api/index.js'

export const useAccountStore = defineStore('account', () => {
  // ========== 状态 ==========

  // CEX 账户
  const cexAccounts = ref([])
  const cexLoading = ref(false)
  const cexError = ref(null)

  // 链上钱包
  const walletAddresses = ref([])
  const walletLoading = ref(false)
  const walletError = ref(null)

  // 手动资产
  const manualAssets = ref([])
  const manualLoading = ref(false)
  const manualError = ref(null)

  // 全局操作状态
  const isRefreshing = ref(false)
  const isDeleting = ref(false)
  const isAdding = ref(false)

  // ========== 计算属性 ==========

  // CEX 总价值
  const cexTotalValue = computed(() => {
    return cexAccounts.value.reduce((sum, acc) => sum + (acc.balance || 0), 0)
  })

  // 链上总价值
  const walletTotalValue = computed(() => {
    return walletAddresses.value.reduce((sum, w) => sum + (w.balance || 0), 0)
  })

  // 法币→USD 简易汇率（后续对接实时汇率 API）
  const fiatToUsdRates = {
    USD: 1, CNY: 0.14, HKD: 0.128, EUR: 1.08,
    GBP: 1.27, JPY: 0.0067, SGD: 0.75
  }

  // 手动资产总价值（多币种转换为 USD）
  const manualTotalValue = computed(() => {
    return manualAssets.value.reduce((sum, a) => {
      const rate = fiatToUsdRates[a.currency] || 1
      return sum + (a.balance || 0) * rate
    }, 0)
  })

  // 总资产价值
  const totalValue = computed(() => {
    return cexTotalValue.value + walletTotalValue.value + manualTotalValue.value
  })

  // 是否有任何加载中
  const isLoading = computed(() => {
    return cexLoading.value || walletLoading.value || manualLoading.value
  })

  // ========== CEX 账户操作 ==========

  /**
   * 加载所有 CEX 账户
   */
  async function loadCexAccounts() {
    cexLoading.value = true
    cexError.value = null

    try {
      cexAccounts.value = await cexService.getAccounts()
    } catch (err) {
      cexError.value = err.message || '加载 CEX 账户失败'
      console.error('加载 CEX 账户失败:', err)
    } finally {
      cexLoading.value = false
    }
  }

  /**
   * 添加 CEX 账户
   * @param {object} accountData - 账户数据
   */
  async function addCexAccount(accountData) {
    isAdding.value = true
    cexError.value = null

    try {
      const newAccount = await cexService.addAccount(accountData)
      cexAccounts.value.push(newAccount)
      return newAccount
    } catch (err) {
      cexError.value = err.message || '添加 CEX 账户失败'
      throw err
    } finally {
      isAdding.value = false
    }
  }

  /**
   * 更新 CEX 账户
   * @param {number} id - 账户ID
   * @param {object} updates - 更新数据
   */
  async function updateCexAccount(id, updates) {
    cexError.value = null

    try {
      const updated = await cexService.updateAccount(id, updates)
      const index = cexAccounts.value.findIndex(a => a.id === id)
      if (index !== -1) {
        cexAccounts.value[index] = updated
      }
      return updated
    } catch (err) {
      cexError.value = err.message || '更新 CEX 账户失败'
      throw err
    }
  }

  /**
   * 删除 CEX 账户
   * @param {number} id - 账户ID
   */
  async function deleteCexAccount(id) {
    isDeleting.value = true
    cexError.value = null

    try {
      await cexService.deleteAccount(id)
      cexAccounts.value = cexAccounts.value.filter(a => a.id !== id)
      return true
    } catch (err) {
      cexError.value = err.message || '删除 CEX 账户失败'
      throw err
    } finally {
      isDeleting.value = false
    }
  }

  /**
   * 刷新 CEX 账户余额
   * @param {number} id - 账户ID
   */
  async function refreshCexAccount(id) {
    isRefreshing.value = true
    cexError.value = null

    try {
      const updated = await cexService.refreshAccount(id)
      const index = cexAccounts.value.findIndex(a => a.id === id)
      if (index !== -1) {
        cexAccounts.value[index] = updated
      }
      return updated
    } catch (err) {
      cexError.value = err.message || '刷新 CEX 账户失败'
      throw err
    } finally {
      isRefreshing.value = false
    }
  }

  // ========== 链上钱包操作 ==========

  /**
   * 加载所有钱包地址
   */
  async function loadWalletAddresses() {
    walletLoading.value = true
    walletError.value = null

    try {
      walletAddresses.value = await walletService.getAddresses()
    } catch (err) {
      walletError.value = err.message || '加载钱包地址失败'
      console.error('加载钱包地址失败:', err)
    } finally {
      walletLoading.value = false
    }
  }

  /**
   * 添加钱包地址
   * @param {object} walletData - 钱包数据
   */
  async function addWalletAddress(walletData) {
    isAdding.value = true
    walletError.value = null

    try {
      const newWallet = await walletService.addAddress(walletData)
      walletAddresses.value.push(newWallet)
      return newWallet
    } catch (err) {
      walletError.value = err.message || '添加钱包地址失败'
      throw err
    } finally {
      isAdding.value = false
    }
  }

  /**
   * 更新钱包地址
   * @param {number} id - 钱包ID
   * @param {object} updates - 更新数据
   */
  async function updateWalletAddress(id, updates) {
    walletError.value = null

    try {
      const updated = await walletService.updateAddress(id, updates)
      const index = walletAddresses.value.findIndex(w => w.id === id)
      if (index !== -1) {
        walletAddresses.value[index] = updated
      }
      return updated
    } catch (err) {
      walletError.value = err.message || '更新钱包地址失败'
      throw err
    }
  }

  /**
   * 删除钱包地址
   * @param {number} id - 钱包ID
   */
  async function deleteWalletAddress(id) {
    isDeleting.value = true
    walletError.value = null

    try {
      await walletService.deleteAddress(id)
      walletAddresses.value = walletAddresses.value.filter(w => w.id !== id)
      return true
    } catch (err) {
      walletError.value = err.message || '删除钱包地址失败'
      throw err
    } finally {
      isDeleting.value = false
    }
  }

  /**
   * 刷新钱包余额
   * @param {number} id - 钱包ID
   */
  async function refreshWalletAddress(id) {
    isRefreshing.value = true
    walletError.value = null

    try {
      const updated = await walletService.refreshAddress(id)
      const index = walletAddresses.value.findIndex(w => w.id === id)
      if (index !== -1) {
        walletAddresses.value[index] = updated
      }
      return updated
    } catch (err) {
      walletError.value = err.message || '刷新钱包地址失败'
      throw err
    } finally {
      isRefreshing.value = false
    }
  }

  // ========== 手动资产操作 ==========

  /**
   * 加载所有手动资产
   */
  async function loadManualAssets() {
    manualLoading.value = true
    manualError.value = null

    try {
      manualAssets.value = await manualAssetService.getAssets()
    } catch (err) {
      manualError.value = err.message || '加载手动资产失败'
      console.error('加载手动资产失败:', err)
    } finally {
      manualLoading.value = false
    }
  }

  /**
   * 添加手动资产
   * @param {object} assetData - 资产数据
   */
  async function addManualAsset(assetData) {
    isAdding.value = true
    manualError.value = null

    try {
      const newAsset = await manualAssetService.addAsset(assetData)
      manualAssets.value.push(newAsset)
      return newAsset
    } catch (err) {
      manualError.value = err.message || '添加手动资产失败'
      throw err
    } finally {
      isAdding.value = false
    }
  }

  /**
   * 更新手动资产
   * @param {number} id - 资产ID
   * @param {object} updates - 更新数据
   */
  async function updateManualAsset(id, updates) {
    manualError.value = null

    try {
      const updated = await manualAssetService.updateAsset(id, updates)
      const index = manualAssets.value.findIndex(a => a.id === id)
      if (index !== -1) {
        manualAssets.value[index] = updated
      }
      return updated
    } catch (err) {
      manualError.value = err.message || '更新手动资产失败'
      throw err
    }
  }

  /**
   * 删除手动资产
   * @param {number} id - 资产ID
   */
  async function deleteManualAsset(id) {
    isDeleting.value = true
    manualError.value = null

    try {
      await manualAssetService.deleteAsset(id)
      manualAssets.value = manualAssets.value.filter(a => a.id !== id)
      return true
    } catch (err) {
      manualError.value = err.message || '删除手动资产失败'
      throw err
    } finally {
      isDeleting.value = false
    }
  }

  // ========== 通用操作 ==========

  /**
   * 加载所有账户数据
   */
  async function loadAll() {
    await Promise.all([
      loadCexAccounts(),
      loadWalletAddresses(),
      loadManualAssets()
    ])
  }

  /**
   * 根据类型删除账户
   * @param {string} type - 类型 (cex, blockchain, manual)
   * @param {number} id - 账户ID
   */
  async function deleteAccount(type, id) {
    switch (type) {
      case 'cex':
        return deleteCexAccount(id)
      case 'blockchain':
        return deleteWalletAddress(id)
      case 'manual':
        return deleteManualAsset(id)
      default:
        throw new Error('未知的账户类型')
    }
  }

  /**
   * 根据类型刷新账户
   * @param {string} type - 类型 (cex, blockchain)
   * @param {number} id - 账户ID
   */
  async function refreshAccount(type, id) {
    switch (type) {
      case 'cex':
        return refreshCexAccount(id)
      case 'blockchain':
        return refreshWalletAddress(id)
      default:
        throw new Error('该类型账户不支持刷新')
    }
  }

  /**
   * 清除所有错误
   */
  function clearErrors() {
    cexError.value = null
    walletError.value = null
    manualError.value = null
  }

  return {
    // 状态
    cexAccounts,
    walletAddresses,
    manualAssets,
    cexLoading,
    walletLoading,
    manualLoading,
    cexError,
    walletError,
    manualError,
    isRefreshing,
    isDeleting,
    isAdding,

    // 计算属性
    cexTotalValue,
    walletTotalValue,
    manualTotalValue,
    totalValue,
    isLoading,

    // CEX 操作
    loadCexAccounts,
    addCexAccount,
    updateCexAccount,
    deleteCexAccount,
    refreshCexAccount,

    // 钱包操作
    loadWalletAddresses,
    addWalletAddress,
    updateWalletAddress,
    deleteWalletAddress,
    refreshWalletAddress,

    // 手动资产操作
    loadManualAssets,
    addManualAsset,
    updateManualAsset,
    deleteManualAsset,

    // 通用操作
    loadAll,
    deleteAccount,
    refreshAccount,
    clearErrors
  }
})
