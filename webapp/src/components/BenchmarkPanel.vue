<script setup>
/**
 * 基准对比面板
 * 展示用户收益率与 BTC/ETH 等基准指数的对比
 * 适用于自托管单用户场景——与市场比，不与他人比
 */
import { ref, onMounted } from 'vue'
import { PhChartLineUp, PhArrowUp, PhArrowDown, PhEquals } from '@phosphor-icons/vue'
import { benchmarkService } from '../api/benchmarkService.js'
import { useI18n } from '../composables/useI18n'

const { t } = useI18n()

const benchmarkData = ref(null)
const isLoading = ref(false)
const selectedPeriod = ref('30d')
const periods = ['7d', '30d', '90d', '1y']

// 加载基准对比数据
const loadBenchmark = async () => {
  isLoading.value = true
  try {
    benchmarkData.value = await benchmarkService.getBenchmark(selectedPeriod.value)
  } catch (e) {
    console.error('加载基准对比失败:', e)
  } finally {
    isLoading.value = false
  }
}

// 切换周期
const changePeriod = (period) => {
  selectedPeriod.value = period
  loadBenchmark()
}

// 格式化收益率
const formatReturn = (val) => {
  if (val === 0 || val === undefined || val === null) return '0.0%'
  const sign = val >= 0 ? '+' : ''
  return `${sign}${val.toFixed(1)}%`
}

onMounted(loadBenchmark)
</script>

<template>
  <section class="benchmark-panel">
    <!-- 标题栏 -->
    <div class="panel-header">
      <h3 class="panel-title">
        <PhChartLineUp :size="16" />
        {{ t('benchmark.title') }}
      </h3>
      <div class="period-selector">
        <button
          v-for="p in periods"
          :key="p"
          class="period-btn"
          :class="{ active: selectedPeriod === p }"
          @click="changePeriod(p)"
        >
          {{ t(`benchmark.period_${p}`) }}
        </button>
      </div>
    </div>

    <!-- 加载中 -->
    <div v-if="isLoading" class="loading-state">{{ t('common.loading') }}</div>

    <template v-else-if="benchmarkData">
      <!-- 用户收益率 -->
      <div class="user-return-card">
        <span class="return-label">{{ t('benchmark.myReturn') }}</span>
        <span class="return-value font-mono" :class="benchmarkData.user_return >= 0 ? 'positive' : 'negative'">
          {{ formatReturn(benchmarkData.user_return) }}
        </span>
      </div>

      <!-- 基准指数对比 -->
      <div class="benchmark-list">
        <div
          v-for="idx in benchmarkData.benchmarks"
          :key="idx.name"
          class="benchmark-item"
        >
          <div class="benchmark-info">
            <span class="benchmark-name">{{ idx.name }}</span>
            <span class="benchmark-return font-mono" :class="idx.return_percent >= 0 ? 'positive' : 'negative'">
              {{ formatReturn(idx.return_percent) }}
            </span>
          </div>
          <div class="outperform-row">
            <component
              :is="idx.user_outperform > 0 ? PhArrowUp : idx.user_outperform < 0 ? PhArrowDown : PhEquals"
              :size="12"
              :class="idx.user_outperform > 0 ? 'positive' : idx.user_outperform < 0 ? 'negative' : ''"
            />
            <span
              class="outperform-value font-mono"
              :class="idx.user_outperform > 0 ? 'positive' : idx.user_outperform < 0 ? 'negative' : ''"
            >
              {{ idx.user_outperform > 0 ? t('benchmark.outperform') : idx.user_outperform < 0 ? t('benchmark.underperform') : t('benchmark.equal') }}
              {{ Math.abs(idx.user_outperform).toFixed(1) }}%
            </span>
          </div>
        </div>
      </div>

      <!-- 提示 -->
      <div class="benchmark-hint">
        {{ t('benchmark.hint') }}
      </div>
    </template>

    <!-- 无数据 -->
    <div v-else class="empty-state">
      {{ t('benchmark.noData') }}
    </div>
  </section>
</template>

<style scoped>
.benchmark-panel {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: var(--gap-sm);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.period-selector {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.period-btn {
  padding: 4px 10px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.period-btn:hover {
  color: var(--color-text-primary);
}

.period-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

.loading-state,
.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 用户收益率卡片 */
.user-return-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-md) var(--gap-lg);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.return-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.return-value {
  font-size: 1.25rem;
  font-weight: 600;
}

.return-value.positive { color: var(--color-success); }
.return-value.negative { color: var(--color-error); }

/* 基准指数列表 */
.benchmark-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.benchmark-item {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.benchmark-item:hover {
  background: var(--color-bg-tertiary);
}

.benchmark-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.benchmark-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.benchmark-return {
  font-size: 0.8125rem;
  font-weight: 600;
}

.benchmark-return.positive { color: var(--color-success); }
.benchmark-return.negative { color: var(--color-error); }

.outperform-row {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.outperform-value {
  font-size: 0.6875rem;
  font-weight: 500;
}

.outperform-value.positive { color: var(--color-success); }
.outperform-value.negative { color: var(--color-error); }

.benchmark-hint {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  padding-top: var(--gap-sm);
  border-top: 1px solid var(--color-border);
}
</style>
