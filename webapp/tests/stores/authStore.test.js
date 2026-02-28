/**
 * authStore 单元测试
 * 测试认证状态管理的核心功能
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useAuthStore } from "@/stores/authStore";

// 使用 importOriginal 进行部分 mock
vi.mock("@/api/index", async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    authService: {
      ...actual.authService,
      getStatus: vi.fn(),
      setupPIN: vi.fn(),
      login: vi.fn(),
      changePIN: vi.fn(),
      setup2FA: vi.fn(),
      enable2FA: vi.fn(),
      disable2FA: vi.fn(),
    },
  };
});

describe("authStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    localStorage.clear();
  });

  it("初始状态应该正确设置", () => {
    const store = useAuthStore();
    expect(store.isAuthenticated).toBe(false);
    expect(store.pinSet).toBe(false);
    expect(store.token).toBe(null);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBe(null);
  });

  it("checkAuthStatus 应该更新 pinSet 状态", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = { pin_set: true };
    vi.mocked(authService.getStatus).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.checkAuthStatus();

    expect(result).toBe(true);
    expect(store.pinSet).toBe(true);
    expect(authService.getStatus).toHaveBeenCalledOnce();
  });

  it("setupPIN 应该成功设置 PIN 并更新状态", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = { token: "test-jwt-token" };
    vi.mocked(authService.setupPIN).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.setupPIN("1234");

    expect(result).toBe(true);
    expect(store.isAuthenticated).toBe(true);
    expect(store.pinSet).toBe(true);
    expect(store.token).toBe("test-jwt-token");
    expect(localStorage.getItem("allfi-auth")).toBeTruthy();
  });

  it("login 应该成功登录并返回 token", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = { token: "login-jwt-token" };
    vi.mocked(authService.login).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.login("1234");

    expect(result).toBe(true);
    expect(store.isAuthenticated).toBe(true);
    expect(store.token).toBe("login-jwt-token");
  });

  it("logout 应该清除所有状态", () => {
    localStorage.setItem("allfi-auth", "test-token");

    const store = useAuthStore();
    store.logout();

    expect(store.isAuthenticated).toBe(false);
    expect(store.token).toBe(null);
    expect(localStorage.getItem("allfi-auth")).toBeNull();
  });

  it("setup2FA 应该获取密钥和二维码 URL", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = {
      secret: "JBSWY3DPEHPK3PXP",
      qr_url: "otpauth://totp/...",
    };
    vi.mocked(authService.setup2FA).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.setup2FA();

    expect(result).toEqual(mockResult);
    expect(authService.setup2FA).toHaveBeenCalledOnce();
  });

  it("enable2FA 应该成功启用 2FA", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = { success: true };
    vi.mocked(authService.enable2FA).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.enable2FA("123456");

    expect(result).toBe(true);
    expect(authService.enable2FA).toHaveBeenCalledWith("123456");
  });

  it("disable2FA 应该成功禁用 2FA", async () => {
    const { authService } = await import("@/api/index");
    const mockResult = { success: true };
    vi.mocked(authService.disable2FA).mockResolvedValue(mockResult);

    const store = useAuthStore();
    const result = await store.disable2FA("123456");

    expect(result).toBe(true);
    expect(authService.disable2FA).toHaveBeenCalledWith("123456");
  });
});
