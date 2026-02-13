// Package controller 基准对比模块路由注册
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	benchmarkApi "your-finance/allfi/api/v1/benchmark"
	"your-finance/allfi/internal/app/benchmark/service"
)

// BenchmarkController 基准对比控制器
type BenchmarkController struct{}

// Get 获取基准对比数据
//
// 对应路由: GET /benchmark
// 查询参数: range — 时间范围（7d/30d/90d/1y），默认 30d
func (c *BenchmarkController) Get(ctx context.Context, req *benchmarkApi.GetReq) (res *benchmarkApi.GetRes, err error) {
	// 调用 Service 层
	result, err := service.Benchmark().GetBenchmark(ctx, req.Range)
	if err != nil {
		return nil, gerror.Wrap(err, "获取基准对比数据失败")
	}

	// 将业务 DTO 转换为 API 响应
	// 构建时间序列（Portfolio + 各基准）
	var series []benchmarkApi.BenchmarkSeries

	// 用户组合序列
	series = append(series, benchmarkApi.BenchmarkSeries{
		Name:   "Portfolio",
		Points: []benchmarkApi.DataPoint{}, // 时间序列数据待扩展
		Return: result.UserReturn,
	})

	// 各基准指数序列
	for _, bm := range result.Benchmarks {
		series = append(series, benchmarkApi.BenchmarkSeries{
			Name:   bm.Name,
			Points: []benchmarkApi.DataPoint{},
			Return: bm.Return,
		})
	}

	res = &benchmarkApi.GetRes{
		Series:    series,
		Range:     result.Period,
		StartDate: result.StartDate,
		EndDate:   result.EndDate,
	}

	return res, nil
}

// Register 注册基准对比路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&BenchmarkController{})
}
