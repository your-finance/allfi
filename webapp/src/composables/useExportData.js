/**
 * 数据导出组合式函数
 * 支持将资产数据导出为 CSV 或 JSON 格式
 */
import { useAssetStore } from '../stores/assetStore'
import { useTransactionStore } from '../stores/transactionStore'
import { useI18n } from './useI18n'

export function useExportData() {
  const assetStore = useAssetStore()
  const { t } = useI18n()

  /**
   * 收集所有持仓数据（合并后的）
   * 返回统一的数据结构供导出使用
   */
  function collectHoldingsData() {
    const holdings = []
    const categories = [
      { id: 'cex', accounts: assetStore.cexAccounts },
      { id: 'blockchain', accounts: assetStore.walletAddresses },
      { id: 'manual', accounts: assetStore.manualAssets }
    ]

    for (const category of categories) {
      for (const account of category.accounts) {
        if (account.holdings) {
          for (const h of account.holdings) {
            holdings.push({
              symbol: h.symbol,
              name: h.name,
              source: account.name,
              sourceType: category.id,
              price: h.price,
              change24h: h.change24h,
              balance: h.balance,
              value: h.value,
              percentage: assetStore.totalValue
                ? ((h.value / assetStore.totalValue) * 100).toFixed(2)
                : '0.00'
            })
          }
        }
      }
    }

    return holdings
  }

  /**
   * 导出为 CSV 格式
   * @param {string} filename - 文件名（不含扩展名）
   */
  function exportAsCSV(filename) {
    const holdings = collectHoldingsData()

    // CSV 表头
    const headers = [
      t('dashboard.asset'),
      'Symbol',
      t('dashboard.sources'),
      t('export.sourceType'),
      t('dashboard.price'),
      t('dashboard.change'),
      t('dashboard.balance'),
      t('dashboard.value'),
      t('dashboard.percentage')
    ]

    // 转换为 CSV 行
    const rows = holdings.map(h => [
      `"${h.name}"`,
      h.symbol,
      `"${h.source}"`,
      h.sourceType,
      h.price,
      h.change24h,
      h.balance,
      h.value,
      `${h.percentage}%`
    ])

    // 添加汇总行
    rows.push([])
    rows.push([`"${t('dashboard.totalAssets')}"`, '', '', '', '', '', '', assetStore.totalValue, '100%'])

    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.join(','))
    ].join('\n')

    // 添加 BOM 以支持 Excel 正确显示中文
    const bom = '\uFEFF'
    downloadFile(bom + csvContent, `${filename}.csv`, 'text/csv;charset=utf-8')
  }

  /**
   * 导出为 JSON 格式
   * @param {string} filename - 文件名（不含扩展名）
   */
  function exportAsJSON(filename) {
    const holdings = collectHoldingsData()

    const exportData = {
      exportedAt: new Date().toISOString(),
      totalValue: assetStore.totalValue,
      change24h: assetStore.change24h,
      holdings
    }

    const jsonContent = JSON.stringify(exportData, null, 2)
    downloadFile(jsonContent, `${filename}.json`, 'application/json')
  }

  /**
   * 触发文件下载
   */
  function downloadFile(content, filename, mimeType) {
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  /**
   * 生成默认文件名
   */
  function getDefaultFilename() {
    const date = new Date().toISOString().split('T')[0]
    return `allfi-export-${date}`
  }

  /**
   * 导出交易记录为 CSV 格式
   * @param {string} filename - 文件名（不含扩展名）
   */
  function exportTransactionsAsCSV(filename) {
    const txStore = useTransactionStore()
    const txList = txStore.transactions

    // CSV 表头
    const headers = [
      t('transaction.date'),
      t('transaction.time'),
      t('transaction.txType'),
      t('transaction.fromAsset'),
      t('transaction.fromAmount'),
      t('transaction.toAsset'),
      t('transaction.toAmount'),
      t('transaction.fee'),
      t('transaction.feeCurrency'),
      t('transaction.source'),
      t('transaction.sourceType'),
      t('transaction.chain'),
      t('transaction.note'),
    ]

    const rows = txList.map(tx => {
      const d = new Date(tx.timestamp)
      const dateStr = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
      const timeStr = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
      return [
        dateStr,
        timeStr,
        tx.type,
        tx.from.symbol,
        tx.from.amount,
        tx.to.symbol,
        tx.to.amount,
        tx.fee?.amount || 0,
        tx.fee?.currency || '',
        `"${tx.source}"`,
        tx.sourceType,
        tx.chain || '',
        `"${tx.note || ''}"`,
      ]
    })

    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.join(','))
    ].join('\n')

    const bom = '\uFEFF'
    downloadFile(bom + csvContent, `${filename}.csv`, 'text/csv;charset=utf-8')
  }

  return {
    exportAsCSV,
    exportAsJSON,
    exportTransactionsAsCSV,
    getDefaultFilename
  }
}
