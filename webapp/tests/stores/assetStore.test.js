/**
 * assetStore 单元测试
 * 测试资产状态管理功能
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useAssetStore } from "@/stores/assetStore";

// 使用 importOriginal 进行部分 mock
vi.mock("@/api/index", async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    assetService: {
      ...actual.assetService,
      getSummary: vi.fn(),
      getHistory: vi.fn(),
    },
  };
});

describe("assetStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("初始状态应该正确设置", () => {
    const store = useAssetStore();
    expect(store.summary).toBeNull();
    expect(store.summaryLoading).toBe(false);
    expect(store.historyData).toBeNull();
    expect(store.historyLoading).toBe(false);
    expect(store.isRefreshing).toBe(false);
    expect(store.currentTimeRange).toBe("30D");
  });

  it("loadSummary 应该获取资产总览", async () => {
    const { assetService } = await import("@/api/index");
    const mockSummary = {
      totalValue: 100000,
      exchangeValue: 50000,
      blockchainValue: 30000,
      manualValue: 20000,
      change24h: 5000,
      changeValue: 2500,
    };
    vi.mocked(assetService.getSummary).mockResolvedValue(mockSummary);

    const store = useAssetStore();
    await store.loadSummary();

    expect(assetService.getSummary).toHaveBeenCalledOnce();
    expect(store.totalValue).toBe(100000);
  });

  it("loadHistory 应该获取历史数据", async () => {
    const { assetService } = await import("@/api/index");
    const mockHistory = [
      { date: "2026-02-27", value: 95000 },
      { date: "2026-02-26", value: 94000 },
    ];
    vi.mocked(assetService.getHistory).mockResolvedValue(mockHistory);

    const store = useAssetStore();
    await store.loadHistory("30D");

    expect(store.currentTimeRange).toBe("30D");
    expect(assetService.getHistory).toHaveBeenCalledWith("30D");
  });

  it("totalValue 计算属性应该正确工作", () => {
    const store = useAssetStore();
    expect(store.totalValue).toBe(0); // summary 为 null 时返回 0
  });

  it("change24h 计算属性应该正确工作", () => {
    const store = useAssetStore();
    expect(store.change24h).toBe(0); // summary 为 null 时返回 0
  });

  it("manualAssets 计算属性应该返回空数组", () => {
    const store = useAssetStore();
    expect(store.manualAssets).toEqual([]); // summary 为 null 时返回空数组
  });

  it("isRefreshing 初始状态应该为 false", () => {
    const store = useAssetStore();
    expect(store.isRefreshing).toBe(false);
  });

  it("summaryLoading 初始状态应该为 false", () => {
    const store = useAssetStore();
    expect(store.summaryLoading).toBe(false);
  });
});
