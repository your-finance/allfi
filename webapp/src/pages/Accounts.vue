<script setup>
/**
 * Accounts 页面 - 账户管理
 * 表格式布局：Tab切换 + 数据表格
 */
import { ref, computed, onMounted } from 'vue'
import {
  PhPlus,
  PhPencilSimple,
  PhTrash,
  PhCopy,
  PhCheck,
  PhWarning,
  PhLinkSimple,
  PhBank,
  PhCurrencyBtc,
  PhArrowsClockwise,
  PhList,
  PhGridFour,
  PhUploadSimple,
  PhCubeTransparent,
  PhImageSquare
} from '@phosphor-icons/vue'
import { useFormatters } from '../composables/useFormatters'
import { useAccountStore } from '../stores/accountStore'
import { useAssetStore } from '../stores/assetStore'
import { useNFTStore } from '../stores/nftStore'
import { useI18n } from '../composables/useI18n'
import { useToast } from '../composables/useToast'
import AddAccountDialog from '../components/AddAccountDialog.vue'
import BatchImportDialog from '../components/BatchImportDialog.vue'
import DeFiPositionCard from '../components/DeFiPositionCard.vue'
import NFTGallery from '../components/NFTGallery.vue'
import WalletConnectButton from '../components/WalletConnectButton.vue'
import WalletConnectDialog from '../components/WalletConnectDialog.vue'

const { formatNumber, formatRelativeTime, shortenAddress, pricingCurrency } = useFormatters()
const accountStore = useAccountStore()
const assetStore = useAssetStore()
const nftStore = useNFTStore()
const { t } = useI18n()
const { showToast } = useToast()

// Tab 状态
const activeTab = ref('cex')
// DeFi 仓位数据
const defiPositions = computed(() => assetStore.defiPositions)

// DeFi 仓位筛选类型
const defiFilterType = ref('all')
const defiTypeFilters = ['all', 'lp', 'staking', 'lending']

// 筛选后的 DeFi 仓位
const filteredDefiPositions = computed(() => {
  if (defiFilterType.value === 'all') return defiPositions.value
  return defiPositions.value.filter(p => p.type === defiFilterType.value)
})

// DeFi 仓位按协议分组
const groupedDefiPositions = computed(() => {
  const groups = {}
  for (const p of filteredDefiPositions.value) {
    if (!groups[p.protocol]) {
      groups[p.protocol] = { protocol: p.protocol, protocolIcon: p.protocolIcon, positions: [], totalValue: 0 }
    }
    groups[p.protocol].positions.push(p)
    groups[p.protocol].totalValue += p.valueUSD
  }
  return Object.values(groups).sort((a, b) => b.totalValue - a.totalValue)
})

// DeFi 仓位总值
const defiTotalValue = computed(() => assetStore.defiTotalValue)

// DeFi 类型筛选标签
const getDefiTypeLabel = (type) => {
  const labels = { all: t('defi.filterAll'), lp: t('defi.typeLp'), staking: t('defi.typeStaking'), lending: t('defi.typeLending') }
  return labels[type] || type
}

const tabs = computed(() => [
  { id: 'cex', labelKey: 'accounts.cexTab', icon: PhCurrencyBtc, count: cexAccounts.value.length },
  { id: 'blockchain', labelKey: 'accounts.blockchainTab', icon: PhLinkSimple, count: walletAddresses.value.length },
  { id: 'defi', labelKey: 'accounts.defiTab', icon: PhCubeTransparent, count: defiPositions.value.length },
  { id: 'nft', labelKey: 'accounts.nftTab', icon: PhImageSquare, count: nftStore.totalCount },
  { id: 'manual', labelKey: 'accounts.manualTab', icon: PhBank, count: manualAssets.value.length }
])

// 视图模式：'list'(表格) 或 'card'(卡片)
const viewMode = ref('card')

// 对话框状态
const showAddDialog = ref(false)
const showDeleteConfirm = ref(false)
const selectedAccount = ref(null)
const editingAccount = ref(null)
const copiedAddress = ref(null)
const showBatchImport = ref(false)
const showWalletConnect = ref(false)

// 操作状态
const isRefreshing = ref(false)
const isDeleting = ref(false)

