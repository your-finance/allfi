/**
 * 命令面板 Store
 * 管理 Cmd+K 命令面板的状态、命令注册、资产搜索和快捷筛选
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAssetStore } from './assetStore'
import { useNFTStore } from './nftStore'
import { useStrategyStore } from './strategyStore'

export const useCommandStore = defineStore('command', () => {
  // 面板状态
  const isOpen = ref(false)
  const searchQuery = ref('')
  const selectedIndex = ref(0)

  // 快捷命令
  const commands = ref([
    { id: 'refresh', label: '刷新所有资产', labelEn: 'Refresh all assets', icon: 'ArrowsClockwise', category: 'action' },
    { id: 'addExchange', label: '添加 CEX 账户', labelEn: 'Add CEX account', icon: 'Plus', category: 'action' },
    { id: 'addWallet', label: '添加钱包地址', labelEn: 'Add wallet address', icon: 'Wallet', category: 'action' },
    { id: 'toggleTheme', label: '切换主题', labelEn: 'Switch theme', icon: 'PaintBrush', category: 'action' },
    { id: 'togglePrivacy', label: '切换隐私模式', labelEn: 'Toggle privacy mode', icon: 'EyeSlash', category: 'action' },
  ])

  // 页面跳转命令
  const pages = ref([
    { id: 'dashboard', label: '仪表盘', labelEn: 'Dashboard', path: '/', icon: 'ChartPieSlice' },
    { id: 'accounts', label: '账户管理', labelEn: 'Accounts', path: '/accounts', icon: 'Wallet' },
    { id: 'history', label: '历史记录', labelEn: 'History', path: '/history', icon: 'ClockCounterClockwise' },
    { id: 'analytics', label: '数据分析', labelEn: 'Analytics', path: '/analytics', icon: 'ChartLine' },
    { id: 'settings', label: '设置', labelEn: 'Settings', path: '/settings', icon: 'Gear' },
  ])

  /**
   * 从 assetStore 获取资产搜索源数据
   * 汇总 CEX、链上、手动资产和 DeFi 仓位的所有持仓
   */
  const assetSearchItems = computed(() => {
    let assetStoreRef
    try {
      assetStoreRef = useAssetStore()
    } catch {
      return []
    }

    const items = []

    // CEX 账户持仓
    for (const acc of (assetStoreRef.cexAccounts || [])) {
      if (acc.holdings) {
        for (const h of acc.holdings) {
          items.push({
            id: `cex-${acc.id}-${h.symbol}`,
            type: 'asset',
            symbol: h.symbol,
            name: h.name || h.symbol,
            label: `${h.symbol} — ${h.name || h.symbol}`,
            labelEn: `${h.symbol} — ${h.name || h.symbol}`,
            icon: 'CurrencyBtc',
            balance: h.balance,
            value: h.value,
            source: acc.name,
            sourceType: 'cex',
            exchange: acc.exchange
          })
        }
      }
    }

    // 链上钱包持仓
    for (const w of (assetStoreRef.walletAddresses || [])) {
      if (w.holdings) {
        for (const h of w.holdings) {
          items.push({
            id: `chain-${w.id}-${h.symbol}`,
            type: 'asset',
            symbol: h.symbol,
            name: h.name || h.symbol,
            label: `${h.symbol} — ${h.name || h.symbol}`,
            labelEn: `${h.symbol} — ${h.name || h.symbol}`,
            icon: 'CurrencyBtc',
            balance: h.balance,
            value: h.value,
            source: w.name,
            sourceType: 'blockchain',
            blockchain: w.blockchain
          })
        }
      }
    }

    // DeFi 仓位
    for (const p of (assetStoreRef.defiPositions || [])) {
      items.push({
        id: `defi-${p.id}`,
        type: 'asset',
        symbol: p.tokens?.map(t => t.symbol).join('/') || p.protocol,
        name: `${p.protocol} ${p.type.toUpperCase()}`,
        label: `${p.protocol} — ${p.tokens?.map(t => t.symbol).join('/')}`,
        labelEn: `${p.protocol} — ${p.tokens?.map(t => t.symbol).join('/')}`,
        icon: 'CubeTransparent',
        balance: 0,
        value: p.valueUSD,
        source: p.protocol,
        sourceType: 'defi'
      })
    }

    // NFT 收藏集
    try {
      const nftStore = useNFTStore()
      for (const col of (nftStore.collections || [])) {
        items.push({
          id: `nft-${col.slug || col.name}`,
          type: 'asset',
          symbol: col.name,
          name: `NFT — ${col.name}`,
          label: `${col.name} (${col.count} NFTs)`,
          labelEn: `${col.name} (${col.count} NFTs)`,
          icon: 'ImageSquare',
          balance: col.count,
          value: col.totalFloorUSD || 0,
          source: col.chain || 'Ethereum',
          sourceType: 'nft'
        })
      }
    } catch {}

    // 策略
    try {
      const strategyStore = useStrategyStore()
      for (const s of (strategyStore.strategies || [])) {
        const typeLabel = { rebalance: '再平衡', dca: '定投', alert: '止盈止损' }
        items.push({
          id: `strategy-${s.id}`,
          type: 'asset',
          symbol: s.name || s.type,
          name: `${typeLabel[s.type] || s.type} — ${s.name}`,
          label: `${typeLabel[s.type] || s.type} — ${s.name}`,
          labelEn: `${s.type} — ${s.name}`,
          icon: 'Strategy',
          balance: 0,
          value: 0,
          source: typeLabel[s.type] || s.type,
          sourceType: 'strategy'
        })
      }
    } catch {}

    return items
  })

  /**
   * 解析搜索查询中的快捷筛选前缀
   * @param {string} query - 原始查询
   * @returns {{ prefix: string|null, keyword: string }}
   */
  function parseSearchPrefix(query) {
    const trimmed = query.trim()

    // @binance → 按交易所/来源筛选
    if (trimmed.startsWith('@')) {
      return { prefix: 'source', keyword: trimmed.slice(1).toLowerCase() }
    }
    // #defi → 按类型筛选 (cex/blockchain/defi/manual)
    if (trimmed.startsWith('#')) {
      return { prefix: 'type', keyword: trimmed.slice(1).toLowerCase() }
    }
    // >1000 → 按价值筛选
    if (trimmed.startsWith('>')) {
      const num = parseFloat(trimmed.slice(1))
      if (!isNaN(num)) {
        return { prefix: 'minValue', keyword: String(num) }
      }
    }

    return { prefix: null, keyword: trimmed.toLowerCase() }
  }

  // 搜索过滤结果（分组展示）
  const filteredItems = computed(() => {
    const rawQuery = searchQuery.value
    const { prefix, keyword } = parseSearchPrefix(rawQuery)

    // 如果有快捷前缀，只搜索资产
    if (prefix) {
      let assets = assetSearchItems.value

      if (prefix === 'source') {
        assets = assets.filter(a =>
          a.source.toLowerCase().includes(keyword) ||
          (a.exchange || '').toLowerCase().includes(keyword)
        )
      } else if (prefix === 'type') {
        assets = assets.filter(a => a.sourceType.includes(keyword))
      } else if (prefix === 'minValue') {
        const minVal = parseFloat(keyword)
        assets = assets.filter(a => a.value >= minVal)
      }

      return assets.map(a => ({ ...a, type: 'asset' }))
    }

    // 普通搜索：合并命令 + 页面 + 资产
    const commandItems = commands.value.map(c => ({ ...c, type: 'command' }))
    const pageItems = pages.value.map(p => ({ ...p, type: 'page' }))

    if (!keyword) {
      // 无关键词时只显示命令和页面
      return [...commandItems, ...pageItems]
    }

    // 有关键词时搜索全部
    const filteredCommands = commandItems.filter(item =>
      item.label.toLowerCase().includes(keyword) ||
      item.labelEn.toLowerCase().includes(keyword) ||
      item.id.toLowerCase().includes(keyword)
    )

    const filteredPages = pageItems.filter(item =>
      item.label.toLowerCase().includes(keyword) ||
      item.labelEn.toLowerCase().includes(keyword) ||
      item.id.toLowerCase().includes(keyword)
    )

    // 资产模糊搜索：按代币名称、符号、来源
    const filteredAssets = assetSearchItems.value.filter(a =>
      a.symbol.toLowerCase().includes(keyword) ||
      a.name.toLowerCase().includes(keyword) ||
      a.source.toLowerCase().includes(keyword)
    ).slice(0, 10) // 最多展示 10 个资产结果

    return [...filteredCommands, ...filteredPages, ...filteredAssets]
  })

  // 操作
  function toggle() {
    isOpen.value = !isOpen.value
    if (isOpen.value) {
      searchQuery.value = ''
      selectedIndex.value = 0
    }
  }

  function open() {
    isOpen.value = true
    searchQuery.value = ''
    selectedIndex.value = 0
  }

  function close() {
    isOpen.value = false
    searchQuery.value = ''
    selectedIndex.value = 0
  }

  // 键盘导航
  function moveUp() {
    if (selectedIndex.value > 0) {
      selectedIndex.value--
    }
  }

  function moveDown() {
    if (selectedIndex.value < filteredItems.value.length - 1) {
      selectedIndex.value++
    }
  }

  // 获取当前选中项
  function getSelectedItem() {
    return filteredItems.value[selectedIndex.value] || null
  }

  return {
    isOpen,
    searchQuery,
    selectedIndex,
    commands,
    pages,
    filteredItems,
    assetSearchItems,
    toggle,
    open,
    close,
    moveUp,
    moveDown,
    getSelectedItem,
  }
})
