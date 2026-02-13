<template>
  <Transition name="modal">
    <div v-if="visible" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ dialogTitle }}</h3>
          <button class="close-btn" @click="emit('close')">
            <PhX :size="20" />
          </button>
        </div>

        <div class="modal-body">
          <!-- CEX 账户表单 -->
          <div v-if="accountType === 'cex'" class="form-section">
            <p class="form-description">{{ t('accounts.addCexDescription') }}</p>
            <div class="form-group">
              <label for="cexName" class="input-label">{{ t('accounts.accountName') }}</label>
              <input type="text" id="cexName" v-model="form.name" class="input-field" :placeholder="t('accounts.accountNamePlaceholder')" />
            </div>
            <div class="form-group">
              <label for="cexExchange" class="input-label">{{ t('accounts.exchange') }}</label>
              <select id="cexExchange" v-model="form.exchange" class="input-field">
                <option value="" disabled>{{ t('accounts.selectExchange') }}</option>
                <option v-for="ex in availableExchanges" :key="ex.id" :value="ex.id">{{ ex.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="cexApiKey" class="input-label">{{ t('accounts.apiKey') }}</label>
              <input type="text" id="cexApiKey" v-model="form.apiKey" class="input-field" :placeholder="t('accounts.apiKeyPlaceholder')" />
              <p class="input-hint">{{ t('accounts.apiKeyHint') }}</p>
            </div>
            <div class="form-group">
              <label for="cexApiSecret" class="input-label">{{ t('accounts.apiSecret') }}</label>
              <input type="password" id="cexApiSecret" v-model="form.apiSecret" class="input-field" :placeholder="t('accounts.apiSecretPlaceholder')" />
              <p class="input-hint">{{ t('accounts.apiSecretHint') }}</p>
            </div>
            <div class="warning-box">
              <PhWarning :size="20" weight="fill" />
              <p>{{ t('accounts.apiPermissionWarning') }}</p>
            </div>
          </div>

          <!-- Web3 钱包表单 -->
          <div v-else-if="accountType === 'blockchain'" class="form-section">
            <p class="form-description">{{ t('accounts.addWalletDescription') }}</p>
            <div class="form-group">
              <label for="walletName" class="input-label">{{ t('accounts.walletName') }}</label>
              <input type="text" id="walletName" v-model="form.name" class="input-field" :placeholder="t('accounts.walletNamePlaceholder')" />
            </div>
            <div class="form-group">
              <label for="walletAddress" class="input-label">{{ t('accounts.walletAddress') }}</label>
              <input type="text" id="walletAddress" v-model="form.address" class="input-field" :placeholder="t('accounts.walletAddressPlaceholder')" />
              <p class="input-hint">{{ t('accounts.walletAddressHint') }}</p>
            </div>
            <div class="form-group">
              <label for="walletBlockchain" class="input-label">{{ t('accounts.blockchain') }}</label>
              <select id="walletBlockchain" v-model="form.blockchain" class="input-field">
                <option value="" disabled>{{ t('accounts.selectBlockchain') }}</option>
                <option v-for="chain in availableBlockchains" :key="chain.id" :value="chain.id">{{ chain.name }}</option>
              </select>
            </div>
            <div class="warning-box">
              <PhWarning :size="20" weight="fill" />
              <p>{{ t('accounts.walletPermissionWarning') }}</p>
            </div>
          </div>

          <!-- 手动资产表单 — 三步流程 -->
          <div v-else-if="accountType === 'manual'" class="form-section">
            <p class="form-description">{{ t('accounts.addManualDescription') }}</p>

            <!-- 步骤一：选择资产类型 -->
            <div v-if="manualStep === 'type'" class="step-type">
              <label class="input-label">{{ t('accounts.assetType') }}</label>
              <div class="type-grid">
                <button
                  v-for="atype in availableAssetTypes"
                  :key="atype.id"
                  class="type-btn"
                  :class="{ active: form.type === atype.id }"
                  @click="selectAssetType(atype.id)"
                >
                  <span class="type-icon">{{ atype.icon }}</span>
                  <span class="type-label">{{ t(atype.labelKey) }}</span>
                </button>
              </div>
            </div>

            <!-- 步骤二：选择机构（bank/stock/fund 需要） -->
            <div v-else-if="manualStep === 'institution'" class="step-institution">
              <button class="back-btn" @click="manualStep = 'type'">
                <PhArrowLeft :size="14" />
                {{ t('accounts.backToTypeSelect') }}
              </button>
              <label class="input-label">{{ t('accounts.selectInstitution') }}</label>
              <input
                type="text"
                v-model="institutionSearch"
                class="input-field search-input"
                :placeholder="t('accounts.searchInstitution')"
              />
              <div class="institution-list">
                <template v-for="group in filteredInstitutions" :key="group.region">
                  <div class="region-label">{{ t('institutions.region.' + group.region) }}</div>
                  <div class="institution-grid">
                    <button
                      v-for="inst in group.items"
                      :key="inst.id"
                      class="institution-btn"
                      @click="selectInstitution(inst.name)"
                    >
                      {{ inst.name }}
                    </button>
                  </div>
                </template>
                <!-- 自定义输入 -->
                <div class="region-label">{{ t('institutions.custom') }}</div>
                <div class="custom-institution">
                  <input
                    type="text"
                    v-model="customInstitution"
                    class="input-field"
                    :placeholder="t('accounts.customInstitution')"
                    @keyup.enter="selectInstitution(customInstitution)"
                  />
                  <button
                    class="btn btn-primary btn-sm"
                    :disabled="!customInstitution.trim()"
                    @click="selectInstitution(customInstitution)"
                  >
                    {{ t('common.confirm') }}
                  </button>
                </div>
              </div>
            </div>

            <!-- 步骤三：填写详情 -->
            <div v-else-if="manualStep === 'detail'" class="step-detail">
              <button class="back-btn" @click="goBackFromDetail">
                <PhArrowLeft :size="14" />
                {{ needsInstitution ? t('accounts.backToInstitution') : t('accounts.backToTypeSelect') }}
              </button>

              <!-- 已选信息摘要 -->
              <div class="selected-summary">
                <span class="summary-badge">{{ currentTypeIcon }} {{ t(currentTypeLabelKey) }}</span>
                <span v-if="form.institution" class="summary-badge">{{ form.institution }}</span>
              </div>

              <div class="form-group">
                <label for="manualName" class="input-label">{{ t('accounts.assetName') }}</label>
                <input type="text" id="manualName" v-model="form.name" class="input-field" :placeholder="t('accounts.assetNamePlaceholder')" />
              </div>
              <div class="form-group">
                <label for="manualBalance" class="input-label">{{ t('accounts.balance') }}</label>
                <input type="number" id="manualBalance" v-model.number="form.balance" class="input-field" :placeholder="t('accounts.balancePlaceholder')" />
              </div>
              <div class="form-group">
                <label for="manualCurrency" class="input-label">{{ t('accounts.currency') }}</label>
                <select id="manualCurrency" v-model="form.currency" class="input-field">
                  <option value="" disabled>{{ t('accounts.selectCurrency') }}</option>
                  <option v-for="curr in availableCurrencies" :key="curr.code" :value="curr.code">
                    {{ curr.code }} — {{ curr.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label for="manualNote" class="input-label">{{ t('accounts.note') }}</label>
                <textarea id="manualNote" v-model="form.note" class="input-field" :placeholder="t('accounts.notePlaceholder')"></textarea>
              </div>
            </div>
          </div>

          <!-- 错误提示 -->
          <div v-if="error" class="error-message">
            <PhWarning :size="16" weight="fill" />
            <span>{{ error }}</span>
          </div>
        </div>

        <!-- 底部按钮：手动资产在 type/institution 步骤不显示提交 -->
        <div class="modal-actions">
          <button class="btn btn-secondary" @click="emit('close')">
            {{ t('common.cancel') }}
          </button>
          <button
            v-if="showSubmitButton"
            class="btn btn-primary"
            @click="handleSubmit"
            :disabled="isSubmitting"
          >
            <PhSpinnerGap v-if="isSubmitting" :size="20" class="animate-spin" />
            <span v-else>{{ submitButtonText }}</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { PhX, PhWarning, PhSpinnerGap, PhArrowLeft } from '@phosphor-icons/vue'
import { useI18n } from '../composables/useI18n'
import { useAccountStore } from '../stores/accountStore'
import { getInstitutionsByType, searchInstitutions } from '../data/institutions'

const emit = defineEmits(['close', 'submitted'])

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  accountType: {
    type: String,
    required: true
  },
  editingAccount: {
    type: Object,
    default: null
  }
})

const { t } = useI18n()
const accountStore = useAccountStore()

const isSubmitting = ref(false)
const error = ref(null)

// 手动资产表单步骤
const manualStep = ref('type') // 'type' | 'institution' | 'detail'
const institutionSearch = ref('')
const customInstitution = ref('')

// 表单数据
const form = reactive({
  name: '',
  exchange: '',
  apiKey: '',
  apiSecret: '',
  address: '',
  blockchain: '',
  type: '',
  institution: '',
  balance: 0,
  currency: '',
  note: ''
})

// 重置表单函数
const resetForm = () => {
  form.name = ''
  form.exchange = ''
  form.apiKey = ''
  form.apiSecret = ''
  form.address = ''
  form.blockchain = ''
  form.type = ''
  form.institution = ''
  form.balance = 0
  form.currency = ''
  form.note = ''
  manualStep.value = 'type'
  institutionSearch.value = ''
  customInstitution.value = ''
  error.value = null
}

// 可用的交易所、区块链
const availableExchanges = ref([
  { id: 'Binance', name: 'Binance' },
  { id: 'OKX', name: 'OKX' },
  { id: 'Coinbase', name: 'Coinbase' }
])

const availableBlockchains = ref([
  { id: 'ETH', name: 'Ethereum' },
  { id: 'BSC', name: 'BNB Chain' },
  { id: 'SOL', name: 'Solana' },
  { id: 'MATIC', name: 'Polygon' }
])

// 资产类型（去掉 other，加入图标）
const availableAssetTypes = ref([
  { id: 'bank', labelKey: 'accounts.bankDeposit', icon: '🏦' },
  { id: 'cash', labelKey: 'accounts.cash', icon: '💵' },
  { id: 'stock', labelKey: 'accounts.stock', icon: '📈' },
  { id: 'fund', labelKey: 'accounts.fund', icon: '📊' },
])

// 法币列表（7 种）
const availableCurrencies = ref([
  { code: 'CNY', symbol: '¥', name: '人民币' },
  { code: 'USD', symbol: '$', name: '美元' },
  { code: 'HKD', symbol: 'HK$', name: '港币' },
  { code: 'EUR', symbol: '€', name: '欧元' },
  { code: 'GBP', symbol: '£', name: '英镑' },
  { code: 'JPY', symbol: '¥', name: '日元' },
  { code: 'SGD', symbol: 'S$', name: '新加坡元' },
])

// 是否需要机构选择步骤
const needsInstitution = computed(() => {
  return ['bank', 'stock', 'fund'].includes(form.type)
})

// 当前类型图标
const currentTypeIcon = computed(() => {
  const type = availableAssetTypes.value.find(t => t.id === form.type)
  return type?.icon || ''
})

// 当前类型标签 Key
const currentTypeLabelKey = computed(() => {
  const type = availableAssetTypes.value.find(t => t.id === form.type)
  return type?.labelKey || ''
})

// 过滤后的机构列表
const filteredInstitutions = computed(() => {
  if (!form.type) return []
  const groups = getInstitutionsByType(form.type)
  if (!institutionSearch.value) return groups
  return searchInstitutions(groups, institutionSearch.value)
})

// 是否显示提交按钮
const showSubmitButton = computed(() => {
  if (props.accountType !== 'manual') return true
  return manualStep.value === 'detail'
})

// 选择资产类型
const selectAssetType = (typeId) => {
  form.type = typeId
  form.institution = ''
  if (needsInstitution.value) {
    manualStep.value = 'institution'
  } else {
    // cash 类型直接进入详情
    manualStep.value = 'detail'
  }
}

// 选择机构
const selectInstitution = (name) => {
  if (!name || !name.trim()) return
  form.institution = name.trim()
  institutionSearch.value = ''
  customInstitution.value = ''
  manualStep.value = 'detail'
}

// 从详情步骤返回
const goBackFromDetail = () => {
  if (needsInstitution.value) {
    form.institution = ''
    manualStep.value = 'institution'
  } else {
    form.type = ''
    manualStep.value = 'type'
  }
}

// 对话框标题
const dialogTitle = computed(() => {
  if (props.editingAccount) {
    return t('accounts.editAccountTitle')
  }
  switch (props.accountType) {
    case 'cex': return t('accounts.addCexTitle')
    case 'blockchain': return t('accounts.addWalletTitle')
    case 'manual': return t('accounts.addManualTitle')
    default: return t('accounts.addAccount')
  }
})

// 提交按钮文本
const submitButtonText = computed(() => {
  return props.editingAccount ? t('common.update') : t('common.add')
})

// 监听 editingAccount 变化，填充表单
watch(() => props.editingAccount, (newVal) => {
  if (newVal) {
    Object.assign(form, newVal)
    // 编辑模式直接进入详情步骤
    manualStep.value = 'detail'
  } else {
    resetForm()
  }
}, { immediate: true })

// 监听 visible 变化，在对话框关闭时重置表单和错误
watch(() => props.visible, (newVal) => {
  if (!newVal) {
    resetForm()
    error.value = null
  }
})

// 提交表单
const handleSubmit = async () => {
  isSubmitting.value = true
  error.value = null

  try {
    if (props.editingAccount) {
      if (props.accountType === 'cex') {
        await accountStore.updateCexAccount(form.id, form)
      } else if (props.accountType === 'blockchain') {
        await accountStore.updateWalletAddress(form.id, form)
      } else if (props.accountType === 'manual') {
        await accountStore.updateManualAsset(form.id, form)
      }
    } else {
      if (props.accountType === 'cex') {
        await accountStore.addCexAccount(form)
      } else if (props.accountType === 'blockchain') {
        await accountStore.addWalletAddress(form)
      } else if (props.accountType === 'manual') {
        await accountStore.addManualAsset(form)
      }
    }
    emit('submitted')
    emit('close')
  } catch (err) {
    error.value = err.message || t('common.operationFailed')
  } finally {
    isSubmitting.value = false
  }
}
</script>

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
  max-width: 520px;
  max-height: 85vh;
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
  font-family: var(--font-heading);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.close-btn:hover {
  color: var(--color-text-primary);
}

.modal-body {
  margin-bottom: var(--gap-lg);
  overflow-y: auto;
  flex: 1;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.form-description {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-sm);
  line-height: 1.5;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--gap-xs);
}

.input-label {
  font-size: 12px;
  color: var(--color-text-primary);
  font-weight: 500;
}

.input-field {
  width: 100%;
  padding: 8px 12px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 13px;
  transition: border-color var(--transition-fast);
}

.input-field:focus {
  outline: none;
  border-color: var(--color-accent-primary);
}

.input-field::placeholder {
  color: var(--color-text-muted);
}

.input-field[type="number"] {
  -moz-appearance: textfield;
}
.input-field::-webkit-outer-spin-button,
.input-field::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

textarea.input-field {
  min-height: 72px;
  resize: vertical;
}

.input-hint {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
}

.warning-box {
  display: flex;
  align-items: flex-start;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: color-mix(in srgb, var(--color-warning) 10%, transparent);
  border: 1px solid color-mix(in srgb, var(--color-warning) 30%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-warning);
  font-size: 12px;
  line-height: 1.4;
  margin-top: var(--gap-sm);
}

.warning-box p {
  flex: 1;
}

.error-message {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-sm) var(--gap-md);
  background: color-mix(in srgb, var(--color-error) 10%, transparent);
  border-radius: var(--radius-sm);
  color: var(--color-error);
  font-size: 12px;
  margin-top: var(--gap-md);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap-sm);
  padding-top: var(--gap-md);
  border-top: 1px solid var(--color-border);
}

