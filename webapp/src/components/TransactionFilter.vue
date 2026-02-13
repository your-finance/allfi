<script setup>
/**
 * 交易记录筛选器
 * 提供类型、来源、时间范围和搜索筛选
 */
import { ref, watch } from 'vue'
import { PhMagnifyingGlass, PhFunnel } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  filters: { type: Object, required: true }
})

const emit = defineEmits(['update'])
const { t } = useI18n()

// 本地筛选状态
const localFilters = ref({ ...props.filters })
const searchInput = ref(props.filters.search || '')
let searchTimer = null

// 交易类型选项
const typeOptions = [
  { id: 'all', labelKey: 'transaction.filterAll' },
  { id: 'buy', labelKey: 'transaction.typeBuy' },
  { id: 'sell', labelKey: 'transaction.typeSell' },
  { id: 'swap', labelKey: 'transaction.typeSwap' },
  { id: 'transfer', labelKey: 'transaction.typeTransfer' },
  { id: 'deposit', labelKey: 'transaction.typeDeposit' },
  { id: 'withdraw', labelKey: 'transaction.typeWithdraw' },
]

// 来源类型选项
const sourceOptions = [
  { id: 'all', labelKey: 'transaction.sourceAll' },
  { id: 'cex', labelKey: 'transaction.sourceCex' },
  { id: 'blockchain', labelKey: 'transaction.sourceChain' },
  { id: 'manual', labelKey: 'transaction.sourceManual' },
]

// 时间范围选项
const timeOptions = [
  { id: 'all', labelKey: 'transaction.timeAll' },
  { id: '7', labelKey: 'transaction.time7d' },
  { id: '30', labelKey: 'transaction.time30d' },
  { id: '90', labelKey: 'transaction.time90d' },
]

// 类型筛选
const setType = (type) => {
  localFilters.value.type = type
  emit('update', { ...localFilters.value })
}

// 来源筛选
const setSource = (sourceType) => {
  localFilters.value.sourceType = sourceType
  emit('update', { ...localFilters.value })
}

// 时间范围筛选
const setTimeRange = (range) => {
  localFilters.value.timeRange = range
  emit('update', { ...localFilters.value })
}

// 搜索（防抖）
watch(searchInput, (val) => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    localFilters.value.search = val
    emit('update', { ...localFilters.value })
  }, 300)
})

// 外部 filters 变化时同步
watch(() => props.filters, (val) => {
  localFilters.value = { ...val }
  searchInput.value = val.search || ''
}, { deep: true })
</script>

<template>
  <div class="tx-filter">
    <!-- 搜索框 -->
    <div class="filter-search">
      <PhMagnifyingGlass :size="14" class="search-icon" />
      <input
        type="text"
        v-model="searchInput"
        :placeholder="t('transaction.searchPlaceholder')"
        class="search-input"
      />
    </div>

    <!-- 筛选标签栏 -->
    <div class="filter-rows">
      <!-- 交易类型 -->
      <div class="filter-row">
        <PhFunnel :size="12" class="filter-icon" />
        <div class="filter-pills">
          <button
            v-for="opt in typeOptions"
            :key="opt.id"
            class="filter-pill"
            :class="{ active: localFilters.type === opt.id }"
            @click="setType(opt.id)"
          >
            {{ t(opt.labelKey) }}
          </button>
        </div>
      </div>

      <!-- 来源 + 时间 -->
      <div class="filter-row">
        <div class="filter-pills">
          <button
            v-for="opt in sourceOptions"
            :key="opt.id"
            class="filter-pill source-pill"
            :class="{ active: localFilters.sourceType === opt.id }"
            @click="setSource(opt.id)"
          >
            {{ t(opt.labelKey) }}
          </button>
        </div>
        <span class="filter-divider">|</span>
        <div class="filter-pills">
          <button
            v-for="opt in timeOptions"
            :key="opt.id"
            class="filter-pill time-pill"
            :class="{ active: localFilters.timeRange === opt.id }"
            @click="setTimeRange(opt.id)"
          >
            {{ t(opt.labelKey) }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tx-filter {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

/* 搜索 */
.filter-search {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: 6px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  max-width: 320px;
}

.search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.search-input {
  background: none;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-size: 0.8125rem;
  width: 100%;
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

/* 筛选标签 */
.filter-rows {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.filter-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  flex-wrap: wrap;
}

.filter-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.filter-pills {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.filter-pill {
  padding: 3px 10px;
  font-size: 0.6875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.filter-pill:hover {
  border-color: var(--color-accent-primary);
  color: var(--color-text-primary);
}

.filter-pill.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

.filter-divider {
  color: var(--color-border);
  font-size: 0.75rem;
}
</style>
