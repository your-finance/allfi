package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

// Context 上下文中间件
// 为每个请求生成唯一的TraceID，用于日志追踪
func Context(r *ghttp.Request) {
	// 生成唯一的 TraceID
	traceID := guid.S()

	// 将 TraceID 存入 Context
	ctx := r.Context()
	r.SetCtx(ctx)

	// 设置响应头
	r.Response.Header().Set("X-Trace-ID", traceID)

	// 继续处理请求
	r.Middleware.Next()
}
