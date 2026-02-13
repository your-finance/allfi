<script setup>
/**
 * 交易记录时间线组件
 * 按日期分组的时间线布局，支持无限滚动加载更多
 */
import { onMounted, ref } from 'vue'
import { PhSpinnerGap } from '@phosphor-icons/vue'
import { useTransactionStore } from '../stores/transactionStore'
import { useI18n } from '../composables/useI18n'
import TransactionItem from './TransactionItem.vue'
import TransactionFilter from './TransactionFilter.vue'

const txStore = useTransactionStore()
const { t } = useI18n()

// 无限滚动观察目标
const sentinel = ref(null)

// 格式化日期标题
const formatDateLabel = (dateStr) => {
  const d = new Date(dateStr)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  if (dateStr === formatDateKey(today)) return t('transaction.today')
  if (dateStr === formatDateKey(yesterday)) return t('transaction.yesterday')

  const month = d.getMonth() + 1
  const day = d.getDate()
  return `${month}月${day}日`
}

const formatDateKey = (d) => {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 筛选更新
const handleFilterUpdate = (newFilters) => {
  txStore.updateFilters(newFilters)
}

// 加载更多
const handleLoadMore = () => {
  if (txStore.hasMore && !txStore.isLoading) {
    txStore.loadMore()
  }
}

// 初始化加载
onMounted(() => {
  txStore.fetchTransactions(1)

  // IntersectionObserver 实现无限滚动
  if (sentinel.value) {
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        handleLoadMore()
      }
    }, { rootMargin: '100px' })
    observer.observe(sentinel.value)
  }
})
</script>

<template>
  <div class="tx-timeline">
    <!-- 筛选器 -->
    <TransactionFilter
      :filters="txStore.filters"
      @update="handleFilterUpdate"
    />

    <!-- 总条数 -->
    <div class="tx-summary">
      <span class="tx-total font-mono">{{ txStore.totalCount }} {{ t('transaction.records') }}</span>
    </div>

    <!-- 无数据 -->
    <div v-if="!txStore.isLoading && txStore.transactions.length === 0" class="empty-state">
      <p>{{ t('transaction.noRecords') }}</p>
    </div>

    <!-- 时间线 -->
    <div v-else class="timeline-list">
      <div
        v-for="group in txStore.groupedByDate"
        :key="group.date"
        class="timeline-group"
      >
        <!-- 日期标题 -->
        <div class="date-header">
          <div class="date-dot"></div>
          <span class="date-label">{{ formatDateLabel(group.date) }}</span>
          <span class="date-count font-mono">{{ group.transactions.length }} {{ t('transaction.items') }}</span>
        </div>

        <!-- 交易列表 -->
        <div class="tx-list">
          <TransactionItem
            v-for="tx in group.transactions"
            :key="tx.id"
            :tx="tx"
          />
        </div>
      </div>

      <!-- 加载更多/加载中 -->
      <div ref="sentinel" class="sentinel">
        <div v-if="txStore.isLoading" class="loading-more">
          <PhSpinnerGap :size="20" class="spinner" />
          <span>{{ t('common.loading') }}</span>
        </div>
        <button
          v-else-if="txStore.hasMore"
          class="btn btn-ghost load-more-btn"
          @click="handleLoadMore"
        >
          {{ t('transaction.loadMore') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tx-timeline {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* 摘要 */
.tx-summary {
  display: flex;
  align-items: center;
}

.tx-total {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* 无数据 */
.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 时间线列表 */
.timeline-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* 日期分组 */
.timeline-group {
  display: flex;
  flex-direction: column;
}

.date-header {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
  margin-bottom: var(--gap-sm);
}

.date-dot {
  width: 6px;
  height: 6px;
  background: var(--color-accent-primary);
  border-radius: 50%;
  flex-shrink: 0;
}

.date-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.date-count {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 交易列表 */
.tx-list {
  display: flex;
  flex-direction: column;
}

/* 加载更多 */
.sentinel {
  display: flex;
  justify-content: center;
  padding: var(--gap-md);
}

.loading-more {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  color: var(--color-text-muted);
  font-size: 0.75rem;
}

.spinner {
  animation: spin 1s linear infinite;
  color: var(--color-accent-primary);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.load-more-btn {
  font-size: 0.75rem;
}
</style>
