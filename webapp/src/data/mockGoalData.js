/**
 * 目标追踪 Mock 数据
 * 用于开发和演示
 */

export const mockGoals = [
  {
    id: '1',
    title: '总资产达到 $300,000',
    type: 'asset_value',
    targetValue: 300000,
    currency: 'USDC',
    deadline: '2026-12-31',
    createdAt: '2026-01-15T10:00:00Z',
  },
  {
    id: '2',
    title: '持有 1 个 BTC',
    type: 'holding_amount',
    targetValue: 1,
    currency: 'BTC',
    deadline: '2026-06-30',
    createdAt: '2026-02-01T08:00:00Z',
  },
  {
    id: '3',
    title: '年化收益率达到 20%',
    type: 'return_rate',
    targetValue: 20,
    currency: 'USDC',
    deadline: null,
    createdAt: '2026-02-10T14:00:00Z',
  },
];

export default mockGoals;
