/**
 * 风险管理 API 服务
 * 提供风险指标、回撤分析、Beta 系数等数据
 */

import { get } from "./client.js";
import * as mockData from "./mockData.js";

// 是否使用 Mock 数据
const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== "false";

// 模拟网络延迟
const simulateDelay = (ms = 500) =>
  new Promise((resolve) => setTimeout(resolve, ms));

/**
 * 解析 period 字符串为天数
 * @param {string} period - 时间周期 (7d/30d/90d/1y)
 * @returns {number} 天数
 */
const parsePeriodToDays = (period) => {
  const match = period.match(/^(\d+)([dmy])$/);
  if (!match) return 30;
  const value = parseInt(match[1], 10);
  const unit = match[2];
  switch (unit) {
    case "d":
      return value;
    case "m":
      return value * 30;
    case "y":
      return value * 365;
    default:
      return 30;
  }
};

export const riskService = {
  /**
   * 获取风险总览数据
   * @param {string} period - 时间周期 (7d/30d/90d/1y)
   * @returns {Promise<Object>} 风险总览数据
   */
  async getRiskOverview(period = "30d") {
    if (USE_MOCK) {
      await simulateDelay(300);
      return mockData.getRiskOverview(period);
    }
    // 后端 /risk/overview 不接受 period 参数，返回最新指标
    const result = await get("/risk/overview");
    const metrics = result.metrics || result;

    // 转换为组件期望的格式，补充缺失字段
    return {
      // 直接使用后端返回的字段
      volatility: metrics.volatility ?? 0,
      sharpe_ratio: metrics.sharpe_ratio ?? 0,
      max_drawdown: metrics.max_drawdown ?? 0,
      var_95: metrics.var_95 ?? 0,
      var_99: metrics.var_99 ?? 0,
      // 后端只有 beta（相对 BTC），映射到 beta_btc
      beta_btc: metrics.beta ?? 0,
      // 后端暂无 ETH Beta 和 CVaR，使用默认值
      beta_eth: 0,
      cvar_95: metrics.var_99 ?? metrics.var_95 ?? 0, // 用 VaR 99% 近似 CVaR
      // 保留原始数据
      ...metrics,
    };
  },

  /**
   * 获取风险指标历史趋势
   * @param {string} period - 时间周期
   * @returns {Promise<Object>} 风险指标历史数据
   */
  async getRiskMetrics(period = "30d") {
    if (USE_MOCK) {
      await simulateDelay(400);
      return mockData.getRiskMetrics(period);
    }
    const days = parsePeriodToDays(period);
    const result = await get(`/risk/metrics?days=${days}`);
    const history = result.history || result;

    // 转换为 RiskMetricsChart 组件期望的格式
    if (!Array.isArray(history) || history.length === 0) {
      return { labels: [], volatility: [], sharpe: [] };
    }

    // 按日期升序排列
    const sorted = [...history].sort(
      (a, b) => new Date(a.metric_date) - new Date(b.metric_date),
    );

    return {
      labels: sorted.map((item) => item.metric_date),
      volatility: sorted.map((item) => item.volatility ?? 0),
      sharpe: sorted.map((item) => item.sharpe_ratio ?? 0),
    };
  },

  /**
   * 获取回撤曲线数据
   * @param {string} period - 时间周期
   * @returns {Promise<Object>} 回撤曲线数据
   */
  async getDrawdown(period = "30d") {
    if (USE_MOCK) {
      await simulateDelay(350);
      return mockData.getDrawdown(period);
    }
    // 后端没有独立的 drawdown 接口，从 metrics 历史中提取
    const days = parsePeriodToDays(period);
    const result = await get(`/risk/metrics?days=${days}`);
    const history = result.history || result;

    // 转换为前端期望的 drawdown 格式
    // 从历史数据中提取 max_drawdown 字段构建回撤曲线
    const drawdownData = {
      labels: [],
      drawdown: [],
      maxDrawdown: 0,
      maxDrawdownDate: null,
    };

    if (Array.isArray(history)) {
      history.forEach((item) => {
        drawdownData.labels.push(item.metric_date);
        drawdownData.drawdown.push(item.max_drawdown);
        if (item.max_drawdown > drawdownData.maxDrawdown) {
          drawdownData.maxDrawdown = item.max_drawdown;
          drawdownData.maxDrawdownDate = item.metric_date;
        }
      });
      // 反转使日期升序
      drawdownData.labels.reverse();
      drawdownData.drawdown.reverse();
    }

    return drawdownData;
  },

  /**
   * 获取 Beta 系数对比数据
   * @param {string} period - 时间周期
   * @returns {Promise<Object>} Beta 系数数据
   */
  async getBeta(period = "30d") {
    if (USE_MOCK) {
      await simulateDelay(300);
      return mockData.getBeta(period);
    }
    // 后端没有独立的 beta 接口，从 overview 中获取最新 beta
    const result = await get("/risk/overview");
    const metrics = result.metrics || result;
    const portfolioBeta = metrics.beta ?? 0;

    // 转换为 BetaComparisonCard 组件期望的格式
    return {
      portfolio_beta: portfolioBeta,
      // 基准对比列表
      benchmarks: [
        {
          name: "BTC",
          beta: 1.0,
          correlation: 1.0, // BTC 与自身相关性为 100%
        },
        {
          name: "投资组合",
          beta: portfolioBeta,
          correlation: portfolioBeta > 0 ? 0.8 : 0.5, // 估算相关性
        },
      ],
    };
  },

  /**
   * 获取风险预警列表
   * @returns {Promise<Array>} 风险预警列表
   */
  async getRiskAlerts() {
    if (USE_MOCK) {
      await simulateDelay(250);
      return mockData.getRiskAlerts();
    }
    // 后端暂无 alerts 接口，返回空数组
    // TODO: 后续实现风险预警功能
    return [];
  },
};
