<script setup>
/**
 * 价格预警对话框
 * 支持创建、查看、暂停/恢复、删除价格预警
 */
import { ref, computed, onMounted } from 'vue'
import {
  PhX,
  PhBell,
  PhPlus,
  PhTrash,
  PhPause,
  PhPlay,
  PhArrowUp,
  PhArrowDown,
  PhCheck
} from '@phosphor-icons/vue'
import { priceAlertService } from '../api/index.js'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close'])

const { t } = useI18n()

// 状态
const alerts = ref([])
const loading = ref(false)
const showForm = ref(false)
const submitting = ref(false)

// 新预警表单
const form = ref({
  symbol: 'BTC',
  condition: 'above',
  target_price: '',
  note: ''
})

// 常见币种列表
const commonSymbols = ['BTC', 'ETH', 'BNB', 'SOL', 'DOGE', 'ADA', 'XRP', 'DOT']

// 活跃预警数
const activeCount = computed(() => alerts.value.filter(a => a.is_active && !a.triggered).length)

// 加载预警列表
const loadAlerts = async () => {
  loading.value = true
  try {
    alerts.value = await priceAlertService.getAlerts()
  } catch (e) {
    console.error('加载预警失败:', e)
  } finally {
    loading.value = false
  }
}

// 创建预警
const createAlert = async () => {
  if (!form.value.target_price || form.value.target_price <= 0) return

  submitting.value = true
  try {
    const newAlert = await priceAlertService.createAlert({
      symbol: form.value.symbol,
      condition: form.value.condition,
      target_price: parseFloat(form.value.target_price),
      note: form.value.note
    })
    alerts.value.unshift(newAlert)
    resetForm()
    showForm.value = false
  } catch (e) {
    console.error('创建预警失败:', e)
  } finally {
    submitting.value = false
  }
}

// 切换预警状态（暂停/恢复）
const toggleAlert = async (alert) => {
  try {
    await priceAlertService.updateAlert(alert.id, { is_active: !alert.is_active })
    alert.is_active = !alert.is_active
  } catch (e) {
    console.error('更新预警失败:', e)
  }
}

// 删除预警
const deleteAlert = async (alert) => {
  try {
    await priceAlertService.deleteAlert(alert.id)
    alerts.value = alerts.value.filter(a => a.id !== alert.id)
  } catch (e) {
    console.error('删除预警失败:', e)
  }
}

// 重置表单
const resetForm = () => {
  form.value = { symbol: 'BTC', condition: 'above', target_price: '', note: '' }
}

// 格式化时间
const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  if (props.visible) {
    loadAlerts()
  }
})

