// Package provider Provider Manager 单元测试
package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProviderManager(t *testing.T) {
	manager := NewProviderManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.providers)
	assert.NotNil(t, manager.lastErrors)
	assert.NotNil(t, manager.healthCache)

	// 验证已注册的 Provider 数量
	assert.Equal(t, 4, len(manager.providers), "应该注册4个Provider")
}

func TestProviderManager_RegisterProvider(t *testing.T) {
	manager := &ProviderManager{
		providers:   make([]Provider, 0),
		lastErrors:  make(map[string]error),
		healthCache: make(map[string]bool),
	}

	// 注册一个 Provider
	p := NewLocalProvider()
	manager.RegisterProvider(p)

	assert.Equal(t, 1, len(manager.providers))
	assert.Equal(t, p.Name(), manager.providers[0].Name())
}

func TestProviderManager_ProviderPriority(t *testing.T) {
	manager := NewProviderManager()

	// 验证 Provider 按优先级排序
	assert.Equal(t, "Binance", manager.providers[0].Name(), "第一个应该是Binance")
	assert.Equal(t, "Gate.io", manager.providers[1].Name(), "第二个应该是Gate.io")
	assert.Equal(t, "Frankfurter", manager.providers[2].Name(), "第三个应该是Frankfurter")
	assert.Equal(t, "Local", manager.providers[3].Name(), "第四个应该是Local")

	// 验证优先级
	assert.Equal(t, 1, manager.providers[0].Priority())
	assert.Equal(t, 2, manager.providers[1].Priority())
	assert.Equal(t, 3, manager.providers[2].Priority())
	assert.Equal(t, 999, manager.providers[3].Priority())
}

func TestProviderManager_FetchRate_LocalFallback(t *testing.T) {
	manager := NewProviderManager()
	ctx := context.Background()

	// 测试获取汇率（即使网络失败，也应该通过Local Provider返回）
	rate, warning, err := manager.FetchRate(ctx, "BTC")

	// 即使前面的 Provider 失败，Local Provider 应该成功返回
	if err != nil {
		t.Logf("Warning: %s", warning)
	}
	require.NoError(t, err)
	require.NotNil(t, rate)

	assert.Equal(t, "BTC", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0)
}

func TestProviderManager_FetchRate_CNY(t *testing.T) {
	manager := NewProviderManager()
	ctx := context.Background()

	// 测试获取 CNY 汇率（应该由Frankfurter或Local提供）
	rate, warning, err := manager.FetchRate(ctx, "CNY")

	if err != nil {
		t.Logf("Warning: %s", warning)
	}
	require.NoError(t, err)
	require.NotNil(t, rate)

	assert.Equal(t, "CNY", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0)
	// CNY应该由Frankfurter（优先级3）或Local（优先级999）提供
	assert.Contains(t, []string{"Frankfurter", "Local"}, rate.Source)
}

func TestProviderManager_FetchRate_UnsupportedSymbol(t *testing.T) {
	manager := NewProviderManager()
	ctx := context.Background()

	// 测试不支持的币种
	_, _, err := manager.FetchRate(ctx, "UNSUPPORTED")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "所有数据源获取")
}

func TestProviderManager_FetchRates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	manager := NewProviderManager()
	ctx := context.Background()

	symbols := []string{"BTC", "ETH", "USDT"}

	// 测试批量获取汇率
	rates, warnings, err := manager.FetchRates(ctx, symbols)

	require.NoError(t, err)
	require.NotNil(t, rates)

	// 验证返回的汇率
	for _, symbol := range symbols {
		rate, ok := rates[symbol]
		assert.True(t, ok, "应该包含%s的汇率", symbol)
		if ok {
			assert.Equal(t, symbol, rate.Symbol)
			assert.Greater(t, rate.PriceUSD, 0.0)
		}
	}

	// 如果有警告，记录日志
	if len(warnings) > 0 {
		t.Logf("Warnings: %v", warnings)
	}
}

func TestProviderManager_FetchRates_EmptyList(t *testing.T) {
	manager := NewProviderManager()
	ctx := context.Background()

	// 测试空币种列表
	rates, warnings, err := manager.FetchRates(ctx, []string{})

	assert.Error(t, err)
	assert.Nil(t, rates)
	assert.Contains(t, err.Error(), "所有币种汇率获取失败")
	t.Logf("Warnings: %v", warnings)
}

func TestProviderManager_RefreshHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	manager := NewProviderManager()
	ctx := context.Background()

	// 刷新健康状态
	manager.RefreshHealth(ctx)

	// 验证健康状态缓存已更新
	assert.NotEmpty(t, manager.healthCache)

	// 至少Local Provider应该是健康的
	localHealth, ok := manager.healthCache["Local"]
	assert.True(t, ok, "应该包含Local Provider的健康状态")
	assert.True(t, localHealth, "Local Provider应该是健康的")
}

func TestProviderManager_GetHealthStatus(t *testing.T) {
	manager := NewProviderManager()
	ctx := context.Background()

	// 先刷新健康状态
	manager.RefreshHealth(ctx)

	// 获取健康状态
	healthStatus := manager.GetHealthStatus()

	assert.NotEmpty(t, healthStatus)
	assert.Contains(t, healthStatus, "Local", "应该包含Local Provider")
	assert.True(t, healthStatus["Local"], "Local Provider应该是健康的")
}

func TestProviderManager_GetLastErrors(t *testing.T) {
	manager := NewProviderManager()

	// 初始应该没有错误
	errors := manager.GetLastErrors()
	assert.Empty(t, errors)
}

func TestProviderManager_AutoDegradation(t *testing.T) {
	// 这个测试验证当高优先级 Provider 失败时，自动降级到低优先级 Provider
	manager := NewProviderManager()
	ctx := context.Background()

	// 获取支持的币种（应该能从某个Provider成功获取）
	rate, warning, err := manager.FetchRate(ctx, "BTC")

	// 即使高优先级的Provider失败，也应该能从低优先级获取
	require.NoError(t, err)
	require.NotNil(t, rate)

	if warning != "" {
		t.Logf("Degradation occurred: %s", warning)
		// 验证警告信息表明发生了降级
		assert.Contains(t, warning, "数据源失败")
	}

	assert.Equal(t, "BTC", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0)
}

func TestProviderManager_ConcurrentFetchRates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	manager := NewProviderManager()
	ctx := context.Background()

	symbols := []string{"BTC", "ETH", "BNB", "SOL", "USDT"}

	// 测试并发获取（内部使用了并发）
	rates, warnings, err := manager.FetchRates(ctx, symbols)

	require.NoError(t, err)
	require.NotNil(t, rates)

	// 验证至少获取了一部分汇率
	assert.Greater(t, len(rates), 0, "应该至少获取到一些汇率")

	if len(warnings) > 0 {
		t.Logf("Warnings: %v", warnings)
	}
}
