// Package middleware 提供HTTP中间件
package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Logger 日志中间件
// 记录每个HTTP请求的详细信息：方法、路径、耗时、状态码
func Logger(r *ghttp.Request) {
	ctx := r.Context()
	startTime := time.Now()

	// 请求开始日志
	g.Log().Info(ctx, "请求开始",
		"method", r.Method,
		"path", r.URL.Path,
		"clientIP", r.GetClientIp(),
		"userAgent", r.Header.Get("User-Agent"),
	)

	// 继续处理请求
	r.Middleware.Next()

	// 计算耗时
	elapsed := time.Since(startTime)

	// 请求结束日志
	g.Log().Info(ctx, "请求完成",
		"method", r.Method,
		"path", r.URL.Path,
		"status", r.Response.Status,
		"elapsed", elapsed.Milliseconds(),
		"elapsedStr", elapsed.String(),
	)
}
