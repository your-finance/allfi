// =================================================================================
// 基准对比服务接口定义
// 提供用户投资组合收益率与基准指数（BTC/ETH/S&P500）对比能力
// =================================================================================

package service

import (
	"context"

	"your-finance/allfi/internal/app/benchmark/model"
)

// IBenchmark 基准对比服务接口
type IBenchmark interface {
	// GetBenchmark 获取用户收益率与基准指数的对比
	// period: 时间范围（7d/30d/90d/1y）
	// 返回基准对比结果，包含用户收益率、各基准指数收益率和超额收益
	GetBenchmark(ctx context.Context, period string) (*model.BenchmarkResult, error)
}

// localBenchmark 基准对比服务实例（延迟注入）
var localBenchmark IBenchmark

// Benchmark 获取基准对比服务实例
// 如果服务未注册，会触发 panic
func Benchmark() IBenchmark {
	if localBenchmark == nil {
		panic("IBenchmark 服务未注册，请检查 logic/benchmark 包的 init 函数")
	}
	return localBenchmark
}

// RegisterBenchmark 注册基准对比服务实现
// 由 logic 层在 init 函数中调用
func RegisterBenchmark(i IBenchmark) {
	localBenchmark = i
}
