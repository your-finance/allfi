/**
 * 交易记录状态管理 Store
 * 管理统一交易时间线的数据、筛选条件和分页
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { transactionService } from '../api/transactionService.js'

export const useTransactionStore = defineStore('transaction', () => {
  // 状态
  const transactions = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const currentPage = ref(1)
  const totalPages = ref(1)
  const totalCount = ref(0)

  // 筛选条件
  const filters = ref({
    type: 'all',        // all / buy / sell / transfer / swap / deposit / withdraw
    sourceType: 'all',  // all / cex / blockchain / manual
    timeRange: 'all',   // all / 7 / 30 / 90
    search: '',
  })

  // 按日期分组的交易记录
  const groupedByDate = computed(() => {
    const groups = {}
    for (const tx of transactions.value) {
      const date = new Date(tx.timestamp)
      const dateKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
      if (!groups[dateKey]) {
        groups[dateKey] = { date: dateKey, transactions: [] }
      }
      groups[dateKey].transactions.push(tx)
    }
    // 按日期降序排列
    return Object.values(groups).sort((a, b) => b.date.localeCompare(a.date))
  })

  /**
   * 获取交易记录
   */
  async function fetchTransactions(page = 1) {
    isLoading.value = true
    error.value = null

    try {
      const result = await transactionService.getTransactions({
        ...filters.value,
        page,
        pageSize: 20,
      })
      if (page === 1) {
        transactions.value = result.items
      } else {
        // 追加加载（无限滚动）
        transactions.value = [...transactions.value, ...result.items]
      }
      currentPage.value = result.page
      totalPages.value = result.totalPages
      totalCount.value = result.total
    } catch (err) {
      error.value = err.message || '加载交易记录失败'
      console.error('加载交易记录失败:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 加载更多
   */
  async function loadMore() {
    if (currentPage.value < totalPages.value && !isLoading.value) {
      await fetchTransactions(currentPage.value + 1)
    }
  }

  /**
   * 是否还有更多数据
   */
  const hasMore = computed(() => currentPage.value < totalPages.value)

  /**
   * 更新筛选条件并重新加载
   */
  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
    fetchTransactions(1)
  }

  /**
   * 重置
   */
  function reset() {
    transactions.value = []
    error.value = null
    currentPage.value = 1
    totalPages.value = 1
    totalCount.value = 0
    filters.value = { type: 'all', sourceType: 'all', timeRange: 'all', search: '' }
  }

  return {
    transactions,
    isLoading,
    error,
    currentPage,
    totalPages,
    totalCount,
    filters,
    groupedByDate,
    hasMore,
    fetchTransactions,
    loadMore,
    updateFilters,
    reset,
  }
})
