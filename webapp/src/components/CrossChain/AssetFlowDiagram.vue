<template>
  <div class="asset-flow-diagram">
    <div class="diagram-header">
      <h3>资产流向分析</h3>
      <div class="time-range">
        <input type="date" v-model="startDate" @change="loadFlowData" />
        <span>至</span>
        <input type="date" v-model="endDate" @change="loadFlowData" />
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <div v-else-if="!hasData" class="empty">
      <p>暂无流向数据</p>
    </div>

    <div v-else ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import * as echarts from 'echarts'
import { crossChainApi } from '@/api/cross-chain'

const chartRef = ref(null)
const loading = ref(false)
const startDate = ref('')
const endDate = ref('')
const flowData = ref({ nodes: [], links: [] })
let chartInstance = null

const hasData = computed(() => {
  return flowData.value.nodes.length > 0 && flowData.value.links.length > 0
})

const loadFlowData = async () => {
  loading.value = true
  try {
    const params = {}
    if (startDate.value) params.start_time = startDate.value
    if (endDate.value) params.end_time = endDate.value

    const res = await crossChainApi.getAssetFlow(params)
    flowData.value = res

    if (hasData.value) {
      renderChart()
    }
  } catch (error) {
    console.error('加载资产流向失败:', error)
  } finally {
    loading.value = false
  }
}

const renderChart = () => {
  if (!chartRef.value) return

  if (!chartInstance) {
    chartInstance = echarts.init(chartRef.value)
  }

  // 转换数据格式为 ECharts Sankey 所需格式
  const nodes = flowData.value.nodes.map(node => ({
    name: formatChainName(node.name)
  }))

  const links = flowData.value.links.map(link => ({
    source: formatChainName(link.source),
    target: formatChainName(link.target),
    value: link.value
  }))

  const option = {
    title: {
      text: '跨链资产流向',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 600
      }
    },
    tooltip: {
      trigger: 'item',
      triggerOn: 'mousemove',
      formatter: (params) => {
        if (params.dataType === 'edge') {
          return `${params.data.source} → ${params.data.target}<br/>金额: $${params.value.toFixed(2)}`
        } else {
          return `${params.name}<br/>交易次数: ${params.value}`
        }
      }
    },
    series: [
      {
        type: 'sankey',
        layout: 'none',
        emphasis: {
          focus: 'adjacency'
        },
        data: nodes,
        links: links,
        lineStyle: {
          color: 'gradient',
          curveness: 0.5
        },
        label: {
          fontSize: 12,
          fontWeight: 500
        },
        itemStyle: {
          borderWidth: 1,
          borderColor: '#fff'
        }
      }
    ]
  }

  chartInstance.setOption(option)
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
  if (chartInstance) {
    chartInstance.resize()
  }
}

onMounted(() => {
  // 设置默认时间范围（最近30天）
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 30)

  endDate.value = end.toISOString().split('T')[0]
  startDate.value = start.toISOString().split('T')[0]

  loadFlowData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.asset-flow-diagram {
  padding: 20px;
}

.diagram-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.diagram-header h3 {
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

.chart-container {
  width: 100%;
  height: 500px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}
</style>
