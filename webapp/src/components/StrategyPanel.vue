<script setup>
/**
 * 策略面板组件
 * 展示策略列表、状态开关、创建入口、再平衡详情
 */
import { ref, onMounted } from 'vue'
import {
  PhPlus,
  PhArrowsClockwise,
  PhCalendarBlank,
  PhBell,
  PhTrash,
  PhPause,
  PhPlay
} from '@phosphor-icons/vue'
import { useStrategyStore } from '../stores/strategyStore'
import { useFormatters } from '../composables/useFormatters'
import { useI18n } from '../composables/useI18n'
import AddStrategyDialog from './AddStrategyDialog.vue'
import RebalanceView from './RebalanceView.vue'

const strategyStore = useStrategyStore()
const { formatNumber } = useFormatters()
const { t } = useI18n()

const showAddDialog = ref(false)
const expandedStrategy = ref(null)

// 策略类型图标
const typeIcon = {
  rebalance: PhArrowsClockwise,
  dca: PhCalendarBlank,
  alert: PhBell,
}

// 策略类型颜色
const typeColor = {
  rebalance: '#3B82F6',
  dca: '#10B981',
  alert: '#F59E0B',
}

// 状态配置
const statusConfig = {
  active: { labelKey: 'strategy.statusActive', color: '#10B981' },
  paused: { labelKey: 'strategy.statusPaused', color: '#6B7280' },
  triggered: { labelKey: 'strategy.statusTriggered', color: '#F59E0B' },
}

// 展开/收起策略详情
const toggleExpand = (id) => {
  expandedStrategy.value = expandedStrategy.value === id ? null : id
}

// 处理策略创建
const handleCreated = (strategy) => {
  strategyStore.addStrategy(strategy)
}

onMounted(() => {
  if (strategyStore.strategies.length === 0) {
    strategyStore.fetchStrategies()
  }
})
</script>

<template>
  <section class="strategy-panel">
    <!-- 标题栏 -->
    <div class="panel-header">
      <h3 class="panel-title">{{ t('strategy.title') }}</h3>
      <button class="btn btn-ghost btn-sm" @click="showAddDialog = true">
        <PhPlus :size="14" />
        {{ t('strategy.addStrategy') }}
      </button>
    </div>

    <!-- 无策略 -->
    <div v-if="strategyStore.strategies.length === 0 && !strategyStore.isLoading" class="empty-state">
      <p>{{ t('strategy.noStrategies') }}</p>
    </div>

    <!-- 策略列表 -->
    <div v-else class="strategy-list">
      <div
        v-for="s in strategyStore.strategies"
        :key="s.id"
        class="strategy-card"
      >
        <div class="card-main" @click="toggleExpand(s.id)">
          <!-- 类型图标 -->
          <div class="stg-icon" :style="{ background: typeColor[s.type] + '18', color: typeColor[s.type] }">
            <component :is="typeIcon[s.type]" :size="16" weight="bold" />
          </div>

          <!-- 策略信息 -->
          <div class="stg-info">
            <div class="stg-name">{{ s.name }}</div>
            <div class="stg-meta">
              <span class="stg-type">{{ t(`strategy.type${s.type.charAt(0).toUpperCase() + s.type.slice(1)}`) }}</span>
              <span v-if="s.lastTriggeredAt" class="stg-last-trigger">
                {{ t('strategy.lastTriggered') }}: {{ s.lastTriggeredAt }}
              </span>
            </div>
          </div>

          <!-- 状态标签 -->
          <span class="status-badge" :style="{ color: statusConfig[s.status].color }">
            {{ t(statusConfig[s.status].labelKey) }}
          </span>

          <!-- 操作按钮 -->
          <div class="stg-actions">
            <button
              class="action-btn"
              :title="s.status === 'active' ? t('strategy.pause') : t('strategy.resume')"
              @click.stop="strategyStore.toggleStrategy(s.id)"
            >
              <PhPause v-if="s.status === 'active'" :size="14" />
              <PhPlay v-else :size="14" />
            </button>
            <button
              class="action-btn danger"
              :title="t('common.delete')"
              @click.stop="strategyStore.deleteStrategy(s.id)"
            >
              <PhTrash :size="14" />
            </button>
          </div>
        </div>

        <!-- 再平衡详情（展开） -->
        <div v-if="expandedStrategy === s.id && s.type === 'rebalance'" class="card-expand">
          <RebalanceView :strategy="s" />
        </div>

        <!-- 定投详情 -->
        <div v-if="expandedStrategy === s.id && s.type === 'dca'" class="card-expand">
          <div class="dca-detail">
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.dcaSymbol') }}</span>
              <span class="detail-value font-mono">{{ s.config.symbol }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.dcaAmount') }}</span>
              <span class="detail-value font-mono">${{ formatNumber(s.config.amount, 0) }} / {{ t(`strategy.${s.config.frequency}`) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.totalInvested') }}</span>
              <span class="detail-value font-mono">${{ formatNumber(s.config.totalInvested, 0) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.avgPrice') }}</span>
              <span class="detail-value font-mono">${{ formatNumber(s.config.avgPrice, 2) }}</span>
            </div>
          </div>
        </div>

        <!-- 止盈止损详情 -->
        <div v-if="expandedStrategy === s.id && s.type === 'alert'" class="card-expand">
          <div class="alert-detail">
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.alertSymbol') }}</span>
              <span class="detail-value font-mono">{{ s.config.symbol }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.alertDirection') }}</span>
              <span class="detail-value">{{ s.config.direction === 'above' ? t('strategy.priceAbove') : t('strategy.priceBelow') }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('strategy.targetPrice') }}</span>
              <span class="detail-value font-mono">${{ formatNumber(s.config.targetPrice, 0) }}</span>
            </div>
            <div v-if="s.config.note" class="detail-row">
              <span class="detail-label">{{ t('strategy.note') }}</span>
              <span class="detail-value">{{ s.config.note }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加策略对话框 -->
    <AddStrategyDialog
      :visible="showAddDialog"
      @close="showAddDialog = false"
      @created="handleCreated"
    />
  </section>
</template>

<style scoped>
.strategy-panel {
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
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.btn-sm {
  font-size: 0.6875rem;
  padding: 4px 8px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-state {
  padding: var(--gap-xl);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.8125rem;
}

/* 策略列表 */
.strategy-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.strategy-card {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.card-main {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.card-main:hover {
  background: var(--color-bg-tertiary);
}

.stg-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stg-info {
  flex: 1;
  min-width: 0;
}

.stg-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.stg-meta {
  display: flex;
  gap: var(--gap-md);
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

.status-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  flex-shrink: 0;
}

.stg-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.action-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius-xs);
  transition: all var(--transition-fast);
}

.action-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-tertiary);
}

.action-btn.danger:hover {
  color: var(--color-error);
}

/* 展开详情 */
.card-expand {
  padding: var(--gap-md);
  border-top: 1px solid var(--color-border);
  background: var(--color-bg-tertiary);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  font-size: 0.75rem;
}

.detail-label {
  color: var(--color-text-muted);
}

.detail-value {
  color: var(--color-text-primary);
  font-weight: 500;
}
</style>
