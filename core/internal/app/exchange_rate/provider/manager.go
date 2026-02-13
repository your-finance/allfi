// Package provider Provider 管理器
package provider

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// ProviderManager Provider 管理器
type ProviderManager struct {
	providers   []Provider
	mu          sync.RWMutex
	lastErrors  map[string]error // 记录每个提供者的最后错误
	healthCache map[string]bool  // 健康状态缓存
}

var (
	once            sync.Once
	managerInstance *ProviderManager
)

// GetProviderManager 获取 Provider 管理器单例
func GetProviderManager() *ProviderManager {
	once.Do(func() {
		managerInstance = NewProviderManager()
	})
	return managerInstance
}

// NewProviderManager 创建 Provider 管理器
func NewProviderManager() *ProviderManager {
	pm := &ProviderManager{
		providers:   make([]Provider, 0),
		lastErrors:  make(map[string]error),
		healthCache: make(map[string]bool),
	}

	// 注册所有 Provider（按优先级排序）
	pm.RegisterProvider(NewBinanceProvider())      // Priority: 1（加密货币）
	pm.RegisterProvider(NewGateioProvider())       // Priority: 2（Binance 降级）
	pm.RegisterProvider(NewFrankfurterProvider())  // Priority: 3（法币 CNY）
	pm.RegisterProvider(NewLocalProvider())        // Priority: 999（兜底）

	return pm
}

// RegisterProvider 注册 Provider
func (m *ProviderManager) RegisterProvider(p Provider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers = append(m.providers, p)

	// 按优先级排序
	sort.Slice(m.providers, func(i, j int) bool {
		return m.providers[i].Priority() < m.providers[j].Priority()
	})

	g.Log().Info(context.Background(), "注册 Provider", g.Map{
		"provider": p.Name(),
		"priority": p.Priority(),
	})
}

// FetchRate 获取汇率（自动降级）
func (m *ProviderManager) FetchRate(ctx context.Context, symbol string) (*RateInfo, string, error) {
	m.mu.RLock()
	providers := m.providers
	m.mu.RUnlock()

	var lastErr error
	var warnings []string

	// 按优先级依次尝试
	for _, provider := range providers {
		// 检查是否支持该币种
		if !provider.SupportsSymbol(symbol) {
			continue
		}

		// 检查健康状态
		if !m.isHealthy(ctx, provider) {
			warning := fmt.Sprintf("%s 不健康，跳过", provider.Name())
			warnings = append(warnings, warning)
			g.Log().Warning(ctx, warning, g.Map{
				"provider": provider.Name(),
				"symbol":   symbol,
			})
			continue
		}

		// 尝试获取汇率（记录响应时间）
		timer := time.Now()
		rate, err := provider.FetchRate(ctx, symbol)
		elapsed := time.Since(timer).Seconds()
		providerResponseTime.WithLabelValues(provider.Name()).Observe(elapsed)

		if err != nil {
			lastErr = err
			warning := fmt.Sprintf("%s 获取失败: %v", provider.Name(), err)
			warnings = append(warnings, warning)

			// 记录错误
			m.mu.Lock()
			m.lastErrors[provider.Name()] = err
			m.mu.Unlock()

			// 记录失败调用
			providerCalls.WithLabelValues(provider.Name(), "failure").Inc()

			g.Log().Warning(ctx, "Provider 获取汇率失败", g.Map{
				"provider": provider.Name(),
				"symbol":   symbol,
				"error":    err.Error(),
			})
			continue
		}

		// 成功获取
		providerCalls.WithLabelValues(provider.Name(), "success").Inc()

		g.Log().Info(ctx, "Provider 获取汇率成功", g.Map{
			"provider": provider.Name(),
			"symbol":   symbol,
			"price":    rate.PriceUSD,
		})

		// 清除错误记录
		m.mu.Lock()
		delete(m.lastErrors, provider.Name())
		m.mu.Unlock()

		// 如果有警告，附加到返回值
		warningMsg := ""
		if len(warnings) > 0 {
			warningMsg = fmt.Sprintf("部分数据源失败: %v", warnings)
		}

		return rate, warningMsg, nil
	}

	// 所有 Provider 都失败
	if lastErr != nil {
		return nil, "", fmt.Errorf("所有数据源获取 %s 汇率失败: %w", symbol, lastErr)
	}
	return nil, "", fmt.Errorf("所有数据源获取 %s 汇率失败", symbol)
}

// FetchRates 批量获取汇率
func (m *ProviderManager) FetchRates(ctx context.Context, symbols []string) (map[string]*RateInfo, map[string]string, error) {
	rates := make(map[string]*RateInfo)
	warnings := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并发获取（限制并发数）
	semaphore := make(chan struct{}, 10) // 最多10个并发

	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			rate, warning, err := m.FetchRate(ctx, sym)
			if err != nil {
				g.Log().Warning(ctx, "批量获取汇率失败", g.Map{
					"symbol": sym,
					"error":  err.Error(),
				})
				return
			}

			mu.Lock()
			rates[sym] = rate
			if warning != "" {
				warnings[sym] = warning
			}
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()

	if len(rates) == 0 {
		return nil, nil, fmt.Errorf("所有币种汇率获取失败")
	}

	return rates, warnings, nil
}

// isHealthy 检查 Provider 健康状态（带缓存）
func (m *ProviderManager) isHealthy(ctx context.Context, provider Provider) bool {
	// 检查缓存
	m.mu.RLock()
	cached, ok := m.healthCache[provider.Name()]
	m.mu.RUnlock()

	if ok {
		return cached
	}

	// 执行健康检查
	healthy := provider.IsHealthy(ctx)

	// 更新缓存
	m.mu.Lock()
	m.healthCache[provider.Name()] = healthy
	m.mu.Unlock()

	return healthy
}

// RefreshHealth 刷新所有 Provider 的健康状态
func (m *ProviderManager) RefreshHealth(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.healthCache = make(map[string]bool)

	for _, provider := range m.providers {
		healthy := provider.IsHealthy(ctx)
		m.healthCache[provider.Name()] = healthy

		// 更新 Prometheus 健康状态指标
		if healthy {
			providerHealth.WithLabelValues(provider.Name()).Set(1)
		} else {
			providerHealth.WithLabelValues(provider.Name()).Set(0)
		}

		g.Log().Info(ctx, "Provider 健康检查", g.Map{
			"provider": provider.Name(),
			"healthy":  healthy,
		})
	}
}

// GetLastErrors 获取所有 Provider 的最后错误
func (m *ProviderManager) GetLastErrors() map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	errors := make(map[string]error)
	for k, v := range m.lastErrors {
		errors[k] = v
	}
	return errors
}

// GetHealthStatus 获取所有 Provider 的健康状态
func (m *ProviderManager) GetHealthStatus() map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]bool)
	for k, v := range m.healthCache {
		status[k] = v
	}
	return status
}
