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
    const result = await get(`/risk/overview?period=${period}`);
    return result;
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
    const result = await get(`/risk/metrics?period=${period}`);
    return result;
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
    const result = await get(`/risk/drawdown?period=${period}`);
    return result;
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
    const result = await get(`/risk/beta?period=${period}`);
    return result;
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
    const result = await get("/risk/alerts");
    return result;
  },
};