/* ========== 手动资产 — 步骤一：资产类型选择 ========== */
.type-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--gap-sm);
  margin-top: var(--gap-sm);
}

.type-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-xs);
  padding: var(--gap-md) var(--gap-sm);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: border-color var(--transition-fast), color var(--transition-fast), background var(--transition-fast);
}

.type-btn:hover {
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
  background: var(--color-bg-elevated);
}

.type-btn.active {
  border-color: var(--color-accent-primary);
  color: var(--color-accent-primary);
  background: rgba(75, 131, 240, 0.08);
}

.type-icon {
  font-size: 1.25rem;
  line-height: 1;
}

.type-label {
  font-size: 0.8125rem;
  font-weight: 500;
}

/* ========== 手动资产 — 步骤二：机构选择器 ========== */
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--gap-xs);
  background: none;
  border: none;
  color: var(--color-text-secondary);
  font-size: 0.75rem;
  cursor: pointer;
  padding: 0;
  margin-bottom: var(--gap-sm);
  transition: color var(--transition-fast);
}

.back-btn:hover {
  color: var(--color-accent-primary);
}

.search-input {
  margin-top: var(--gap-sm);
  margin-bottom: var(--gap-sm);
}

.institution-list {
  max-height: 320px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.region-label {
  font-size: 0.6875rem;
  font-weight: 500;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding-top: var(--gap-xs);
}

.institution-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap-xs);
}

.institution-btn {
  padding: 6px 10px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 0.75rem;
  cursor: pointer;
  transition: border-color var(--transition-fast), background var(--transition-fast);
  white-space: nowrap;
}

.institution-btn:hover {
  border-color: var(--color-accent-primary);
  background: rgba(75, 131, 240, 0.08);
}

.custom-institution {
  display: flex;
  gap: var(--gap-xs);
  align-items: center;
}

.custom-institution .input-field {
  flex: 1;
}

/* ========== 手动资产 — 步骤三：详情 ========== */
.selected-summary {
  display: flex;
  gap: var(--gap-xs);
  margin-bottom: var(--gap-md);
  flex-wrap: wrap;
}

.summary-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 3px 8px;
  background: rgba(75, 131, 240, 0.1);
  border-radius: var(--radius-sm);
  color: var(--color-accent-primary);
  font-size: 0.75rem;
  font-weight: 500;
}

/* ========== 模态框动画 ========== */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.97);
}
</style>
