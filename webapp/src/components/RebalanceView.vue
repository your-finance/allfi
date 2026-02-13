<script setup>
/**
 * 再平衡分析视图
 * 显示当前配比 vs 目标配比对比柱状图，偏离高亮，建议操作
 */
import { ref, computed, onMounted } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js'
import { PhWarning } from '@phosphor-icons/vue'
import { useStrategyStore } from '../stores/strategyStore'
import { useThemeStore } from '../stores/themeStore'
import { useI18n } from '../composables/useI18n'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const props = defineProps({
  strategy: { type: Object, required: true }
})

const strategyStore = useStrategyStore()
const themeStore = useThemeStore()
const { t } = useI18n()

const suggestions = ref([])
const colors = computed(() => themeStore.currentTheme.colors)
const targets = computed(() => props.strategy.config.targets || [])
const threshold = computed(() => props.strategy.config.deviationThreshold || 5)

// 加载再平衡建议
onMounted(async () => {
  const result = await strategyStore.getRebalanceSuggestion(props.strategy.id)
  if (result) suggestions.value = result.suggestions
})

// 柱状图数据
const chartData = computed(() => ({
  labels: targets.value.map(t => t.symbol),
  datasets: [
    {
      label: t('strategy.currentAllocation'),
      data: targets.value.map(t => t.currentPct),
      backgroundColor: colors.value.accentSecondary,
      borderRadius: 3,
      barPercentage: 0.4,
    },
    {
      label: t('strategy.targetAllocation'),
      data: targets.value.map(t => t.targetPct),
      backgroundColor: `${colors.value.accentPrimary}66`,
      borderColor: colors.value.accentPrimary,
      borderWidth: 1,
      borderRadius: 3,
      barPercentage: 0.4,
    },
  ]
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      align: 'end',
      labels: {
        color: colors.value.textSecondary,
        usePointStyle: true,
        pointStyle: 'rect',
        padding: 12,
        font: { size: 10 },
      }
    },
    tooltip: {
      backgroundColor: colors.value.bgElevated,
      titleColor: colors.value.textPrimary,
      bodyColor: colors.value.textSecondary,
      borderColor: colors.value.border,
      borderWidth: 1,
      callbacks: {
        label: (ctx) => `${ctx.dataset.label}: ${ctx.raw}%`
      }
    }
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: colors.value.textMuted, font: { size: 11 } }
    },
    y: {
      max: 100,
      grid: { color: colors.value.border, lineWidth: 0.5 },
      ticks: {
        color: colors.value.textSecondary,
        font: { size: 10 },
        callback: (v) => `${v}%`
      }
    }
  }
}))
</script>

<template>
  <div class="rebalance-view">
    <!-- 对比图 -->
    <div class="chart-wrapper">
      <Bar :data="chartData" :options="chartOptions" />
    </div>

    <!-- 偏离明细 -->
    <div class="deviation-list">
      <div
        v-for="tgt in targets"
        :key="tgt.symbol"
        class="deviation-item"
        :class="{ warning: Math.abs(tgt.currentPct - tgt.targetPct) > threshold }"
      >
        <span class="dev-symbol">{{ tgt.symbol }}</span>
        <span class="dev-current font-mono">{{ tgt.currentPct }}%</span>
        <span class="dev-arrow">&rarr;</span>
        <span class="dev-target font-mono">{{ tgt.targetPct }}%</span>
        <span class="dev-diff font-mono" :class="tgt.currentPct > tgt.targetPct ? 'negative' : 'positive'">
          {{ tgt.currentPct > tgt.targetPct ? '+' : '' }}{{ (tgt.currentPct - tgt.targetPct).toFixed(1) }}%
        </span>
      </div>
    </div>

    <!-- 建议操作 -->
    <div v-if="suggestions.length > 0" class="suggestions">
      <h4 class="suggestion-title">{{ t('strategy.suggestedActions') }}</h4>
      <ul class="suggestion-list">
        <li v-for="s in suggestions" :key="s.symbol" class="suggestion-item">
          <span :class="s.action === 'sell' ? 'negative' : 'positive'">
            {{ s.description }}
          </span>
        </li>
      </ul>
    </div>

    <!-- 免责声明 -->
    <div class="disclaimer">
      <PhWarning :size="12" />
      <span>{{ t('strategy.disclaimer') }}</span>
    </div>
  </div>
</template>

<style scoped>
.rebalance-view {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.chart-wrapper {
  height: 200px;
}

/* 偏离明细 */
.deviation-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.deviation-item {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-xs) var(--gap-md);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
}

.deviation-item.warning {
  background: rgba(245, 158, 11, 0.06);
}

.dev-symbol {
  font-weight: 600;
  color: var(--color-text-primary);
  width: 50px;
}

.dev-current, .dev-target {
  color: var(--color-text-secondary);
}

.dev-arrow {
  color: var(--color-text-muted);
  font-size: 0.6875rem;
}

.dev-diff {
  font-weight: 600;
  margin-left: auto;
}

.dev-diff.positive { color: var(--color-success); }
.dev-diff.negative { color: var(--color-error); }

/* 建议 */
.suggestion-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-xs);
}

.suggestion-list {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.suggestion-item {
  font-size: 0.75rem;
  padding: var(--gap-xs) var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
}

.suggestion-item .positive { color: var(--color-success); }
.suggestion-item .negative { color: var(--color-error); }

/* 免责 */
.disclaimer {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.6875rem;
  color: var(--color-warning);
  padding: var(--gap-xs) var(--gap-md);
  background: rgba(245, 158, 11, 0.06);
  border-radius: var(--radius-sm);
}
</style>
