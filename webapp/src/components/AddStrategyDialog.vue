<script setup>
/**
 * 添加策略对话框
 * 步骤 1：选择策略类型（再平衡/定投/止盈止损）
 * 步骤 2：配置参数
 * 步骤 3：确认并启用
 */
import { ref, watch } from 'vue'
import { PhX, PhArrowsClockwise, PhCalendarBlank, PhBell } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'

const props = defineProps({
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'created'])
const { t } = useI18n()

// 步骤
const step = ref(1)
const selectedType = ref('')

// 表单数据 — 再平衡
const rebalanceForm = ref({
  name: '',
  targets: [
    { symbol: 'BTC', targetPct: 50 },
    { symbol: 'ETH', targetPct: 30 },
    { symbol: 'USDC', targetPct: 20 },
  ],
  deviationThreshold: 5,
})

// 表单数据 — 定投
const dcaForm = ref({
  name: '',
  symbol: 'BTC',
  amount: 1000,
  currency: 'USDC',
  frequency: 'monthly',
})

// 表单数据 — 止盈止损
const alertForm = ref({
  name: '',
  symbol: 'BTC',
  direction: 'above',
  targetPrice: 0,
  note: '',
})

// 策略类型选项
const typeOptions = [
  { id: 'rebalance', icon: PhArrowsClockwise, labelKey: 'strategy.typeRebalance', descKey: 'strategy.typeRebalanceDesc' },
  { id: 'dca', icon: PhCalendarBlank, labelKey: 'strategy.typeDca', descKey: 'strategy.typeDcaDesc' },
  { id: 'alert', icon: PhBell, labelKey: 'strategy.typeAlert', descKey: 'strategy.typeAlertDesc' },
]

const frequencyOptions = [
  { id: 'weekly', labelKey: 'strategy.weekly' },
  { id: 'biweekly', labelKey: 'strategy.biweekly' },
  { id: 'monthly', labelKey: 'strategy.monthly' },
]

// 重置
const resetDialog = () => {
  step.value = 1
  selectedType.value = ''
  rebalanceForm.value = { name: '', targets: [{ symbol: 'BTC', targetPct: 50 }, { symbol: 'ETH', targetPct: 30 }, { symbol: 'USDC', targetPct: 20 }], deviationThreshold: 5 }
  dcaForm.value = { name: '', symbol: 'BTC', amount: 1000, currency: 'USDC', frequency: 'monthly' }
  alertForm.value = { name: '', symbol: 'BTC', direction: 'above', targetPrice: 0, note: '' }
}

watch(() => props.visible, (val) => {
  if (!val) resetDialog()
})

// 选择类型
const selectType = (type) => {
  selectedType.value = type
  step.value = 2
}

// 返回上一步
const goBack = () => {
  if (step.value > 1) step.value--
}

// 下一步
const goNext = () => {
  if (step.value === 2) step.value = 3
}

// 确认创建
const confirmCreate = () => {
  let strategy = { type: selectedType.value }
  if (selectedType.value === 'rebalance') {
    strategy.name = rebalanceForm.value.name || t('strategy.typeRebalance')
    strategy.config = { ...rebalanceForm.value, targets: rebalanceForm.value.targets.map(t => ({ ...t, currentPct: t.targetPct })) }
  } else if (selectedType.value === 'dca') {
    strategy.name = dcaForm.value.name || `${t('strategy.typeDca')} ${dcaForm.value.symbol}`
    strategy.config = { ...dcaForm.value, totalInvested: 0, totalPurchased: 0, avgPrice: 0, nextDate: '' }
  } else {
    strategy.name = alertForm.value.name || `${alertForm.value.symbol} ${t('strategy.typeAlert')}`
    strategy.config = { ...alertForm.value, currentPrice: 0 }
  }
  emit('created', strategy)
  emit('close')
}
</script>

