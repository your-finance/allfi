<script setup>
/**
 * 盈亏归因分析面板
 * 展示各资产对整体收益的贡献拆解（价格效应 / 数量效应）
 */
import { ref, computed, onMounted, watch } from 'vue'
import { PhTrendUp, PhTrendDown, PhEquals, PhArrowsLeftRight } from '@phosphor-icons/vue'
import { analyticsService } from '../api/index.js'
import { useFormatters } from '../composables/useFormatters'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

const { formatNumber, currencySymbol } = useFormatters()
const themeStore = useThemeStore()
const { t } = useI18n()

// 时间范围
const selectedRange = ref('7D')
const rangeOptions = ['7D', '30D', '90D']

// 数据
const data = ref(null)
const isLoading = ref(false)
const hasError = ref(false)

const rangeToApi = { '7D': '7d', '30D': '30d', '90D': '90d' }

const loadData = async () => {
  isLoading.value = true
  hasError.value = false
  try {
    data.value = await analyticsService.getAttribution(rangeToApi[selectedRange.value])
  } catch (err) {
    console.error('加载归因数据失败:', err)
    hasError.value = true
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)
watch(selectedRange, loadData)

const colors = computed(() => themeStore.currentTheme.colors)

// 汇总
const totalReturn = computed(() => data.value?.total_return || 0)
const totalPercent = computed(() => data.value?.total_percent || 0)
const attributions = computed(() => data.value?.attributions || [])

// 按贡献排序
const sortedAttributions = computed(() => {
  return [...attributions.value].sort((a, b) => Math.abs(b.contribution) - Math.abs(a.contribution))
})

// 贡献占比（基于绝对贡献和）
const totalAbsContribution = computed(() =>
  attributions.value.reduce((sum, a) => sum + Math.abs(a.contribution), 0) || 1
)

const getContributionWidth = (contribution) => {
  return Math.abs(contribution) / totalAbsContribution.value * 100
}
</script>

<template>
  <section class="attribution-panel">
    <!-- 标题栏 + 时间范围选择 -->
    <div class="panel-header">
      <h3 class="panel-title">
        <PhArrowsLeftRight :size="16" weight="bold" />
        {{ t('analytics.attribution.title') }}
      </h3>
      <div class="range-selector">
        <button
          v-for="r in rangeOptions"
          :key="r"
          class="range-btn"
          :class="{ active: selectedRange === r }"
          @click="selectedRange = r"
        >
          {{ r }}
        </button>
      </div>
    </div>

    <!-- 加载中 -->
    <div v-if="isLoading" class="loading-state">
      {{ t('common.loading') }}
    </div>

    <!-- 错误 -->
    <div v-else-if="hasError" class="empty-state">
      <p>{{ t('analytics.attribution.error') }}</p>
    </div>

    <!-- 无数据 -->
    <div v-else-if="!data || attributions.length === 0" class="empty-state">
      <p>{{ t('analytics.attribution.noData') }}</p>
    </div>

    <!-- 主体 -->
    <template v-else>
      <!-- 汇总卡片 -->
      <div class="summary-row">
        <div class="summary-card">
          <span class="summary-label">{{ t('analytics.attribution.totalReturn') }}</span>
          <span class="summary-value font-mono" :class="totalReturn >= 0 ? 'positive' : 'negative'">
            {{ totalReturn >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(Math.abs(totalReturn), 2) }}
          </span>
          <span class="summary-pct font-mono" :class="totalPercent >= 0 ? 'positive' : 'negative'">
            {{ totalPercent >= 0 ? '+' : '' }}{{ totalPercent.toFixed(2) }}%
          </span>
        </div>
      </div>

      <!-- 归因表格 -->
      <div class="table-wrapper">
        <table class="attr-table">
          <thead>
            <tr>
              <th class="col-asset">{{ t('analytics.attribution.asset') }}</th>
              <th class="col-num right">{{ t('analytics.attribution.weight') }}</th>
              <th class="col-num right">{{ t('analytics.attribution.returnPct') }}</th>
              <th class="col-contribution">{{ t('analytics.attribution.contribution') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in sortedAttributions" :key="a.symbol">
              <td class="col-asset">
                <span class="asset-symbol">{{ a.symbol }}</span>
                <span v-if="a.source" class="asset-source">{{ a.source }}</span>
              </td>
              <td class="col-num right font-mono">{{ a.weight }}%</td>
              <td class="col-num right font-mono" :class="a.return >= 0 ? 'positive' : 'negative'">
                {{ a.return >= 0 ? '+' : '' }}{{ a.return.toFixed(2) }}%
              </td>
              <td class="col-contribution">
                <div class="contribution-cell">
                  <span class="contribution-value font-mono" :class="a.contribution >= 0 ? 'positive' : 'negative'">
                    {{ a.contribution >= 0 ? '+' : '' }}{{ currencySymbol }}{{ formatNumber(Math.abs(a.contribution), 0) }}
                  </span>
                  <div class="contribution-bar-bg">
                    <div
                      class="contribution-bar-fill"
                      :class="a.contribution >= 0 ? 'bar-positive' : 'bar-negative'"
                      :style="{ width: getContributionWidth(a.contribution) + '%' }"
                    />
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 底部提示 -->
      <div class="attribution-hint">
        <PhEquals :size="12" />
        <span>{{ t('analytics.attribution.hint') }}</span>
      </div>
    </template>
  </section>
</template>

<style scoped>
.attribution-panel {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
  padding: var(--gap-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

/* 标题栏 */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.range-selector {
  display: flex;
  gap: 4px;
}

.range-btn {
  padding: 3px 10px;
  font-size: 0.6875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.range-btn:hover {
  border-color: var(--color-accent-primary);
}

.range-btn.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

/* 汇总卡片 */
.summary-row {
  display: flex;
  gap: var(--gap-md);
}

.summary-card {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--gap-md) var(--gap-lg);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.summary-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.summary-value {
  font-size: 1.25rem;
  font-weight: 600;
}

.summary-pct {
  font-size: 0.75rem;
  font-weight: 500;
}

/* 表格 */
.table-wrapper {
  overflow-x: auto;
}

.attr-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.75rem;
}

.attr-table th {
  text-align: left;
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap;
}

.attr-table td {
  padding: var(--gap-sm) var(--gap-md);
  border-bottom: 1px solid color-mix(in srgb, var(--color-border) 50%, transparent);
  color: var(--color-text-primary);
}

.attr-table tbody tr:hover {
  background: var(--color-bg-tertiary);
}

.right {
  text-align: right;
}

.col-asset {
  min-width: 100px;
}

.col-num {
  min-width: 60px;
}

.col-contribution {
  min-width: 160px;
}

.asset-symbol {
  font-weight: 600;
  color: var(--color-text-primary);
}

.asset-source {
  margin-left: 6px;
  font-size: 0.625rem;
  color: var(--color-text-muted);
  padding: 1px 4px;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
}

/* 贡献条 */
.contribution-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.contribution-value {
  font-size: 0.75rem;
  font-weight: 500;
}

.contribution-bar-bg {
  height: 4px;
  background: var(--color-bg-tertiary);
  border-radius: 2px;
  width: 100%;
}

.contribution-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.4s ease;
}

.bar-positive {
  background: var(--color-success);
}

.bar-negative {
  background: var(--color-error);
}

/* 状态 */
.positive {
  color: var(--color-success);
}

.negative {
  color: var(--color-error);
}

/* 底部提示 */
.attribution-hint {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

/* 空状态 / 加载 */
.loading-state,
.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .summary-row {
    flex-direction: column;
  }
}
</style>