// 从 Store 获取数据
const cexAccounts = computed(() => accountStore.cexAccounts)
const walletAddresses = computed(() => accountStore.walletAddresses)
const manualAssets = computed(() => accountStore.manualAssets)

// 交易所配置
const exchangeConfig = {
  Binance: { color: '#F3BA2F', logo: '₿' },
  OKX: { color: '#000000', logo: 'OK' },
  Coinbase: { color: '#0052FF', logo: 'CB' }
}

// 区块链配置
const blockchainConfig = {
  ethereum: { color: '#627EEA', name: 'Ethereum', symbol: 'ETH' },
  bsc: { color: '#F3BA2F', name: 'BNB Chain', symbol: 'BSC' },
  polygon: { color: '#8247E5', name: 'Polygon', symbol: 'MATIC' },
  arbitrum: { color: '#28A0F0', name: 'Arbitrum', symbol: 'ARB' },
  optimism: { color: '#FF0420', name: 'Optimism', symbol: 'OP' },
  base: { color: '#0052FF', name: 'Base', symbol: 'BASE' }
}

// 资产类型配置
const assetTypeConfig = {
  bank: { color: '#3B82F6', labelKey: 'accounts.bankDeposit', icon: '🏦' },
  cash: { color: '#10B981', labelKey: 'accounts.cash', icon: '💵' },
  stock: { color: '#8B5CF6', labelKey: 'accounts.stock', icon: '📈' },
  fund: { color: '#F59E0B', labelKey: 'accounts.fund', icon: '📊' },
}

// 按机构分组手动资产
const groupedManualAssets = computed(() => {
  const groups = {}
  manualAssets.value.forEach(asset => {
    // 无机构的（cash 类型）用资产类型标签作为组名
    const key = asset.institution || t(assetTypeConfig[asset.type]?.labelKey || 'accounts.cash')
    if (!groups[key]) {
      groups[key] = {
        institution: key,
        assetType: asset.type,
        assets: [],
        totalUSD: 0
      }
    }
    groups[key].assets.push(asset)
    groups[key].totalUSD += convertManualAssetBalance(asset)
  })
  return Object.values(groups)
})

