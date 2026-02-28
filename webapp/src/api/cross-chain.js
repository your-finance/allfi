// 跨链交易 API
export const crossChainApi = {
  // 获取跨链交易列表
  getTransactions: (params = {}) => {
    return client.get('/cross-chain/transactions', { params })
  },

  // 获取资产流向数据
  getAssetFlow: (params = {}) => {
    return client.get('/cross-chain/flow', { params })
  },

  // 获取跨链手续费统计
  getFeeStats: (params = {}) => {
    return client.get('/cross-chain/fees', { params })
  },

  // 获取支持的跨链桥列表
  getBridges: () => {
    return client.get('/cross-chain/bridges')
  }
}