// 监听 visible 变化
import { watch } from 'vue'
watch(() => props.visible, (val) => {
  if (val) loadAlerts()
})
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="visible" class="dialog-overlay" @click.self="emit('close')">
        <div class="dialog-panel">
          <!-- 标题栏 -->
          <div class="dialog-header">
            <div class="header-left">
              <PhBell :size="18" weight="bold" />
              <h3>{{ t('priceAlert.title') }}</h3>
              <span v-if="activeCount > 0" class="active-badge">{{ activeCount }}</span>
            </div>
            <button class="close-btn" @click="emit('close')">
              <PhX :size="16" />
            </button>
          </div>

          <!-- 内容区域 -->
          <div class="dialog-body">
            <!-- 新建预警按钮 -->
            <button
              v-if="!showForm"
              class="add-btn"
              @click="showForm = true"
            >
              <PhPlus :size="14" />
              {{ t('priceAlert.addAlert') }}
            </button>

            <!-- 新建预警表单 -->
            <div v-if="showForm" class="alert-form">
              <div class="form-row">
                <label>{{ t('priceAlert.symbol') }}</label>
                <div class="symbol-selector">
                  <button
                    v-for="sym in commonSymbols"
                    :key="sym"
                    class="symbol-chip"
                    :class="{ active: form.symbol === sym }"
                    @click="form.symbol = sym"
                  >
                    {{ sym }}
                  </button>
                </div>
              </div>

              <div class="form-row">
                <label>{{ t('priceAlert.condition') }}</label>
                <div class="condition-selector">
                  <button
                    class="condition-btn"
                    :class="{ active: form.condition === 'above' }"
                    @click="form.condition = 'above'"
                  >
                    <PhArrowUp :size="14" />
                    {{ t('priceAlert.above') }}
                  </button>
                  <button
                    class="condition-btn"
                    :class="{ active: form.condition === 'below' }"
                    @click="form.condition = 'below'"
                  >
                    <PhArrowDown :size="14" />
                    {{ t('priceAlert.below') }}
                  </button>
                </div>
              </div>

              <div class="form-row">
                <label>{{ t('priceAlert.targetPrice') }}</label>
                <div class="price-input-group">
                  <span class="price-prefix">$</span>
                  <input
                    v-model="form.target_price"
                    type="number"
                    step="0.01"
                    min="0"
                    :placeholder="t('priceAlert.pricePlaceholder')"
                    class="price-input"
                  />
                </div>
              </div>

              <div class="form-row">
                <label>{{ t('priceAlert.note') }}</label>
                <input
                  v-model="form.note"
                  type="text"
                  :placeholder="t('priceAlert.notePlaceholder')"
                  class="note-input"
                  maxlength="100"
                />
              </div>

              <div class="form-actions">
                <button class="btn-cancel" @click="showForm = false; resetForm()">
                  {{ t('common.cancel') }}
                </button>
                <button
                  class="btn-submit"
                  :disabled="!form.target_price || form.target_price <= 0 || submitting"
                  @click="createAlert"
                >
                  <PhCheck :size="14" />
                  {{ t('priceAlert.create') }}
                </button>
              </div>
            </div>

            <!-- 预警列表 -->
            <div class="alert-list">
              <div v-if="loading" class="list-loading">
                {{ t('common.loading') }}
              </div>

              <div v-else-if="alerts.length === 0" class="list-empty">
                {{ t('priceAlert.empty') }}
              </div>

              <div
                v-else
                v-for="alert in alerts"
                :key="alert.id"
                class="alert-item"
                :class="{ triggered: alert.triggered, paused: !alert.is_active && !alert.triggered }"
              >
                <div class="alert-main">
                  <div class="alert-symbol">{{ alert.symbol }}</div>
                  <div class="alert-info">
                    <span class="alert-condition">
                      <PhArrowUp v-if="alert.condition === 'above'" :size="12" />
                      <PhArrowDown v-else :size="12" />
                      {{ alert.condition === 'above' ? t('priceAlert.above') : t('priceAlert.below') }}
                    </span>
                    <span class="alert-price font-mono">${{ parseFloat(alert.target_price).toLocaleString() }}</span>
                  </div>
                  <div v-if="alert.note" class="alert-note">{{ alert.note }}</div>
                </div>

                <div class="alert-status">
                  <span v-if="alert.triggered" class="status-triggered">
                    {{ t('priceAlert.triggered') }}
                    <span class="triggered-time">{{ formatTime(alert.triggered_at) }}</span>
                  </span>
                  <span v-else-if="!alert.is_active" class="status-paused">
                    {{ t('priceAlert.paused') }}
                  </span>
                  <span v-else class="status-active">
                    {{ t('priceAlert.active') }}
                  </span>
                </div>

                <div class="alert-actions">
                  <button
                    v-if="!alert.triggered"
                    class="action-btn"
                    :title="alert.is_active ? t('priceAlert.pause') : t('priceAlert.resume')"
                    @click="toggleAlert(alert)"
                  >
                    <PhPause v-if="alert.is_active" :size="14" />
                    <PhPlay v-else :size="14" />
                  </button>
                  <button
                    class="action-btn action-delete"
                    :title="t('common.delete')"
                    @click="deleteAlert(alert)"
                  >
                    <PhTrash :size="14" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
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

