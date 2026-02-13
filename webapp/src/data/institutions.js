/**
 * 传统资产机构预置数据
 * 按类型分组：银行 / 券商 / 基金平台
 * 每类按地区分组，末尾支持自定义输入
 */

// ========== 银行列表（按地区分组） ==========
export const bankInstitutions = [
  {
    region: 'china',
    items: [
      { id: 'icbc', name: '工商银行', nameEn: 'ICBC' },
      { id: 'ccb', name: '建设银行', nameEn: 'CCB' },
      { id: 'abc', name: '农业银行', nameEn: 'ABC' },
      { id: 'boc', name: '中国银行', nameEn: 'BOC' },
      { id: 'bocom', name: '交通银行', nameEn: 'BOCOM' },
      { id: 'cmb', name: '招商银行', nameEn: 'CMB' },
      { id: 'spdb', name: '浦发银行', nameEn: 'SPDB' },
      { id: 'citic', name: '中信银行', nameEn: 'CITIC' },
      { id: 'cmbc', name: '民生银行', nameEn: 'CMBC' },
      { id: 'cib', name: '兴业银行', nameEn: 'CIB' },
      { id: 'pab', name: '平安银行', nameEn: 'PAB' },
    ]
  },
  {
    region: 'hongkong',
    items: [
      { id: 'hsbc', name: 'HSBC 汇丰银行', nameEn: 'HSBC' },
      { id: 'hangseng', name: '恒生银行', nameEn: 'Hang Seng' },
      { id: 'bochk', name: '中银香港', nameEn: 'BOCHK' },
      { id: 'sc', name: '渣打银行', nameEn: 'Standard Chartered' },
      { id: 'bea', name: '东亚银行', nameEn: 'BEA' },
    ]
  },
  {
    region: 'singapore',
    items: [
      { id: 'dbs', name: 'DBS 星展银行', nameEn: 'DBS' },
      { id: 'ocbc', name: 'OCBC 华侨银行', nameEn: 'OCBC' },
      { id: 'uob', name: 'UOB 大华银行', nameEn: 'UOB' },
    ]
  },
  {
    region: 'usa',
    items: [
      { id: 'jpmorgan', name: 'JPMorgan Chase 摩根大通', nameEn: 'JPMorgan Chase' },
      { id: 'boa', name: 'Bank of America 美国银行', nameEn: 'Bank of America' },
      { id: 'wellsfargo', name: 'Wells Fargo 富国银行', nameEn: 'Wells Fargo' },
      { id: 'citibank', name: 'Citibank 花旗银行', nameEn: 'Citibank' },
    ]
  },
  {
    region: 'europe',
    items: [
      { id: 'deutsche', name: 'Deutsche Bank 德意志银行', nameEn: 'Deutsche Bank' },
      { id: 'barclays', name: 'Barclays 巴克莱银行', nameEn: 'Barclays' },
      { id: 'ubs', name: 'UBS 瑞银集团', nameEn: 'UBS' },
      { id: 'bnp', name: 'BNP Paribas 法国巴黎银行', nameEn: 'BNP Paribas' },
    ]
  }
]

// ========== 券商列表 ==========
export const stockInstitutions = [
  {
    region: 'global',
    items: [
      { id: 'ibkr', name: '盈透证券', nameEn: 'Interactive Brokers' },
      { id: 'futu', name: '富途证券', nameEn: 'Futu / MooMoo' },
      { id: 'tiger', name: '老虎证券', nameEn: 'Tiger Brokers' },
      { id: 'schwab', name: '嘉信理财', nameEn: 'Charles Schwab' },
      { id: 'huatai', name: '华泰证券', nameEn: 'Huatai Securities' },
      { id: 'citic_sec', name: '中信证券', nameEn: 'CITIC Securities' },
      { id: 'gtja', name: '国泰君安', nameEn: 'Guotai Junan' },
      { id: 'eastmoney', name: '东方财富', nameEn: 'East Money' },
    ]
  }
]

// ========== 基金平台列表 ==========
export const fundInstitutions = [
  {
    region: 'global',
    items: [
      { id: 'ttfund', name: '天天基金', nameEn: 'Tiantian Fund' },
      { id: 'antwealth', name: '蚂蚁财富', nameEn: 'Ant Wealth' },
      { id: 'vanguard', name: 'Vanguard 先锋集团', nameEn: 'Vanguard' },
      { id: 'blackrock', name: 'BlackRock / iShares 贝莱德', nameEn: 'BlackRock' },
      { id: 'fidelity', name: 'Fidelity 富达', nameEn: 'Fidelity' },
      { id: 'huaxia', name: '华夏基金', nameEn: 'China AMC' },
      { id: 'efunds', name: '易方达', nameEn: 'E Fund' },
      { id: 'harvest', name: '嘉实基金', nameEn: 'Harvest Fund' },
    ]
  }
]

/**
 * 根据资产类型返回对应的机构列表
 * @param {string} assetType - 资产类型 (bank/stock/fund)
 * @returns {Array} 分组后的机构列表
 */
export function getInstitutionsByType(assetType) {
  switch (assetType) {
    case 'bank': return bankInstitutions
    case 'stock': return stockInstitutions
    case 'fund': return fundInstitutions
    default: return []
  }
}

/**
 * 在机构列表中搜索匹配项
 * @param {Array} groups - 分组后的机构列表
 * @param {string} query - 搜索关键词
 * @returns {Array} 过滤后的分组列表（空组自动移除）
 */
export function searchInstitutions(groups, query) {
  if (!query || !query.trim()) return groups
  const q = query.trim().toLowerCase()
  return groups
    .map(group => ({
      ...group,
      items: group.items.filter(inst =>
        inst.name.toLowerCase().includes(q) ||
        inst.nameEn.toLowerCase().includes(q) ||
        inst.id.toLowerCase().includes(q)
      )
    }))
    .filter(group => group.items.length > 0)
}
