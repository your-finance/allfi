/**
 * 年度报告 API 服务
 * 管理年度报告数据的获取，支持 Mock/Real 自动切换
 */

import { get } from "./client.js";
import * as mockData from "../data/mockAnnualReportData.js";

const USE_MOCK = import.meta.env.VITE_USE_MOCK_API !== "false";
const simulateDelay = (ms = 500) =>
  new Promise((resolve) => setTimeout(resolve, ms));

export const annualReportService = {
  /**
   * 获取年度报告
   * @param {number} year - 年份
   * @returns {Promise<Object|null>} 返回报告数据或 null（如果没有数据）
   */
  async getAnnualReport(year) {
    if (USE_MOCK) {
      await simulateDelay(400);
      return mockData.getAnnualReport(year);
    }
    const response = await get(`/reports/annual/${year}`);
    // 后端返回 { report: {...} } 或 { report: null }
    return response;
  },
};
