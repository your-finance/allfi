// Package utils 限流工具
// 提供 Token Bucket 限流器，用于控制外部 API 调用频率
package utils

import (
	"context"
	"sync"
	"time"
)

// RateLimiter Token Bucket 限流器
type RateLimiter struct {
	rate       float64   // 每秒产生的令牌数
	bucketSize float64   // 桶的最大容量
	tokens     float64   // 当前令牌数
	lastUpdate time.Time // 上次更新时间
	mutex      sync.Mutex
}

// NewRateLimiter 创建新的限流器
// rate: 每秒允许的请求数
// bucketSize: 突发请求的最大数量
func NewRateLimiter(rate float64, bucketSize float64) *RateLimiter {
	if bucketSize < 1 {
		bucketSize = 1
	}
	return &RateLimiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastUpdate: time.Now(),
	}
}

// Allow 检查是否允许请求（非阻塞）
// 返回 true 表示允许，false 表示被限流
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN 检查是否允许 n 个请求（非阻塞）
func (rl *RateLimiter) AllowN(n float64) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.refill()

	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}
	return false
}

// Wait 等待直到允许请求（阻塞）
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.WaitN(ctx, 1)
}

// WaitN 等待直到允许 n 个请求（阻塞）
func (rl *RateLimiter) WaitN(ctx context.Context, n float64) error {
	for {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if rl.AllowN(n) {
			return nil
		}

		// 计算需要等待的时间
		waitTime := rl.waitTime(n)
		if waitTime <= 0 {
			waitTime = time.Millisecond * 10
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
}

// refill 补充令牌（内部方法，需要在持有锁的情况下调用）
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.lastUpdate = now

	// 根据时间补充令牌
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.bucketSize {
		rl.tokens = rl.bucketSize
	}
}

// waitTime 计算需要等待的时间
func (rl *RateLimiter) waitTime(n float64) time.Duration {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.refill()

	if rl.tokens >= n {
		return 0
	}

	need := n - rl.tokens
	return time.Duration(need / rl.rate * float64(time.Second))
}

// Reset 重置限流器
func (rl *RateLimiter) Reset() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.tokens = rl.bucketSize
	rl.lastUpdate = time.Now()
}

// Tokens 获取当前令牌数
func (rl *RateLimiter) Tokens() float64 {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.refill()
	return rl.tokens
}

// RateLimiterManager 限流器管理器
// 为不同的 API 维护独立的限流器
type RateLimiterManager struct {
	limiters map[string]*RateLimiter
	mutex    sync.RWMutex
	configs  map[string]RateLimitConfig
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Rate       float64 // 每秒请求数
	BucketSize float64 // 突发请求数
}

// 预定义的 API 限流配置
var DefaultAPIRateLimits = map[string]RateLimitConfig{
	"binance":    {Rate: 10, BucketSize: 20},   // Binance: 1200/min = 20/s
	"okx":        {Rate: 10, BucketSize: 10},   // OKX: 较保守
	"coinbase":   {Rate: 10, BucketSize: 10},   // Coinbase
	"etherscan":   {Rate: 5, BucketSize: 5},     // Etherscan 免费版: 5/s
	"bscscan":     {Rate: 5, BucketSize: 5},     // BscScan 免费版: 5/s
	"arbiscan":    {Rate: 5, BucketSize: 5},     // Arbiscan 免费版: 5/s
	"optimism":    {Rate: 5, BucketSize: 5},     // Optimism Explorer 免费版: 5/s
	"polygonscan": {Rate: 5, BucketSize: 5},     // Polygonscan 免费版: 5/s
	"basescan":    {Rate: 5, BucketSize: 5},     // Basescan 免费版: 5/s
	"coingecko":   {Rate: 10, BucketSize: 50},   // CoinGecko 免费版: 10-50/min
	"yahoo":      {Rate: 5, BucketSize: 10},    // Yahoo Finance
}

// NewRateLimiterManager 创建限流器管理器
func NewRateLimiterManager() *RateLimiterManager {
	return &RateLimiterManager{
		limiters: make(map[string]*RateLimiter),
		configs:  DefaultAPIRateLimits,
	}
}

// GetLimiter 获取指定 API 的限流器
func (m *RateLimiterManager) GetLimiter(name string) *RateLimiter {
	m.mutex.RLock()
	limiter, exists := m.limiters[name]
	m.mutex.RUnlock()

	if exists {
		return limiter
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 双重检查
	if limiter, exists = m.limiters[name]; exists {
		return limiter
	}

	// 获取配置
	config, ok := m.configs[name]
	if !ok {
		// 使用默认配置
		config = RateLimitConfig{Rate: 5, BucketSize: 10}
	}

	// 创建限流器
	limiter = NewRateLimiter(config.Rate, config.BucketSize)
	m.limiters[name] = limiter
	return limiter
}

// Wait 等待指定 API 的限流
func (m *RateLimiterManager) Wait(ctx context.Context, name string) error {
	return m.GetLimiter(name).Wait(ctx)
}

// Allow 检查指定 API 是否允许请求
func (m *RateLimiterManager) Allow(name string) bool {
	return m.GetLimiter(name).Allow()
}

// 全局限流器管理器
var globalLimiterManager *RateLimiterManager
var limiterOnce sync.Once

// GetGlobalLimiterManager 获取全局限流器管理器
func GetGlobalLimiterManager() *RateLimiterManager {
	limiterOnce.Do(func() {
		globalLimiterManager = NewRateLimiterManager()
	})
	return globalLimiterManager
}

// WaitForAPI 等待 API 限流（使用全局管理器）
func WaitForAPI(ctx context.Context, apiName string) error {
	return GetGlobalLimiterManager().Wait(ctx, apiName)
}

// AllowAPI 检查 API 是否允许请求（使用全局管理器）
func AllowAPI(apiName string) bool {
	return GetGlobalLimiterManager().Allow(apiName)
}
