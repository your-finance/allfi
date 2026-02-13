<script setup>
/**
 * 批量导入钱包地址对话框
 * 支持逐行粘贴地址（最多 100 个）
 */
import { ref, computed } from 'vue'
import { PhX, PhUploadSimple, PhCheckCircle, PhWarningCircle } from '@phosphor-icons/vue'
import { walletService } from '../api/index.js'
import { useI18n } from '../composables/useI18n'

const { t } = useI18n()

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'imported'])

// 状态
const addressText = ref('')
const blockchain = ref('ethereum')
const importing = ref(false)
const result = ref(null)

// 支持的区块链
const blockchains = [
  { value: 'ethereum', label: 'Ethereum' },
  { value: 'bsc', label: 'BNB Chain' },
  { value: 'polygon', label: 'Polygon' }
]

// 解析地址列表
const parsedAddresses = computed(() => {
  if (!addressText.value.trim()) return []
  return addressText.value
    .split('\n')
    .map(line => line.trim())
    .filter(line => line.length > 0)
})

// 地址数量
const addressCount = computed(() => parsedAddresses.value.length)

// 是否超出限制
const isOverLimit = computed(() => addressCount.value > 100)

// 执行批量导入
const handleImport = async () => {
  if (addressCount.value === 0 || isOverLimit.value) return

  importing.value = true
  result.value = null

  try {
    result.value = await walletService.batchImport(parsedAddresses.value, blockchain.value)
    if (result.value.imported > 0) {
      emit('imported')
    }
  } catch (e) {
    console.error('批量导入失败:', e)
    result.value = { imported: 0, failed: addressCount.value, error: e.message }
  } finally {
    importing.value = false
  }
}

// 重置状态
const handleClose = () => {
  addressText.value = ''
  blockchain.value = 'ethereum'
  result.value = null
  importing.value = false
  emit('close')
}

// 重新开始
const resetResult = () => {
  result.value = null
  addressText.value = ''
}
</script>

<template>
  <Transition name="fade">
    <div v-if="visible" class="dialog-overlay" @click.self="handleClose">
      <div class="dialog">
        <!-- 标题栏 -->
        <div class="dialog-header">
          <h3>{{ t('batchImport.title') }}</h3>
          <button class="close-btn" @click="handleClose">
            <PhX :size="16" />
          </button>
        </div>

        <!-- 导入结果 -->
        <div v-if="result" class="dialog-body">
          <div class="result-summary">
            <div class="result-stat success">
              <PhCheckCircle :size="20" weight="bold" />
              <span class="font-mono">{{ result.imported }}</span>
              <span>{{ t('batchImport.successCount') }}</span>
            </div>
            <div class="result-stat failed" v-if="result.failed > 0">
              <PhWarningCircle :size="20" weight="bold" />
              <span class="font-mono">{{ result.failed }}</span>
              <span>{{ t('batchImport.failedCount') }}</span>
            </div>
          </div>

          <!-- 错误信息 -->
          <div v-if="result.error" class="failed-list">
            <div class="failed-header">{{ t('batchImport.failedDetails') }}</div>
            <div class="failed-item">
              <span class="failed-reason">{{ result.error }}</span>
            </div>
          </div>

          <div class="dialog-actions">
            <button class="btn btn-secondary" @click="resetResult">
              {{ t('batchImport.importMore') }}
            </button>
            <button class="btn btn-primary" @click="handleClose">
              {{ t('common.confirm') }}
            </button>
          </div>
        </div>

        <!-- 导入表单 -->
        <div v-else class="dialog-body">
          <p class="dialog-desc">{{ t('batchImport.description') }}</p>

          <!-- 区块链选择 -->
          <div class="form-group">
            <label>{{ t('batchImport.blockchain') }}</label>
            <select v-model="blockchain" class="form-select">
              <option v-for="chain in blockchains" :key="chain.value" :value="chain.value">
                {{ chain.label }}
              </option>
            </select>
          </div>

          <!-- 地址输入 -->
          <div class="form-group">
            <label>
              {{ t('batchImport.addresses') }}
              <span class="count-badge" :class="{ 'over-limit': isOverLimit }">
                {{ addressCount }} / 100
              </span>
            </label>
            <textarea
              v-model="addressText"
              class="form-textarea font-mono"
              :placeholder="t('batchImport.placeholder')"
              rows="10"
            />
          </div>

          <div class="dialog-actions">
            <button class="btn btn-secondary" @click="handleClose">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn-primary"
              :disabled="addressCount === 0 || isOverLimit || importing"
              @click="handleImport"
            >
              <PhUploadSimple :size="14" />
              {{ importing ? t('batchImport.importing') : t('batchImport.import') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  width: 520px;
  max-width: 90vw;
  max-height: 85vh;
  overflow-y: auto;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.dialog-header h3 {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
}

.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dialog-body {
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.dialog-desc {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.form-group label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.count-badge {
  font-size: 0.6875rem;
  padding: 1px 6px;
  border-radius: var(--radius-xs);
  background: var(--color-bg-tertiary);
  color: var(--color-text-muted);
}

.count-badge.over-limit {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-error);
}

.form-select {
  padding: 8px 12px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.8125rem;
}

.form-textarea {
  padding: 10px 12px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.75rem;
  line-height: 1.6;
  resize: vertical;
}

.form-textarea::placeholder {
  color: var(--color-text-muted);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  padding-top: var(--gap-sm);
}

.btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  font-size: 0.8125rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-primary {
  background: var(--color-accent-primary);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--color-bg-tertiary);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

/* 结果 */
.result-summary {
  display: flex;
  gap: var(--gap-lg);
  padding: var(--gap-md);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
}

.result-stat {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  font-size: 0.8125rem;
}

.result-stat.success {
  color: var(--color-success);
}

.result-stat.failed {
  color: var(--color-error);
}

.result-stat .font-mono {
  font-weight: 600;
  font-size: 1rem;
}

.failed-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
}

.failed-header {
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-tertiary);
  border-bottom: 1px solid var(--color-border);
}

.failed-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap-xs) var(--gap-md);
  font-size: 0.75rem;
  border-bottom: 1px solid var(--color-border);
}

.failed-item:last-child {
  border-bottom: none;
}

.failed-addr {
  color: var(--color-text-primary);
  max-width: 240px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.failed-reason {
  color: var(--color-error);
  font-size: 0.6875rem;
}

/* 过渡 */
.fade-enter-active, .fade-leave-active {
  transition: opacity var(--transition-fast);
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
