/**
 * CryptoIcon 组件单元测试
 * 测试加密货币图标组件的基本功能
 */
import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import CryptoIcon from "@/components/CryptoIcon.vue";

describe("CryptoIcon 组件", () => {
  it("应该正确渲染默认图标", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "BTC",
      },
    });
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.text()).toBe("BTC");
  });

  it("应该显示正确的符号（最多3个字符）", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "ETH",
      },
    });
    expect(wrapper.text()).toBe("ETH");
  });

  it("应该正确显示缩写（USDC显示USD，因组件只显示前3个字符）", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "USDC",
      },
    });
    // 组件使用 slice(0, 3)，所以 USDC 显示为 USD
    expect(wrapper.text()).toBe("USD");
  });

  it("长符号应该截断为前3个字符", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "BITCOIN",
      },
    });
    expect(wrapper.text()).toBe("BIT");
  });

  it("应该支持不同尺寸", () => {
    const sizes = ["sm", "md", "lg"];
    sizes.forEach((size) => {
      const wrapper = mount(CryptoIcon, {
        props: {
          symbol: "USDC",
          size: size,
        },
      });
      const sizeClass = size === "md" ? "w-10" : size === "lg" ? "w-12" : "w-6";
      expect(wrapper.classes()).toContain(sizeClass);
    });
  });

  it("应该处理未知符号（使用默认颜色）", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "UNKNOWN",
      },
    });
    expect(wrapper.exists()).toBe(true);
    const div = wrapper.find("div");
    if (div.exists()) {
      // 颜色会被转换为 rgb，所以检查是否包含默认颜色的 rgb 值
      expect(div.attributes("style")).toContain("rgb(100, 116, 139)");
    }
  });

  it("应该正确映射 BTC 颜色", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "BTC",
      },
    });
    const div = wrapper.find("div");
    if (div.exists()) {
      expect(div.attributes("style")).toContain("rgb(247, 147, 26)");
    }
  });

  it("应该正确映射 ETH 颜色", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "ETH",
      },
    });
    const div = wrapper.find("div");
    if (div.exists()) {
      expect(div.attributes("style")).toContain("rgb(98, 126, 234)");
    }
  });

  it("应该小写转大写", () => {
    const wrapper = mount(CryptoIcon, {
      props: {
        symbol: "btc",
      },
    });
    expect(wrapper.text()).toBe("BTC");
  });
});
