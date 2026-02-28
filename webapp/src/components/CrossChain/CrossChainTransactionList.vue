<template>
  <div class="cross-chain-transactions">
    <div class="transactions-header">
      <h3>跨链交易记录</h3>
      <div class="filters">
        <select v-model="filters.status" @change="loadTransactions">
          <option value="">全部状态</option>
          <option value="pending">进行中</option>
          <option value="confirmed">已确认</option>
          <option value="completed">已完成</option>
          <option value="failed">失败</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="transactions.length === 0" class="empty">
      <p>暂无跨链交易记录</p>
    </div>

    <div v-else class="transactions-list">
      <div v-for="tx in transactions" :key="tx.id" class="transaction-item">
        <div class="tx-header">
          <span class="tx-hash">{{ formatTxHash(tx.tx_hash) }}</span>
          <span :class="['status', tx.status]">{{ formatStatus(tx.status) }}</span>
        </div>

        <div class="tx-body">
          <div class="tx-route">
            <div class="chain-info">
              <span class="chain-name">{{ tx.source_chain }}</span>
              <span class="amount">{{ tx.source_amount }} {{ tx.source_token }}</span>
            </div>
            <div class="arrow">→</div>
            <div class="chain-info">
              <span class="chain-name">{{ tx.dest_chain }}</span>
              <span class="amount">{{ tx.dest_amount }} {{ tx.dest_token }}</span>
            </div>
          </div>

          <div class="tx-details">
            <div class="detail-item">
              <span class="label">跨链桥:</span>
              <span class="value">{{ tx.bridge_protocol }}</span>
            </div>
            <div class="detail-item">
              <span class="label">手续费:</span>
              <span class="value">${{ tx.total_fee_usd?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="detail-item">
              <span class="label">发起时间:</span>
              <span class="value">{{ formatTime(tx.initiated_at) }}</span>
            </div>
            <div v-if="tx.completed_at" class="detail-item">
              <span class="label">完成时间:</span>
              <span class="value">{{ formatTime(tx.completed_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="total > pageSize" class="pagination">
      <button @click="prevPage" :disabled="page === 1">上一页</button>
      <span>第 {{ page }} 页 / 共 {{ totalPages }} 页</span>
      <button @click="nextPage" :disabled="page >= totalPages">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { crossChainApi } from '@/api/cross-chain'

const transactions = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filters = ref({
  status: ''
})

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const loadTransactions = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filters.value.status) {
      params.status = filters.value.status
    }

    const res = await crossChainApi.getTransactions(params)
    transactions.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('加载跨链交易失败:', error)
  } finally {
    loading.value = false
  }
}

const prevPage = () => {
  if (page.value > 1) {
    page.value--
    loadTransactions()
  }
}

const nextPage = () => {
  if (page.value < totalPages.value) {
    page.value++
    loadTransactions()
  }
}

const formatTxHash = (hash) => {
  if (!hash) return ''
  return `${hash.slice(0, 10)}...${hash.slice(-8)}`
}

const formatStatus = (status) => {
  const statusMap = {
    pending: '进行中',
    confirmed: '已确认',
    completed: '已完成',
    failed: '失败'
  }
  return statusMap[status] || status
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadTransactions()
})
</script>

<style scoped>
.cross-chain-transactions {
  padding: 20px;
}

.transactions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.transactions-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.filters select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #666;
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.transaction-item {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  transition: box-shadow 0.2s;
}

.transaction-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tx-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.tx-hash {
  font-family: monospace;
  font-size: 14px;
  color: #666;
}

.status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status.pending {
  background: #fef3c7;
  color: #92400e;
}

.status.confirmed {
  background: #dbeafe;
  color: #1e40af;
}

.status.completed {
  background: #d1fae5;
  color: #065f46;
}

.status.failed {
  background: #fee2e2;
  color: #991b1b;
}

.tx-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tx-route {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 6px;
}

.chain-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.chain-name {
  font-weight: 600;
  color: #111827;
  text-transform: capitalize;
}

.amount {
  font-size: 14px;
  color: #6b7280;
}

.arrow {
  font-size: 20px;
  color: #9ca3af;
}

.tx-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.detail-item {
  display: flex;
  gap: 8px;
  font-size: 14px;
}

.detail-item .label {
  color: #6b7280;
}

.detail-item .value {
  color: #111827;
  font-weight: 500;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination span {
  color: #6b7280;
  font-size: 14px;
}
</style>