.dialog-panel {
  width: 480px;
  max-width: 90vw;
  max-height: 80vh;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--gap-md) var(--gap-lg);
  border-bottom: 1px solid var(--color-border);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  color: var(--color-text-primary);
}

.header-left h3 {
  font-size: 0.875rem;
  font-weight: 600;
}

.active-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  background: var(--color-accent-primary);
  color: #fff;
  padding: 1px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast), color var(--transition-fast);
}

.close-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.dialog-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--gap-lg);
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

/* 新建按钮 */
.add-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-xs);
  width: 100%;
  padding: var(--gap-sm) var(--gap-md);
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-accent-primary);
  background: transparent;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
}

.add-btn:hover {
  border-color: var(--color-accent-primary);
  background: rgba(75, 131, 240, 0.05);
}

/* 表单 */
.alert-form {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.form-row label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.symbol-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.symbol-chip {
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.symbol-chip:hover {
  color: var(--color-text-primary);
  border-color: var(--color-border-hover, var(--color-border));
}

.symbol-chip.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

.condition-selector {
  display: flex;
  gap: var(--gap-xs);
}

.condition-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: var(--gap-sm);
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.condition-btn:hover {
  color: var(--color-text-primary);
}

.condition-btn.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

.price-input-group {
  display: flex;
  align-items: center;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.price-input-group:focus-within {
  border-color: var(--color-accent-primary);
}

.price-prefix {
  padding: 0 var(--gap-sm);
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  font-family: var(--font-mono);
}

.price-input {
  flex: 1;
  padding: var(--gap-sm);
  font-size: 0.8125rem;
  font-family: var(--font-mono);
  color: var(--color-text-primary);
  background: transparent;
  border: none;
  outline: none;
}

.note-input {
  padding: var(--gap-sm);
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  outline: none;
}

.note-input:focus {
  border-color: var(--color-accent-primary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  padding-top: var(--gap-xs);
}

.btn-cancel {
  padding: var(--gap-xs) var(--gap-md);
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-cancel:hover {
  background: var(--color-bg-elevated);
}

.btn-submit {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: var(--gap-xs) var(--gap-md);
  font-size: 0.8125rem;
  font-weight: 500;
  color: #fff;
  background: var(--color-accent-primary);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: opacity var(--transition-fast);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-submit:hover:not(:disabled) {
  opacity: 0.9;
}

/* 预警列表 */
.alert-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.list-loading,
.list-empty {
  text-align: center;
  padding: var(--gap-xl);
  font-size: 0.8125rem;
  color: var(--color-text-muted);
}

.alert-item {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  transition: border-color var(--transition-fast);
}

.alert-item:hover {
  border-color: var(--color-border-hover, var(--color-border));
}

.alert-item.triggered {
  opacity: 0.6;
}

.alert-item.paused {
  opacity: 0.7;
}

.alert-main {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  min-width: 0;
}

.alert-symbol {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  min-width: 36px;
}

.alert-info {
  display: flex;
  align-items: center;
  gap: 4px;
}

.alert-condition {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.alert-price {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.alert-note {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100px;
}

.alert-status {
  font-size: 0.6875rem;
  white-space: nowrap;
}

.status-active {
  color: var(--color-success);
}

.status-paused {
  color: var(--color-warning);
}

.status-triggered {
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
}

.triggered-time {
  font-size: 0.625rem;
  color: var(--color-text-muted);
}

.alert-actions {
  display: flex;
  gap: 2px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: background var(--transition-fast), color var(--transition-fast);
}

.action-btn:hover {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

.action-delete:hover {
  color: var(--color-error);
}

/* 过渡动画 */
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 180ms ease;
}

.dialog-enter-active .dialog-panel,
.dialog-leave-active .dialog-panel {
  transition: transform 180ms ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-from .dialog-panel {
  transform: scale(0.96);
}

.dialog-leave-to .dialog-panel {
  transform: scale(0.96);
}
</style>
