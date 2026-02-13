<script setup>
/**
 * 添加目标对话框
 * 支持 3 种目标类型：总资产、持仓数量、收益率
 */
import { ref } from 'vue'
import { useGoalStore } from '../stores/goalStore'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: Boolean,
})

const emit = defineEmits(['close'])

const goalStore = useGoalStore()
const { t } = useI18n()

// 表单状态
const goalType = ref('asset_value')
const title = ref('')
const targetValue = ref('')
const currency = ref('BTC')
const deadline = ref('')

const goalTypes = [
  { id: 'asset_value', labelKey: 'goals.type.asset_value' },
  { id: 'holding_amount', labelKey: 'goals.type.holding_amount' },
  { id: 'return_rate', labelKey: 'goals.type.return_rate' },
]

const cryptoCurrencies = ['BTC', 'ETH', 'SOL', 'USDC', 'USDT']

// 提交
const handleSubmit = () => {
  if (!title.value.trim() || !targetValue.value) return

  goalStore.addGoal({
    type: goalType.value,
    title: title.value.trim(),
    targetValue: parseFloat(targetValue.value),
    currency: goalType.value === 'holding_amount' ? currency.value : 'USDC',
    deadline: deadline.value || null,
  })

  // 重置表单
  title.value = ''
  targetValue.value = ''
  currency.value = 'BTC'
  deadline.value = ''
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="goal-overlay" @click.self="emit('close')">
        <div class="goal-dialog">
          <h3 class="dialog-title">{{ t('goals.addGoal') }}</h3>

          <!-- 目标类型 -->
          <div class="form-group">
            <label class="form-label">{{ t('goals.goalType') }}</label>
            <div class="type-buttons">
              <button
                v-for="gt in goalTypes"
                :key="gt.id"
                class="type-btn"
                :class="{ active: goalType === gt.id }"
                @click="goalType = gt.id"
              >
                {{ t(gt.labelKey) }}
              </button>
            </div>
          </div>

          <!-- 目标名称 -->
          <div class="form-group">
            <label class="form-label">{{ t('goals.goalTitle') }}</label>
            <input
              v-model="title"
              type="text"
              class="form-input"
              :placeholder="t('goals.titlePlaceholder')"
            />
          </div>

          <!-- 目标值 -->
          <div class="form-group">
            <label class="form-label">{{ t('goals.targetValue') }}</label>
            <div class="input-with-suffix">
              <input
                v-model="targetValue"
                type="number"
                class="form-input"
                :placeholder="t('goals.valuePlaceholder')"
                step="any"
                min="0"
              />
              <span v-if="goalType === 'return_rate'" class="input-suffix">%</span>
            </div>
          </div>

          <!-- 币种选择（仅持仓数量目标） -->
          <div v-if="goalType === 'holding_amount'" class="form-group">
            <label class="form-label">{{ t('goals.currency') }}</label>
            <select v-model="currency" class="form-select">
              <option v-for="c in cryptoCurrencies" :key="c" :value="c">{{ c }}</option>
            </select>
          </div>

          <!-- 截止日期（可选） -->
          <div class="form-group">
            <label class="form-label">{{ t('goals.deadline') }} <span class="optional">({{ t('goals.optional') }})</span></label>
            <input
              v-model="deadline"
              type="date"
              class="form-input"
            />
          </div>

          <!-- 按钮 -->
          <div class="dialog-actions">
            <button class="btn btn-ghost" @click="emit('close')">{{ t('common.cancel') }}</button>
            <button class="btn btn-primary" @click="handleSubmit" :disabled="!title.trim() || !targetValue">
              {{ t('common.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.goal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 400;
  display: flex;
  justify-content: center;
  padding-top: 15vh;
}

.goal-dialog {
  width: 420px;
  max-width: 90vw;
  max-height: 80vh;
  overflow-y: auto;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--gap-xl);
  align-self: flex-start;
}

.dialog-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-lg);
}

/* 表单 */
.form-group {
  margin-bottom: var(--gap-md);
}

.form-label {
  display: block;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-xs);
}

.optional {
  font-weight: 400;
  color: var(--color-text-muted);
}

.form-input,
.form-select {
  width: 100%;
  padding: var(--gap-sm) var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.8125rem;
  font-family: var(--font-body);
  transition: border-color var(--transition-fast);
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.form-input::placeholder {
  color: var(--color-text-muted);
}

.input-with-suffix {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.input-with-suffix .form-input {
  flex: 1;
}

.input-suffix {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  font-weight: 500;
}

/* 类型按钮 */
.type-buttons {
  display: flex;
  gap: 2px;
  background: var(--color-bg-tertiary);
  padding: 2px;
  border-radius: var(--radius-sm);
}

.type-btn {
  flex: 1;
  padding: 6px 8px;
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.type-btn:hover {
  color: var(--color-text-primary);
}

.type-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* 按钮区 */
.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  margin-top: var(--gap-lg);
}

/* 过渡 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 120ms ease;
}

.modal-enter-active .goal-dialog,
.modal-leave-active .goal-dialog {
  transition: transform 120ms ease, opacity 120ms ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .goal-dialog {
  transform: scale(0.96) translateY(-8px);
  opacity: 0;
}

.modal-leave-to .goal-dialog {
  transform: scale(0.96);
  opacity: 0;
}
</style>
