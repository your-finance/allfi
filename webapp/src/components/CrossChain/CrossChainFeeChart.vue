<template>
  <div class="cross-chain-fee-chart">
    <div class="chart-header">
      <h3>跨链手续费统计</h3>
      <div class="time-range">
        <input type="date" v-model="startDate" @change="loadFeeStats" />
        <span>至</span>
        <input type="date" v-model="endDate" @change="loadFeeStats" />
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="!hasData" class="empty">
      <p>暂无手续费数据</p>
    </div>

    <div v-else class="stats-container">
      <div class="summary-cards">
        <div class="card">
          <div class="card-label">总手续费</div>
          <div class="card-value">${{ stats.total_fee?.toFixed(2) || '0.00' }}</div>
        </div>
        <div class="card">
          <div class="card-label">平均手续费</div>
          <div class="card-value">${{ stats.avg_fee?.toFixed(2) || '0.00' }}</div>
        </div>
        <div class="card">
          <div class="card-label">交易次数</div>
          <div class="card-value">{{ stats.transaction_count || 0 }}</div>
        </div>
      </div>

      <div class="charts-grid">
        <div class="chart-wrapper">
          <h4>按跨链桥统计</h4>
          <div ref="bridgeChartRef" class="chart"></div>
        </div>
        <div class="chart-wrapper">
          <h4>按链统计</h4>
          <div ref="chainChartRef" class="chart"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import * as echarts from 'echarts'
import { crossChainApi } from '@/api/cross-chain'

const bridgeChartRef = ref(null)
const chainChartRef = ref(null)
const loading = ref(false)
const startDate = ref('')
const endDate = ref('')
const stats = ref({
  total_fee: 0,
  avg_fee: 0,
  fee_by_bridge: {},
  fee_by_chain: {},
  transaction_count: 0
})

let bridgeChartInstance = null
let chainChartInstance = null

const hasData = computed(() => {
  return stats.value.transaction_count > 0
})

const loadFeeStats = async () => {
  loading.value = true
  try {
    const params = {}
    if (startDate.value) params.start_time = startDate.value
    if (endDate.value) params.end_time = endDate.value

    const res = await crossChainApi.getFeeStats(params)
    stats.value = res

    if (hasData.value) {
      renderCharts()
    }
  } catch (error) {
    console.error('加载手续费统计失败:', error)
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  renderBridgeChart()
  renderChainChart()
}

const renderBridgeChart = () => {
  if (!bridgeChartRef.value) return

  if (!bridgeChartInstance) {
    bridgeChartInstance = echarts.init(bridgeChartRef.value)
  }

  const data = Object.entries(stats.value.fee_by_bridge || {}).map(([name, value]) => ({
    name,
    value
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ${c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data
      }
    ]
  }

  bridgeChartInstance.setOption(option)
}

const renderChainChart = () => {
  if (!chainChartRef.value) return

  if (!chainChartInstance) {
    chainChartInstance = echarts.init(chainChartRef.value)
  }

  const data = Object.entries(stats.value.fee_by_chain || {})
    .map(([name, value]) => ({
      name: formatChainName(name),
      value
    }))
    .sort((a, b) => b.value - a.value)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: '{b}: ${c}'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLabel: {
        rotate: 45,
        fontSize: 12
      }
    },
    yAxis: {
      type: 'value',
      name: '手续费 (USD)',
      axisLabel: {
        formatter: '${value}'
      }
    },
    series: [
      {
        type: 'bar',
        data: data.map(item => item.value),
        itemStyle: {
          borderRadius: [4, 4, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#2378f7' },
              { offset: 0.7, color: '#2378f7' },
              { offset: 1, color: '#83bff6' }
            ])
          }
        }
      }
    ]
  }

  chainChartInstance.setOption(option)
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
    base: 'Base'
  }
  return chainNames[chain] || chain
}

const handleResize = () => {
  if (bridgeChartInstance) bridgeChartInstance.resize()
  if (chainChartInstance) chainChartInstance.resize()
}

onMounted(() => {
  // 设置默认时间范围（最近30天）
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 30)

  endDate.value = end.toISOString().split('T')[0]
  startDate.value = start.toISOString().split('T')[0]

  loadFeeStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  if (bridgeChartInstance) {
    bridgeChartInstance.dispose()
    bridgeChartInstance = null
  }
  if (chainChartInstance) {
    chainChartInstance.dispose()
    chainChartInstance = null
  }
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.cross-chain-fee-chart {
  padding: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.chart-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.time-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-range input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.loading, .empty {
  text-align: center;
  padding: 60px;
  color: #666;
}

.stats-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}

.card-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.chart-wrapper {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
}

.chart-wrapper h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.chart {
  width: 100%;
  height: 300px;
}
</style>