<template>
  <Transition name="modal">
    <div v-if="visible" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-content">
        <!-- 头部 -->
        <div class="modal-header">
          <h3>{{ t('strategy.addStrategy') }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="20" />
          </button>
        </div>

        <div class="modal-body">
          <!-- 步骤 1：选择类型 -->
          <div v-if="step === 1" class="step-select">
            <p class="step-desc">{{ t('strategy.selectType') }}</p>
            <div class="type-list">
              <button
                v-for="opt in typeOptions"
                :key="opt.id"
                class="type-btn"
                @click="selectType(opt.id)"
              >
                <component :is="opt.icon" :size="24" class="type-icon" />
                <div class="type-info">
                  <span class="type-label">{{ t(opt.labelKey) }}</span>
                  <span class="type-desc">{{ t(opt.descKey) }}</span>
                </div>
              </button>
            </div>
          </div>

          <!-- 步骤 2：配置参数 -->
          <div v-else-if="step === 2" class="step-config">
            <!-- 再平衡 -->
            <template v-if="selectedType === 'rebalance'">
              <div class="form-group">
                <label class="input-label">{{ t('strategy.strategyName') }}</label>
                <input type="text" v-model="rebalanceForm.name" class="input-field" :placeholder="t('strategy.typeRebalance')" />
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.targetAllocation') }}</label>
                <div v-for="(tgt, idx) in rebalanceForm.targets" :key="idx" class="target-row">
                  <input type="text" v-model="tgt.symbol" class="input-field input-sm" />
                  <input type="number" v-model.number="tgt.targetPct" class="input-field input-sm font-mono" min="0" max="100" />
                  <span class="pct-sign">%</span>
                </div>
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.deviationThreshold') }}</label>
                <div class="inline-input">
                  <input type="number" v-model.number="rebalanceForm.deviationThreshold" class="input-field input-sm font-mono" min="1" max="50" />
                  <span class="pct-sign">%</span>
                </div>
              </div>
            </template>

            <!-- 定投 -->
            <template v-else-if="selectedType === 'dca'">
              <div class="form-group">
                <label class="input-label">{{ t('strategy.strategyName') }}</label>
                <input type="text" v-model="dcaForm.name" class="input-field" :placeholder="t('strategy.typeDca')" />
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.dcaSymbol') }}</label>
                <input type="text" v-model="dcaForm.symbol" class="input-field" placeholder="BTC" />
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.dcaAmount') }}</label>
                <div class="inline-input">
                  <span class="currency-prefix">$</span>
                  <input type="number" v-model.number="dcaForm.amount" class="input-field font-mono" min="1" />
                </div>
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.dcaFrequency') }}</label>
                <div class="freq-options">
                  <button
                    v-for="f in frequencyOptions"
                    :key="f.id"
                    class="freq-btn"
                    :class="{ active: dcaForm.frequency === f.id }"
                    @click="dcaForm.frequency = f.id"
                  >
                    {{ t(f.labelKey) }}
                  </button>
                </div>
              </div>
            </template>

            <!-- 止盈止损 -->
            <template v-else>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.strategyName') }}</label>
                <input type="text" v-model="alertForm.name" class="input-field" :placeholder="t('strategy.typeAlert')" />
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.alertSymbol') }}</label>
                <input type="text" v-model="alertForm.symbol" class="input-field" placeholder="BTC" />
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.alertDirection') }}</label>
                <div class="freq-options">
                  <button class="freq-btn" :class="{ active: alertForm.direction === 'above' }" @click="alertForm.direction = 'above'">
                    {{ t('strategy.priceAbove') }}
                  </button>
                  <button class="freq-btn" :class="{ active: alertForm.direction === 'below' }" @click="alertForm.direction = 'below'">
                    {{ t('strategy.priceBelow') }}
                  </button>
                </div>
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.targetPrice') }}</label>
                <div class="inline-input">
                  <span class="currency-prefix">$</span>
                  <input type="number" v-model.number="alertForm.targetPrice" class="input-field font-mono" min="0" />
                </div>
              </div>
              <div class="form-group">
                <label class="input-label">{{ t('strategy.note') }}</label>
                <input type="text" v-model="alertForm.note" class="input-field" />
              </div>
            </template>
          </div>

          <!-- 步骤 3：确认 -->
          <div v-else-if="step === 3" class="step-confirm">
            <p class="confirm-text">{{ t('strategy.confirmHint') }}</p>
            <p class="disclaimer">{{ t('strategy.disclaimer') }}</p>
          </div>
        </div>

        <!-- 底部按钮 -->
        <div class="modal-actions">
          <button v-if="step > 1" class="btn btn-secondary" @click="goBack">
            {{ t('common.back') }}
          </button>
          <button v-else class="btn btn-secondary" @click="emit('close')">
            {{ t('common.cancel') }}
          </button>
          <button v-if="step === 2" class="btn btn-primary" @click="goNext">
            {{ t('strategy.next') }}
          </button>
          <button v-if="step === 3" class="btn btn-primary" @click="confirmCreate">
            {{ t('strategy.createAndEnable') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--gap-lg);
}

.modal-content {
  width: 100%;
  max-width: 460px;
  padding: var(--gap-xl);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-lg);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
}

.modal-header h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
}

.close-btn:hover {
  color: var(--color-text-primary);
}

.modal-body {
  margin-bottom: var(--gap-lg);
  min-height: 200px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

/* 步骤 1 */
.step-desc {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-md);
}

.type-list {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.type-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-md);
  padding: var(--gap-md);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  text-align: left;
  transition: border-color var(--transition-fast);
}

.type-btn:hover {
  border-color: var(--color-accent-primary);
}

.type-icon {
  color: var(--color-accent-primary);
  flex-shrink: 0;
}

.type-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.type-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.type-desc {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
}

/* 步骤 2：表单 */
.form-group {
  margin-bottom: var(--gap-md);
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.input-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.input-field {
  padding: 8px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.8125rem;
}

.input-field:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.input-sm {
  width: 80px;
}

.target-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-xs);
}

.pct-sign {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.inline-input {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.currency-prefix {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  font-weight: 600;
}

.freq-options {
  display: flex;
  gap: 4px;
}

.freq-btn {
  padding: 4px 12px;
  font-size: 0.6875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.freq-btn:hover {
  border-color: var(--color-accent-primary);
}

.freq-btn.active {
  background: var(--color-accent-primary);
  border-color: var(--color-accent-primary);
  color: #fff;
}

/* 步骤 3 */
.confirm-text {
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  margin-bottom: var(--gap-md);
}

.disclaimer {
  font-size: 0.6875rem;
  color: var(--color-warning);
  padding: var(--gap-sm) var(--gap-md);
  background: rgba(245, 158, 11, 0.08);
  border-radius: var(--radius-sm);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
