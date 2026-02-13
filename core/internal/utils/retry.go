// Package utils 重试工具
// 提供带指数退避的重试机制，用于外部 API 调用
package utils

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	MaxDelay      time.Duration // 最大延迟
	Multiplier    float64       // 延迟倍数
	RandomFactor  float64       // 随机因子（0-1之间）
	RetryOn       func(error) bool // 判断是否重试的函数
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		RandomFactor: 0.5,
		RetryOn: func(err error) bool {
			return err != nil // 默认所有错误都重试
		},
	}
}

// Retry 执行带重试的操作
// operation: 需要执行的操作，返回错误时会重试
// config: 重试配置，nil 使用默认配置
func Retry(ctx context.Context, operation func() error, config *RetryConfig) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 执行操作
		err := operation()
		if err == nil {
			return nil // 成功
		}

		lastErr = err

		// 判断是否需要重试
		if config.RetryOn != nil && !config.RetryOn(err) {
			return err // 不重试
		}

		// 如果已达最大重试次数，返回错误
		if attempt >= config.MaxRetries {
			break
		}

		// 计算下次延迟（带抖动）
		jitter := time.Duration(float64(delay) * config.RandomFactor * (rand.Float64()*2 - 1))
		sleepDuration := delay + jitter

		// 限制最大延迟
		if sleepDuration > config.MaxDelay {
			sleepDuration = config.MaxDelay
		}

		// 等待
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepDuration):
		}

		// 更新延迟（指数增长）
		delay = time.Duration(float64(delay) * config.Multiplier)
	}

	return lastErr
}

// RetryWithResult 执行带重试的操作（带返回值）
func RetryWithResult[T any](ctx context.Context, operation func() (T, error), config *RetryConfig) (T, error) {
	var result T
	var lastErr error

	if config == nil {
		config = DefaultRetryConfig()
	}

	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// 执行操作
		res, err := operation()
		if err == nil {
			return res, nil // 成功
		}

		lastErr = err

		// 判断是否需要重试
		if config.RetryOn != nil && !config.RetryOn(err) {
			return result, err // 不重试
		}

		// 如果已达最大重试次数，返回错误
		if attempt >= config.MaxRetries {
			break
		}

		// 计算下次延迟（带抖动）
		jitter := time.Duration(float64(delay) * config.RandomFactor * (rand.Float64()*2 - 1))
		sleepDuration := delay + jitter

		// 限制最大延迟
		if sleepDuration > config.MaxDelay {
			sleepDuration = config.MaxDelay
		}

		// 等待
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(sleepDuration):
		}

		// 更新延迟（指数增长）
		delay = time.Duration(float64(delay) * config.Multiplier)
	}

	return result, lastErr
}

// 常用的重试判断函数

// RetryOnNetworkError 仅在网络错误时重试
func RetryOnNetworkError(err error) bool {
	if err == nil {
		return false
	}
	// 这里可以根据具体的错误类型判断
	// 简单实现：检查错误消息中是否包含网络相关关键字
	errMsg := err.Error()
	networkKeywords := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"no such host",
		"network unreachable",
		"temporary failure",
		"EOF",
	}
	for _, keyword := range networkKeywords {
		if containsIgnoreCase(errMsg, keyword) {
			return true
		}
	}
	return false
}

// RetryOnRateLimit 在遇到限流时重试
func RetryOnRateLimit(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	rateLimitKeywords := []string{
		"rate limit",
		"too many requests",
		"429",
		"throttle",
	}
	for _, keyword := range rateLimitKeywords {
		if containsIgnoreCase(errMsg, keyword) {
			return true
		}
	}
	return false
}

// containsIgnoreCase 忽略大小写检查字符串包含
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 containsLower(toLower(s), toLower(substr)))
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func containsLower(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ExponentialBackoff 计算指数退避时间
func ExponentialBackoff(attempt int, baseDelay time.Duration, maxDelay time.Duration) time.Duration {
	delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
	if delay > maxDelay {
		delay = maxDelay
	}
	// 添加随机抖动
	jitter := time.Duration(rand.Float64() * float64(delay) * 0.5)
	return delay + jitter
}

// ErrMaxRetriesExceeded 最大重试次数错误
var ErrMaxRetriesExceeded = errors.New("已达到最大重试次数")