// 复制地址
const copyAddress = async (address) => {
  try {
    await navigator.clipboard.writeText(address)
    copiedAddress.value = address
    setTimeout(() => { copiedAddress.value = null }, 2000)
    showToast(t('common.copySuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  }
}

// 打开添加对话框
const openAddDialog = (type) => {
  editingAccount.value = null
  if (typeof type === 'string' && ['cex', 'blockchain', 'manual'].includes(type)) {
    activeTab.value = type
  }
  showAddDialog.value = true
}

// 编辑账户
const editAccount = (account) => {
  editingAccount.value = account
  activeTab.value = account.type === 'CEX' ? 'cex' : (account.type === 'Blockchain' ? 'blockchain' : 'manual')
  showAddDialog.value = true
}

// 账户提交成功
const handleAccountSubmitted = async () => {
  await accountStore.loadAll()
  showAddDialog.value = false
}

// 删除确认
const confirmDelete = (account) => {
  selectedAccount.value = account
  showDeleteConfirm.value = true
}

// 执行删除
const deleteAccount = async () => {
  if (!selectedAccount.value) return
  isDeleting.value = true
  try {
    await accountStore.deleteAccount(activeTab.value, selectedAccount.value.id)
    showDeleteConfirm.value = false
    selectedAccount.value = null
    showToast(t('accounts.deleteSuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    isDeleting.value = false
  }
}

// 刷新账户
const refreshAccount = async (account) => {
  if (activeTab.value === 'manual') return
  isRefreshing.value = true
  try {
    await accountStore.refreshAccount(activeTab.value, account.id)
    showToast(t('accounts.refreshSuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  } finally {
    isRefreshing.value = false
  }
}

// 当前列表数据
const currentList = computed(() => {
  switch (activeTab.value) {
    case 'cex': return cexAccounts.value
    case 'blockchain': return walletAddresses.value
    case 'manual': return manualAssets.value
    default: return []
  }
})

// 当前总价值
const currentTotal = computed(() => {
  return currentList.value.reduce((sum, item) => sum + item.balance, 0)
})

// 计价货币符号
const pricingCurrencySymbol = computed(() => {
  const symbols = { USDC: '$', BTC: '₿', ETH: 'Ξ', CNY: '¥' }
  return symbols[pricingCurrency.value] || '$'
})

// 转换手动资产余额为统一计价货币
const convertManualAssetBalance = (asset) => {
  // 如果资产货币与计价货币相同，直接返回
  if (asset.currency === pricingCurrency.value) {
    return asset.balance
  }
  // 否则转换为统一计价货币
  return assetStore.convertValue(asset.balance, pricingCurrency.value)
}

// 钱包连接成功后添加到系统
const handleWalletConnected = async (walletData) => {
  try {
    await accountStore.addWalletAddress(walletData)
    showToast(t('walletConnect.addSuccess'), 'success')
  } catch (err) {
    showToast(t('common.operationFailed'), 'error')
  }
}

onMounted(async () => {
  await accountStore.loadAll()
  // 加载 DeFi 仓位数据
  if (assetStore.defiPositions.length === 0) {
    assetStore.loadDefiPositions()
  }
  // 加载 NFT 数据
  nftStore.initialize()
  if (nftStore.nfts.length === 0) {
    nftStore.fetchNFTs()
  }
})
</script>

<template>
  <div class="accounts-page">
    <!-- 头部：Tab + 添加按钮 -->
    <div class="page-toolbar">
      <div class="tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          <component :is="tab.icon" :size="14" />
          <span class="tab-label">{{ t(tab.labelKey) }}</span>
          <span class="tab-count">{{ tab.count }}</span>
        </button>
      </div>
      <div class="toolbar-right">
        <span class="total-label">{{ t('accounts.total') }}:</span>
        <span class="total-value font-mono">${{ formatNumber(currentTotal) }}</span>

        <!-- 视图切换按钮组 -->
        <div class="view-toggle">
          <button
            class="view-btn"
            :class="{ active: viewMode === 'list' }"
            :title="t('accounts.listView')"
            @click="viewMode = 'list'"
          >
            <PhList :size="16" />
          </button>
          <button
            class="view-btn"
            :class="{ active: viewMode === 'card' }"
            :title="t('accounts.cardView')"
            @click="viewMode = 'card'"
          >
            <PhGridFour :size="16" />
          </button>
        </div>

        <WalletConnectButton
          v-if="activeTab === 'blockchain'"
          @connect-click="showWalletConnect = true"
        />

        <button
          v-if="activeTab === 'blockchain'"
          class="btn btn-secondary btn-sm"
          @click="showBatchImport = true"
        >
          <PhUploadSimple :size="14" />
          {{ t('batchImport.title') }}
        </button>

        <button class="btn btn-primary btn-sm" @click="openAddDialog(activeTab)">
          <PhPlus :size="14" weight="bold" />
          {{ t('accounts.addAccount') }}
        </button>
      </div>
    </div>

    <!-- CEX 账户 -->
    <template v-if="activeTab === 'cex'">
      <!-- 列表视图 -->
      <div v-if="viewMode === 'list'" class="glass-card table-panel">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('accounts.accountName') }}</th>
              <th>{{ t('accounts.exchange') }}</th>
              <th>{{ t('accounts.accountBalance') }}</th>
              <th>{{ t('accounts.status') }}</th>
              <th>{{ t('accounts.apiKey') }}</th>
              <th>{{ t('accounts.lastSync') }}</th>
              <th class="col-actions">{{ t('accounts.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="account in cexAccounts" :key="account.id">
              <td>
                <div class="name-cell">
                  <div
                    class="logo-badge"
                    :style="{ backgroundColor: exchangeConfig[account.exchange]?.color }"
                  >{{ exchangeConfig[account.exchange]?.logo }}</div>
                  <span class="name-text">{{ account.name }}</span>
                </div>
              </td>
              <td>{{ account.exchange }}</td>
              <td class="font-mono">${{ formatNumber(account.balance) }}</td>
              <td>
                <span class="status-dot" :class="account.status" />
                <span class="status-text" :class="account.status">
                  {{ account.status === 'connected' ? t('accounts.connected') : t('accounts.connectionError') }}
                </span>
              </td>
              <td class="font-mono text-muted">{{ account.apiKeyMasked }}</td>
              <td class="text-muted">{{ formatRelativeTime(account.lastSync) }}</td>
              <td class="col-actions">
                <div class="action-btns">
                  <button class="icon-btn" :title="t('accounts.refresh')" @click="refreshAccount(account)">
                    <PhArrowsClockwise :size="14" />
                  </button>
                  <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(account)">
                    <PhPencilSimple :size="14" />
                  </button>
                  <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(account)">
                    <PhTrash :size="14" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="cexAccounts.length === 0">
              <td colspan="7" class="empty-row">{{ t('accounts.noCex') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 卡片视图 -->
      <div v-else class="card-grid">
        <div
          v-for="account in cexAccounts"
          :key="account.id"
          class="glass-card account-card"
        >
          <div class="card-header">
            <div class="card-identity">
              <div
                class="logo-badge card-logo"
                :style="{ backgroundColor: exchangeConfig[account.exchange]?.color }"
              >{{ exchangeConfig[account.exchange]?.logo }}</div>
              <div class="card-name-group">
                <span class="card-name">{{ account.name }}</span>
                <span class="card-sub">{{ account.exchange }}</span>
              </div>
            </div>
            <span class="status-dot" :class="account.status" />
          </div>
          <div class="card-body">
            <span class="card-label">{{ t('accounts.accountBalance') }}</span>
            <span class="card-value font-mono">${{ formatNumber(account.balance) }}</span>
          </div>
          <div class="card-meta">
            <span class="text-muted">{{ t('accounts.apiKey') }}: {{ account.apiKeyMasked }}</span>
            <span class="text-muted">{{ formatRelativeTime(account.lastSync) }}</span>
          </div>
          <div class="card-actions">
            <button class="icon-btn" :title="t('accounts.refresh')" @click="refreshAccount(account)">
              <PhArrowsClockwise :size="14" />
            </button>
            <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(account)">
              <PhPencilSimple :size="14" />
            </button>
            <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(account)">
              <PhTrash :size="14" />
            </button>
          </div>
        </div>
        <!-- 空状态 -->
        <div v-if="cexAccounts.length === 0" class="glass-card account-card empty-card">
          <span class="text-muted">{{ t('accounts.noCex') }}</span>
        </div>

        <!-- 添加账户卡片 -->
        <div class="glass-card account-card add-account-card" @click="openAddDialog('cex')">
          <PhPlus :size="32" weight="light" class="add-icon" />
          <span class="add-label">{{ t('accounts.addAccount') }}</span>
        </div>
      </div>
    </template>

    <!-- 链上钱包 -->
    <template v-if="activeTab === 'blockchain'">
      <!-- 列表视图 -->
      <div v-if="viewMode === 'list'" class="glass-card table-panel">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('accounts.walletName') }}</th>
              <th>{{ t('accounts.blockchain') }}</th>
              <th>{{ t('accounts.address') }}</th>
              <th>{{ t('accounts.walletBalance') }}</th>
              <th>{{ t('accounts.status') }}</th>
              <th>{{ t('accounts.lastSync') }}</th>
              <th class="col-actions">{{ t('accounts.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wallet in walletAddresses" :key="wallet.id">
              <td>
                <div class="name-cell">
                  <div
                    class="logo-badge"
                    :style="{ backgroundColor: blockchainConfig[wallet.blockchain]?.color }"
                  >{{ blockchainConfig[wallet.blockchain]?.symbol || wallet.blockchain.substring(0,3).toUpperCase() }}</div>
                  <span class="name-text">{{ wallet.name }}</span>
                </div>
              </td>
              <td>{{ blockchainConfig[wallet.blockchain]?.name || wallet.blockchain }}</td>
              <td>
                <div class="address-cell">
                  <span class="font-mono">{{ shortenAddress(wallet.address, 6) }}</span>
                  <button
                    class="copy-btn"
                    @click="copyAddress(wallet.address)"
                  >
                    <PhCheck v-if="copiedAddress === wallet.address" :size="12" weight="bold" />
                    <PhCopy v-else :size="12" />
                  </button>
                </div>
              </td>
              <td class="font-mono">${{ formatNumber(wallet.balance) }}</td>
              <td>
                <span class="status-dot" :class="wallet.status" />
                <span class="status-text" :class="wallet.status">
                  {{ wallet.status === 'active' ? t('accounts.active') : t('accounts.inactive') }}
                </span>
              </td>
              <td class="text-muted">{{ formatRelativeTime(wallet.lastSync) }}</td>
              <td class="col-actions">
                <div class="action-btns">
                  <button class="icon-btn" :title="t('accounts.refresh')" @click="refreshAccount(wallet)">
                    <PhArrowsClockwise :size="14" />
                  </button>
                  <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(wallet)">
                    <PhPencilSimple :size="14" />
                  </button>
                  <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(wallet)">
                    <PhTrash :size="14" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="walletAddresses.length === 0">
              <td colspan="7" class="empty-row">{{ t('accounts.noWallet') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 卡片视图 -->
      <div v-else class="card-grid">
        <div
          v-for="wallet in walletAddresses"
          :key="wallet.id"
          class="glass-card account-card"
        >
          <div class="card-header">
            <div class="card-identity">
              <div
                class="logo-badge card-logo"
                :style="{ backgroundColor: blockchainConfig[wallet.blockchain]?.color }"
              >{{ blockchainConfig[wallet.blockchain]?.symbol || wallet.blockchain.substring(0,3).toUpperCase() }}</div>
              <div class="card-name-group">
                <span class="card-name">{{ wallet.name }}</span>
                <span class="card-sub">{{ blockchainConfig[wallet.blockchain]?.name || wallet.blockchain }}</span>
              </div>
            </div>
            <span class="status-dot" :class="wallet.status" />
          </div>
          <div class="card-body">
            <span class="card-label">{{ t('accounts.walletBalance') }}</span>
            <span class="card-value font-mono">${{ formatNumber(wallet.balance) }}</span>
          </div>
          <div class="card-meta">
            <div class="address-cell">
              <span class="font-mono">{{ shortenAddress(wallet.address, 6) }}</span>
              <button
                class="copy-btn"
                @click="copyAddress(wallet.address)"
              >
                <PhCheck v-if="copiedAddress === wallet.address" :size="12" weight="bold" />
                <PhCopy v-else :size="12" />
              </button>
            </div>
            <span class="text-muted">{{ formatRelativeTime(wallet.lastSync) }}</span>
          </div>
          <div class="card-actions">
            <button class="icon-btn" :title="t('accounts.refresh')" @click="refreshAccount(wallet)">
              <PhArrowsClockwise :size="14" />
            </button>
            <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(wallet)">
              <PhPencilSimple :size="14" />
            </button>
            <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(wallet)">
              <PhTrash :size="14" />
            </button>
          </div>
        </div>
        <!-- 空状态 -->
        <div v-if="walletAddresses.length === 0" class="glass-card account-card empty-card">
          <span class="text-muted">{{ t('accounts.noWallet') }}</span>
        </div>

        <!-- 添加账户卡片 -->
        <div class="glass-card account-card add-account-card" @click="openAddDialog('blockchain')">
          <PhPlus :size="32" weight="light" class="add-icon" />
          <span class="add-label">{{ t('accounts.addAccount') }}</span>
        </div>
      </div>
    </template>

    <!-- DeFi 仓位 -->
    <template v-if="activeTab === 'defi'">
      <!-- 类型筛选器 -->
      <div class="defi-filter-bar">
        <button
          v-for="type in defiTypeFilters"
          :key="type"
          class="filter-btn"
          :class="{ active: defiFilterType === type }"
          @click="defiFilterType = type"
        >
          {{ getDefiTypeLabel(type) }}
          <span class="filter-count">
            {{ type === 'all' ? defiPositions.length : defiPositions.filter(p => p.type === type).length }}
          </span>
        </button>
      </div>

      <!-- 按协议分组展示 -->
      <div v-if="groupedDefiPositions.length > 0" class="defi-protocol-groups">
        <div
          v-for="group in groupedDefiPositions"
          :key="group.protocol"
          class="glass-card defi-protocol-section"
        >
          <div class="defi-protocol-header">
            <span class="defi-protocol-icon">{{ group.protocolIcon }}</span>
            <span class="defi-protocol-name">{{ group.protocol }}</span>
            <span class="defi-protocol-total font-mono">${{ formatNumber(group.totalValue) }}</span>
          </div>
          <div class="defi-positions-grid">
            <DeFiPositionCard
              v-for="pos in group.positions"
              :key="pos.id"
              :position="pos"
            />
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="glass-card empty-card">
        <span class="text-muted">{{ t('defi.noPositions') }}</span>
      </div>
    </template>

    <!-- NFT 资产画廊 -->
    <template v-if="activeTab === 'nft'">
      <NFTGallery />
    </template>

    <!-- 手动资产（传统资产） -->
    <template v-if="activeTab === 'manual'">
      <!-- 列表视图 — 表格含机构列 -->
      <div v-if="viewMode === 'list'" class="glass-card table-panel">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('accounts.assetName') }}</th>
              <th>{{ t('accounts.institutionColumn') }}</th>
              <th>{{ t('accounts.assetType') }}</th>
              <th>{{ t('accounts.currency') }}</th>
              <th>{{ t('accounts.assetValue') }}</th>
              <th>{{ t('accounts.note') }}</th>
              <th class="col-actions">{{ t('accounts.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="asset in manualAssets" :key="asset.id">
              <td>
                <div class="name-cell">
                  <div
                    class="logo-badge"
                    :style="{ backgroundColor: assetTypeConfig[asset.type]?.color }"
                  >
                    <span>{{ assetTypeConfig[asset.type]?.icon }}</span>
                  </div>
                  <span class="name-text">{{ asset.name }}</span>
                </div>
              </td>
              <td>{{ asset.institution || '-' }}</td>
              <td>{{ t(assetTypeConfig[asset.type]?.labelKey) }}</td>
              <td>{{ asset.currency }}</td>
              <td class="font-mono">
                {{ pricingCurrencySymbol }}{{ formatNumber(convertManualAssetBalance(asset)) }}
              </td>
              <td class="text-muted note-col">{{ asset.note || '-' }}</td>
              <td class="col-actions">
                <div class="action-btns">
                  <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(asset)">
                    <PhPencilSimple :size="14" />
                  </button>
                  <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(asset)">
                    <PhTrash :size="14" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="manualAssets.length === 0">
              <td colspan="7" class="empty-row">{{ t('accounts.noManual') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 卡片视图 — 按机构分组 -->
      <div v-else class="institution-groups">
        <div
          v-for="group in groupedManualAssets"
          :key="group.institution"
          class="glass-card institution-card"
        >
          <!-- 机构头部 -->
          <div class="institution-header">
            <span class="institution-icon">{{ assetTypeConfig[group.assetType]?.icon }}</span>
            <span class="institution-name">{{ group.institution }}</span>
          </div>
          <!-- 资产行 -->
          <div class="institution-assets">
            <div v-for="asset in group.assets" :key="asset.id" class="asset-row">
              <span class="asset-name">{{ asset.name }}</span>
              <span class="asset-currency">{{ asset.currency }}</span>
              <span class="asset-balance font-mono">{{ formatNumber(asset.balance) }}</span>
              <div class="asset-actions">
                <button class="icon-btn" :title="t('accounts.edit')" @click="editAccount(asset)">
                  <PhPencilSimple :size="14" />
                </button>
                <button class="icon-btn danger" :title="t('accounts.delete')" @click="confirmDelete(asset)">
                  <PhTrash :size="14" />
                </button>
              </div>
            </div>
          </div>
          <!-- 多条时显示小计 -->
          <div v-if="group.assets.length > 1" class="institution-total">
            <span>{{ t('accounts.subtotal') }}</span>
            <span class="font-mono">≈ {{ pricingCurrencySymbol }}{{ formatNumber(group.totalUSD) }}</span>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="manualAssets.length === 0" class="glass-card institution-card empty-card">
          <span class="text-muted">{{ t('accounts.noManual') }}</span>
        </div>

        <!-- 添加账户卡片 -->
        <div class="glass-card institution-card add-account-card" @click="openAddDialog('manual')">
          <PhPlus :size="32" weight="light" class="add-icon" />
          <span class="add-label">{{ t('accounts.addAccount') }}</span>
        </div>
      </div>
    </template>

    <!-- 删除确认对话框 -->
    <Transition name="modal">
      <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="showDeleteConfirm = false">
        <div class="modal-content glass-card">
          <div class="modal-header">
            <PhWarning :size="32" weight="fill" class="warning-icon" />
            <h3>{{ t('accounts.confirmDelete') }}</h3>
          </div>
          <p class="modal-body">
            {{ t('accounts.deleteWarning') }} <strong>{{ selectedAccount?.name }}</strong>{{ t('accounts.deleteCannotUndo') }}
          </p>
          <div class="modal-actions">
            <button class="btn btn-secondary btn-sm" @click="showDeleteConfirm = false">
              {{ t('common.cancel') }}
            </button>
            <button class="btn btn-danger btn-sm" @click="deleteAccount">
              {{ t('accounts.confirmDeleteBtn') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 添加/编辑账户对话框 -->
    <AddAccountDialog
      :visible="showAddDialog"
      :accountType="activeTab"
      :editingAccount="editingAccount"
      @close="showAddDialog = false"
      @submitted="handleAccountSubmitted"
    />

    <!-- 批量导入对话框 -->
    <BatchImportDialog
      :visible="showBatchImport"
      @close="showBatchImport = false"
      @imported="handleAccountSubmitted"
    />

    <!-- 钱包连接对话框 -->
    <WalletConnectDialog
      :visible="showWalletConnect"
      @close="showWalletConnect = false"
      @connected="handleWalletConnected"
    />
  </div>
</template>

<style scoped>
.accounts-page {
  display: flex;
  flex-direction: column;
  gap: var(--gap-lg);
}

/* ========== 工具栏 ========== */
.page-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-md);
  flex-wrap: wrap;
}

.tabs {
  display: flex;
  gap: 2px;
  background: var(--color-bg-secondary);
  padding: 2px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  padding: 6px 12px;
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  color: var(--color-text-secondary);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.tab-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.tab-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

.tab-count {
  font-size: 0.6875rem;
  padding: 0 4px;
  opacity: 0.7;
}

.tab-btn.active .tab-count {
  opacity: 1;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.total-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.total-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-right: var(--gap-sm);
}

/* ========== 视图切换 ========== */
.view-toggle {
  display: flex;
  gap: 2px;
  background: var(--color-bg-secondary);
  padding: 2px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
}

.view-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.view-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.view-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

/* ========== 表格面板 ========== */
.table-panel {
  padding: 0;
  overflow-x: auto;
}

.table-panel .table th,
.table-panel .table td {
  white-space: nowrap;
}

/* 名称列 */
.name-cell {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
}

.logo-badge {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-xs);
  font-size: 0.5625rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.name-text {
  font-weight: 500;
}

/* 地址列 */
.address-cell {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.copy-btn {
  width: 20px;
  height: 20px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: var(--radius-xs);
  transition: color var(--transition-fast);
}

.copy-btn:hover {
  color: var(--color-accent-primary);
}

/* 状态 */
.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 4px;
  vertical-align: middle;
}

.status-dot.connected,
.status-dot.active {
  background: var(--color-success);
}

.status-dot.error {
  background: var(--color-error);
}

.status-dot.inactive {
  background: var(--color-text-muted);
}

.status-text {
  font-size: 0.75rem;
}

.status-text.connected,
.status-text.active {
  color: var(--color-success);
}

.status-text.error {
  color: var(--color-error);
}

.text-muted {
  color: var(--color-text-muted);
  font-size: 0.75rem;
}

.note-col {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 操作列 */
.col-actions {
  text-align: right;
  width: 100px;
}

.action-btns {
  display: flex;
  justify-content: flex-end;
  gap: 2px;
  opacity: 0.5;
  transition: opacity var(--transition-fast);
}

tr:hover .action-btns {
  opacity: 1;
}

.icon-btn {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.icon-btn:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
}

.icon-btn.danger:hover {
  background: rgba(226, 92, 92, 0.12);
  color: var(--color-error);
}

.empty-row {
  text-align: center;
  color: var(--color-text-muted);
  padding: var(--gap-xl) !important;
}

/* ========== 卡片网格 ========== */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--gap-md);
}

.account-card {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
  padding: var(--gap-lg);
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap-sm);
}

.card-identity {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  flex: 1;
  min-width: 0;
}

.card-logo {
  width: 32px;
  height: 32px;
  font-size: 0.625rem;
}

.card-name-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  min-width: 0;
}

.card-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-sub {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 卡片主体 */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: var(--gap-sm) 0;
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
}

.card-label {
  font-size: 0.6875rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.card-value {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 卡片元数据 */
.card-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.6875rem;
}

.card-meta .address-cell {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.card-meta .text-muted {
  font-size: 0.6875rem;
}

.card-meta .note-col {
  max-width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 卡片操作按钮 */
.card-actions {
  display: flex;
  gap: 4px;
  padding-top: 4px;
  opacity: 0.5;
  transition: opacity var(--transition-fast);
}

.account-card:hover .card-actions {
  opacity: 1;
}

/* 空状态卡片 */
.empty-card {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 160px;
  text-align: center;
}

/* 添加账户卡片 */
.add-account-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  min-height: 160px;
  cursor: pointer;
  border: 1px dashed var(--color-border);
  background: transparent;
  transition: all var(--transition-fast);
}

.add-account-card:hover {
  border-color: var(--color-accent-primary);
  background: rgba(75, 131, 240, 0.04);
}

.add-icon {
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}

.add-account-card:hover .add-icon {
  color: var(--color-accent-primary);
}

.add-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  transition: color var(--transition-fast);
}

.add-account-card:hover .add-label {
  color: var(--color-accent-primary);
}

/* ========== DeFi 仓位 ========== */
.defi-filter-bar {
  display: flex;
  gap: 2px;
  background: var(--color-bg-secondary);
  padding: 2px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
  align-self: flex-start;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-xs);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.filter-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-tertiary);
}

.filter-btn.active {
  background: var(--color-accent-primary);
  color: #fff;
}

.filter-count {
  font-size: 0.625rem;
  opacity: 0.7;
}

.filter-btn.active .filter-count {
  opacity: 1;
}

.defi-protocol-groups {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.defi-protocol-section {
  padding: var(--gap-lg);
}

.defi-protocol-header {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
  margin-bottom: var(--gap-md);
}

.defi-protocol-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: 0.5625rem;
  font-weight: 700;
  color: var(--color-text-secondary);
}

.defi-protocol-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  flex: 1;
}

.defi-protocol-total {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.defi-positions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--gap-sm);
}

/* ========== 机构分组视图 ========== */
.institution-groups {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

.institution-card {
  padding: var(--gap-lg);
}

.institution-header {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding-bottom: var(--gap-sm);
  border-bottom: 1px solid var(--color-border);
  margin-bottom: var(--gap-sm);
}

.institution-icon {
  font-size: 1rem;
  line-height: 1;
}

.institution-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.institution-assets {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.asset-row {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  padding: var(--gap-xs) 0;
}

.asset-row:hover {
  background: var(--color-bg-tertiary);
  border-radius: var(--radius-xs);
}

.asset-name {
  flex: 1;
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.asset-currency {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  width: 36px;
  text-align: center;
}

.asset-balance {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
  text-align: right;
  min-width: 80px;
}

.asset-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.asset-row:hover .asset-actions {
  opacity: 1;
}

.institution-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--gap-sm);
  margin-top: var(--gap-xs);
  border-top: 1px solid var(--color-border);
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

/* ========== 模态框 ========== */
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
  max-width: 380px;
  padding: var(--gap-xl);
  text-align: center;
}

.modal-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-md);
}

.warning-icon {
  color: var(--color-warning);
}

.modal-header h3 {
  font-size: 0.9375rem;
}

.modal-body {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  margin-bottom: var(--gap-lg);
}

.modal-body strong {
  color: var(--color-text-primary);
}

.modal-actions {
  display: flex;
  gap: var(--gap-sm);
  justify-content: center;
}

/* 模态框动画 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-base);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

/* ========== 响应式 ========== */
@media (max-width: 768px) {
  .page-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-right {
    justify-content: space-between;
  }

  .tab-label {
    display: none;
  }

  /* 移动端隐藏视图切换（强制卡片模式） */
  .view-toggle {
    display: none;
  }

  /* 移动端隐藏表格视图 */
  .table-panel {
    display: none;
  }

  /* 确保卡片网格显示 */
  .card-grid {
    display: grid !important;
  }

  /* 触摸区域优化 */
  .tab-btn {
    min-height: 44px;
    min-width: 44px;
  }

  .view-btn {
    min-width: 44px;
    min-height: 44px;
  }

  .card-action-btn {
    min-width: 44px;
    min-height: 44px;
  }

  .card-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .card-grid {
    grid-template-columns: 1fr;
  }
}
</style>
