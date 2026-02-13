/**
 * 成就系统 Mock 数据
 * 里程碑、坚持、投资三类成就
 */

/**
 * 获取所有成就
 * @returns {Array} 成就列表
 */
export function getAchievements() {
  return [
    // 里程碑成就
    {
      id: 'ach-001',
      name: '初学者',
      description: '首次添加账户',
      condition: '添加至少 1 个账户',
      icon: 'seed',
      category: 'milestone',
      unlocked: true,
      unlockedAt: '2025-01-15',
    },
    {
      id: 'ach-002',
      name: '收藏家',
      description: '管理 5 个以上账户',
      condition: '关联账户数 >= 5',
      icon: 'collection',
      category: 'milestone',
      unlocked: true,
      unlockedAt: '2025-03-20',
    },
    {
      id: 'ach-003',
      name: '整币达人',
      description: '持有 1 个完整 BTC',
      condition: 'BTC 持仓 >= 1.0',
      icon: 'btc',
      category: 'milestone',
      unlocked: false,
      unlockedAt: null,
    },
    {
      id: 'ach-004',
      name: '万元户',
      description: '总资产突破 $10,000',
      condition: '总资产 >= $10,000',
      icon: 'money',
      category: 'milestone',
      unlocked: true,
      unlockedAt: '2025-02-01',
    },
    {
      id: 'ach-005',
      name: '十万大关',
      description: '总资产突破 $100,000',
      condition: '总资产 >= $100,000',
      icon: 'rocket',
      category: 'milestone',
      unlocked: false,
      unlockedAt: null,
    },

    // 坚持成就
    {
      id: 'ach-101',
      name: '早起鸟',
      description: '连续 7 天查看资产',
      condition: '连续登录 7 天',
      icon: 'bird',
      category: 'persistence',
      unlocked: true,
      unlockedAt: '2025-01-22',
    },
    {
      id: 'ach-102',
      name: '数据控',
      description: '导出过 3 次报告',
      condition: '导出报告次数 >= 3',
      icon: 'data',
      category: 'persistence',
      unlocked: true,
      unlockedAt: '2025-04-10',
    },
    {
      id: 'ach-103',
      name: '年度玩家',
      description: '使用 AllFi 满 1 年',
      condition: '注册满 365 天',
      icon: 'calendar',
      category: 'persistence',
      unlocked: false,
      unlockedAt: null,
    },
    {
      id: 'ach-104',
      name: '定投达人',
      description: '设置并坚持定投策略 3 个月',
      condition: '定投策略运行 >= 90 天',
      icon: 'timer',
      category: 'persistence',
      unlocked: false,
      unlockedAt: null,
    },

    // 投资成就
    {
      id: 'ach-201',
      name: '逆势操作',
      description: '在市场下跌时加仓并最终获利',
      condition: '市场下跌期间买入，后续涨幅 > 20%',
      icon: 'trend',
      category: 'investment',
      unlocked: false,
      unlockedAt: null,
    },
    {
      id: 'ach-202',
      name: '跑赢大盘',
      description: '月收益率超过 BTC',
      condition: '当月收益率 > BTC 收益率',
      icon: 'crown',
      category: 'investment',
      unlocked: true,
      unlockedAt: '2025-03-31',
    },
    {
      id: 'ach-203',
      name: '稳如老狗',
      description: '连续 3 个月正收益',
      condition: '连续 3 个月收益率 > 0%',
      icon: 'shield',
      category: 'investment',
      unlocked: true,
      unlockedAt: '2025-05-31',
    },
    {
      id: 'ach-204',
      name: '分散专家',
      description: '持仓 HHI 指数低于 1500',
      condition: '资产集中度 HHI < 1500',
      icon: 'scatter',
      category: 'investment',
      unlocked: false,
      unlockedAt: null,
    },
  ]
}
