<template>
  <div class="bridge-comparison">
    <div class="comparison-header">
      <h3>跨链桥对比</h3>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="bridges.length === 0" class="empty">
      <p>暂无跨链桥数据</p>
    </div>

    <div v-else class="comparison-table">
      <table>
        <thead>
          <tr>
            <th>跨链桥</th>
            <th>支持链数量</th>
            <th>费率</th>
            <th>平均完成时间</th>
            <th>支持的链</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="bridge in bridges" :key="bridge.protocol">
            <td class="bridge-name">
              <div class="name-cell">
                <span class="protocol-badge">{{ bridge.protocol }}</span>
                <span class="full-name">{{ bridge.name }}</span>
              </div>
            </td>
            <td class="text-center">{{ bridge.support_chains?.length || 0 }}</td>
            <td class="text-center">
              <span class="fee-rate">{{ (bridge.fee_rate * 100).toFixed(2) }}%</span>
            </td>
            <td class="text-center">
              <span class="avg-time">{{ formatTime(bridge.avg_time) }}</span>
            </td>
            <td>
              <div class="chains-list">
                <span
                  v-for="chain in bridge.support_chains"
                  :key="chain"
                  class="chain-tag"
                >
                  {{ formatChainName(chain) }}
                </span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="comparison-notes">
      <h4>说明</h4>
      <ul>
        <li>费率为跨链桥收取的手续费百分比，不包括 Gas 费用</li>
        <li>平均完成时间为统计数据，实际时间可能因网络拥堵而变化</li>
        <li>建议根据实际需求选择合适的跨链桥</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { crossChainApi } from '@/api/cross-chain'

const bridges = ref([])
const loading = ref(false)

const loadBridges = async () => {
  loading.value = true
  try {
    const res = await crossChainApi.getBridges()
    bridges.value = res.bridges || []
  } catch (error) {
    console.error('加载跨链桥列表失败:', error)
  } finally {
    loading.value = false
  }
}

const formatTime = (seconds) => {
  if (!seconds) return '-'
  if (seconds < 60) return `${seconds}秒`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const remainMinutes = minutes % 60
  return remainMinutes > 0 ? `${hours}小时${remainMinutes}分钟` : `${hours}小时`
}

const formatChainName = (chain) => {
  const chainNames = {
    ethereum: 'Ethereum',
    bsc: 'BSC',
    polygon: 'Polygon',
    arbitrum: 'Arbitrum',
    optimism: 'Optimism',
    avalanche: 'Avalanche',
    fantom: 'Fantom',
    base: 'Base',
    gnosis: 'Gnosis'
  }
  return chainNames[chain] || chain
}

onMounted(() => {
  loadBridges()
})
</script>

<style scoped>
.bridge-comparison {
  padding: 20px;
}

.comparison-header {
  margin-bottom: 20px;
}

.comparison-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: #666;
}

.comparison-table {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: #f9fafb;
}

th {
  padding: 12px 16px;
  text-align: left;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

td {
  padding: 16px;
  font-size: 14px;
  color: #111827;
  border-bottom: 1px solid #f3f4f6;
}

tbody tr:last-child td {
  border-bottom: none;
}

tbody tr:hover {
  background: #f9fafb;
}

.text-center {
  text-align: center;
}

.bridge-name {
  min-width: 200px;
}

.name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.protocol-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  width: fit-content;
}

.full-name {
  font-size: 13px;
  color: #6b7280;
}

.fee-rate {
  display: inline-block;
  padding: 4px 12px;
  background: #fef3c7;
  color: #92400e;
  border-radius: 12px;
  font-weight: 500;
}

.avg-time {
  color: #059669;
  font-weight: 500;
}

.chains-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.chain-tag {
  display: inline-block;
  padding: 4px 10px;
  background: #f3f4f6;
  color: #374151;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.comparison-notes {
  margin-top: 24px;
  padding: 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 8px;
}

.comparison-notes h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #92400e;
}

.comparison-notes ul {
  margin: 0;
  padding-left: 20px;
}

.comparison-notes li {
  font-size: 13px;
  color: #78350f;
  line-height: 1.6;
}
</style>
